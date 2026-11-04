import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import * as Vue from 'vue'
function tableHarness(overrides = {}) {
    const source = readFileSync(
        new URL(
            '../../resources/js/components/AdminTable.vue',
            import.meta.url
        ),
        'utf8'
    )
    let script = source
        .match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
        .replace(/^import .*$/gm, '')
    script +=
        '\nglobalThis.h={widths,height,columnWidth,resizeColumnKey,resizeHeightKey,startResize,moveResize,endResize,cancelResize,sortColumn};'
    const props = {
        caption: 'Orders',
        busy: false,
        columns: [
            {
                key: 'name',
                label: 'Name',
                width: 200,
                minWidth: 100,
                maxWidth: 400,
                sortable: true,
            },
        ],
        sort: { key: 'name', direction: 'asc' },
        resizableColumns: true,
        resizableHeight: true,
        viewportHeight: 400,
        ...overrides,
    }
    const emitted = []
    const context = {
        ...Vue,
        onBeforeUnmount: () => {},
        withDefaults: (x) => x,
        defineProps: () => props,
        defineEmits:
            () =>
            (...x) =>
                emitted.push(x),
        useId: () => 'id',
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
    return { ...context.h, props, emitted, source }
}
const keyEvent = (key, shiftKey = false) => ({
    key,
    shiftKey,
    preventDefault() {},
})
test('table keyboard resizing respects column and viewport boundaries', () => {
    const h = tableHarness()
    const col = h.props.columns[0]
    h.resizeColumnKey(col, keyEvent('End'))
    assert.equal(h.columnWidth(col), 400)
    h.resizeColumnKey(col, keyEvent('ArrowRight'))
    assert.equal(h.columnWidth(col), 400)
    h.resizeColumnKey(col, keyEvent('Home'))
    assert.equal(h.columnWidth(col), 100)
    h.resizeColumnKey(col, keyEvent('ArrowRight', true))
    assert.equal(h.columnWidth(col), 140)
    h.resizeHeightKey(keyEvent('Home'))
    assert.equal(h.height.value, 240)
    h.resizeHeightKey(keyEvent('End'))
    assert.equal(h.height.value, 900)
})
test('pointer resize captures primary pointer, clamps, restores on cancel and emits final dimensions', () => {
    const h = tableHarness()
    let captured
    const target = {
        setPointerCapture: (id) => (captured = id),
        releasePointerCapture() {},
        hasPointerCapture: () => true,
    }
    const event = {
        pointerId: 7,
        button: 0,
        isPrimary: true,
        clientX: 200,
        clientY: 400,
        currentTarget: target,
        preventDefault() {},
    }
    h.startResize('name', event)
    assert.equal(captured, 7)
    h.moveResize({ ...event, clientX: 900 })
    assert.equal(h.columnWidth(h.props.columns[0]), 400)
    h.cancelResize()
    assert.equal(h.columnWidth(h.props.columns[0]), 200)
    h.startResize(null, event)
    h.moveResize({ ...event, clientY: 500 })
    h.endResize(event)
    assert.equal(h.height.value, 500)
    assert.ok(
        h.emitted.some((x) => x[0] === 'update:viewportHeight' && x[1] === 500)
    )
})
test('busy and opted-out features cannot sort or resize', () => {
    const h = tableHarness({ busy: true })
    const col = h.props.columns[0]
    h.sortColumn(col)
    h.resizeColumnKey(col, keyEvent('End'))
    assert.equal(h.emitted.length, 0)
    assert.equal(h.columnWidth(col), 200)
    h.props.busy = false
    h.sortColumn(col)
    assert.equal(h.emitted[0][1].direction, 'desc')
    h.props.resizableColumns = false
    h.resizeColumnKey(col, keyEvent('End'))
    assert.equal(h.columnWidth(col), 200)
})

test('pagination validates bounds and page size instead of implying client lazy loading', () => {
    const source = readFileSync(
        new URL(
            '../../resources/js/components/AdminPagination.vue',
            import.meta.url
        ),
        'utf8'
    )
    let script = source
        .match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
        .replace(/^import .*$/gm, '')
    script += '\nglobalThis.h={go,changeSize};'
    const props = {
        page: 2,
        perPage: 25,
        total: 60,
        lastPage: 3,
        from: 26,
        to: 50,
        busy: false,
        pageSizes: [25, 50, 100],
    }
    const emitted = []
    const context = {
        ...Vue,
        onBeforeUnmount: () => {},
        withDefaults: (x) => x,
        defineProps: () => props,
        defineEmits:
            () =>
            (...x) =>
                emitted.push(x),
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
    h.go(0)
    h.go(4)
    h.go(2)
    assert.equal(emitted.length, 0)
    h.go(3)
    assert.deepEqual(Array.from(emitted[0]), ['page', 3])
    h.changeSize({ target: { value: '1000' } })
    assert.equal(emitted.length, 1)
    h.changeSize({ target: { value: '50' } })
    assert.equal(emitted[1][1], 50)
    props.busy = true
    h.go(1)
    assert.equal(emitted.length, 2)
})

test('bulk selection cannot include hidden, duplicate or ineligible records', async () => {
    const { useAdminSelection } =
        await import('../../resources/js/composables/useAdminSelection.ts')
    const loaded = Vue.ref(['a', 'b', 'b'])
    const busy = Vue.ref(false)
    const h = useAdminSelection(
        () => loaded.value,
        () => busy.value
    )
    h.toggle('unknown', true)
    assert.equal(h.selected.value.length, 0)
    h.toggleAll(true)
    assert.deepEqual(h.selected.value, ['a', 'b'])
    assert.equal(h.all.value, true)
    loaded.value = ['b', 'c']
    await Vue.nextTick()
    assert.deepEqual(h.selected.value, ['b'])
    assert.equal(h.mixed.value, true)
    busy.value = true
    h.toggle('c', true)
    h.toggleAll(false)
    assert.deepEqual(h.selected.value, ['b'])
    busy.value = false
    h.clear()
    assert.equal(h.selected.value.length, 0)
})

test('rendered native table retains slot compatibility and opts into controls only when configured', async () => {
    const { createRequire } = await import('node:module')
    const require = createRequire(import.meta.url)
    const { parse, compileScript } = require('@vue/compiler-sfc')
    const { renderToString } = require('@vue/server-renderer')
    const source = readFileSync(
        new URL(
            '../../resources/js/components/AdminTable.vue',
            import.meta.url
        ),
        'utf8'
    )
    const { descriptor } = parse(source)
    const compiled = compileScript(descriptor, {
        id: 'admin-table-render',
        inlineTemplate: true,
        templateOptions: { ssr: true },
    })
    const js = ts.transpileModule(compiled.content, {
        compilerOptions: {
            target: ts.ScriptTarget.ES2022,
            module: ts.ModuleKind.CommonJS,
        },
    }).outputText
    const module = { exports: {} }
    new Function('require', 'module', 'exports', js)(
        require,
        module,
        module.exports
    )
    const Component = module.exports.default
    const old = await renderToString(
        Vue.createSSRApp({
            render: () =>
                Vue.h(
                    Component,
                    { caption: 'Legacy table' },
                    {
                        header: () => Vue.h('tr', [Vue.h('th', 'Product')]),
                        default: () => Vue.h('tr', [Vue.h('td', 'Actual row')]),
                    }
                ),
        })
    )
    assert.match(
        old.replace(/ data-v-[a-z-]+/g, '').replace(/<!--[\s\S]*?-->/g, ''),
        /<thead><tr><th>Product<\/th><\/tr><\/thead>/
    )
    assert.match(old, /Actual row/)
    assert.doesNotMatch(old, /role="separator"|Table height|>Filters</)
    const advanced = await renderToString(
        Vue.createSSRApp({
            render: () =>
                Vue.h(
                    Component,
                    {
                        caption: 'Orders',
                        columns: [
                            {
                                key: 'date',
                                label: 'Received',
                                sortable: true,
                                width: 200,
                            },
                        ],
                        sort: { key: 'date', direction: 'desc' },
                        resizableColumns: true,
                        resizableHeight: true,
                        viewportHeight: 420,
                    },
                    {
                        filters: () =>
                            Vue.h('input', { 'aria-label': 'Status' }),
                    }
                ),
        })
    )
    assert.match(advanced, /aria-sort="descending"/)
    assert.match(advanced, /Resize Received column/)
    assert.match(advanced, /max-height:min\(420px, 75dvh\)/)
    assert.match(advanced, /<details[^>]*>/)
    assert.doesNotMatch(advanced, /<details[^>]* open/)
    assert.match(advanced, /Resize table viewport height/)
})

test('a column named height is independent from viewport resizing', () => {
    const h = tableHarness({
        columns: [
            {
                key: 'height',
                label: 'Product height',
                width: 200,
                minWidth: 100,
                maxWidth: 400,
            },
        ],
        resizableHeight: false,
    })
    const event = {
        pointerId: 1,
        button: 0,
        isPrimary: true,
        clientX: 200,
        clientY: 400,
        currentTarget: {
            setPointerCapture() {},
            hasPointerCapture: () => false,
        },
        preventDefault() {},
    }
    h.startResize('height', event)
    h.moveResize({ ...event, clientX: 250 })
    h.endResize(event)
    assert.equal(h.columnWidth(h.props.columns[0]), 250)
    assert.equal(h.height.value, 400)
})
