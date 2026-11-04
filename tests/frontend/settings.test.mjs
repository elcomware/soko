import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import { reactive, ref, computed } from 'vue'

function settingsHarness() {
    const source = readFileSync(new URL('../../resources/js/Pages/Dashboard/Settings.vue', import.meta.url), 'utf8')
    let script = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1].replace(/^import .*$/gm, '')
    script += '\nglobalThis.harness = { addPaymentRoute, saveBrand, savePayments, savePermissions, loadProfile, selected, scopes, brandForm, paymentForm, accessForm };'
    const requests = []
    const props = { currentUser: { id: 'merchant' }, paymentMerchantID: 'merchant', brand: { display_name: 'Initial', version: 1 }, paymentPreference: { provider: 'none', version: 1, country: 'KE', currency: 'KES', routes: [{ country: 'KE', area: '', currency: 'KES', method: 'card', providers: ['stripe'] }] }, profiles: [{ id: 'other', grants: [], version: 1 }], permissionKeys: ['products.read'] }
    const context = { reactive, ref, computed, defineProps: () => props, usePage: () => ({ props: {} }), useForm: (data) => ({ ...data, post: (url, options) => requests.push({ url, options }) }) }
    vm.runInNewContext(ts.transpileModule(script, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.None } }).outputText, context)
    const h = { ...context.harness, props, requests }
    h.selected.value = 'other'
    h.loadProfile()
    return h
}
for (const [save, form, record] of [['saveBrand', 'brandForm', 'brand'], ['savePayments', 'paymentForm', 'paymentPreference'], ['savePermissions', 'accessForm', 'profiles']]) {
    test(`${save} preserves stale version after rejected redirect, including retry`, () => {
        const h = settingsHarness()
        h[save]()
        const updated = record === 'profiles' ? h.props.profiles[0] : h.props[record]
        updated.version = 2
        h.requests[0].options.onSuccess({ props: { flash: { error: 'Settings changed. Reload before saving.' } } })
        assert.equal(h[form].version, 1)
        h[save]()
        assert.equal(h[form].version, 1)
        h.requests[1].options.onSuccess({ props: { flash: {} } })
        assert.equal(h[form].version, 1)
    })
    test(`${save} advances version only after confirmed save`, () => {
        const h = settingsHarness()
        h[save]()
        const updated = record === 'profiles' ? h.props.profiles[0] : h.props[record]
        updated.version = 2
        h.requests[0].options.onSuccess({ props: { flash: { success: 'Saved.' } } })
        assert.equal(h[form].version, 2)
    })
}

test('regional route editor isolates props, bounds draft rows and preserves failed versions', () => {
    const h = settingsHarness()
    h.paymentForm.routes[0].providers.push('flutterwave')
    assert.equal(h.props.paymentPreference.routes[0].providers.length, 1)
    h.addPaymentRoute()
    assert.equal(h.paymentForm.routes.length, 2)
    assert.equal(h.paymentForm.routes[1].country, 'KE')
    h.paymentForm.processing = true
    h.addPaymentRoute()
    assert.equal(h.paymentForm.routes.length, 2)
    h.paymentForm.processing = false
    while (h.paymentForm.routes.length < 50) h.addPaymentRoute()
    h.addPaymentRoute()
    assert.equal(h.paymentForm.routes.length, 50)
    h.savePayments()
    assert.equal(h.requests[0].url, '/dashboard/settings/payments?panel=payments&merchant_id=merchant')
})
