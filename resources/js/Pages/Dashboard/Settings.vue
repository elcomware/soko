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
        <p v-if="flash.error" class="form-error" role="alert">
            {{ flash.error }}
        </p>

        <section
            v-if="activePanel === 'brand' && canBrand"
            class="workspace-panel"
        >
            <h2>Brand identity</h2>
            <p>
                Changes apply across the storefront and workspace. Fonts are
                served locally.
            </p>
            <form @submit.prevent="saveBrand()" class="settings-grid">
                <label
                    >Display name<input
                        v-model="brandForm.display_name"
                        maxlength="50"
                        required /></label
                ><label
                    >Tagline<input v-model="brandForm.tagline" maxlength="100"
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
                        Choose what customers see below each product. Ratings
                        are currently hidden because verified reviews are not
                        connected. Selecting rating will show it only when a
                        verified review source is available.
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
                        Uncheck all details to hide the footer. Reduced-motion
                        preferences are respected.
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
                    v-model="selected"
                    @change="loadProfile"
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
                <p class="settings-notice">
                    {{ Object.values(scopes).filter(Boolean).length }} actions
                    selected. Read access and edit access are separate; grant
                    both when the person needs to browse and change records.
                </p>
                <div class="permission-grid">
                    <label v-for="key in permissionKeys" :key="key"
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
                    {{ accessForm.processing ? 'Saving…' : 'Save permissions' }}
                </button>
            </form>
        </section>
        <section
            v-if="activePanel === 'payments' && canPayments"
            class="workspace-panel"
        >
            <h2>Merchant payment preferences</h2>
            <p>Account: {{ paymentMerchantID }}</p>
            <div v-if="canConfigureAllPayments">
                <label
                    >Merchant account
                    <select v-model="merchantTarget">
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
                        v-for="m in ['card', 'mobile_money', 'bank_transfer']"
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
                        account eligibility must be verified before collection.
                        An exact area rule takes priority over its country-wide
                        rule.
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
                            }}<select v-model="route.providers[providerIndex]">
                                <option value="stripe">Stripe</option>
                                <option value="flutterwave">Flutterwave</option>
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
                            @click="route.providers.push('flutterwave')"
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
                        paymentForm.processing ? 'Saving…' : 'Save preferences'
                    }}
                </button>
            </form>
        </section>
        <section
            v-if="activePanel === 'brand' && canBrand"
            class="workspace-panel"
        >
            <h2>Brand icon library</h2>
            <SokoIconGallery />
        </section>
    </div>
</template>
<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { Link, useForm, usePage } from '@inertiajs/vue3'
import '../../../css/workspace.css'
import SokoIconGallery from '@/components/SokoIconGallery.vue'
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
                brandForm.version = props.brand.version
            },
        }
    )
}
function savePayments() {
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
                paymentForm.version = props.paymentPreference.version
            },
        }
    )
}
const selected = ref('')
const scopes = reactive<Record<string, string>>({})
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
function loadProfile() {
    const p = props.profiles.find((p) => p.id === selected.value)
    for (const key of props.permissionKeys) scopes[key] = ''
    if (p) {
        for (const g of p.grants || []) scopes[g.permission] = g.scope
        accessForm.version = p.version
        accessForm.user_id = p.id
    }
}
function savePermissions() {
    const target = selected.value
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
                const p = props.profiles.find((p) => p.id === selected.value)
                if (p) accessForm.version = p.version
            },
        }
    )
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
    border: 1px solid var(--soko-border);
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
    border: 1px solid var(--soko-border);
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
</style>
