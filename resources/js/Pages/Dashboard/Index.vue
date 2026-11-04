<template>
    <div class="workspace site-container">
        <aside class="workspace-sidebar">
            <div class="workspace-identity">
                <span class="workspace-avatar">{{ initials }}</span>
                <strong>{{
                    currentUser.company_name || currentUser.name
                }}</strong>
                <span>{{
                    isAdmin ? 'Platform administration' : 'Supplier workspace'
                }}</span>
            </div>
            <nav aria-label="Workspace" class="workspace-nav">
                <Link
                    v-for="item in navigation"
                    :key="item.id"
                    :href="sectionHref(item.id)"
                    :aria-current="section === item.id ? 'page' : undefined"
                    :class="{ 'is-current': section === item.id }"
                    preserve-state
                    preserve-scroll
                >
                    <SokoIcon :name="item.icon" /><span>{{ item.label }}</span
                    ><small v-if="item.count !== undefined">{{
                        item.count
                    }}</small>
                </Link>
            </nav>
            <div class="workspace-nav workspace-secondary">
                <Link href="/account"><SokoIcon name="user" />My account</Link
                ><Link href="/"
                    ><SokoIcon name="globe" />Open marketplace<SokoIcon
                        name="arrow"
                /></Link>
            </div>
            <p class="workspace-sidebar-note">
                Your products, conversations and orders, in one place.
            </p>
        </aside>
        <div class="workspace-content">
            <header class="workspace-heading">
                <div>
                    <p class="eyebrow">
                        {{
                            isAdmin
                                ? 'Marketplace operations'
                                : 'Your business at Soko'
                        }}
                    </p>
                    <h1>{{ activeNavigation.label }}</h1>
                    <p>{{ sectionDescription }}</p>
                </div>
                <button
                    v-if="
                        canManageCatalog &&
                        ['overview', 'products'].includes(section)
                    "
                    class="button button-navy"
                    @click="openProduct()"
                >
                    <SokoIcon name="box" />Create listing
                </button>
            </header>
            <p v-if="localNotice" role="status" class="workspace-feedback">
                {{ localNotice }}
            </p>
            <p v-if="pageError" role="alert" class="form-error">
                {{ pageError }}
            </p>

            <template v-if="section === 'overview'">
                <div class="workspace-metrics">
                    <article v-for="metric in metrics" :key="metric.label">
                        <span>{{ metric.label }}</span
                        ><strong>{{ metric.value }}</strong
                        ><small>{{ metric.note }}</small>
                    </article>
                </div>
                <section class="workspace-panel workspace-next">
                    <div>
                        <p class="eyebrow">Keep business moving</p>
                        <h2>Your next steps</h2>
                        <p>
                            Start with the conversations and records that need
                            attention.
                        </p>
                    </div>
                    <div class="workspace-task-list">
                        <Link
                            v-if="canReviewOffers"
                            :href="sectionHref('offers')"
                            preserve-state
                            preserve-scroll
                            ><SokoIcon name="document" /><span
                                ><strong
                                    >{{ pendingOffers.length }} pending
                                    offers</strong
                                ><small
                                    >Review pricing proposals from buyers</small
                                ></span
                            ><SokoIcon name="arrow" /></Link
                        ><Link
                            v-if="canQuote"
                            :href="sectionHref('rfqs')"
                            preserve-state
                            preserve-scroll
                            ><SokoIcon name="globe" /><span
                                ><strong
                                    >{{ rfqList.length }} sourcing
                                    opportunities</strong
                                ><small
                                    >Respond with pricing and lead times</small
                                ></span
                            ><SokoIcon name="arrow" /></Link
                        ><Link
                            :href="sectionHref('orders')"
                            preserve-state
                            preserve-scroll
                            ><SokoIcon name="box" /><span
                                ><strong
                                    >{{ orderList.length }} orders in your
                                    workspace</strong
                                ><small
                                    >Review items, buyers and order
                                    status</small
                                ></span
                            ><SokoIcon name="arrow"
                        /></Link>
                    </div>
                </section>
                <section class="workspace-panel">
                    <div class="workspace-panel-heading">
                        <h2>Recent orders</h2>
                        <Link
                            :href="sectionHref('orders')"
                            class="text-button"
                            preserve-state
                            preserve-scroll
                            >View all orders →</Link
                        >
                    </div>
                    <div v-if="!orderList.length" class="workspace-empty">
                        <SokoIcon name="box" />
                        <h3>No orders yet</h3>
                        <p>
                            Orders will appear here when buyers submit orders
                            for your products.
                        </p>
                        <Link href="/" class="button button-outline"
                            >Explore marketplace</Link
                        >
                    </div>
                    <div v-else class="workspace-recent">
                        <article
                            v-for="order in orderList.slice(0, 5)"
                            :key="order.id"
                        >
                            <div>
                                <strong>{{ order.order_number }}</strong
                                ><small
                                    >{{ order.customer_name || 'Buyer' }} ·
                                    {{ dateLabel(order.created_at) }}</small
                                >
                            </div>
                            <span
                                class="workspace-status"
                                :data-status="order.status"
                                >{{ humanize(order.status) }}</span
                            ><strong>{{
                                formatPrice(order.total_amount)
                            }}</strong>
                        </article>
                    </div>
                </section>
            </template>

            <section
                v-else-if="section === 'content' && isAdmin"
                class="workspace-panel workspace-content-editor"
            >
                <h2>Marketplace copy</h2>
                <p>
                    Changes appear on the public storefront after you save. Keep
                    the copy specific, useful and free of unsupported
                    guarantees.
                </p>
                <form
                    class="soko-form workspace-form"
                    novalidate
                    @submit.prevent="saveContent"
                >
                    <p
                        v-if="contentError"
                        id="content-error"
                        role="alert"
                        class="form-error"
                    >
                        {{ contentError }}
                    </p>
                    <div
                        v-for="field in contentFields"
                        :key="field.key"
                        class="soko-field"
                    >
                        <label :for="`content-${field.key}`">{{
                            field.label
                        }}</label
                        ><textarea
                            class="resize-none"
                            v-if="field.multiline"
                            :id="`content-${field.key}`"
                            v-model="contentForm[field.key]"
                            :name="field.key"
                            :maxlength="field.max"
                            rows="4"
                            :aria-invalid="invalidFields.includes(field.key)"
                            :aria-describedby="
                                contentError ? 'content-error' : undefined
                            "
                        /><input
                            v-else
                            :id="`content-${field.key}`"
                            v-model="contentForm[field.key]"
                            :name="field.key"
                            :required="field.key !== 'announcement'"
                            :maxlength="field.max"
                            :aria-invalid="invalidFields.includes(field.key)"
                            :aria-describedby="
                                contentError ? 'content-error' : undefined
                            "
                        />
                        <p class="form-note">
                            {{ contentForm[field.key].length }} /
                            {{ field.max }} characters
                        </p>
                    </div>
                    <div class="workspace-form-actions">
                        <Link href="/" class="button button-outline"
                            >Preview storefront</Link
                        ><button
                            type="submit"
                            class="button button-navy"
                            :disabled="busy"
                        >
                            {{
                                busy
                                    ? 'Saving content…'
                                    : 'Save storefront content'
                            }}
                        </button>
                    </div>
                </form>
            </section>
            <template v-else>
                <div class="workspace-toolbar">
                    <label class="workspace-search"
                        ><SokoIcon name="search" /><span class="sr-only"
                            >Search
                            {{ activeNavigation.label.toLowerCase() }}</span
                        ><input
                            v-model="search"
                            type="search"
                            :placeholder="`Search ${activeNavigation.label.toLowerCase()}`" /><button
                            v-if="search"
                            type="button"
                            class="icon-button"
                            aria-label="Clear workspace search"
                            @click="search = ''"
                        >
                            <SokoIcon name="close" /></button
                    ></label>
                    <p role="status">
                        {{ visibleRecords.length }}
                        {{ visibleRecords.length === 1 ? 'record' : 'records'
                        }}<span v-if="search"> found</span>
                    </p>
                </div>
                <div
                    v-if="!visibleRecords.length"
                    class="workspace-panel workspace-empty"
                >
                    <SokoIcon :name="activeNavigation.icon" />
                    <h2>{{ search ? 'No matching records' : emptyTitle }}</h2>
                    <p>
                        {{
                            search
                                ? 'Try a different name, reference or keyword.'
                                : emptyDescription
                        }}
                    </p>
                    <button
                        v-if="search"
                        class="button button-outline"
                        @click="search = ''"
                    >
                        Clear search</button
                    ><button
                        v-else-if="section === 'products' && canManageCatalog"
                        class="button button-navy"
                        @click="openProduct()"
                    >
                        Create your first listing
                    </button>
                </div>
                <div
                    v-else-if="section === 'products'"
                    class="workspace-panel workspace-table-frame"
                >
                    <table class="workspace-table">
                        <caption class="sr-only">
                            Product catalog
                        </caption>
                        <thead>
                            <tr>
                                <th scope="col">Product</th>
                                <th scope="col">Price</th>
                                <th scope="col">Inventory</th>
                                <th scope="col">Market</th>
                                <th scope="col">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr
                                v-for="product in visibleRecords"
                                :key="product.id"
                            >
                                <td>
                                    <div class="workspace-product">
                                        <img
                                            v-if="product.image_url"
                                            :src="product.image_url"
                                            alt=""
                                            loading="lazy"
                                            @error="hideBrokenImage"
                                        /><span
                                            v-else
                                            class="workspace-product-placeholder"
                                            ><SokoIcon name="box"
                                        /></span>
                                        <div>
                                            <Link
                                                :href="`/products/${product.id}`"
                                                ><strong>{{
                                                    product.title
                                                }}</strong></Link
                                            ><small
                                                >{{ product.sku || 'No SKU' }} ·
                                                {{ humanize(product.format) }} ·
                                                {{
                                                    humanize(product.condition)
                                                }}</small
                                            ><small v-if="isAdmin">{{
                                                product.vendor_name
                                            }}</small>
                                        </div>
                                    </div>
                                </td>
                                <td>
                                    <strong>{{
                                        formatPrice(product.price)
                                    }}</strong
                                    ><small v-if="product.wholesale_price > 0"
                                        >Bulk
                                        {{
                                            formatPrice(product.wholesale_price)
                                        }}</small
                                    >
                                </td>
                                <td>
                                    <strong>{{
                                        product.stock < 0
                                            ? 'Unlimited'
                                            : `${product.stock} in stock`
                                    }}</strong
                                    ><small
                                        >Minimum {{ product.moq }} units</small
                                    >
                                </td>
                                <td>
                                    {{
                                        humanize(
                                            product.vendor_region || 'global'
                                        )
                                    }}<small>{{
                                        product.vendor_country
                                    }}</small>
                                </td>
                                <td>
                                    <div class="workspace-row-actions">
                                        <Link
                                            :href="`/products/${product.id}`"
                                            class="workspace-action"
                                            >View</Link
                                        ><button
                                            v-if="canManageCatalog"
                                            class="workspace-action"
                                            @click="openProduct(product)"
                                        >
                                            Edit<span class="sr-only">
                                                {{ product.title }}</span
                                            ></button
                                        ><button
                                            v-if="canManageCatalog"
                                            class="workspace-action is-danger"
                                            @click="openDelete(product)"
                                        >
                                            Delete<span class="sr-only">
                                                {{ product.title }}</span
                                            >
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div v-else-if="section === 'orders'" class="workspace-records">
                    <article
                        v-for="order in visibleRecords"
                        :key="order.id"
                        class="workspace-panel workspace-order"
                    >
                        <header>
                            <div>
                                <h2>{{ order.order_number }}</h2>
                                <p>
                                    {{ dateLabel(order.created_at) }} ·
                                    {{ order.customer_name || 'Buyer' }}
                                </p>
                            </div>
                            <span
                                class="workspace-status"
                                :data-status="order.status"
                                >{{ humanize(order.status) }}</span
                            >
                        </header>
                        <ul class="workspace-order-items">
                            <li
                                v-for="item in order.items || []"
                                :key="item.id || item.product_id"
                            >
                                <span
                                    ><strong>{{ item.product_title }}</strong
                                    ><small
                                        >{{ item.quantity }} ×
                                        {{
                                            formatPrice(item.unit_price)
                                        }}</small
                                    ></span
                                ><strong>{{
                                    formatPrice(item.total_price)
                                }}</strong>
                                <p
                                    v-if="
                                        item.product_type &&
                                        item.product_type !== 'physical'
                                    "
                                    class="workspace-fulfillment-note"
                                >
                                    Digital fulfillment is arranged separately
                                    with the supplier.
                                </p>
                            </li>
                        </ul>
                        <footer>
                            <span
                                >{{
                                    order.payment_method
                                        ? humanize(order.payment_method)
                                        : 'Payment method not recorded'
                                }}<small v-if="order.customer_email">{{
                                    order.customer_email
                                }}</small></span
                            ><strong
                                >Total
                                {{ formatPrice(order.total_amount) }}</strong
                            >
                        </footer>
                    </article>
                </div>
                <div v-else-if="section === 'rfqs'" class="workspace-records">
                    <article
                        v-for="rfq in visibleRecords"
                        :key="rfq.id"
                        class="workspace-panel workspace-rfq"
                    >
                        <header>
                            <div>
                                <p class="eyebrow">{{ rfq.rfq_number }}</p>
                                <h2>{{ rfq.title }}</h2>
                                <p>
                                    {{ rfq.buyer_company || rfq.buyer_name }} ·
                                    {{ rfq.buyer_country || 'Global' }}
                                </p>
                            </div>
                            <span
                                class="workspace-status"
                                :data-status="rfq.status"
                                >{{ humanize(rfq.status) }}</span
                            >
                        </header>
                        <p class="workspace-description">
                            {{
                                rfq.specifications ||
                                'No additional specifications provided.'
                            }}
                        </p>
                        <dl class="workspace-facts">
                            <div>
                                <dt>Quantity</dt>
                                <dd>
                                    {{ Number(rfq.quantity).toLocaleString() }}
                                    units
                                </dd>
                            </div>
                            <div>
                                <dt>Target unit price</dt>
                                <dd>
                                    {{
                                        rfq.target_price > 0
                                            ? formatPrice(rfq.target_price)
                                            : 'Open to quotations'
                                    }}
                                </dd>
                            </div>
                            <div>
                                <dt>Destination</dt>
                                <dd>{{ rfq.destination_country }}</dd>
                            </div>
                            <div>
                                <dt>Quotations</dt>
                                <dd>{{ (rfq.quotes || []).length }}</dd>
                            </div>
                        </dl>
                        <footer>
                            <span>{{ rfq.category }}</span
                            ><button
                                v-if="canQuote && rfq.status === 'open'"
                                class="button button-navy"
                                @click="openQuote(rfq)"
                            >
                                {{
                                    ownQuote(rfq)
                                        ? 'Update quotation'
                                        : 'Submit quotation'
                                }}<SokoIcon name="arrow" />
                            </button>
                        </footer>
                        <div v-if="ownQuote(rfq)" class="workspace-own-quote">
                            <SokoIcon name="check" /><span
                                >Your quotation:
                                {{ formatPrice(ownQuote(rfq).unit_price) }} per
                                unit · {{ ownQuote(rfq).lead_time_days }} days ·
                                {{ humanize(ownQuote(rfq).status) }}</span
                            >
                        </div>
                    </article>
                </div>
                <div v-else-if="section === 'offers'" class="workspace-records">
                    <article
                        v-for="offer in visibleRecords"
                        :key="offer.id"
                        class="workspace-panel workspace-offer"
                    >
                        <header>
                            <div>
                                <p class="eyebrow">
                                    Buyer offer ·
                                    {{ dateLabel(offer.created_at) }}
                                </p>
                                <h2>
                                    <Link
                                        :href="`/products/${offer.product_id}`"
                                        >{{ offer.product_title }}</Link
                                    >
                                </h2>
                                <p>
                                    {{
                                        offer.buyer_company || offer.buyer_name
                                    }}
                                    <span v-if="offer.buyer_email"
                                        >· {{ offer.buyer_email }}</span
                                    >
                                </p>
                            </div>
                            <span
                                class="workspace-status"
                                :data-status="offer.status"
                                >{{ humanize(offer.status) }}</span
                            >
                        </header>
                        <dl class="workspace-facts">
                            <div>
                                <dt>Offered unit price</dt>
                                <dd>{{ formatPrice(offer.offered_price) }}</dd>
                            </div>
                            <div>
                                <dt>Quantity</dt>
                                <dd>{{ offer.quantity }} units</dd>
                            </div>
                            <div>
                                <dt>Total offer</dt>
                                <dd>{{ formatPrice(offer.total_offered) }}</dd>
                            </div>
                        </dl>
                        <p v-if="offer.message" class="workspace-offer-message">
                            “{{ offer.message }}”
                        </p>
                        <footer
                            v-if="canReviewOffers && offer.status === 'pending'"
                        >
                            <span>Choose how to respond to this proposal.</span>
                            <div class="workspace-row-actions">
                                <button
                                    class="button button-outline"
                                    :disabled="busy"
                                    @click="decideOffer(offer, 'declined')"
                                >
                                    {{
                                        busyKey === `offer-${offer.id}-declined`
                                            ? 'Declining…'
                                            : 'Decline'
                                    }}</button
                                ><button
                                    class="button button-navy"
                                    :disabled="busy"
                                    @click="decideOffer(offer, 'accepted')"
                                >
                                    {{
                                        busyKey === `offer-${offer.id}-accepted`
                                            ? 'Accepting…'
                                            : 'Accept offer'
                                    }}
                                </button>
                            </div>
                        </footer>
                    </article>
                </div>
                <div
                    v-else-if="section === 'users'"
                    class="workspace-panel workspace-table-frame"
                >
                    <table class="workspace-table">
                        <caption class="sr-only">
                            Marketplace users
                        </caption>
                        <thead>
                            <tr>
                                <th scope="col">User</th>
                                <th scope="col">Role</th>
                                <th scope="col">Company</th>
                                <th scope="col">Country</th>
                                <th scope="col">Verification</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="user in visibleRecords" :key="user.id">
                                <td>
                                    <strong>{{ user.name }}</strong
                                    ><small>{{ user.email }}</small>
                                </td>
                                <td>{{ humanize(user.role) }}</td>
                                <td>
                                    {{ user.company_name || 'Not provided' }}
                                </td>
                                <td>{{ user.country || 'Not provided' }}</td>
                                <td>
                                    <span
                                        class="workspace-status"
                                        :data-status="
                                            user.verified
                                                ? 'completed'
                                                : 'pending'
                                        "
                                        >{{
                                            user.verified
                                                ? 'Verified'
                                                : 'Unverified'
                                        }}</span
                                    >
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </template>
        </div>

        <SokoDialog
            :open="productOpen"
            :title="
                editingProduct
                    ? 'Edit product listing'
                    : 'Create product listing'
            "
            @close="closeProduct"
        >
            <form
                class="soko-form workspace-form"
                novalidate
                @submit.prevent="saveProduct"
            >
                <p
                    v-if="formError"
                    id="product-errors"
                    role="alert"
                    class="form-error"
                >
                    {{ formError }}
                </p>
                <div class="soko-field">
                    <label for="product-title"
                        >Product title <span>(required)</span></label
                    ><input
                        id="product-title"
                        v-model="productForm.title"
                        name="title"
                        required
                        minlength="3"
                        maxlength="255"
                        :aria-invalid="invalidFields.includes('title')"
                        aria-describedby="product-errors"
                    />
                </div>
                <div class="form-row">
                    <div class="soko-field">
                        <label for="product-format">Selling format</label
                        ><select
                            id="product-format"
                            v-model="productForm.format"
                            @change="changeFormat"
                        >
                            <option value="wholesale">Wholesale</option>
                            <option value="retail">Ready to buy</option>
                            <option value="make_offer">Make an offer</option>
                            <option value="digital">Digital product</option>
                        </select>
                    </div>
                    <div class="soko-field">
                        <label for="product-category">Category</label
                        ><select
                            id="product-category"
                            v-model="productForm.category"
                        >
                            <option
                                v-for="category in categories"
                                :key="category"
                            >
                                {{ category }}
                            </option>
                        </select>
                    </div>
                </div>
                <div class="form-row">
                    <div class="soko-field">
                        <label for="product-sku"
                            >SKU
                            <span>(optional; generated if empty)</span></label
                        ><input
                            id="product-sku"
                            v-model="productForm.sku"
                            maxlength="64"
                        />
                    </div>
                    <div class="soko-field">
                        <label for="product-subcategory"
                            >Subcategory <span>(optional)</span></label
                        ><input
                            id="product-subcategory"
                            v-model="productForm.sub_category"
                            maxlength="128"
                        />
                    </div>
                </div>
                <div class="form-row">
                    <div class="soko-field">
                        <label for="product-price"
                            >Unit price (USD) <span>(required)</span></label
                        ><input
                            id="product-price"
                            v-model.number="productForm.price"
                            name="price"
                            type="number"
                            required
                            min="0.01"
                            step="0.01"
                            :aria-invalid="invalidFields.includes('price')"
                            aria-describedby="product-errors"
                        />
                    </div>
                    <div
                        v-if="productForm.format === 'wholesale'"
                        class="soko-field"
                    >
                        <label for="product-bulk">Bulk unit price (USD)</label
                        ><input
                            id="product-bulk"
                            v-model.number="productForm.wholesale_price"
                            name="wholesale_price"
                            type="number"
                            min="0"
                            :max="productForm.price"
                            step="0.01"
                            :aria-invalid="
                                invalidFields.includes('wholesale_price')
                            "
                            aria-describedby="product-errors"
                        />
                    </div>
                </div>
                <fieldset class="workspace-fieldset workspace-tier-editor">
                    <legend>Volume pricing <span>(optional)</span></legend>
                    <p class="form-note">
                        Set the unit price for each quantity range. Ranges must
                        be in ascending order and cannot overlap. Use 0 as the
                        last maximum quantity for no upper limit.
                    </p>
                    <p v-if="!productForm.price_tiers.length" class="form-note">
                        No volume tiers: the base price applies. For wholesale
                        quantities of 10 or more, an explicit bulk unit price
                        applies when greater than zero.
                    </p>
                    <div
                        v-for="(tier, index) in productForm.price_tiers"
                        :key="index"
                        class="workspace-tier-row"
                    >
                        <div class="soko-field">
                            <label :for="`tier-min-${index}`"
                                >Minimum units
                                <span class="sr-only"
                                    >for tier {{ index + 1 }}</span
                                ></label
                            ><input
                                :id="`tier-min-${index}`"
                                v-model.number="tier.min_qty"
                                :name="`price_tiers_${index}_min_qty`"
                                data-tier-field
                                type="number"
                                required
                                :min="productForm.moq"
                                max="1000000"
                                step="1"
                                :aria-invalid="
                                    invalidFields.includes(
                                        `price_tiers_${index}_min_qty`
                                    )
                                "
                                aria-describedby="product-errors"
                            />
                        </div>
                        <div class="soko-field">
                            <label :for="`tier-max-${index}`"
                                >Maximum units
                                <span class="sr-only"
                                    >for tier {{ index + 1 }}</span
                                ></label
                            ><input
                                :id="`tier-max-${index}`"
                                v-model.number="tier.max_qty"
                                :name="`price_tiers_${index}_max_qty`"
                                data-tier-field
                                type="number"
                                required
                                min="0"
                                max="1000000"
                                step="1"
                                :aria-invalid="
                                    invalidFields.includes(
                                        `price_tiers_${index}_max_qty`
                                    )
                                "
                                aria-describedby="product-errors"
                            />
                        </div>
                        <div class="soko-field">
                            <label :for="`tier-price-${index}`"
                                >Unit price (USD)
                                <span class="sr-only"
                                    >for tier {{ index + 1 }}</span
                                ></label
                            ><input
                                :id="`tier-price-${index}`"
                                v-model.number="tier.price"
                                :name="`price_tiers_${index}_price`"
                                data-tier-field
                                type="number"
                                required
                                min="0.01"
                                step="0.01"
                                :aria-invalid="
                                    invalidFields.includes(
                                        `price_tiers_${index}_price`
                                    )
                                "
                                aria-describedby="product-errors"
                            />
                        </div>
                        <button
                            type="button"
                            class="workspace-action is-danger"
                            :disabled="busy"
                            @click="productForm.price_tiers.splice(index, 1)"
                        >
                            Remove tier {{ index + 1 }}
                        </button>
                    </div>
                    <button
                        type="button"
                        class="button button-outline"
                        :disabled="busy || !canAddTier"
                        @click="addPriceTier"
                    >
                        Add price tier
                    </button>
                    <p
                        v-if="productForm.price_tiers.length >= 20"
                        class="form-note"
                    >
                        A listing supports up to 20 price tiers.
                    </p>
                    <p
                        v-else-if="
                            productForm.price_tiers.length && !canAddTier
                        "
                        class="form-note"
                    >
                        Give the last tier an upper limit before adding another
                        range.
                    </p>
                </fieldset>
                <div class="form-row">
                    <div class="soko-field">
                        <label for="product-moq">Minimum order quantity</label
                        ><input
                            id="product-moq"
                            v-model.number="productForm.moq"
                            name="moq"
                            required
                            type="number"
                            min="1"
                            step="1"
                            :aria-invalid="invalidFields.includes('moq')"
                            aria-describedby="product-errors"
                        />
                    </div>
                    <div class="soko-field">
                        <label for="product-stock">Stock available</label
                        ><input
                            id="product-stock"
                            v-model.number="productForm.stock"
                            name="stock"
                            required
                            type="number"
                            :min="productForm.format === 'digital' ? -1 : 0"
                            step="1"
                            :aria-invalid="invalidFields.includes('stock')"
                            :aria-describedby="
                                productForm.format === 'digital'
                                    ? 'stock-guidance product-errors'
                                    : 'product-errors'
                            "
                        />
                        <p
                            v-if="productForm.format === 'digital'"
                            id="stock-guidance"
                            class="form-note"
                        >
                            Enter available units, or -1 for unlimited digital
                            inventory.
                        </p>
                    </div>
                </div>
                <div class="form-row">
                    <div class="soko-field">
                        <label for="product-condition">Condition</label
                        ><select
                            id="product-condition"
                            v-model="productForm.condition"
                        >
                            <option value="new">New</option>
                            <option value="refurbished">Refurbished</option>
                            <option value="open_box">Open box</option>
                        </select>
                    </div>
                    <div class="soko-field">
                        <label for="product-region">Market region</label
                        ><select
                            id="product-region"
                            v-model="productForm.region"
                        >
                            <option
                                v-for="region in marketRegions"
                                :key="region.id"
                                :value="region.id"
                            >
                                {{ region.name }}
                            </option>
                        </select>
                    </div>
                </div>
                <div class="soko-field">
                    <label for="product-image"
                        >Product image URL <span>(optional)</span></label
                    ><input
                        id="product-image"
                        v-model="productForm.image_url"
                        type="url"
                        name="image_url"
                        placeholder="https://…"
                        :aria-invalid="invalidFields.includes('image_url')"
                        aria-describedby="product-errors"
                    />
                </div>
                <div class="soko-field">
                    <label for="product-description"
                        >Description and specifications</label
                    ><textarea
                        class="resize-none"
                        id="product-description"
                        v-model="productForm.description"
                        rows="4"
                    />
                </div>
                <fieldset
                    v-if="productForm.format === 'digital'"
                    class="workspace-fieldset"
                >
                    <legend>Digital delivery</legend>
                    <div class="soko-field">
                        <label for="product-type">Product type</label
                        ><select id="product-type" v-model="productForm.type">
                            <option value="digital_software">Software</option>
                            <option value="digital_asset">Digital asset</option>
                            <option value="digital_license">License</option>
                        </select>
                    </div>
                    <div class="soko-field">
                        <label for="product-download"
                            >Download URL <span>(optional)</span></label
                        ><input
                            id="product-download"
                            v-model="productForm.digital_download_url"
                            type="url"
                            name="digital_download_url"
                            placeholder="https://…"
                            :aria-invalid="
                                invalidFields.includes('digital_download_url')
                            "
                            aria-describedby="product-errors"
                        />
                    </div>
                    <div class="form-row">
                        <div class="soko-field">
                            <label for="product-license"
                                >License key format
                                <span>(optional)</span></label
                            ><input
                                id="product-license"
                                v-model="productForm.digital_license_format"
                                placeholder="SOKO-XXXX-XXXX"
                            />
                        </div>
                        <div class="soko-field">
                            <label for="product-size"
                                >File size <span>(optional)</span></label
                            ><input
                                id="product-size"
                                v-model="productForm.digital_file_size"
                                placeholder="e.g. 24 MB"
                            />
                        </div>
                    </div>
                </fieldset>
                <label
                    v-if="productForm.format !== 'digital'"
                    class="workspace-checkbox"
                    ><input
                        v-model="productForm.accepts_offers"
                        type="checkbox"
                    />Accept buyer offers on this listing</label
                >
                <p class="form-note">
                    Prices are entered in USD. Display currency conversion does
                    not change the listing's base price.
                </p>
                <p v-if="editingProduct" class="form-note">
                    If inventory changes while you are editing, reopen the
                    listing before trying again.
                </p>
                <div class="workspace-form-actions">
                    <button
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="closeProduct"
                    >
                        Cancel</button
                    ><button
                        type="submit"
                        class="button button-navy"
                        :disabled="busy"
                    >
                        {{
                            busy
                                ? 'Saving listing…'
                                : editingProduct
                                  ? 'Save changes'
                                  : 'Publish listing'
                        }}
                    </button>
                </div>
            </form>
        </SokoDialog>
        <SokoDialog
            :open="!!deleteTarget"
            title="Delete product listing?"
            @close="closeDelete"
            ><form class="soko-form" novalidate @submit.prevent="removeProduct">
                <p>
                    This removes <strong>{{ deleteTarget?.title }}</strong> from
                    the marketplace. Existing order records will remain.
                </p>
                <p v-if="formError" role="alert" class="form-error">
                    {{ formError }}
                </p>
                <div class="workspace-form-actions">
                    <button
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="closeDelete"
                    >
                        Keep listing</button
                    ><button
                        type="submit"
                        class="button workspace-button-danger"
                        :disabled="busy"
                    >
                        {{ busy ? 'Deleting…' : 'Delete listing' }}
                    </button>
                </div>
            </form></SokoDialog
        >
        <SokoDialog
            :open="!!quoteTarget"
            :title="
                quoteTarget && ownQuote(quoteTarget)
                    ? 'Update your quotation'
                    : 'Submit a quotation'
            "
            @close="closeQuote"
            ><form
                class="soko-form workspace-form"
                novalidate
                @submit.prevent="saveQuote"
            >
                <div class="workspace-quote-context">
                    <strong>{{ quoteTarget?.title }}</strong>
                    <p>
                        {{ quoteTarget?.quantity }} units ·
                        {{ quoteTarget?.destination_country }}
                    </p>
                </div>
                <p
                    v-if="formError"
                    id="quote-errors"
                    role="alert"
                    class="form-error"
                >
                    {{ formError }}
                </p>
                <div class="form-row">
                    <div class="soko-field">
                        <label for="quote-price">Unit price (USD)</label
                        ><input
                            id="quote-price"
                            v-model.number="quoteForm.unit_price"
                            type="number"
                            name="unit_price"
                            required
                            min="0.01"
                            step="0.01"
                            :aria-invalid="invalidFields.includes('unit_price')"
                            aria-describedby="quote-errors"
                        />
                    </div>
                    <div class="soko-field">
                        <label for="quote-lead">Lead time (days)</label
                        ><input
                            id="quote-lead"
                            v-model.number="quoteForm.lead_time_days"
                            type="number"
                            name="lead_time_days"
                            required
                            min="1"
                            step="1"
                            :aria-invalid="
                                invalidFields.includes('lead_time_days')
                            "
                            aria-describedby="quote-errors"
                        />
                    </div>
                </div>
                <div class="soko-field">
                    <label for="quote-notes">Notes and shipping terms</label
                    ><textarea
                        class="resize-none"
                        id="quote-notes"
                        v-model="quoteForm.notes"
                        rows="4"
                        placeholder="Describe what's included, shipping terms and any conditions."
                    />
                </div>
                <div class="workspace-form-actions">
                    <button
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="closeQuote"
                    >
                        Cancel</button
                    ><button
                        type="submit"
                        class="button button-navy"
                        :disabled="busy"
                    >
                        {{ busy ? 'Sending quotation…' : 'Send quotation' }}
                    </button>
                </div>
            </form></SokoDialog
        >
    </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { Link, router, usePage } from '@inertiajs/vue3'
import SokoDialog from '@/components/SokoDialog.vue'
import SokoIcon from '@/components/SokoIcon.vue'
import {
    useCurrency,
    type Currency,
    type MarketRegion,
} from '@/composables/useCurrency'
import '../../../css/workspace.css'

type RecordData = Record<string, any>
const props = defineProps<{
    currentUser: RecordData
    stats: Record<string, number>
    products?: RecordData[] | null
    orders?: RecordData[] | null
    users?: RecordData[] | null
    rfqs?: RecordData[] | null
    offers?: RecordData[] | null
    currencies?: Currency[]
    regions?: MarketRegion[]
    storefrontContent?: Record<string, string>
    capabilities?: {
        manageCatalog?: boolean
        manageUsers?: boolean
        reviewOffers?: boolean
        quoteRfqs?: boolean
    }
}>()
const page = usePage()
const { formatPrice } = useCurrency(props.currencies)
const isAdmin = computed(() => props.currentUser.role === 'admin')
const canManageCatalog = computed(
    () =>
        props.capabilities?.manageCatalog ??
        ['admin', 'merchant'].includes(props.currentUser.role)
)
const canReviewOffers = computed(
    () =>
        props.capabilities?.reviewOffers ??
        ['admin', 'merchant'].includes(props.currentUser.role)
)
const canQuote = computed(
    () => props.capabilities?.quoteRfqs ?? props.currentUser.role === 'merchant'
)
const canManageUsers = computed(
    () => props.capabilities?.manageUsers ?? isAdmin.value
)
const productList = computed(() => props.products || [])
const orderList = computed(() => props.orders || [])
const rfqList = computed(() => props.rfqs || [])
const offerList = computed(() => props.offers || [])
const pendingOffers = computed(() =>
    offerList.value.filter((item) => item.status === 'pending')
)
const initials = computed(() =>
    String(props.currentUser.company_name || props.currentUser.name || 'Soko')
        .split(' ')
        .slice(0, 2)
        .map((part: string) => part[0])
        .join('')
        .toUpperCase()
)
const navigation = computed(() => [
    { id: 'overview', label: 'Overview', icon: 'layers' },
    ...(canManageCatalog.value
        ? [
              {
                  id: 'products',
                  label: 'Products',
                  icon: 'box',
                  count: productList.value.length,
              },
          ]
        : []),
    {
        id: 'orders',
        label: 'Orders',
        icon: 'document',
        count: orderList.value.length,
    },
    {
        id: 'rfqs',
        label: 'Sourcing requests',
        icon: 'globe',
        count: rfqList.value.length,
    },
    {
        id: 'offers',
        label: 'Buyer offers',
        icon: 'document',
        count: offerList.value.length,
    },
    ...(canManageUsers.value
        ? [
              {
                  id: 'users',
                  label: 'Users',
                  icon: 'user',
                  count: (props.users || []).length,
              },
              { id: 'content', label: 'Storefront content', icon: 'document' },
          ]
        : []),
])
const section = computed(() => {
    const requested =
        new URL(page.url, 'http://soko.local').searchParams.get('section') ||
        'overview'
    return navigation.value.some((item) => item.id === requested)
        ? requested
        : 'overview'
})
const activeNavigation = computed(
    () =>
        navigation.value.find((item) => item.id === section.value) ||
        navigation.value[0]
)
const descriptions: Record<string, string> = {
    overview: 'A clear view of your marketplace activity.',
    products: 'Manage the details buyers need to choose your products.',
    orders: 'Review purchases, ordered items and delivery information.',
    rfqs: 'Discover what buyers need and prepare a quotation.',
    offers: 'Review buyer proposals and keep negotiations moving.',
    users: 'View registered marketplace accounts and verification status.',
    content: 'Update the words buyers see in the marketplace.',
}
const sectionDescription = computed(() => descriptions[section.value])
const sectionHref = (id: string) => `/dashboard?section=${id}`
const search = ref('')
const pageError = ref('')
const localNotice = ref('')
watch(section, () => {
    search.value = ''
    pageError.value = ''
    localNotice.value = ''
})
const records = computed<RecordData[]>(
    () =>
        ({
            products: productList.value,
            orders: orderList.value,
            rfqs: rfqList.value,
            offers: offerList.value,
            users: props.users || [],
        })[section.value] || []
)
const visibleRecords = computed(() => {
    const query = search.value.trim().toLowerCase()
    return query
        ? records.value.filter((item) =>
              [
                  'title',
                  'sku',
                  'category',
                  'vendor_name',
                  'order_number',
                  'customer_name',
                  'customer_email',
                  'rfq_number',
                  'buyer_company',
                  'buyer_name',
                  'product_title',
                  'name',
                  'email',
                  'company_name',
                  'status',
                  'role',
              ].some((key) =>
                  String(item[key] || '')
                      .toLowerCase()
                      .includes(query)
              )
          )
        : records.value
})
const emptyTitle = computed(
    () =>
        ({
            products: 'Your catalog starts here',
            orders: 'No orders yet',
            rfqs: 'No open sourcing requests',
            offers: 'No buyer offers yet',
            users: 'No users to display',
        })[section.value]
)
const emptyDescription = computed(
    () =>
        ({
            products:
                'Create a listing with clear pricing, inventory and specifications.',
            orders: 'New orders will appear here when they are placed.',
            rfqs: 'Return here to discover new buyer requirements.',
            offers: 'Buyer proposals will appear when they submit an offer on your listings.',
            users: 'Registered accounts will appear in this directory.',
        })[section.value]
)
const metrics = computed(() => [
    {
        label: isAdmin.value ? 'Marketplace order value' : 'Order value',
        value: formatPrice(props.stats?.total_revenue ?? 0),
        note: 'Completed and shipped orders',
    },
    {
        label: 'Pending orders',
        value: props.stats?.active_orders ?? 0,
        note: 'Awaiting review',
    },
    {
        label: 'Product listings',
        value: props.stats?.total_products ?? productList.value.length,
        note: isAdmin.value ? 'Across the marketplace' : 'In your catalog',
    },
    {
        label: 'Open conversations',
        value:
            (props.stats?.pending_rfqs ?? 0) +
            (props.stats?.active_negotiations ?? 0),
        note: `${props.stats?.pending_rfqs ?? 0} requests · ${props.stats?.active_negotiations ?? 0} offers`,
    },
])
const humanize = (value: unknown) =>
    String(value || 'Not recorded')
        .replace(/[_-]/g, ' ')
        .replace(/^./, (char) => char.toUpperCase())
const dateLabel = (value: string) => {
    const date = new Date(value)
    return Number.isNaN(date.getTime())
        ? 'Date not recorded'
        : date.toLocaleDateString(undefined, {
              day: 'numeric',
              month: 'short',
              year: 'numeric',
          })
}
const hideBrokenImage = (event: Event) => {
    const image = event.target as HTMLImageElement
    image.style.display = 'none'
}
const busy = ref(false)
const busyKey = ref('')
const formError = ref('')
const invalidFields = ref<string[]>([])
const validate = (event: Event) => {
    formError.value = ''
    invalidFields.value = []
    const form = event.target as HTMLFormElement
    const invalid = Array.from(form.elements)
        .filter(
            (
                field
            ): field is
                HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement =>
                field instanceof HTMLInputElement ||
                field instanceof HTMLTextAreaElement ||
                field instanceof HTMLSelectElement
        )
        .filter((field) => !field.checkValidity())
    if (!invalid.length) return true
    invalidFields.value = invalid.map((field) => field.name)
    const first = invalid[0]
    const label =
        form
            .querySelector(`label[for="${first.id}"]`)
            ?.textContent?.replace(/\s*\(.*?\)/g, '')
            .trim() || 'this field'
    formError.value = first.validity.customError
        ? first.validationMessage
        : first.validity.valueMissing
          ? `Enter ${label.toLowerCase()}.`
          : first.validity.tooShort
            ? `Enter at least ${'minLength' in first ? first.minLength : 0} characters for ${label.toLowerCase()}.`
            : `Check ${label.toLowerCase()} and enter a valid value.`
    first.focus()
    return false
}
const mutate = (
    url: string,
    data: RecordData,
    key: string,
    complete: () => void,
    modal = false
) => {
    if (busy.value) return
    busy.value = true
    busyKey.value = key
    formError.value = ''
    pageError.value = ''
    router.post(url, data, {
        preserveScroll: true,
        preserveState: true,
        onSuccess: (response) => {
            const error = (response.props.flash as RecordData)?.error
            if (error) {
                if (modal) formError.value = String(error)
                else if (section.value === 'content')
                    contentError.value = String(error)
                else pageError.value = String(error)
            } else complete()
        },
        onError: (errors) => {
            const message =
                Object.values(errors).join(' ') ||
                'Your changes could not be saved. Please try again.'
            if (modal) formError.value = message
            else if (section.value === 'content') contentError.value = message
            else pageError.value = message
        },
        onFinish: () => {
            busy.value = false
            busyKey.value = ''
        },
    })
}
const productOpen = ref(false)
const editingProduct = ref<RecordData | null>(null)
const emptyProduct = () => ({
    version: '',
    price_tiers: [] as { min_qty: number; max_qty: number; price: number }[],
    title: '',
    category: 'Industrial & Machinery',
    sub_category: '',
    type: 'physical',
    format: 'wholesale',
    condition: 'new',
    region: 'global',
    price: 0,
    wholesale_price: 0,
    moq: 1,
    stock: 0,
    sku: '',
    description: '',
    image_url: '',
    digital_license_format: '',
    digital_file_size: '',
    digital_download_url: '',
    accepts_offers: false,
})
const productForm = reactive(emptyProduct())
const categories = [
    'Industrial & Machinery',
    'Electronics & Telecom',
    'Digital Software',
    'Digital Assets',
    'Agriculture & Processing',
    'Energy & Utilities',
]
const marketRegions = computed(() =>
    props.regions?.length
        ? props.regions
        : [
              { id: 'global', name: 'Global' },
              { id: 'east-africa', name: 'East Africa' },
              { id: 'west-africa', name: 'West Africa' },
              { id: 'north-america', name: 'North America' },
              { id: 'europe', name: 'Europe & UK' },
              { id: 'asia-pacific', name: 'Asia-Pacific' },
              { id: 'middle-east', name: 'Middle East' },
          ]
)
const openProduct = (product?: RecordData) => {
    editingProduct.value = product || null
    Object.assign(productForm, emptyProduct())
    if (product) {
        for (const key of Object.keys(emptyProduct()))
            if (product[key] !== undefined)
                (productForm as RecordData)[key] = product[key]
        productForm.region = product.vendor_region || 'global'
        productForm.version = String(product.updated_at || '')
        productForm.price_tiers = (product.price_tiers || []).map(
            (tier: { min_qty: number; max_qty: number; price: number }) => ({
                ...tier,
            })
        )
    }
    formError.value = ''
    invalidFields.value = []
    productOpen.value = true
}
const canAddTier = computed(() => {
    const tiers = productForm.price_tiers
    return (
        tiers.length < 20 &&
        (!tiers.length || tiers[tiers.length - 1].max_qty > 0)
    )
})
const addPriceTier = () => {
    if (!canAddTier.value || busy.value) return
    const tiers = productForm.price_tiers
    const minimum = tiers.length
        ? Math.max(productForm.moq, Number(tiers[tiers.length - 1].max_qty) + 1)
        : productForm.moq
    tiers.push({ min_qty: minimum, max_qty: 0, price: productForm.price })
}
const validatePriceTiers = (form: HTMLFormElement) => {
    form.querySelectorAll<HTMLInputElement>('[data-tier-field]').forEach(
        (input) => input.setCustomValidity('')
    )
    productForm.price_tiers.forEach((tier, index) => {
        const previous = productForm.price_tiers[index - 1]
        const minimum = form.querySelector<HTMLInputElement>(
            `#tier-min-${index}`
        )
        const maximum = form.querySelector<HTMLInputElement>(
            `#tier-max-${index}`
        )
        if (
            previous &&
            (previous.max_qty === 0 || tier.min_qty <= previous.max_qty)
        )
            minimum?.setCustomValidity(
                'Price tier ranges must be in ascending order and cannot overlap. Only the last tier can have no upper limit.'
            )
        if (tier.max_qty !== 0 && tier.max_qty < tier.min_qty)
            maximum?.setCustomValidity(
                'Maximum quantity must be zero for no upper limit, or at least the minimum quantity.'
            )
    })
}
const changeFormat = () => {
    if (productForm.format === 'digital') {
        productForm.type = 'digital_software'
        productForm.accepts_offers = false
    } else {
        if (productForm.type !== 'physical')
            productForm.stock = Math.max(0, productForm.stock)
        productForm.type = 'physical'
        productForm.digital_download_url = ''
        productForm.digital_file_size = ''
        productForm.digital_license_format = ''
        productForm.accepts_offers = productForm.format === 'make_offer'
    }
}
const closeProduct = () => {
    if (!busy.value) productOpen.value = false
}
const saveProduct = (event: Event) => {
    if (busy.value) return
    validatePriceTiers(event.target as HTMLFormElement)
    if (!validate(event)) return
    const id = editingProduct.value?.id
    mutate(
        `${id ? `/dashboard/products/${id}` : '/dashboard/products'}?section=products`,
        { ...productForm },
        'product',
        () => {
            productOpen.value = false
            localNotice.value = id
                ? 'Listing changes saved.'
                : 'Listing published.'
            router.get(
                sectionHref('products'),
                {},
                { preserveState: true, preserveScroll: true }
            )
        },
        true
    )
}
const deleteTarget = ref<RecordData | null>(null)
const openDelete = (product: RecordData) => {
    deleteTarget.value = product
    formError.value = ''
}
const closeDelete = () => {
    if (!busy.value) deleteTarget.value = null
}
const removeProduct = () => {
    if (!deleteTarget.value) return
    mutate(
        `/dashboard/products/${deleteTarget.value.id}/delete?section=products`,
        {},
        'delete',
        () => {
            deleteTarget.value = null
            localNotice.value = 'Listing deleted.'
        },
        true
    )
}
const quoteTarget = ref<RecordData | null>(null)
const quoteForm = reactive({ unit_price: 0, lead_time_days: 7, notes: '' })
const ownQuote = (rfq: RecordData) =>
    (rfq.quotes || []).find(
        (quote: RecordData) => quote.vendor_id === props.currentUser.id
    )
const openQuote = (rfq: RecordData) => {
    quoteTarget.value = rfq
    const quote = ownQuote(rfq)
    Object.assign(quoteForm, {
        unit_price: quote?.unit_price ?? rfq.target_price ?? 0,
        lead_time_days: quote?.lead_time_days ?? 7,
        notes: quote?.notes ?? '',
    })
    formError.value = ''
    invalidFields.value = []
}
const closeQuote = () => {
    if (!busy.value) quoteTarget.value = null
}
const saveQuote = (event: Event) => {
    if (busy.value || !quoteTarget.value || !validate(event)) return
    mutate(
        '/dashboard/rfqs/quote?section=rfqs',
        { rfq_id: quoteTarget.value.id, ...quoteForm },
        'quote',
        () => {
            quoteTarget.value = null
            localNotice.value = 'Quotation sent.'
        },
        true
    )
}
const decideOffer = (offer: RecordData, status: 'accepted' | 'declined') =>
    mutate(
        `/dashboard/offers/${offer.id}/decision?section=offers`,
        { status },
        `offer-${offer.id}-${status}`,
        () => {
            localNotice.value =
                status === 'accepted' ? 'Offer accepted.' : 'Offer declined.'
        }
    )
const contentFields = [
    { key: 'hero_heading', label: 'Hero heading', max: 120 },
    { key: 'hero_accent', label: 'Hero accent line', max: 120 },
    {
        key: 'hero_description',
        label: 'Hero description',
        max: 500,
        multiline: true,
    },
    { key: 'announcement', label: 'Announcement', max: 160 },
    { key: 'sourcing_heading', label: 'Sourcing heading', max: 120 },
    {
        key: 'sourcing_description',
        label: 'Sourcing description',
        max: 500,
        multiline: true,
    },
]
const contentForm = reactive<Record<string, string>>(
    Object.fromEntries(
        contentFields.map((field) => [
            field.key,
            props.storefrontContent?.[field.key] || '',
        ])
    )
)
const contentError = ref('')
const saveContent = (event: Event) => {
    if (busy.value || !validate(event)) {
        contentError.value = formError.value
        return
    }
    contentError.value = ''
    mutate(
        '/dashboard/content?section=content',
        { ...contentForm },
        'content',
        () => {
            localNotice.value = 'Storefront content saved.'
        }
    )
}
</script>
