import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import { reactive, ref, computed, watch } from 'vue'

function cartHarness() {
    const source = readFileSync(new URL('../../resources/js/Pages/Store/Cart.vue', import.meta.url), 'utf8')
    let script = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
    script = script.replace(/^import .*$/gm, '')
    script += '\n globalThis.harness = { checkout, schedule, flushItem, removeItem, busy, quantities, failedRemovals, error };'
    const requests = []
    const item = { product_id: 'product', quantity: 1, product: { moq: 1, stock: 20 } }
    const props = { cart: { items: [item], checkout_token: 'stable-token' } }
    const context = { reactive, ref, computed, watch, onBeforeUnmount: () => {}, defineProps: () => props,
        usePage: () => ({ props: { auth: { user: { id: 'buyer' } } } }),
        useCurrency: () => ({ formatPrice: String, activeCurrency: ref({ code: 'USD' }) }),
        router: { post: (url, data, options) => requests.push({ url, data, options }) },
        setTimeout, clearTimeout,
    }
    vm.runInNewContext(ts.transpileModule(script, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.None } }).outputText, context)
    return { ...context.harness, item, requests, source }
}
const tick = async () => { for (let i = 0; i < 8; i++) await Promise.resolve() }
test('checkout freezes mutation handlers and sends stable idempotency token', async () => {
    const h = cartHarness()
    const checkout = h.checkout()
    await tick()
    assert.equal(h.busy.value, true)
    assert.equal(h.requests.length, 1)
    assert.equal(h.requests[0].url, '/checkout')
    assert.equal(h.requests[0].data.checkout_token, 'stable-token')
    h.schedule(h.item); h.flushItem(h.item); h.removeItem(h.item)
    await tick()
    assert.equal(h.requests.length, 1)
    h.requests[0].options.onNetworkError()
    h.requests[0].options.onFinish()
    await checkout
    assert.match(h.error.value, /Check your orders first/)
    assert.match(h.source, /v-model="quantities\[item.product_id\]"\s+:disabled="busy"/)
})
test('failed removal retains intent and retry sends deletion, not quantity update', async () => {
    const h = cartHarness()
    h.removeItem(h.item)
    assert.equal(h.requests[0].url, '/cart/items/product/delete')
    h.requests[0].options.onNetworkError()
    h.requests[0].options.onFinish()
    await tick()
    assert.equal(h.failedRemovals.value.has('product'), true)
    h.flushItem(h.item)
    assert.equal(h.requests.length, 2)
    assert.equal(h.requests[1].url, '/cart/items/product/delete')
    h.requests[1].options.onSuccess({ props: { flash: {} } })
    h.requests[1].options.onFinish()
    await tick()
    assert.equal(h.failedRemovals.value.size, 0)
})

test('regular Remove after failed deletion clears residual failure state', async () => {
    const h = cartHarness()
    h.removeItem(h.item)
    h.requests[0].options.onNetworkError()
    h.requests[0].options.onFinish()
    await tick()
    assert.equal(h.failedRemovals.value.has('product'), true)
    h.removeItem(h.item)
    assert.equal(h.requests[1].url, '/cart/items/product/delete')
    h.requests[1].options.onSuccess({ props: { flash: {} } })
    h.requests[1].options.onFinish()
    await tick()
    assert.equal(h.failedRemovals.value.size, 0)
})
