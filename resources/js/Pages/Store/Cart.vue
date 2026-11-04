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
                        <span
                            v-if="item.discount_per_unit > 0"
                            class="commerce-muted"
                            ><del>{{
                                formatPrice(item.original_unit_price)
                            }}</del>
                            · {{ item.pricing?.promotion_name }}</span
                        >
                        <span class="commerce-muted"
                            >{{ formatPrice(item.unit_price) }} per unit · Min.
                            {{ item.product.moq || 1 }}</span
                        >
                        <form
                            class="cart-quantity-form"
                            novalidate
                            @submit.prevent="flushItem(item)"
                        >
                            <label :for="'cart-qty-' + item.product_id"
                                >Quantity</label
                            ><input
                                :id="'cart-qty-' + item.product_id"
                                v-model="quantities[item.product_id]"
                                :disabled="busy"
                                @input="schedule(item)"
                                @blur="flushItem(item)"
                                @keydown.enter.prevent="flushItem(item)"
                                :aria-invalid="
                                    !!lineErrors[item.product_id] || undefined
                                "
                                :aria-describedby="
                                    'cart-status-' + item.product_id
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
                                type="button"
                                class="commerce-remove"
                                :disabled="busy"
                                @click="removeItem(item)"
                                :aria-label="'Remove ' + item.product.title"
                            >
                                Remove
                            </button>
                        </form>
                        <p
                            :id="'cart-status-' + item.product_id"
                            role="status"
                            aria-live="polite"
                            class="commerce-small"
                            :class="{
                                'form-error': lineErrors[item.product_id],
                            }"
                        >
                            {{
                                lineErrors[item.product_id] ||
                                (saving === item.product_id
                                    ? 'Saving quantity…'
                                    : dirty[item.product_id]
                                      ? 'Waiting to save…'
                                      : 'Quantity saved automatically')
                            }}
                        </p>
                        <button
                            v-if="lineErrors[item.product_id]"
                            type="button"
                            class="commerce-text-link"
                            @click="flushItem(item)"
                            :disabled="busy"
                        >
                            {{
                                failedRemovals.has(item.product_id)
                                    ? 'Retry removing item'
                                    : 'Retry saving quantity'
                            }}
                        </button>
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
                        :disabled="busy || !!cart.issues?.length || hasErrors"
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
import { reactive, watch, computed, ref, onBeforeUnmount } from 'vue'
import { Link, usePage, router } from '@inertiajs/vue3'
import SokoIcon from '@/components/SokoIcon.vue'
import { useCurrency, type Currency } from '@/composables/useCurrency'
import '@/../css/commerce.css'
const props = defineProps<{ cart: any; currencies?: Currency[] }>()
const { formatPrice, activeCurrency } = useCurrency(props.currencies)
const user = computed(() => (usePage().props.auth as any)?.user)
const quantities = reactive<Record<string, string | number>>({})
const dirty = reactive<Record<string, boolean>>({})
const lineErrors = reactive<Record<string, string>>({})
const removals = new Set<string>()
const saving = ref(''),
    busy = ref(false),
    error = ref('')
let timer: ReturnType<typeof setTimeout> | undefined
let running: Promise<void> | undefined
let disposed = false
const hasErrors = computed(() => Object.values(lineErrors).some(Boolean))
watch(
    () => props.cart,
    (cart) => {
        for (const item of cart.items || []) {
            if (!dirty[item.product_id] && saving.value !== item.product_id)
                quantities[item.product_id] = item.quantity
        }
    },
    { immediate: true }
)
const usd = (n: number) =>
    Number(n).toLocaleString('en-US', { style: 'currency', currency: 'USD' })
function validQuantity(item: any) {
    const raw = String(quantities[item.product_id] ?? '').trim()
    const value = Number(raw),
        minimum = item.product.moq || 1
    const maximum = item.product.stock >= 0 ? item.product.stock : 1000000
    lineErrors[item.product_id] =
        raw === '' ||
        !Number.isInteger(value) ||
        value < minimum ||
        value > maximum
            ? `Enter a whole number between ${minimum} and ${maximum}.`
            : ''
    return !lineErrors[item.product_id]
}
const failedRemovals = ref(new Set<string>())
function schedule(item: any) {
    if (busy.value) return
    failedRemovals.value.delete(item.product_id)
    dirty[item.product_id] = true
    validQuantity(item)
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
        void drain()
    }, 450)
}
function flushItem(item: any) {
    if (busy.value) return
    if (timer) clearTimeout(timer)
    if (failedRemovals.value.has(item.product_id)) {
        failedRemovals.value.delete(item.product_id)
        removeItem(item)
        return
    }
    if (validQuantity(item)) void drain()
}
function request(url: string, data: Record<string, any>): Promise<string> {
    return new Promise((resolve) => {
        let message = ''
        router.post(url, data, {
            preserveScroll: true,
            preserveState: true,
            onSuccess: (page) => {
                message = String((page.props.flash as any)?.error || '')
            },
            onError: (errors) => {
                message =
                    Object.values(errors).join(' ') ||
                    'Unable to save. Try again.'
            },
            onHttpException: () => {
                message =
                    'The server could not save this change. Please try again.'
                return false
            },
            onNetworkError: () => {
                message =
                    url === '/checkout'
                        ? 'Order submission is unconfirmed. Check your orders first. Retrying this submission uses the same reference to avoid a duplicate.'
                        : 'Connection lost. Your change has not been confirmed. Try again.'
                return false
            },
            onCancel: () => {
                message = 'Update interrupted. Please try again.'
            },
            onFinish: () => resolve(message),
        })
    })
}
function drain(): Promise<void> {
    if (running) return running
    running = (async () => {
        while (!disposed) {
            const removeID = removals.values().next().value as
                string | undefined
            const item = props.cart.items?.find(
                (row: any) =>
                    dirty[row.product_id] &&
                    !lineErrors[row.product_id] &&
                    !removals.has(row.product_id)
            )
            if (!removeID && !item) break
            const id = removeID || item.product_id
            const sent = removeID ? 0 : Number(quantities[id])
            saving.value = id
            let message = ''
            try {
                message = await request(
                    '/cart/items/' +
                        encodeURIComponent(id) +
                        (removeID ? '/delete' : ''),
                    removeID ? {} : { quantity: sent }
                )
            } catch {
                message =
                    'Could not save this change. Check your connection and try again.'
            }
            saving.value = ''
            if (removeID) {
                removals.delete(id)
                if (message) failedRemovals.value.add(id)
                if (!message) {
                    failedRemovals.value.delete(id)
                    delete dirty[id]
                    delete quantities[id]
                    delete lineErrors[id]
                }
            } else if (Number(quantities[id]) === sent) {
                if (!message) dirty[id] = false
            }
            if (message) {
                lineErrors[id] = message
                error.value = message
            }
        }
    })().finally(() => {
        running = undefined
    })
    return running
}
function removeItem(item: any) {
    if (busy.value) return
    const id = item.product_id
    removals.add(id)
    delete lineErrors[id]
    dirty[id] = false
    void drain()
}
async function checkout() {
    if (busy.value) return
    busy.value = true
    error.value = ''
    if (timer) clearTimeout(timer)
    for (const item of props.cart.items || []) validQuantity(item)
    await drain()
    if (
        hasErrors.value ||
        failedRemovals.value.size > 0 ||
        Object.values(dirty).some(Boolean)
    ) {
        error.value = 'Check the quantities above before submitting your order.'
        busy.value = false
        return
    }
    try {
        error.value = await request('/checkout', {
            checkout_token: props.cart.checkout_token,
        })
    } catch {
        error.value =
            'Order submission could not be confirmed. Check your orders before trying again.'
    } finally {
        busy.value = false
    }
}
onBeforeUnmount(() => {
    disposed = true
    if (timer) clearTimeout(timer)
})
</script>
