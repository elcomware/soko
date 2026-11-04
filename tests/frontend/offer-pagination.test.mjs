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
    const props = Vue.reactive({ currentUser: { id: 'm', role: 'merchant' }, stats: {}, products: [], offers: [{ id: 'o', offer_number: 'SOKO-1', status: 'pending', total_amount: 20, items: [{ id: 'line', quantity: 1, unit_price: 20 }] }], offerPage: meta(), offerVisibility: 'supplier', recordCapabilities: { offers: { o: { decide: false } } }, capabilities: { readOffers: true, readProducts: true }, ...extra })
    const page = Vue.reactive({ url: '/dashboard?section=offers', props: {} })
    let script = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
    script += '\n globalThis.h = { offerFilters, offerFiltersPending, offerListUrl, applyOfferFilters, resetOfferFilters, changeOfferPage, changeOfferPageSize, sortOfferColumn, offerListLoading, offerListError, visibleRecords, offerColumns, openRecord, detailRecord, detailKind, pendingOperation, refreshSelectedOffer, selectedOfferVisibility, operationActions, confirmOperation, formError, busy };'
    const context = { exports: {}, defineProps: () => props, URL, URLSearchParams, AbortController,
        setTimeout: (callback, delay) => { timers.push({ callback, delay, cleared: false }); return timers.length }, clearTimeout: id => { if(timers[id-1]) timers[id-1].cleared = true },
        fetch: (url, options) => new Promise(resolve => fetches.push({ url, options, resolve })),
        require: name => name === 'vue' ? { ...Vue, onBeforeUnmount: fn => cleanups.push(fn) } : name === '@inertiajs/vue3' ? { usePage: () => page, router: { get: (...args) => gets.push(args), post: (...args) => posts.push(args), restore: () => undefined, remember: () => {} } } : name.includes('productDraft') ? { safeProductDraft } : name.includes('useCurrency') ? { useCurrency: () => ({ formatPrice: String }) } : {},
    }
    vm.runInNewContext(ts.transpileModule(script, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.CommonJS } }).outputText, context)
    return { ...context.h, props, page, gets, posts, timers, fetches, cleanup: () => cleanups.forEach(fn => fn()) }
}
const flush = async () => { await Vue.nextTick(); await new Promise(resolve => setImmediate(resolve)) }
const respond = (request, record, decide = false, visibility = 'supplier') => request.resolve({ ok: true, json: async () => ({ props: { selectedOffer: record, selectedOfferStatus: record ? 'available' : 'unavailable', offerVisibility: visibility, recordCapabilities: { offers: { o: { decide } } } } }) })
test('offers use server results/totals and encode search/status/sort without local filtering', () => {
    const h = setup({ offerPage: meta('ÉCOLE', 3, 'pending') })
    assert.equal(h.visibleRecords.value.length, 1)
    h.offerFilters.q = 'Équipement & buyer'
    h.offerFilters.sort = 'oldest'
    h.applyOfferFilters()
    const query = new URL(h.gets[0][0], 'https://example.com').searchParams
    assert.equal(query.get('section'), 'offers'); assert.equal(query.get('q'), 'Équipement & buyer'); assert.equal(query.get('status'), 'pending'); assert.equal(query.get('sort'), 'oldest'); assert.equal(query.get('page'), '1')
    assert.equal(h.offerColumns.value.find(c => c.key === 'created_at').sortable, true)
    assert.equal(h.operationActions('offers', h.props.offers[0]).length, 1)
    h.cleanup()
})
test('debounced offer search serializes latest query and immediate clear before first response', async () => {
    const h = setup()
    h.offerFilters.q = 'first'; await flush(); assert.equal(h.timers.at(-1).delay, 350); h.timers.at(-1).callback()
    h.offerFilters.q = ''; await flush(); assert.equal(h.timers.at(-1).delay, 0)
    h.props.offerPage = meta('first'); assert.equal(h.offerFilters.q, '')
    h.gets[0][2].onFinish(); await flush()
    assert.equal(h.gets.length, 2); assert.equal(new URL(h.gets[1][0], 'https://example.com').searchParams.has('q'), false)
    h.cleanup()
})
test('server sort and page size reset page; unapplied inputs block page navigation', () => {
    const h = setup({ offerPage: meta('buyer', 3) })
    h.offerFilters.q = 'new'; h.changeOfferPage(4); assert.equal(h.gets.length, 0)
    h.offerFilters.q = 'buyer'; h.sortOfferColumn({ key: 'buyer', direction: 'asc' }); assert.equal(h.gets.length, 0)
    h.sortOfferColumn({ key: 'created_at', direction: 'asc' }); assert.match(h.gets[0][0], /sort=oldest/); assert.match(h.gets[0][0], /page=1/)
    h.gets[0][2].onFinish(); h.changeOfferPageSize(50); assert.match(h.gets[1][0], /per_page=50/)
    h.cleanup()
})
test('selected offer refresh uses off-page lookup and preserves review until success, then requires confirmation again', async () => {
    const h = setup({ offerVisibility: 'full', recordCapabilities: { offers: { o: { decide: true } } } })
    h.openRecord('offers', h.props.offers[0], 'accepted'); assert.equal(h.posts.length, 0)
    const loading = h.refreshSelectedOffer(); assert.match(h.fetches[0].url, /selected_offer=o/)
    assert.equal(h.pendingOperation.value, 'accepted')
    respond(h.fetches[0], { id: 'o', status: 'pending', total_amount: 20, items: [] }, true, 'full'); await loading
    assert.equal(h.pendingOperation.value, ''); assert.equal(h.selectedOfferVisibility.value, 'full')
    h.confirmOperation(); assert.equal(h.posts.length, 0)
    h.cleanup()
})
test('403 and unavailable selected-offer refresh revoke cached offer decision capability', async () => {
    for (const status of ['403', 'unavailable']) {
        const h = setup({ offerVisibility: 'full', recordCapabilities: { offers: { o: { decide: true } } } })
        h.openRecord('offers', h.props.offers[0], 'accepted')
        const loading = h.refreshSelectedOffer()
        if (status === '403') h.fetches[0].resolve({ ok: false, status: 403 })
        else respond(h.fetches[0], null, true)
        await loading
        assert.equal(h.detailRecord.value.id, 'o', 'keep the reviewed snapshot rather than silently switching records')
        assert.equal(h.operationActions('offers', h.detailRecord.value).length, 1)
        h.confirmOperation(); assert.equal(h.posts.length, 0)
        assert.ok(h.formError.value)
        h.cleanup()
    }
})
test('leaving Offers during selected refresh aborts and removes stale drawer without reopening it', async () => {
    const h = setup()
    h.openRecord('offers', h.props.offers[0]); const loading = h.refreshSelectedOffer()
    h.page.url = '/dashboard?section=products'; await flush()
    assert.equal(h.fetches[0].options.signal.aborted, true)
    assert.equal(h.detailRecord.value, null)
    respond(h.fetches[0], { id: 'o', status: 'pending', items: [] })
    await loading
    assert.equal(h.detailRecord.value, null); assert.equal(h.pendingOperation.value, ''); assert.equal(h.busy.value, false)
    h.cleanup()
})
test('confirmed offer decision returns to captured offer query/page after static redirect', async () => {
    const h = setup({ offerPage: meta('buyer', 3, 'pending'), offerVisibility: 'full', recordCapabilities: { offers: { o: { decide: true } } } })
    h.openRecord('offers', h.props.offers[0], 'accepted'); h.confirmOperation()
    assert.equal(h.posts[0][0], '/dashboard/offers/o/decision?section=offers')
    h.props.offerPage = meta()
    h.posts[0][2].onSuccess({ props: { flash: {} } }); h.posts[0][2].onFinish(); await flush()
    const query = new URL(h.gets[0][0], 'https://example.com').searchParams
    assert.equal(query.get('q'), 'buyer'); assert.equal(query.get('page'), '3'); assert.equal(query.get('status'), 'pending'); assert.equal(h.detailRecord.value, null)
    h.cleanup()
})

test('reviewing offer decision after off-page refresh retains fresh decision grants and the original return filters', async () => {
    const h = setup({ offerPage: meta('buyer', 3, 'pending'), offerVisibility: 'full', recordCapabilities: { offers: { o: { decide: true } } } })
    h.openRecord('offers', h.props.offers[0], 'accepted')
    // Static mutation error redirect has changed page props; selected record remains reviewed.
    h.props.offerPage = meta(); h.props.recordCapabilities.offers = {}
    const check = h.refreshSelectedOffer()
    respond(h.fetches[0], { id: 'o', status: 'pending', items: [] }, true, 'full')
    await check
    h.openRecord('offers', h.detailRecord.value, 'accepted')
    assert.equal(h.pendingOperation.value, 'accepted')
    assert.equal(h.selectedOfferVisibility.value, 'full')
    h.confirmOperation()
    assert.equal(h.posts.length, 1)
    h.posts[0][2].onSuccess({ props: { flash: {} } }); h.posts[0][2].onFinish(); await flush()
    assert.equal(new URL(h.gets[0][0], 'https://example.com').searchParams.get('q'), 'buyer')
    assert.equal(new URL(h.gets[0][0], 'https://example.com').searchParams.get('page'), '3')
    h.cleanup()
})

test('offer decisions require explicit review, preserve review on conflict and recheck current grant', async()=>{
 const h=setup({recordCapabilities:{offers:{o:{decide:true}}}})
 h.openRecord('offers',h.props.offers[0],'declined')
 assert.equal(h.posts.length,0);assert.equal(h.pendingOperation.value,'declined')
 h.confirmOperation();assert.equal(h.posts[0][1].status,'declined')
 h.posts[0][2].onError({error:'This offer is no longer pending.'});h.posts[0][2].onFinish()
 assert.equal(h.detailRecord.value.id,'o');assert.equal(h.pendingOperation.value,'declined');assert.ok(h.formError.value)
 const loading=h.refreshSelectedOffer()
 respond(h.fetches[0],{id:'o',status:'accepted',offered_price:20,quantity:1},false)
 await loading
 assert.equal(h.pendingOperation.value,'');assert.equal(h.operationActions('offers',h.detailRecord.value).length,1)
 h.openRecord('offers',h.detailRecord.value,'declined');h.confirmOperation();assert.equal(h.posts.length,1)
 h.cleanup()
})
