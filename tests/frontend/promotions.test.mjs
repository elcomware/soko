import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import { reactive, computed, ref } from 'vue'
function campaignHarness(canManageAll, products = [], campaigns = []) {
    const source = readFileSync(
        new URL(
            '../../resources/js/Pages/Dashboard/Promotions.vue',
            import.meta.url
        ),
        'utf8'
    )
    let script = source
        .match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
        .replace(/^import .*$/gm, '')
    script +=
        '\n globalThis.harness = { form, edit, reset, resetDraft, busy, error, eligibleProducts, merchants, selectableProducts, targetSearch, editorOpen, dirty, createCampaign, closeEditor, save, currentTime, status, visibleCampaigns, campaignSearch, statusFilter };'
    let mount, unmount
    const clockTimers = []
    const context = {
        onMounted: fn => { mount = fn },
        onBeforeUnmount: fn => { unmount = fn },
        setInterval: (callback, delay) => { clockTimers.push({ callback, delay, cleared: false }); return clockTimers.length },
        clearInterval: id => { clockTimers[id - 1].cleared = true },
        reactive,
        computed,
        ref,
        defineProps: () => ({
            canManageAll,
            actorId: 'merchant-own',
            products,
            campaigns,
        }),
        document: { getElementById: () => null },
        useMarketplaceForm: () => ({
            busy: ref(false),
            error: ref('old error'),
            post: () => {},
            validate: () => true,
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
    return { ...context.harness, mount: () => mount(), unmount: () => unmount(), clockTimers }
}
for (const canManageAll of [true, false])
    test(`campaign reset restores every default (all permission: ${canManageAll})`, () => {
        const h = campaignHarness(canManageAll)
        const defaults = JSON.stringify(h.form)
        const scope = {
            all: false,
            product_ids: ['old-product'],
            categories: ['Clothes'],
            sub_categories: [{ category: 'Clothes', sub_category: 'Shoes' }],
            exclude_product_ids: ['excluded'],
            exclude_categories: ['Clearance'],
            exclude_sub_categories: [
                { category: 'Clothes', sub_category: 'Coats' },
            ],
        }
        h.edit({
            id: 'old',
            name: 'Past campaign',
            kind: 'fixed',
            amount: 77,
            enabled: true,
            vendor_id: 'old-merchant',
            updated_at: 'version',
            starts_at: '2026-12-01T00:00:00Z',
            ends_at: '2026-12-31T00:00:00Z',
            scope,
        })
        h.form.product_ids.push('new-selection')
        assert.deepEqual(
            scope.product_ids,
            ['old-product'],
            'editing must not mutate campaign history'
        )
        h.error.value = 'Failed request'
        h.reset()
        assert.equal(JSON.stringify(h.form), defaults)
        assert.equal(h.error.value, '')
        assert.equal(h.form.vendor_id, canManageAll ? '' : 'merchant-own')
        h.form.categories.push('temporary')
        h.reset()
        assert.equal(
            h.form.categories.length,
            0,
            'reset must allocate fresh target arrays'
        )
    })

test('preview respects merchant, category, subcategory and exclusion scope without changing on search', () => {
    const products = [
        {
            id: 'a',
            title: 'Coat',
            category: 'Clothes',
            sub_category: 'Winter',
            vendor_id: 'one',
            vendor_name: 'One',
        },
        {
            id: 'b',
            title: 'Shirt',
            category: 'Clothes',
            sub_category: 'Summer',
            vendor_id: 'one',
            vendor_name: 'One',
        },
        {
            id: 'c',
            title: 'Other coat',
            category: 'Clothes',
            sub_category: 'Winter',
            vendor_id: 'two',
            vendor_name: 'Two',
        },
    ]
    const h = campaignHarness(true, products)
    assert.equal(h.merchants.value.length, 2)
    h.form.vendor_id = 'one'
    assert.equal(h.eligibleProducts.value.length, 2)
    h.form.exclude_sub_categories = [
        JSON.stringify({ category: 'Clothes', sub_category: 'Summer' }),
    ]
    assert.equal(h.eligibleProducts.value[0].id, 'a')
    assert.equal(h.eligibleProducts.value.length, 1)
    h.targetSearch.value = 'Shirt'
    assert.equal(h.selectableProducts.value.length, 1)
    assert.equal(
        h.eligibleProducts.value[0].id,
        'a',
        'search filters choices, not campaign coverage'
    )
    h.form.all = false
    assert.equal(h.eligibleProducts.value.length, 0)
    h.form.product_ids = ['a', 'c']
    assert.equal(
        h.eligibleProducts.value.length,
        1,
        'explicit products cannot escape merchant boundary'
    )
    h.form.exclude_categories = ['Clothes']
    assert.equal(h.eligibleProducts.value.length, 0)
})

test('in-flight campaign saves cannot switch or cancel the current draft', () => {
    const h = campaignHarness(true)
    h.form.name = 'Pending campaign'
    h.busy.value = true
    h.edit({ id: 'other' })
    h.reset()
    assert.equal(h.form.name, 'Pending campaign')
    h.resetDraft()
    assert.equal(
        h.form.name,
        '',
        'confirmed success must reset even before onFinish clears busy'
    )
})

test('campaign drawer starts clean, tracks edits, and success reset closes it', () => {
    const h = campaignHarness(true)
    assert.equal(h.editorOpen.value, false)
    h.createCampaign()
    assert.equal(h.editorOpen.value, true)
    assert.equal(h.dirty.value, false)
    h.form.name = 'Holiday'
    assert.equal(h.dirty.value, true)
    h.createCampaign()
    assert.equal(h.form.name, 'Holiday', 'another create request must not silently reset an open draft')
    h.busy.value = true
    h.closeEditor()
    assert.equal(h.editorOpen.value, true)
    h.resetDraft()
    assert.equal(h.editorOpen.value, false)
    assert.equal(h.dirty.value, false)
})
test('reactive campaign clock updates status filters at start and exclusive end', () => {
    const c = { name: 'Holiday', vendor_id: '', kind: 'percentage', enabled: true, starts_at: '2030-01-01T00:00:00Z', ends_at: '2030-01-02T00:00:00Z' }
    const h = campaignHarness(true, [], [c])
    h.currentTime.value = Date.parse(c.starts_at) - 1
    h.statusFilter.value = 'Scheduled'
    assert.equal(h.visibleCampaigns.value.length, 1)
    h.currentTime.value = Date.parse(c.starts_at)
    assert.equal(h.visibleCampaigns.value.length, 0)
    h.statusFilter.value = 'Active'
    assert.equal(h.visibleCampaigns.value.length, 1)
    h.currentTime.value = Date.parse(c.ends_at)
    assert.equal(h.status(c), 'Ended')
    assert.equal(h.visibleCampaigns.value.length, 0)
    h.statusFilter.value = 'Ended'
    h.campaignSearch.value = 'holiday'
    assert.equal(h.visibleCampaigns.value.length, 1)
    h.campaignSearch.value = 'other'
    assert.equal(h.visibleCampaigns.value.length, 0)
})

test('campaign status clock is bounded and cleaned up with page', () => {
    const h = campaignHarness(true)
    h.mount()
    assert.equal(h.clockTimers.length, 1)
    assert.equal(h.clockTimers[0].delay, 1000)
    h.unmount()
    assert.equal(h.clockTimers[0].cleared, true)
})
