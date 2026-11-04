<script setup lang="ts">
import {
    computed,
    reactive,
    ref,
    onMounted,
    onBeforeUnmount,
    watch,
    nextTick,
} from 'vue'
import { Head, router } from '@inertiajs/vue3'
import AdminPagination from '@/components/AdminPagination.vue'
import { useAdminDraftRecovery } from '@/design/adminDraftRecovery'
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
    campaignPage?: {
        page: number
        per_page: number
        total: number
        filtered_total: number
        last_page: number
        from: number
        to: number
        filters: { q: string; status: string; sort: string }
    }
    campaignCounts?: {
        total: number
        active: number
        scheduled: number
        ended: number
        disabled: number
    }
    serverTime?: string
    nextStatusChange?: string
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
const sourceSchedule = reactive({
    starts_at: '',
    ends_at: '',
    local_start: '',
    local_end: '',
})
const draftRecovered = ref(false)
const pendingCampaignRecovery = ref<Record<string, any> | null>(null)
const targetSearch = ref('')
const editorOpen = ref(false)
const savedDraft = ref(JSON.stringify(form))
const dirty = computed(() => JSON.stringify(form) !== savedDraft.value)
let serverOffset = props.serverTime
    ? Date.parse(props.serverTime) - Date.now()
    : 0
const currentTime = ref(Date.now() + serverOffset)
let statusClock: ReturnType<typeof setInterval> | undefined
onMounted(() => {
    statusClock = setInterval(() => {
        currentTime.value = Date.now() + serverOffset
    }, 1000)
})
onBeforeUnmount(() => {
    if (statusClock) clearInterval(statusClock)
})
const campaignSearch = ref(props.campaignPage?.filters.q || '')
const statusFilter = ref(props.campaignPage?.filters.status || 'all')
const campaignSort = ref(props.campaignPage?.filters.sort || 'newest')
const visibleCampaigns = computed(() =>
    props.campaignPage
        ? props.campaigns
        : props.campaigns.filter((c) => {
              const query = campaignSearch.value.trim().toLowerCase()
              return (
                  (statusFilter.value === 'all' ||
                      status(c).toLowerCase() ===
                          statusFilter.value.toLowerCase()) &&
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
    applyCampaignFilters()
}
function createCampaign() {
    if (
        busy.value ||
        editorOpen.value ||
        listLoading.value ||
        filtersPending.value ||
        pendingCampaignRecovery.value
    )
        return
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
const campaignCounts = computed(
    () =>
        props.campaignCounts || {
            active: props.campaigns.filter((c) => status(c) === 'Active')
                .length,
            scheduled: props.campaigns.filter((c) => status(c) === 'Scheduled')
                .length,
            total: props.campaigns.length,
        }
)
const categories = computed(() =>
    [
        ...new Set(
            scopedProducts.value
                .map((p) => p.category)
                .filter((value) => value.trim())
        ),
    ].sort()
)
const subcategories = computed(() =>
    [
        ...new Set(
            scopedProducts.value
                .filter((p) => p.category.trim() && p.sub_category.trim())
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
type TargetKind = 'categories' | 'sub_categories' | 'product_ids'
type TargetBucket = 'include' | 'exclude'
type TargetKey =
    | TargetKind
    | 'exclude_categories'
    | 'exclude_sub_categories'
    | 'exclude_product_ids'
const targetKinds: { key: TargetKind; label: string }[] = [
    { key: 'categories', label: 'Categories' },
    { key: 'sub_categories', label: 'Subcategories' },
    { key: 'product_ids', label: 'Products' },
]
const targetPanels = computed(() =>
    form.all
        ? (['exclude'] as TargetBucket[])
        : (['include', 'exclude'] as TargetBucket[])
)
const activeTargetKind = reactive<Record<TargetBucket, TargetKind>>({
    include: 'categories',
    exclude: 'categories',
})
const targetKey = (bucket: TargetBucket, kind: TargetKind): TargetKey =>
    bucket === 'exclude' ? (`exclude_${kind}` as TargetKey) : kind
function targetOptions(kind: TargetKind) {
    if (kind === 'categories')
        return categories.value.map((value) => {
            const count = scopedProducts.value.filter(
                (p) => p.category === value
            ).length
            return {
                value,
                label: value || 'Uncategorized',
                context: `${count} ${count === 1 ? 'product' : 'products'}`,
            }
        })
    if (kind === 'sub_categories')
        return subcategories.value.map((value) => ({
            value,
            label: sublabel(value),
            context: 'Category / subcategory',
        }))
    return scopedProducts.value.map((p) => ({
        value: p.id,
        label: p.title,
        context: [p.category, p.sub_category, p.vendor_name]
            .filter(Boolean)
            .join(' / '),
    }))
}
function filteredTargetOptions(kind: TargetKind) {
    const query = targetSearch.value.trim().toLowerCase()
    return targetOptions(kind).filter((option) =>
        `${option.label} ${option.context}`.toLowerCase().includes(query)
    )
}
function selectedTargets(bucket: TargetBucket) {
    return targetKinds.flatMap((kind) => {
        const key = targetKey(bucket, kind.key)
        const options = targetOptions(kind.key)
        return form[key].map((value) => ({
            key,
            value,
            kind: kind.label,
            label:
                options.find((option) => option.value === value)?.label ||
                (kind.key === 'sub_categories' ? sublabel(value) : value),
            unavailable: !options.some((option) => option.value === value),
        }))
    })
}
function removeTarget(key: TargetKey, value: string) {
    if (busy.value) return
    form[key] = form[key].filter((item) => item !== value)
}
const listLoading = ref(false)
const listError = ref('')
const pageSize = ref(props.campaignPage?.per_page || 25)
const listHeight = ref(480)
const campaignColumns = [
    { key: 'name', label: 'Campaign', width: 220 },
    { key: 'coverage', label: 'Coverage', width: 220 },
    { key: 'amount', label: 'Discount', width: 150 },
    { key: 'period', label: 'Period', width: 280 },
    { key: 'status', label: 'Status', width: 130 },
    { key: 'action', label: 'Action', width: 100 },
]
const filterSignature = () =>
    JSON.stringify([
        campaignSearch.value.trim(),
        statusFilter.value === 'all' ? '' : statusFilter.value.toLowerCase(),
        campaignSort.value,
        pageSize.value,
    ])
const filtersPending = computed(() =>
    Boolean(
        props.campaignPage &&
        filterSignature() !==
            JSON.stringify([
                props.campaignPage.filters.q,
                props.campaignPage.filters.status,
                props.campaignPage.filters.sort,
                props.campaignPage.per_page,
            ])
    )
)
const statusSnapshotExpired = computed(() =>
    Boolean(
        props.nextStatusChange &&
        currentTime.value >= Date.parse(props.nextStatusChange)
    )
)
let searchTimer: ReturnType<typeof setTimeout> | undefined
let queuedSearch = false,
    disposed = false,
    submittedFilters = '',
    syncingFilters = false
function campaignListUrl(page = 1) {
    const query = new URLSearchParams({
        page: String(page),
        per_page: String(pageSize.value),
        sort: campaignSort.value,
    })
    if (campaignSearch.value.trim()) query.set('q', campaignSearch.value.trim())
    if (statusFilter.value !== 'all')
        query.set('status', statusFilter.value.toLowerCase())
    return '/dashboard/promotions?' + query.toString()
}
function visitCampaignList(
    url: string,
    success?: (response: any) => void,
    recovering = false
) {
    if (
        !props.campaignPage ||
        disposed ||
        busy.value ||
        editorOpen.value ||
        (pendingCampaignRecovery.value && !recovering)
    )
        return
    if (listLoading.value) {
        queuedSearch = true
        return
    }
    if (searchTimer) clearTimeout(searchTimer)
    queuedSearch = false
    submittedFilters = filterSignature()
    listLoading.value = true
    listError.value = ''
    router.get(
        url,
        {},
        {
            preserveState: true,
            preserveScroll: true,
            onSuccess: (response) => {
                if (
                    !disposed &&
                    !queuedSearch &&
                    filterSignature() === submittedFilters
                )
                    success?.(response)
            },
            onError: () => {
                listError.value =
                    'Campaigns could not be loaded. Previous results are still shown.'
            },
            onNetworkError: () => {
                listError.value = 'Connection lost. Retry to refresh campaigns.'
                return false
            },
            onHttpException: () => {
                listError.value =
                    'Campaigns could not be loaded. Previous results are still shown.'
                return false
            },
            onFinish: () => {
                listLoading.value = false
                if (queuedSearch && !disposed) {
                    queuedSearch = false
                    void nextTick(applyCampaignFilters)
                }
            },
        }
    )
}
function applyCampaignFilters() {
    visitCampaignList(campaignListUrl(1))
}
function refreshCampaigns() {
    visitCampaignList(
        campaignListUrl(
            filtersPending.value ? 1 : props.campaignPage?.page || 1
        )
    )
}
function changeCampaignPage(page: number) {
    if (
        !props.campaignPage ||
        filtersPending.value ||
        page < 1 ||
        page > props.campaignPage.last_page
    )
        return
    visitCampaignList(campaignListUrl(page))
}
function changeCampaignSize(size: number) {
    if (![25, 50, 100].includes(size) || busy.value || editorOpen.value) return
    pageSize.value = size
    applyCampaignFilters()
}
watch(
    () => props.campaignPage,
    (value) => {
        if (
            !value ||
            (listLoading.value && filterSignature() !== submittedFilters)
        )
            return
        syncingFilters = true
        campaignSearch.value = value.filters.q
        statusFilter.value = value.filters.status || 'all'
        campaignSort.value = value.filters.sort
        pageSize.value = value.per_page
        syncingFilters = false
    },
    { deep: true, flush: 'sync' }
)
watch(
    campaignSearch,
    () => {
        if (
            syncingFilters ||
            !props.campaignPage ||
            pendingCampaignRecovery.value
        )
            return
        if (searchTimer) clearTimeout(searchTimer)
        if (listLoading.value && filterSignature() !== submittedFilters)
            queuedSearch = true
        searchTimer = setTimeout(
            applyCampaignFilters,
            campaignSearch.value ? 350 : 0
        )
    },
    { flush: 'sync' }
)
watch(
    () => props.serverTime,
    (value) => {
        if (value) {
            serverOffset = Date.parse(value) - Date.now()
            currentTime.value = Date.now() + serverOffset
        }
    }
)
watch(
    [statusSnapshotExpired, editorOpen, busy, listLoading],
    ([expired, open, saving]) => {
        if (
            expired &&
            !open &&
            !saving &&
            !listLoading.value &&
            !listError.value
        )
            refreshCampaigns()
    }
)
onBeforeUnmount(() => {
    disposed = true
    if (searchTimer) clearTimeout(searchTimer)
})
const localTime = (value: string) => {
    const d = new Date(value)
    return new Date(d.getTime() - d.getTimezoneOffset() * 60000)
        .toISOString()
        .slice(0, 23)
}
function edit(c: any) {
    if (
        busy.value ||
        editorOpen.value ||
        listLoading.value ||
        filtersPending.value ||
        pendingCampaignRecovery.value
    )
        return
    if (props.campaignPage) {
        visitCampaignList(
            campaignListUrl(props.campaignPage.page) +
                '&selected_campaign=' +
                encodeURIComponent(c.id),
            (response) => {
                if (
                    response.props.selectedCampaignStatus !== 'available' ||
                    !response.props.selectedCampaign ||
                    response.props.selectedCampaign.id !== c.id
                ) {
                    listError.value =
                        'This campaign is no longer available. Refresh the list.'
                    return
                }
                editCampaignData(response.props.selectedCampaign)
            }
        )
        return
    }
    editCampaignData(c)
}
function editCampaignData(c: any) {
    if (busy.value || editorOpen.value) return
    draftRecovered.value = false
    Object.assign(sourceSchedule, {
        starts_at: c.starts_at,
        ends_at: c.ends_at,
        local_start: localTime(c.starts_at),
        local_end: localTime(c.ends_at),
    })
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
    for (const bucket of ['include', 'exclude'] as TargetBucket[]) {
        activeTargetKind[bucket] =
            targetKinds.find((kind) => form[targetKey(bucket, kind.key)].length)
                ?.key || 'categories'
    }
    editorOpen.value = true
}
function reset() {
    if (busy.value) return
    resetDraft()
}
function resetDraft() {
    draftRecovered.value = false
    Object.assign(sourceSchedule, {
        starts_at: '',
        ends_at: '',
        local_start: '',
        local_end: '',
    })
    Object.assign(form, freshCampaign())
    activeTargetKind.include = 'categories'
    activeTargetKind.exclude = 'categories'
    error.value = ''
    targetSearch.value = ''
    savedDraft.value = JSON.stringify(form)
    editorOpen.value = false
}

function scheduleTimestamp(key: 'starts_at' | 'ends_at') {
    const originalLocal =
        key === 'starts_at'
            ? sourceSchedule.local_start
            : sourceSchedule.local_end
    return sourceSchedule[key] && form[key] === originalLocal
        ? sourceSchedule[key]
        : new Date(form[key]).toISOString()
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
            starts_at: scheduleTimestamp('starts_at'),
            ends_at: scheduleTimestamp('ends_at'),
            vendor_id: form.vendor_id,
            version: form.version,
            return_q: props.campaignPage?.filters.q || '',
            return_status: props.campaignPage?.filters.status || '',
            return_sort: props.campaignPage?.filters.sort || 'newest',
            return_page: props.campaignPage?.page || 1,
            return_per_page: props.campaignPage?.per_page || 25,
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
const recoveryKey = `soko:campaign-draft:${props.actorId}`
useAdminDraftRecovery({
    key: recoveryKey,
    allowed: () => Boolean(props.actorId), // This route requires promotions.manage on the server.
    read: () =>
        pendingCampaignRecovery.value || {
            form: Object.fromEntries(
                Object.keys(freshCampaign()).map((key) => [
                    key,
                    (form as any)[key],
                ])
            ),
            baseline: savedDraft.value,
            schedule: { ...sourceSchedule },
        },
    dirty: () =>
        Boolean(pendingCampaignRecovery.value) ||
        (editorOpen.value && dirty.value),
    apply: restoreCampaignDraft,
})
function discardPendingRecovery() {
    if (listLoading.value) return
    pendingCampaignRecovery.value = null
    router.remember(null, recoveryKey)
    listError.value = ''
}
function retryCampaignRecovery() {
    if (pendingCampaignRecovery.value && !listLoading.value)
        restoreCampaignDraft(pendingCampaignRecovery.value)
}
function restoreCampaignDraft(saved: Record<string, any>) {
    const draft = saved.form as ReturnType<typeof freshCampaign> | undefined
    if (
        !draft ||
        typeof draft !== 'object' ||
        (!props.canManageAll && draft.vendor_id !== props.actorId)
    )
        return
    pendingCampaignRecovery.value = saved
    const restore = () => {
        for (const key of Object.keys(freshCampaign())) {
            if (Object.prototype.hasOwnProperty.call(draft, key))
                (form as any)[key] = JSON.parse(
                    JSON.stringify((draft as any)[key])
                )
        }
        if (typeof saved.baseline === 'string')
            savedDraft.value = saved.baseline
        const schedule = saved.schedule as Record<string, unknown> | undefined
        for (const key of Object.keys(sourceSchedule)) {
            if (typeof schedule?.[key] === 'string')
                (sourceSchedule as any)[key] = schedule[key]
        }
        draftRecovered.value = true
        editorOpen.value = true
        pendingCampaignRecovery.value = null
    }
    if (!draft.id) {
        restore()
        return
    }
    visitCampaignList(
        campaignListUrl(props.campaignPage?.page || 1) +
            '&selected_campaign=' +
            encodeURIComponent(draft.id),
        (response) => {
            const current = response.props.selectedCampaign
            if (
                response.props.selectedCampaignStatus !== 'available' ||
                !current ||
                current.id !== draft.id ||
                (!props.canManageAll && current.vendor_id !== props.actorId)
            ) {
                listError.value =
                    'The saved campaign draft is no longer accessible.'
                pendingCampaignRecovery.value = null
                router.remember(null, recoveryKey)
                return
            }
            // Preserve the recovered version. A concurrent change must conflict, not be overwritten.
            restore()
        },
        true
    )
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
                :disabled="
                    busy ||
                    listLoading ||
                    filtersPending ||
                    !!pendingCampaignRecovery
                "
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
                    >{{
                        props.campaignPage?.filtered_total ??
                        visibleCampaigns.length
                    }}
                    matching campaigns</span
                >
            </div>
            <form
                class="campaign-list-filters"
                @submit.prevent="applyCampaignFilters"
            >
                <label
                    >Search campaigns<input
                        v-model="campaignSearch"
                        :disabled="!!pendingCampaignRecovery || busy"
                        type="search"
                        placeholder="Campaign or merchant name" /></label
                ><label
                    >Status<select
                        v-model="statusFilter"
                        :disabled="
                            busy || listLoading || !!pendingCampaignRecovery
                        "
                        @change="applyCampaignFilters"
                    >
                        <option value="all">All statuses</option>
                        <option
                            v-for="state in [
                                'Active',
                                'Scheduled',
                                'Disabled',
                                'Ended',
                            ]"
                            :key="state"
                            :value="state.toLowerCase()"
                        >
                            {{ state }}
                        </option>
                    </select></label
                >
                <label
                    >Sort<select
                        v-model="campaignSort"
                        :disabled="
                            busy || listLoading || !!pendingCampaignRecovery
                        "
                        @change="applyCampaignFilters"
                    >
                        <option value="newest">Newest first</option>
                        <option value="oldest">Oldest first</option>
                    </select></label
                >
                <button
                    type="submit"
                    class="button button-outline"
                    :disabled="busy || listLoading || !!pendingCampaignRecovery"
                >
                    Search
                </button>
                <button
                    type="button"
                    class="button button-outline"
                    :disabled="busy || listLoading || !!pendingCampaignRecovery"
                    @click="refreshCampaigns"
                >
                    Refresh
                </button>
                <button
                    v-if="campaignSearch || statusFilter !== 'all'"
                    type="button"
                    class="button button-outline"
                    @click="clearCampaignFilters"
                >
                    Clear filters
                </button>
            </form>
            <div
                v-if="pendingCampaignRecovery"
                class="campaign-recovery"
                role="status"
            >
                <p>
                    Unsaved campaign draft found.
                    {{
                        listLoading
                            ? 'Checking access before restoring…'
                            : 'Restore it before starting another campaign.'
                    }}
                </p>
                <button
                    type="button"
                    class="button button-navy"
                    :disabled="listLoading"
                    @click="retryCampaignRecovery"
                >
                    Retry restoring draft
                </button>
                <button
                    type="button"
                    class="button button-outline"
                    :disabled="listLoading"
                    @click="discardPendingRecovery"
                >
                    Discard recovered draft
                </button>
            </div>
            <p v-if="listError" role="alert" class="form-error">
                {{ listError }}
            </p>
            <p
                v-if="listLoading || filtersPending || statusSnapshotExpired"
                role="status"
                class="commerce-small"
            >
                {{
                    listLoading
                        ? 'Updating campaigns…'
                        : statusSnapshotExpired
                          ? 'A schedule boundary passed. Refresh to update counts and matching rows.'
                          : 'Filters changed. Updating results…'
                }}
            </p>
            <AdminTable
                caption="Promotion campaigns"
                :columns="campaignColumns"
                resizable-columns
                resizable-height
                v-model:viewport-height="listHeight"
                :empty="!visibleCampaigns.length"
                :busy="busy || listLoading"
                :empty-text="
                    campaignCounts.total
                        ? 'No campaigns match these filters. Try another name or status.'
                        : 'No campaigns yet. Create one when you are ready; regular prices stay unchanged.'
                "
            >
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
                            :disabled="
                                busy ||
                                listLoading ||
                                filtersPending ||
                                !!pendingCampaignRecovery
                            "
                            :items="[{ id: 'edit', label: 'Edit campaign' }]"
                            @select="edit(c)"
                        />
                    </td>
                </tr>
                <template v-if="campaignPage" #footer>
                    <AdminPagination
                        :page="campaignPage.page"
                        :per-page="campaignPage.per_page"
                        :total="campaignPage.filtered_total"
                        :last-page="campaignPage.last_page"
                        :from="campaignPage.from"
                        :to="campaignPage.to"
                        :busy="busy || listLoading || filtersPending"
                        @page="changeCampaignPage"
                        @per-page="changeCampaignSize"
                    />
                </template>
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
            <p v-if="draftRecovered" role="status" class="commerce-small">
                Unsaved draft restored. Review before saving; newer server
                changes will require a reload.
            </p>
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
                                step="0.001"
                                required
                            />
                        </div>
                        <div>
                            <label for="campaign-end">Ends</label
                            ><input
                                id="campaign-end"
                                v-model="form.ends_at"
                                type="datetime-local"
                                step="0.001"
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
                        >Search coverage choices</label
                    >
                    <input
                        id="campaign-product-search"
                        v-model="targetSearch"
                        type="search"
                        placeholder="Category, subcategory or product"
                    />
                    <section
                        v-for="bucket in targetPanels"
                        :key="bucket"
                        class="target-picker"
                        :aria-labelledby="`target-heading-${bucket}`"
                    >
                        <div class="target-heading">
                            <h3 :id="`target-heading-${bucket}`">
                                {{
                                    bucket === 'include'
                                        ? 'Include in campaign'
                                        : 'Exceptions · keep regular price'
                                }}
                            </h3>
                            <span role="status"
                                >{{
                                    selectedTargets(bucket).length
                                }}
                                selected</span
                            >
                        </div>
                        <p v-if="bucket === 'exclude'" class="commerce-small">
                            Exceptions override inclusions. Check an option to
                            exclude it.
                        </p>
                        <p v-else class="commerce-small">
                            Check any category, subcategory or individual
                            product to include it.
                        </p>
                        <div
                            class="target-kind-switch"
                            role="group"
                            :aria-label="`${bucket === 'include' ? 'Include' : 'Exception'} selection type`"
                        >
                            <button
                                v-for="kind in targetKinds"
                                :key="kind.key"
                                type="button"
                                :aria-pressed="
                                    activeTargetKind[bucket] === kind.key
                                "
                                :disabled="busy"
                                @click="activeTargetKind[bucket] = kind.key"
                            >
                                {{ kind.label }}
                                <span>{{
                                    form[targetKey(bucket, kind.key)].length
                                }}</span>
                            </button>
                        </div>
                        <fieldset class="target-choice-list">
                            <legend class="target-list-label">
                                {{
                                    targetKinds.find(
                                        (kind) =>
                                            kind.key ===
                                            activeTargetKind[bucket]
                                    )?.label
                                }}
                                ·
                                {{
                                    filteredTargetOptions(
                                        activeTargetKind[bucket]
                                    ).length
                                }}
                                choices
                            </legend>
                            <label
                                v-for="option in filteredTargetOptions(
                                    activeTargetKind[bucket]
                                )"
                                :key="option.value"
                                class="target-choice"
                            >
                                <input
                                    v-model="
                                        form[
                                            targetKey(
                                                bucket,
                                                activeTargetKind[bucket]
                                            )
                                        ]
                                    "
                                    type="checkbox"
                                    :value="option.value"
                                    :disabled="busy"
                                />
                                <span
                                    ><strong>{{ option.label }}</strong
                                    ><small>{{ option.context }}</small></span
                                >
                            </label>
                            <p
                                v-if="
                                    !filteredTargetOptions(
                                        activeTargetKind[bucket]
                                    ).length
                                "
                                class="target-empty"
                            >
                                No matching choices. Try another search or
                                selection type.
                            </p>
                        </fieldset>
                        <div
                            v-if="selectedTargets(bucket).length"
                            class="target-selected"
                            :aria-label="`${bucket === 'include' ? 'Included' : 'Excluded'} selections`"
                        >
                            <button
                                v-for="item in selectedTargets(bucket)"
                                :key="item.key + item.value"
                                type="button"
                                :disabled="busy"
                                :aria-label="`Remove ${bucket === 'include' ? 'inclusion' : 'exception'}: ${item.kind}, ${item.label}`"
                                @click="removeTarget(item.key, item.value)"
                            >
                                <span
                                    ><small>{{ item.kind }}</small
                                    >{{ item.label
                                    }}<small v-if="item.unavailable"
                                        >Not in current merchant catalog</small
                                    ></span
                                ><span aria-hidden="true">×</span>
                            </button>
                        </div>
                        <p v-else class="target-empty">
                            {{
                                bucket === 'include'
                                    ? 'Nothing selected yet.'
                                    : 'No exceptions. All included products can receive this discount.'
                            }}
                        </p>
                    </section>
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
.target-picker {
    margin-top: 1.25rem;
    min-width: 0;
}
.target-heading {
    display: flex;
    gap: 0.75rem;
    justify-content: space-between;
    align-items: baseline;
    flex-wrap: wrap;
}
.target-heading h3 {
    font-size: 1rem;
    margin: 0;
}
.target-heading > span {
    color: var(--soko-muted);
    font-size: 0.8rem;
}
.target-kind-switch {
    display: flex;
    flex-wrap: wrap;
    gap: 0.35rem;
    margin: 0.75rem 0;
}
.target-kind-switch button {
    min-height: 44px;
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    background: var(--soko-surface);
    color: var(--soko-ink);
}
.target-kind-switch button[aria-pressed='true'] {
    background: var(--soko-navy);
    color: var(--soko-surface);
    border-color: var(--soko-navy);
}
.target-kind-switch button span {
    margin-left: 0.35rem;
    opacity: 0.75;
}
fieldset.target-choice-list {
    padding: 0.25rem 0.5rem;
    max-height: 220px;
    overflow-y: auto;
    overscroll-behavior: contain;
    border-color: var(--soko-line);
}
.target-list-label {
    padding: 0 0.3rem;
    font-size: 0.75rem;
    color: var(--soko-muted);
}
.target-choice-list .target-choice {
    display: flex;
    align-items: flex-start;
    gap: 0.7rem;
    min-height: 48px;
    padding: 0.65rem 0.35rem;
    border-bottom: 1px solid var(--soko-line);
    cursor: pointer;
    margin: 0;
}
.target-choice:last-child {
    border-bottom: 0;
}
.target-choice input {
    flex: 0 0 auto;
    margin-top: 0.2rem;
    width: 18px;
    height: 18px;
    accent-color: var(--soko-navy);
}
.target-choice > span {
    display: grid;
    gap: 0.2rem;
    min-width: 0;
    overflow-wrap: anywhere;
}
.target-choice strong {
    font-size: 0.875rem;
    font-weight: 600;
}
.target-choice small,
.target-selected small {
    font-size: 0.72rem;
    color: var(--soko-muted);
}
.target-selected {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-top: 0.65rem;
}
.target-selected button {
    display: flex;
    gap: 0.75rem;
    align-items: center;
    min-height: 44px;
    max-width: 100%;
    padding: 0.4rem 0.65rem;
    background: var(--soko-canvas);
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    color: var(--soko-ink);
    text-align: left;
}
.target-selected button > span:first-child {
    display: grid;
    min-width: 0;
    overflow-wrap: anywhere;
    font-size: 0.8rem;
}
.target-empty {
    color: var(--soko-muted);
    font-size: 0.8rem;
    margin: 0.65rem 0;
}
.target-picker button:focus-visible,
.target-choice:focus-within {
    outline: 2px solid var(--soko-amber);
    outline-offset: 2px;
}
</style>
