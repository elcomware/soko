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
    const props = Vue.reactive({ currentUser: { id: 'm', role: 'merchant' }, stats: {}, products: [], rfqs: [{ id: 'o', rfq_number: 'SOKO-1', status: 'open', target_price: 20, items: [{ id: 'line', quantity: 1, unit_price: 20 }] }], rfqPage: meta(), rfqVisibility: 'supplier', recordCapabilities: { rfqs: { o: { close: false } } }, capabilities: { readRfqs: true, readProducts: true }, ...extra })
    const page = Vue.reactive({ url: '/dashboard?section=rfqs', props: {} })
    let script = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
    script += '\n globalThis.h = { rfqFilters, rfqFiltersPending, rfqListUrl, applyRfqFilters, resetRfqFilters, changeRfqPage, changeRfqPageSize, sortRfqColumn, rfqListLoading, rfqListError, visibleRecords, rfqColumns, openRecord, detailRecord, detailKind, pendingOperation, refreshSelectedRfq, selectedRfqVisibility, operationActions, confirmOperation, formError, busy, prepareOperation, pendingQuote, recordCan, quoteForm, quoteTarget, pendingQuoteDraft, recoverQuoteDraft, quoteRecoveryError };'
    const context = { exports: {}, defineProps: () => props, URL, URLSearchParams, AbortController,
        setTimeout: (callback, delay) => { timers.push({ callback, delay, cleared: false }); return timers.length }, clearTimeout: id => { if(timers[id-1]) timers[id-1].cleared = true },
        fetch: (url, options) => new Promise(resolve => fetches.push({ url, options, resolve })),
        require: name => name === 'vue' ? { ...Vue, onBeforeUnmount: fn => cleanups.push(fn) } : name === '@inertiajs/vue3' ? { usePage: () => page, router: { get: (...args) => gets.push(args), post: (...args) => posts.push(args), restore: () => undefined, remember: () => {} } } : name.includes('productDraft') ? { safeProductDraft } : name.includes('useCurrency') ? { useCurrency: () => ({ formatPrice: String }) } : {},
    }
    vm.runInNewContext(ts.transpileModule(script, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.CommonJS } }).outputText, context)
    return { ...context.h, props, page, gets, posts, timers, fetches, cleanup: () => cleanups.forEach(fn => fn()) }
}
const flush = async () => { await Vue.nextTick(); await new Promise(resolve => setImmediate(resolve)) }
const respond = (request, record, close = false, visibility = 'supplier') => request.resolve({ ok: true, json: async () => ({ props: { selectedRfq: record, selectedRfqStatus: record ? 'available' : 'unavailable', rfqVisibility: visibility, recordCapabilities: { rfqs: { o: { close } } } } }) })
test('rfqs use server results/totals and encode search/status/sort without local filtering', () => {
    const h = setup({ rfqPage: meta('ÉCOLE', 3, 'open') })
    assert.equal(h.visibleRecords.value.length, 1)
    h.rfqFilters.q = 'Équipement & buyer'
    h.rfqFilters.sort = 'oldest'
    h.applyRfqFilters()
    const query = new URL(h.gets[0][0], 'https://example.com').searchParams
    assert.equal(query.get('section'), 'rfqs'); assert.equal(query.get('q'), 'Équipement & buyer'); assert.equal(query.get('status'), 'open'); assert.equal(query.get('sort'), 'oldest'); assert.equal(query.get('page'), '1')
    assert.equal(h.rfqColumns.value.find(c => c.key === 'created_at').sortable, true)
    assert.equal(h.operationActions('rfqs', h.props.rfqs[0]).length, 1)
    h.cleanup()
})
test('debounced rfq search serializes latest query and immediate clear before first response', async () => {
    const h = setup()
    h.rfqFilters.q = 'first'; await flush(); assert.equal(h.timers.at(-1).delay, 350); h.timers.at(-1).callback()
    h.rfqFilters.q = ''; await flush(); assert.equal(h.timers.at(-1).delay, 0)
    h.props.rfqPage = meta('first'); assert.equal(h.rfqFilters.q, '')
    h.gets[0][2].onFinish(); await flush()
    assert.equal(h.gets.length, 2); assert.equal(new URL(h.gets[1][0], 'https://example.com').searchParams.has('q'), false)
    h.cleanup()
})
test('server sort and page size reset page; unapplied inputs block page navigation', () => {
    const h = setup({ rfqPage: meta('buyer', 3) })
    h.rfqFilters.q = 'new'; h.changeRfqPage(4); assert.equal(h.gets.length, 0)
    h.rfqFilters.q = 'buyer'; h.sortRfqColumn({ key: 'buyer', direction: 'asc' }); assert.equal(h.gets.length, 0)
    h.sortRfqColumn({ key: 'created_at', direction: 'asc' }); assert.match(h.gets[0][0], /sort=oldest/); assert.match(h.gets[0][0], /page=1/)
    h.gets[0][2].onFinish(); h.changeRfqPageSize(50); assert.match(h.gets[1][0], /per_page=50/)
    h.cleanup()
})
test('selected rfq refresh uses off-page lookup and preserves review until success, then requires confirmation again', async () => {
    const h = setup({ rfqVisibility: 'full', recordCapabilities: { rfqs: { o: { close: true } } } })
    h.openRecord('rfqs', h.props.rfqs[0], 'close'); assert.equal(h.posts.length, 0)
    const loading = h.refreshSelectedRfq(); assert.match(h.fetches[0].url, /selected_rfq=o/)
    assert.equal(h.pendingOperation.value, 'close')
    respond(h.fetches[0], { id: 'o', status: 'open', target_price: 20, items: [] }, true, 'full'); await loading
    assert.equal(h.pendingOperation.value, ''); assert.equal(h.selectedRfqVisibility.value, 'full')
    h.confirmOperation(); assert.equal(h.posts.length, 0)
    h.cleanup()
})
test('403 and unavailable selected-rfq refresh revoke cached request closure capability', async () => {
    for (const status of ['403', 'unavailable']) {
        const h = setup({ rfqVisibility: 'full', recordCapabilities: { rfqs: { o: { close: true } } } })
        h.openRecord('rfqs', h.props.rfqs[0], 'close')
        const loading = h.refreshSelectedRfq()
        if (status === '403') h.fetches[0].resolve({ ok: false, status: 403 })
        else respond(h.fetches[0], null, true)
        await loading
        assert.equal(h.detailRecord.value.id, 'o', 'keep the reviewed snapshot rather than silently switching records')
        assert.equal(h.operationActions('rfqs', h.detailRecord.value).length, 1)
        h.confirmOperation(); assert.equal(h.posts.length, 0)
        assert.ok(h.formError.value)
        h.cleanup()
    }
})
test('leaving RFQs during selected refresh aborts and removes stale drawer without reopening it', async () => {
    const h = setup()
    h.openRecord('rfqs', h.props.rfqs[0]); const loading = h.refreshSelectedRfq()
    h.page.url = '/dashboard?section=products'; await flush()
    assert.equal(h.fetches[0].options.signal.aborted, true)
    assert.equal(h.detailRecord.value, null)
    respond(h.fetches[0], { id: 'o', status: 'open', items: [] })
    await loading
    assert.equal(h.detailRecord.value, null); assert.equal(h.pendingOperation.value, ''); assert.equal(h.busy.value, false)
    h.cleanup()
})
test('confirmed request closure returns to captured rfq query/page after static redirect', async () => {
    const h = setup({ rfqPage: meta('buyer', 3, 'open'), rfqVisibility: 'full', recordCapabilities: { rfqs: { o: { close: true } } } })
    h.openRecord('rfqs', h.props.rfqs[0], 'close'); h.confirmOperation()
    assert.equal(h.posts[0][0], '/dashboard/rfqs/o/close')
    h.props.rfqPage = meta()
    h.posts[0][2].onSuccess({ props: { flash: {} } }); h.posts[0][2].onFinish(); await flush()
    const query = new URL(h.gets[0][0], 'https://example.com').searchParams
    assert.equal(query.get('q'), 'buyer'); assert.equal(query.get('page'), '3'); assert.equal(query.get('status'), 'open'); assert.equal(h.detailRecord.value, null)
    h.cleanup()
})

test('reviewing request closure after off-page refresh retains fresh grants and the original return filters', async () => {
    const h = setup({ rfqPage: meta('buyer', 3, 'open'), rfqVisibility: 'full', recordCapabilities: { rfqs: { o: { close: true } } } })
    h.openRecord('rfqs', h.props.rfqs[0], 'close')
    // Static mutation error redirect has changed page props; selected record remains reviewed.
    h.props.rfqPage = meta(); h.props.recordCapabilities.rfqs = {}
    const check = h.refreshSelectedRfq()
    respond(h.fetches[0], { id: 'o', status: 'open', items: [] }, true, 'full')
    await check
    h.openRecord('rfqs', h.detailRecord.value, 'close')
    assert.equal(h.pendingOperation.value, 'close')
    assert.equal(h.selectedRfqVisibility.value, 'full')
    h.confirmOperation()
    assert.equal(h.posts.length, 1)
    h.posts[0][2].onSuccess({ props: { flash: {} } }); h.posts[0][2].onFinish(); await flush()
    assert.equal(new URL(h.gets[0][0], 'https://example.com').searchParams.get('q'), 'buyer')
    assert.equal(new URL(h.gets[0][0], 'https://example.com').searchParams.get('page'), '3')
    h.cleanup()
})

test('RFQ award retains reviewed version and loses all decision grants after denied refresh', async () => {
    const quote = {id:'q1',updated_at:'2026-11-04T12:00:00Z',unit_price:12,status:'pending'}
    const h = setup({recordCapabilities:{rfqs:{o:{close:true}},quotes:{q1:{accept:true}}}})
    h.props.rfqs[0].quotes = [quote]
    h.openRecord('rfqs',h.props.rfqs[0]); h.prepareOperation('award',quote)
    quote.updated_at = '2026-11-04T13:00:00Z'
    h.confirmOperation()
    assert.equal(h.posts[0][1].quote_version,'2026-11-04T12:00:00Z')
    h.posts[0][2].onError({error:'Quotation changed'}); h.posts[0][2].onFinish()
    assert.equal(h.pendingQuote.value.updated_at,'2026-11-04T12:00:00Z')
    const loading=h.refreshSelectedRfq(); h.fetches[0].resolve({ok:false,status:403}); await loading
    h.confirmOperation(); assert.equal(h.posts.length,1)
    assert.equal(h.recordCan('quotes','q1','accept'),false)
    h.cleanup()
})
test('RFQ list filters cannot replace an open quote draft', () => {
    const h=setup({recordCapabilities:{rfqs:{o:{quote:true}}}})
    h.openRecord('rfqs',h.props.rfqs[0],'quote')
    h.quoteForm.notes='Keep my terms'; h.rfqFilters.q='new buyer'; h.applyRfqFilters()
    assert.equal(h.gets.length,0); assert.equal(h.quoteForm.notes,'Keep my terms')
    h.openRecord('rfqs',{id:'different'},'view')
    assert.equal(h.quoteTarget.value.id,'o')
    h.cleanup()
})


test('leaving RFQs saves a dirty quote without leaking the drawer into another section; recovery checks current scope', async()=>{
 const h=setup({recordCapabilities:{rfqs:{o:{quote:true}}}})
 h.openRecord('rfqs',h.props.rfqs[0],'quote');h.quoteForm.notes='My shipping terms';h.quoteForm.lead_time_days=0
 h.page.url='/dashboard?section=orders';await flush()
 assert.equal(h.quoteTarget.value,null);assert.equal(h.pendingQuoteDraft.value.form.notes,'My shipping terms')
 h.page.url='/dashboard?section=rfqs';await flush()
 let recovery=h.recoverQuoteDraft();assert.match(h.fetches[0].url,/selected_rfq=o/)
 h.fetches[0].resolve({ok:true,json:async()=>({props:{selectedRfqStatus:'unavailable'}})});await recovery
 assert.equal(h.quoteTarget.value,null);assert.equal(h.pendingQuoteDraft.value.form.notes,'My shipping terms');assert.ok(h.quoteRecoveryError.value)
 recovery=h.recoverQuoteDraft()
 h.fetches[1].resolve({ok:true,json:async()=>({props:{selectedRfqStatus:'available',selectedRfq:{id:'o',title:'Current title'},recordCapabilities:{rfqs:{o:{quote:true}}}}})});await recovery
 assert.equal(h.quoteForm.lead_time_days,0);assert.equal(h.quoteTarget.value.title,'Current title');assert.equal(h.quoteForm.notes,'My shipping terms');assert.equal(h.pendingQuoteDraft.value,null)
 h.cleanup()
})
test('leaving RFQs during draft recovery aborts lookup and retains saved terms', async()=>{
 const h=setup({recordCapabilities:{rfqs:{o:{quote:true}}}})
 h.openRecord('rfqs',h.props.rfqs[0],'quote');h.quoteForm.notes='Keep'
 h.page.url='/dashboard?section=orders';await flush();h.page.url='/dashboard?section=rfqs';await flush()
 const recovery=h.recoverQuoteDraft();h.page.url='/dashboard?section=orders';await flush()
 assert.equal(h.fetches[0].options.signal.aborted,true)
 h.fetches[0].resolve({ok:true,json:async()=>({props:{selectedRfqStatus:'available',selectedRfq:{id:'o'},recordCapabilities:{rfqs:{o:{quote:true}}}}})});await recovery
 assert.equal(h.quoteTarget.value,null);assert.equal(h.pendingQuoteDraft.value.form.notes,'Keep');h.cleanup()
})
