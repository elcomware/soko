<template>
    <div class="commerce-page detail-page">
        <nav class="commerce-breadcrumb" aria-label="Breadcrumb">
            <Link href="/">Marketplace</Link><span>/</span
            ><Link
                :href="
                    '/?category=' +
                    encodeURIComponent(product.category) +
                    '#catalog'
                "
                >{{ product.category }}</Link
            ><span>/</span><span>Product details</span>
        </nav>
        <div class="detail-grid">
            <div class="detail-media-column">
                <div class="detail-media">
                    <img
                        v-if="displayImage"
                        :src="displayImage"
                        :alt="product.title"
                        @error="onMediaError"
                    />
                    <div
                        v-else
                        class="detail-image-placeholder"
                        role="img"
                        aria-label="Product image not provided"
                    >
                        <SokoIcon name="package" />
                        <p>Product image not provided</p>
                    </div>
                    <span
                        v-if="isIllustrativeImage({ image_url: displayImage })"
                        class="detail-image-disclosure"
                        >Illustrative stock image</span
                    >
                    <span class="commerce-badge">{{
                        label(product.format)
                    }}</span>
                </div>
                <div
                    v-if="gallery.length > 1"
                    class="detail-gallery"
                    aria-label="Product images"
                >
                    <button
                        v-for="(image, i) in gallery"
                        :key="image"
                        :aria-label="'View image ' + (i + 1)"
                        :aria-pressed="
                            (selectedImage || product.image_url) === image
                        "
                        @click="selectedImage = image"
                    >
                        <img
                            :src="image"
                            alt=""
                            @error="failedImages.add(image)"
                        />
                    </button>
                </div>
                <section class="commerce-panel supplier-panel">
                    <div class="supplier-avatar">
                        <SokoIcon name="factory" />
                    </div>
                    <div>
                        <p class="commerce-eyebrow">Listed by</p>
                        <h2>{{ product.vendor_name }}</h2>
                        <p>
                            {{
                                product.vendor_country ||
                                'Country not specified'
                            }}
                        </p>
                    </div>
                    <span v-if="product.vendor_verified" class="commerce-badge"
                        ><SokoIcon name="check" />Verified account</span
                    >
                </section>
                <section class="commerce-panel product-description">
                    <p class="commerce-eyebrow">The details matter</p>
                    <h2>About this product</h2>
                    <p>
                        {{
                            product.description ||
                            'The supplier has not added a description. Include your requirements when making an inquiry.'
                        }}
                    </p>
                    <dl class="commerce-specs">
                        <div>
                            <dt>Category</dt>
                            <dd>{{ product.category }}</dd>
                        </div>
                        <div>
                            <dt>Condition</dt>
                            <dd>{{ label(product.condition) }}</dd>
                        </div>
                        <div>
                            <dt>SKU</dt>
                            <dd>{{ product.sku || 'Not provided' }}</dd>
                        </div>
                        <div v-if="product.lead_time_days > 0">
                            <dt>Supplier lead time</dt>
                            <dd>
                                {{ product.lead_time_days }} days · confirm with
                                supplier
                            </dd>
                        </div>
                        <div v-if="product.incoterm">
                            <dt>Trade terms</dt>
                            <dd>{{ product.incoterm }}</dd>
                        </div>
                        <div v-if="product.digital_file_size">
                            <dt>File size</dt>
                            <dd>{{ product.digital_file_size }}</dd>
                        </div>
                        <div v-if="product.digital_license_format">
                            <dt>License format</dt>
                            <dd>{{ product.digital_license_format }}</dd>
                        </div>
                    </dl>
                </section>
            </div>
            <div class="detail-buy commerce-panel">
                <p class="commerce-eyebrow">{{ product.category }}</p>
                <h1>{{ product.title }}</h1>
                <p class="detail-short">
                    {{
                        product.short_desc ||
                        'Review the listing, choose your quantity, and take the next step for your business.'
                    }}
                </p>
                <div class="detail-price">
                    {{
                        unitPrice === null
                            ? quoteError
                                ? 'Price unavailable'
                                : 'Updating price…'
                            : formatPrice(unitPrice)
                    }}<span>/ {{ digital ? 'license' : 'unit' }}</span>
                </div>
                <p
                    v-if="currentQuote?.discount_per_unit > 0"
                    class="commerce-small"
                >
                    <del>{{
                        formatPrice(currentQuote.original_unit_price)
                    }}</del>
                    · {{ currentQuote.promotion_name }} · Save
                    {{ formatPrice(currentQuote.discount_per_unit) }} per unit
                </p>
                <p v-if="quoteError" role="alert" class="form-error">
                    {{ quoteError }}
                    <button type="button" @click="refreshQuote">
                        Retry price
                    </button>
                </p>
                <p class="commerce-muted">
                    {{
                        activeCurrency.code === 'USD'
                            ? 'Listed in USD'
                            : 'Indicative ' +
                              activeCurrency.code +
                              ' conversion · orders recorded in USD'
                    }}
                </p>
                <p v-if="product.price_tiers?.length" class="commerce-small">
                    Base volume prices below. Current promotions are applied to
                    your selected quantity above.
                </p>
                <div
                    v-if="product.price_tiers?.length"
                    class="detail-tiers"
                    aria-label="Base volume pricing before any current promotion"
                >
                    <div
                        v-for="tier in product.price_tiers"
                        :key="tier.min_qty"
                        :class="{
                            selected:
                                quantity >= tier.min_qty &&
                                (!tier.max_qty || quantity <= tier.max_qty),
                        }"
                    >
                        <span
                            >{{ tier.min_qty
                            }}{{
                                tier.max_qty ? '–' + tier.max_qty : '+'
                            }}
                            units</span
                        ><strong>{{ formatPrice(tier.price) }}</strong>
                    </div>
                </div>
                <div class="detail-availability">
                    <span :class="{ 'out-of-stock': product.stock === 0 }">{{
                        product.stock === 0
                            ? 'Currently unavailable'
                            : product.stock < 0
                              ? 'Digital listing available'
                              : product.stock + ' units available'
                    }}</span
                    ><span>Min. {{ product.moq || 1 }} units</span>
                </div>
                <form
                    class="soko-form"
                    novalidate
                    :aria-busy="busy"
                    @submit.prevent="addToCart"
                >
                    <div class="soko-field">
                        <label for="detail-quantity">Quantity</label
                        ><input
                            id="detail-quantity"
                            v-model.number="quantity"
                            v-bind="fieldState('detail-quantity')"
                            type="number"
                            :min="product.moq || 1"
                            :max="product.stock >= 0 ? product.stock : 1000000"
                            step="1"
                            required
                        />
                    </div>
                    <p
                        v-if="error && !showOffer"
                        id="market-form-error"
                        class="form-error"
                        role="alert"
                    >
                        {{ error }}
                    </p>
                    <div class="detail-total">
                        <span>Estimated total</span
                        ><strong>{{
                            unitPrice === null
                                ? '—'
                                : formatPrice(
                                      unitPrice * (Number(quantity) || 0)
                                  )
                        }}</strong>
                    </div>
                    <button
                        class="button button-amber button-full"
                        :disabled="
                            busy || product.stock === 0 || unitPrice === null
                        "
                    >
                        {{ busy ? 'Adding…' : 'Add to cart'
                        }}<SokoIcon name="arrow" /></button
                    ><button
                        v-if="product.accepts_offers"
                        type="button"
                        class="button button-subtle button-full"
                        :disabled="
                            busy || product.stock === 0 || unitPrice === null
                        "
                        @click="openOffer"
                    >
                        Make an offer
                    </button>
                </form>
                <p class="commerce-small">
                    <SokoIcon name="document" />Adding to cart does not charge
                    you. Orders are submitted for review; payment and delivery
                    are arranged separately.
                </p>
                <Link href="/#rfq" class="commerce-text-link"
                    >Need something different? Create a sourcing request
                    <SokoIcon name="arrow"
                /></Link>
            </div>
        </div>
        <section v-if="related?.length" class="related-section">
            <div class="commerce-section-heading">
                <div>
                    <p class="commerce-eyebrow">Keep exploring</p>
                    <h2>More in this category</h2>
                </div>
                <Link
                    :href="
                        '/?category=' +
                        encodeURIComponent(product.category) +
                        '#catalog'
                    "
                    >View category <SokoIcon name="arrow"
                /></Link>
            </div>
            <div class="related-grid">
                <Link
                    v-for="item in related"
                    :key="item.id"
                    :href="'/products/' + item.id"
                    class="commerce-panel related-card"
                    ><img
                        v-if="productImage(item)"
                        :src="productImage(item)"
                        :alt="item.title"
                        loading="lazy"
                    />
                    <div>
                        <span class="commerce-eyebrow">{{
                            label(item.format)
                        }}</span>
                        <h3>{{ item.title }}</h3>
                        <p
                            v-if="isIllustrativeImage(item)"
                            class="commerce-small"
                        >
                            Illustrative stock image
                        </p>
                        <span>View product and pricing</span>
                    </div></Link
                >
            </div>
        </section>
        <SokoDialog
            :open="showOffer"
            title="Start a conversation"
            @close="closeOffer"
            ><form
                class="soko-form"
                novalidate
                :aria-busy="busy"
                @submit.prevent="submitOffer"
            >
                <p class="form-note">
                    Offer a unit price in USD for {{ product.title }}.
                    Acceptance does not collect payment.
                </p>
                <p
                    v-if="error"
                    id="market-form-error"
                    class="form-error"
                    role="alert"
                >
                    {{ error }}
                </p>
                <div class="form-row">
                    <div class="soko-field">
                        <label for="detail-offer-price">Unit price (USD)</label
                        ><input
                            id="detail-offer-price"
                            v-model.number="offeredPrice"
                            v-bind="fieldState('detail-offer-price')"
                            type="number"
                            min="0.01"
                            step="0.01"
                            required
                        />
                    </div>
                    <div class="soko-field">
                        <label for="detail-offer-qty">Quantity</label
                        ><input
                            id="detail-offer-qty"
                            v-model.number="offerQuantity"
                            v-bind="fieldState('detail-offer-qty')"
                            type="number"
                            :min="product.moq || 1"
                            :max="product.stock >= 0 ? product.stock : 1000000"
                            step="1"
                            required
                        />
                    </div>
                </div>
                <div class="soko-field">
                    <label for="detail-offer-note">Message (optional)</label
                    ><textarea
                        class="resize-none"
                        id="detail-offer-note"
                        v-model="message"
                        rows="4"
                        maxlength="2000"
                        placeholder="Share timing, delivery requirements or questions."
                    ></textarea>
                </div>
                <p v-if="!user" class="form-note">
                    Sign in is required to submit and track an offer.
                </p>
                <button class="button button-navy button-full" :disabled="busy">
                    {{
                        busy
                            ? 'Sending…'
                            : user
                              ? 'Send offer'
                              : 'Continue to sign in'
                    }}
                </button>
            </form></SokoDialog
        >
    </div>
</template>
<script setup lang="ts">
import { computed, ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { Link, usePage, router } from '@inertiajs/vue3'
import SokoIcon from '@/components/SokoIcon.vue'
import SokoDialog from '@/components/SokoDialog.vue'
import { useCurrency, type Currency } from '@/composables/useCurrency'
import { useMarketplaceForm } from '@/composables/useMarketplaceForm'
import '@/../css/commerce.css'
import { productImage, isIllustrativeImage } from '@/design/merchandising'
const props = defineProps<{
    product: any
    related: any[]
    currencies?: Currency[]
}>()
const { activeCurrency, formatPrice } = useCurrency(props.currencies)
const { busy, error, validate, post, fieldState } = useMarketplaceForm()
const user = computed(() => (usePage().props.auth as any)?.user)
const quantity = ref(props.product.moq || 1),
    selectedImage = ref(''),
    showOffer = ref(false),
    offeredPrice = ref(props.product.price),
    offerQuantity = ref(props.product.moq || 1),
    message = ref('')
const failedImages = ref(new Set<string>())
const label = (s: string) => String(s || '').replace(/_/g, ' ')
const digital = computed(() => props.product.type !== 'physical')
const gallery = computed(
    () =>
        Array.from(
            new Set([
                props.product.image_url,
                ...(props.product.gallery_urls || []),
            ])
        ).filter(
            (image) =>
                productImage({ image_url: image }) &&
                !failedImages.value.has(image)
        ) as string[]
)
const displayImage = computed(() =>
    gallery.value.includes(selectedImage.value)
        ? selectedImage.value
        : gallery.value[0] || ''
)
function onMediaError(event: Event) {
    const src = (event.currentTarget as HTMLImageElement).getAttribute('src')
    if (src) failedImages.value.add(src)
}
// Quotes are quantity-specific: the best discount can differ between volume tiers.
const quote = ref(props.product.pricing || null)
const quoteError = ref('')
const clock = ref(Date.now())
let quoteSequence = 0
let quoteTimer: ReturnType<typeof setTimeout> | undefined
let expiryTimer: ReturnType<typeof setInterval> | undefined
let quoteAbort: AbortController | undefined
const currentQuote = computed(() => {
    const value = quote.value
    if (!value || value.quantity !== Number(quantity.value)) return null
    if (
        value.promotion_ends_at &&
        Date.parse(value.promotion_ends_at) <= clock.value
    )
        return null
    return value
})
const unitPrice = computed(() => currentQuote.value?.unit_price ?? null)
async function refreshQuote() {
    const sequence = ++quoteSequence
    quoteAbort?.abort()
    quoteAbort = new AbortController()
    quoteError.value = ''
    const qty = Number(quantity.value)
    if (
        !Number.isInteger(qty) ||
        qty < (props.product.moq || 1) ||
        qty > 1000000 ||
        (props.product.stock >= 0 && qty > props.product.stock)
    ) {
        quote.value = null
        quoteError.value =
            'Enter an available whole quantity to see your price.'
        return
    }
    try {
        const response = await fetch(
            '/products/' +
                encodeURIComponent(props.product.id) +
                '/price?quantity=' +
                qty,
            {
                signal: quoteAbort.signal,
                headers: { Accept: 'application/json' },
                cache: 'no-store',
            }
        )
        const body = await response.json()
        if (sequence !== quoteSequence) return
        if (!response.ok)
            throw new Error(body.error || 'Unable to confirm this price.')
        quote.value = body.pricing
        clock.value = Date.now()
    } catch (failure) {
        if (
            sequence !== quoteSequence ||
            (failure instanceof Error && failure.name === 'AbortError')
        )
            return
        quote.value = null
        quoteError.value =
            failure instanceof Error
                ? failure.message
                : 'Unable to confirm this price. Please retry.'
    }
}
watch([quantity, () => props.product.id], () => {
    ++quoteSequence
    quoteAbort?.abort()
    quote.value = null
    quoteError.value = ''
    if (quoteTimer) clearTimeout(quoteTimer)
    quoteTimer = setTimeout(refreshQuote, 300)
})
onMounted(() => {
    void refreshQuote()
    let ticks = 0
    expiryTimer = setInterval(() => {
        clock.value = Date.now()
        const expired =
            quote.value?.promotion_ends_at &&
            Date.parse(quote.value.promotion_ends_at) <= clock.value
        if (expired) quote.value = null
        if (expired || ++ticks % 30 === 0) void refreshQuote()
    }, 1000)
})
onBeforeUnmount(() => {
    ++quoteSequence
    quoteAbort?.abort()
    if (quoteTimer) clearTimeout(quoteTimer)
    if (expiryTimer) clearInterval(expiryTimer)
})

watch(
    () => props.product.id,
    () => {
        quantity.value = props.product.moq || 1
        selectedImage.value = ''
        error.value = ''
        showOffer.value = false
    }
)
function addToCart(e: Event) {
    if (unitPrice.value !== null && validate(e))
        post('/cart/items', {
            product_id: props.product.id,
            quantity: quantity.value,
        })
}
function closeOffer() {
    showOffer.value = false
    error.value = ''
}
function openOffer() {
    error.value = ''
    offeredPrice.value = props.product.price
    offerQuantity.value = quantity.value
    showOffer.value = true
}
function submitOffer(e: Event) {
    if (!validate(e)) return
    if (!user.value) {
        router.visit('/login')
        return
    }
    post(
        '/make-offer',
        {
            product_id: props.product.id,
            offered_price: offeredPrice.value,
            quantity: offerQuantity.value,
            message: message.value,
        },
        () => {
            showOffer.value = false
        }
    )
}
</script>

<style scoped>
.detail-image-placeholder {
    display: grid;
    place-content: center;
    justify-items: center;
    width: 100%;
    min-height: 360px;
    color: var(--soko-muted);
    background: var(--soko-paper, #f5f4ef);
}
.detail-image-placeholder svg {
    width: 96px;
    height: 96px;
}
.detail-image-disclosure {
    position: absolute;
    bottom: 1rem;
    left: 1rem;
    padding: 0.4rem 0.6rem;
    border-radius: 0.4rem;
    background: var(--soko-paper, #f5f4ef);
    color: var(--soko-ink);
    font-size: 0.75rem;
}
</style>
