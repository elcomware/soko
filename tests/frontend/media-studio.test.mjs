import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import { computed, reactive, ref, nextTick } from 'vue'

// Exercise the component's real handlers with a controlled HTTP boundary.
function mediaHarness(props = {}) {
    const source = readFileSync(
        new URL(
            '../../resources/js/components/ProductMediaStudio.vue',
            import.meta.url
        ),
        'utf8'
    )
    let script = source
        .match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
        .replace(/^import .*$/gm, '')
    script +=
        '\n globalThis.harness = { customizeCopy, mockup, mockupPanel, busy, error, notice, previewAsset, previewPanel, activeView, mediaViews, loadLibrary, clearFilters, search, kind, policy };'
    const requests = [],
        events = [],
        emissions = []
    const context = {
        computed,
        reactive,
        ref,
        nextTick,
        onMounted: () => {},
        useId: () => 'media-test',
        defineProps: () => ({ modelValue: '', gallery: [], ...props }),
        defineEmits:
            () =>
            (...args) =>
                emissions.push(args),
        document: {
            getElementById: (id) => ({
                focus: () => events.push(['focus', id]),
            }),
        },
        fetch: (url, options) =>
            new Promise((resolve, reject) =>
                requests.push({ url, options, resolve, reject })
            ),
        URLSearchParams,
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
    h.mockupPanel.value = { open: false }
    h.mockup.logo_data = 'previous-raster-data'
    return { ...h, source, requests, events, emissions }
}
const recipe = {
    template: 'box',
    title: 'Coffee',
    brand: 'Soko',
    primary: '#123456',
    accent: '#654321',
    background: '#ffffff',
}

test('Customize a copy loads the recipe, clears logo, opens and focuses editor, and releases busy state', async () => {
    const h = mediaHarness()
    const pending = h.customizeCopy({ id: 'a/b' })
    assert.equal(h.busy.value, true)
    assert.equal(h.requests[0].url, '/dashboard/media/a%2Fb/recipe')
    assert.equal(h.requests[0].options.method, 'GET')
    assert.equal(h.requests[0].options.credentials, 'same-origin')
    await h.customizeCopy({ id: 'duplicate' })
    assert.equal(
        h.requests.length,
        1,
        'in-flight requests cannot be duplicated'
    )
    h.requests[0].resolve({
        ok: true,
        redirected: false,
        json: async () => ({ recipe }),
    })
    await pending
    for (const [key, value] of Object.entries(recipe))
        assert.equal(h.mockup[key], value)
    assert.equal(h.mockup.logo_data, '')
    assert.equal(h.mockupPanel.value.open, true)
    assert.deepEqual(h.events, [['focus', 'media-test-title']])
    assert.equal(h.busy.value, false)
    assert.equal(h.error.value, '')
    assert.match(h.notice.value, /original image stays unchanged/)
    assert.match(h.notice.value, /logo again/)
    assert.equal(
        h.emissions.length,
        0,
        'loading settings must not replace a product image'
    )
    assert.match(h.source, /ref="mockupPanel"/)
    assert.match(h.source, /@click="customizeCopy\(asset\)"/)
})

test('failed recipe request preserves current editor and releases busy state for retry', async () => {
    const h = mediaHarness()
    const pending = h.customizeCopy({ id: 'forbidden' })
    h.requests[0].resolve({
        ok: false,
        redirected: false,
        json: async () => ({ error: 'Mockup not found.' }),
    })
    await pending
    assert.equal(h.busy.value, false)
    assert.equal(h.error.value, 'Mockup not found.')
    assert.equal(h.mockup.logo_data, 'previous-raster-data')
    assert.equal(h.mockupPanel.value.open, false)
    assert.equal(h.events.length, 0)
    const retry = h.customizeCopy({ id: 'allowed' })
    assert.equal(h.busy.value, true)
    assert.equal(h.requests.length, 2)
    h.requests[1].resolve({
        ok: true,
        redirected: false,
        json: async () => ({ recipe }),
    })
    await retry
    assert.equal(h.busy.value, false)
    assert.equal(h.error.value, '')
    assert.equal(h.mockup.title, 'Coffee')
})

test('standalone preview moves to library and focuses the enlarged preview', async () => {
    const h = mediaHarness({ standalone: true })
    h.activeView.value = 'upload'
    h.previewPanel.value = { focus: () => h.events.push(['preview-focus']) }
    await h.previewAsset({ url: '/media/selected' })
    assert.equal(h.activeView.value, 'library')
    assert.deepEqual(h.emissions, [['update:modelValue', '/media/selected']])
    assert.deepEqual(h.events, [['preview-focus']])
})

test('policy-only delegates land on policy; disabled mockups have no navigation entry', async () => {
    const h = mediaHarness({ standalone: true })
    const pending = h.loadLibrary(1, true)
    h.requests[0].resolve({
        ok: true,
        json: async () => ({
            assets: [],
            page: 1,
            per_page: 24,
            total: 0,
            policy: { allow_mockups: false },
            can_create: false,
            can_manage_policy: true,
        }),
    })
    await pending
    assert.equal(h.activeView.value, 'policy')
    assert.equal(
        h.mediaViews.value.some((view) => view.key === 'mockup'),
        false
    )
})

test('clear filters clears both query and image kind before requesting page one', async () => {
    const h = mediaHarness()
    h.search.value = 'coffee'
    h.kind.value = 'mockup'
    const pending = h.clearFilters()
    assert.equal(h.search.value, '')
    assert.equal(h.kind.value, '')
    assert.equal(h.requests[0].url, '/dashboard/media?page=1&q=&kind=')
    h.requests[0].resolve({
        ok: true,
        json: async () => ({
            assets: [],
            page: 1,
            total: 0,
            per_page: 24,
            can_create: true,
        }),
    })
    await pending
})
