import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import * as Vue from 'vue'
import { safeProductDraft } from '../../resources/js/design/productDraft.ts'
const source = readFileSync('resources/js/Pages/Dashboard/Index.vue', 'utf8')
const meta = (q = '', page = 1, status = '') => ({ page, per_page: 25, total: 101, filtered_total: 76, last_page: 4, from: (page-1)*25+1, to: Math.min(page*25,76), filters: { q, status, sort: 'newest' } })
function setup(extra = {}) {
    const gets = [], posts = [], timers = [], fetches = [], cleanups = []
    const props = Vue.reactive({ currentUser: { id: 'm', role: 'merchant' }, stats: {}, products: [], orders: [{ id: 'o', order_number: 'SOKO-1', status: 'pending', total_amount: 20, items: [{ id: 'line', quantity: 1, unit_price: 20 }] }], orderPage: meta(), orderVisibility: 'supplier', recordCapabilities: { orders: { o: { cancel: false } } }, capabilities: { readOrders: true, readProducts: true }, ...extra })
    const page = Vue.reactive({ url: '/dashboard?section=orders', props: {} })
    let script = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
    script += '\n globalThis.h = { orderFilters, orderFiltersPending, orderListUrl, applyOrderFilters, resetOrderFilters, changeOrderPage, changeOrderPageSize, sortOrderColumn, orderListLoading, orderListError, visibleRecords, orderColumns, openRecord, detailRecord, detailKind, pendingOperation, refreshSelectedOrder, selectedOrderVisibility, operationActions, confirmOperation, formError, busy };'
    const context = { exports: {}, defineProps: () => props, URL, URLSearchParams, AbortController,
        setTimeout: (callback, delay) => { timers.push({ callback, delay, cleared: false }); return timers.length }, clearTimeout: id => { if(timers[id-1]) timers[id-1].cleared = true },
        fetch: (url, options) => new Promise(resolve => fetches.push({ url, options, resolve })),
        require: name => name === 'vue' ? { ...Vue, onBeforeUnmount: fn => cleanups.push(fn) } : name === '@inertiajs/vue3' ? { usePage: () => page, router: { get: (...args) => gets.push(args), post: (...args) => posts.push(args), restore: () => undefined, remember: () => {} } } : name.includes('productDraft') ? { safeProductDraft } : name.includes('useCurrency') ? { useCurrency: () => ({ formatPrice: String }) } : {},
    }
    vm.runInNewContext(ts.transpileModule(script, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.CommonJS } }).outputText, context)
    return { ...context.h, props, page, gets, posts, timers, fetches, cleanup: () => cleanups.forEach(fn => fn()) }
}
const flush = async () => { await Vue.nextTick(); await new Promise(resolve => setImmediate(resolve)) }
const respond = (request, record, cancel = false, visibility = 'supplier') => request.resolve({ ok: true, json: async () => ({ props: { selectedOrder: record, selectedOrderStatus: record ? 'available' : 'unavailable', orderVisibility: visibility, recordCapabilities: { orders: { o: { cancel } } } } }) })
test('orders use server results/totals and encode search/status/sort without local filtering', () => {
    const h = setup({ orderPage: meta('ÉCOLE', 3, 'pending') })
    assert.equal(h.visibleRecords.value.length, 1)
    h.orderFilters.q = 'Équipement & buyer'
    h.orderFilters.sort = 'oldest'
    h.applyOrderFilters()
    const query = new URL(h.gets[0][0], 'https://example.com').searchParams
    assert.equal(query.get('section'), 'orders'); assert.equal(query.get('q'), 'Équipement & buyer'); assert.equal(query.get('status'), 'pending'); assert.equal(query.get('sort'), 'oldest'); assert.equal(query.get('page'), '1')
    assert.equal(h.orderColumns.value.find(c => c.key === 'total').label, 'Visible total')
    assert.equal(h.operationActions('orders', h.props.orders[0]).length, 1)
    h.cleanup()
})
test('debounced order search serializes latest query and immediate clear before first response', async () => {
    const h = setup()
    h.orderFilters.q = 'first'; await flush(); assert.equal(h.timers.at(-1).delay, 350); h.timers.at(-1).callback()
    h.orderFilters.q = ''; await flush(); assert.equal(h.timers.at(-1).delay, 0)
    h.props.orderPage = meta('first'); assert.equal(h.orderFilters.q, '')
    h.gets[0][2].onFinish(); await flush()
    assert.equal(h.gets.length, 2); assert.equal(new URL(h.gets[1][0], 'https://example.com').searchParams.has('q'), false)
    h.cleanup()
})
test('server sort and page size reset page; unapplied inputs block page navigation', () => {
    const h = setup({ orderPage: meta('buyer', 3) })
    h.orderFilters.q = 'new'; h.changeOrderPage(4); assert.equal(h.gets.length, 0)
    h.orderFilters.q = 'buyer'; h.sortOrderColumn({ key: 'buyer', direction: 'asc' }); assert.equal(h.gets.length, 0)
    h.sortOrderColumn({ key: 'created_at', direction: 'asc' }); assert.match(h.gets[0][0], /sort=oldest/); assert.match(h.gets[0][0], /page=1/)
    h.gets[0][2].onFinish(); h.changeOrderPageSize(50); assert.match(h.gets[1][0], /per_page=50/)
    h.cleanup()
})
test('selected order refresh uses off-page lookup and preserves review until success, then requires confirmation again', async () => {
    const h = setup({ orderVisibility: 'full', recordCapabilities: { orders: { o: { cancel: true } } } })
    h.openRecord('orders', h.props.orders[0], 'cancel'); assert.equal(h.posts.length, 0)
    const loading = h.refreshSelectedOrder(); assert.match(h.fetches[0].url, /selected_order=o/)
    assert.equal(h.pendingOperation.value, 'cancel')
    respond(h.fetches[0], { id: 'o', status: 'pending', total_amount: 20, items: [] }, true, 'full'); await loading
    assert.equal(h.pendingOperation.value, ''); assert.equal(h.selectedOrderVisibility.value, 'full')
    h.confirmOperation(); assert.equal(h.posts.length, 0)
    h.cleanup()
})
test('403 and unavailable selected-order refresh revoke cached cancellation capability', async () => {
    for (const status of ['403', 'unavailable']) {
        const h = setup({ orderVisibility: 'full', recordCapabilities: { orders: { o: { cancel: true } } } })
        h.openRecord('orders', h.props.orders[0], 'cancel')
        const loading = h.refreshSelectedOrder()
        if (status === '403') h.fetches[0].resolve({ ok: false, status: 403 })
        else respond(h.fetches[0], null, true)
        await loading
        assert.equal(h.detailRecord.value.id, 'o', 'keep the reviewed snapshot rather than silently switching records')
        assert.equal(h.operationActions('orders', h.detailRecord.value).length, 1)
        h.confirmOperation(); assert.equal(h.posts.length, 0)
        assert.ok(h.formError.value)
        h.cleanup()
    }
})
test('leaving Orders during selected refresh aborts and removes stale drawer without reopening it', async () => {
    const h = setup()
    h.openRecord('orders', h.props.orders[0]); const loading = h.refreshSelectedOrder()
    h.page.url = '/dashboard?section=products'; await flush()
    assert.equal(h.fetches[0].options.signal.aborted, true)
    assert.equal(h.detailRecord.value, null)
    respond(h.fetches[0], { id: 'o', status: 'pending', items: [] })
    await loading
    assert.equal(h.detailRecord.value, null); assert.equal(h.pendingOperation.value, ''); assert.equal(h.busy.value, false)
    h.cleanup()
})
test('confirmed cancellation returns to captured order query/page after static redirect', async () => {
    const h = setup({ orderPage: meta('buyer', 3, 'pending'), orderVisibility: 'full', recordCapabilities: { orders: { o: { cancel: true } } } })
    h.openRecord('orders', h.props.orders[0], 'cancel'); h.confirmOperation()
    assert.equal(h.posts[0][0], '/dashboard/orders/o/cancel')
    h.props.orderPage = meta()
    h.posts[0][2].onSuccess({ props: { flash: {} } }); h.posts[0][2].onFinish(); await flush()
    const query = new URL(h.gets[0][0], 'https://example.com').searchParams
    assert.equal(query.get('q'), 'buyer'); assert.equal(query.get('page'), '3'); assert.equal(query.get('status'), 'pending'); assert.equal(h.detailRecord.value, null)
    h.cleanup()
})

test('reviewing cancellation after off-page refresh retains fresh grants and the original return filters', async () => {
    const h = setup({ orderPage: meta('buyer', 3, 'pending'), orderVisibility: 'full', recordCapabilities: { orders: { o: { cancel: true } } } })
    h.openRecord('orders', h.props.orders[0], 'cancel')
    // Static mutation error redirect has changed page props; selected record remains reviewed.
    h.props.orderPage = meta(); h.props.recordCapabilities.orders = {}
    const check = h.refreshSelectedOrder()
    respond(h.fetches[0], { id: 'o', status: 'pending', items: [] }, true, 'full')
    await check
    h.openRecord('orders', h.detailRecord.value, 'cancel')
    assert.equal(h.pendingOperation.value, 'cancel')
    assert.equal(h.selectedOrderVisibility.value, 'full')
    h.confirmOperation()
    assert.equal(h.posts.length, 1)
    h.posts[0][2].onSuccess({ props: { flash: {} } }); h.posts[0][2].onFinish(); await flush()
    assert.equal(new URL(h.gets[0][0], 'https://example.com').searchParams.get('q'), 'buyer')
    assert.equal(new URL(h.gets[0][0], 'https://example.com').searchParams.get('page'), '3')
    h.cleanup()
})
