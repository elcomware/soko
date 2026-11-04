import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import * as Vue from 'vue'
import { safeProductDraft } from '../../resources/js/design/productDraft.ts'
const source = readFileSync(
    new URL('../../resources/js/Pages/Dashboard/Index.vue', import.meta.url),
    'utf8'
)
function harness(grants, navigation = [], initial = {}) {
    const requests = []
    const reloads = []
    const props = {
        currentUser: { id: 'viewer', role: 'buyer' },
        stats: {},
        products: [],
        orders: [],
        currencies: [],
        regions: [],
        recordCapabilities: { products: grants },
        capabilities: { readProducts: true, createProducts: false },
        ...initial,
    }
    let script = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
    script +=
        '\n globalThis.harness = { productActions, productAction, productOpen, deleteTarget, busy, productForm, productDirty, openProduct, closeProduct, operationActions, openRecord, detailRecord, detailKind, pendingOperation, prepareOperation, confirmOperation, accessSettingsHref, quoteForm, quoteDirty, openQuote, pendingQuote, formError, refreshRecord, visibleRecords, accountRole, accountVerification, accountRoles, search, hasActiveFilters, clearRecordFilters, activeNavigation, recordStatus, statusOptions };'
    const context = {
        exports: {},
        defineProps: () => props,
        URL,
        URLSearchParams,
        setTimeout, clearTimeout,
        require: (name) =>
            name === 'vue'
                ? { ...Vue, onBeforeUnmount: () => {} }
                : name === '@inertiajs/vue3'
                  ? {
                        useRemember: data => Vue.reactive(data),
                        usePage: () => ({
                            url: '/dashboard?section=' + (initial.section || 'products'),
                            props: { workspaceNavigation: navigation },
                        }),
                        router: { post: (...args) => requests.push(args), reload: options => reloads.push(options), restore: () => undefined, remember: () => {} },
                    }
                  : name.includes('productDraft') ? { safeProductDraft } : name.includes('useCurrency')
                    ? { useCurrency: () => ({ formatPrice: String }) }
                    : {},
    }
    vm.runInNewContext(
        ts.transpileModule(script, {
            compilerOptions: {
                target: ts.ScriptTarget.ES2022,
                module: ts.ModuleKind.CommonJS,
            },
        }).outputText,
        context
    )
    return { ...context.harness, requests, reloads, props }
}
test('product menus only expose per-record permissions, never inferred role permissions', () => {
    const h = harness({ allowed: { update: true, delete: true }, readonly: {} })
    assert.deepEqual(
        Array.from(h.productActions({ id: 'allowed' }), (x) => x.id),
        ['edit', 'delete']
    )
    assert.equal(h.productActions({ id: 'readonly' }).length, 0)
    assert.equal(h.productActions({ id: 'unknown' }).length, 0)
    h.productAction({ id: 'readonly' }, 'delete')
    assert.equal(h.deleteTarget.value, null)
    h.productAction({ id: 'readonly' }, 'edit')
    assert.equal(h.productOpen.value, false)
})
test('delete menu opens existing confirmation without sending a mutation, and busy blocks dispatch', () => {
    const h = harness({ p: { update: true, delete: true } })
    const product = { id: 'p', title: 'Machine' }
    h.productAction(product, 'delete')
    assert.equal(h.deleteTarget.value.id, 'p')
    assert.equal(h.requests.length, 0)
    h.deleteTarget.value = null
    h.busy.value = true
    h.productAction(product, 'delete')
    h.productAction(product, 'edit')
    assert.equal(h.deleteTarget.value, null)
    assert.equal(h.productOpen.value, false)
})
test('lists use shared table and create CTA requires creation grant', () => {
    assert.equal((source.match(/<AdminTable/g) || []).length, 5)
    assert.doesNotMatch(source, /<table class="workspace-table"/)
    assert.match(
        source,
        /section === 'products' &&\s+capabilities\?\.createProducts/
    )
    assert.match(source, /:label="`Actions for \$\{product.title\}`"/)
    assert.match(source, /@submit.prevent="removeProduct"/)
    assert.match(source, /Orders · authorized lines only/)
})

test('product and user identities are semantic row headers', () => {
    for (const record of ['product', 'user']) {
        const row = source.indexOf(
            `<tr v-for="${record} in visibleRecords" :key="${record}.id">`
        )
        assert.ok(row >= 0)
        assert.match(source.slice(row, row + 160), /<th scope="row">/)
    }
})

test('product drawer snapshots initial data and tracks nested edits without switching drafts', () => {
    const h = harness({ p: { update: true }, other: { update: true } })
    h.openProduct({
        id: 'p',
        title: 'Original',
        price_tiers: [{ min_qty: 1, max_qty: 5, price: 4 }],
    })
    assert.equal(h.productDirty.value, false)
    h.productForm.price_tiers[0].price = 3
    assert.equal(h.productDirty.value, true)
    h.openProduct({ id: 'other', title: 'Other' })
    assert.equal(h.productForm.title, 'Original')
    h.busy.value = true
    h.closeProduct()
    assert.equal(h.productOpen.value, true)
    assert.match(source, /<AdminDrawer[\s\S]*?:dirty="productDirty"/)
    assert.match(source, /form="product-editor-form"/)
    assert.match(source, /@click="requestClose"/)
})

test('operations only expose granted actions and require explicit confirmation before posting', () => {
    const h = harness({})
    h.props.recordCapabilities.orders = { o: { cancel: true } }
    const order = { id: 'o', order_number: 'ORDER-1', status: 'pending', items: [] }
    assert.deepEqual(Array.from(h.operationActions('orders', order), x => x.id), ['view', 'cancel'])
    assert.equal(h.operationActions('orders', { id: 'other' }).length, 1)
    h.openRecord('orders', order, 'cancel')
    assert.equal(h.pendingOperation.value, 'cancel')
    assert.equal(h.requests.length, 0)
    h.confirmOperation()
    assert.equal(h.requests[0][0], '/dashboard/orders/o/cancel')
    h.confirmOperation()
    assert.equal(h.requests.length, 1, 'busy guard prevents repeat submission')
})
test('quotation award requires matching visible quote and per-quote grant', () => {
    const h = harness({})
    const rfq = { id: 'r', title: 'Request', quotes: [{ id: 'q', updated_at: '2030-01-01T00:00:00Z', unit_price: 10 }] }
    h.props.recordCapabilities.quotes = { q: { accept: true }, foreign: { accept: true } }
    h.openRecord('rfqs', rfq)
    h.prepareOperation('award', { id: 'foreign' })
    assert.equal(h.pendingOperation.value, '')
    h.prepareOperation('award', rfq.quotes[0])
    assert.equal(h.requests.length, 0)
    h.confirmOperation()
    assert.equal(h.requests[0][0], '/dashboard/rfqs/r/quotes/q/accept')
})
test('offer decisions require confirmation and account settings link requires exact navigation grant', () => {
    const h = harness({})
    h.props.recordCapabilities.offers = { offer: { decide: true } }
    h.openRecord('offers', { id: 'offer', product_title: 'Product' }, 'accepted')
    assert.equal(h.requests.length, 0)
    h.confirmOperation()
    assert.equal(h.requests[0][0], '/dashboard/offers/offer/decision?section=offers')
    assert.equal(h.requests[0][1].status, 'accepted')
    assert.equal(h.accessSettingsHref.value, '')
    const staff = harness({}, [{ href: '/dashboard/settings?panel=access' }])
    assert.equal(staff.accessSettingsHref.value, '/dashboard/settings?panel=access')
    assert.deepEqual(Array.from(staff.operationActions('users', { id: 'account' }), x => x.id), ['view'])
})

test('RFQ close requires scoped grant and quotation drawer tracks dirty values', () => {
    const h = harness({})
    h.props.recordCapabilities.rfqs = { r: { close: true, quote: true } }
    const rfq = { id: 'r', title: 'Request', target_price: 10, quotes: [] }
    h.openRecord('rfqs', rfq, 'close')
    assert.equal(h.requests.length, 0)
    h.confirmOperation()
    assert.equal(h.requests[0][0], '/dashboard/rfqs/r/close')
    const quote = harness({})
    quote.props.recordCapabilities.rfqs = { r: { quote: true } }
    quote.openQuote(rfq)
    assert.equal(quote.quoteDirty.value, false)
    quote.quoteForm.unit_price = 11
    assert.equal(quote.quoteDirty.value, true)
})

test('quote award posts reviewed version, keeps snapshot on conflict, and requires fresh review after explicit reload', () => {
    const h = harness({})
    h.props.recordCapabilities.quotes = { q: { accept: true } }
    const quote = { id: 'q', updated_at: '2030-01-01T00:00:00Z', unit_price: 10 }
    h.openRecord('rfqs', { id: 'r', quotes: [quote] })
    h.prepareOperation('award', quote)
    quote.unit_price = 99
    quote.updated_at = '2030-01-02T00:00:00Z'
    h.confirmOperation()
    assert.equal(h.requests[0][1].quote_version, '2030-01-01T00:00:00Z')
    assert.equal(h.pendingQuote.value.unit_price, 10)
    h.requests[0][2].onSuccess({ props: { flash: { error: 'Quotation changed. Review current terms.' } } })
    h.requests[0][2].onFinish()
    assert.equal(h.pendingOperation.value, 'award')
    assert.equal(h.pendingQuote.value.unit_price, 10)
    h.refreshRecord()
    assert.equal(h.pendingQuote.value.unit_price, 10, 'do not replace terms before refresh succeeds')
    h.reloads[0].onSuccess({ props: { rfqs: [{ id: 'r', quotes: [quote] }] } })
    h.reloads[0].onFinish()
    assert.equal(h.pendingQuote.value, null)
    assert.equal(h.pendingOperation.value, '')
    h.confirmOperation()
    assert.equal(h.requests.length, 1, 'refreshed terms require a new explicit review selection')
})

test('whole-order cancellation disclosure includes unseen supplier lines', () => {
    assert.match(source, /entire order, including any supplier lines you cannot see/)
})

test('Accounts filters intersect role, verification and country search without adding mutation actions', () => {
    const h = harness({}, [], { section: 'users', capabilities: { manageUsers: true }, users: [
        { id: 'one', name: 'Alice', email: 'a@example.com', role: 'merchant', verified: true, country: 'Kenya' },
        { id: 'two', name: 'Bob', email: 'b@example.com', role: 'merchant', verified: false, country: 'Kenya' },
        { id: 'three', name: 'Carol', email: 'c@example.com', role: 'buyer', verified: true, country: 'Ghana' },
    ] })
    assert.equal(h.activeNavigation.value.label, 'Accounts')
    assert.equal(h.visibleRecords.value.length, 3)
    h.accountRole.value = 'merchant'
    h.accountVerification.value = 'verified'
    assert.equal(h.visibleRecords.value.length, 1)
    assert.equal(h.visibleRecords.value[0].id, 'one')
    h.search.value = 'ghana'
    assert.equal(h.visibleRecords.value.length, 0)
    assert.equal(h.hasActiveFilters.value, true)
    h.clearRecordFilters()
    assert.equal(h.visibleRecords.value.length, 3)
    assert.equal(h.hasActiveFilters.value, false)
    h.search.value = 'kenya'
    assert.equal(h.visibleRecords.value.length, 2)
    h.search.value = 'verified'
    assert.equal(h.visibleRecords.value.length, 2, 'verified search must not also match unverified')
    h.search.value = 'unverified'
    assert.equal(h.visibleRecords.value.length, 1)
    assert.equal(h.visibleRecords.value[0].id, 'two')
    assert.equal(h.operationActions('users', h.visibleRecords.value[0]).length, 1)
})
test('status-only filter is an active filter, and clear resets it', () => {
    const h = harness({}, [], { section: 'offers', capabilities: { readOffers: true }, offers: [{ id: 'p', status: 'pending' }] })
    h.recordStatus.value = 'declined'
    assert.equal(h.visibleRecords.value.length, 0)
    assert.equal(h.hasActiveFilters.value, true)
    h.clearRecordFilters()
    assert.equal(h.visibleRecords.value.length, 1)
    assert.equal(h.hasActiveFilters.value, false)
    assert.match(source, /No records match your current filters/)
})
test('first-column identities open read-only detail entry points with contextual names and touch-sized controls', () => {
    for (const [record, kind] of [['rfq', 'rfqs'], ['offer', 'offers'], ['user', 'users']]) {
        const start = source.indexOf(`<tr v-for="${record} in visibleRecords"`)
        const header = source.slice(start, source.indexOf('</th>', start))
        assert.match(header, /class="operation-identity"/)
        assert.ok(header.includes(`openRecord('${kind}', ${record})`))
        assert.match(header, /:aria-label=/)
        assert.match(header, /type="button"/)
    }
    assert.match(source, /\.workspace-search \.icon-button \{\s*min-height: 44px;\s*min-width: 44px;/)
})

test('selected status and role remain visible after last matching record changes until filters reset', () => {
    const offers = Vue.reactive([{ id: 'one', status: 'pending' }])
    const h = harness({}, [], { section: 'offers', capabilities: { readOffers: true }, offers })
    h.recordStatus.value = 'pending'
    assert.ok(h.statusOptions.value.includes('pending'))
    offers[0].status = 'accepted'
    assert.equal(h.visibleRecords.value.length, 0)
    assert.ok(h.statusOptions.value.includes('pending'))
    assert.ok(h.statusOptions.value.includes('accepted'))
    h.clearRecordFilters()
    assert.equal(h.statusOptions.value.includes('pending'), false)
    const users = Vue.reactive([{ id: 'a', role: 'merchant', verified: true }])
    const accounts = harness({}, [], { section: 'users', capabilities: { manageUsers: true }, users })
    accounts.accountRole.value = 'merchant'
    assert.ok(accounts.accountRoles.value.includes('merchant'))
    users[0].role = 'buyer'
    assert.equal(accounts.visibleRecords.value.length, 0)
    assert.ok(accounts.accountRoles.value.includes('merchant'))
    accounts.clearRecordFilters()
    assert.equal(accounts.accountRoles.value.includes('merchant'), false)
})
