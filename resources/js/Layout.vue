<template>
    <div
        class="soko-app"
        :class="{ 'workspace-layout': page.component === 'Dashboard/Index' }"
    >
        <Head :title="pageTitle" />
        <a href="#main-content" class="skip-link">Skip to content</a>
        <header class="site-header">
            <div class="utility-bar">
                <div class="site-container utility-inner">
                    <span class="utility-message"
                        ><SokoIcon name="globe" />
                        {{
                            (page.props.storefrontContent as any)
                                ?.announcement ??
                            'A world of suppliers. One marketplace.'
                        }}</span
                    >
                    <div class="market-controls">
                        <label class="utility-select"
                            ><SokoIcon name="globe" /><span class="sr-only"
                                >Market region</span
                            >
                            <select
                                v-model="selectedRegion"
                                @change="onRegionChange"
                            >
                                <option
                                    v-for="region in regions"
                                    :key="region.id"
                                    :value="region.id"
                                >
                                    {{ region.label }}
                                </option>
                            </select>
                        </label>
                        <label class="utility-select currency-select"
                            ><span class="sr-only">Display currency</span>
                            <select v-model="selectedCurrency">
                                <option
                                    v-for="currency in currencies"
                                    :key="currency.code"
                                    :value="currency.code"
                                >
                                    {{ currency.code }} {{ currency.symbol }}
                                </option>
                            </select>
                        </label>
                        <Link href="/dashboard" class="utility-supplier"
                            >Supplier workspace <SokoIcon name="arrow"
                        /></Link>
                    </div>
                </div>
            </div>
            <div class="site-container main-header">
                <Link href="/" class="brand-link"><SokoLogo /></Link>
                <form
                    class="header-search"
                    role="search"
                    novalidate
                    @submit.prevent="applySearch"
                >
                    <label class="sr-only" for="catalog-search"
                        >Search products and suppliers</label
                    >
                    <SokoIcon name="search" />
                    <input
                        id="catalog-search"
                        ref="searchInput"
                        v-model="searchQuery"
                        type="search"
                        placeholder="Search products, equipment and software"
                        autocomplete="off"
                    />
                    <button
                        v-if="searchQuery"
                        type="button"
                        class="search-clear icon-button"
                        aria-label="Clear search"
                        @click="clearSearch"
                    >
                        <SokoIcon name="close" />
                    </button>
                    <button
                        type="submit"
                        class="search-submit"
                        aria-label="Search Soko"
                    >
                        <SokoIcon name="arrow" />
                    </button>
                </form>
                <div class="header-actions">
                    <Link
                        href="/cart"
                        class="header-cart-link"
                        aria-label="View cart"
                        ><SokoIcon name="cart" /><span>Cart</span
                        ><span
                            v-if="Number(page.props.cartCount) > 0"
                            class="header-cart-count"
                            >{{ page.props.cartCount }}</span
                        ></Link
                    >
                    <template v-if="user"
                        ><Link
                            :href="
                                user.role === 'buyer'
                                    ? '/account'
                                    : '/dashboard'
                            "
                            class="account-link"
                            ><SokoIcon name="user" /><span>{{
                                user.name
                            }}</span></Link
                        ><Link
                            href="/logout"
                            method="post"
                            as="button"
                            class="icon-button"
                            aria-label="Sign out"
                            ><SokoIcon name="logout" /></Link
                    ></template>
                    <template v-else
                        ><Link href="/login" class="account-link">Sign in</Link
                        ><Link
                            href="/register"
                            class="button button-navy header-register"
                            >Join Soko <SokoIcon name="arrow" /></Link
                    ></template>
                </div>
            </div>
            <div class="navigation-border">
                <nav
                    class="site-container main-nav"
                    aria-label="Main navigation"
                >
                    <Link
                        v-for="item in navigation"
                        :key="item.label"
                        :href="catalogHref({ format: item.format })"
                        :class="{ active: isActive(item.format) }"
                        :aria-current="
                            isActive(item.format) ? 'page' : undefined
                        "
                        ><SokoIcon v-if="item.icon" :name="item.icon" />{{
                            item.label
                        }}</Link
                    >
                    <Link href="/#rfq" class="nav-rfq"
                        >Request a quotation <SokoIcon name="arrow"
                    /></Link>
                </nav>
            </div>
        </header>
        <div
            v-if="flash.success || flash.error"
            class="site-container flash-container"
            aria-live="polite"
        >
            <div
                v-if="flash.success"
                role="status"
                class="notice notice-success"
            >
                <SokoIcon name="check" /><span>{{ flash.success }}</span
                ><button
                    class="icon-button"
                    aria-label="Dismiss success message"
                    @click="dismissedSuccess = flash.success"
                >
                    <SokoIcon name="close" />
                </button>
            </div>
            <div v-if="flash.error" role="alert" class="notice notice-error">
                <span>{{ flash.error }}</span
                ><button
                    class="icon-button"
                    aria-label="Dismiss error message"
                    @click="dismissedError = flash.error"
                >
                    <SokoIcon name="close" />
                </button>
            </div>
        </div>
        <main id="main-content" class="site-main" tabindex="-1"><slot /></main>
        <footer class="site-footer">
            <div class="site-container">
                <div class="footer-grid">
                    <div class="footer-brand">
                        <SokoLogo />
                        <p>
                            Connecting businesses to the products, partners and
                            possibilities that move them forward.
                        </p>
                        <span class="footer-location"
                            ><SokoIcon name="globe" /> Built for trade. Open to
                            the world.</span
                        >
                    </div>
                    <div>
                        <h2>Explore Soko</h2>
                        <Link href="/">All products</Link
                        ><Link href="/?format=wholesale"
                            >Wholesale sourcing</Link
                        ><Link href="/?format=retail">Ready to buy</Link
                        ><Link href="/?format=digital">Digital products</Link
                        ><Link href="/?format=make_offer"
                            >Negotiable offers</Link
                        >
                    </div>
                    <div>
                        <h2>For your business</h2>
                        <Link href="/#rfq">Request a quotation</Link
                        ><Link href="/register">Create an account</Link
                        ><Link href="/dashboard">Supplier workspace</Link
                        ><Link href="/login">Sign in</Link>
                    </div>
                    <div>
                        <h2>Your marketplace</h2>
                        <button
                            v-for="region in regions.slice(1, 5)"
                            :key="region.id"
                            @click="switchRegion(region.id)"
                        >
                            {{ region.label }}
                        </button>
                    </div>
                </div>
                <div class="footer-bottom">
                    <span
                        >© {{ new Date().getFullYear() }} Soko. Global commerce,
                        connected.</span
                    ><span>Wholesale · Retail · Digital</span>
                </div>
            </div>
        </footer>
    </div>
</template>
<script setup lang="ts">
import { Head, Link, usePage } from '@inertiajs/vue3'
import { computed, nextTick, ref, watch } from 'vue'
import SokoLogo from './components/SokoLogo.vue'
import '../css/commerce.css'
import SokoIcon from './components/SokoIcon.vue'
import { useCurrency } from '@/composables/useCurrency'
import { useCatalogNavigation } from '@/composables/useCatalogNavigation'
const page = usePage()
const pageTitle = computed(() => {
    const titles: Record<string, string> = {
        'Store/Index': 'Global marketplace',
        'Auth/Login': 'Sign in',
        'Auth/Register': 'Create an account',
        'Dashboard/Index': 'Business workspace',
        'Store/Cart': 'Your cart',
        'Account/Index': 'Your account',
    }
    return page.component === 'Store/ProductDetail'
        ? String((page.props.product as any)?.title || 'Product details')
        : titles[page.component] || 'Global commerce'
})
const { selectedCurrency, selectedRegion, currencies, setRegion } =
    useCurrency()
const { visitCatalog, catalogHref } = useCatalogNavigation()
const user = computed(() => (page.props.auth as any)?.user)
const dismissedSuccess = ref('')
const dismissedError = ref('')
const flash = computed(() => {
    const value = (page.props.flash as Record<string, string>) || {}
    return {
        success: value.success === dismissedSuccess.value ? '' : value.success,
        error: value.error === dismissedError.value ? '' : value.error,
    }
})
watch(
    () => page.props.flash,
    () => {
        dismissedSuccess.value = ''
        dismissedError.value = ''
    }
)
const searchQuery = ref(String(page.props.searchQuery || ''))
const searchInput = ref<HTMLInputElement | null>(null)
watch(
    () => page.props.searchQuery,
    (query) => {
        searchQuery.value = String(query || '')
    }
)
watch(
    () => page.props.activeRegion,
    (region) => {
        if (region) selectedRegion.value = String(region)
    },
    { immediate: true }
)
const regions = [
    { id: 'global', label: 'Global market', currency: 'USD' },
    { id: 'north-america', label: 'North America', currency: 'USD' },
    { id: 'europe', label: 'Europe & UK', currency: 'EUR' },
    { id: 'east-africa', label: 'East Africa', currency: 'KES' },
    { id: 'west-africa', label: 'West Africa', currency: 'NGN' },
    { id: 'middle-east', label: 'Middle East', currency: 'AED' },
    { id: 'asia-pacific', label: 'Asia-Pacific', currency: 'CNY' },
]
const navigation = [
    { label: 'Marketplace', href: '/', format: 'all', icon: 'layers' },
    { label: 'Wholesale', href: '/?format=wholesale', format: 'wholesale' },
    { label: 'Ready to buy', href: '/?format=retail', format: 'retail' },
    { label: 'Digital products', href: '/?format=digital', format: 'digital' },
    {
        label: 'Make an offer',
        href: '/?format=make_offer',
        format: 'make_offer',
    },
]
const isActive = (format: string) =>
    page.component === 'Store/Index' &&
    String(page.props.activeFormat || 'all') === format
const applySearch = () => visitCatalog({ q: searchQuery.value.trim() })
const clearSearch = async () => {
    searchQuery.value = ''
    applySearch()
    await nextTick()
    searchInput.value?.focus()
}
const onRegionChange = () => {
    setRegion(
        selectedRegion.value,
        regions.find((region) => region.id === selectedRegion.value)
            ?.currency || 'USD'
    )
    visitCatalog({ region: selectedRegion.value })
}
const switchRegion = (region: string) => {
    selectedRegion.value = region
    onRegionChange()
}
</script>
