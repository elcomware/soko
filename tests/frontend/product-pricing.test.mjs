import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import { reactive, ref, computed, watch, nextTick } from 'vue'
import { productImage, isIllustrativeImage } from '../../resources/js/design/merchandising.ts'
function harness() {
    const source = readFileSync(new URL('../../resources/js/Pages/Store/ProductDetail.vue', import.meta.url), 'utf8')
    let script = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1].replace(/^import .*$/gm, '')
    script += '\n globalThis.harness = { refreshQuote, quantity, quote, currentQuote, unitPrice, clock };'
    const pending = []
    const props = reactive({ product: { id: 'p', moq: 1, stock: 100, price: 100, pricing: { quantity: 1, unit_price: 80, original_unit_price: 100, promotion_ends_at: '2030-01-01T00:00:00Z' } } })
    const context = { ref, computed, watch, productImage, isIllustrativeImage, onMounted: () => {}, onBeforeUnmount: () => {},
        defineProps: () => props, useCurrency: () => ({}), usePage: () => ({ props: {} }),
        useMarketplaceForm: () => ({ busy: ref(false), error: ref(''), validate: () => true, post: () => {}, fieldState: () => ({}) }),
        setTimeout: () => 1, clearTimeout: () => {}, setInterval: () => 1, clearInterval: () => {}, AbortController,
        fetch: (url, options) => new Promise(resolve => pending.push({ url, options, resolve })),
    }
    vm.runInNewContext(ts.transpileModule(script, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.None } }).outputText, context)
    return { ...context.harness, pending }
}
test('quantity-specific detail quote ignores stale responses', async () => {
    const h = harness()
    assert.equal(h.unitPrice.value, 80)
    const first = h.refreshQuote()
    h.quantity.value = 10
    await nextTick()
    assert.equal(h.unitPrice.value, null)
    assert.equal(h.pending[0].options.signal.aborted, true)
    const second = h.refreshQuote()
    assert.match(h.pending[1].url, /quantity=10$/)
    h.pending[1].resolve({ ok: true, json: async () => ({ pricing: { quantity: 10, unit_price: 64, original_unit_price: 80 } }) })
    await second
    h.pending[0].resolve({ ok: true, json: async () => ({ pricing: { quantity: 1, unit_price: 80 } }) })
    await first
    assert.equal(h.unitPrice.value, 64)
})
test('expired quote cannot display a stale sale price', () => {
    const h = harness()
    h.clock.value = Date.parse('2030-01-01T00:00:00Z')
    assert.equal(h.currentQuote.value, null)
    assert.equal(h.unitPrice.value, null)
})
