import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import * as Vue from 'vue'
const source = readFileSync(new URL('../../resources/js/Pages/Dashboard/Index.vue', import.meta.url), 'utf8')
function harness(grants) {
    const requests = []
    const props = { currentUser: { id: 'viewer', role: 'buyer' }, stats: {}, products: [], orders: [], currencies: [], regions: [], recordCapabilities: { products: grants }, capabilities: { readProducts: true, createProducts: false } }
    let script = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
    script += '\n globalThis.harness = { productActions, productAction, productOpen, deleteTarget, busy };'
    const context = { exports: {}, defineProps: () => props, URL,
        require: name => name === 'vue' ? Vue : name === '@inertiajs/vue3' ? { usePage: () => ({ url: '/dashboard?section=products', props: {} }), router: { post: (...args) => requests.push(args) } } : name.includes('useCurrency') ? { useCurrency: () => ({ formatPrice: String }) } : {},
    }
    vm.runInNewContext(ts.transpileModule(script, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.CommonJS } }).outputText, context)
    return { ...context.harness, requests, props }
}
test('product menus only expose per-record permissions, never inferred role permissions', () => {
    const h = harness({ allowed: { update: true, delete: true }, readonly: {} })
    assert.deepEqual(Array.from(h.productActions({ id: 'allowed' }), x => x.id), ['edit', 'delete'])
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
    assert.equal((source.match(/<AdminTable/g) || []).length, 3)
    assert.doesNotMatch(source, /<table class="workspace-table"/)
    assert.match(source, /section === 'products' &&\s+capabilities\?\.createProducts/)
    assert.match(source, /:label="`Actions for \$\{product.title\}`"/)
    assert.match(source, /@submit.prevent="removeProduct"/)
    assert.match(source, /View items and payment for \$\{order.order_number\}/)
})

test('product and user identities are semantic row headers', () => {
    for (const record of ['product', 'user']) {
        const row = source.indexOf(`<tr v-for="${record} in visibleRecords" :key="${record}.id">`)
        assert.ok(row >= 0)
        assert.match(source.slice(row, row + 160), /<th scope="row">/)
    }
})
