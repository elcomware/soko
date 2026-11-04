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
function harness(restored, authorized = true) {
    const saved = []
    const requests = []
    const props = {
        currentUser: { id: 'merchant', role: 'merchant' },
        stats: {},
        products: [
            {
                id: 'p',
                title: 'Server title',
                digital_download_url: 'https://private.example/original',
            },
        ],
        currencies: [],
        regions: [],
        recordCapabilities: { products: { p: { update: authorized } } },
        capabilities: { createProducts: authorized },
    }
    let script = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
    script +=
        '\nglobalThis.h={productOpen,productForm,productRecovered,productDirty,openProduct,closeProduct,saveProduct,busy,formError};'
    const context = {
        exports: {},
        defineProps: () => props,
        URL, URLSearchParams, setTimeout, clearTimeout,
        require: (name) =>
            name === 'vue'
                ? { ...Vue, onBeforeUnmount: () => {} }
                : name === '@inertiajs/vue3'
                  ? {
                        usePage: () => ({
                            url: '/dashboard?section=products',
                            props: {},
                        }),
                        router: {
                            restore: () => restored,
                            remember: (...x) => saved.push(x),
                            post: (...x) => requests.push(x),
                            get() {},
                        },
                    }
                  : name.includes('productDraft')
                    ? { safeProductDraft }
                    : name.includes('useCurrency')
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
    return { ...context.h, saved, requests }
}
test('history recovery excludes delivery secrets, File-like values and arbitrary fields', () => {
    const safe = safeProductDraft({
        title: 'Draft',
        digital_download_url: 'secret',
        api_key: 'secret',
        file: { name: 'image.png' },
        gallery_urls: ['/image.png', { secret: true }],
        price_tiers: [],
    })
    assert.equal(safe.title, 'Draft')
    assert.equal('digital_download_url' in safe, false)
    assert.equal('file' in safe, false)
    assert.equal('api_key' in safe, false)
    assert.deepEqual(safe.gallery_urls, ['/image.png'])
})
test('returning through history restores authorized safe draft and warns about non-restored uploads', async () => {
    const h = harness({
        open: true,
        productId: 'p',
        form: {
            title: 'Unsaved title',
            digital_download_url: 'unpublished-secret',
            price_tiers: [],
        },
    })
    assert.equal(h.productOpen.value, true)
    assert.equal(h.productForm.title, 'Unsaved title')
    assert.equal(h.productRecovered.value, true)
    assert.equal(
        h.productForm.digital_download_url,
        'https://private.example/original'
    )
    h.productForm.title = 'Continued'
    await Vue.nextTick()
    const remembered = h.saved.at(-1)[0]
    assert.equal(remembered.form.title, 'Continued')
    assert.equal('digital_download_url' in remembered.form, false)
    h.busy.value = true
    h.closeProduct()
    assert.equal(h.productOpen.value, true)
    h.busy.value = false
    h.closeProduct()
    assert.equal(h.saved.at(-1)[0], null)
    assert.equal(
        harness({ open: true, productId: 'p', form: { title: 'Draft' } }, false)
            .productOpen.value,
        false
    )
})
test('validation/server errors keep product draft; successful save clears recovery', async () => {
    const h = harness(undefined)
    h.openProduct({ id: 'p', title: 'Original' })
    h.productForm.title = 'Changed'
    const event = { target: { elements: [], querySelectorAll: () => [] } }
    h.saveProduct(event)
    let callbacks = h.requests[0][2]
    callbacks.onSuccess({ props: { flash: { error: 'Stock changed' } } })
    callbacks.onFinish()
    assert.equal(h.productOpen.value, true)
    assert.equal(h.productForm.title, 'Changed')
    assert.equal(h.formError.value, 'Stock changed')
    assert.notEqual(h.saved.at(-1)?.[0], null)
    h.saveProduct(event)
    callbacks = h.requests[1][2]
    callbacks.onError({ title: 'Invalid title' })
    callbacks.onFinish()
    assert.equal(h.productOpen.value, true)
    assert.equal(h.formError.value, 'Invalid title')
    h.saveProduct(event)
    callbacks = h.requests[2][2]
    callbacks.onSuccess({ props: { flash: { success: 'Saved' } } })
    callbacks.onFinish()
    await Vue.nextTick()
    assert.equal(h.productOpen.value, false)
    assert.equal(h.saved.at(-1)[0], null)
    assert.match(
        source,
        /<fieldset\s+class="product-editor-fields"\s+:disabled="busy \|\| productRefreshing"/
    )
})
