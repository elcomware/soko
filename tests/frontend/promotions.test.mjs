import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import { reactive, computed, ref, watch, nextTick } from 'vue'
function campaignHarness(
    canManageAll,
    products = [],
    campaigns = [],
    extra = {}
) {
    const source = readFileSync(
        new URL(
            '../../resources/js/Pages/Dashboard/Promotions.vue',
            import.meta.url
        ),
        'utf8'
    )
    let script = source
        .match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
        .replace(/^import[\s\S]*?(?:from\s*)?['"][^'"]+['"]\s*$/gm, '')
    script +=
        '\n globalThis.harness = { form, edit, reset, resetDraft, busy, error, eligibleProducts, merchants, selectableProducts, targetSearch, targetPanels, targetKey, targetOptions, filteredTargetOptions, selectedTargets, removeTarget, activeTargetKind, editorOpen, dirty, createCampaign, closeEditor, save, currentTime, status, visibleCampaigns, campaignSearch, statusFilter, listLoading, listError, filtersPending, campaignListUrl, applyCampaignFilters, refreshCampaigns, changeCampaignPage, changeCampaignSize, statusSnapshotExpired, campaignCounts, sourceSchedule, scheduleTimestamp, editCampaignData, draftRecovered, pendingCampaignRecovery, retryCampaignRecovery, discardPendingRecovery, loadTargets, targetPage, targetPageSize, changeTargetPage, changeTargetPageSize, currentTargets, targetChoices, targetLoading, targetError, coverageTotal, campaignScope };'
    const mounts = [],
        unmounts = [],
        requests = [],
        posts = [],
        timers = []
    const props = reactive({
        canManageAll,
        actorId: 'merchant-own',
        products,
        campaigns,
        ...extra,
    })
    const clockTimers = [],
        recoveries = []
    const fetches = []
    const context = {
        AbortController,
        fetch: (url, options) =>
            new Promise((resolve, reject) =>
                fetches.push({ url, options, resolve, reject })
            ),
        onMounted: (fn) => {
            mounts.push(fn)
        },
        onBeforeUnmount: (fn) => {
            unmounts.push(fn)
        },
        setInterval: (callback, delay) => {
            clockTimers.push({ callback, delay, cleared: false })
            return clockTimers.length
        },
        clearInterval: (id) => {
            clockTimers[id - 1].cleared = true
        },
        useAdminDraftRecovery: (options) => recoveries.push(options),
        watch,
        nextTick,
        URLSearchParams,
        setTimeout: (callback, delay) => {
            timers.push({ callback, delay, cleared: false })
            return timers.length
        },
        clearTimeout: (id) => {
            if (timers[id - 1]) timers[id - 1].cleared = true
        },
        router: {
            remember: () => {},
            get: (url, data, options) => requests.push({ url, data, options }),
        },
        reactive,
        computed,
        ref,
        defineProps: () => props,
        document: { getElementById: () => null, querySelector: () => null },
        useMarketplaceForm: () => ({
            busy: ref(false),
            error: ref('old error'),
            post: (url, data, success) => posts.push({ url, data, success }),
            validate: () => true,
        }),
    }
    vm.runInNewContext(
        ts.transpileModule(script, {
            compilerOptions: {
                target: ts.ScriptTarget.ES2022,
                module: ts.ModuleKind.None,
            },
        }).outputText,
        context
    )
    return {
        ...context.harness,
        mount: () => mounts.forEach((fn) => fn()),
        unmount: () => unmounts.forEach((fn) => fn()),
        props,
        requests,
        posts,
        timers,
        clockTimers,
        fetches,
        recoveries,
    }
}
for (const canManageAll of [true, false])
    test(`campaign reset restores every default (all permission: ${canManageAll})`, () => {
        const h = campaignHarness(canManageAll)
        const defaults = JSON.stringify(h.form)
        const scope = {
            all: false,
            product_ids: ['old-product'],
            categories: ['Clothes'],
            sub_categories: [{ category: 'Clothes', sub_category: 'Shoes' }],
            exclude_product_ids: ['excluded'],
            exclude_categories: ['Clearance'],
            exclude_sub_categories: [
                { category: 'Clothes', sub_category: 'Coats' },
            ],
        }
        h.edit({
            id: 'old',
            name: 'Past campaign',
            kind: 'fixed',
            amount: 77,
            enabled: true,
            vendor_id: 'old-merchant',
            updated_at: 'version',
            starts_at: '2026-12-01T00:00:00Z',
            ends_at: '2026-12-31T00:00:00Z',
            scope,
        })
        h.form.product_ids.push('new-selection')
        assert.deepEqual(
            scope.product_ids,
            ['old-product'],
            'editing must not mutate campaign history'
        )
        h.error.value = 'Failed request'
        h.reset()
        assert.equal(JSON.stringify(h.form), defaults)
        assert.equal(h.error.value, '')
        assert.equal(h.form.vendor_id, canManageAll ? '' : 'merchant-own')
        h.form.categories.push('temporary')
        h.reset()
        assert.equal(
            h.form.categories.length,
            0,
            'reset must allocate fresh target arrays'
        )
    })

test('preview respects merchant, category, subcategory and exclusion scope without changing on search', () => {
    const products = [
        {
            id: 'a',
            title: 'Coat',
            category: 'Clothes',
            sub_category: 'Winter',
            vendor_id: 'one',
            vendor_name: 'One',
        },
        {
            id: 'b',
            title: 'Shirt',
            category: 'Clothes',
            sub_category: 'Summer',
            vendor_id: 'one',
            vendor_name: 'One',
        },
        {
            id: 'c',
            title: 'Other coat',
            category: 'Clothes',
            sub_category: 'Winter',
            vendor_id: 'two',
            vendor_name: 'Two',
        },
    ]
    const h = campaignHarness(true, products)
    assert.equal(h.merchants.value.length, 2)
    h.form.vendor_id = 'one'
    assert.equal(h.eligibleProducts.value.length, 2)
    h.form.exclude_sub_categories = [
        JSON.stringify({ category: 'Clothes', sub_category: 'Summer' }),
    ]
    assert.equal(h.eligibleProducts.value[0].id, 'a')
    assert.equal(h.eligibleProducts.value.length, 1)
    h.targetSearch.value = 'Shirt'
    assert.equal(h.selectableProducts.value.length, 1)
    assert.equal(
        h.eligibleProducts.value[0].id,
        'a',
        'search filters choices, not campaign coverage'
    )
    h.form.all = false
    assert.equal(h.eligibleProducts.value.length, 0)
    h.form.product_ids = ['a', 'c']
    assert.equal(
        h.eligibleProducts.value.length,
        1,
        'explicit products cannot escape merchant boundary'
    )
    h.form.exclude_categories = ['Clothes']
    assert.equal(h.eligibleProducts.value.length, 0)
})

test('in-flight campaign saves cannot switch or cancel the current draft', () => {
    const h = campaignHarness(true)
    h.form.name = 'Pending campaign'
    h.busy.value = true
    h.edit({ id: 'other' })
    h.reset()
    assert.equal(h.form.name, 'Pending campaign')
    h.resetDraft()
    assert.equal(
        h.form.name,
        '',
        'confirmed success must reset even before onFinish clears busy'
    )
})

test('campaign drawer starts clean, tracks edits, and success reset closes it', () => {
    const h = campaignHarness(true)
    assert.equal(h.editorOpen.value, false)
    h.createCampaign()
    assert.equal(h.editorOpen.value, true)
    assert.equal(h.dirty.value, false)
    h.form.name = 'Holiday'
    assert.equal(h.dirty.value, true)
    h.createCampaign()
    assert.equal(
        h.form.name,
        'Holiday',
        'another create request must not silently reset an open draft'
    )
    h.busy.value = true
    h.closeEditor()
    assert.equal(h.editorOpen.value, true)
    h.resetDraft()
    assert.equal(h.editorOpen.value, false)
    assert.equal(h.dirty.value, false)
})
test('reactive campaign clock updates status filters at start and exclusive end', () => {
    const c = {
        name: 'Holiday',
        vendor_id: '',
        kind: 'percentage',
        enabled: true,
        starts_at: '2030-01-01T00:00:00Z',
        ends_at: '2030-01-02T00:00:00Z',
    }
    const h = campaignHarness(true, [], [c])
    h.currentTime.value = Date.parse(c.starts_at) - 1
    h.statusFilter.value = 'Scheduled'
    assert.equal(h.visibleCampaigns.value.length, 1)
    h.currentTime.value = Date.parse(c.starts_at)
    assert.equal(h.visibleCampaigns.value.length, 0)
    h.statusFilter.value = 'Active'
    assert.equal(h.visibleCampaigns.value.length, 1)
    h.currentTime.value = Date.parse(c.ends_at)
    assert.equal(h.status(c), 'Ended')
    assert.equal(h.visibleCampaigns.value.length, 0)
    h.statusFilter.value = 'Ended'
    h.campaignSearch.value = 'holiday'
    assert.equal(h.visibleCampaigns.value.length, 1)
    h.campaignSearch.value = 'other'
    assert.equal(h.visibleCampaigns.value.length, 0)
})

test('campaign status clock is bounded and cleaned up with page', () => {
    const h = campaignHarness(true)
    h.mount()
    assert.equal(h.clockTimers.length, 1)
    assert.equal(h.clockTimers[0].delay, 1000)
    h.unmount()
    assert.equal(h.clockTimers[0].cleared, true)
})

test('coverage checklists preserve selections across searching and expose hierarchy and removable exceptions', () => {
    const products = [
        {
            id: 'a',
            title: 'Rain coat',
            category: 'Clothes',
            sub_category: 'Coats',
            vendor_id: 'merchant-own',
            vendor_name: 'Acme',
        },
        {
            id: 'b',
            title: 'Work boot',
            category: 'Clothes',
            sub_category: 'Shoes',
            vendor_id: 'merchant-own',
        },
    ]
    const h = campaignHarness(false, products)
    h.form.all = false
    h.form.categories = ['Clothes']
    h.form.exclude_product_ids = ['a', 'archived']
    h.targetSearch.value = 'Shoes'
    assert.equal(h.filteredTargetOptions('product_ids').length, 1)
    assert.equal(h.filteredTargetOptions('product_ids')[0].value, 'b')
    assert.equal(
        h.filteredTargetOptions('sub_categories')[0].label,
        'Clothes / Shoes'
    )
    assert.equal(h.selectedTargets('exclude').length, 2)
    assert.equal(h.selectedTargets('exclude')[0].label, 'Rain coat')
    assert.equal(h.selectedTargets('exclude')[1].unavailable, true)
    assert.equal(h.eligibleProducts.value.length, 1)
    h.busy.value = true
    h.removeTarget('exclude_product_ids', 'a')
    assert.equal(h.form.exclude_product_ids.length, 2)
    h.busy.value = false
    h.removeTarget('exclude_product_ids', 'a')
    assert.equal(h.form.exclude_product_ids.length, 1)
    assert.equal(h.eligibleProducts.value.length, 2)
    assert.equal(
        h.form.exclude_product_ids[0],
        'archived',
        'unknown saved IDs remain until explicitly removed'
    )
    h.form.all = true
    assert.deepEqual([...h.targetPanels.value], ['exclude'])
    assert.equal(
        h.form.categories[0],
        'Clothes',
        'all-products mode does not destroy saved targeting'
    )
})
test('coverage uses explicit checkbox choices without platform multi-select gestures', () => {
    const source = readFileSync(
        new URL(
            '../../resources/js/Pages/Dashboard/Promotions.vue',
            import.meta.url
        ),
        'utf8'
    )
    assert.doesNotMatch(source, /Ctrl\/Cmd|\bmultiple\b/)
    assert.ok(
        source
            .replace(/\s+/g, '')
            .includes(
                'v-model="form[targetKey(bucket,activeTargetKind[bucket])]"'
            )
    )
    assert.ok(
        source
            .replace(/\s+/g, '')
            .includes(':aria-pressed="activeTargetKind[bucket]===kind.key"')
    )
    assert.match(source, /Remove \$\{bucket === 'include'/)
})

test('category choice counts use singular and plural product labels', () => {
    const h = campaignHarness(false, [
        {
            id: 'a',
            title: 'A',
            category: 'Machines',
            sub_category: '',
            vendor_id: 'merchant-own',
        },
        {
            id: 'b',
            title: 'B',
            category: 'Clothes',
            sub_category: '',
            vendor_id: 'merchant-own',
        },
        {
            id: 'c',
            title: 'C',
            category: 'Clothes',
            sub_category: '',
            vendor_id: 'merchant-own',
        },
    ])
    const options = h.targetOptions('categories')
    assert.equal(
        options.find((option) => option.value === 'Machines').context,
        '1 product'
    )
    assert.equal(
        options.find((option) => option.value === 'Clothes').context,
        '2 products'
    )
})

const campaignPage = (overrides = {}) => ({
    page: 2,
    per_page: 25,
    total: 80,
    filtered_total: 40,
    last_page: 2,
    from: 26,
    to: 40,
    filters: { q: '', status: '', sort: 'newest' },
    ...overrides,
})
test('server campaign results and counts stay authoritative with bounded paging', () => {
    const h = campaignHarness(
        true,
        [],
        [
            {
                id: 'server-match',
                name: 'Server matched merchant',
                enabled: false,
            },
        ],
        {
            campaignPage: campaignPage(),
            campaignCounts: {
                total: 80,
                active: 30,
                scheduled: 10,
                ended: 20,
                disabled: 20,
            },
        }
    )
    h.campaignSearch.value = 'pending query'
    assert.equal(
        h.visibleCampaigns.value.length,
        1,
        'do not locally refilter a server page'
    )
    assert.equal(h.campaignCounts.value.total, 80)
    h.changeCampaignPage(1)
    assert.equal(h.requests.length, 0, 'no paging with unapplied filters')
    h.applyCampaignFilters()
    assert.match(h.requests[0].url, /page=1/)
    assert.match(h.requests[0].url, /q=pending\+query/)
    h.createCampaign()
    assert.equal(
        h.editorOpen.value,
        false,
        'list request cannot race draft opening'
    )
    h.requests[0].options.onFinish()
    h.unmount()
})
test('debounced campaign search queues latest query and preserves typing across stale responses', async () => {
    const h = campaignHarness(true, [], [], { campaignPage: campaignPage() })
    h.campaignSearch.value = 'coat'
    assert.equal(h.timers.at(-1).delay, 350)
    h.timers.at(-1).callback()
    assert.equal(h.requests.length, 1)
    h.campaignSearch.value = 'coats'
    h.props.campaignPage = campaignPage({
        page: 1,
        filters: { q: 'coat', status: '', sort: 'newest' },
    })
    assert.equal(h.campaignSearch.value, 'coats')
    h.requests[0].options.onFinish()
    await nextTick()
    assert.equal(h.requests.length, 2)
    assert.match(h.requests[1].url, /q=coats/)
    h.requests[1].options.onFinish()
    h.campaignSearch.value = ''
    assert.equal(
        h.timers.at(-1).delay,
        0,
        'clearing search refreshes immediately'
    )
    h.unmount()
})
test('schedule boundary waits for draft closure and failed refresh remains visible without retry loop', async () => {
    const boundary = Date.now() + 10000
    const h = campaignHarness(true, [], [], {
        campaignPage: campaignPage(),
        nextStatusChange: new Date(boundary).toISOString(),
    })
    h.createCampaign()
    h.form.name = 'Unsaved campaign'
    h.currentTime.value = boundary
    await nextTick()
    assert.equal(h.requests.length, 0)
    assert.equal(h.form.name, 'Unsaved campaign')
    h.closeEditor()
    await nextTick()
    assert.equal(h.requests.length, 1)
    h.requests[0].options.onNetworkError()
    h.requests[0].options.onFinish()
    h.currentTime.value += 1000
    await nextTick()
    assert.equal(h.requests.length, 1)
    assert.match(h.listError.value, /Connection lost/)
    h.unmount()
})
test('editing refreshes selected campaign and unavailable records never open stale drafts', () => {
    const h = campaignHarness(true, [], [], { campaignPage: campaignPage() })
    h.edit({ id: 'old-row' })
    assert.match(h.requests[0].url, /selected_campaign=old-row/)
    h.requests[0].options.onSuccess({
        props: { selectedCampaignStatus: 'unavailable' },
    })
    assert.equal(h.editorOpen.value, false)
    assert.match(h.listError.value, /no longer available/)
    h.requests[0].options.onFinish()
    h.createCampaign()
    h.form.starts_at = '2026-12-01T00:00'
    h.form.ends_at = '2026-12-02T00:00'
    h.save({})
    assert.equal(h.posts[0].data.return_page, 2)
    assert.equal(h.posts[0].data.return_sort, 'newest')
    h.unmount()
})

function savedCampaign() {
    return {
        id: 'campaign-1',
        name: 'Scheduled',
        kind: 'percentage',
        amount: 10,
        enabled: true,
        vendor_id: 'merchant-own',
        updated_at: 'original-version',
        starts_at: '2026-12-01T12:30:47.123456789Z',
        ends_at: '2026-12-02T15:16:59.987654321Z',
        scope: {
            all: true,
            product_ids: [],
            categories: [],
            sub_categories: [],
            exclude_product_ids: ['removed-product'],
            exclude_categories: [],
            exclude_sub_categories: [],
        },
    }
}
test('editing unrelated campaign fields preserves original schedule precision byte-for-byte', () => {
    const h = campaignHarness(false)
    const campaign = savedCampaign()
    h.edit(campaign)
    h.form.name = 'Renamed'
    h.save({})
    assert.equal(h.posts[0].data.starts_at, campaign.starts_at)
    assert.equal(h.posts[0].data.ends_at, campaign.ends_at)
    h.form.starts_at = '2026-12-01T12:31:48.456'
    assert.equal(new Date(h.scheduleTimestamp('starts_at')).getSeconds(), 48)
    assert.equal(
        new Date(h.scheduleTimestamp('starts_at')).getMilliseconds(),
        456
    )
})
test('empty category and parentless subcategory cannot be selected but uncategorized products can', () => {
    const h = campaignHarness(false, [
        {
            id: 'uncategorized',
            title: 'Uncategorized item',
            category: '',
            sub_category: 'Invalid parent',
            vendor_id: 'merchant-own',
        },
    ])
    assert.equal(h.targetOptions('categories').length, 0)
    assert.equal(h.targetOptions('sub_categories').length, 0)
    assert.equal(h.targetOptions('product_ids').length, 1)
    h.form.categories = ['']
    assert.equal(
        h.selectedTargets('include').length,
        1,
        'existing scope is not silently changed'
    )
})
test('campaign history recovery verifies current ownership without refreshing the saved concurrency version', () => {
    const source = campaignHarness(false)
    source.edit(savedCampaign())
    source.form.name = 'Unsaved name'
    const snapshot = JSON.parse(JSON.stringify(source.recoveries[0].read()))
    const h = campaignHarness(false, [], [], { campaignPage: campaignPage() })
    assert.equal(h.recoveries[0].key, 'soko:campaign-draft:merchant-own')
    h.recoveries[0].apply(snapshot)
    assert.equal(h.editorOpen.value, false)
    h.requests[0].options.onSuccess({
        props: {
            selectedCampaignStatus: 'available',
            selectedCampaign: {
                ...savedCampaign(),
                updated_at: 'new-server-version',
            },
        },
    })
    assert.equal(h.form.name, 'Unsaved name')
    assert.equal(h.form.version, 'original-version')
    assert.equal(h.dirty.value, true)
    assert.equal(h.form.exclude_product_ids[0], 'removed-product')
    assert.equal(h.scheduleTimestamp('starts_at'), savedCampaign().starts_at)
    assert.equal(h.draftRecovered.value, true)
    const revoked = campaignHarness(false, [], [], {
        campaignPage: campaignPage(),
    })
    snapshot.form.vendor_id = 'another-merchant'
    revoked.recoveries[0].apply(snapshot)
    assert.equal(revoked.requests.length, 0)
    assert.equal(revoked.editorOpen.value, false)
})

test('failed history restore keeps pending draft and blocks competing actions until retry or explicit discard', () => {
    const source = campaignHarness(false)
    source.edit(savedCampaign())
    source.form.name = 'Important unsaved change'
    const snapshot = JSON.parse(JSON.stringify(source.recoveries[0].read()))
    const h = campaignHarness(false, [], [], { campaignPage: campaignPage() })
    h.recoveries[0].apply(snapshot)
    h.requests[0].options.onNetworkError()
    h.requests[0].options.onFinish()
    assert.equal(
        h.pendingCampaignRecovery.value.form.name,
        'Important unsaved change'
    )
    assert.equal(h.recoveries[0].dirty(), true)
    h.createCampaign()
    h.applyCampaignFilters()
    h.edit({ id: 'another' })
    assert.equal(h.requests.length, 1)
    assert.equal(h.editorOpen.value, false)
    h.retryCampaignRecovery()
    assert.equal(h.requests.length, 2)
    h.requests[1].options.onNetworkError()
    h.requests[1].options.onFinish()
    h.discardPendingRecovery()
    assert.equal(h.pendingCampaignRecovery.value, null)
    h.createCampaign()
    assert.equal(h.editorOpen.value, true)
})
test('selected record identity is checked and schedule expiry during a request refreshes after completion', async () => {
    const h = campaignHarness(true, [], [], {
        campaignPage: campaignPage(),
        nextStatusChange: new Date(Date.now() + 10000).toISOString(),
    })
    h.edit({ id: 'wanted' })
    h.requests[0].options.onSuccess({
        props: {
            selectedCampaignStatus: 'available',
            selectedCampaign: savedCampaign(),
        },
    })
    assert.equal(h.editorOpen.value, false)
    assert.match(h.listError.value, /no longer available/)
    h.requests[0].options.onFinish()
    h.refreshCampaigns()
    h.currentTime.value = Date.parse(h.props.nextStatusChange)
    await nextTick()
    assert.equal(h.requests.length, 2)
    h.requests[1].options.onFinish()
    await nextTick()
    assert.equal(h.requests.length, 3)
    h.unmount()
})

function targetResult(overrides = {}) {
    const product = {
        id: 'page-product',
        title: 'Page product',
        category: 'Tools',
        sub_category: 'Hand tools',
        vendor_id: 'merchant-own',
    }
    return {
        vendor_id: 'merchant-own',
        products: [product],
        productPage: {
            page: 1,
            per_page: 25,
            total: 70,
            filtered_total: 70,
            last_page: 3,
            from: 1,
            to: 25,
        },
        categories: [{ category: 'Tools', count: 70 }],
        subcategories: [
            { category: 'Tools', sub_category: 'Hand tools', count: 70 },
        ],
        selectedProducts: [],
        preview: { total: 70, products: [product] },
        ...overrides,
    }
}
function resolveTargets(request, data) {
    request.resolve({ ok: true, redirected: false, json: async () => data })
}
test('remote target pages keep whole-scope coverage and off-page selected labels independently', async () => {
    const h = campaignHarness(false, [], [], {
        targetingApi: true,
        merchants: [],
    })
    h.createCampaign()
    h.form.exclude_product_ids = ['off-page', 'deleted']
    const pending = h.loadTargets()
    const saved = {
        id: 'off-page',
        title: 'Chosen elsewhere',
        category: 'Tools',
        sub_category: '',
        vendor_id: 'merchant-own',
    }
    resolveTargets(
        h.fetches[0],
        targetResult({
            selectedProducts: [
                { id: 'off-page', status: 'available', product: saved },
                { id: 'deleted', status: 'unavailable', product: null },
            ],
            preview: { total: 68, products: [] },
        })
    )
    await pending
    assert.equal(h.coverageTotal.value, 68)
    assert.equal(h.targetOptions('product_ids').length, 1)
    assert.equal(h.selectedTargets('exclude')[0].label, 'Chosen elsewhere')
    assert.equal(h.selectedTargets('exclude')[1].unavailable, true)
    h.form.product_ids.push('page-product')
    assert.equal(
        h.coverageTotal.value,
        null,
        'changed scope never shows old count as current'
    )
    assert.equal(
        h.targetOptions('product_ids').length,
        1,
        'scope changes retain stable visible choices for multi-select'
    )
    assert.equal(h.selectedTargets('include')[0].label, 'Page product')
    h.targetSearch.value = 'no-result'
    const searched = h.loadTargets()
    resolveTargets(
        h.fetches[1],
        targetResult({
            products: [],
            productPage: {
                page: 1,
                per_page: 25,
                total: 70,
                filtered_total: 0,
                last_page: 1,
                from: 0,
                to: 0,
            },
            selectedProducts: [
                { id: 'off-page', status: 'available', product: saved },
            ],
            preview: { total: 68, products: [] },
        })
    )
    await searched
    assert.equal(h.targetOptions('product_ids').length, 0)
    assert.equal(
        h.coverageTotal.value,
        68,
        'zero search results do not mean zero campaign coverage'
    )
    assert.equal(h.form.exclude_product_ids.length, 2)
    h.unmount()
})
test('stale targeting replies cannot cross merchant, scope, page or closed editor', async () => {
    const h = campaignHarness(true, [], [], {
        targetingApi: true,
        merchants: [],
    })
    h.createCampaign()
    h.form.vendor_id = 'merchant-own'
    const first = h.loadTargets()
    h.form.vendor_id = 'different-owner'
    const second = h.loadTargets()
    resolveTargets(
        h.fetches[1],
        targetResult({
            vendor_id: 'different-owner',
            preview: { total: 4, products: [] },
        })
    )
    await second
    resolveTargets(
        h.fetches[0],
        targetResult({ preview: { total: 999, products: [] } })
    )
    await first
    assert.equal(h.coverageTotal.value, 4)
    assert.equal(h.fetches[0].options.signal.aborted, true)
    const third = h.loadTargets()
    h.form.exclude_categories.push('Tools')
    resolveTargets(
        h.fetches[2],
        targetResult({
            vendor_id: 'different-owner',
            preview: { total: 1000, products: [] },
        })
    )
    await third
    assert.equal(h.coverageTotal.value, null)
    const fourth = h.loadTargets()
    h.closeEditor()
    resolveTargets(h.fetches[3], targetResult({ vendor_id: 'different-owner' }))
    await fourth
    assert.equal(h.currentTargets.value, null)
    h.unmount()
})
test('target lookup failure retains choices in draft, hides coverage and allows retry', async () => {
    const h = campaignHarness(false, [], [], {
        targetingApi: true,
        merchants: [],
    })
    h.createCampaign()
    h.form.categories = ['Tools']
    const pending = h.loadTargets()
    h.fetches[0].reject(new Error('offline'))
    await pending
    assert.match(h.targetError.value, /selections are unchanged/)
    assert.equal(h.coverageTotal.value, null)
    assert.equal(h.form.categories[0], 'Tools')
    const retried = h.loadTargets()
    resolveTargets(h.fetches[1], targetResult())
    await retried
    assert.equal(h.coverageTotal.value, 70)
    h.busy.value = true
    h.changeTargetPage(2)
    h.changeTargetPageSize(50)
    assert.equal(h.targetPage.value, 1)
    assert.equal(h.targetPageSize.value, 25)
    h.unmount()
})
