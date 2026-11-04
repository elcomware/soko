<template>
    <div class="site-container settings-page">
        <header class="workspace-heading">
            <div>
                <p class="eyebrow">Workspace configuration</p>
                <h1>{{ panelTitle }}</h1>
                <p>{{ panelDescription }}</p>
            </div>
            <Link href="/dashboard" class="button button-outline"
                >Back to workspace</Link
            >
        </header>
        <AdminNavigationGuard
            :busy="
                brandForm.processing ||
                paymentForm.processing ||
                accessForm.processing
            "
            :dirty="activeDirty"
            @discard="discardActiveDraft"
        />
        <SokoDialog
            :open="pendingProfile !== null"
            title="Discard permission changes?"
            @close="pendingProfile = null"
        >
            <p>
                You have unsaved permission changes for this account. Keep
                editing or discard them before switching.
            </p>
            <div class="settings-dialog-actions">
                <button
                    type="button"
                    class="button button-outline"
                    @click="pendingProfile = null"
                >
                    Keep editing</button
                ><button
                    type="button"
                    class="button button-navy"
                    @click="confirmProfileSwitch"
                >
                    Discard and switch
                </button>
            </div>
        </SokoDialog>
        <div class="settings-draft-bar" :class="{ 'is-dirty': activeDirty }">
            <span role="status">{{
                activeDirty ? 'Unsaved changes' : 'Saved configuration'
            }}</span
            ><button
                v-if="activeDirty"
                type="button"
                class="text-button"
                :disabled="
                    brandForm.processing ||
                    paymentForm.processing ||
                    accessForm.processing
                "
                @click="showDiscard = true"
            >
                Discard changes
            </button>
        </div>
        <SokoDialog
            :open="showDiscard"
            title="Discard this draft?"
            @close="showDiscard = false"
            ><p>
                Saved settings will be restored. This does not change the live
                store.
            </p>
            <div class="settings-dialog-actions">
                <button
                    type="button"
                    class="button button-outline"
                    @click="showDiscard = false"
                >
                    Keep editing</button
                ><button
                    type="button"
                    class="button button-navy"
                    @click="confirmDiscardDraft"
                >
                    Discard draft
                </button>
            </div></SokoDialog
        >
        <p v-if="flash.error" class="form-error" role="alert">
            {{ flash.error }}
        </p>

        <div
            v-if="Object.keys(activeErrors).length"
            class="form-error"
            role="alert"
            id="settings-validation"
        >
            <strong>Check these fields:</strong>
            <ul>
                <li v-for="(error, field) in activeErrors" :key="field">
                    {{ error }}
                </li>
            </ul>
        </div>
        <section
            v-if="activePanel === 'brand' && canBrand"
            class="workspace-panel"
        >
            <h2>Brand identity</h2>
            <p>
                Changes apply across the storefront and workspace. Fonts are
                served locally.
            </p>
            <div
                class="brand-preview"
                :style="brandTokens(normalizeBrand(brandForm))"
                aria-label="Unsaved brand preview"
            >
                <div>
                    <span class="brand-preview-mark">s</span
                    ><strong>{{ brandForm.display_name }}</strong
                    ><small>{{ brandForm.tagline }}</small>
                </div>
                <span class="brand-preview-action">Shop the collection</span>
                <p>
                    Preview only · Save identity to apply these colors and
                    fonts.
                </p>
            </div>
            <form @submit.prevent="saveBrand()" class="settings-grid">
                <fieldset
                    class="settings-form-lock"
                    :disabled="brandForm.processing"
                >
                    <label
                        >Display name<input
                            v-model="brandForm.display_name"
                            :aria-invalid="
                                brandForm.errors.display_name ? true : undefined
                            "
                            :aria-describedby="
                                brandForm.errors.display_name
                                    ? 'settings-validation'
                                    : undefined
                            "
                            maxlength="50"
                            required /></label
                    ><label
                        >Tagline<input
                            v-model="brandForm.tagline"
                            maxlength="100"
                    /></label>
                    <label
                        >Primary color<input
                            v-model="brandForm.primary_color"
                            type="color" /></label
                    ><label
                        >Accent color<input
                            v-model="brandForm.accent_color"
                            type="color"
                    /></label>
                    <label
                        >Body font<select v-model="brandForm.font_body">
                            <option>Manrope</option>
                            <option value="system">System</option>
                        </select></label
                    ><label
                        >Display font<select v-model="brandForm.font_display">
                            <option>Manrope</option>
                            <option value="system">System</option>
                        </select></label
                    >
                    <fieldset class="settings-wide">
                        <legend>Product card details</legend>
                        <p>
                            Choose what customers see below each product.
                            Ratings are currently hidden because verified
                            reviews are not connected. Selecting rating will
                            show it only when a verified review source is
                            available.
                        </p>
                        <label
                            ><input
                                v-model="brandForm.card_footer_price"
                                type="checkbox"
                            />Price</label
                        >
                        <label
                            ><input
                                v-model="brandForm.card_footer_rating"
                                type="checkbox"
                            />Rating</label
                        >
                        <label
                            ><input
                                v-model="brandForm.card_footer_availability"
                                type="checkbox"
                            />Availability</label
                        >
                        <label
                            >Display mode<select
                                v-model="brandForm.card_footer_mode"
                            >
                                <option value="together">
                                    Show selected details together
                                </option>
                                <option value="rotate">
                                    Rotate selected details
                                </option>
                            </select></label
                        >
                        <label v-if="brandForm.card_footer_mode === 'rotate'"
                            >Seconds per detail<input
                                v-model.number="
                                    brandForm.card_footer_interval_seconds
                                "
                                type="number"
                                min="5"
                                max="60"
                                step="1"
                                required
                        /></label>
                        <p>
                            Uncheck all details to hide the footer.
                            Reduced-motion preferences are respected.
                        </p>
                    </fieldset>
                    <p class="settings-wide">
                        Check contrast and readability after changing colors.
                    </p>
                    <button
                        class="button button-navy"
                        :disabled="brandForm.processing"
                    >
                        {{ brandForm.processing ? 'Saving…' : 'Save identity' }}
                    </button>
                </fieldset>
            </form>
        </section>
        <section
            v-if="activePanel === 'access' && canPermissions"
            class="workspace-panel"
        >
            <h2>Resource permissions</h2>
            <p>
                Grant only the actions needed. “Own” limits access to the user’s
                resources. Empty grants revoke access. You cannot change your
                own permissions.
            </p>
            <label
                >Account<select
                    :value="selected"
                    @change="
                        requestProfile(
                            ($event.target as HTMLSelectElement).value,
                            $event
                        )
                    "
                    :disabled="accessForm.processing"
                >
                    <option value="">Choose an account</option>
                    <option
                        v-for="p in profiles"
                        :key="p.id"
                        :value="p.id"
                        :disabled="p.id === currentUser.id"
                    >
                        {{ p.name }} — {{ p.email
                        }}{{ p.id === currentUser.id ? ' (you)' : '' }}
                    </option>
                </select></label
            >
            <form v-if="selected" @submit.prevent="savePermissions">
                <fieldset
                    class="settings-form-lock"
                    :disabled="accessForm.processing"
                >
                    <p class="settings-notice">
                        {{ Object.values(scopes).filter(Boolean).length }}
                        actions selected. Read access and edit access are
                        separate; grant both when the person needs to browse and
                        change records.
                    </p>
                    <label class="permission-search"
                        >Find an action<input
                            v-model="permissionQuery"
                            type="search"
                            placeholder="Search products, orders, media…"
                    /></label>
                    <p
                        v-if="!filteredPermissionKeys.length"
                        class="settings-notice"
                    >
                        No actions match this search. Existing selections are
                        unchanged.
                    </p>
                    <div class="permission-grid">
                        <label v-for="key in filteredPermissionKeys" :key="key"
                            >{{ permissionLabel(key) }}<small>{{ key }}</small
                            ><select v-model="scopes[key]">
                                <option value="">No access</option>
                                <option
                                    v-if="!platformOnly.includes(key)"
                                    value="own"
                                >
                                    Own resources
                                </option>
                                <option value="all">All resources</option>
                            </select></label
                        >
                    </div>
                    <button
                        class="button button-navy"
                        :disabled="accessForm.processing"
                    >
                        {{
                            accessForm.processing
                                ? 'Saving…'
                                : 'Save permissions'
                        }}
                    </button>
                </fieldset>
            </form>
        </section>
        <section
            v-if="activePanel === 'payments' && canPayments"
            class="workspace-panel"
        >
            <div class="settings-section-title">
                <h2>Merchant payment preferences</h2>
                <span class="configuration-state">Not live</span>
            </div>
            <p>Account: {{ paymentMerchantID }}</p>
            <div v-if="canConfigureAllPayments">
                <label
                    >Merchant account
                    <select
                        v-model="merchantTarget"
                        :disabled="paymentForm.processing"
                    >
                        <option
                            v-for="merchant in merchantChoices"
                            :key="merchant.id"
                            :value="merchant.id"
                        >
                            {{ merchant.company_name || merchant.name }} ·
                            {{ merchant.id }}
                        </option>
                    </select>
                </label>
                <Link
                    :href="
                        '/dashboard/settings?panel=payments&merchant_id=' +
                        encodeURIComponent(merchantTarget)
                    "
                    class="button button-outline"
                    as="button"
                    type="button"
                    :disabled="paymentForm.processing"
                    >Open merchant settings</Link
                >
            </div>
            <p class="settings-notice">
                <strong>Not activated.</strong> These preferences do not collect
                payments. Provider onboarding, credentials, country/currency
                eligibility, settlement and reconciliation must be verified
                before activation.
            </p>
            <form @submit.prevent="savePayments()" class="settings-grid">
                <fieldset
                    class="settings-form-lock"
                    :disabled="paymentForm.processing"
                >
                    <label
                        >Provider<select v-model="paymentForm.provider">
                            <option value="none">Not selected</option>
                            <option value="stripe">Stripe</option>
                            <option value="flutterwave">Flutterwave</option>
                            <option value="manual_bank_transfer">
                                Manual bank transfer
                            </option>
                        </select></label
                    ><label
                        >Merchant country (ISO code)<input
                            v-model="paymentForm.country"
                            placeholder="US"
                            pattern="[A-Z]{2}"
                            maxlength="2"
                            required /></label
                    ><label
                        >Currency (ISO code)<input
                            v-model="paymentForm.currency"
                            pattern="[A-Z]{3}"
                            maxlength="3"
                            required
                    /></label>
                    <fieldset>
                        <legend>Requested methods</legend>
                        <label
                            v-for="m in [
                                'card',
                                'mobile_money',
                                'bank_transfer',
                            ]"
                            :key="m"
                            ><input
                                v-model="paymentForm.methods"
                                type="checkbox"
                                :value="m"
                            />{{ m.replace(/_/g, ' ') }}</label
                        >
                    </fieldset>
                    <fieldset class="settings-wide">
                        <legend>Regional provider routes</legend>
                        <p>
                            Draft choices only. Provider country, currency and
                            account eligibility must be verified before
                            collection. An exact area rule takes priority over
                            its country-wide rule.
                        </p>
                        <section
                            v-for="(route, index) in paymentForm.routes"
                            :key="index"
                            class="settings-route"
                        >
                            <h3>Route {{ index + 1 }}</h3>
                            <label
                                >Country<input
                                    v-model="route.country"
                                    pattern="[A-Z]{2}"
                                    maxlength="2"
                                    required
                            /></label>
                            <label
                                >Area code (optional)<input
                                    v-model="route.area"
                                    pattern="[A-Z0-9-]{0,32}"
                                    maxlength="32"
                                    placeholder="Country-wide"
                            /></label>
                            <label
                                >Currency<input
                                    v-model="route.currency"
                                    pattern="[A-Z]{3}"
                                    maxlength="3"
                                    required
                            /></label>
                            <label
                                >Method<select v-model="route.method">
                                    <option value="card">Card</option>
                                    <option value="mobile_money">
                                        Mobile money
                                    </option>
                                    <option value="bank_transfer">
                                        Bank transfer
                                    </option>
                                </select></label
                            >
                            <label
                                v-for="(_, providerIndex) in route.providers"
                                :key="providerIndex"
                                >Provider preference {{ providerIndex + 1
                                }}<select
                                    v-model="route.providers[providerIndex]"
                                >
                                    <option value="stripe">Stripe</option>
                                    <option value="flutterwave">
                                        Flutterwave
                                    </option>
                                    <option value="manual_bank_transfer">
                                        Manual bank transfer
                                    </option></select
                                ><button
                                    type="button"
                                    :disabled="paymentForm.processing"
                                    @click="
                                        route.providers.splice(providerIndex, 1)
                                    "
                                >
                                    Remove provider
                                </button></label
                            >
                            <button
                                type="button"
                                :disabled="
                                    paymentForm.processing ||
                                    route.providers.length >= 3
                                "
                                @click="addRouteProvider(route)"
                            >
                                Add provider preference
                            </button>
                            <button
                                type="button"
                                :disabled="paymentForm.processing"
                                @click="paymentForm.routes.splice(index, 1)"
                            >
                                Remove route
                            </button>
                        </section>
                        <button
                            type="button"
                            class="button button-outline"
                            :disabled="
                                paymentForm.processing ||
                                paymentForm.routes.length >= 50
                            "
                            @click="addPaymentRoute"
                        >
                            Add regional route
                        </button>
                    </fieldset>
                    <button
                        class="button button-navy"
                        :disabled="paymentForm.processing"
                    >
                        {{
                            paymentForm.processing
                                ? 'Saving…'
                                : 'Save preferences'
                        }}
                    </button>
                </fieldset>
            </form>
        </section>
        <section
            v-if="activePanel === 'brand' && canBrand"
            class="workspace-panel"
        >
            <details>
                <summary>Brand icon library · Preview the full system</summary>
                <SokoIconGallery />
            </details>
        </section>
    </div>
</template>
<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { Link, useForm, usePage } from '@inertiajs/vue3'
import '../../../css/workspace.css'
import SokoIconGallery from '@/components/SokoIconGallery.vue'
import AdminNavigationGuard from '@/components/AdminNavigationGuard.vue'
import SokoDialog from '@/components/SokoDialog.vue'
import { brandTokens, normalizeBrand } from '@/design/brand'
import { useAdminDraftRecovery } from '@/design/adminDraftRecovery'
type PaymentRoute = {
    country: string
    area: string
    currency: string
    method: string
    providers: string[]
}
type Grant = { permission: string; scope: string }
type Profile = {
    id: string
    name: string
    email: string
    grants: Grant[]
    version: number
}
const props = defineProps<{
    currentUser: { id: string }
    canBrand: boolean
    canPermissions: boolean
    canPayments: boolean
    merchantChoices?: { id: string; name: string; company_name: string }[]
    canConfigureAllPayments: boolean
    paymentMerchantID: string
    profiles: Profile[]
    permissionKeys: string[]
    brand: {
        display_name: string
        tagline: string
        primary_color: string
        accent_color: string
        font_body: string
        font_display: string
        card_footer_price: boolean
        card_footer_rating: boolean
        card_footer_availability: boolean
        card_footer_mode: 'together' | 'rotate'
        card_footer_interval_seconds: number
        version: number
    }
    paymentPreference: {
        provider: string
        country: string
        currency: string
        methods: string[]
        routes: PaymentRoute[]
        version: number
    }
}>()
const page = usePage()
const activePanel = computed(() => {
    const requested = new URL(
        page.url || '/dashboard/settings',
        'https://workspace.invalid'
    ).searchParams.get('panel')
    const allowed = {
        brand: props.canBrand,
        access: props.canPermissions,
        payments: props.canPayments,
    }
    if (requested && allowed[requested as keyof typeof allowed])
        return requested
    return (
        Object.keys(allowed).find(
            (key) => allowed[key as keyof typeof allowed]
        ) || ''
    )
})
const panelTitle = computed(
    () =>
        ({
            brand: 'Brand & presentation',
            access: 'Access & permissions',
            payments: 'Payment configuration',
        })[activePanel.value] || 'Settings'
)
const panelDescription = computed(
    () =>
        ({
            brand: 'Manage the identity and product presentation shared across your storefront.',
            access: 'Give each person the resources and actions they need, with explicit scope.',
            payments:
                'Prepare merchant payment routes. Collection remains disabled until verified onboarding.',
        })[activePanel.value] || 'Manage your workspace.'
)
const flash = computed(
    () => (page.props.flash as { error?: string; success?: string }) || {}
)
const brandForm = useForm({ ...props.brand })
const merchantTarget = ref(
    props.paymentMerchantID || props.currentUser?.id || ''
)
const paymentForm = useForm({
    ...props.paymentPreference,
    routes: (props.paymentPreference.routes || []).map((route) => ({
        ...route,
        providers: [...route.providers],
    })),
    enabled: false,
})
const brandBaseline = ref(JSON.stringify(props.brand))
const paymentBaseline = ref(
    JSON.stringify({
        ...props.paymentPreference,
        routes: props.paymentPreference.routes || [],
        enabled: false,
    })
)
const brandDirty = computed(
    () =>
        JSON.stringify(
            Object.fromEntries(
                Object.keys(props.brand).map((key) => [
                    key,
                    (brandForm as any)[key],
                ])
            )
        ) !== brandBaseline.value
)
const paymentDirty = computed(
    () =>
        JSON.stringify(
            Object.fromEntries(
                Object.keys(JSON.parse(paymentBaseline.value)).map((key) => [
                    key,
                    (paymentForm as any)[key],
                ])
            )
        ) !== paymentBaseline.value
)
const activeDirty = computed(() =>
    activePanel.value === 'brand'
        ? brandDirty.value
        : activePanel.value === 'payments'
          ? paymentDirty.value
          : accessDirty.value
)
const showDiscard = ref(false)
const activeErrors = computed(
    () =>
        (activePanel.value === 'brand'
            ? brandForm.errors
            : activePanel.value === 'payments'
              ? paymentForm.errors
              : accessForm.errors) || {}
)
function confirmDiscardDraft() {
    discardActiveDraft()
    showDiscard.value = false
}
function discardActiveDraft() {
    if (brandForm.processing || paymentForm.processing || accessForm.processing)
        return
    if (activePanel.value === 'brand')
        Object.assign(brandForm, JSON.parse(brandBaseline.value))
    else if (activePanel.value === 'payments')
        Object.assign(paymentForm, JSON.parse(paymentBaseline.value))
    else {
        const restored = JSON.parse(accessBaseline.value)
        for (const key of props.permissionKeys)
            scopes[key] = restored[key] || ''
    }
}
function addRouteProvider(route: PaymentRoute) {
    if (paymentForm.processing) return
    const next = ['stripe', 'flutterwave', 'manual_bank_transfer'].find(
        (provider) => !route.providers.includes(provider)
    )
    if (next) route.providers.push(next)
}
function addPaymentRoute() {
    if (paymentForm.processing || paymentForm.routes.length >= 50) return
    paymentForm.routes.push({
        country: paymentForm.country || '',
        area: '',
        currency: paymentForm.currency || 'USD',
        method: 'card',
        providers: ['stripe'],
    })
}
function saveBrand() {
    if (brandForm.processing) return
    const submitted = JSON.parse(
        JSON.stringify(
            Object.fromEntries(
                Object.keys(props.brand).map((key) => [
                    key,
                    (brandForm as any)[key],
                ])
            )
        )
    )
    brandForm.post(
        '/dashboard/settings/brand?panel=brand&merchant_id=' +
            encodeURIComponent(
                props.paymentMerchantID || props.currentUser?.id || ''
            ),
        {
            preserveScroll: true,
            onSuccess: (response) => {
                // Redirects also resolve successfully after server-side rejection.
                // Never pair unsaved form values with a newer concurrency version.
                const result = response.props.flash as
                    { error?: string; success?: string } | undefined
                if (result?.error || !result?.success) return
                const current = Object.fromEntries(
                    Object.keys(submitted).map((key) => [
                        key,
                        (brandForm as any)[key],
                    ])
                )
                const unchanged =
                    JSON.stringify(current) === JSON.stringify(submitted)
                const acknowledged = JSON.parse(
                    JSON.stringify(
                        response.props.brand || {
                            ...submitted,
                            version: props.brand.version,
                        }
                    )
                )
                if (unchanged) Object.assign(brandForm, acknowledged)
                else brandForm.version = acknowledged.version
                brandBaseline.value = JSON.stringify(acknowledged)
            },
        }
    )
}
function savePayments() {
    if (paymentForm.processing) return
    const submitted = JSON.parse(
        JSON.stringify(
            Object.fromEntries(
                Object.keys(JSON.parse(paymentBaseline.value)).map((key) => [
                    key,
                    (paymentForm as any)[key],
                ])
            )
        )
    )
    paymentForm.post(
        '/dashboard/settings/payments?panel=payments&merchant_id=' +
            encodeURIComponent(
                props.paymentMerchantID || props.currentUser?.id || ''
            ),
        {
            preserveScroll: true,
            onSuccess: (response) => {
                // Redirects also resolve successfully after server-side rejection.
                // Never pair unsaved form values with a newer concurrency version.
                const result = response.props.flash as
                    { error?: string; success?: string } | undefined
                if (result?.error || !result?.success) return
                const current = Object.fromEntries(
                    Object.keys(submitted).map((key) => [
                        key,
                        (paymentForm as any)[key],
                    ])
                )
                const unchanged =
                    JSON.stringify(current) === JSON.stringify(submitted)
                const acknowledged = JSON.parse(
                    JSON.stringify({
                        ...submitted,
                        ...((response.props.paymentPreference as any) || {}),
                        version:
                            (response.props.paymentPreference as any)
                                ?.version ?? props.paymentPreference.version,
                    })
                )
                if (unchanged) Object.assign(paymentForm, acknowledged)
                else paymentForm.version = acknowledged.version
                paymentBaseline.value = JSON.stringify(acknowledged)
            },
        }
    )
}
const selected = ref('')
const loadedProfile = ref('')
const pendingProfile = ref<string | null>(null)
const permissionQuery = ref('')
const filteredPermissionKeys = computed(() =>
    props.permissionKeys.filter((key) =>
        `${permissionLabel(key)} ${key}`
            .toLowerCase()
            .includes(permissionQuery.value.trim().toLowerCase())
    )
)
const accessBaseline = ref('{}')
const scopes = reactive<Record<string, string>>({})
const accessDirty = computed(
    () => JSON.stringify(scopes) !== accessBaseline.value
)
const accessForm = useForm({ user_id: '', grants: [] as Grant[], version: 0 })
function permissionLabel(key: string) {
    const [resource, ...actions] = key.split('.')
    const action = actions.join(' ').replace(/_/g, ' ')
    return `${action.charAt(0).toUpperCase() + action.slice(1)} ${resource}`
}
const platformOnly = [
    'permissions.manage',
    'brand.update',
    'content.update',
    'users.read',
    'media.settings.manage',
]
function requestProfile(target: string, event?: Event) {
    if (accessForm.processing) return
    if (accessDirty.value && loadedProfile.value) {
        pendingProfile.value = target
        if (event?.target)
            (event.target as HTMLSelectElement).value = loadedProfile.value
        selected.value = loadedProfile.value
        return
    }
    selected.value = target
    loadProfile()
}
function confirmProfileSwitch() {
    selected.value = pendingProfile.value || ''
    pendingProfile.value = null
    loadProfile()
}
function loadProfile() {
    const p = props.profiles.find((p) => p.id === selected.value)
    for (const key of props.permissionKeys) scopes[key] = ''
    if (p) {
        for (const g of p.grants || []) scopes[g.permission] = g.scope
        accessForm.version = p.version
        accessForm.user_id = p.id
    }
    loadedProfile.value = selected.value
    accessBaseline.value = JSON.stringify(scopes)
}
function savePermissions() {
    if (accessForm.processing) return
    const target = selected.value
    const submittedScopes = JSON.stringify(scopes)
    accessForm.grants = Object.entries(scopes)
        .filter(([, scope]) => scope)
        .map(([permission, scope]) => ({ permission, scope }))
    accessForm.post(
        '/dashboard/settings/permissions?panel=access&merchant_id=' +
            encodeURIComponent(
                props.paymentMerchantID || props.currentUser?.id || ''
            ),
        {
            preserveScroll: true,
            onSuccess: (response) => {
                // Redirects also resolve successfully after server-side rejection.
                // Never pair unsaved form values with a newer concurrency version.
                const result = response.props.flash as
                    { error?: string; success?: string } | undefined
                if (result?.error || !result?.success) return
                if (selected.value !== target) return
                const returned = (
                    response.props.profiles as Profile[] | undefined
                )?.find((p) => p.id === target)
                const p =
                    returned || props.profiles.find((p) => p.id === target)
                const acknowledged = returned
                    ? Object.fromEntries(
                          props.permissionKeys.map((key) => [
                              key,
                              returned.grants.find(
                                  (grant) => grant.permission === key
                              )?.scope || '',
                          ])
                      )
                    : JSON.parse(submittedScopes)
                if (JSON.stringify(scopes) === submittedScopes)
                    Object.assign(scopes, acknowledged)
                if (p) accessForm.version = p.version
                accessBaseline.value = JSON.stringify(acknowledged)
            },
        }
    )
}
const recoveryActor = props.currentUser?.id || ''
const recoveryMerchant = props.paymentMerchantID || recoveryActor
const pickFields = (form: any, baseline: string) =>
    Object.fromEntries(
        Object.keys(JSON.parse(baseline)).map((key) => [key, form[key]])
    )
useAdminDraftRecovery({
    key: `soko:brand-draft:${recoveryActor}`,
    allowed: () => Boolean(recoveryActor && props.canBrand),
    read: () => pickFields(brandForm, brandBaseline.value),
    apply: (draft) => Object.assign(brandForm, draft),
    dirty: () => brandDirty.value,
})
useAdminDraftRecovery({
    key: `soko:payment-draft:${recoveryActor}:${recoveryMerchant}`,
    allowed: () => Boolean(recoveryActor && props.canPayments),
    read: () => pickFields(paymentForm, paymentBaseline.value),
    apply: (draft) => Object.assign(paymentForm, draft, { enabled: false }),
    dirty: () => paymentDirty.value,
})
for (const profile of props.profiles) {
    useAdminDraftRecovery({
        key: `soko:access-draft:${recoveryActor}:${profile.id}`,
        allowed: () =>
            Boolean(
                recoveryActor &&
                props.canPermissions &&
                profile.id !== recoveryActor &&
                props.profiles.some((p) => p.id === profile.id)
            ),
        read: () => ({ version: accessForm.version, scopes: { ...scopes } }),
        apply: (draft) => {
            if (
                activePanel.value !== 'access' ||
                selected.value ||
                !draft.scopes ||
                typeof draft.scopes !== 'object'
            )
                return
            selected.value = profile.id
            loadProfile()
            for (const key of props.permissionKeys) {
                const value = (draft.scopes as Record<string, unknown>)[key]
                if (
                    typeof value === 'string' &&
                    ['', 'own', 'all'].includes(value)
                )
                    scopes[key] = value
            }
            if (Number.isInteger(draft.version))
                accessForm.version = Number(draft.version)
        },
        dirty: () => loadedProfile.value === profile.id && accessDirty.value,
    })
}
</script>
<style scoped>
.settings-page {
    padding-block: 2rem;
    display: grid;
    gap: 1.5rem;
}
.settings-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 1rem;
    margin-top: 1.25rem;
}
.settings-page label {
    min-width: 0;
    overflow-wrap: anywhere;
    display: grid;
    gap: 0.4rem;
    font-weight: 600;
    font-size: 0.875rem;
}
.settings-page input,
.settings-page select {
    min-width: 0;
    min-height: 44px;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    padding: 0.6rem;
    background: var(--soko-surface, #fff);
    color: inherit;
    max-width: 100%;
}
.settings-page input[type='checkbox'] {
    min-height: auto;
}
.settings-page fieldset {
    min-width: 0;
}
.settings-page fieldset label {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
}
.settings-page fieldset label > select,
.settings-page fieldset label > input:not([type='checkbox']) {
    flex: 1 1 10rem;
}
.settings-page fieldset label > button {
    max-width: 100%;
    white-space: normal;
}
.settings-wide {
    grid-column: 1/-1;
}
.permission-grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 1rem;
    margin-block: 1rem;
}
.settings-notice {
    padding: 1rem;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
}
.settings-page h2 {
    margin-bottom: 0.75rem;
}
.settings-page button {
    justify-self: start;
}
@media (max-width: 760px) {
    .settings-grid,
    .permission-grid {
        grid-template-columns: 1fr;
    }
}

.settings-draft-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
    font-size: 12px;
    padding: 10px 14px;
    border: 1px solid var(--soko-line);
    background: var(--soko-surface);
    border-radius: 6px;
}
.settings-draft-bar.is-dirty {
    border-left: 3px solid var(--soko-amber);
}
.settings-dialog-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    justify-content: flex-end;
    margin-top: 20px;
}
.brand-preview {
    padding: 24px;
    border: 1px solid var(--soko-line);
    border-radius: 8px;
    background: var(--soko-canvas);
    font-family: var(--soko-font);
    margin: 20px 0;
}
.brand-preview > div {
    display: flex;
    gap: 12px;
    align-items: center;
    flex-wrap: wrap;
    color: var(--soko-navy);
}
.brand-preview strong {
    font: 800 24px var(--soko-font-display);
}
.brand-preview small {
    color: var(--soko-muted);
}
.brand-preview-mark {
    display: grid;
    place-items: center;
    background: var(--soko-amber);
    color: var(--soko-navy);
    width: 36px;
    height: 40px;
    font-size: 28px;
    font-weight: 800;
    border-radius: 8px 8px 8px 2px;
}
.brand-preview-action {
    display: inline-block;
    background: var(--soko-navy);
    color: var(--soko-surface);
    padding: 8px 12px;
    border-radius: 5px;
    font-size: 12px;
    margin: 18px 0 10px;
}
.brand-preview p {
    font-size: 11px;
    color: var(--soko-muted);
}
.permission-search {
    margin-top: 16px;
    max-width: 420px;
}
.settings-section-title {
    display: flex;
    gap: 12px;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
}
.configuration-state {
    font-size: 11px;
    font-weight: 750;
    padding: 5px 9px;
    background: var(--soko-canvas);
    border: 1px solid var(--soko-line);
    border-radius: 5px;
}
.settings-page summary {
    cursor: pointer;
    font-weight: 750;
    padding: 8px 0;
}
.settings-page .settings-form-lock {
    display: contents;
}
</style>
