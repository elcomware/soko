<template>
    <div class="site-container merchandising-editor">
        <header class="workspace-heading">
            <div>
                <h1>Storefront</h1>
                <p>Arrange sections and choose featured products.</p>
            </div>
            <Link href="/dashboard" class="button button-outline"
                >Back to workspace</Link
            >
        </header>
        <AdminNavigationGuard
            :busy="form.processing"
            :dirty="draftDirty && !editingKey"
            @discard="discardDraft"
        />
        <SokoDialog
            :open="confirmDiscard"
            title="Discard storefront changes?"
            @close="confirmDiscard = false"
            ><p>Your published storefront will stay unchanged.</p>
            <div class="editor-actions">
                <button
                    type="button"
                    class="button button-outline"
                    @click="confirmDiscard = false"
                >
                    Keep editing</button
                ><button
                    type="button"
                    class="button button-navy"
                    @click="confirmDiscardDraft"
                >
                    Discard draft
                </button>
            </div></SokoDialog
        >
        <div class="editor-status">
            <span role="status">{{
                draftDirty
                    ? 'Draft changes · Not published'
                    : 'Published configuration'
            }}</span
            ><span>Version {{ form.version }}</span
            ><button
                v-if="draftDirty"
                class="text-button"
                type="button"
                :disabled="form.processing"
                @click="confirmDiscard = true"
            >
                Discard changes</button
            ><Link href="/" target="_blank" rel="noopener" class="text-button"
                >View live storefront ↗</Link
            >
        </div>
        <p v-if="flash.error" class="form-error" role="alert">
            {{ flash.error }} Reload this page before retrying if another
            administrator changed these settings.
        </p>
        <div
            v-if="Object.keys(form.errors || {}).length"
            class="form-error"
            role="alert"
        >
            <strong>Check these fields before publishing:</strong>
            <ul>
                <li v-for="(error, field) in form.errors" :key="field">
                    {{ error }}
                </li>
            </ul>
        </div>
        <form @submit.prevent="save" class="merchandising-form">
            <fieldset
                class="merchandising-form-lock"
                :disabled="form.processing"
            >
                <section
                    class="workspace-panel"
                    aria-labelledby="section-order-title"
                >
                    <h2 id="section-order-title">Section order</h2>
                    <p>
                        Reorder sections. Hidden or empty sections leave no gap.
                    </p>
                    <ol class="section-order-list">
                        <li
                            v-for="(key, index) in form.section_order"
                            :key="key"
                        >
                            <span class="section-position" aria-hidden="true">{{
                                index + 1
                            }}</span>
                            <span class="section-order-copy"
                                ><strong>{{ sectionLabels[key] }}</strong
                                ><small>{{
                                    key === 'showcase'
                                        ? 'Main product carousel and side cards stay together'
                                        : key === 'catalog'
                                          ? 'Required: products, filters and buying paths'
                                          : 'Visibility follows the section’s settings and available content'
                                }}</small></span
                            >
                            <div class="section-order-controls">
                                <button
                                    type="button"
                                    :aria-label="
                                        'Move ' + sectionLabels[key] + ' up'
                                    "
                                    :disabled="form.processing || index === 0"
                                    @click="moveSection(index, -1)"
                                >
                                    Move up
                                </button>
                                <button
                                    type="button"
                                    :aria-label="
                                        'Move ' + sectionLabels[key] + ' down'
                                    "
                                    :disabled="
                                        form.processing ||
                                        index === form.section_order.length - 1
                                    "
                                    @click="moveSection(index, 1)"
                                >
                                    Move down
                                </button>
                            </div>
                        </li>
                    </ol>
                    <p class="sr-only" role="status" aria-live="polite">
                        {{ orderAnnouncement }}
                    </p>
                    <button
                        type="button"
                        class="button button-outline"
                        :disabled="form.processing"
                        @click="resetSectionOrder"
                    >
                        Reset section order
                    </button>
                </section>
                <section class="workspace-panel">
                    <h2>Storefront announcement</h2>
                    <label
                        >Announcement text<input
                            v-model="form.announcement"
                            maxlength="160"
                    /></label>
                    <p>
                        Plain text only. Keep claims accurate; leave blank to
                        show no announcement text.
                    </p>
                </section>
                <section
                    v-for="section in sections"
                    :key="section.key"
                    class="workspace-panel placement-summary"
                >
                    <div>
                        <span
                            class="placement-state"
                            :class="{ disabled: !form[section.key].enabled }"
                            >{{
                                form[section.key].enabled
                                    ? 'Visible when eligible products exist'
                                    : 'Hidden'
                            }}</span
                        >
                        <h2>{{ section.label }}</h2>
                        <p>{{ section.help }}</p>
                    </div>
                    <dl>
                        <div>
                            <dt>Selection</dt>
                            <dd>
                                {{
                                    strategyLabels[
                                        form[section.key].strategy
                                    ] || form[section.key].strategy
                                }}
                            </dd>
                        </div>
                        <div>
                            <dt>Products</dt>
                            <dd>Up to {{ form[section.key].limit }}</dd>
                        </div>
                        <div>
                            <dt>Category</dt>
                            <dd>
                                {{
                                    form[section.key].category || 'Any category'
                                }}
                            </dd>
                        </div>
                    </dl>
                    <button
                        type="button"
                        class="button button-outline"
                        :disabled="form.processing"
                        @click="openPlacement(section.key)"
                    >
                        Edit placement
                    </button>
                    <AdminDrawer
                        :open="editingKey === section.key"
                        :title="section.label"
                        :description="section.help"
                        :busy="form.processing"
                        :dirty="placementDirty"
                        @close="discardPlacement"
                    >
                        <div class="merchandising-form placement-editor">
                            <p
                                v-if="placementError"
                                class="form-error"
                                role="alert"
                            >
                                {{ placementError }}
                            </p>
                            <label class="check"
                                ><input
                                    type="checkbox"
                                    v-model="form[section.key].enabled"
                                />Show section</label
                            >
                            <div class="editor-grid">
                                <label
                                    >Section title<input
                                        v-model="form[section.key].title"
                                        maxlength="100"
                                /></label>
                                <label
                                    >Selection rule<select
                                        v-model="form[section.key].strategy"
                                    >
                                        <option value="featured">
                                            Featured products first
                                        </option>
                                        <option value="newest">
                                            Newest products
                                        </option>
                                        <option value="promotions">
                                            Active discounted products
                                        </option>
                                        <option value="manual">
                                            Selected products in your order
                                        </option>
                                    </select></label
                                >
                                <label
                                    >Maximum products<input
                                        type="number"
                                        v-model.number="form[section.key].limit"
                                        min="1"
                                        :max="
                                            section.key === 'hero_side' ? 2 : 12
                                        "
                                        required
                                /></label>
                                <label
                                    >If there are not enough matching
                                    products<select
                                        v-model="form[section.key].fallback"
                                    >
                                        <option value="newest">
                                            Fill with newest eligible products
                                        </option>
                                        <option value="hide">
                                            Show only matches; hide if empty
                                        </option>
                                    </select></label
                                >
                                <label
                                    >Restrict to category (optional)<select
                                        v-model="form[section.key].category"
                                    >
                                        <option value="">Any category</option>
                                        <option
                                            v-for="category in categories"
                                            :key="category"
                                        >
                                            {{ category }}
                                        </option>
                                    </select></label
                                >
                            </div>
                            <fieldset
                                v-if="form[section.key].strategy === 'manual'"
                            >
                                <legend>Ordered product selection</legend>
                                <div
                                    v-for="(id, index) in form[section.key]
                                        .product_ids"
                                    :key="index"
                                    class="pin-row"
                                >
                                    <label
                                        >Position {{ index + 1
                                        }}<select
                                            v-model="
                                                form[section.key].product_ids[
                                                    index
                                                ]
                                            "
                                        >
                                            <option value="">
                                                Choose a product
                                            </option>
                                            <option
                                                v-for="product in choices(
                                                    section.key
                                                )"
                                                :key="product.id"
                                                :value="product.id"
                                            >
                                                {{ product.title }} ·
                                                {{ product.category
                                                }}{{
                                                    !product.digital &&
                                                    product.stock === 0
                                                        ? ' (out of stock)'
                                                        : ''
                                                }}
                                            </option>
                                            <option
                                                v-if="
                                                    id &&
                                                    !choices(section.key).some(
                                                        (p) => p.id === id
                                                    )
                                                "
                                                :value="id"
                                            >
                                                Unavailable product: {{ id }}
                                            </option>
                                        </select></label
                                    >
                                    <button
                                        type="button"
                                        @click="
                                            form[
                                                section.key
                                            ].product_ids.splice(index, 1)
                                        "
                                        :disabled="form.processing"
                                    >
                                        Remove
                                    </button>
                                </div>
                                <button
                                    type="button"
                                    class="button button-outline"
                                    @click="addPin(section.key)"
                                    :disabled="
                                        form.processing ||
                                        form[section.key].product_ids.length >=
                                            24
                                    "
                                >
                                    Add product
                                </button>
                            </fieldset>
                        </div>
                        <template #footer="{ requestClose }"
                            ><button
                                type="button"
                                class="button button-outline"
                                @click="requestClose"
                            >
                                Cancel</button
                            ><button
                                type="button"
                                class="button button-navy"
                                @click="applyPlacement"
                            >
                                Apply to draft
                            </button></template
                        >
                    </AdminDrawer>
                </section>
                <section class="workspace-panel">
                    <h2>Category navigation</h2>
                    <label class="check"
                        ><input
                            type="checkbox"
                            v-model="form.category_navigation.enabled"
                        />Show category navigation</label
                    ><label
                        >Section title<input
                            v-model="form.category_navigation.title"
                            maxlength="100"
                    /></label>
                </section>
                <details class="selection-rules">
                    <summary>How automatic selection works</summary>
                    <p>
                        Automatic selections prioritize usable product imagery.
                        Only products matching the shopper’s current filters and
                        available at their minimum order quantity can appear.
                        Side cards exclude every hero product. Expired discounts
                        are not treated as active promotions. Empty sections
                        disappear. These rules never invent sales or ratings.
                    </p>
                </details>
                <div class="editor-publish-bar">
                    <span
                        >Publishing applies this draft to the live
                        storefront.</span
                    >
                    <button
                        class="button button-navy"
                        :disabled="form.processing || !draftDirty"
                    >
                        {{
                            form.processing
                                ? 'Publishing…'
                                : 'Publish storefront settings'
                        }}
                    </button>
                </div>
            </fieldset>
        </form>
    </div>
</template>
<script setup lang="ts">
import { computed, ref } from 'vue'
import AdminDrawer from '@/components/AdminDrawer.vue'
import { useAdminDraftRecovery } from '@/design/adminDraftRecovery'
import AdminNavigationGuard from '@/components/AdminNavigationGuard.vue'
import SokoDialog from '@/components/SokoDialog.vue'
import { Link, useForm, usePage } from '@inertiajs/vue3'
import {
    normalizeSectionOrder,
    storefrontSectionKeys,
} from '@/design/sectionOrder'
import '../../../css/workspace.css'
type Placement = {
    enabled: boolean
    title: string
    strategy: string
    product_ids: string[]
    category: string
    fallback: string
    limit: number
}
type PlacementKey = 'hero' | 'hero_side' | 'discovery'
type Settings = {
    section_order: string[]
    announcement: string
    version: number
    hero: Placement
    hero_side: Placement
    discovery: Placement
    category_navigation: { enabled: boolean; title: string }
}
type Product = {
    id: string
    title: string
    category: string
    digital: boolean
    stock: number
}
const props = defineProps<{ settings: Settings; products: Product[] }>()
const page = usePage()
const flash = computed(
    () => (page.props.flash as { error?: string; success?: string }) || {}
)
const sectionLabels: Record<string, string> = {
    showcase: 'Product showcase',
    category_navigation: 'Category navigation',
    catalog: 'Product catalog',
    discovery: 'Digital discovery',
    sourcing: 'Sourcing help',
}
const defaultSectionOrder = () => [...storefrontSectionKeys]
const orderAnnouncement = ref('')
const form = useForm<Settings>({
    ...JSON.parse(JSON.stringify(props.settings)),
    section_order: [...normalizeSectionOrder(props.settings.section_order)],
})
const baseline = ref(
    JSON.stringify({
        ...JSON.parse(JSON.stringify(props.settings)),
        section_order: [...normalizeSectionOrder(props.settings.section_order)],
    })
)
const draftValue = () =>
    Object.fromEntries(
        Object.keys(JSON.parse(baseline.value)).map((key) => [
            key,
            (form as any)[key],
        ])
    )
const draftDirty = computed(
    () => JSON.stringify(draftValue()) !== baseline.value
)
const confirmDiscard = ref(false)
const editingKey = ref<PlacementKey | null>(null)
const placementBaseline = ref('')
const placementError = ref('')
const placementDirty = computed(() =>
    Boolean(
        editingKey.value &&
        JSON.stringify(form[editingKey.value]) !== placementBaseline.value
    )
)
const strategyLabels: Record<string, string> = {
    featured: 'Featured first',
    newest: 'Newest first',
    promotions: 'Active offers',
    manual: 'Hand-picked order',
}
function openPlacement(key: PlacementKey) {
    editingKey.value = key
    placementBaseline.value = JSON.stringify(form[key])
    placementError.value = ''
}
function discardPlacement() {
    if (editingKey.value && placementBaseline.value)
        form[editingKey.value] = JSON.parse(placementBaseline.value)
    editingKey.value = null
    placementError.value = ''
}
function applyPlacement() {
    if (form.processing) return
    const key = editingKey.value
    if (!key) return
    const placement = form[key]
    if (
        !Number.isInteger(placement.limit) ||
        placement.limit < 1 ||
        placement.limit > (key === 'hero_side' ? 2 : 12)
    ) {
        placementError.value = 'Choose a valid maximum product count.'
        return
    }
    if (placement.enabled && placement.strategy === 'manual') {
        if (placement.product_ids.some((id) => !id)) {
            placementError.value =
                'Choose a product for every position, or remove unused positions.'
            return
        }
        if (
            new Set(placement.product_ids).size !== placement.product_ids.length
        ) {
            placementError.value =
                'Each product can appear only once in this placement.'
            return
        }
        if (
            placement.product_ids.some(
                (id) => !choices(key).some((product) => product.id === id)
            )
        ) {
            placementError.value =
                'One selected product no longer matches this category or placement. Choose a matching product.'
            return
        }
    }
    editingKey.value = null
    placementError.value = ''
}
function confirmDiscardDraft() {
    discardDraft()
    confirmDiscard.value = false
}
function discardDraft() {
    if (form.processing) return
    Object.assign(form, JSON.parse(baseline.value))
    editingKey.value = null
    placementError.value = ''
}
const sections: { key: PlacementKey; label: string; help: string }[] = [
    {
        key: 'hero',
        label: 'Main product showcase',
        help: 'Choose the products in the main carousel.',
    },
    {
        key: 'hero_side',
        label: 'Side product cards',
        help: 'Up to two additional products. These never repeat a product in the main showcase.',
    },
    {
        key: 'discovery',
        label: 'Digital discovery',
        help: 'A shelf of software, assets and licenses. Only digital listings are eligible.',
    },
]
const categories = computed(() =>
    [...new Set(props.products.map((p) => p.category))].filter(Boolean).sort()
)
function choices(key: PlacementKey) {
    return props.products.filter(
        (p) =>
            (key !== 'discovery' || p.digital) &&
            (!form[key].category || p.category === form[key].category)
    )
}
function addPin(key: PlacementKey) {
    if (!form.processing && form[key].product_ids.length < 24)
        form[key].product_ids.push('')
}
function moveSection(index: number, direction: number) {
    const next = index + direction
    if (form.processing || next < 0 || next >= form.section_order.length) return
    const [key] = form.section_order.splice(index, 1)
    form.section_order.splice(next, 0, key)
    orderAnnouncement.value = `${sectionLabels[key]} moved to position ${next + 1} of ${form.section_order.length}. Publish to apply.`
}
function resetSectionOrder() {
    if (form.processing) return
    form.section_order = defaultSectionOrder()
    orderAnnouncement.value =
        'Default section order restored in this draft. Publish to apply.'
}
function save() {
    if (form.processing) return
    if (editingKey.value) {
        applyPlacement()
        return
    }
    const submitted = JSON.parse(JSON.stringify(draftValue()))
    form.post('/dashboard/merchandising', {
        preserveScroll: true,
        onSuccess: (response) => {
            const feedback = response.props.flash as
                { error?: string; success?: string } | undefined
            if (feedback?.success && !feedback.error) {
                const unchanged =
                    JSON.stringify(draftValue()) === JSON.stringify(submitted)
                const acknowledged = JSON.parse(
                    JSON.stringify(
                        response.props.settings || {
                            ...submitted,
                            version: props.settings.version,
                        }
                    )
                )
                if (unchanged) Object.assign(form, acknowledged)
                else form.version = acknowledged.version
                baseline.value = JSON.stringify(acknowledged)
            }
        },
    })
}
const recoveryActor = String((page.props.auth as any)?.user?.id || '')
useAdminDraftRecovery({
    key: `soko:merchandising-draft:${recoveryActor}`,
    // This page is served only after content.update/all authorization.
    allowed: () =>
        Boolean(
            recoveryActor &&
            (page.props.workspaceNavigation as any[])?.some(
                (link) => link.href === '/dashboard/merchandising'
            )
        ),
    read: draftValue,
    apply: (draft) => Object.assign(form, draft),
    dirty: () => draftDirty.value,
})
</script>
<style scoped>
.merchandising-editor {
    padding-block: 2rem;
}
.merchandising-form {
    display: grid;
    gap: 1.5rem;
}
.editor-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 1rem;
}
.merchandising-form label {
    display: grid;
    gap: 0.4rem;
    margin-block: 0.7rem;
    font-weight: 600;
}
.merchandising-form input,
.merchandising-form select {
    padding: 0.65rem;
    min-height: 44px;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    background: var(--soko-surface);
    color: var(--soko-ink);
    max-width: 100%;
}
.merchandising-form .check {
    display: flex;
    align-items: center;
}
.check input {
    min-height: auto;
}
.pin-row {
    display: flex;
    gap: 1rem;
    align-items: center;
}
.pin-row label {
    flex: 1;
    min-width: 0;
}
.merchandising-form button {
    justify-self: start;
}
@media (max-width: 650px) {
    .editor-grid {
        grid-template-columns: 1fr;
    }
    .pin-row {
        align-items: stretch;
        flex-direction: column;
    }
}
.section-order-list {
    display: grid;
    gap: 0.5rem;
    padding: 0;
    list-style: none;
    margin: 1.5rem 0;
}
.section-order-list li {
    display: flex;
    gap: 1rem;
    align-items: center;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    padding: 1rem;
}
.section-position {
    font-variant-numeric: tabular-nums;
    color: var(--soko-muted);
    font-size: 1.1rem;
    min-width: 1.5rem;
}
.section-order-copy {
    display: grid;
    gap: 0.25rem;
    flex: 1;
}
.section-order-copy small {
    color: var(--soko-muted);
    font-size: 0.75rem;
}
.section-order-controls {
    display: flex;
    gap: 0.5rem;
}
.section-order-controls button {
    padding: 0.5rem 0.75rem;
    min-height: 44px;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    font-size: 0.75rem;
}
.section-order-controls button:disabled {
    opacity: 0.45;
    cursor: not-allowed;
}
@media (max-width: 650px) {
    .section-order-list li {
        flex-wrap: wrap;
    }
    .section-order-controls {
        width: 100%;
        padding-left: 2.5rem;
    }
}

.editor-status {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 14px;
    padding: 12px 16px;
    background: var(--soko-surface);
    border: 1px solid var(--soko-line);
    border-left: 3px solid var(--soko-amber);
    border-radius: 6px;
    margin-bottom: 20px;
    font-size: 12px;
}
.editor-status > span:nth-child(2) {
    color: var(--soko-muted);
    margin-right: auto;
}
.placement-summary {
    display: grid;
    grid-template-columns: minmax(0, 1fr) minmax(240px, 1fr) auto;
    gap: 22px;
    align-items: center;
}
.placement-state {
    font-size: 10px;
    color: var(--soko-success);
    font-weight: 750;
}
.placement-state.disabled {
    color: var(--soko-muted);
}
.placement-summary h2 {
    margin: 6px 0;
    font-size: 17px;
}
.placement-summary p {
    font-size: 12px;
    color: var(--soko-muted);
}
.placement-summary dl {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 12px;
}
.placement-summary dt {
    font-size: 10px;
    color: var(--soko-muted);
    margin-bottom: 5px;
}
.placement-summary dd {
    font-size: 12px;
    font-weight: 650;
    margin: 0;
    overflow-wrap: anywhere;
}
.editor-publish-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
    flex-wrap: wrap;
    padding: 16px;
    border: 1px solid var(--soko-line);
    background: var(--soko-surface);
    border-radius: 6px;
}
.editor-publish-bar > span {
    font-size: 12px;
    color: var(--soko-muted);
}
.editor-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 20px;
}
@media (max-width: 1050px) {
    .placement-summary {
        grid-template-columns: 1fr;
    }
    .placement-summary > .button {
        justify-self: start;
    }
}
.merchandising-form-lock {
    display: contents;
}
</style>
