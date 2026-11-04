import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import { computed, nextTick, ref, markRaw } from 'vue'

function menuHarness() {
    const source = readFileSync(
        new URL(
            '../../resources/js/components/AdminActionMenu.vue',
            import.meta.url
        ),
        'utf8'
    )
    let script = source
        .match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
        .replace(/^import .*$/gm, '')
    script +=
        '\nglobalThis.harness = { show, close, select, keydown, outside, scrolled, trigger, menu, open, position };'
    const props = {
        label: 'Product actions',
        items: [
            { id: 'edit', label: 'Edit' },
            { id: 'delete', label: 'Delete', disabled: true },
            { id: 'copy', label: 'Copy' },
        ],
    }
    const listeners = new Map()
    const document = {
        activeElement: null,
        addEventListener: (name, fn) => listeners.set(name, fn),
        removeEventListener: (name) => listeners.delete(name),
    }
    const window = {
        innerWidth: 320,
        innerHeight: 300,
        addEventListener: (name, fn) => listeners.set(name, fn),
        removeEventListener: (name) => listeners.delete(name),
    }
    const emitted = []
    const context = {
        computed,
        nextTick,
        ref,
        useId: () => 'menu-id',
        watch: () => {},
        onBeforeUnmount: () => {},
        defineProps: () => props,
        defineEmits:
            () =>
            (...args) =>
                emitted.push(args),
        document,
        window,
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
    const h = context.harness
    const buttons = ['Edit', 'Copy'].map((textContent) => ({
        textContent,
        focus() {
            document.activeElement = this
        },
    }))
    const trigger = {
        focus() {
            document.activeElement = this
        },
        contains: () => false,
        getBoundingClientRect: () => ({ right: 315, bottom: 280, top: 236 }),
    }
    let visible = false
    const menu = {
        offsetWidth: 180,
        offsetHeight: 100,
        contains: (node) => buttons.includes(node),
        querySelectorAll: () => buttons,
        showPopover: () => {
            visible = true
        },
        hidePopover: () => {
            visible = false
        },
    }
    h.trigger.value = markRaw(trigger)
    h.menu.value = markRaw(menu)
    const key = (value) => {
        let prevented = false
        h.keydown({
            key: value,
            preventDefault: () => {
                prevented = true
            },
            stopPropagation() {},
        })
        return prevented
    }
    return {
        ...h,
        props,
        emitted,
        document,
        window,
        listeners,
        buttons,
        triggerElement: trigger,
        key,
        visible: () => visible,
        source,
    }
}
test('menu opens in viewport and skips disabled items with keyboard wrap/Home/End/typeahead', async () => {
    const h = menuHarness()
    await h.show()
    assert.equal(h.open.value, true)
    assert.equal(h.visible(), true)
    assert.equal(h.position.value.left, '132px')
    assert.equal(h.position.value.top, '132px')
    assert.equal(h.document.activeElement, h.buttons[0])
    h.key('ArrowUp')
    assert.equal(h.document.activeElement, h.buttons[1])
    h.key('ArrowDown')
    assert.equal(h.document.activeElement, h.buttons[0])
    h.key('End')
    assert.equal(h.document.activeElement, h.buttons[1])
    h.key('Home')
    assert.equal(h.document.activeElement, h.buttons[0])
    h.key('c')
    assert.equal(h.document.activeElement, h.buttons[1])
    assert.equal(h.key('Escape'), true)
    assert.equal(h.open.value, false)
    assert.equal(h.visible(), false)
    assert.equal(h.document.activeElement, h.triggerElement)
    assert.equal(h.listeners.size, 0)
})
test('menu actions cannot execute disabled/absent items and restore focus before emitting', async () => {
    const h = menuHarness()
    await h.show(true)
    assert.equal(h.document.activeElement, h.buttons[1])
    h.select('delete')
    h.select('unknown')
    assert.equal(h.emitted.length, 0)
    h.select('edit')
    assert.equal(h.emitted[0][1], 'edit')
    assert.equal(h.open.value, false)
    assert.equal(h.document.activeElement, h.triggerElement)
    h.props.disabled = true
    await h.show()
    h.select('edit')
    assert.equal(h.open.value, false)
    assert.equal(h.emitted.length, 1)
})
test('outside focus, scroll, Tab and unmount cleanup do not leave a floating popup', async () => {
    const h = menuHarness()
    await h.show()
    h.outside({ target: h.buttons[0] })
    assert.equal(h.open.value, true)
    h.outside({ target: {} })
    assert.equal(h.open.value, false)
    assert.equal(h.listeners.size, 0)
    await h.show()
    h.scrolled({ target: h.buttons[0] })
    assert.equal(h.open.value, true)
    h.scrolled({ target: {} })
    assert.equal(h.open.value, false)
    await h.show()
    assert.equal(h.key('Tab'), false, 'native tab navigation must continue')
    assert.equal(h.open.value, false)
    assert.match(h.source, /onBeforeUnmount\(\(\) => close\(\)\)/)
    assert.doesNotMatch(h.source, /<Teleport/)
    assert.match(h.source, /popover="manual"/)
})
test('shared table preserves native semantics, caption and keyboard horizontal scrolling', () => {
    const source = readFileSync(
        new URL(
            '../../resources/js/components/AdminTable.vue',
            import.meta.url
        ),
        'utf8'
    )
    assert.match(
        source,
        /role="region"[\s\S]*?:aria-labelledby="captionId"[\s\S]*?tabindex="0"/
    )
    assert.match(source, /<caption :id="captionId">/)
    assert.match(source, /<thead>/)
    assert.match(source, /<tbody>/)
    assert.match(source, /overflow-x: auto/)
    assert.match(source, /font-size: 13px/)
    assert.doesNotMatch(source, /min-width: 700px/)
})
