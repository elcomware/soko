import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import { ref, nextTick, markRaw } from 'vue'
const source = readFileSync(
    new URL('../../resources/js/components/AdminDrawer.vue', import.meta.url),
    'utf8'
)
function harness() {
    const props = {
        open: true,
        title: 'Edit product',
        dirty: false,
        busy: false,
    }
    const emitted = []
    const listeners = new Map()
    let beforeVisit
    let visits = []
    let openChanged
    let unmount
    const document = {
        body: { style: { overflow: 'auto' } },
        activeElement: null,
    }
    const focusable = () =>
        markRaw({
            isConnected: true,
            focus() {
                document.activeElement = this
            },
        })
    const opener = focusable()
    document.activeElement = opener
    const native = Object.assign(focusable(), {
        open: false,
        showModal() {
            this.open = true
        },
        close() {
            this.open = false
        },
    })
    let script = source
        .match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
        .replace(/^import .*$/gm, '')
    script +=
        '\nglobalThis.h = { dialog, heading, keepButton, requestClose, keepEditing, discard, escape, confirming, beforeUnload };'
    const context = {
        ref,
        nextTick,
        useId: () => 'test-id',
        defineProps: () => props,
        withDefaults: (x) => x,
        defineEmits:
            () =>
            (...x) =>
                emitted.push(x),
        defineExpose() {},
        onMounted: (callback) => callback(),
        router: {
            on: (_event, callback) => {
                beforeVisit = callback
                return () => {}
            },
            visit: (...args) => visits.push(args),
        },
        watch: (_get, callback) => {
            openChanged = callback
        },
        onBeforeUnmount: (callback) => {
            unmount = callback
        },
        document,
        window: {
            addEventListener: (key, fn) => listeners.set(key, fn),
            removeEventListener: (key) => listeners.delete(key),
        },
    }
    vm.runInNewContext(
        ts.transpileModule(script, {
            compilerOptions: {
                target: ts.ScriptTarget.ES2022,
                module: ts.ModuleKind.None,
            },
        }).outputText,
        context
    )
    const h = context.h
    h.dialog.value = native
    h.heading.value = focusable()
    h.keepButton.value = focusable()
    return {
        ...h,
        props,
        emitted,
        visits,
        beforeVisit: (event) => beforeVisit(event),
        document,
        opener,
        native,
        listeners,
        openChanged: (...x) => openChanged(...x),
        unmount: () => unmount(),
    }
}
test('drawer uses native modal containment, heading focus, scroll lock and opener return', async () => {
    const h = harness()
    await h.openChanged(true)
    assert.equal(h.native.open, true)
    assert.equal(h.document.body.style.overflow, 'hidden')
    assert.equal(h.document.activeElement, h.heading.value)
    assert.equal(h.listeners.size, 1)
    h.props.open = false
    await h.openChanged(false)
    assert.equal(h.native.open, false)
    assert.equal(h.document.body.style.overflow, 'auto')
    assert.equal(h.document.activeElement, h.opener)
    assert.equal(h.listeners.size, 0)
})
test('dirty close requires explicit discard; Escape from warning keeps edits', async () => {
    const h = harness()
    await h.openChanged(true)
    h.props.dirty = true
    await h.requestClose()
    assert.equal(h.confirming.value, true)
    assert.equal(h.emitted.length, 0)
    assert.equal(h.document.activeElement, h.keepButton.value)
    h.escape()
    await nextTick()
    assert.equal(h.confirming.value, false)
    assert.equal(h.emitted.length, 0)
    await h.requestClose()
    h.discard()
    assert.equal(h.emitted[0][0], 'close')
})
test('busy blocks clean close, Escape and discard; clean close emits immediately', async () => {
    const h = harness()
    h.props.busy = true
    await h.requestClose()
    h.escape()
    h.discard()
    assert.equal(h.emitted.length, 0)
    h.props.busy = false
    await h.requestClose()
    assert.equal(h.emitted[0][0], 'close')
})
test('dirty browser unload warns and unmount restores scroll without leaving listeners', async () => {
    const h = harness()
    await h.openChanged(true)
    let prevented = 0
    const event = { preventDefault: () => prevented++, returnValue: undefined }
    h.beforeUnload(event)
    assert.equal(prevented, 0)
    h.props.dirty = true
    h.beforeUnload(event)
    assert.equal(prevented, 1)
    assert.equal(event.returnValue, '')
    h.unmount()
    assert.equal(h.native.open, false)
    assert.equal(h.document.body.style.overflow, 'auto')
    assert.equal(h.listeners.size, 0)
})
test('drawer stays within brand ancestry and has responsive bottom-sheet and safe footer semantics', () => {
    assert.doesNotMatch(source, /<Teleport/)
    assert.match(source, /@cancel.prevent="escape"/)
    assert.match(source, /max-width: 640px/)
    assert.match(source, /height: 92dvh/)
    assert.match(source, /safe-area-inset-bottom/)
    assert.match(source, /prefers-reduced-motion/)
    assert.match(source, /:inert="confirming \|\| undefined"/)
    assert.match(source, /name="footer" :request-close="requestClose"/)
})

test('dirty Inertia GET navigation waits for discard and busy navigation is blocked', async () => {
    const h = harness()
    await h.openChanged(true)
    h.props.dirty = true
    let cancelled = 0
    const visit = { url: '/dashboard?section=orders', method: 'get' }
    h.beforeVisit({ detail: { visit }, preventDefault: () => cancelled++ })
    await nextTick()
    assert.equal(cancelled, 1)
    assert.equal(h.confirming.value, true)
    await h.keepEditing()
    assert.equal(h.visits.length, 0)
    h.beforeVisit({ detail: { visit }, preventDefault: () => cancelled++ })
    await nextTick()
    h.discard()
    assert.equal(h.visits.length, 1)
    assert.equal(h.visits[0][0], visit.url)
    h.props.busy = true
    h.beforeVisit({ detail: { visit }, preventDefault: () => cancelled++ })
    assert.equal(cancelled, 3)
})
