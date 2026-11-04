<template>
    <div class="commerce-page cart-page">
        <nav class="commerce-breadcrumb" aria-label="Breadcrumb">
            <Link href="/">Marketplace</Link><span>/</span
            ><span>Your cart</span>
        </nav>
        <div class="commerce-page-heading">
            <div>
                <p class="commerce-eyebrow">A good next step</p>
                <h1>Your business, in the making.</h1>
                <p>Review your selection before submitting an order.</p>
            </div>
            <Link href="/#catalog" class="commerce-text-link"
                >Keep exploring <SokoIcon name="arrow"
            /></Link>
        </div>
        <div v-if="!cart.items?.length" class="commerce-panel commerce-empty">
            <SokoIcon name="box" />
            <h2>Your cart is ready for possibility.</h2>
            <p>
                Explore equipment, products and digital tools for your next
                project.
            </p>
            <Link href="/#catalog" class="button button-navy"
                >Explore products <SokoIcon name="arrow"
            /></Link>
        </div>
        <div v-else class="cart-layout">
            <section aria-label="Cart items">
                <p
                    v-if="error"
                    id="market-form-error"
                    class="form-error"
                    role="alert"
                >
                    {{ error }}
                </p>
                <div
                    v-if="cart.issues?.length"
                    class="commerce-panel cart-issues"
                    role="alert"
                >
                    <h2>Your cart needs attention</h2>
                    <p v-for="(issue, i) in cart.issues" :key="i">
                        {{ typeof issue === 'string' ? issue : issue.message }}
                    </p>
                </div>
                <article
                    v-for="item in cart.items"
                    :key="item.product_id"
                    class="commerce-panel cart-item"
                >
                    <Link :href="'/products/' + item.product_id"
                        ><img
                            :src="
                                item.product.image_url ||
                                '/public/images/soko-industry.jpg'
                            "
                            :alt="item.product.title"
                    /></Link>
                    <div class="cart-item-main">
                        <p class="commerce-eyebrow">
                            {{ item.product.category }}
                        </p>
                        <h2>
                            <Link :href="'/products/' + item.product_id">{{
                                item.product.title
                            }}</Link>
                        </h2>
                        <p>{{ item.product.vendor_name }}</p>
                        <span class="commerce-muted"
                            >{{ formatPrice(item.unit_price) }} per unit · Min.
                            {{ item.product.moq || 1 }}</span
                        >
                        <form
                            class="cart-quantity-form"
                            novalidate
                            @submit.prevent="update($event, item)"
                        >
                            <label :for="'cart-qty-' + item.product_id"
                                >Quantity</label
                            ><input
                                :id="'cart-qty-' + item.product_id"
                                v-model.number="quantities[item.product_id]"
                                v-bind="
                                    fieldState('cart-qty-' + item.product_id)
                                "
                                type="number"
                                :min="item.product.moq || 1"
                                :max="
                                    item.product.stock >= 0
                                        ? item.product.stock
                                        : 1000000
                                "
                                step="1"
                                required
                            /><button
                                class="button button-subtle"
                                :disabled="busy"
                            >
                                Update</button
                            ><button
                                type="button"
                                class="commerce-remove"
                                :disabled="busy"
                                @click="
                                    post(
                                        '/cart/items/' +
                                            item.product_id +
                                            '/delete',
                                        {}
                                    )
                                "
                                :aria-label="'Remove ' + item.product.title"
                            >
                                Remove
                            </button>
                        </form>
                    </div>
                    <strong class="cart-item-total">{{
                        formatPrice(item.total_price)
                    }}</strong>
                </article>
            </section>
            <aside class="commerce-panel cart-summary">
                <p class="commerce-eyebrow">Ready when you are</p>
                <h2>Order summary</h2>
                <dl class="commerce-specs">
                    <div>
                        <dt>Items</dt>
                        <dd>{{ cart.items.length }}</dd>
                    </div>
                    <div>
                        <dt>Subtotal</dt>
                        <dd>{{ formatPrice(cart.subtotal) }}</dd>
                    </div>
                    <div>
                        <dt>Shipping & taxes</dt>
                        <dd>To be confirmed</dd>
                    </div>
                    <div class="cart-summary-total">
                        <dt>Product total</dt>
                        <dd>{{ formatPrice(cart.subtotal) }}</dd>
                    </div>
                </dl>
                <p v-if="activeCurrency.code !== 'USD'" class="commerce-muted">
                    Display conversion is indicative. Order amount:
                    {{ usd(cart.subtotal) }} USD.
                </p>
                <div class="commerce-callout">
                    <SokoIcon name="document" />
                    <p>
                        Orders are recorded as pending. Payment, shipping and
                        digital fulfillment are arranged separately. No payment
                        is collected here.
                    </p>
                </div>
                <form
                    v-if="user"
                    class="soko-form"
                    novalidate
                    @submit.prevent="checkout"
                >
                    <button
                        class="button button-amber button-full"
                        :disabled="busy || !!cart.issues?.length"
                    >
                        {{ busy ? 'Submitting…' : 'Submit order for review'
                        }}<SokoIcon name="arrow" />
                    </button>
                    <p class="commerce-small">
                        Submitted as {{ user.name }} · {{ user.email }}
                    </p>
                </form>
                <Link
                    v-else
                    href="/login"
                    class="button button-navy button-full"
                    >Sign in to submit order <SokoIcon name="arrow" /></Link
                ><Link href="/account" class="commerce-text-link"
                    >View your order history</Link
                >
            </aside>
        </div>
    </div>
</template>
<script setup lang="ts">
import { reactive, watch } from 'vue'
import { Link, usePage } from '@inertiajs/vue3'
import { computed } from 'vue'
import SokoIcon from '@/components/SokoIcon.vue'
import { useCurrency, type Currency } from '@/composables/useCurrency'
import { useMarketplaceForm } from '@/composables/useMarketplaceForm'
import '@/../css/commerce.css'
const props = defineProps<{ cart: any; currencies?: Currency[] }>()
const { formatPrice, activeCurrency } = useCurrency(props.currencies),
    { busy, error, validate, post, fieldState } = useMarketplaceForm()
const user = computed(() => (usePage().props.auth as any)?.user)
const quantities = reactive<Record<string, number>>({})
watch(
    () => props.cart,
    (cart) => {
        for (const item of cart.items || [])
            quantities[item.product_id] = item.quantity
    },
    { immediate: true }
)
const usd = (n: number) =>
    Number(n).toLocaleString('en-US', { style: 'currency', currency: 'USD' })
function update(e: Event, item: any) {
    if (validate(e))
        post('/cart/items/' + item.product_id, {
            quantity: quantities[item.product_id],
        })
}
function checkout() {
    post('/checkout', {})
}
</script>
