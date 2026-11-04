<template>
    <fieldset
        class="workspace-fieldset media-studio"
        :class="{ 'media-standalone': standalone }"
    >
        <legend>
            {{ standalone ? 'Media library & studio' : 'Product imagery' }}
        </legend>
        <nav
            v-if="standalone && ready"
            class="media-view-nav"
            aria-label="Media workspace views"
        >
            <template v-if="canCreate">
                <button
                    v-for="view in mediaViews"
                    :key="view.key"
                    type="button"
                    class="button"
                    :class="
                        activeView === view.key
                            ? 'button-navy'
                            : 'button-secondary'
                    "
                    :aria-pressed="activeView === view.key"
                    @click="activeView = view.key"
                >
                    {{ view.label }}
                </button>
            </template>
            <button
                v-if="canManagePolicy"
                type="button"
                class="button"
                :class="
                    activeView === 'policy' ? 'button-navy' : 'button-secondary'
                "
                :aria-pressed="activeView === 'policy'"
                @click="activeView = 'policy'"
            >
                Image policy
            </button>
        </nav>
        <p v-if="!ready && loading" role="status" class="form-note">
            Opening your media workspace…
        </p>
        <button
            v-if="!ready && !loading"
            type="button"
            class="button button-secondary"
            @click="loadLibrary(1, true)"
        >
            Retry loading media
        </button>
        <p class="form-note">
            Show the real product when possible. A mockup is an illustration,
            not evidence of its appearance or specifications.
        </p>
        <p v-if="error" role="alert" class="form-error">{{ error }}</p>
        <p v-if="notice" role="status" class="form-note">{{ notice }}</p>
        <div
            ref="previewPanel"
            tabindex="-1"
            class="media-cover"
            v-if="modelValue && (!standalone || activeView === 'library')"
        >
            <img :src="modelValue" alt="Current product cover preview" />
            <button
                type="button"
                class="button button-secondary"
                @click="emit('update:modelValue', '')"
            >
                {{ standalone ? 'Close preview' : 'Remove cover' }}
            </button>
        </div>
        <div class="soko-field" v-if="policy.allow_remote && !standalone">
            <label :for="`${id}-url`">Image URL or local asset path</label>
            <input
                :id="`${id}-url`"
                name="image_url"
                :value="modelValue"
                @input="
                    emit(
                        'update:modelValue',
                        ($event.target as HTMLInputElement).value
                    )
                "
                placeholder="https://…"
                maxlength="2048"
            />
            <p class="form-note">
                Remote image dimensions cannot be verified here. Upload for
                quality validation.
            </p>
        </div>
        <div
            v-if="canCreate"
            v-show="!standalone || activeView === 'upload'"
            class="soko-field media-upload-panel"
        >
            <label :for="`${id}-upload`">Upload a product photo</label>
            <input
                :id="`${id}-upload`"
                type="file"
                accept="image/png,image/jpeg"
                :disabled="busy || !ready"
                @change="uploadPhoto"
            />
            <p class="form-note">
                PNG or JPEG · up to
                {{ Math.round((policy.max_bytes / 1048576) * 10) / 10 }} MB ·
                minimum {{ policy.min_width }} × {{ policy.min_height }} px ·
                maximum {{ policy.max_dimension }} px per side. Images are
                re-encoded to remove metadata.
            </p>
        </div>
        <details
            ref="mockupPanel"
            :open="standalone"
            v-show="!standalone || activeView === 'mockup'"
            v-if="canCreate && policy.allow_mockups"
            class="media-details"
        >
            <summary>Create an illustrative product mockup</summary>
            <div class="form-row">
                <div class="soko-field">
                    <label :for="`${id}-template`">Product canvas</label
                    ><select :id="`${id}-template`" v-model="mockup.template">
                        <option
                            v-for="option in templates"
                            :key="option.value"
                            :value="option.value"
                        >
                            {{ option.label }}
                        </option>
                    </select>
                </div>
                <div class="soko-field">
                    <label :for="`${id}-brand`">Brand name</label
                    ><input
                        :id="`${id}-brand`"
                        v-model="mockup.brand"
                        maxlength="40"
                    />
                </div>
            </div>
            <div class="soko-field">
                <label :for="`${id}-title`">Short title on mockup</label
                ><input
                    :id="`${id}-title`"
                    v-model="mockup.title"
                    maxlength="80"
                />
            </div>
            <div class="media-colors">
                <div
                    class="soko-field"
                    v-for="color in colors"
                    :key="color.key"
                >
                    <label :for="`${id}-${color.key}`">{{ color.label }}</label
                    ><input
                        :id="`${id}-${color.key}`"
                        type="color"
                        v-model="mockup[color.key]"
                    />
                </div>
            </div>
            <div class="soko-field">
                <label :for="`${id}-logo`"
                    >Your logo (optional PNG or JPEG)</label
                ><input
                    :id="`${id}-logo`"
                    type="file"
                    accept="image/png,image/jpeg"
                    :disabled="busy"
                    @change="chooseLogo"
                />
                <p class="form-note">
                    Minimum 64 × 64 px. Only upload a logo you are authorized to
                    use.
                </p>
            </div>
            <button
                type="button"
                class="button button-navy"
                :disabled="busy || !ready"
                @click="createMockup"
            >
                {{
                    busy
                        ? 'Saving image…'
                        : standalone
                          ? 'Create mockup'
                          : 'Create and use mockup'
                }}
            </button>
        </details>
        <details
            v-if="canCreate"
            v-show="!standalone || activeView === 'library'"
            class="media-details"
            :open="standalone"
        >
            <summary>Your image library ({{ total }})</summary>
            <div class="form-row media-filters">
                <div class="soko-field">
                    <label :for="`${id}-search`"
                        >Search mockup text or image ID</label
                    ><input
                        :id="`${id}-search`"
                        v-model="search"
                        maxlength="100"
                        type="search"
                        @keydown.enter.prevent="loadLibrary(1)"
                    />
                </div>
                <div class="soko-field">
                    <label :for="`${id}-kind`">Image type</label
                    ><select
                        :id="`${id}-kind`"
                        v-model="kind"
                        @change="loadLibrary(1)"
                    >
                        <option value="">All images</option>
                        <option value="upload">Photos</option>
                        <option value="mockup">Mockups</option>
                    </select>
                </div>
                <button
                    type="button"
                    class="button button-secondary"
                    :disabled="loading"
                    @click="loadLibrary(1)"
                >
                    Search library
                </button>
                <button
                    v-if="search || kind"
                    type="button"
                    class="button button-secondary"
                    :disabled="loading"
                    @click="clearFilters"
                >
                    Clear filters
                </button>
            </div>
            <p v-if="loading" role="status" class="form-note">
                Loading images…
            </p>
            <div
                v-if="!loading && ready && !assets.length"
                class="media-empty"
                role="status"
            >
                <strong>{{
                    search || kind
                        ? 'No images match your filters'
                        : 'Your image library starts here'
                }}</strong>
                <p>
                    {{
                        search || kind
                            ? 'Clear the filters or try another mockup title, brand or image ID.'
                            : 'Upload a product photo or create a branded mockup. New images remain private until used in a saved listing.'
                    }}
                </p>
                <button
                    v-if="standalone && !search && !kind"
                    type="button"
                    class="button button-navy"
                    @click="activeView = 'upload'"
                >
                    Upload your first photo
                </button>
            </div>
            <div class="media-library" :aria-busy="loading">
                <article v-for="asset in assets" :key="asset.id">
                    <img
                        :src="asset.url"
                        :alt="
                            asset.kind === 'mockup'
                                ? 'Illustrative product mockup'
                                : 'Uploaded product photo'
                        "
                        loading="lazy"
                    /><small
                        >{{
                            asset.kind === 'mockup'
                                ? 'Illustrative mockup'
                                : 'Product photo'
                        }}
                        · {{ asset.width }} × {{ asset.height }}</small
                    >
                    <details class="media-asset-meta">
                        <summary>Image details</summary>
                        <small class="media-asset-id">ID: {{ asset.id }}</small
                        ><small v-if="asset.created_at"
                            >Added {{ formatDate(asset.created_at) }}</small
                        >
                    </details>
                    <div>
                        <button
                            v-if="
                                asset.kind === 'mockup' && policy.allow_mockups
                            "
                            type="button"
                            class="button button-secondary"
                            :disabled="busy"
                            @click="customizeCopy(asset)"
                        >
                            Customize a copy
                        </button>
                        <button
                            type="button"
                            class="button button-secondary"
                            @click="previewAsset(asset)"
                        >
                            {{
                                standalone ? 'Preview image' : 'Use cover'
                            }}</button
                        ><button
                            v-if="!standalone"
                            type="button"
                            class="button button-secondary"
                            :disabled="
                                gallery.length >= 10 ||
                                gallery.includes(asset.url)
                            "
                            @click="addGallery(asset.url)"
                        >
                            Add to gallery
                        </button>
                    </div>
                </article>
            </div>
            <nav
                v-if="total > 0"
                class="media-pagination"
                aria-label="Image library pages"
            >
                <button
                    type="button"
                    class="button button-secondary"
                    :disabled="loading || page <= 1"
                    @click="loadLibrary(page - 1)"
                >
                    Previous
                </button>
                <span role="status"
                    >Page {{ page }} of
                    {{ Math.max(1, Math.ceil(total / perPage)) }} ·
                    {{ total }} images</span
                >
                <button
                    type="button"
                    class="button button-secondary"
                    :disabled="loading || page * perPage >= total"
                    @click="loadLibrary(page + 1)"
                >
                    Next
                </button>
            </nav>
        </details>
        <div v-if="gallery.length && !standalone" class="media-gallery">
            <strong>Gallery · {{ gallery.length }}/10</strong>
            <div v-for="(url, index) in gallery" :key="`${index}-${url}`">
                <img :src="url" :alt="`Gallery image ${index + 1}`" /><button
                    type="button"
                    class="button button-secondary"
                    @click="
                        emit(
                            'update:gallery',
                            gallery.filter((_, i) => i !== index)
                        )
                    "
                >
                    Remove image {{ index + 1 }}
                </button>
            </div>
        </div>
        <details
            v-if="canManagePolicy"
            class="media-details"
            :open="standalone"
            v-show="!standalone || activeView === 'policy'"
        >
            <summary>Administrator image policy</summary>
            <p class="form-note">
                Applies to new uploads and new or edited image links. Existing
                images are retained.
            </p>
            <div class="form-row">
                <div class="soko-field">
                    <label :for="`${id}-max`">Maximum upload MB</label
                    ><input
                        :id="`${id}-max`"
                        type="number"
                        min="0.1"
                        max="10"
                        step="0.1"
                        :value="policy.max_bytes / 1048576"
                        @input="
                            policy.max_bytes = Math.round(
                                Number(
                                    ($event.target as HTMLInputElement).value
                                ) * 1048576
                            )
                        "
                    />
                </div>
                <div class="soko-field">
                    <label :for="`${id}-dimension`"
                        >Maximum dimension (px)</label
                    ><input
                        :id="`${id}-dimension`"
                        type="number"
                        v-model.number="policy.max_dimension"
                        min="100"
                        max="6000"
                    />
                </div>
            </div>
            <div class="form-row">
                <div class="soko-field">
                    <label :for="`${id}-width`">Minimum width (px)</label
                    ><input
                        :id="`${id}-width`"
                        type="number"
                        v-model.number="policy.min_width"
                        min="100"
                        max="2048"
                    />
                </div>
                <div class="soko-field">
                    <label :for="`${id}-height`">Minimum height (px)</label
                    ><input
                        :id="`${id}-height`"
                        type="number"
                        v-model.number="policy.min_height"
                        min="100"
                        max="2048"
                    />
                </div>
            </div>
            <label class="media-toggle"
                ><input type="checkbox" v-model="policy.allow_remote" />Allow
                HTTPS image links</label
            ><label class="media-toggle"
                ><input type="checkbox" v-model="policy.allow_mockups" />Allow
                illustrative mockups</label
            >
            <button
                type="button"
                class="button button-navy"
                :disabled="busy"
                @click="savePolicy"
            >
                Save image policy
            </button>
        </details>
    </fieldset>
</template>
<script setup lang="ts">
import { computed, nextTick, onMounted, reactive, ref, useId } from 'vue'
const props = defineProps<{
    modelValue: string
    gallery: string[]
    title?: string
    standalone?: boolean
}>()
const emit = defineEmits<{
    'update:modelValue': [value: string]
    'update:gallery': [value: string[]]
}>()
const id = useId()
const mockupPanel = ref<HTMLDetailsElement | null>(null)
const previewPanel = ref<HTMLDivElement | null>(null)
const activeView = ref('library')
const mediaViews = computed(() => [
    { key: 'library', label: 'Image library' },
    { key: 'upload', label: 'Upload photos' },
    ...(policy.allow_mockups
        ? [{ key: 'mockup', label: 'Create mockup' }]
        : []),
])
const formatDate = (value: string) =>
    new Date(value).toLocaleDateString(undefined, {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
    })
async function previewAsset(asset: Asset) {
    emit('update:modelValue', asset.url)
    if (props.standalone) {
        activeView.value = 'library'
        await nextTick()
        previewPanel.value?.focus()
    }
}
async function clearFilters() {
    search.value = ''
    kind.value = ''
    await loadLibrary(1)
}
interface Asset {
    id: string
    url: string
    kind: string
    width: number
    height: number
    created_at?: string
}
const page = ref(1),
    total = ref(0),
    perPage = ref(24),
    search = ref(''),
    kind = ref(''),
    loading = ref(false),
    canCreate = ref(false)
let requestVersion = 0
const assets = ref<Asset[]>([]),
    busy = ref(false),
    ready = ref(false),
    error = ref(''),
    notice = ref(''),
    canManagePolicy = ref(false)
const policy = reactive({
    max_bytes: 5242880,
    min_width: 600,
    min_height: 600,
    max_dimension: 4096,
    allow_remote: true,
    allow_mockups: true,
})
const mockup = reactive({
    template: 'software',
    title: props.title?.slice(0, 80) || '',
    brand: '',
    primary: '#14283d',
    accent: '#d5a64d',
    background: '#f4f5f7',
    logo_data: '',
})
const templates = [
    { value: 'software', label: 'Software / screen' },
    { value: 'box', label: 'Packaging / box' },
    { value: 'bottle', label: 'Bottle / container' },
    { value: 'apparel', label: 'Apparel / shirt' },
    { value: 'book', label: 'Book / report' },
    { value: 'equipment', label: 'Equipment / generic machine' },
]
const colors: { key: 'primary' | 'accent' | 'background'; label: string }[] = [
    { key: 'primary', label: 'Main color' },
    { key: 'accent', label: 'Accent color' },
    { key: 'background', label: 'Background' },
]
async function api(path: string, body?: unknown) {
    const response = await fetch(path, {
        method: body === undefined ? 'GET' : 'POST',
        credentials: 'same-origin',
        headers: {
            Accept: 'application/json',
            ...(body === undefined
                ? {}
                : { 'Content-Type': 'application/json' }),
        },
        body: body === undefined ? undefined : JSON.stringify(body),
    })
    if (response.redirected)
        throw new Error('Please sign in again before managing images.')
    const data = await response.json()
    if (!response.ok) throw new Error(data.error || 'Image request failed.')
    return data
}
async function loadLibrary(targetPage = 1, initial = false) {
    const version = ++requestVersion
    loading.value = true
    error.value = ''
    try {
        const params = new URLSearchParams({
            page: String(targetPage),
            q: search.value,
            kind: kind.value,
        })
        const data = await api(`/dashboard/media?${params}`)
        if (version !== requestVersion) return
        assets.value = data.assets || []
        page.value = data.page
        total.value = data.total
        perPage.value = data.per_page
        if (initial) Object.assign(policy, data.policy)
        canCreate.value = !!data.can_create
        canManagePolicy.value = !!data.can_manage_policy
        if (!canCreate.value && canManagePolicy.value)
            activeView.value = 'policy'
        ready.value = true
    } catch (e) {
        if (version === requestVersion) error.value = message(e)
    } finally {
        if (version === requestVersion) loading.value = false
    }
}
onMounted(() => loadLibrary(1, true))
const message = (e: unknown) =>
    e instanceof Error ? e.message : 'Could not save the image.'
function readFile(file: File): Promise<string> {
    if (!['image/png', 'image/jpeg'].includes(file.type))
        return Promise.reject(new Error('Choose a PNG or JPEG image.'))
    if (file.size > policy.max_bytes)
        return Promise.reject(
            new Error('The file exceeds the configured upload limit.')
        )
    return new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.onload = () => resolve(String(reader.result))
        reader.onerror = () => reject(new Error('The image could not be read.'))
        reader.readAsDataURL(file)
    })
}
async function work(action: () => Promise<void>) {
    if (busy.value) return
    busy.value = true
    error.value = ''
    notice.value = ''
    try {
        await action()
    } catch (e) {
        error.value = message(e)
    } finally {
        busy.value = false
    }
}
async function selectAsset(asset: Asset) {
    search.value = ''
    kind.value = ''
    await loadLibrary(1)
    await previewAsset(asset)
    notice.value = props.standalone
        ? 'Image saved privately to your library. Select it in a product listing and save the listing to publish it.'
        : 'Image saved to your library and selected as the cover. Save the listing to publish this selection.'
}
async function uploadPhoto(event: Event) {
    const input = event.target as HTMLInputElement
    const file = input.files?.[0]
    if (!file) return
    await work(async () => {
        const data = await api('/dashboard/media/upload', {
            data: await readFile(file),
        })
        await selectAsset(data.asset)
    })
    input.value = ''
}
async function chooseLogo(event: Event) {
    const input = event.target as HTMLInputElement
    const file = input.files?.[0]
    if (!file) return
    await work(async () => {
        mockup.logo_data = await readFile(file)
        notice.value = 'Logo selected. Create the mockup to apply it.'
    })
    input.value = ''
}
async function customizeCopy(asset: Asset) {
    await work(async () => {
        const data = await api(
            `/dashboard/media/${encodeURIComponent(asset.id)}/recipe`
        )
        Object.assign(mockup, data.recipe, { logo_data: '' })
        activeView.value = 'mockup'
        if (mockupPanel.value) mockupPanel.value.open = true
        await nextTick()
        document.getElementById(`${id}-title`)?.focus()
        notice.value =
            'Settings loaded for a new copy. The original image stays unchanged. Add your logo again if needed, then create the mockup.'
    })
}
async function createMockup() {
    await work(async () => {
        const data = await api('/dashboard/media/mockup', mockup)
        await selectAsset(data.asset)
    })
}
function addGallery(url: string) {
    if (props.gallery.length < 10 && !props.gallery.includes(url))
        emit('update:gallery', [...props.gallery, url])
}
async function savePolicy() {
    await work(async () => {
        const data = await api('/dashboard/media/policy', policy)
        Object.assign(policy, data.policy)
        notice.value = 'Image policy saved.'
    })
}
</script>
<style scoped>
.media-view-nav {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid var(--soko-border);
    margin-bottom: 1rem;
}
.media-standalone {
    border: 0;
    padding: 0;
}
.media-standalone > legend {
    font-size: 1rem;
    padding: 0;
    margin-bottom: 1rem;
}
.media-standalone .media-details {
    border: 0;
    margin-top: 0;
}
.media-standalone .media-library {
    grid-template-columns: repeat(auto-fill, minmax(min(100%, 230px), 1fr));
    align-items: start;
}
.media-standalone .media-library article {
    border: 1px solid var(--soko-border);
    border-radius: var(--soko-radius);
    overflow: hidden;
    padding: 1rem;
    gap: 0.8rem;
    background: var(--soko-surface);
}
.media-standalone .media-library img {
    height: auto;
    aspect-ratio: 4/3;
    border-radius: var(--soko-radius);
    border: 1px solid var(--soko-border);
}
.media-standalone .media-library .button {
    min-height: 40px;
    font-size: 0.8rem;
}
.media-asset-meta {
    font-size: 0.75rem;
    color: var(--soko-muted);
}
.media-asset-meta summary {
    cursor: pointer;
}
.media-asset-meta small {
    display: block;
    margin-top: 0.4rem;
}
.media-empty {
    padding: 2rem;
    border: 1px dashed var(--soko-border);
    border-radius: var(--soko-radius);
    text-align: center;
}
.media-empty p {
    max-width: 54ch;
    margin: 0.7rem auto 1rem;
    color: var(--soko-muted);
}
.media-standalone .media-cover {
    padding: 1rem;
    margin: 1rem 0;
    border: 1px solid var(--soko-border);
    border-radius: var(--soko-radius);
}
.media-standalone .media-cover img {
    width: min(100%, 480px);
    height: 280px;
}
.media-cover:focus {
    outline: 2px solid var(--soko-primary);
    outline-offset: 3px;
}
.media-upload-panel {
    max-width: 640px;
}
.media-upload-panel input[type='file'] {
    padding: 1rem;
    border: 1px dashed var(--soko-border);
    border-radius: var(--soko-radius);
}

.media-pagination {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    flex-wrap: wrap;
    margin-top: 1rem;
}
.media-filters {
    align-items: end;
}
.media-standalone .media-library {
    max-height: none;
    overflow: visible;
}
.media-asset-id {
    overflow-wrap: anywhere;
    color: var(--soko-muted);
}
.media-studio {
    min-width: 0;
}
.media-cover {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-wrap: wrap;
}
.media-cover img {
    width: 180px;
    height: 140px;
    object-fit: contain;
    border: 1px solid var(--soko-border);
    border-radius: var(--soko-radius, 8px);
    background: var(--soko-surface, #fff);
}
.media-details {
    border-top: 1px solid var(--soko-border);
    padding-top: 1rem;
    margin-top: 1rem;
}
.media-details summary {
    font-weight: 700;
    cursor: pointer;
    margin-bottom: 1rem;
}
.media-colors {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 1rem;
}
.media-colors input {
    width: 100%;
    min-height: 44px;
    padding: 4px;
}
.media-library {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 1rem;
    max-height: 420px;
    overflow: auto;
}
.media-library article {
    min-width: 0;
    display: grid;
    gap: 0.5rem;
}
.media-library img {
    width: 100%;
    height: 140px;
    object-fit: contain;
    background: var(--soko-surface, #fff);
}
.media-library article > div {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
}
.media-library .button {
    padding: 0.5rem;
    font-size: 0.75rem;
}
.media-gallery {
    display: grid;
    gap: 0.5rem;
}
.media-gallery > div {
    display: flex;
    align-items: center;
    gap: 1rem;
}
.media-gallery img {
    width: 64px;
    height: 64px;
    object-fit: contain;
}
.media-toggle {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    margin: 0.8rem 0;
}
@media (max-width: 480px) {
    .media-colors {
        grid-template-columns: 1fr;
    }
    .media-library {
        grid-template-columns: 1fr 1fr;
    }
}
</style>
