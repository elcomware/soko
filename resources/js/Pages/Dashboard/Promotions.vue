<script setup lang="ts">
import { computed, reactive, ref, onMounted, onBeforeUnmount } from 'vue'
import { Head } from '@inertiajs/vue3'
import AdminDrawer from '@/components/AdminDrawer.vue'
import AdminTable from '@/components/AdminTable.vue'
import AdminActionMenu from '@/components/AdminActionMenu.vue'
import { useMarketplaceForm } from '@/composables/useMarketplaceForm'
import '@/../css/commerce.css'
type Target = {
    id: string
    title: string
    category: string
    sub_category: string
    vendor_id?: string
    vendor_name?: string
}
const props = defineProps<{
    campaigns: any[]
    products: Target[]
    canManageAll: boolean
    actorId: string
}>()
const { busy, error, post, validate } = useMarketplaceForm()
const freshCampaign = () => ({
    id: '',
    name: '',
    kind: 'percentage',
    amount: 10,
    enabled: false,
    starts_at: '',
    ends_at: '',
    vendor_id: props.canManageAll ? '' : props.actorId,
    version: '',
    all: true,
    product_ids: [] as string[],
    categories: [] as string[],
    sub_categories: [] as string[],
    exclude_product_ids: [] as string[],
    exclude_categories: [] as string[],
    exclude_sub_categories: [] as string[],
})
const form = reactive(freshCampaign())
const targetSearch = ref('')
const editorOpen = ref(false)
const savedDraft = ref(JSON.stringify(form))
const dirty = computed(() => JSON.stringify(form) !== savedDraft.value)
const currentTime = ref(Date.now())
let statusClock: ReturnType<typeof setInterval> | undefined
onMounted(() => {
    statusClock = setInterval(() => {
        currentTime.value = Date.now()
    }, 1000)
})
onBeforeUnmount(() => {
    if (statusClock) clearInterval(statusClock)
})
const campaignSearch = ref('')
const statusFilter = ref('all')
const visibleCampaigns = computed(() =>
    props.campaigns.filter((c) => {
        const query = campaignSearch.value.trim().toLowerCase()
        return (
            (statusFilter.value === 'all' ||
                status(c) === statusFilter.value) &&
            (!query ||
                [c.name, merchantName(c.vendor_id), c.kind]
                    .join(' ')
                    .toLowerCase()
                    .includes(query))
        )
    })
)
function clearCampaignFilters() {
    campaignSearch.value = ''
    statusFilter.value = 'all'
}
function createCampaign() {
    if (busy.value || editorOpen.value) return
    resetDraft()
    editorOpen.value = true
}
function closeEditor() {
    if (busy.value) return
    resetDraft()
}

const merchants = computed(() => {
    const owners = new Map<string, string>()
    for (const p of props.products)
        if (p.vendor_id) owners.set(p.vendor_id, p.vendor_name || p.vendor_id)
    return [...owners]
        .map(([id, name]) => ({ id, name }))
        .sort((a, b) => a.name.localeCompare(b.name))
})
const merchantName = (id: string) =>
    merchants.value.find((m) => m.id === id)?.name || id
const scopedProducts = computed(() =>
    props.products.filter(
        (p) => !form.vendor_id || p.vendor_id === form.vendor_id
    )
)
const selectableProducts = computed(() =>
    scopedProducts.value.filter(
        (p) =>
            !targetSearch.value ||
            (p.title + ' ' + p.category)
                .toLowerCase()
                .includes(targetSearch.value.toLowerCase())
    )
)
const eligibleProducts = computed(() =>
    scopedProducts.value.filter((p) => {
        const sub = JSON.stringify({
            category: p.category,
            sub_category: p.sub_category,
        })
        if (
            form.exclude_product_ids.includes(p.id) ||
            form.exclude_categories.includes(p.category) ||
            form.exclude_sub_categories.includes(sub)
        )
            return false
        return (
            form.all ||
            form.product_ids.includes(p.id) ||
            form.categories.includes(p.category) ||
            form.sub_categories.includes(sub)
        )
    })
)
const campaignCounts = computed(() => ({
    active: props.campaigns.filter((c) => status(c) === 'Active').length,
    scheduled: props.campaigns.filter((c) => status(c) === 'Scheduled').length,
    total: props.campaigns.length,
}))
const categories = computed(() =>
    [...new Set(scopedProducts.value.map((p) => p.category))].sort()
)
const subcategories = computed(() =>
    [
        ...new Set(
            scopedProducts.value
                .filter((p) => p.sub_category)
                .map((p) =>
                    JSON.stringify({
                        category: p.category,
                        sub_category: p.sub_category,
                    })
                )
        ),
    ].sort()
)
const sublabel = (value: string) => {
    const s = JSON.parse(value)
    return `${s.category} / ${s.sub_category}`
}
const localTime = (value: string) => {
    const d = new Date(value)
    return new Date(d.getTime() - d.getTimezoneOffset() * 60000)
        .toISOString()
        .slice(0, 16)
}
function edit(c: any) {
    if (busy.value) return
    error.value = ''
    Object.assign(form, {
        id: c.id,
        name: c.name,
        kind: c.kind,
        amount: c.amount,
        enabled: c.enabled,
        starts_at: localTime(c.starts_at),
        ends_at: localTime(c.ends_at),
        vendor_id: c.vendor_id,
        version: c.updated_at,
        all: c.scope.all,
        product_ids: [...(c.scope.product_ids || [])],
        categories: [...(c.scope.categories || [])],
        sub_categories: (c.scope.sub_categories || []).map((s: any) =>
            JSON.stringify(s)
        ),
        exclude_product_ids: [...(c.scope.exclude_product_ids || [])],
        exclude_categories: [...(c.scope.exclude_categories || [])],
        exclude_sub_categories: (c.scope.exclude_sub_categories || []).map(
            (s: any) => JSON.stringify(s)
        ),
    })
    savedDraft.value = JSON.stringify(form)
    targetSearch.value = ''
    editorOpen.value = true
}
function reset() {
    if (busy.value) return
    resetDraft()
}
function resetDraft() {
    Object.assign(form, freshCampaign())
    error.value = ''
    targetSearch.value = ''
    savedDraft.value = JSON.stringify(form)
    editorOpen.value = false
}

function save(e: Event) {
    if (busy.value || !validate(e)) return
    if (
        !Number.isFinite(Date.parse(form.starts_at)) ||
        !Number.isFinite(Date.parse(form.ends_at))
    ) {
        error.value = 'Choose a valid start and end time.'
        return
    }
    const scope = {
        all: form.all,
        product_ids: form.product_ids,
        categories: form.categories,
        sub_categories: form.sub_categories.map((s) => JSON.parse(s)),
        exclude_product_ids: form.exclude_product_ids,
        exclude_categories: form.exclude_categories,
        exclude_sub_categories: form.exclude_sub_categories.map((s) =>
            JSON.parse(s)
        ),
    }
    post(
        '/dashboard/promotions' + (form.id ? '/' + form.id : ''),
        {
            name: form.name,
            kind: form.kind,
            amount: form.amount,
            enabled: String(form.enabled),
            starts_at: new Date(form.starts_at).toISOString(),
            ends_at: new Date(form.ends_at).toISOString(),
            vendor_id: form.vendor_id,
            version: form.version,
            scope: JSON.stringify(scope),
        },
        resetDraft
    )
}
function status(c: any) {
    const now = currentTime.value
    return !c.enabled
        ? 'Disabled'
        : now < Date.parse(c.starts_at)
          ? 'Scheduled'
          : now >= Date.parse(c.ends_at)
            ? 'Ended'
            : 'Active'
}
</script>
<template>
    <Head title="Promotions" />
    <div class="promotions-admin">
        <div class="commerce-page-heading">
            <div>
                <p class="commerce-eyebrow">Merchandising</p>
                <h1>Promotions</h1>
                <p>
                    Optional, time-limited campaigns. Exclusions always take
                    precedence.
                </p>
            </div>
            <button
                type="button"
                class="button button-navy"
                :disabled="busy"
                @click="createCampaign"
            >
                Create campaign
            </button>
        </div>
        <div class="campaign-metrics" aria-label="Campaign overview">
            <div>
                <span>Active now</span
                ><strong>{{ campaignCounts.active }}</strong
                ><small>Eligible discounts in effect</small>
            </div>
            <div>
                <span>Scheduled</span
                ><strong>{{ campaignCounts.scheduled }}</strong
                ><small>Start and end automatically</small>
            </div>
            <div>
                <span>Total campaigns</span
                ><strong>{{ campaignCounts.total }}</strong
                ><small>Including disabled and ended</small>
            </div>
        </div>
        <section class="commerce-panel" aria-labelledby="campaign-list">
            <div class="campaign-list-heading">
                <div>
                    <h2 id="campaign-list">Campaigns</h2>
                    <p>
                        Review live promotions, prepare upcoming events, or
                        reopen a campaign to adjust its scope.
                    </p>
                </div>
                <span role="status"
                    >{{ visibleCampaigns.length }} of
                    {{ campaigns.length }}</span
                >
            </div>
            <div class="campaign-list-filters">
                <label
                    >Search campaigns<input
                        v-model="campaignSearch"
                        type="search"
                        placeholder="Campaign or merchant name" /></label
                ><label
                    >Status<select v-model="statusFilter">
                        <option value="all">All statuses</option>
                        <option
                            v-for="state in [
                                'Active',
                                'Scheduled',
                                'Disabled',
                                'Ended',
                            ]"
                            :key="state"
                            :value="state"
                        >
                            {{ state }}
                        </option>
                    </select></label
                ><button
                    v-if="campaignSearch || statusFilter !== 'all'"
                    type="button"
                    class="button button-outline"
                    @click="clearCampaignFilters"
                >
                    Clear filters
                </button>
            </div>
            <AdminTable
                caption="Promotion campaigns"
                :empty="!visibleCampaigns.length"
                :busy="busy"
                :empty-text="
                    campaigns.length
                        ? 'No campaigns match these filters. Try another name or status.'
                        : 'No campaigns yet. Create one when you are ready; regular prices stay unchanged.'
                "
            >
                <template #header>
                    <tr>
                        <th scope="col">Campaign</th>
                        <th scope="col">Coverage</th>
                        <th scope="col">Discount</th>
                        <th scope="col">Period</th>
                        <th scope="col">Status</th>
                        <th scope="col">Action</th>
                    </tr>
                </template>
                <tr v-for="c in visibleCampaigns" :key="c.id">
                    <th scope="row">{{ c.name }}</th>
                    <td>
                        {{
                            c.vendor_id
                                ? merchantName(c.vendor_id)
                                : 'All merchants'
                        }}<small class="coverage-caption">{{
                            c.scope.all
                                ? 'All products, with exceptions'
                                : 'Selected products and categories'
                        }}</small>
                    </td>
                    <td>
                        {{ c.amount
                        }}{{ c.kind === 'percentage' ? '%' : ' USD per unit' }}
                    </td>
                    <td>
                        {{ new Date(c.starts_at).toLocaleString() }} —
                        {{ new Date(c.ends_at).toLocaleString() }}
                    </td>
                    <td>
                        <span
                            class="campaign-status"
                            :data-status="status(c)"
                            >{{ status(c) }}</span
                        >
                    </td>
                    <td>
                        <AdminActionMenu
                            :label="`Actions for ${c.name}`"
                            :disabled="busy"
                            :items="[{ id: 'edit', label: 'Edit campaign' }]"
                            @select="edit(c)"
                        />
                    </td>
                </tr>
            </AdminTable>
        </section>
        <AdminDrawer
            :open="editorOpen"
            :title="form.id ? 'Edit campaign' : 'Create campaign'"
            description="Choose the audience, discount and schedule. Campaigns remain optional and stop automatically."
            :busy="busy"
            :dirty="dirty"
            @close="closeEditor"
        >
            <p class="commerce-muted">
                Times use your browser’s local time zone and are stored in UTC.
                Discounts end automatically at the end time. The best eligible
                single discount applies after quantity pricing; campaigns do not
                stack.
            </p>
            <p class="commerce-small">
                Merchant coverage:
                {{
                    form.vendor_id
                        ? 'Only ' + merchantName(form.vendor_id)
                        : 'All merchants'
                }}. Category and product selections below narrow this coverage.
            </p>
            <form
                id="campaign-form"
                class="soko-form campaign-editor-form"
                novalidate
                @submit.prevent="save"
            >
                <fieldset :disabled="busy" class="campaign-editor-fields">
                    <p
                        v-if="error"
                        id="market-form-error"
                        class="form-error"
                        role="alert"
                    >
                        {{ error }}
                    </p>
                    <div v-if="canManageAll" class="merchant-selector">
                        <label for="campaign-merchant">Merchant coverage</label>
                        <select id="campaign-merchant" v-model="form.vendor_id">
                            <option value="">All merchants</option>
                            <option
                                v-if="
                                    form.vendor_id &&
                                    !merchants.some(
                                        (m) => m.id === form.vendor_id
                                    )
                                "
                                :value="form.vendor_id"
                            >
                                {{ merchantName(form.vendor_id) }}
                            </option>
                            <option
                                v-for="merchant in merchants"
                                :key="merchant.id"
                                :value="merchant.id"
                            >
                                {{ merchant.name }}
                            </option>
                        </select>
                        <p class="commerce-small">
                            Limit this campaign to one merchant, or manage a
                            store-wide event. Exclusions still apply.
                        </p>
                    </div>
                    <h3 class="campaign-step">1. Campaign and schedule</h3>
                    <label for="campaign-name">Campaign name</label
                    ><input
                        id="campaign-name"
                        v-model="form.name"
                        maxlength="120"
                        required
                    />
                    <div class="campaign-grid">
                        <div>
                            <label for="campaign-kind">Discount type</label
                            ><select id="campaign-kind" v-model="form.kind">
                                <option value="percentage">Percentage</option>
                                <option value="fixed">
                                    Fixed USD amount per unit
                                </option>
                            </select>
                        </div>
                        <div>
                            <label for="campaign-amount">{{
                                form.kind === 'percentage'
                                    ? 'Percentage off'
                                    : 'USD off per unit'
                            }}</label
                            ><input
                                id="campaign-amount"
                                v-model.number="form.amount"
                                type="number"
                                min="0.01"
                                :max="
                                    form.kind === 'percentage'
                                        ? 100
                                        : 1000000000000
                                "
                                step="0.01"
                                required
                            />
                        </div>
                        <div>
                            <label for="campaign-start">Starts</label
                            ><input
                                id="campaign-start"
                                v-model="form.starts_at"
                                type="datetime-local"
                                required
                            />
                        </div>
                        <div>
                            <label for="campaign-end">Ends</label
                            ><input
                                id="campaign-end"
                                v-model="form.ends_at"
                                type="datetime-local"
                                :min="form.starts_at"
                                required
                            />
                        </div>
                    </div>
                    <label
                        ><input v-model="form.enabled" type="checkbox" /> Enable
                        this campaign during its scheduled period</label
                    >
                    <label
                        ><input v-model="form.all" type="checkbox" /> Include
                        all eligible products</label
                    >
                    <h3 class="campaign-step">2. Product coverage</h3>
                    <label for="campaign-product-search"
                        >Find products in selection lists</label
                    >
                    <input
                        id="campaign-product-search"
                        v-model="targetSearch"
                        type="search"
                        placeholder="Search product names or categories"
                    />
                    <div v-if="!form.all" class="campaign-grid">
                        <div>
                            <label for="campaign-categories">Categories</label
                            ><select
                                id="campaign-categories"
                                v-model="form.categories"
                                multiple
                                size="5"
                            >
                                <option
                                    v-for="c in categories"
                                    :key="c"
                                    :value="c"
                                >
                                    {{ c }}
                                </option>
                            </select>
                        </div>
                        <div>
                            <label for="campaign-subcategories"
                                >Subcategories</label
                            ><select
                                id="campaign-subcategories"
                                v-model="form.sub_categories"
                                multiple
                                size="5"
                            >
                                <option
                                    v-for="s in subcategories"
                                    :key="s"
                                    :value="s"
                                >
                                    {{ sublabel(s) }}
                                </option>
                            </select>
                        </div>
                        <div>
                            <label for="campaign-products"
                                >Individual products</label
                            ><select
                                id="campaign-products"
                                v-model="form.product_ids"
                                multiple
                                size="5"
                            >
                                <option
                                    v-for="p in selectableProducts"
                                    :key="p.id"
                                    :value="p.id"
                                >
                                    {{ p.title }}
                                </option>
                            </select>
                        </div>
                    </div>
                    <fieldset>
                        <legend>Exceptions</legend>
                        <p class="commerce-small">
                            Select any products, categories or subcategories
                            that must keep their usual price. Use Ctrl/Cmd to
                            select more than one option.
                        </p>
                        <div class="campaign-grid">
                            <div>
                                <label for="campaign-exclude-categories"
                                    >Exclude categories</label
                                ><select
                                    id="campaign-exclude-categories"
                                    v-model="form.exclude_categories"
                                    multiple
                                    size="4"
                                >
                                    <option
                                        v-for="c in categories"
                                        :key="c"
                                        :value="c"
                                    >
                                        {{ c }}
                                    </option>
                                </select>
                            </div>
                            <div>
                                <label for="campaign-exclude-subcategories"
                                    >Exclude subcategories</label
                                ><select
                                    id="campaign-exclude-subcategories"
                                    v-model="form.exclude_sub_categories"
                                    multiple
                                    size="4"
                                >
                                    <option
                                        v-for="s in subcategories"
                                        :key="s"
                                        :value="s"
                                    >
                                        {{ sublabel(s) }}
                                    </option>
                                </select>
                            </div>
                            <div>
                                <label for="campaign-exclude-products"
                                    >Exclude products</label
                                ><select
                                    id="campaign-exclude-products"
                                    v-model="form.exclude_product_ids"
                                    multiple
                                    size="4"
                                >
                                    <option
                                        v-for="p in selectableProducts"
                                        :key="p.id"
                                        :value="p.id"
                                    >
                                        {{ p.title }}
                                    </option>
                                </select>
                            </div>
                        </div>
                    </fieldset>
                    <section
                        class="campaign-preview"
                        aria-labelledby="coverage-preview"
                    >
                        <div>
                            <p class="commerce-eyebrow">Coverage preview</p>
                            <h3 id="coverage-preview">
                                {{ eligibleProducts.length }} current products
                                match
                            </h3>
                            <p>
                                Based on selected merchant, inclusions and
                                exclusions. This is targeting coverage, not a
                                price quote; stock, quantity tiers and competing
                                campaigns are checked when pricing.
                            </p>
                        </div>
                        <ul v-if="eligibleProducts.length">
                            <li
                                v-for="p in eligibleProducts.slice(0, 6)"
                                :key="p.id"
                            >
                                {{ p.title }}<span>{{ p.category }}</span>
                            </li>
                        </ul>
                        <p v-if="eligibleProducts.length > 6">
                            And {{ eligibleProducts.length - 6 }} more matching
                            products.
                        </p>
                        <p v-if="!eligibleProducts.length">
                            No current products match. A category or
                            all-products campaign can apply to matching products
                            added later.
                        </p>
                        <p>
                            <strong>{{
                                form.enabled
                                    ? 'Enabled for the selected period'
                                    : 'Disabled until you enable it'
                            }}</strong>
                            ·
                            {{
                                form.kind === 'percentage'
                                    ? form.amount + '% off'
                                    : form.amount + ' USD off per unit'
                            }}
                        </p>
                    </section>
                </fieldset>
            </form>
            <template #footer="{ requestClose }">
                <p class="campaign-save-note">
                    {{
                        form.enabled
                            ? 'Enabled only within its scheduled period.'
                            : 'Saved as disabled. Enable when ready.'
                    }}
                </p>
                <div class="campaign-actions">
                    <button
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="requestClose"
                    >
                        Cancel</button
                    ><button
                        type="submit"
                        form="campaign-form"
                        class="button button-navy"
                        :disabled="busy"
                    >
                        {{
                            busy
                                ? 'Saving…'
                                : form.id
                                  ? 'Save changes'
                                  : 'Create campaign'
                        }}
                    </button>
                </div>
            </template>
        </AdminDrawer>
    </div>
</template>
<style scoped>
.campaign-editor-fields {
    display: grid;
    gap: 1rem;
    border: 0;
    padding: 0;
    margin: 0;
    min-width: 0;
}
.campaign-save-note {
    color: var(--soko-muted);
    font-size: 0.8125rem;
    margin: 0 0 0.75rem;
}
.campaign-list-heading {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    align-items: start;
    margin-bottom: 1.25rem;
}
.campaign-list-heading p {
    color: var(--soko-muted);
    font-size: 0.875rem;
    max-width: 65ch;
}
.campaign-list-heading > span {
    color: var(--soko-muted);
    white-space: nowrap;
    font-size: 0.8125rem;
}
.campaign-list-filters {
    display: flex;
    align-items: end;
    gap: 1rem;
    margin-bottom: 1.25rem;
    flex-wrap: wrap;
}
.campaign-list-filters label {
    display: grid;
    gap: 0.4rem;
    font-size: 0.8125rem;
    font-weight: 650;
    flex: 1;
    min-width: min(100%, 12rem);
}
.campaign-list-filters input,
.campaign-list-filters select {
    min-height: 44px;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    padding: 0.65rem 0.8rem;
    background: var(--soko-surface);
    color: var(--soko-ink);
}

.promotions-admin {
    width: 100%;
    min-width: 0;
}
.campaign-metrics {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 1rem;
}
.campaign-metrics > div {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    background: var(--soko-white, #fff);
    border: 1px solid var(--soko-border);
    padding: 1.25rem;
    border-radius: var(--soko-radius);
}
.campaign-metrics strong {
    font-size: 2rem;
    color: var(--soko-navy);
}
.campaign-metrics small,
.coverage-caption {
    display: block;
    color: var(--soko-muted);
    font-size: 0.8rem;
}
.merchant-selector {
    display: grid;
    gap: 0.5rem;
    padding: 1rem;
    background: var(--soko-paper, #f5f4ef);
    border-radius: var(--soko-radius);
}
.campaign-step {
    margin: 1.5rem 0 0.75rem;
    border-top: 1px solid var(--soko-border);
    padding-top: 1.25rem;
}
.campaign-preview {
    margin-top: 1rem;
    padding: 1.25rem;
    border: 1px solid var(--soko-border);
    background: var(--soko-paper, #f5f4ef);
    border-radius: var(--soko-radius);
}
.campaign-preview ul {
    list-style: none;
    padding: 0;
    display: grid;
    gap: 0.5rem;
}
.campaign-preview li {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    border-bottom: 1px solid var(--soko-border);
    padding: 0.5rem 0;
}
.campaign-preview li span {
    color: var(--soko-muted);
    font-size: 0.85rem;
}
.campaign-status {
    display: inline-flex;
    border-radius: 999px;
    padding: 0.3rem 0.65rem;
    background: var(--soko-paper, #f5f4ef);
    font-size: 0.8rem;
    border: 1px solid var(--soko-border);
}
.campaign-status[data-status='Active'] {
    background: var(--soko-navy);
    color: white;
}
@media (max-width: 640px) {
    .campaign-metrics {
        grid-template-columns: 1fr;
    }
    .campaign-preview li {
        flex-direction: column;
        gap: 0.25rem;
    }
}

.commerce-panel {
    margin-block: 1.5rem;
    padding: clamp(1rem, 3vw, 2rem);
}
.campaign-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(min(100%, 15rem), 1fr));
    gap: 1rem;
}
.campaign-grid > div {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    min-width: 0;
}
.campaign-grid select {
    width: 100%;
}
fieldset:not(.campaign-editor-fields) {
    border: 1px solid var(--soko-border);
    padding: 1rem;
    border-radius: var(--soko-radius);
}
.campaign-actions {
    display: flex;
    gap: 0.75rem;
}
</style>
