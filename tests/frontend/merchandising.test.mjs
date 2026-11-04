import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import { computed, ref, reactive } from 'vue'
import {
    normalizeSectionOrder,
    storefrontSectionKeys,
} from '../../resources/js/design/sectionOrder.ts'
function harness() {
    const source = readFileSync(
        new URL(
            '../../resources/js/Pages/Dashboard/Merchandising.vue',
            import.meta.url
        ),
        'utf8'
    )
    let script = source
        .match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
        .replace(/^import[\s\S]*?(?:from\s*)?['"][^'"]+['"]\s*$/gm, '')
    script +=
        '\nglobalThis.harness={form,save,addPin,choices,moveSection,resetSectionOrder,orderAnnouncement,openPlacement,applyPlacement,discardPlacement,editingKey,placementError,placementDirty,draftDirty,discardDraft};'
    const placement = {
        enabled: true,
        title: 'Test',
        strategy: 'manual',
        product_ids: [],
        category: '',
        fallback: 'hide',
        limit: 2,
    }
    const props = {
        settings: {
            version: 1,
            hero: structuredClone(placement),
            hero_side: structuredClone(placement),
            discovery: structuredClone(placement),
            category_navigation: { enabled: true, title: 'Categories' },
        },
        products: [
            {
                id: 'p',
                title: 'Physical',
                category: 'Machines',
                digital: false,
                stock: 5,
            },
            {
                id: 'd',
                title: 'Digital',
                category: 'Software',
                digital: true,
                stock: -1,
            },
        ],
    }
    const requests = [],
        recoveries = []
    const context = {
        useAdminDraftRecovery: (options) => recoveries.push(options),
        computed,
        ref,
        normalizeSectionOrder,
        storefrontSectionKeys,
        defineProps: () => props,
        usePage: () => ({ props: {} }),
        useForm: (data) =>
            reactive({
                ...data,
                post: (url, options) => requests.push({ url, options }),
            }),
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
    return { ...context.harness, props, requests, recoveries }
}
test('editor filters digital choices and keeps unsaved pins separate from server props', () => {
    const h = harness()
    assert.equal(h.choices('discovery').length, 1)
    assert.equal(h.choices('discovery')[0].id, 'd')
    h.addPin('hero')
    assert.equal(h.form.hero.product_ids.length, 1)
    assert.equal(h.props.settings.hero.product_ids.length, 0)
    h.form.processing = true
    h.addPin('hero')
    assert.equal(h.form.hero.product_ids.length, 1)
})
test('publication conflict does not pair stale edits with a fresh version', () => {
    const h = harness()
    h.form.hero.title = 'Draft'
    h.save()
    h.props.settings.version = 2
    h.requests[0].options.onSuccess({ props: { flash: { error: 'Conflict' } } })
    assert.equal(h.form.version, 1)
    assert.equal(h.form.hero.title, 'Draft')
    h.requests[0].options.onSuccess({
        props: { flash: { success: 'Published' } },
    })
    assert.equal(h.form.version, 2)
})

test('section ordering moves actual draft sequence, respects bounds and pending state', () => {
    const h = harness()
    assert.equal(h.form.section_order[0], 'showcase')
    h.moveSection(0, -1)
    assert.equal(h.form.section_order[0], 'showcase')
    h.moveSection(0, 1)
    assert.equal(h.form.section_order[0], 'category_navigation')
    assert.match(h.orderAnnouncement.value, /position 2/)
    h.form.processing = true
    h.moveSection(1, 1)
    assert.equal(h.form.section_order[1], 'showcase')
    h.form.processing = false
    h.resetSectionOrder()
    assert.equal(h.form.section_order[0], 'showcase')
    assert.ok(h.form.section_order.includes('catalog'))
})
test('malformed storefront section order falls back without dropping catalog', () => {
    for (const bad of [
        undefined,
        [],
        ['catalog'],
        ['catalog', 'catalog', 'showcase', 'discovery', 'sourcing'],
        ['showcase', 'category_navigation', 'catalog', 'discovery', '<script>'],
    ])
        assert.deepEqual(normalizeSectionOrder(bad), [...storefrontSectionKeys])
    const custom = [
        'catalog',
        'showcase',
        'category_navigation',
        'discovery',
        'sourcing',
    ]
    assert.deepEqual(normalizeSectionOrder(custom), custom)
})

test('placement drawer stages changes, validates pins, and discards without publication', () => {
    const h = harness()
    h.openPlacement('hero')
    h.form.hero.title = 'Unapplied'
    assert.equal(h.placementDirty.value, true)
    h.discardPlacement()
    assert.equal(h.form.hero.title, 'Test')
    assert.equal(h.requests.length, 0)
    h.openPlacement('hero')
    h.form.hero.product_ids = ['p', 'p']
    h.applyPlacement()
    assert.match(h.placementError.value, /only once/)
    assert.equal(h.editingKey.value, 'hero')
    h.form.hero.product_ids = ['p']
    h.applyPlacement()
    assert.equal(h.editingKey.value, null)
    assert.equal(h.draftDirty.value, true)
    assert.equal(h.requests.length, 0)
    h.discardDraft()
    assert.equal(h.form.hero.product_ids.length, 0)
})
test('Enter-submit inside an open placement stages it instead of publishing', () => {
    const h = harness()
    h.openPlacement('hero')
    h.form.hero.title = 'Staged'
    h.save()
    assert.equal(h.requests.length, 0)
    assert.equal(h.editingKey.value, null)
    assert.equal(h.form.hero.title, 'Staged')
})

test('publish snapshots survive late edits and use acknowledged normalized settings', () => {
    const h = harness()
    h.form.hero.title = 'Submitted'
    h.save()
    h.form.hero.title = 'Later edit'
    const acknowledged = {
        ...structuredClone(h.props.settings),
        hero: { ...h.props.settings.hero, title: 'Submitted' },
        section_order: [...h.form.section_order],
        version: 2,
    }
    h.requests[0].options.onSuccess({
        props: { flash: { success: 'Published' }, settings: acknowledged },
    })
    assert.equal(h.form.hero.title, 'Later edit')
    assert.equal(h.draftDirty.value, true)
    h.form.processing = true
    h.discardDraft()
    assert.equal(h.form.hero.title, 'Later edit')
    h.form.processing = false
    h.discardDraft()
    assert.equal(h.form.hero.title, 'Submitted')
    h.form.hero.title = ' Normalized '
    h.save()
    h.requests[1].options.onSuccess({
        props: {
            flash: { success: 'Published' },
            settings: {
                ...acknowledged,
                hero: { ...acknowledged.hero, title: 'Normalized' },
                version: 3,
            },
        },
    })
    assert.equal(h.form.hero.title, 'Normalized')
    assert.equal(h.draftDirty.value, false)
})
