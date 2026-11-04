import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import { computed, ref } from 'vue'
function layoutHarness(url) {
    const source = readFileSync(
        new URL('../../resources/js/AdminLayout.vue', import.meta.url),
        'utf8'
    )
    let script = source
        .match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
        .replace(/^import .*$/gm, '')
    script += '\nglobalThis.harness={active,groups,currentTitle,menuOpen};'
    const page = {
        url,
        component: 'Dashboard/Settings',
        props: {
            workspaceNavigation: [
                {
                    label: 'Overview',
                    href: '/dashboard',
                    group: 'Workspace',
                    icon: 'grid',
                },
                {
                    label: 'Payments',
                    href: '/dashboard/settings?panel=payments',
                    group: 'Administration',
                    icon: 'card',
                },
            ],
        },
    }
    const context = {
        computed,
        ref,
        usePage: () => page,
        normalizeBrand: () => ({}),
        URL,
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
    return { ...context.harness, source }
}
test('admin navigation distinguishes sections while ignoring selected merchant identity', () => {
    const h = layoutHarness(
        '/dashboard/settings?panel=payments&merchant_id=merchant'
    )
    assert.equal(h.active('/dashboard/settings?panel=payments'), true)
    assert.equal(h.active('/dashboard/settings?panel=brand'), false)
    assert.equal(h.active('/dashboard'), false)
    assert.equal(h.currentTitle.value, 'Payments')
    assert.equal(h.groups.value.length, 2)
    assert.match(h.source, /aria-controls="admin-navigation"/)
    assert.match(h.source, /SharedFlashNotifications/)
})
test('client and SSR isolate all Dashboard pages from buyer/storefront chrome', () => {
    for (const entry of ['app.ts', 'ssr.ts']) {
        const source = readFileSync(
            new URL('../../resources/js/' + entry, import.meta.url),
            'utf8'
        )
        assert.match(source, /name.startsWith\('Dashboard\/'\)/)
        assert.match(source, /AdminLayout/)
    }
    const index = readFileSync(
        new URL(
            '../../resources/js/Pages/Dashboard/Index.vue',
            import.meta.url
        ),
        'utf8'
    )
    assert.doesNotMatch(index, /<aside class="workspace-sidebar">/)
    assert.match(index, /recordCan\('products', product.id, 'update'\)/)
    assert.match(index, /recordCan\('offers', offer.id, 'decide'\)/)
})

test('legacy settings route highlights its first permitted panel', () => {
    const h = layoutHarness('/dashboard/settings')
    assert.equal(h.active('/dashboard/settings?panel=payments'), true)
    assert.equal(h.currentTitle.value, 'Payments')
})

test('closing mobile navigation returns focus to the visible menu toggle', () => {
    const source = readFileSync(
        new URL('../../resources/js/AdminLayout.vue', import.meta.url),
        'utf8'
    )
    let script = source
        .match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
        .replace(/^import .*$/gm, '')
    script += '\nglobalThis.harness={menuOpen,closeNavigation};'
    let focused = 0
    let visible = true
    const context = {
        computed,
        ref,
        usePage: () => ({ url: '/dashboard', props: {} }),
        normalizeBrand: () => ({}),
        URL,
        document: {
            getElementById: (id) =>
                id === 'admin-menu-toggle'
                    ? {
                          focus: () => focused++,
                          getClientRects: () => (visible ? [{}] : []),
                      }
                    : null,
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
    context.harness.menuOpen.value = true
    context.harness.closeNavigation()
    assert.equal(context.harness.menuOpen.value, false)
    assert.equal(focused, 1)
    context.harness.closeNavigation()
    visible = false
    context.harness.menuOpen.value = true
    context.harness.closeNavigation()
    assert.equal(
        focused,
        1,
        'desktop navigation must not move focus to a hidden toggle'
    )
})
