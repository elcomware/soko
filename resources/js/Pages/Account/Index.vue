<template>
    <div class="commerce-page account-page">
        <div class="commerce-page-heading">
            <div>
                <p class="commerce-eyebrow">Your Soko</p>
                <h1>Good to see you, {{ currentUser.name.split(' ')[0] }}.</h1>
                <p>
                    Your orders, conversations and sourcing requests in one
                    place.
                </p>
            </div>
            <Link href="/#catalog" class="button button-navy"
                >Explore marketplace <SokoIcon name="arrow"
            /></Link>
        </div>
        <nav class="account-tabs" aria-label="Account sections">
            <Link
                v-for="tab in tabs"
                :key="tab.id"
                :href="'/account?section=' + tab.id"
                :aria-current="section === tab.id ? 'page' : undefined"
                >{{ tab.label
                }}<span v-if="tab.count !== undefined">{{
                    tab.count
                }}</span></Link
            >
        </nav>
        <p v-if="error && !confirmation" class="form-error" role="alert">
            {{ error }}
        </p>
        <template v-if="section === 'overview'"
            ><div class="account-overview">
                <section class="account-welcome">
                    <SokoIcon name="globe" />
                    <p class="commerce-eyebrow">
                        From discovery to your next move
                    </p>
                    <h2>Make room for what’s next.</h2>
                    <p>
                        Keep an eye on pending orders, review supplier
                        quotations and continue your conversations.
                    </p>
                    <Link href="/#rfq" class="button button-amber"
                        >Create a sourcing request <SokoIcon name="arrow"
                    /></Link>
                </section>
                <section class="commerce-panel">
                    <p class="commerce-eyebrow">Account details</p>
                    <h2>{{ currentUser.name }}</h2>
                    <dl class="commerce-specs">
                        <div>
                            <dt>Email</dt>
                            <dd>{{ currentUser.email }}</dd>
                        </div>
                        <div>
                            <dt>Company</dt>
                            <dd>
                                {{ currentUser.company_name || 'Not provided' }}
                            </dd>
                        </div>
                        <div>
                            <dt>Country</dt>
                            <dd>{{ currentUser.country || 'Not provided' }}</dd>
                        </div>
                        <div>
                            <dt>Account</dt>
                            <dd>{{ currentUser.role }}</dd>
                        </div>
                    </dl>
                    <Link
                        v-if="currentUser.role !== 'buyer'"
                        href="/dashboard"
                        class="commerce-text-link"
                        >Open business workspace <SokoIcon name="arrow"
                    /></Link>
                </section>
            </div>
            <div class="account-metrics">
                <Link href="/account?section=orders" class="commerce-panel"
                    ><span>Orders</span
                    ><strong>{{ orders?.length || 0 }}</strong
                    ><small>{{ pendingOrders }} awaiting review</small></Link
                ><Link href="/account?section=offers" class="commerce-panel"
                    ><span>Offers</span
                    ><strong>{{ offers?.length || 0 }}</strong
                    ><small
                        >{{
                            (offers || []).filter((o) => o.status === 'pending')
                                .length
                        }}
                        pending responses</small
                    ></Link
                ><Link href="/account?section=requests" class="commerce-panel"
                    ><span>Sourcing requests</span
                    ><strong>{{ rfqs?.length || 0 }}</strong
                    ><small
                        >{{
                            (rfqs || []).filter((r) => r.status === 'open')
                                .length
                        }}
                        open requests</small
                    ></Link
                >
            </div></template
        >
        <section v-else-if="section === 'orders'" class="account-records">
            <div class="commerce-section-heading">
                <div>
                    <p class="commerce-eyebrow">Your purchases</p>
                    <h2>Orders</h2>
                </div>
            </div>
            <div v-if="!orders?.length" class="commerce-panel commerce-empty">
                <SokoIcon name="box" />
                <h2>No orders yet</h2>
                <p>Items you submit for review will appear here.</p>
                <Link href="/#catalog" class="button button-navy"
                    >Find your next product</Link
                >
            </div>
            <article
                v-for="order in orders"
                :key="order.id"
                class="commerce-panel order-record"
            >
                <div class="record-heading">
                    <div>
                        <p class="commerce-eyebrow">
                            {{ date(order.created_at) }}
                        </p>
                        <h3>{{ order.order_number }}</h3>
                    </div>
                    <span
                        class="commerce-status"
                        :class="'status-' + order.status"
                        >{{ order.status }}</span
                    ><strong>{{ formatPrice(order.total_amount) }}</strong>
                </div>
                <div
                    v-for="item in order.items || []"
                    :key="item.id"
                    class="order-line"
                >
                    <div>
                        <h4>{{ item.product_title }}</h4>
                        <p>
                            {{ item.quantity }} ×
                            {{ formatPrice(item.unit_price) }}
                        </p>
                    </div>
                    <strong>{{ formatPrice(item.total_price) }}</strong>
                </div>
                <p class="commerce-small">
                    Recorded in {{ order.currency }} ·
                    {{
                        order.status === 'pending'
                            ? 'Awaiting review. No payment has been collected through Soko.'
                            : 'Status recorded by the marketplace.'
                    }}
                </p>
                <button
                    v-if="order.status === 'pending'"
                    class="commerce-remove"
                    :disabled="busy"
                    @click="openConfirmation({ kind: 'cancel', order })"
                >
                    Cancel pending order
                </button>
            </article>
        </section>
        <section v-else-if="section === 'offers'" class="account-records">
            <div class="commerce-section-heading">
                <div>
                    <p class="commerce-eyebrow">Your conversations</p>
                    <h2>Offers</h2>
                </div>
            </div>
            <div v-if="!offers?.length" class="commerce-panel commerce-empty">
                <SokoIcon name="document" />
                <h2>No offers yet</h2>
                <p>
                    Make an offer on an eligible listing to start a
                    conversation.
                </p>
                <Link
                    href="/?format=make_offer#catalog"
                    class="button button-navy"
                    >Explore negotiable listings</Link
                >
            </div>
            <article
                v-for="offer in offers"
                :key="offer.id"
                class="commerce-panel offer-record"
            >
                <div class="record-heading">
                    <div>
                        <p class="commerce-eyebrow">
                            {{ date(offer.created_at) }}
                        </p>
                        <h3>
                            <Link :href="'/products/' + offer.product_id">{{
                                offer.product_title
                            }}</Link>
                        </h3>
                    </div>
                    <span
                        class="commerce-status"
                        :class="'status-' + offer.status"
                        >{{ offer.status }}</span
                    >
                </div>
                <div class="offer-figures">
                    <span>{{ offer.quantity }} units</span
                    ><strong
                        >{{ formatPrice(offer.offered_price) }} per unit</strong
                    ><span>{{ formatPrice(offer.total_offered) }} total</span>
                </div>
                <p v-if="offer.message" class="record-note">
                    {{ offer.message }}
                </p>
                <p class="commerce-small">
                    {{
                        offer.status === 'accepted'
                            ? 'The supplier accepted your proposal. Arrange payment and fulfillment with the supplier.'
                            : 'Your proposal is recorded. Acceptance does not collect payment.'
                    }}
                </p>
            </article>
        </section>
        <section v-else class="account-records">
            <div class="commerce-section-heading">
                <div>
                    <p class="commerce-eyebrow">
                        Requirements to possibilities
                    </p>
                    <h2>Sourcing requests</h2>
                </div>
                <Link href="/#rfq" class="button button-subtle"
                    >New request <SokoIcon name="arrow"
                /></Link>
            </div>
            <div v-if="!rfqs?.length" class="commerce-panel commerce-empty">
                <SokoIcon name="document" />
                <h2>Start with what you need.</h2>
                <p>Publish a clear brief and review quotations here.</p>
                <Link href="/#rfq" class="button button-navy"
                    >Create a sourcing request</Link
                >
            </div>
            <article
                v-for="rfq in rfqs"
                :key="rfq.id"
                class="commerce-panel request-record"
            >
                <div class="record-heading">
                    <div>
                        <p class="commerce-eyebrow">
                            {{ rfq.rfq_number }} · {{ rfq.category }}
                        </p>
                        <h3>{{ rfq.title }}</h3>
                    </div>
                    <span
                        class="commerce-status"
                        :class="'status-' + rfq.status"
                        >{{ rfq.status }}</span
                    >
                </div>
                <p>{{ rfq.quantity }} units · {{ rfq.destination_country }}</p>
                <p v-if="rfq.specifications" class="record-note">
                    {{ rfq.specifications }}
                </p>
                <details class="request-quotes">
                    <summary>
                        {{ rfq.quotes?.length || 0 }} supplier quotations
                    </summary>
                    <p v-if="!rfq.quotes?.length" class="commerce-muted">
                        Quotations will appear here when suppliers respond.
                    </p>
                    <div
                        v-for="quote in rfq.quotes || []"
                        :key="quote.id"
                        class="quote-record"
                    >
                        <div>
                            <h4>{{ quote.vendor_name }}</h4>
                            <p>
                                {{ quote.vendor_country }} ·
                                {{ quote.lead_time_days }} day lead time
                            </p>
                            <p v-if="quote.notes">{{ quote.notes }}</p>
                        </div>
                        <div>
                            <strong
                                >{{ formatPrice(quote.unit_price) }} /
                                unit</strong
                            ><span class="commerce-status">{{
                                quote.status
                            }}</span
                            ><button
                                v-if="
                                    rfq.status === 'open' &&
                                    quote.status === 'pending'
                                "
                                class="button button-navy"
                                :disabled="busy"
                                @click="
                                    confirmation = { kind: 'award', rfq, quote }
                                "
                            >
                                Select quotation
                            </button>
                        </div>
                    </div>
                </details>
                <button
                    v-if="rfq.status === 'open'"
                    class="commerce-remove"
                    :disabled="busy"
                    @click="confirmation = { kind: 'close', rfq }"
                >
                    Close request
                </button>
            </article>
        </section>
        <SokoDialog
            :open="!!confirmation"
            :title="
                confirmation?.kind === 'cancel'
                    ? 'Cancel this pending order?'
                    : confirmation?.kind === 'award'
                      ? 'Select this quotation?'
                      : 'Close this sourcing request?'
            "
            @close="closeConfirmation"
            ><template v-if="confirmation"
                ><p class="form-note">
                    {{
                        confirmation.kind === 'cancel'
                            ? 'This cancels your pending order and releases its reserved inventory. The order stays in your history.'
                            : confirmation.kind === 'award'
                              ? 'This marks the request as awarded and records the selected supplier. It does not create an order or collect payment.'
                              : 'Suppliers will no longer be able to submit quotations for this request.'
                    }}
                </p>
                <p v-if="error" class="form-error" role="alert">{{ error }}</p>
                <div class="commerce-dialog-actions">
                    <button
                        class="button button-subtle"
                        :disabled="busy"
                        @click="closeConfirmation"
                    >
                        Keep reviewing</button
                    ><button
                        class="button button-navy"
                        :disabled="busy"
                        @click="confirmAction"
                    >
                        {{
                            busy
                                ? 'Saving…'
                                : confirmation.kind === 'cancel'
                                  ? 'Cancel order'
                                  : confirmation.kind === 'award'
                                    ? 'Select quotation'
                                    : 'Close request'
                        }}
                    </button>
                </div></template
            ></SokoDialog
        >
    </div>
</template>
<script setup lang="ts">
import { computed, ref } from 'vue'
import { Link, usePage } from '@inertiajs/vue3'
import SokoIcon from '@/components/SokoIcon.vue'
import SokoDialog from '@/components/SokoDialog.vue'
import { useCurrency, type Currency } from '@/composables/useCurrency'
import { useMarketplaceForm } from '@/composables/useMarketplaceForm'
import '@/../css/commerce.css'
const props = defineProps<{
    currentUser: any
    orders: any[] | null
    offers: any[] | null
    rfqs: any[] | null
    currencies?: Currency[]
}>()
const page = usePage(),
    { formatPrice } = useCurrency(props.currencies),
    { busy, error, post } = useMarketplaceForm()
const confirmation = ref<any>(null)
function closeConfirmation() {
    confirmation.value = null
    error.value = ''
}
function openConfirmation(value: any) {
    error.value = ''
    confirmation.value = value
}
const section = computed(() => {
    const value =
        new URL(page.url, 'http://soko.local').searchParams.get('section') ||
        'overview'
    return ['overview', 'orders', 'offers', 'requests'].includes(value)
        ? value
        : 'overview'
})
const tabs = computed(() => [
    { id: 'overview', label: 'Overview' },
    { id: 'orders', label: 'Orders', count: props.orders?.length || 0 },
    { id: 'offers', label: 'Offers', count: props.offers?.length || 0 },
    {
        id: 'requests',
        label: 'Sourcing requests',
        count: props.rfqs?.length || 0,
    },
])
const pendingOrders = computed(
    () => (props.orders || []).filter((o) => o.status === 'pending').length
)
function date(value: string) {
    return new Date(value).toLocaleDateString(undefined, {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
    })
}
function confirmAction() {
    const c = confirmation.value
    if (!c) return
    post(
        c.kind === 'cancel'
            ? '/account/orders/' + c.order.id + '/cancel'
            : c.kind === 'award'
              ? '/account/rfqs/' +
                c.rfq.id +
                '/quotes/' +
                c.quote.id +
                '/accept'
              : '/account/rfqs/' + c.rfq.id + '/close',
        {},
        () => {
            confirmation.value = null
        }
    )
}
</script>
