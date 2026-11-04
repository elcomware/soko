import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import * as Vue from 'vue'
import { safeProductDraft } from '../../resources/js/design/productDraft.ts'
const source = readFileSync('resources/js/Pages/Dashboard/Index.vue', 'utf8')
const meta = (q = '', page = 1) => ({ page, per_page: 25, total: 101, filtered_total: 76, last_page: 4, from: (page-1)*25+1, to: Math.min(page*25,76), filters: { q, category: '', sort: 'newest' } })
function setup({ draft, props: extra = {} } = {}) {
    const gets = [], posts = [], timers = [], fetches = []
    const props = Vue.reactive({ currentUser: { id: 'm', role: 'merchant' }, stats: { total_products: 101 }, products: [{ id: 'p', title: 'Original', updated_at: 'v1', price_tiers: [], gallery_urls: [] }], productPage: meta(), productCategories: ['Clothing', 'Hardware'], recordCapabilities: { products: { p: { update: true, delete: true } } }, capabilities: { readProducts: true, createProducts: true }, ...extra })
    const page = Vue.reactive({ url: '/dashboard?section=products', props: {} })
    let cleanup
    let script = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
    script += '\n globalThis.h = { productFilters, productFiltersPending, productListUrl, applyProductFilters, resetProductFilters, changeProductPage, productListLoading, productListError, visibleRecords, openProduct, productForm, productOpen, productDirty, pendingProductDraft, recoverOffpageProduct, checkSavedProduct, freshProduct, replaceWithSavedProduct, productRefreshMessage, saveProduct, openDelete, removeProduct, productReturnUrl, editingProduct, pendingProductDiscard, discardPendingProductDraft, sortProductColumn, changeProductPageSize, productColumns, productTableSort };'
    const context = { exports: {}, defineProps: () => props, URL, URLSearchParams,
        setTimeout: (callback, delay) => { timers.push({ callback, delay, cleared: false }); return timers.length }, clearTimeout: id => { if(timers[id-1]) timers[id-1].cleared = true },
        fetch: (url, options) => new Promise(resolve => fetches.push({ url, options, resolve })),
        require: name => name === 'vue' ? { ...Vue, onBeforeUnmount: fn => { cleanup = fn } } : name === '@inertiajs/vue3' ? { usePage: () => page, router: { get: (...args) => gets.push(args), post: (...args) => posts.push(args), restore: () => draft, remember: () => {} } } : name.includes('productDraft') ? { safeProductDraft } : name.includes('useCurrency') ? { useCurrency: () => ({ formatPrice: String }) } : {},
    }
    vm.runInNewContext(ts.transpileModule(script, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.CommonJS } }).outputText, context)
    return { ...context.h, props, page, gets, posts, timers, fetches, cleanup: () => cleanup() }
}
const flush = async () => { await Vue.nextTick(); await new Promise(resolve => setImmediate(resolve)) }
const respond = (request, id, record, allowed = true) => request.resolve({ ok: true, json: async () => ({ props: { selectedProduct: record, selectedProductStatus: record ? 'available' : 'unavailable', recordCapabilities: { products: { [id]: { update: allowed } } } } }) })
test('server product results are not re-filtered locally and filters reset page in an encoded URL', async () => {
    const h = setup({ props: { productPage: meta('server search', 3) } })
    assert.equal(h.visibleRecords.value.length, 1)
    h.productFilters.q = 'Équipement & tools'
    h.productFilters.category = 'Hardware'
    h.productFilters.sort = 'title_asc'
    h.applyProductFilters()
    const query = new URL(h.gets[0][0], 'https://example.com').searchParams
    assert.equal(query.get('q'), 'Équipement & tools')
    assert.equal(query.get('category'), 'Hardware')
    assert.equal(query.get('sort'), 'title_asc')
    assert.equal(query.get('page'), '1')
    h.cleanup()
})
test('debounced searches serialize requests and never replace a newer input with older response filters', async () => {
    const h = setup()
    h.productFilters.q = 'first'
    await flush()
    assert.equal(h.timers.at(-1).delay, 350)
    h.timers.at(-1).callback()
    assert.equal(h.gets.length, 1)
    h.productFilters.q = 'latest'
    await flush()
    h.props.productPage = meta('first')
    assert.equal(h.productFilters.q, 'latest')
    h.gets[0][2].onFinish()
    await flush()
    assert.equal(h.gets.length, 2)
    assert.equal(new URL(h.gets[1][0], 'https://example.com').searchParams.get('q'), 'latest')
    h.productFilters.q = ''
    await flush()
    assert.equal(h.timers.at(-1).delay, 0)
    h.props.productPage = meta('latest')
    h.gets[1][2].onFinish()
    await flush()
    assert.equal(new URL(h.gets[2][0], 'https://example.com').searchParams.has('q'), false)
    h.cleanup()
})
test('page changes preserve committed query and are blocked while unapplied filters or editor are active', () => {
    const h = setup({ props: { productPage: meta('needle', 2) } })
    h.changeProductPage(3)
    assert.equal(new URL(h.gets[0][0], 'https://example.com').searchParams.get('q'), 'needle')
    h.gets[0][2].onFinish()
    h.productFilters.q = 'new'
    h.changeProductPage(4)
    assert.equal(h.gets.length, 1)
    h.productFilters.q = 'needle'
    h.openProduct(h.props.products[0])
    h.applyProductFilters()
    assert.equal(h.gets.length, 1)
    h.cleanup()
})
test('off-page draft is retained and only reopened after scoped selected-product authorization', async () => {
    const h = setup({ draft: { open: true, productId: 'off', form: { title: 'Unsaved work', version: 'old' } } })
    assert.ok(h.pendingProductDraft.value)
    assert.equal(h.productOpen.value, false)
    const pending = h.recoverOffpageProduct()
    assert.match(h.fetches[0].url, /selected_product=off/)
    respond(h.fetches[0], 'off', { id: 'off', title: 'Saved', updated_at: 'new' })
    await pending
    assert.equal(h.productOpen.value, true)
    assert.equal(h.productForm.title, 'Unsaved work')
    assert.equal(h.productForm.version, 'old')
    assert.equal(h.productDirty.value, true)
    h.cleanup()
    const denied = setup({ draft: { open: true, productId: 'off', form: { title: 'Unsaved work' } } })
    const check = denied.recoverOffpageProduct()
    respond(denied.fetches[0], 'off', null, false)
    await check
    assert.equal(denied.productOpen.value, false)
    assert.ok(denied.pendingProductDraft.value)
    assert.match(denied.productRefreshMessage.value, /not been discarded/)
    denied.cleanup()
})
test('selected refresh preserves dirty form until explicit replace', async () => {
    const h = setup()
    h.openProduct(h.props.products[0])
    h.productForm.title = 'Draft edit'
    const check = h.checkSavedProduct()
    respond(h.fetches[0], 'p', { id: 'p', title: 'New server value', updated_at: 'v2', price_tiers: [] })
    await check
    assert.equal(h.productForm.title, 'Draft edit')
    assert.equal(h.productForm.version, 'v1')
    h.replaceWithSavedProduct()
    assert.equal(h.productForm.title, 'New server value')
    assert.equal(h.productForm.version, 'v2')
    assert.equal(h.productDirty.value, false)
    h.cleanup()
})
test('successful product save returns to captured filters/page despite static redirect props', async () => {
    const h = setup({ props: { productPage: meta('needle', 3) } })
    h.openProduct(h.props.products[0])
    h.saveProduct({ target: { elements: [], querySelectorAll: () => [] } })
    assert.equal(h.posts.length, 1)
    h.props.productPage = meta()
    h.posts[0][2].onSuccess({ props: { flash: {} } })
    h.posts[0][2].onFinish()
    await flush()
    assert.equal(h.productOpen.value, false)
    const query = new URL(h.gets[0][0], 'https://example.com').searchParams
    assert.equal(query.get('q'), 'needle')
    assert.equal(query.get('page'), '3')
    h.cleanup()
})

test('clearing while first search is in flight queues the empty query, even though old committed query was empty', async () => {
    const h = setup()
    h.productFilters.q = 'first'
    h.applyProductFilters()
    await flush()
    h.productFilters.q = ''
    await flush()
    assert.equal(h.timers.at(-1).delay, 0)
    h.props.productPage = meta('first')
    h.gets[0][2].onFinish()
    await flush()
    assert.equal(h.gets.length, 2)
    assert.equal(new URL(h.gets[1][0], 'https://example.com').searchParams.has('q'), false)
    h.cleanup()
})
test('immediate reset does not enqueue a redundant second request from its input watcher', async () => {
    const h = setup({ props: { productPage: meta('old') } })
    h.resetProductFilters()
    await flush()
    h.props.productPage = meta()
    h.gets[0][2].onFinish()
    await flush()
    assert.equal(h.gets.length, 1)
    h.cleanup()
})

test('off-page recovery blocks competing edits/deletes/navigation, and unresolved draft requires explicit discard', async () => {
    const h = setup({ draft: { open: true, productId: 'off', form: { title: 'Keep me' } } })
    h.openProduct()
    assert.equal(h.productOpen.value, false)
    assert.ok(h.pendingProductDraft.value)
    const recovery = h.recoverOffpageProduct()
    h.openProduct(h.props.products[0]); h.openDelete(h.props.products[0]); h.applyProductFilters()
    assert.equal(h.productOpen.value, false)
    assert.equal(h.gets.length, 0)
    respond(h.fetches[0], 'off', null, false)
    await recovery
    h.discardPendingProductDraft()
    assert.ok(h.pendingProductDraft.value, 'discard requires explicit confirmation state')
    h.pendingProductDiscard.value = true
    h.discardPendingProductDraft()
    assert.equal(h.pendingProductDraft.value, null)
    h.openProduct()
    assert.equal(h.productOpen.value, true)
    h.cleanup()
})

test('shared product header sort and pagination map only supported server operations',()=>{
 const h=setup({props:{productPage:meta('',3)}})
 h.sortProductColumn({key:'price',direction:'asc'});assert.equal(h.gets.length,0)
 h.sortProductColumn({key:'title',direction:'asc'});assert.match(h.gets[0][0],/sort=title_asc/);assert.match(h.gets[0][0],/page=1/)
 assert.equal(h.productTableSort.value.key,'title');assert.equal(h.productTableSort.value.direction,'asc')
 h.gets[0][2].onFinish();h.changeProductPageSize(50);assert.match(h.gets[1][0],/per_page=50/);assert.match(h.gets[1][0],/page=1/)
 h.gets[1][2].onFinish();h.sortProductColumn({key:'created_at',direction:'desc'});assert.match(h.gets[2][0],/sort=newest/)
 assert.equal(h.productColumns.value.filter(c=>c.sortable).length,2);h.cleanup()
})
test('shared product sorting/page size cannot mutate or replace an open draft',()=>{
 const h=setup();h.openProduct(h.props.products[0]);h.productForm.title='Unsaved changes'
 h.sortProductColumn({key:'title',direction:'desc'});h.changeProductPageSize(100)
 assert.equal(h.gets.length,0);assert.equal(h.productFilters.sort,'newest');assert.equal(h.productFilters.per_page,25);assert.equal(h.productForm.title,'Unsaved changes');h.cleanup()
})
test('Products uses shared table slots and paging even when filtered results are empty',()=>{
 const table=source.slice(source.indexOf('v-else-if="section === \'products\'"\n                    caption="Product catalog"'),source.indexOf('v-else-if="section === \'orders\'"'))
 assert.match(table,/:columns="productColumns"/);assert.match(table,/resizable-columns/);assert.match(table,/resizable-height/);assert.match(table,/#filters/);assert.match(table,/<AdminPagination/);assert.doesNotMatch(table,/#header|product-pagination/)
 assert.match(table,/productPage.total === 0 &&[\s\S]*?capabilities\?\.createProducts/)
 const h=setup({props:{products:[],productPage:{...meta('absent'),filtered_total:0,from:0,to:0,last_page:1}}});assert.equal(h.visibleRecords.value.length,0);h.resetProductFilters();assert.equal(h.gets.length,1);h.cleanup()
})
