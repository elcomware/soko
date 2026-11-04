import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import { reactive, ref, computed } from 'vue'

function settingsHarness() {
    const source = readFileSync(
        new URL(
            '../../resources/js/Pages/Dashboard/Settings.vue',
            import.meta.url
        ),
        'utf8'
    )
    let script = source
        .match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
        .replace(/^import .*$/gm, '')
    script +=
        '\nglobalThis.harness = { addPaymentRoute, saveBrand, savePayments, savePermissions, loadProfile, selected, scopes, brandForm, paymentForm, accessForm, requestProfile, confirmProfileSwitch, pendingProfile, loadedProfile, accessDirty, brandDirty, paymentDirty, activeDirty, discardActiveDraft, addRouteProvider };'
    const requests = [],
        recoveries = []
    const props = {
        currentUser: { id: 'merchant' },
        paymentMerchantID: 'merchant',
        brand: { display_name: 'Initial', version: 1 },
        paymentPreference: {
            provider: 'none',
            version: 1,
            country: 'KE',
            currency: 'KES',
            routes: [
                {
                    country: 'KE',
                    area: '',
                    currency: 'KES',
                    method: 'card',
                    providers: ['stripe'],
                },
            ],
        },
        profiles: [{ id: 'other', grants: [], version: 1 }],
        permissionKeys: ['products.read'],
    }
    const page = reactive({ url: '/dashboard/settings?panel=brand', props: {} })
    props.canBrand = true
    props.canPermissions = true
    props.canPayments = true
    const context = {
        useAdminDraftRecovery: (options) => recoveries.push(options),
        URL,
        reactive,
        ref,
        computed,
        defineProps: () => props,
        usePage: () => page,
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
    const h = { ...context.harness, props, requests, page, recoveries }
    h.selected.value = 'other'
    h.loadProfile()
    return h
}
for (const [save, form, record] of [
    ['saveBrand', 'brandForm', 'brand'],
    ['savePayments', 'paymentForm', 'paymentPreference'],
    ['savePermissions', 'accessForm', 'profiles'],
]) {
    test(`${save} preserves stale version after rejected redirect, including retry`, () => {
        const h = settingsHarness()
        h[save]()
        const updated =
            record === 'profiles' ? h.props.profiles[0] : h.props[record]
        updated.version = 2
        h.requests[0].options.onSuccess({
            props: {
                flash: { error: 'Settings changed. Reload before saving.' },
            },
        })
        assert.equal(h[form].version, 1)
        h[save]()
        assert.equal(h[form].version, 1)
        h.requests[1].options.onSuccess({ props: { flash: {} } })
        assert.equal(h[form].version, 1)
    })
    test(`${save} advances version only after confirmed save`, () => {
        const h = settingsHarness()
        h[save]()
        const updated =
            record === 'profiles' ? h.props.profiles[0] : h.props[record]
        updated.version = 2
        h.requests[0].options.onSuccess({
            props: { flash: { success: 'Saved.' } },
        })
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
    assert.equal(
        h.requests[0].url,
        '/dashboard/settings/payments?panel=payments&merchant_id=merchant'
    )
})

test('switching account protects unsaved scopes until explicit discard', () => {
    const h = settingsHarness()
    h.page.url = '/dashboard/settings?panel=access'
    h.props.profiles.push({ id: 'next', grants: [], version: 3 })
    h.scopes['products.read'] = 'all'
    assert.equal(h.accessDirty.value, true)
    const event = { target: { value: 'next' } }
    h.requestProfile('next', event)
    assert.equal(h.selected.value, 'other')
    assert.equal(h.pendingProfile.value, 'next')
    assert.equal(event.target.value, 'other')
    assert.equal(h.scopes['products.read'], 'all')
    h.confirmProfileSwitch()
    assert.equal(h.selected.value, 'next')
    assert.equal(h.scopes['products.read'], '')
    assert.equal(h.accessDirty.value, false)
})
test('brand and merchant drafts track unsaved changes and reset without publishing', () => {
    const h = settingsHarness()
    h.brandForm.display_name = 'Draft'
    assert.equal(h.brandDirty.value, true)
    h.discardActiveDraft()
    assert.equal(h.brandForm.display_name, 'Initial')
    assert.equal(h.requests.length, 0)
    h.page.url = '/dashboard/settings?panel=payments'
    h.paymentForm.routes[0].country = 'UG'
    assert.equal(h.paymentDirty.value, true)
    h.discardActiveDraft()
    assert.equal(h.paymentForm.routes[0].country, 'KE')
    assert.equal(h.requests.length, 0)
})
test('provider preference additions never create duplicate options', () => {
    const h = settingsHarness()
    const route = h.paymentForm.routes[0]
    h.addRouteProvider(route)
    h.addRouteProvider(route)
    h.addRouteProvider(route)
    assert.equal(new Set(route.providers).size, 3)
    assert.equal(route.providers.length, 3)
})

test('brand save preserves a late edit as dirty and accepts normalized response only when unchanged', () => {
    const h = settingsHarness()
    h.brandForm.display_name = 'Submitted'
    h.saveBrand()
    h.brandForm.display_name = 'Later edit'
    h.requests[0].options.onSuccess({
        props: {
            flash: { success: 'Saved' },
            brand: { display_name: 'Submitted', version: 2 },
        },
    })
    assert.equal(h.brandForm.display_name, 'Later edit')
    assert.equal(h.brandDirty.value, true)
    h.discardActiveDraft()
    assert.equal(h.brandForm.display_name, 'Submitted')
    h.brandForm.display_name = ' Trimmed '
    h.saveBrand()
    h.requests[1].options.onSuccess({
        props: {
            flash: { success: 'Saved' },
            brand: { display_name: 'Trimmed', version: 3 },
        },
    })
    assert.equal(h.brandForm.display_name, 'Trimmed')
    assert.equal(h.brandDirty.value, false)
})
test('payment and permission late edits remain dirty, while pending discards and switches are blocked', () => {
    const h = settingsHarness()
    h.page.url = '/dashboard/settings?panel=payments'
    h.paymentForm.country = 'UG'
    h.savePayments()
    h.paymentForm.country = 'TZ'
    h.requests[0].options.onSuccess({
        props: {
            flash: { success: 'Saved' },
            paymentPreference: {
                ...h.props.paymentPreference,
                country: 'UG',
                version: 2,
            },
        },
    })
    assert.equal(h.paymentForm.country, 'TZ')
    assert.equal(h.paymentDirty.value, true)
    h.paymentForm.processing = true
    h.discardActiveDraft()
    assert.equal(h.paymentForm.country, 'TZ')
    h.paymentForm.processing = false
    h.page.url = '/dashboard/settings?panel=access'
    h.scopes['products.read'] = 'own'
    h.savePermissions()
    h.scopes['products.read'] = 'all'
    h.requests[1].options.onSuccess({
        props: {
            flash: { success: 'Saved' },
            profiles: [
                {
                    id: 'other',
                    version: 2,
                    grants: [{ permission: 'products.read', scope: 'own' }],
                },
            ],
        },
    })
    assert.equal(h.scopes['products.read'], 'all')
    assert.equal(h.accessDirty.value, true)
    h.accessForm.processing = true
    h.requestProfile('')
    h.discardActiveDraft()
    assert.equal(h.selected.value, 'other')
    assert.equal(h.scopes['products.read'], 'all')
})

test('settings recovery scopes actor and target and restores only authorized permission keys', () => {
    const h = settingsHarness()
    const brand = h.recoveries.find(
        (r) => r.key === 'soko:brand-draft:merchant'
    )
    const payment = h.recoveries.find(
        (r) => r.key === 'soko:payment-draft:merchant:merchant'
    )
    assert.equal(brand.allowed(), true)
    assert.equal(payment.allowed(), true)
    h.props.canPayments = false
    assert.equal(payment.allowed(), false)
    const access = h.recoveries.find(
        (r) => r.key === 'soko:access-draft:merchant:other'
    )
    h.page.url = '/dashboard/settings?panel=access'
    h.selected.value = ''
    access.apply({
        version: 0,
        scopes: { 'products.read': 'own', 'unknown.permission': 'all' },
    })
    assert.equal(h.selected.value, 'other')
    assert.equal(h.scopes['products.read'], 'own')
    assert.equal(h.scopes['unknown.permission'], undefined)
    assert.equal(h.accessForm.version, 0)
    assert.equal(h.accessDirty.value, true)
    h.props.canPermissions = false
    assert.equal(access.allowed(), false)
})

test('brand preview wraps maximum-length unbroken name and tagline without shrinking the mark', () => {
    const source = readFileSync(new URL('../../resources/js/Pages/Dashboard/Settings.vue', import.meta.url), 'utf8')
    const wrapping = source.match(/\.brand-preview strong,\s*\.brand-preview small\s*\{([^}]+)\}/)?.[1]
    assert.ok(wrapping, 'both editable preview text elements share overflow protection')
    assert.match(wrapping, /min-width:\s*0/)
    assert.match(wrapping, /max-width:\s*100%/)
    assert.match(wrapping, /overflow-wrap:\s*anywhere/)
    assert.match(source, /\.brand-preview-mark\s*\{[^}]*flex-shrink:\s*0/)
})
