<template>
    <div class="workspace site-container">
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
                        capabilities?.createProducts &&
                        ['overview', 'products'].includes(section)
                    "
                    class="button button-navy"
                    @click="openProduct()"
                >
                    <SokoIcon name="box" />Create listing
                </button>
            </header>
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
                v-else-if="section === 'content' && capabilities?.manageContent"
                class="workspace-panel workspace-content-editor"
            >
                <h2>Storefront content and product selection</h2>
                <p>
                    Control the announcement, section titles, category
                    navigation, main product showcase, side cards and digital
                    discovery from the versioned storefront editor.
                </p>
                <Link href="/dashboard/merchandising" class="button button-navy"
                    >Open storefront editor</Link
                >
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
                    <label
                        v-if="['orders', 'rfqs', 'offers'].includes(section)"
                        class="operation-status-filter"
                        >Status<select v-model="recordStatus">
                            <option value="all">All statuses</option>
                            <option
                                v-for="status in statusOptions"
                                :key="status"
                                :value="status"
                            >
                                {{ humanize(status) }}
                            </option>
                        </select></label
                    >
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
                    <h2>
                        {{
                            search || recordStatus !== 'all'
                                ? 'No matching records'
                                : emptyTitle
                        }}
                    </h2>
                    <p>
                        {{
                            search
                                ? 'Try a different name, reference or keyword.'
                                : emptyDescription
                        }}
                    </p>
                    <button
                        v-if="search || recordStatus !== 'all'"
                        class="button button-outline"
                        @click="clearRecordFilters"
                    >
                        Clear filters</button
                    ><button
                        v-else-if="
                            section === 'products' &&
                            capabilities?.createProducts
                        "
                        class="button button-navy"
                        @click="openProduct()"
                    >
                        Create your first listing
                    </button>
                </div>
                <AdminTable
                    v-else-if="section === 'products'"
                    caption="Product catalog"
                    :busy="busy"
                    class="workspace-admin-list"
                    ><template #header>
                        <tr>
                            <th scope="col">Product</th>
                            <th scope="col">Price</th>
                            <th scope="col">Inventory</th>
                            <th scope="col">Market</th>
                            <th scope="col">Actions</th>
                        </tr>
                    </template>
                    <tr v-for="product in visibleRecords" :key="product.id">
                        <th scope="row">
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
                                    <Link :href="`/products/${product.id}`"
                                        ><strong>{{
                                            product.title
                                        }}</strong></Link
                                    ><small
                                        >{{ product.sku || 'No SKU' }} ·
                                        {{ humanize(product.format) }} ·
                                        {{ humanize(product.condition) }}</small
                                    ><small v-if="isAdmin">{{
                                        product.vendor_name
                                    }}</small>
                                </div>
                            </div>
                        </th>
                        <td>
                            <strong>{{ formatPrice(product.price) }}</strong
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
                            ><small>Minimum {{ product.moq }} units</small>
                        </td>
                        <td>
                            {{ humanize(product.vendor_region || 'global')
                            }}<small>{{ product.vendor_country }}</small>
                        </td>
                        <td>
                            <div class="workspace-row-actions">
                                <Link
                                    :href="`/products/${product.id}`"
                                    class="admin-view-link"
                                    :aria-label="`View ${product.title} on storefront`"
                                    >View</Link
                                >
                                <AdminActionMenu
                                    :label="`Actions for ${product.title}`"
                                    :items="productActions(product)"
                                    :disabled="busy"
                                    @select="(id) => productAction(product, id)"
                                />
                            </div>
                        </td>
                    </tr>
                </AdminTable>
                <AdminTable
                    v-else-if="section === 'orders'"
                    caption="Orders · authorized lines only"
                    class="workspace-admin-list"
                >
                    <template #header
                        ><tr>
                            <th scope="col">Order</th>
                            <th scope="col">Buyer</th>
                            <th scope="col">Visible lines</th>
                            <th scope="col">Status</th>
                            <th scope="col">Visible total</th>
                            <th scope="col">Actions</th>
                        </tr></template
                    >
                    <tr v-for="order in visibleRecords" :key="order.id">
                        <th scope="row">
                            <strong>{{ order.order_number }}</strong
                            ><small>{{ dateLabel(order.created_at) }}</small>
                        </th>
                        <td>{{ order.customer_name || 'Buyer' }}</td>
                        <td>{{ (order.items || []).length }}</td>
                        <td>
                            <span
                                class="workspace-status"
                                :data-status="order.status"
                                >{{ humanize(order.status) }}</span
                            >
                        </td>
                        <td>{{ formatPrice(order.total_amount) }}</td>
                        <td>
                            <AdminActionMenu
                                :label="`Actions for order ${order.order_number}`"
                                :items="operationActions('orders', order)"
                                :disabled="busy"
                                @select="
                                    (id) => openRecord('orders', order, id)
                                "
                            />
                        </td>
                    </tr>
                </AdminTable>
                <AdminTable
                    v-else-if="section === 'rfqs'"
                    caption="Sourcing requests"
                    class="workspace-admin-list"
                >
                    <template #header
                        ><tr>
                            <th scope="col">Request</th>
                            <th scope="col">Buyer</th>
                            <th scope="col">Quantity / target</th>
                            <th scope="col">Destination</th>
                            <th scope="col">Status</th>
                            <th scope="col">Actions</th>
                        </tr></template
                    >
                    <tr v-for="rfq in visibleRecords" :key="rfq.id">
                        <th scope="row">
                            <strong>{{ rfq.title }}</strong
                            ><small
                                >{{ rfq.rfq_number }} ·
                                {{ rfq.category }}</small
                            >
                        </th>
                        <td>{{ rfq.buyer_company || rfq.buyer_name }}</td>
                        <td>
                            {{ Number(rfq.quantity).toLocaleString() }}
                            units<small>{{
                                rfq.target_price > 0
                                    ? formatPrice(rfq.target_price) + ' / unit'
                                    : 'Open to quotations'
                            }}</small>
                        </td>
                        <td>
                            {{ rfq.destination_country || 'Not specified' }}
                        </td>
                        <td>
                            <span
                                class="workspace-status"
                                :data-status="rfq.status"
                                >{{ humanize(rfq.status) }}</span
                            >
                        </td>
                        <td>
                            <AdminActionMenu
                                :label="`Actions for request ${rfq.title}`"
                                :items="operationActions('rfqs', rfq)"
                                :disabled="busy"
                                @select="(id) => openRecord('rfqs', rfq, id)"
                            />
                        </td>
                    </tr>
                </AdminTable>
                <AdminTable
                    v-else-if="section === 'offers'"
                    caption="Buyer offers"
                    class="workspace-admin-list"
                >
                    <template #header
                        ><tr>
                            <th scope="col">Product / received</th>
                            <th scope="col">Buyer</th>
                            <th scope="col">Proposal</th>
                            <th scope="col">Status</th>
                            <th scope="col">Actions</th>
                        </tr></template
                    >
                    <tr v-for="offer in visibleRecords" :key="offer.id">
                        <th scope="row">
                            <strong>{{ offer.product_title }}</strong
                            ><small>{{ dateLabel(offer.created_at) }}</small>
                        </th>
                        <td>{{ offer.buyer_company || offer.buyer_name }}</td>
                        <td>
                            {{ formatPrice(offer.offered_price) }} ×
                            {{ offer.quantity
                            }}<small
                                >Total
                                {{ formatPrice(offer.total_offered) }}</small
                            >
                        </td>
                        <td>
                            <span
                                class="workspace-status"
                                :data-status="offer.status"
                                >{{ humanize(offer.status) }}</span
                            >
                        </td>
                        <td>
                            <AdminActionMenu
                                :label="`Actions for offer on ${offer.product_title}`"
                                :items="operationActions('offers', offer)"
                                :disabled="busy"
                                @select="
                                    (id) => openRecord('offers', offer, id)
                                "
                            />
                        </td>
                    </tr>
                </AdminTable>
                <AdminTable
                    v-else-if="section === 'users'"
                    caption="Marketplace users"
                    :busy="busy"
                    class="workspace-admin-list"
                    ><template #header>
                        <tr>
                            <th scope="col">User</th>
                            <th scope="col">Role</th>
                            <th scope="col">Company</th>
                            <th scope="col">Country</th>
                            <th scope="col">Verification</th>
                            <th scope="col">Actions</th>
                        </tr>
                    </template>
                    <tr v-for="user in visibleRecords" :key="user.id">
                        <th scope="row">
                            <strong>{{ user.name }}</strong
                            ><small>{{ user.email }}</small>
                        </th>
                        <td>{{ humanize(user.role) }}</td>
                        <td>
                            {{ user.company_name || 'Not provided' }}
                        </td>
                        <td>{{ user.country || 'Not provided' }}</td>
                        <td>
                            <span
                                class="workspace-status"
                                :data-status="
                                    user.verified ? 'completed' : 'pending'
                                "
                                >{{
                                    user.verified ? 'Verified' : 'Unverified'
                                }}</span
                            >
                        </td>
                        <td>
                            <button
                                type="button"
                                class="admin-view-link"
                                :aria-label="`View account ${user.name}`"
                                @click="openRecord('users', user)"
                            >
                                View account
                            </button>
                        </td>
                    </tr>
                </AdminTable>
            </template>
        </div>

        <AdminDrawer
            :open="!!detailRecord"
            :title="detailTitle"
            :description="
                detailKind === 'users'
                    ? 'Account directory · read-only'
                    : 'Only information authorized for your account is shown.'
            "
            :busy="busy"
            @close="closeRecord"
        >
            <div v-if="detailRecord" class="operation-details">
                <div v-if="formError" class="operation-error">
                    <p role="alert" class="form-error">{{ formError }}</p>
                    <p class="operation-scope-note">
                        The reviewed terms have not been changed automatically.
                        Refresh to load current details, then review and confirm
                        again.
                    </p>
                    <button
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="refreshRecord"
                    >
                        Refresh details and review again
                    </button>
                </div>
                <template v-if="detailKind === 'orders'">
                    <p class="operation-scope-note">
                        This view contains the order lines you are authorized to
                        see. The visible total may represent only your portion
                        of an order.
                    </p>
                    <dl class="operation-facts">
                        <div>
                            <dt>Buyer</dt>
                            <dd>
                                {{ detailRecord.customer_name || 'Buyer'
                                }}<small>{{
                                    detailRecord.customer_email
                                }}</small>
                            </dd>
                        </div>
                        <div>
                            <dt>Status</dt>
                            <dd>{{ humanize(detailRecord.status) }}</dd>
                        </div>
                        <div>
                            <dt>Payment</dt>
                            <dd>
                                {{
                                    detailRecord.payment_method ||
                                    'Payment method not recorded'
                                }}
                            </dd>
                        </div>
                        <div>
                            <dt>Visible total</dt>
                            <dd>
                                {{ formatPrice(detailRecord.total_amount) }}
                            </dd>
                        </div>
                    </dl>
                    <h3>Visible order lines</h3>
                    <ul class="operation-lines">
                        <li
                            v-for="item in detailRecord.items || []"
                            :key="item.id || item.product_id"
                        >
                            <strong>{{ item.product_title }}</strong
                            ><span
                                >{{ item.quantity }} ×
                                {{ formatPrice(item.unit_price) }} ·
                                {{ formatPrice(item.total_price) }}</span
                            ><small
                                v-if="
                                    item.product_type &&
                                    item.product_type !== 'physical'
                                "
                                >Digital fulfillment is arranged separately with
                                the supplier.</small
                            >
                        </li>
                    </ul>
                </template>
                <template v-else-if="detailKind === 'rfqs'">
                    <dl class="operation-facts">
                        <div>
                            <dt>Buyer</dt>
                            <dd>
                                {{
                                    detailRecord.buyer_company ||
                                    detailRecord.buyer_name
                                }}<small>{{
                                    detailRecord.buyer_country
                                }}</small>
                            </dd>
                        </div>
                        <div>
                            <dt>Quantity</dt>
                            <dd>{{ detailRecord.quantity }} units</dd>
                        </div>
                        <div>
                            <dt>Destination</dt>
                            <dd>
                                {{
                                    detailRecord.destination_country ||
                                    'Not specified'
                                }}
                            </dd>
                        </div>
                        <div>
                            <dt>Status</dt>
                            <dd>{{ humanize(detailRecord.status) }}</dd>
                        </div>
                    </dl>
                    <h3>Requirements</h3>
                    <p class="operation-message">
                        {{
                            detailRecord.specifications ||
                            'No additional specifications provided.'
                        }}
                    </p>
                    <h3>Visible quotations</h3>
                    <p class="operation-scope-note">
                        Only quotations your account can access are listed. This
                        is not necessarily the full set of competing quotes.
                    </p>
                    <p v-if="!(detailRecord.quotes || []).length">
                        No quotations are visible to your account.
                    </p>
                    <ul class="operation-lines">
                        <li
                            v-for="quote in detailRecord.quotes || []"
                            :key="quote.id"
                        >
                            <strong>{{
                                quote.vendor_company ||
                                quote.vendor_name ||
                                'Supplier quotation'
                            }}</strong
                            ><span
                                >{{ formatPrice(quote.unit_price) }} / unit ·
                                {{ quote.lead_time_days }} days ·
                                {{ humanize(quote.status) }}</span
                            >
                            <p v-if="quote.notes" class="operation-message">
                                {{ quote.notes }}
                            </p>
                            <button
                                v-if="recordCan('quotes', quote.id, 'accept')"
                                type="button"
                                class="button button-outline"
                                :disabled="busy"
                                @click="prepareOperation('award', quote)"
                            >
                                Review award
                            </button>
                        </li>
                    </ul>
                </template>
                <template v-else-if="detailKind === 'offers'">
                    <Link
                        :href="`/products/${detailRecord.product_id}`"
                        class="admin-view-link"
                        >View product</Link
                    >
                    <dl class="operation-facts">
                        <div>
                            <dt>Buyer</dt>
                            <dd>
                                {{
                                    detailRecord.buyer_company ||
                                    detailRecord.buyer_name
                                }}<small>{{ detailRecord.buyer_email }}</small>
                            </dd>
                        </div>
                        <div>
                            <dt>Unit proposal</dt>
                            <dd>
                                {{ formatPrice(detailRecord.offered_price) }}
                            </dd>
                        </div>
                        <div>
                            <dt>Quantity</dt>
                            <dd>{{ detailRecord.quantity }} units</dd>
                        </div>
                        <div>
                            <dt>Total proposal</dt>
                            <dd>
                                {{ formatPrice(detailRecord.total_offered) }}
                            </dd>
                        </div>
                        <div>
                            <dt>Status</dt>
                            <dd>{{ humanize(detailRecord.status) }}</dd>
                        </div>
                    </dl>
                    <h3>Buyer message</h3>
                    <p class="operation-message">
                        {{ detailRecord.message || 'No message supplied.' }}
                    </p>
                </template>
                <template v-else-if="detailKind === 'users'">
                    <dl class="operation-facts">
                        <div>
                            <dt>Email</dt>
                            <dd>{{ detailRecord.email }}</dd>
                        </div>
                        <div>
                            <dt>Role</dt>
                            <dd>{{ humanize(detailRecord.role) }}</dd>
                        </div>
                        <div>
                            <dt>Company</dt>
                            <dd>
                                {{
                                    detailRecord.company_name || 'Not provided'
                                }}
                            </dd>
                        </div>
                        <div>
                            <dt>Country</dt>
                            <dd>
                                {{ detailRecord.country || 'Not provided' }}
                            </dd>
                        </div>
                        <div>
                            <dt>Verification</dt>
                            <dd>
                                {{
                                    detailRecord.verified
                                        ? 'Verified'
                                        : 'Unverified'
                                }}
                            </dd>
                        </div>
                    </dl>
                    <p>
                        Account details are read-only here. Permissions are
                        managed separately.
                    </p>
                    <Link
                        v-if="accessSettingsHref"
                        :href="accessSettingsHref"
                        class="button button-outline"
                        >Open access settings</Link
                    >
                </template>
                <section
                    v-if="pendingOperation"
                    class="operation-confirm"
                    aria-labelledby="operation-confirm-title"
                >
                    <h3 id="operation-confirm-title" tabindex="-1">
                        {{ operationConfirmation.title }}
                    </h3>
                    <p>{{ operationConfirmation.message }}</p>
                    <p v-if="pendingQuote">
                        {{ formatPrice(pendingQuote.unit_price) }} per unit ·
                        {{ pendingQuote.lead_time_days }} days ·
                        {{ detailRecord.quantity }} units
                    </p>
                    <div class="workspace-form-actions">
                        <button
                            type="button"
                            class="button button-outline"
                            :disabled="busy"
                            @click="keepReviewing"
                        >
                            Keep reviewing</button
                        ><button
                            type="button"
                            class="button button-navy"
                            :disabled="busy"
                            @click="confirmOperation"
                        >
                            {{
                                busy ? 'Saving…' : operationConfirmation.button
                            }}
                        </button>
                    </div>
                </section>
            </div>
            <template #footer="{ requestClose }"
                ><div class="workspace-form-actions">
                    <button
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="requestClose"
                    >
                        Close details</button
                    ><template v-if="detailRecord && !pendingOperation"
                        ><button
                            v-for="action in operationActions(
                                detailKind,
                                detailRecord
                            ).filter((a) => a.id !== 'view')"
                            :key="action.id"
                            type="button"
                            class="button button-navy"
                            :disabled="busy"
                            @click="
                                openRecord(detailKind, detailRecord, action.id)
                            "
                        >
                            {{ action.label }}
                        </button></template
                    >
                </div></template
            >
        </AdminDrawer>
        <AdminDrawer
            :open="productOpen"
            :title="
                editingProduct
                    ? 'Edit product listing'
                    : 'Create product listing'
            "
            :busy="busy"
            :dirty="productDirty"
            description="Manage product details, pricing, inventory and media. Changes are published only when you save."
            @close="closeProduct"
        >
            <form
                id="product-editor-form"
                class="soko-form workspace-form"
                novalidate
                @submit.prevent="saveProduct"
            >
                <fieldset class="product-editor-fields" :disabled="busy">
                    <p v-if="productRecovered" class="form-note" role="status">
                        Recovered your unsaved listing draft. Reselect
                        unfinished uploads and re-enter any unpublished digital
                        delivery link before saving.
                    </p>
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
                                <option value="make_offer">
                                    Make an offer
                                </option>
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
                                <span
                                    >(optional; generated if empty)</span
                                ></label
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
                            <label for="product-bulk"
                                >Bulk unit price (USD)</label
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
                            Set the unit price for each quantity range. Ranges
                            must be in ascending order and cannot overlap. Use 0
                            as the last maximum quantity for no upper limit.
                        </p>
                        <p
                            v-if="!productForm.price_tiers.length"
                            class="form-note"
                        >
                            No volume tiers: the base price applies. For
                            wholesale quantities of 10 or more, an explicit bulk
                            unit price applies when greater than zero.
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
                                @click="
                                    productForm.price_tiers.splice(index, 1)
                                "
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
                            Give the last tier an upper limit before adding
                            another range.
                        </p>
                    </fieldset>
                    <div class="form-row">
                        <div class="soko-field">
                            <label for="product-moq"
                                >Minimum order quantity</label
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
                                Enter available units, or -1 for unlimited
                                digital inventory.
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
                        <label for="product-value"
                            >What this product helps the buyer achieve</label
                        >
                        <textarea
                            id="product-value"
                            v-model="productForm.short_desc"
                            maxlength="300"
                            rows="2"
                            placeholder="Describe the useful outcome, without unsupported claims."
                        />
                        <p class="form-note">
                            A concise value proposition for product discovery.
                            Up to 300 characters.
                        </p>
                    </div>
                    <ProductMediaStudio
                        v-if="
                            editingProduct
                                ? recordCan(
                                      'products',
                                      editingProduct.id,
                                      'media'
                                  )
                                : capabilities?.createMedia
                        "
                        v-model="productForm.image_url"
                        v-model:gallery="productForm.gallery_urls"
                        :title="productForm.title"
                    />
                    <label class="workspace-checkbox"
                        ><input
                            type="checkbox"
                            v-model="productForm.featured"
                        />Feature this listing in discovery</label
                    >
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
                            ><select
                                id="product-type"
                                v-model="productForm.type"
                            >
                                <option value="digital_software">
                                    Software
                                </option>
                                <option value="digital_asset">
                                    Digital asset
                                </option>
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
                                    invalidFields.includes(
                                        'digital_download_url'
                                    )
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
                        Prices are entered in USD. Display currency conversion
                        does not change the listing's base price.
                    </p>
                    <p v-if="editingProduct" class="form-note">
                        If inventory changes while you are editing, reopen the
                        listing before trying again.
                    </p>
                </fieldset>
            </form>
            <template #footer="{ requestClose }">
                <div class="workspace-form-actions">
                    <button
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="requestClose"
                    >
                        Cancel</button
                    ><button
                        type="submit"
                        form="product-editor-form"
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
            </template>
        </AdminDrawer>
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
        <AdminDrawer
            :open="!!quoteTarget"
            :title="
                quoteTarget && ownQuote(quoteTarget)
                    ? 'Update your quotation'
                    : 'Submit a quotation'
            "
            :busy="busy"
            :dirty="quoteDirty"
            @close="closeQuote"
            ><form
                id="quote-editor-form"
                :inert="busy || undefined"
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
            </form>
            <template #footer="{ requestClose }"
                ><div class="workspace-form-actions">
                    <button
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="requestClose"
                    >
                        Cancel</button
                    ><button
                        type="submit"
                        form="quote-editor-form"
                        class="button button-navy"
                        :disabled="busy"
                    >
                        {{ busy ? 'Sending quotation…' : 'Send quotation' }}
                    </button>
                </div></template
            ></AdminDrawer
        >
    </div>
</template>

<script setup lang="ts">
import { computed, nextTick, reactive, ref, watch } from 'vue'
import { safeProductDraft } from '@/design/productDraft'
import { Link, router, usePage } from '@inertiajs/vue3'
import SokoDialog from '@/components/SokoDialog.vue'
import AdminDrawer from '@/components/AdminDrawer.vue'
import AdminTable from '@/components/AdminTable.vue'
import AdminActionMenu from '@/components/AdminActionMenu.vue'
import SokoIcon from '@/components/SokoIcon.vue'
import ProductMediaStudio from '@/components/ProductMediaStudio.vue'
import {
    useCurrency,
    type Currency,
    type MarketRegion,
} from '@/composables/useCurrency'
import '../../../css/workspace.css'

type RecordData = Record<string, any>
const props = defineProps<{
    recordCapabilities?: Record<string, Record<string, Record<string, boolean>>>
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
        createMedia?: boolean
        createProducts?: boolean
        editProducts?: boolean
        deleteProducts?: boolean
        readProducts?: boolean
        readOrders?: boolean
        readOffers?: boolean
        readRfqs?: boolean
        manageContent?: boolean
        manageSettings?: boolean
        managePromotions?: boolean
        manageUsers?: boolean
        reviewOffers?: boolean
        quoteRfqs?: boolean
    }
}>()
const recordCan = (kind: string, id: string, action: string) =>
    props.recordCapabilities?.[kind]?.[id]?.[action] === true
const detailRecord = ref<RecordData | null>(null)
const detailKind = ref('')
const pendingOperation = ref('')
const pendingQuote = ref<RecordData | null>(null)
const detailTitle = computed(
    () =>
        detailRecord.value?.order_number ||
        detailRecord.value?.title ||
        detailRecord.value?.product_title ||
        detailRecord.value?.name ||
        'Record details'
)
const accessSettingsHref = computed(
    () =>
        ((page.props.workspaceNavigation as RecordData[]) || []).find(
            (item) => item.href === '/dashboard/settings?panel=access'
        )?.href || ''
)
const operationActions = (kind: string, record: RecordData) => [
    { id: 'view', label: 'View details' },
    ...(kind === 'orders' && recordCan(kind, record.id, 'cancel')
        ? [{ id: 'cancel', label: 'Review cancellation', danger: true }]
        : []),
    ...(kind === 'rfqs' && recordCan(kind, record.id, 'quote')
        ? [{ id: 'quote', label: 'Prepare quotation' }]
        : []),
    ...(kind === 'rfqs' && recordCan(kind, record.id, 'close')
        ? [{ id: 'close', label: 'Review closing request', danger: true }]
        : []),
    ...(kind === 'offers' && recordCan(kind, record.id, 'decide')
        ? [
              { id: 'accepted', label: 'Review acceptance' },
              { id: 'declined', label: 'Review decline', danger: true },
          ]
        : []),
]
function openRecord(kind: string, record: RecordData, action = 'view') {
    if (
        busy.value ||
        !operationActions(kind, record).some((item) => item.id === action)
    )
        return
    if (action === 'quote') {
        detailRecord.value = null
        openQuote(record)
        return
    }
    detailKind.value = kind
    detailRecord.value = record
    pendingOperation.value = ''
    pendingQuote.value = null
    formError.value = ''
    if (action !== 'view') prepareOperation(action)
}
function closeRecord() {
    if (!busy.value) {
        detailRecord.value = null
        pendingOperation.value = ''
        pendingQuote.value = null
    }
}
function prepareOperation(action: string, quote: RecordData | null = null) {
    const record = detailRecord.value
    if (busy.value || !record) return
    if (action === 'award') {
        if (
            detailKind.value !== 'rfqs' ||
            !quote ||
            !(record.quotes || []).some((q: RecordData) => q.id === quote.id) ||
            !recordCan('quotes', quote.id, 'accept')
        )
            return
    } else if (
        !operationActions(detailKind.value, record).some(
            (item) => item.id === action && item.id !== 'view'
        )
    )
        return
    pendingOperation.value = action
    pendingQuote.value = quote ? JSON.parse(JSON.stringify(quote)) : null
    void nextTick(() => {
        if (
            typeof document !== 'undefined' &&
            pendingOperation.value === action
        )
            document.getElementById('operation-confirm-title')?.focus()
    })
}
const operationConfirmation = computed(
    () =>
        ({
            cancel: {
                title: 'Cancel the entire pending order?',
                message:
                    'This cancels the entire order, including any supplier lines you cannot see in this view, and releases all reserved inventory. Payment is not refunded or collected by this action.',
                button: 'Confirm cancellation',
            },
            close: {
                title: 'Close this sourcing request?',
                message: 'This ends the open request and stops new quotations.',
                button: 'Close request',
            },
            award: {
                title: 'Award this quotation?',
                message:
                    'This accepts the selected supplier quotation and resolves the sourcing request. Review the price, lead time and terms before confirming. No payment is collected.',
                button: 'Accept quotation',
            },
            accepted: {
                title: 'Accept this buyer offer?',
                message:
                    'This records acceptance of the displayed proposal. No payment or fulfillment is performed.',
                button: 'Accept offer',
            },
            declined: {
                title: 'Decline this buyer offer?',
                message: 'This marks the displayed proposal as declined.',
                button: 'Decline offer',
            },
        })[pendingOperation.value] || { title: '', message: '', button: '' }
)
function confirmOperation() {
    const record = detailRecord.value,
        action = pendingOperation.value
    if (!record || busy.value || !action) return
    if (action === 'accepted' || action === 'declined') {
        if (detailKind.value === 'offers') decideOffer(record, action)
        return
    }
    if (action === 'award' && !pendingQuote.value?.updated_at) {
        formError.value =
            'Quotation version is missing. Refresh details and review the current terms before awarding.'
        return
    }
    let url = ''
    if (
        action === 'cancel' &&
        detailKind.value === 'orders' &&
        recordCan('orders', record.id, 'cancel')
    )
        url = `/dashboard/orders/${record.id}/cancel`
    if (
        action === 'close' &&
        detailKind.value === 'rfqs' &&
        recordCan('rfqs', record.id, 'close')
    )
        url = `/dashboard/rfqs/${record.id}/close`
    if (
        action === 'award' &&
        detailKind.value === 'rfqs' &&
        pendingQuote.value &&
        (record.quotes || []).some(
            (q: RecordData) => q.id === pendingQuote.value?.id
        ) &&
        recordCan('quotes', pendingQuote.value.id, 'accept')
    )
        url = `/dashboard/rfqs/${record.id}/quotes/${pendingQuote.value.id}/accept`
    if (!url) return
    mutate(
        url,
        action === 'award'
            ? { quote_version: String(pendingQuote.value?.updated_at) }
            : {},
        action,
        () => {
            detailRecord.value = null
            pendingOperation.value = ''
            pendingQuote.value = null
        },
        true
    )
}
function refreshRecord() {
    if (busy.value || !detailRecord.value) return
    const kind = detailKind.value,
        id = detailRecord.value.id
    busy.value = true
    router.reload({
        only: [kind, 'recordCapabilities'],
        onSuccess: (response) => {
            const latest = ((response.props[kind] as RecordData[]) || []).find(
                (record) => record.id === id
            )
            if (!latest) {
                formError.value =
                    'This record is no longer available in your workspace. Close these details and review the list.'
                return
            }
            detailRecord.value = latest
            pendingOperation.value = ''
            pendingQuote.value = null
            formError.value = ''
        },
        onError: () => {
            formError.value =
                'Could not refresh these details. Your previous review is still displayed.'
        },
        onNetworkError: () => {
            formError.value =
                'Connection lost while refreshing. Your previous review is still displayed.'
            return false
        },
        onHttpException: () => {
            formError.value =
                'Could not refresh these details. Your previous review is still displayed.'
            return false
        },
        onFinish: () => {
            busy.value = false
        },
    })
}
const productActions = (product: RecordData) => [
    ...(recordCan('products', product.id, 'update')
        ? [{ id: 'edit', label: 'Edit product' }]
        : []),
    ...(recordCan('products', product.id, 'delete')
        ? [{ id: 'delete', label: 'Delete product', danger: true }]
        : []),
]
function productAction(product: RecordData, action: string) {
    if (busy.value) return
    if (action === 'edit') openProduct(product)
    else if (action === 'delete') openDelete(product)
}
const page = usePage()
const { formatPrice } = useCurrency(props.currencies)
const isAdmin = computed(() => props.currentUser.role === 'admin')
const canManageCatalog = computed(
    () => props.capabilities?.manageCatalog === true
)
const canReviewOffers = computed(
    () => props.capabilities?.reviewOffers === true
)
const canQuote = computed(() => props.capabilities?.quoteRfqs === true)
const canManageUsers = computed(() => props.capabilities?.manageUsers === true)
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
    ...(props.capabilities?.readProducts
        ? [
              {
                  id: 'products',
                  label: 'Products',
                  icon: 'box',
                  count: productList.value.length,
              },
          ]
        : []),
    ...(props.capabilities?.readOrders
        ? [
              {
                  id: 'orders',
                  label: 'Orders',
                  icon: 'document',
                  count: orderList.value.length,
              },
          ]
        : []),
    ...(props.capabilities?.readRfqs
        ? [
              {
                  id: 'rfqs',
                  label: 'Sourcing requests',
                  icon: 'globe',
                  count: rfqList.value.length,
              },
          ]
        : []),
    ...(props.capabilities?.readOffers
        ? [
              {
                  id: 'offers',
                  label: 'Buyer offers',
                  icon: 'document',
                  count: offerList.value.length,
              },
          ]
        : []),
    ...(canManageUsers.value
        ? [
              {
                  id: 'users',
                  label: 'Users',
                  icon: 'user',
                  count: (props.users || []).length,
              },
          ]
        : []),
    ...(props.capabilities?.manageContent
        ? [{ id: 'content', label: 'Storefront content', icon: 'document' }]
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
const recordStatus = ref('all')
function keepReviewing() {
    pendingOperation.value = ''
    pendingQuote.value = null
}
function clearRecordFilters() {
    search.value = ''
    recordStatus.value = 'all'
}
const statusOptions = computed(() =>
    [
        ...new Set(
            records.value
                .map((item) => String(item.status || ''))
                .filter(Boolean)
        ),
    ].sort()
)
const pageError = ref('')
const localNotice = ref('')
watch(section, () => {
    search.value = ''
    recordStatus.value = 'all'
    closeRecord()
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
    const scoped =
        recordStatus.value === 'all'
            ? records.value
            : records.value.filter((item) => item.status === recordStatus.value)
    return query
        ? scoped.filter((item) =>
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
                  'buyer_email',
                  'destination_country',
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
        : scoped
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
    short_desc: '',
    gallery_urls: [] as string[],
    featured: false,
    digital_license_format: '',
    digital_file_size: '',
    digital_download_url: '',
    accepts_offers: false,
})
const productForm = reactive(emptyProduct())
const productSnapshot = ref('')
const productRecovered = ref(false)
const productDirty = computed(
    () =>
        productOpen.value &&
        JSON.stringify(productForm) !== productSnapshot.value
)
const productDraftKey = `soko.product-editor.${props.currentUser.id}`
const restoredProductDraft = router.restore(productDraftKey) as
    RecordData | undefined
if (restoredProductDraft?.open && restoredProductDraft.form) {
    const record = restoredProductDraft.productId
        ? (props.products || []).find(
              (p: RecordData) => p.id === restoredProductDraft.productId
          )
        : null
    const allowed = record
        ? recordCan('products', record.id, 'update')
        : !restoredProductDraft.productId && props.capabilities?.createProducts
    if (allowed) {
        Object.assign(productForm, safeProductDraft(restoredProductDraft.form))
        productForm.digital_download_url = record?.digital_download_url || ''
        editingProduct.value = record || null
        productSnapshot.value = ''
        productRecovered.value = true
        productOpen.value = true
    }
}
watch(
    [productForm, productOpen],
    () => {
        router.remember(
            productOpen.value
                ? {
                      open: true,
                      productId: editingProduct.value?.id || '',
                      form: safeProductDraft(productForm),
                  }
                : null,
            productDraftKey
        )
    },
    { deep: true }
)
function clearProductDraft() {
    productOpen.value = false
    productRecovered.value = false
    router.remember(null, productDraftKey)
}
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
    if (busy.value || productOpen.value) return
    if (
        product
            ? !recordCan('products', product.id, 'update')
            : !props.capabilities?.createProducts
    )
        return
    editingProduct.value = product || null
    Object.assign(productForm, emptyProduct())
    if (product) {
        for (const key of Object.keys(emptyProduct()))
            if (product[key] !== undefined)
                (productForm as RecordData)[key] = product[key]
        productForm.region = product.vendor_region || 'global'
        productForm.version = String(product.updated_at || '')
        productForm.gallery_urls = [...(product.gallery_urls || [])]
        productForm.price_tiers = (product.price_tiers || []).map(
            (tier: { min_qty: number; max_qty: number; price: number }) => ({
                ...tier,
            })
        )
    }
    formError.value = ''
    invalidFields.value = []
    productRecovered.value = false
    productSnapshot.value = JSON.stringify(productForm)
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
    if (!busy.value) clearProductDraft()
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
            clearProductDraft()
            localNotice.value = id
                ? 'Listing changes saved.'
                : 'Listing published.'
            void nextTick(() =>
                router.get(
                    sectionHref('products'),
                    {},
                    { preserveState: true, preserveScroll: true }
                )
            )
        },
        true
    )
}
const deleteTarget = ref<RecordData | null>(null)
const openDelete = (product: RecordData) => {
    if (!recordCan('products', product.id, 'delete')) return
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
const quoteSnapshot = ref(JSON.stringify(quoteForm))
const quoteDirty = computed(
    () => JSON.stringify(quoteForm) !== quoteSnapshot.value
)
const ownQuote = (rfq: RecordData) =>
    (rfq.quotes || []).find(
        (quote: RecordData) => quote.vendor_id === props.currentUser.id
    )
const openQuote = (rfq: RecordData) => {
    if (busy.value || !recordCan('rfqs', rfq.id, 'quote')) return
    quoteTarget.value = rfq
    const quote = ownQuote(rfq)
    Object.assign(quoteForm, {
        unit_price: quote?.unit_price ?? rfq.target_price ?? 0,
        lead_time_days: quote?.lead_time_days ?? 7,
        notes: quote?.notes ?? '',
    })
    quoteSnapshot.value = JSON.stringify(quoteForm)
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
const decideOffer = (offer: RecordData, status: 'accepted' | 'declined') => {
    if (busy.value || !recordCan('offers', offer.id, 'decide')) return
    mutate(
        `/dashboard/offers/${offer.id}/decision?section=offers`,
        { status },
        `offer-${offer.id}-${status}`,
        () => {
            localNotice.value =
                status === 'accepted' ? 'Offer accepted.' : 'Offer declined.'
            detailRecord.value = null
            pendingOperation.value = ''
        },
        true
    )
}
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

<style scoped>
.operation-status-filter {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.8125rem;
}
.operation-status-filter select {
    min-height: 44px;
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    background: var(--soko-surface);
    color: var(--soko-ink);
}
.operation-details {
    display: grid;
    gap: 1.25rem;
}
.operation-facts {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 1rem;
    margin: 0;
}
.operation-facts dt {
    font-size: 0.8rem;
    color: var(--soko-muted);
    margin-bottom: 0.4rem;
}
.operation-facts dd {
    margin: 0;
    overflow-wrap: anywhere;
}
.operation-facts small,
.operation-lines small {
    display: block;
    color: var(--soko-muted);
    font-size: 0.8125rem;
}
.operation-lines {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    gap: 0.8rem;
}
.operation-lines li {
    display: grid;
    gap: 0.5rem;
    padding: 1rem;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
}
.operation-scope-note {
    background: var(--soko-canvas);
    padding: 1rem;
    border-radius: var(--soko-radius);
    font-size: 0.875rem;
    line-height: 1.6;
}
.operation-message {
    white-space: pre-wrap;
    overflow-wrap: anywhere;
    line-height: 1.6;
}
.operation-confirm {
    border: 2px solid var(--soko-navy);
    border-radius: var(--soko-radius);
    padding: 1.25rem;
    display: grid;
    gap: 1rem;
}

.workspace-admin-list .workspace-status {
    font-size: 0.8125rem;
}
.workspace-admin-list small {
    display: block;
    font-size: 0.8125rem;
    color: var(--soko-muted);
    line-height: 1.5;
    margin-top: 0.25rem;
}
.workspace-admin-list .workspace-product {
    min-width: 0;
}
.workspace-admin-list .workspace-product strong {
    min-width: 0;
}
.admin-view-link {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: 44px;
    min-width: 44px;
    padding: 0.5rem 0.7rem;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    font-size: 0.8125rem;
    font-weight: 650;
}
.admin-view-link:hover {
    background: var(--soko-canvas);
}
.admin-order-details summary {
    min-height: 44px;
    padding: 0.7rem 0;
    cursor: pointer;
    font-weight: 650;
}
.admin-order-details ul {
    margin: 0.5rem 0;
    padding-left: 1rem;
}
.admin-order-details li {
    margin-bottom: 0.75rem;
}
.admin-order-details p {
    font-size: 0.8125rem;
    line-height: 1.5;
    color: var(--soko-muted);
}
</style>

<style scoped>
.product-editor-fields {
    display: contents;
    border: 0;
    padding: 0;
    margin: 0;
    min-width: 0;
}
</style>
