import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import * as Vue from 'vue'
import { compile } from '@vue/compiler-dom'
import { renderToString } from '@vue/server-renderer'
import { brandTokens, normalizeBrand } from '../../resources/js/design/brand.ts'
const source = readFileSync(new URL('../../resources/js/components/SokoDialog.vue', import.meta.url), 'utf8')

test('native dialog stays under branded admin DOM ancestor, preserving selector and token inheritance', async () => {
    const template = source.match(/<template>([\s\S]*?)<\/template>/)[1]
    const render = new Function('Vue', compile(template, { mode: 'function' }).code)(Vue)
    const Dialog = Vue.defineComponent({
        props: ['open', 'title'], render,
        components: { SokoIcon: { render: () => Vue.h('svg') } },
        setup: () => ({ titleId: 'test-dialog-title', dialog: Vue.ref(null), emit: () => {} }),
    })
    const tokens = brandTokens(normalizeBrand({ primary_color: '#754321', accent_color: '#36ab91', font_body: 'system' }))
    const app = Vue.createSSRApp({ render: () => Vue.h('div', { class: 'admin-shell', style: tokens },
        [Vue.h('main', { class: 'admin-content' }, [Vue.h(Dialog, { open: true, title: 'Edit product' }, { default: () => Vue.h('input', { class: 'merchant-field' }) })])]) })
    app.config.warnHandler = () => {}
    const html = await renderToString(app)
    assert.match(html, /--soko-navy:#754321/)
    assert.match(html, /--soko-font:system-ui/)
    assert.match(html, /<main class="admin-content">[\s\S]*<dialog[\s\S]*class="merchant-field"[\s\S]*<\/dialog>[\s\S]*<\/main>/)
    assert.doesNotMatch(html, /teleport/i)
    assert.match(html, /aria-labelledby="test-dialog-title"/)
})
function lifecycleHarness(initialOpen = true) {
    let script = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1].replace(/^import .*$/gm, '')
    script += '\n globalThis.harness = { dialog };'
    const props = Vue.reactive({ open: initialOpen, title: 'Product' })
    let cleanup
    const body = { style: { overflow: 'scroll' } }
    const context = { ...Vue, defineProps: () => props, defineEmits: () => () => {},
        useId: () => 'dialog-title', onBeforeUnmount: fn => { cleanup = fn }, document: { body },
    }
    vm.runInNewContext(ts.transpileModule(script, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.None } }).outputText, context)
    return { ...context.harness, props, body, cleanup: () => cleanup() }
}
test('native modal opens in top layer and restores scroll when closed', async () => {
    const h = lifecycleHarness()
    const node = { open: false, shown: 0, closed: 0, showModal() { this.open = true; this.shown++ }, close() { this.open = false; this.closed++ } }
    h.dialog.value = node
    await Vue.nextTick(); await new Promise(resolve => setImmediate(resolve))
    assert.equal(node.shown, 1)
    assert.equal(h.body.style.overflow, 'hidden')
    h.props.open = false
    await Vue.nextTick(); await new Promise(resolve => setImmediate(resolve))
    assert.equal(node.closed, 1)
    assert.equal(h.body.style.overflow, 'scroll')
    h.cleanup()
})
test('missing or unmounted dialog never locks document scroll', async () => {
    const h = lifecycleHarness()
    await Vue.nextTick()
    assert.equal(h.body.style.overflow, 'scroll')
    h.cleanup()
})
