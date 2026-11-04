<template>
    <div class="site-container merchandising-editor">
        <header class="workspace-heading">
            <div>
                <p class="eyebrow">Storefront editor</p>
                <h1>Put the right products in view</h1>
                <p>
                    Choose real products, selection rules and the sections
                    shoppers see.
                </p>
            </div>
            <Link href="/dashboard" class="button button-outline"
                >Back to workspace</Link
            >
        </header>
        <p v-if="flash.error" class="form-error" role="alert">
            {{ flash.error }} Reload this page before retrying if another
            administrator changed these settings.
        </p>
        <form @submit.prevent="save" class="merchandising-form">
            <section
                class="workspace-panel"
                aria-labelledby="section-order-title"
            >
                <h2 id="section-order-title">Page section order</h2>
                <p>
                    Move sections into the order shoppers should encounter them.
                    The product catalog stays available. Hidden or empty
                    sections keep their place but do not leave a gap.
                </p>
                <ol class="section-order-list">
                    <li v-for="(key, index) in form.section_order" :key="key">
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
                    Plain text only. Keep claims accurate; leave blank to show
                    no announcement text.
                </p>
            </section>
            <section
                v-for="section in sections"
                :key="section.key"
                class="workspace-panel"
            >
                <h2>{{ section.label }}</h2>
                <p>{{ section.help }}</p>
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
                            <option value="newest">Newest products</option>
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
                            :max="section.key === 'hero_side' ? 2 : 12"
                            required
                    /></label>
                    <label
                        >If there are not enough matching products<select
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
                <fieldset v-if="form[section.key].strategy === 'manual'">
                    <legend>Ordered product selection</legend>
                    <div
                        v-for="(id, index) in form[section.key].product_ids"
                        :key="index"
                        class="pin-row"
                    >
                        <label
                            >Position {{ index + 1
                            }}<select
                                v-model="form[section.key].product_ids[index]"
                            >
                                <option value="">Choose a product</option>
                                <option
                                    v-for="product in choices(section.key)"
                                    :key="product.id"
                                    :value="product.id"
                                >
                                    {{ product.title }} · {{ product.category
                                    }}{{
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
                                form[section.key].product_ids.splice(index, 1)
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
                            form[section.key].product_ids.length >= 24
                        "
                    >
                        Add product
                    </button>
                </fieldset>
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
            <p>
                Automatic selections prioritize usable product imagery. Only
                products matching the shopper’s current filters and available at
                their minimum order quantity can appear. Side cards exclude
                every hero product. Expired discounts are not treated as active
                promotions. Empty sections disappear. These rules never invent
                sales or ratings.
            </p>
            <button class="button button-navy" :disabled="form.processing">
                {{
                    form.processing
                        ? 'Publishing…'
                        : 'Publish storefront settings'
                }}
            </button>
        </form>
    </div>
</template>
<script setup lang="ts">
import { computed, ref } from 'vue'
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
    form.post('/dashboard/merchandising', {
        preserveScroll: true,
        onSuccess: (response) => {
            const feedback = response.props.flash as
                { error?: string; success?: string } | undefined
            if (feedback?.success && !feedback.error)
                form.version = props.settings.version
        },
    })
}
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
    border: 1px solid var(--soko-border);
    border-radius: var(--soko-radius);
    background: var(--soko-surface);
    color: var(--soko-text);
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
</style>
