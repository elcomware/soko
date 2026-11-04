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
                    :disabled="
                        busy ||
                        productListLoading ||
                        productRefreshing ||
                        !!pendingProductDraft
                    "
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
                <form
                    v-if="section === 'products' && productPage"
                    class="product-list-toolbar"
                    @submit.prevent="applyProductFilters"
                    :aria-busy="productListLoading"
                >
                    <label class="product-query"
                        >Search products<input
                            v-model="productFilters.q"
                            type="search"
                            maxlength="120"
                            placeholder="Product name, SKU, category or merchant"
                            :disabled="busy || productRefreshing"
                    /></label>
                    <label
                        >Category<select
                            v-model="productFilters.category"
                            @change="applyProductFilters"
                            :disabled="busy || productRefreshing"
                        >
                            <option value="">All categories</option>
                            <option
                                v-for="category in productFilterCategories"
                                :key="category"
                                :value="category"
                            >
                                {{ category }}
                            </option>
                        </select></label
                    >
                    <label
                        >Sort<select
                            v-model="productFilters.sort"
                            @change="applyProductFilters"
                            :disabled="busy || productRefreshing"
                        >
                            <option value="newest">Newest first</option>
                            <option value="oldest">Oldest first</option>
                            <option value="title_asc">Name A–Z</option>
                            <option value="title_desc">Name Z–A</option>
                        </select></label
                    >
                    <label
                        >Per page<select
                            v-model.number="productFilters.per_page"
                            @change="applyProductFilters"
                            :disabled="busy || productRefreshing"
                        >
                            <option
                                v-if="
                                    ![25, 50, 100].includes(
                                        productFilters.per_page
                                    )
                                "
                                :value="productFilters.per_page"
                            >
                                {{ productFilters.per_page }}
                            </option>
                            <option :value="25">25</option>
                            <option :value="50">50</option>
                            <option :value="100">100</option>
                        </select></label
                    >
                    <button
                        type="submit"
                        class="button button-navy"
                        :disabled="busy || productRefreshing"
                    >
                        {{
                            productListLoading ? 'Loading…' : 'Apply filters'
                        }}</button
                    ><button
                        type="button"
                        class="button button-outline"
                        :disabled="busy || productRefreshing"
                        @click="resetProductFilters"
                    >
                        Reset
                    </button>
                    <p role="status" class="product-results-count">
                        {{
                            productListLoading ? 'Updating results… ' : ''
                        }}Showing {{ productPage.from }}–{{ productPage.to }} of
                        {{ productPage.filtered_total }} matching products ·
                        {{ productPage.total }} total in your authorized
                        catalog.
                    </p>
                    <p v-if="productListError" role="alert" class="form-error">
                        {{ productListError }}
                    </p>
                </form>
                <div
                    v-else-if="
                        !(section === 'orders' && orderPage) &&
                        !(section === 'users' && userPage) &&
                        !(section === 'offers' && offerPage) &&
                        !(section === 'rfqs' && rfqPage)
                    "
                    class="workspace-toolbar"
                >
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
                    <template v-if="section === 'users'">
                        <label class="operation-status-filter"
                            >Role<select v-model="accountRole">
                                <option value="all">All roles</option>
                                <option
                                    v-for="role in accountRoles"
                                    :key="role"
                                    :value="role"
                                >
                                    {{ humanize(role) }}
                                </option>
                            </select></label
                        >
                        <label class="operation-status-filter"
                            >Verification<select v-model="accountVerification">
                                <option value="all">All accounts</option>
                                <option value="verified">Verified</option>
                                <option value="unverified">Unverified</option>
                            </select></label
                        >
                    </template>
                    <button
                        v-if="hasActiveFilters"
                        type="button"
                        class="admin-view-link"
                        @click="clearRecordFilters"
                    >
                        Clear filters
                    </button>
                    <p role="status">
                        {{ visibleRecords.length }}
                        {{ visibleRecords.length === 1 ? 'record' : 'records'
                        }}<span v-if="hasActiveFilters">
                            matched out of {{ records.length }}</span
                        >
                    </p>
                </div>
                <div
                    v-if="section === 'products' && pendingProductDraft"
                    class="product-recovery-banner"
                    role="status"
                >
                    <p>
                        An unsaved product draft is outside this page. It has
                        not been discarded. Check current access to reopen it.
                    </p>
                    <button
                        type="button"
                        class="button button-outline"
                        :disabled="productRefreshing || busy"
                        @click="recoverOffpageProduct"
                    >
                        {{
                            productRefreshing
                                ? 'Checking…'
                                : 'Reopen saved draft'
                        }}
                    </button>
                    <button
                        v-if="!pendingProductDiscard"
                        type="button"
                        class="button button-outline"
                        :disabled="busy || productRefreshing"
                        @click="pendingProductDiscard = true"
                    >
                        Discard saved draft…
                    </button>
                    <div v-else>
                        <p>
                            Discard this unsaved draft? The saved product will
                            not change.
                        </p>
                        <button
                            type="button"
                            class="button button-outline"
                            :disabled="busy || productRefreshing"
                            @click="pendingProductDiscard = false"
                        >
                            Keep draft</button
                        ><button
                            type="button"
                            class="button button-outline"
                            :disabled="busy || productRefreshing"
                            @click="discardPendingProductDraft"
                        >
                            Discard unsaved draft
                        </button>
                    </div>
                    <p v-if="productRefreshMessage">
                        {{ productRefreshMessage }}
                    </p>
                </div>
                <div
                    v-if="
                        !visibleRecords.length &&
                        !(section === 'orders' && orderPage) &&
                        !(section === 'users' && userPage) &&
                        !(section === 'offers' && offerPage) &&
                        !(section === 'rfqs' && rfqPage)
                    "
                    class="workspace-panel workspace-empty"
                >
                    <SokoIcon :name="activeNavigation.icon" />
                    <h2>
                        {{
                            hasActiveFilters
                                ? 'No matching records'
                                : emptyTitle
                        }}
                    </h2>
                    <p>
                        {{
                            hasActiveFilters
                                ? 'No records match your current filters. Try another name, status or filter combination.'
                                : emptyDescription
                        }}
                    </p>
                    <button
                        v-if="hasActiveFilters"
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
                                    :disabled="
                                        busy ||
                                        productListLoading ||
                                        productRefreshing
                                    "
                                    @select="(id) => productAction(product, id)"
                                />
                            </div>
                        </td>
                    </tr>
                    <template v-if="productPage" #footer
                        ><nav
                            class="product-pagination"
                            aria-label="Product catalog pages"
                        >
                            <button
                                type="button"
                                class="button button-outline"
                                :disabled="
                                    productListLoading ||
                                    productRefreshing ||
                                    productFiltersPending ||
                                    busy ||
                                    productPage.page <= 1
                                "
                                @click="changeProductPage(productPage.page - 1)"
                            >
                                Previous</button
                            ><span aria-live="polite"
                                >Page {{ productPage.page }} of
                                {{ productPage.last_page }}</span
                            ><button
                                type="button"
                                class="button button-outline"
                                :disabled="
                                    productListLoading ||
                                    productRefreshing ||
                                    productFiltersPending ||
                                    busy ||
                                    productPage.page >= productPage.last_page
                                "
                                @click="changeProductPage(productPage.page + 1)"
                            >
                                Next
                            </button>
                        </nav></template
                    >
                </AdminTable>
                <AdminTable
                    v-else-if="section === 'orders'"
                    :caption="
                        orderVisibility === 'full'
                            ? 'Orders'
                            : 'Orders · authorized supplier lines'
                    "
                    class="workspace-admin-list"
                    :columns="orderColumns"
                    :sort="{
                        key: 'created_at',
                        direction:
                            orderFilters.sort === 'oldest' ? 'asc' : 'desc',
                    }"
                    resizable-columns
                    :viewport-height="520"
                    resizable-height
                    :empty="!visibleRecords.length"
                    :empty-text="
                        hasActiveFilters
                            ? 'No orders match these filters. Change or reset the filters to continue.'
                            : 'No orders are available in your authorized scope.'
                    "
                    :busy="busy || orderListLoading"
                    @sort="sortOrderColumn"
                >
                    <template v-if="orderPage" #toolbar
                        ><form
                            class="order-search-toolbar"
                            @submit.prevent="applyOrderFilters"
                        >
                            <label
                                >Search orders<input
                                    v-model="orderFilters.q"
                                    type="search"
                                    maxlength="120"
                                    placeholder="Order reference, buyer or status"
                                    :disabled="busy" /></label
                            ><button
                                type="submit"
                                class="button button-outline"
                                :disabled="busy"
                            >
                                Search</button
                            ><button
                                v-if="
                                    orderFilters.q ||
                                    orderFilters.status ||
                                    orderFilters.sort !== 'newest'
                                "
                                type="button"
                                class="button button-outline"
                                :disabled="busy"
                                @click="resetOrderFilters"
                            >
                                Reset filters
                            </button>
                            <p role="status">
                                {{ orderListLoading ? 'Updating… ' : ''
                                }}{{ orderPage.filtered_total }} matching ·
                                {{ orderPage.total }} authorized orders
                            </p>
                            <p
                                v-if="orderListError"
                                role="alert"
                                class="form-error"
                            >
                                {{ orderListError }}
                            </p>
                        </form></template
                    >
                    <template v-if="orderPage" #filters
                        ><div class="order-filter-fields">
                            <label
                                >Status<select
                                    v-model="orderFilters.status"
                                    :disabled="busy"
                                    @change="applyOrderFilters"
                                >
                                    <option value="">All statuses</option>
                                    <option
                                        v-for="status in [
                                            'pending',
                                            'cancelled',
                                            'completed',
                                            'shipped',
                                        ]"
                                        :key="status"
                                        :value="status"
                                    >
                                        {{ humanize(status) }}
                                    </option>
                                </select></label
                            ><label
                                >Received<select
                                    v-model="orderFilters.sort"
                                    :disabled="busy"
                                    @change="applyOrderFilters"
                                >
                                    <option value="newest">Newest first</option>
                                    <option value="oldest">Oldest first</option>
                                </select></label
                            >
                        </div></template
                    >
                    <tr v-for="order in visibleRecords" :key="order.id">
                        <th scope="row">
                            <button
                                type="button"
                                class="operation-identity"
                                :disabled="busy || orderListLoading"
                                :aria-label="`View order ${order.order_number}`"
                                @click="openRecord('orders', order)"
                            >
                                {{ order.order_number }}</button
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
                                :disabled="busy || orderListLoading"
                                @select="
                                    (id) => openRecord('orders', order, id)
                                "
                            />
                        </td>
                    </tr>
                    <template v-if="orderPage" #footer
                        ><AdminPagination
                            :page="orderPage.page"
                            :per-page="orderPage.per_page"
                            :total="orderPage.filtered_total"
                            :last-page="orderPage.last_page"
                            :from="orderPage.from"
                            :to="orderPage.to"
                            :busy="
                                busy || orderListLoading || orderFiltersPending
                            "
                            @page="changeOrderPage"
                            @per-page="changeOrderPageSize"
                    /></template>
                </AdminTable>
                <AdminTable
                    v-else-if="section === 'rfqs'"
                    caption="Sourcing requests"
                    :columns="rfqColumns"
                    :sort="{
                        key: 'created_at',
                        direction:
                            rfqFilters.sort === 'oldest' ? 'asc' : 'desc',
                    }"
                    resizable-columns
                    :viewport-height="520"
                    resizable-height
                    :empty="!visibleRecords.length"
                    empty-text="No requests match these filters."
                    :busy="rfqListLoading"
                    @sort="sortRfqColumn"
                    class="workspace-admin-list"
                >
                    <template v-if="rfqPage" #toolbar
                        ><form
                            class="order-search-toolbar"
                            @submit.prevent="applyRfqFilters"
                        >
                            <label
                                >Search requests<input
                                    v-model="rfqFilters.q"
                                    type="search"
                                    maxlength="120"
                                    placeholder="Request, category or buyer"
                                    :disabled="busy" /></label
                            ><button
                                type="submit"
                                class="button button-outline"
                                :disabled="busy || rfqListLoading"
                            >
                                Search</button
                            ><button
                                v-if="
                                    rfqFilters.q ||
                                    rfqFilters.status ||
                                    rfqFilters.sort !== 'newest'
                                "
                                type="button"
                                class="button button-outline"
                                :disabled="busy"
                                @click="resetRfqFilters"
                            >
                                Reset filters
                            </button>
                            <div
                                v-if="pendingQuoteDraft"
                                class="operation-error"
                            >
                                <p>
                                    An unsent quotation draft is saved in this
                                    browser history.
                                </p>
                                <button
                                    type="button"
                                    class="button button-outline"
                                    :disabled="busy || rfqListLoading"
                                    @click="recoverQuoteDraft"
                                >
                                    Reopen draft
                                </button>
                                <button
                                    type="button"
                                    class="button button-outline"
                                    :disabled="busy"
                                    @click="quoteDiscardPending = true"
                                >
                                    Discard draft…
                                </button>
                                <div v-if="quoteDiscardPending">
                                    <p>Discard the unsent quotation?</p>
                                    <button
                                        type="button"
                                        class="button button-outline"
                                        :disabled="busy"
                                        @click="discardQuoteDraft"
                                    >
                                        Discard quotation</button
                                    ><button
                                        type="button"
                                        class="button button-outline"
                                        :disabled="busy"
                                        @click="quoteDiscardPending = false"
                                    >
                                        Keep draft
                                    </button>
                                </div>
                                <p v-if="quoteRecoveryError" role="alert">
                                    {{ quoteRecoveryError }}
                                </p>
                            </div>
                            <p role="status">
                                {{ rfqListLoading ? 'Updating… ' : ''
                                }}{{ rfqPage.filtered_total }} matching ·
                                {{ rfqPage.total }} authorized requests
                            </p>
                            <p
                                v-if="rfqListError"
                                role="alert"
                                class="form-error"
                            >
                                {{ rfqListError }}
                            </p>
                        </form></template
                    >
                    <template v-if="rfqPage" #filters
                        ><div class="order-filter-fields">
                            <label
                                >Status<select
                                    v-model="rfqFilters.status"
                                    :disabled="busy || rfqListLoading"
                                    @change="applyRfqFilters"
                                >
                                    <option value="">All statuses</option>
                                    <option
                                        v-for="status in [
                                            'open',
                                            'closed',
                                            'awarded',
                                        ]"
                                        :key="status"
                                        :value="status"
                                    >
                                        {{ humanize(status) }}
                                    </option>
                                </select></label
                            ><label
                                >Received<select
                                    v-model="rfqFilters.sort"
                                    :disabled="busy || rfqListLoading"
                                    @change="applyRfqFilters"
                                >
                                    <option value="newest">Newest first</option>
                                    <option value="oldest">Oldest first</option>
                                </select></label
                            >
                        </div></template
                    >
                    <tr v-for="rfq in visibleRecords" :key="rfq.id">
                        <th scope="row">
                            <button
                                type="button"
                                class="operation-identity"
                                :aria-label="`View request ${rfq.title}`"
                                :disabled="busy || rfqListLoading"
                                @click="openRecord('rfqs', rfq)"
                            >
                                {{ rfq.title }}</button
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
                                :disabled="busy || rfqListLoading"
                                @select="(id) => openRecord('rfqs', rfq, id)"
                            />
                        </td>
                    </tr>
                    <template v-if="rfqPage" #footer
                        ><AdminPagination
                            :page="rfqPage.page"
                            :per-page="rfqPage.per_page"
                            :total="rfqPage.filtered_total"
                            :last-page="rfqPage.last_page"
                            :from="rfqPage.from"
                            :to="rfqPage.to"
                            :busy="busy || rfqListLoading || rfqFiltersPending"
                            @page="changeRfqPage"
                            @per-page="changeRfqPageSize"
                    /></template>
                </AdminTable>
                <AdminTable
                    v-else-if="section === 'offers'"
                    caption="Buyer offers"
                    :columns="offerColumns"
                    :sort="{
                        key: 'created_at',
                        direction:
                            offerFilters.sort === 'oldest' ? 'asc' : 'desc',
                    }"
                    resizable-columns
                    :viewport-height="520"
                    resizable-height
                    :empty="!visibleRecords.length"
                    empty-text="No offers match these filters."
                    :busy="offerListLoading"
                    @sort="sortOfferColumn"
                    class="workspace-admin-list"
                >
                    <template v-if="offerPage" #toolbar
                        ><form
                            class="order-search-toolbar"
                            @submit.prevent="applyOfferFilters"
                        >
                            <label
                                >Search offers<input
                                    v-model="offerFilters.q"
                                    type="search"
                                    maxlength="120"
                                    placeholder="Product, buyer or status"
                                    :disabled="busy" /></label
                            ><button
                                type="submit"
                                class="button button-outline"
                                :disabled="busy"
                            >
                                Search</button
                            ><button
                                v-if="
                                    offerFilters.q ||
                                    offerFilters.status ||
                                    offerFilters.sort !== 'newest'
                                "
                                type="button"
                                class="button button-outline"
                                :disabled="busy"
                                @click="resetOfferFilters"
                            >
                                Reset filters
                            </button>
                            <p role="status">
                                {{ offerListLoading ? 'Updating… ' : ''
                                }}{{ offerPage.filtered_total }} matching ·
                                {{ offerPage.total }} authorized offers
                            </p>
                            <p
                                v-if="offerListError"
                                role="alert"
                                class="form-error"
                            >
                                {{ offerListError }}
                            </p>
                        </form></template
                    >
                    <template v-if="offerPage" #filters
                        ><div class="order-filter-fields">
                            <label
                                >Status<select
                                    v-model="offerFilters.status"
                                    :disabled="busy"
                                    @change="applyOfferFilters"
                                >
                                    <option value="">All statuses</option>
                                    <option
                                        v-for="status in [
                                            'pending',
                                            'accepted',
                                            'declined',
                                            'countered',
                                        ]"
                                        :key="status"
                                        :value="status"
                                    >
                                        {{ humanize(status) }}
                                    </option>
                                </select></label
                            ><label
                                >Received<select
                                    v-model="offerFilters.sort"
                                    :disabled="busy"
                                    @change="applyOfferFilters"
                                >
                                    <option value="newest">Newest first</option>
                                    <option value="oldest">Oldest first</option>
                                </select></label
                            >
                        </div></template
                    >
                    <tr v-for="offer in visibleRecords" :key="offer.id">
                        <th scope="row">
                            <button
                                type="button"
                                class="operation-identity"
                                :aria-label="`View offer ${offer.product_title}`"
                                :disabled="busy || offerListLoading"
                                @click="openRecord('offers', offer)"
                            >
                                {{ offer.product_title }}</button
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
                                :disabled="busy || offerListLoading"
                                @select="
                                    (id) => openRecord('offers', offer, id)
                                "
                            />
                        </td>
                    </tr>
                    <template v-if="offerPage" #footer
                        ><AdminPagination
                            :page="offerPage.page"
                            :per-page="offerPage.per_page"
                            :total="offerPage.filtered_total"
                            :last-page="offerPage.last_page"
                            :from="offerPage.from"
                            :to="offerPage.to"
                            :busy="
                                busy || offerListLoading || offerFiltersPending
                            "
                            @page="changeOfferPage"
                            @per-page="changeOfferPageSize"
                    /></template>
                </AdminTable>
                <AdminTable
                    v-else-if="section === 'users'"
                    caption="Marketplace accounts"
                    :columns="userColumns"
                    :sort="{
                        key: 'created_at',
                        direction:
                            userFilters.sort === 'oldest' ? 'asc' : 'desc',
                    }"
                    resizable-columns
                    :viewport-height="520"
                    resizable-height
                    :empty="!visibleRecords.length"
                    empty-text="No accounts match these filters."
                    @sort="sortUserColumn"
                    :busy="busy || userListLoading"
                    class="workspace-admin-list"
                >
                    <template v-if="userPage" #toolbar
                        ><form
                            class="order-search-toolbar"
                            @submit.prevent="applyUserFilters"
                        >
                            <label
                                >Search accounts<input
                                    v-model="userFilters.q"
                                    type="search"
                                    maxlength="120"
                                    placeholder="Name, email, company or country"
                                    :disabled="busy" /></label
                            ><button
                                type="submit"
                                class="button button-outline"
                                :disabled="busy"
                            >
                                Search</button
                            ><button
                                v-if="
                                    userFilters.q ||
                                    userFilters.role ||
                                    userFilters.verification ||
                                    userFilters.sort !== 'newest'
                                "
                                type="button"
                                class="button button-outline"
                                :disabled="busy"
                                @click="resetUserFilters"
                            >
                                Reset filters
                            </button>
                            <p role="status">
                                {{ userListLoading ? 'Updating… ' : ''
                                }}{{ userPage.filtered_total }} matching ·
                                {{ userPage.total }} authorized accounts
                            </p>
                            <p
                                v-if="userListError"
                                role="alert"
                                class="form-error"
                            >
                                {{ userListError }}
                            </p>
                        </form></template
                    >
                    <template v-if="userPage" #filters
                        ><div class="order-filter-fields">
                            <label
                                >Role<select
                                    v-model="userFilters.role"
                                    :disabled="busy"
                                    @change="applyUserFilters"
                                >
                                    <option value="">All roles</option>
                                    <option
                                        v-for="role in [
                                            'admin',
                                            'merchant',
                                            'buyer',
                                        ]"
                                        :key="role"
                                        :value="role"
                                    >
                                        {{ humanize(role) }}
                                    </option>
                                </select></label
                            ><label
                                >Joined<select
                                    v-model="userFilters.sort"
                                    :disabled="busy"
                                    @change="applyUserFilters"
                                >
                                    <option value="newest">Newest first</option>
                                    <option value="oldest">Oldest first</option>
                                </select></label
                            >
                            <label
                                >Verification<select
                                    v-model="userFilters.verification"
                                    :disabled="busy"
                                    @change="applyUserFilters"
                                >
                                    <option value="">All accounts</option>
                                    <option value="verified">Verified</option>
                                    <option value="unverified">
                                        Unverified
                                    </option>
                                </select></label
                            >
                        </div></template
                    >
                    <tr v-for="user in visibleRecords" :key="user.id">
                        <th scope="row">
                            <button
                                type="button"
                                class="operation-identity"
                                :aria-label="`View account ${user.name}`"
                                :disabled="busy || userListLoading"
                                @click="openRecord('users', user)"
                            >
                                {{ user.name }}</button
                            ><small>{{ user.email }}</small
                            ><small>{{ dateLabel(user.created_at) }}</small>
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
                                :disabled="busy || userListLoading"
                                :aria-label="`View account ${user.name}`"
                                @click="openRecord('users', user)"
                            >
                                View account
                            </button>
                        </td>
                    </tr>
                    <template v-if="userPage" #footer
                        ><AdminPagination
                            :page="userPage.page"
                            :per-page="userPage.per_page"
                            :total="userPage.filtered_total"
                            :last-page="userPage.last_page"
                            :from="userPage.from"
                            :to="userPage.to"
                            :busy="
                                busy || userListLoading || userFiltersPending
                            "
                            @page="changeUserPage"
                            @per-page="changeUserPageSize"
                    /></template>
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
                        {{
                            selectedOrderVisibility === 'full'
                                ? 'All order lines are shown.'
                                : 'Only your authorized supplier lines and their subtotal are shown.'
                        }}
                    </p>
                    <button
                        v-if="orderPage"
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="refreshSelectedOrder"
                    >
                        Refresh order details
                    </button>
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
                            <dt>
                                {{
                                    selectedOrderVisibility === 'full'
                                        ? 'Order total'
                                        : 'Visible subtotal'
                                }}
                            </dt>
                            <dd>
                                {{ formatPrice(detailRecord.total_amount) }}
                            </dd>
                        </div>
                    </dl>
                    <h3>
                        {{
                            selectedOrderVisibility === 'full'
                                ? 'Order lines'
                                : 'Visible order lines'
                        }}
                    </h3>
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
                    <button
                        v-if="rfqPage"
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="refreshSelectedRfq"
                    >
                        Refresh request details
                    </button>
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
                        {{
                            selectedRfqVisibility === 'full'
                                ? 'All quotations for this request are shown.'
                                : 'Only your supplier quotations are shown. Competing quotations are not included.'
                        }}
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
                    <button
                        v-if="offerPage"
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="refreshSelectedOffer"
                    >
                        Refresh offer details
                    </button>
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
                    <button
                        v-if="userPage"
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="refreshSelectedUser"
                    >
                        Refresh account details
                    </button>
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
                            :disabled="
                                busy ||
                                (detailKind === 'orders' &&
                                    pendingOperation === 'cancel' &&
                                    !recordCan(
                                        'orders',
                                        detailRecord.id,
                                        'cancel'
                                    ))
                            "
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
            :busy="busy || productRefreshing"
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
                <fieldset
                    class="product-editor-fields"
                    :disabled="busy || productRefreshing"
                >
                    <section
                        v-if="editingProduct"
                        class="product-version-review"
                    >
                        <button
                            type="button"
                            class="button button-outline"
                            :disabled="busy || productRefreshing"
                            @click="checkSavedProduct"
                        >
                            {{
                                productRefreshing
                                    ? 'Checking saved version…'
                                    : 'Check current saved version'
                            }}
                        </button>
                        <p v-if="productRefreshMessage" role="status">
                            {{ productRefreshMessage }}
                        </p>
                        <div v-if="freshProduct">
                            <p>
                                Latest saved version: {{ freshProduct.title }} ·
                                {{ dateLabel(freshProduct.updated_at) }}. Your
                                draft is still unchanged.
                            </p>
                            <button
                                type="button"
                                class="button button-outline"
                                :disabled="busy || productRefreshing"
                                @click="replaceWithSavedProduct"
                            >
                                Discard this draft and load saved version
                            </button>
                        </div>
                    </section>
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
                <fieldset :disabled="busy" class="product-editor-fields">
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
                                :aria-invalid="
                                    invalidFields.includes('unit_price')
                                "
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
                </fieldset>
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
import { computed, nextTick, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { safeProductDraft } from '@/design/productDraft'
import { Link, router, usePage } from '@inertiajs/vue3'
import SokoDialog from '@/components/SokoDialog.vue'
import AdminDrawer from '@/components/AdminDrawer.vue'
import AdminTable from '@/components/AdminTable.vue'
import AdminPagination from '@/components/AdminPagination.vue'
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
    productPage?: {
        page: number
        per_page: number
        total: number
        filtered_total: number
        last_page: number
        from: number
        to: number
        filters: { q: string; category: string; sort: string }
    }
    productCategories?: string[]
    selectedProduct?: RecordData | null
    selectedProductStatus?: 'none' | 'available' | 'unavailable'
    orders?: RecordData[] | null
    orderPage?: {
        page: number
        per_page: number
        total: number
        filtered_total: number
        last_page: number
        from: number
        to: number
        filters: { q: string; status: string; sort: string }
    } | null
    selectedOrder?: RecordData | null
    selectedOrderStatus?: 'none' | 'available' | 'unavailable'
    orderVisibility?: 'full' | 'supplier' | 'none'
    userPage?: {
        page: number
        per_page: number
        total: number
        filtered_total: number
        last_page: number
        from: number
        to: number
        filters: { q: string; role: string; verification: string; sort: string }
    } | null
    selectedUser?: RecordData | null
    selectedUserStatus?: 'none' | 'available' | 'unavailable'
    offerPage?: {
        page: number
        per_page: number
        total: number
        filtered_total: number
        last_page: number
        from: number
        to: number
        filters: { q: string; status: string; sort: string }
    } | null
    selectedOffer?: RecordData | null
    selectedOfferStatus?: 'none' | 'available' | 'unavailable'
    offerVisibility?: 'full' | 'supplier' | 'none'
    rfqPage?: {
        page: number
        per_page: number
        total: number
        filtered_total: number
        last_page: number
        from: number
        to: number
        filters: { q: string; status: string; sort: string }
    } | null
    selectedRfq?: RecordData | null
    selectedRfqStatus?: 'none' | 'available' | 'unavailable'
    rfqVisibility?: 'full' | 'supplier' | 'none'
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
const refreshedProductAccess = ref<{
    id: string
    permissions: Record<string, boolean>
} | null>(null)
const refreshedOrderAccess = ref<{
    id: string
    permissions: Record<string, boolean>
} | null>(null)
const selectedOrderVisibility = ref<'full' | 'supplier'>('supplier')
const refreshedRfqAccess = ref<{
    id: string
    permissions: Record<string, boolean>
    quotes: Record<string, Record<string, boolean>>
} | null>(null)
const selectedRfqVisibility = ref<'full' | 'supplier'>('supplier')
const refreshedOfferAccess = ref<{
    id: string
    permissions: Record<string, boolean>
} | null>(null)
const selectedOfferVisibility = ref<'full' | 'supplier'>('supplier')
const recordCan = (kind: string, id: string, action: string) =>
    kind === 'offers' && refreshedOfferAccess.value?.id === id
        ? refreshedOfferAccess.value.permissions[action] === true
        : kind === 'rfqs' && refreshedRfqAccess.value?.id === id
          ? refreshedRfqAccess.value.permissions[action] === true
          : kind === 'quotes' &&
              refreshedRfqAccess.value &&
              detailKind.value === 'rfqs' &&
              detailRecord.value?.id === refreshedRfqAccess.value.id
            ? refreshedRfqAccess.value.quotes[id]?.[action] === true
            : kind === 'orders' && refreshedOrderAccess.value?.id === id
              ? refreshedOrderAccess.value.permissions[action] === true
              : kind === 'products' && refreshedProductAccess.value?.id === id
                ? refreshedProductAccess.value.permissions[action] === true
                : props.recordCapabilities?.[kind]?.[id]?.[action] === true
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
    { id: 'view', label: 'View details', icon: 'eye' },
    ...(kind === 'orders' && recordCan(kind, record.id, 'cancel')
        ? [
              {
                  id: 'cancel',
                  label: 'Cancel order…',
                  danger: true,
                  icon: 'close',
              },
          ]
        : []),
    ...(kind === 'rfqs' && recordCan(kind, record.id, 'quote')
        ? [{ id: 'quote', label: 'Prepare quote…', icon: 'document' }]
        : []),
    ...(kind === 'rfqs' && recordCan(kind, record.id, 'close')
        ? [
              {
                  id: 'close',
                  label: 'Close request…',
                  danger: true,
                  icon: 'close',
              },
          ]
        : []),
    ...(kind === 'offers' && recordCan(kind, record.id, 'decide')
        ? [
              { id: 'accepted', label: 'Accept offer…', icon: 'check' },
              {
                  id: 'declined',
                  label: 'Decline offer…',
                  danger: true,
                  icon: 'close',
              },
          ]
        : []),
]
function openRecord(kind: string, record: RecordData, action = 'view') {
    if (
        quoteTarget.value ||
        (kind === 'rfqs' && action === 'quote' && pendingQuoteDraft.value)
    )
        return
    if (kind === 'orders') {
        if (busy.value || orderListLoading.value) return
        if (orderSearchTimer) clearTimeout(orderSearchTimer)
        orderQueuedSearch = false
        if (
            detailKind.value !== 'orders' ||
            detailRecord.value?.id !== record.id
        ) {
            orderReturnUrl.value = orderListUrl()
            selectedOrderVisibility.value =
                props.orderVisibility === 'full' ? 'full' : 'supplier'
            refreshedOrderAccess.value = null
        }
    }
    if (kind === 'users') {
        if (busy.value || userListLoading.value) return
        if (userSearchTimer) clearTimeout(userSearchTimer)
        userQueuedSearch = false
        if (
            detailKind.value !== 'users' ||
            detailRecord.value?.id !== record.id
        ) {
            userReturnUrl.value = userListUrl()
        }
    }
    if (kind === 'offers') {
        if (busy.value || offerListLoading.value) return
        if (offerSearchTimer) clearTimeout(offerSearchTimer)
        offerQueuedSearch = false
        if (
            detailKind.value !== 'offers' ||
            detailRecord.value?.id !== record.id
        ) {
            offerReturnUrl.value = offerListUrl()
            selectedOfferVisibility.value =
                props.offerVisibility === 'full' ? 'full' : 'supplier'
            refreshedOfferAccess.value = null
        }
    }
    if (kind === 'rfqs') {
        if (busy.value || rfqListLoading.value) return
        if (rfqSearchTimer) clearTimeout(rfqSearchTimer)
        rfqQueuedSearch = false
        if (
            detailKind.value !== 'rfqs' ||
            detailRecord.value?.id !== record.id
        ) {
            rfqReturnUrl.value = rfqListUrl()
            selectedRfqVisibility.value =
                props.rfqVisibility === 'full' ? 'full' : 'supplier'
            refreshedRfqAccess.value = null
        }
    }
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
    const returnOrderUrl =
        detailKind.value === 'orders'
            ? orderReturnUrl.value || orderListUrl()
            : detailKind.value === 'rfqs'
              ? rfqReturnUrl.value || rfqListUrl()
              : ''
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
            if (returnOrderUrl)
                void nextTick(() =>
                    router.get(
                        returnOrderUrl,
                        {},
                        { preserveState: true, preserveScroll: true }
                    )
                )
        },
        true
    )
}
function refreshRecord() {
    if (detailKind.value === 'rfqs' && props.rfqPage) {
        void refreshSelectedRfq()
        return
    }
    if (detailKind.value === 'offers' && props.offerPage) {
        void refreshSelectedOffer()
        return
    }
    if (detailKind.value === 'users' && props.userPage) {
        void refreshSelectedUser()
        return
    }
    if (detailKind.value === 'orders' && props.orderPage) {
        void refreshSelectedOrder()
        return
    }
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
                  count: props.productPage?.total ?? productList.value.length,
              },
          ]
        : []),
    ...(props.capabilities?.readOrders
        ? [
              {
                  id: 'orders',
                  label: 'Orders',
                  icon: 'document',
                  count: props.orderPage?.total ?? orderList.value.length,
              },
          ]
        : []),
    ...(props.capabilities?.readRfqs
        ? [
              {
                  id: 'rfqs',
                  label: 'Sourcing requests',
                  icon: 'globe',
                  count: props.rfqPage?.total ?? rfqList.value.length,
              },
          ]
        : []),
    ...(props.capabilities?.readOffers
        ? [
              {
                  id: 'offers',
                  label: 'Buyer offers',
                  icon: 'document',
                  count: props.offerPage?.total ?? offerList.value.length,
              },
          ]
        : []),
    ...(canManageUsers.value
        ? [
              {
                  id: 'users',
                  label: 'Accounts',
                  icon: 'user',
                  count: props.userPage?.total ?? (props.users || []).length,
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
    products: 'Products, pricing and inventory.',
    orders: 'Orders and fulfillment details.',
    rfqs: 'Buyer requests and supplier quotes.',
    offers: 'Buyer proposals for your listings.',
    users: 'Read-only account directory.',
    content: 'Update the words buyers see in the marketplace.',
}
const sectionDescription = computed(() => descriptions[section.value])
const sectionHref = (id: string) => `/dashboard?section=${id}`
const orderFilters = reactive({
    q: props.orderPage?.filters.q || '',
    status: props.orderPage?.filters.status || '',
    sort: props.orderPage?.filters.sort || 'newest',
    per_page: props.orderPage?.per_page || 25,
})
const orderListLoading = ref(false)
const orderListError = ref('')
const orderReturnUrl = ref('')
let orderSearchTimer: ReturnType<typeof setTimeout> | undefined
let orderQueuedSearch = false
let orderSubmittedFilters = ''
let orderDisposed = false
let orderRefreshController: AbortController | undefined
const orderSignature = () => JSON.stringify(orderFilters)
const orderFiltersPending = computed(() =>
    props.orderPage
        ? orderFilters.q !== props.orderPage.filters.q ||
          orderFilters.status !== props.orderPage.filters.status ||
          orderFilters.sort !== props.orderPage.filters.sort ||
          orderFilters.per_page !== props.orderPage.per_page
        : false
)
const orderColumns = computed(() => [
    {
        key: 'created_at',
        label: 'Order / received',
        sortable: true,
        width: 240,
        minWidth: 180,
    },
    { key: 'buyer', label: 'Buyer', width: 190, minWidth: 140 },
    {
        key: 'items',
        label: props.orderVisibility === 'full' ? 'Lines' : 'Visible lines',
        width: 110,
        minWidth: 90,
    },
    { key: 'status', label: 'Status', width: 140, minWidth: 110 },
    {
        key: 'total',
        label:
            props.orderVisibility === 'full' ? 'Order total' : 'Visible total',
        width: 150,
        minWidth: 120,
    },
    { key: 'actions', label: 'Actions', width: 100, minWidth: 90 },
])
function orderListUrl(
    pageNumber = props.orderPage?.page || 1,
    filters = props.orderPage?.filters || orderFilters,
    perPage = props.orderPage?.per_page || orderFilters.per_page
) {
    const query = new URLSearchParams({
        section: 'orders',
        sort: filters.sort || 'newest',
        page: String(pageNumber),
        per_page: String(perPage),
    })
    if (filters.q) query.set('q', filters.q)
    if (filters.status) query.set('status', filters.status)
    return '/dashboard?' + query.toString()
}
watch(
    () => props.orderPage,
    (value) => {
        if (
            value &&
            (!orderListLoading.value ||
                orderSignature() === orderSubmittedFilters)
        )
            Object.assign(orderFilters, value.filters, {
                per_page: value.per_page,
            })
    },
    { deep: true, flush: 'sync' }
)
function visitOrderList(url: string) {
    if (
        orderDisposed ||
        section.value !== 'orders' ||
        busy.value ||
        detailRecord.value
    )
        return
    if (orderListLoading.value) {
        orderQueuedSearch = true
        return
    }
    if (orderSearchTimer) clearTimeout(orderSearchTimer)
    orderQueuedSearch = false
    orderSubmittedFilters = orderSignature()
    orderListLoading.value = true
    orderListError.value = ''
    router.get(
        url,
        {},
        {
            preserveState: true,
            preserveScroll: true,
            onError: () => {
                orderListError.value =
                    'Orders could not be loaded. Previous results are still shown.'
            },
            onNetworkError: () => {
                orderListError.value =
                    'Connection lost. Previous order results are still shown. Try again.'
                return false
            },
            onHttpException: () => {
                orderListError.value =
                    'Orders could not be loaded. Previous results are still shown.'
                return false
            },
            onFinish: () => {
                orderListLoading.value = false
                if (orderQueuedSearch) {
                    orderQueuedSearch = false
                    void nextTick(applyOrderFilters)
                }
            },
        }
    )
}
function applyOrderFilters() {
    visitOrderList(orderListUrl(1, orderFilters, orderFilters.per_page))
}
watch(
    () => orderFilters.q,
    (value) => {
        if (orderSearchTimer) clearTimeout(orderSearchTimer)
        if (section.value !== 'orders' || !props.orderPage) return
        if (orderListLoading.value) {
            if (orderSignature() === orderSubmittedFilters) return
            orderQueuedSearch = true
        } else if (value === props.orderPage.filters.q) return
        orderSearchTimer = setTimeout(applyOrderFilters, value ? 350 : 0)
    }
)
function resetOrderFilters() {
    if (busy.value || detailRecord.value) return
    if (orderSearchTimer) clearTimeout(orderSearchTimer)
    Object.assign(orderFilters, {
        q: '',
        status: '',
        sort: 'newest',
        per_page: 25,
    })
    applyOrderFilters()
}
function changeOrderPage(number: number) {
    if (
        !props.orderPage ||
        orderFiltersPending.value ||
        number < 1 ||
        number > props.orderPage.last_page
    )
        return
    visitOrderList(orderListUrl(number))
}
function changeOrderPageSize(number: number) {
    if (![25, 50, 100].includes(number) || busy.value || detailRecord.value)
        return
    orderFilters.per_page = number
    applyOrderFilters()
}
function sortOrderColumn(sort: { key: string; direction: 'asc' | 'desc' }) {
    if (sort.key !== 'created_at' || busy.value || detailRecord.value) return
    orderFilters.sort = sort.direction === 'asc' ? 'oldest' : 'newest'
    applyOrderFilters()
}
async function refreshSelectedOrder() {
    if (busy.value || !detailRecord.value || detailKind.value !== 'orders')
        return
    const id = detailRecord.value.id
    busy.value = true
    orderRefreshController = new AbortController()
    try {
        const response = await fetch(
            orderListUrl() + '&selected_order=' + encodeURIComponent(id),
            {
                headers: {
                    'X-Inertia': 'true',
                    'X-Requested-With': 'XMLHttpRequest',
                    Accept: 'application/json',
                },
                signal: orderRefreshController.signal,
                credentials: 'same-origin',
                cache: 'no-store',
            }
        )
        if (!response.ok) {
            if (response.status === 403 || response.status === 404)
                refreshedOrderAccess.value = { id, permissions: {} }
            throw new Error(
                'Could not refresh this order. The previous review is still displayed.'
            )
        }
        const data = await response.json()
        if (
            orderDisposed ||
            section.value !== 'orders' ||
            detailRecord.value?.id !== id
        )
            return
        const record = data.props?.selectedOrder
        refreshedOrderAccess.value = {
            id,
            permissions: data.props?.recordCapabilities?.orders?.[id] || {},
        }
        if (
            data.props?.selectedOrderStatus !== 'available' ||
            !record ||
            record.id !== id
        ) {
            refreshedOrderAccess.value = { id, permissions: {} }
            throw new Error(
                'This order is unavailable or no longer accessible. No cancellation can be submitted from this review.'
            )
        }
        detailRecord.value = record
        selectedOrderVisibility.value =
            data.props.orderVisibility === 'full' ? 'full' : 'supplier'
        pendingOperation.value = ''
        pendingQuote.value = null
        formError.value = ''
    } catch (error) {
        if (
            orderDisposed ||
            section.value !== 'orders' ||
            detailKind.value !== 'orders' ||
            detailRecord.value?.id !== id
        )
            return
        formError.value =
            error instanceof Error
                ? error.message
                : 'Could not refresh this order. The previous review is still displayed.'
    } finally {
        busy.value = false
    }
}
onBeforeUnmount(() => {
    orderDisposed = true
    orderRefreshController?.abort()
    if (orderSearchTimer) clearTimeout(orderSearchTimer)
})
const userFilters = reactive({
    q: props.userPage?.filters.q || '',
    role: props.userPage?.filters.role || '',
    verification: props.userPage?.filters.verification || '',
    sort: props.userPage?.filters.sort || 'newest',
    per_page: props.userPage?.per_page || 25,
})
const userListLoading = ref(false)
const userListError = ref('')
const userReturnUrl = ref('')
let userSearchTimer: ReturnType<typeof setTimeout> | undefined
let userQueuedSearch = false
let userSubmittedFilters = ''
let userDisposed = false
let userRefreshController: AbortController | undefined
const userSignature = () => JSON.stringify(userFilters)
const userFiltersPending = computed(() =>
    props.userPage
        ? userFilters.q !== props.userPage.filters.q ||
          userFilters.role !== props.userPage.filters.role ||
          userFilters.verification !== props.userPage.filters.verification ||
          userFilters.sort !== props.userPage.filters.sort ||
          userFilters.per_page !== props.userPage.per_page
        : false
)
const userColumns = computed(() => [
    {
        key: 'created_at',
        label: 'Account / joined',
        sortable: true,
        width: 250,
        minWidth: 180,
    },
    { key: 'role', label: 'Role', width: 130, minWidth: 110 },
    { key: 'company', label: 'Company', width: 190, minWidth: 140 },
    { key: 'country', label: 'Country', width: 150, minWidth: 120 },
    { key: 'verification', label: 'Verification', width: 140, minWidth: 110 },
    { key: 'actions', label: 'Actions', width: 120, minWidth: 100 },
])
function userListUrl(
    pageNumber = props.userPage?.page || 1,
    filters = props.userPage?.filters || userFilters,
    perPage = props.userPage?.per_page || userFilters.per_page
) {
    const query = new URLSearchParams({
        section: 'users',
        sort: filters.sort || 'newest',
        page: String(pageNumber),
        per_page: String(perPage),
    })
    if (filters.q) query.set('q', filters.q)
    if (filters.role) query.set('role', filters.role)
    if (filters.verification) query.set('verification', filters.verification)
    return '/dashboard?' + query.toString()
}
watch(
    () => props.userPage,
    (value) => {
        if (
            value &&
            (!userListLoading.value || userSignature() === userSubmittedFilters)
        )
            Object.assign(userFilters, value.filters, {
                per_page: value.per_page,
            })
    },
    { deep: true, flush: 'sync' }
)
function visitUserList(url: string) {
    if (
        userDisposed ||
        section.value !== 'users' ||
        busy.value ||
        detailRecord.value
    )
        return
    if (userListLoading.value) {
        userQueuedSearch = true
        return
    }
    if (userSearchTimer) clearTimeout(userSearchTimer)
    userQueuedSearch = false
    userSubmittedFilters = userSignature()
    userListLoading.value = true
    userListError.value = ''
    router.get(
        url,
        {},
        {
            preserveState: true,
            preserveScroll: true,
            onError: () => {
                userListError.value =
                    'Accounts could not be loaded. Previous results are still shown.'
            },
            onNetworkError: () => {
                userListError.value =
                    'Connection lost. Previous account results are still shown. Try again.'
                return false
            },
            onHttpException: () => {
                userListError.value =
                    'Accounts could not be loaded. Previous results are still shown.'
                return false
            },
            onFinish: () => {
                userListLoading.value = false
                if (userQueuedSearch) {
                    userQueuedSearch = false
                    void nextTick(applyUserFilters)
                }
            },
        }
    )
}
function applyUserFilters() {
    visitUserList(userListUrl(1, userFilters, userFilters.per_page))
}
watch(
    () => userFilters.q,
    (value) => {
        if (userSearchTimer) clearTimeout(userSearchTimer)
        if (section.value !== 'users' || !props.userPage) return
        if (userListLoading.value) {
            if (userSignature() === userSubmittedFilters) return
            userQueuedSearch = true
        } else if (value === props.userPage.filters.q) return
        userSearchTimer = setTimeout(applyUserFilters, value ? 350 : 0)
    }
)
function resetUserFilters() {
    if (busy.value || detailRecord.value) return
    if (userSearchTimer) clearTimeout(userSearchTimer)
    Object.assign(userFilters, {
        q: '',
        role: '',
        verification: '',
        sort: 'newest',
        per_page: 25,
    })
    applyUserFilters()
}
function changeUserPage(number: number) {
    if (
        !props.userPage ||
        userFiltersPending.value ||
        number < 1 ||
        number > props.userPage.last_page
    )
        return
    visitUserList(userListUrl(number))
}
function changeUserPageSize(number: number) {
    if (![25, 50, 100].includes(number) || busy.value || detailRecord.value)
        return
    userFilters.per_page = number
    applyUserFilters()
}
function sortUserColumn(sort: { key: string; direction: 'asc' | 'desc' }) {
    if (sort.key !== 'created_at' || busy.value || detailRecord.value) return
    userFilters.sort = sort.direction === 'asc' ? 'oldest' : 'newest'
    applyUserFilters()
}
async function refreshSelectedUser() {
    if (busy.value || !detailRecord.value || detailKind.value !== 'users')
        return
    const id = detailRecord.value.id
    busy.value = true
    userRefreshController = new AbortController()
    try {
        const response = await fetch(
            userListUrl() + '&selected_user=' + encodeURIComponent(id),
            {
                headers: {
                    'X-Inertia': 'true',
                    'X-Requested-With': 'XMLHttpRequest',
                    Accept: 'application/json',
                },
                signal: userRefreshController.signal,
                credentials: 'same-origin',
                cache: 'no-store',
            }
        )
        if (!response.ok) {
            throw new Error(
                'Could not refresh this account. The previous review is still displayed.'
            )
        }
        const data = await response.json()
        if (
            userDisposed ||
            section.value !== 'users' ||
            detailRecord.value?.id !== id
        )
            return
        const record = data.props?.selectedUser
        if (
            data.props?.selectedUserStatus !== 'available' ||
            !record ||
            record.id !== id
        ) {
            throw new Error(
                'This account is unavailable or no longer accessible. Close these details and review the directory.'
            )
        }
        detailRecord.value = record
        pendingOperation.value = ''
        pendingQuote.value = null
        formError.value = ''
    } catch (error) {
        if (
            userDisposed ||
            section.value !== 'users' ||
            detailKind.value !== 'users' ||
            detailRecord.value?.id !== id
        )
            return
        formError.value =
            error instanceof Error
                ? error.message
                : 'Could not refresh this account. The previous review is still displayed.'
    } finally {
        busy.value = false
    }
}
onBeforeUnmount(() => {
    userDisposed = true
    userRefreshController?.abort()
    if (userSearchTimer) clearTimeout(userSearchTimer)
})
const offerFilters = reactive({
    q: props.offerPage?.filters.q || '',
    status: props.offerPage?.filters.status || '',
    sort: props.offerPage?.filters.sort || 'newest',
    per_page: props.offerPage?.per_page || 25,
})
const offerListLoading = ref(false)
const offerListError = ref('')
const offerReturnUrl = ref('')
let offerSearchTimer: ReturnType<typeof setTimeout> | undefined
let offerQueuedSearch = false
let offerSubmittedFilters = ''
let offerDisposed = false
let offerRefreshController: AbortController | undefined
const offerSignature = () => JSON.stringify(offerFilters)
const offerFiltersPending = computed(() =>
    props.offerPage
        ? offerFilters.q !== props.offerPage.filters.q ||
          offerFilters.status !== props.offerPage.filters.status ||
          offerFilters.sort !== props.offerPage.filters.sort ||
          offerFilters.per_page !== props.offerPage.per_page
        : false
)
const offerColumns = computed(() => [
    {
        key: 'created_at',
        label: 'Product / received',
        sortable: true,
        width: 260,
        minWidth: 180,
    },
    { key: 'buyer', label: 'Buyer', width: 190, minWidth: 140 },
    { key: 'proposal', label: 'Proposal', width: 190, minWidth: 150 },
    { key: 'status', label: 'Status', width: 140, minWidth: 110 },
    { key: 'actions', label: 'Actions', width: 100, minWidth: 90 },
])
function offerListUrl(
    pageNumber = props.offerPage?.page || 1,
    filters = props.offerPage?.filters || offerFilters,
    perPage = props.offerPage?.per_page || offerFilters.per_page
) {
    const query = new URLSearchParams({
        section: 'offers',
        sort: filters.sort || 'newest',
        page: String(pageNumber),
        per_page: String(perPage),
    })
    if (filters.q) query.set('q', filters.q)
    if (filters.status) query.set('status', filters.status)
    return '/dashboard?' + query.toString()
}
watch(
    () => props.offerPage,
    (value) => {
        if (
            value &&
            (!offerListLoading.value ||
                offerSignature() === offerSubmittedFilters)
        )
            Object.assign(offerFilters, value.filters, {
                per_page: value.per_page,
            })
    },
    { deep: true, flush: 'sync' }
)
function visitOfferList(url: string) {
    if (
        offerDisposed ||
        section.value !== 'offers' ||
        busy.value ||
        detailRecord.value
    )
        return
    if (offerListLoading.value) {
        offerQueuedSearch = true
        return
    }
    if (offerSearchTimer) clearTimeout(offerSearchTimer)
    offerQueuedSearch = false
    offerSubmittedFilters = offerSignature()
    offerListLoading.value = true
    offerListError.value = ''
    router.get(
        url,
        {},
        {
            preserveState: true,
            preserveScroll: true,
            onError: () => {
                offerListError.value =
                    'Offers could not be loaded. Previous results are still shown.'
            },
            onNetworkError: () => {
                offerListError.value =
                    'Connection lost. Previous offer results are still shown. Try again.'
                return false
            },
            onHttpException: () => {
                offerListError.value =
                    'Offers could not be loaded. Previous results are still shown.'
                return false
            },
            onFinish: () => {
                offerListLoading.value = false
                if (offerQueuedSearch) {
                    offerQueuedSearch = false
                    void nextTick(applyOfferFilters)
                }
            },
        }
    )
}
function applyOfferFilters() {
    visitOfferList(offerListUrl(1, offerFilters, offerFilters.per_page))
}
watch(
    () => offerFilters.q,
    (value) => {
        if (offerSearchTimer) clearTimeout(offerSearchTimer)
        if (section.value !== 'offers' || !props.offerPage) return
        if (offerListLoading.value) {
            if (offerSignature() === offerSubmittedFilters) return
            offerQueuedSearch = true
        } else if (value === props.offerPage.filters.q) return
        offerSearchTimer = setTimeout(applyOfferFilters, value ? 350 : 0)
    }
)
function resetOfferFilters() {
    if (busy.value || detailRecord.value) return
    if (offerSearchTimer) clearTimeout(offerSearchTimer)
    Object.assign(offerFilters, {
        q: '',
        status: '',
        sort: 'newest',
        per_page: 25,
    })
    applyOfferFilters()
}
function changeOfferPage(number: number) {
    if (
        !props.offerPage ||
        offerFiltersPending.value ||
        number < 1 ||
        number > props.offerPage.last_page
    )
        return
    visitOfferList(offerListUrl(number))
}
function changeOfferPageSize(number: number) {
    if (![25, 50, 100].includes(number) || busy.value || detailRecord.value)
        return
    offerFilters.per_page = number
    applyOfferFilters()
}
function sortOfferColumn(sort: { key: string; direction: 'asc' | 'desc' }) {
    if (sort.key !== 'created_at' || busy.value || detailRecord.value) return
    offerFilters.sort = sort.direction === 'asc' ? 'oldest' : 'newest'
    applyOfferFilters()
}
async function refreshSelectedOffer() {
    if (busy.value || !detailRecord.value || detailKind.value !== 'offers')
        return
    const id = detailRecord.value.id
    busy.value = true
    offerRefreshController = new AbortController()
    try {
        const response = await fetch(
            offerListUrl() + '&selected_offer=' + encodeURIComponent(id),
            {
                headers: {
                    'X-Inertia': 'true',
                    'X-Requested-With': 'XMLHttpRequest',
                    Accept: 'application/json',
                },
                signal: offerRefreshController.signal,
                credentials: 'same-origin',
                cache: 'no-store',
            }
        )
        if (!response.ok) {
            if (response.status === 403 || response.status === 404)
                refreshedOfferAccess.value = { id, permissions: {} }
            throw new Error(
                'Could not refresh this offer. The previous review is still displayed.'
            )
        }
        const data = await response.json()
        if (
            offerDisposed ||
            section.value !== 'offers' ||
            detailRecord.value?.id !== id
        )
            return
        const record = data.props?.selectedOffer
        refreshedOfferAccess.value = {
            id,
            permissions: data.props?.recordCapabilities?.offers?.[id] || {},
        }
        if (
            data.props?.selectedOfferStatus !== 'available' ||
            !record ||
            record.id !== id
        ) {
            refreshedOfferAccess.value = { id, permissions: {} }
            throw new Error(
                'This offer is unavailable or no longer accessible. No decision can be submitted from this review.'
            )
        }
        detailRecord.value = record
        selectedOfferVisibility.value =
            data.props.offerVisibility === 'full' ? 'full' : 'supplier'
        pendingOperation.value = ''
        pendingQuote.value = null
        formError.value = ''
    } catch (error) {
        if (
            offerDisposed ||
            section.value !== 'offers' ||
            detailKind.value !== 'offers' ||
            detailRecord.value?.id !== id
        )
            return
        formError.value =
            error instanceof Error
                ? error.message
                : 'Could not refresh this offer. The previous review is still displayed.'
    } finally {
        busy.value = false
    }
}
onBeforeUnmount(() => {
    offerDisposed = true
    offerRefreshController?.abort()
    if (offerSearchTimer) clearTimeout(offerSearchTimer)
})
const rfqFilters = reactive({
    q: props.rfqPage?.filters.q || '',
    status: props.rfqPage?.filters.status || '',
    sort: props.rfqPage?.filters.sort || 'newest',
    per_page: props.rfqPage?.per_page || 25,
})
const rfqListLoading = ref(false)
const rfqListError = ref('')
const rfqReturnUrl = ref('')
let rfqSearchTimer: ReturnType<typeof setTimeout> | undefined
let rfqQueuedSearch = false
let rfqSubmittedFilters = ''
let rfqDisposed = false
let rfqRefreshController: AbortController | undefined
const rfqSignature = () => JSON.stringify(rfqFilters)
const rfqFiltersPending = computed(() =>
    props.rfqPage
        ? rfqFilters.q !== props.rfqPage.filters.q ||
          rfqFilters.status !== props.rfqPage.filters.status ||
          rfqFilters.sort !== props.rfqPage.filters.sort ||
          rfqFilters.per_page !== props.rfqPage.per_page
        : false
)
const rfqColumns = computed(() => [
    {
        key: 'created_at',
        label: 'Request / received',
        sortable: true,
        width: 260,
        minWidth: 180,
    },
    { key: 'buyer', label: 'Buyer', width: 180, minWidth: 140 },
    { key: 'quantity', label: 'Quantity / target', width: 160, minWidth: 130 },
    { key: 'destination', label: 'Destination', width: 150, minWidth: 120 },
    { key: 'status', label: 'Status', width: 130, minWidth: 100 },
    { key: 'actions', label: 'Actions', width: 100, minWidth: 90 },
])
function rfqListUrl(
    pageNumber = props.rfqPage?.page || 1,
    filters = props.rfqPage?.filters || rfqFilters,
    perPage = props.rfqPage?.per_page || rfqFilters.per_page
) {
    const query = new URLSearchParams({
        section: 'rfqs',
        sort: filters.sort || 'newest',
        page: String(pageNumber),
        per_page: String(perPage),
    })
    if (filters.q) query.set('q', filters.q)
    if (filters.status) query.set('status', filters.status)
    return '/dashboard?' + query.toString()
}
watch(
    () => props.rfqPage,
    (value) => {
        if (
            value &&
            (!rfqListLoading.value || rfqSignature() === rfqSubmittedFilters)
        )
            Object.assign(rfqFilters, value.filters, {
                per_page: value.per_page,
            })
    },
    { deep: true, flush: 'sync' }
)
function visitRfqList(url: string) {
    if (
        rfqDisposed ||
        section.value !== 'rfqs' ||
        busy.value ||
        detailRecord.value ||
        quoteTarget.value
    )
        return
    if (rfqListLoading.value) {
        rfqQueuedSearch = true
        return
    }
    if (rfqSearchTimer) clearTimeout(rfqSearchTimer)
    rfqQueuedSearch = false
    rfqSubmittedFilters = rfqSignature()
    rfqListLoading.value = true
    rfqListError.value = ''
    router.get(
        url,
        {},
        {
            preserveState: true,
            preserveScroll: true,
            onError: () => {
                rfqListError.value =
                    'Requests could not be loaded. Previous results are still shown.'
            },
            onNetworkError: () => {
                rfqListError.value =
                    'Connection lost. Previous request results are still shown. Try again.'
                return false
            },
            onHttpException: () => {
                rfqListError.value =
                    'Requests could not be loaded. Previous results are still shown.'
                return false
            },
            onFinish: () => {
                rfqListLoading.value = false
                if (rfqQueuedSearch) {
                    rfqQueuedSearch = false
                    void nextTick(applyRfqFilters)
                }
            },
        }
    )
}
function applyRfqFilters() {
    visitRfqList(rfqListUrl(1, rfqFilters, rfqFilters.per_page))
}
watch(
    () => rfqFilters.q,
    (value) => {
        if (rfqSearchTimer) clearTimeout(rfqSearchTimer)
        if (section.value !== 'rfqs' || !props.rfqPage) return
        if (rfqListLoading.value) {
            if (rfqSignature() === rfqSubmittedFilters) return
            rfqQueuedSearch = true
        } else if (value === props.rfqPage.filters.q) return
        rfqSearchTimer = setTimeout(applyRfqFilters, value ? 350 : 0)
    }
)
function resetRfqFilters() {
    if (busy.value || detailRecord.value || quoteTarget.value) return
    if (rfqSearchTimer) clearTimeout(rfqSearchTimer)
    Object.assign(rfqFilters, {
        q: '',
        status: '',
        sort: 'newest',
        per_page: 25,
    })
    applyRfqFilters()
}
function changeRfqPage(number: number) {
    if (
        !props.rfqPage ||
        rfqFiltersPending.value ||
        number < 1 ||
        number > props.rfqPage.last_page
    )
        return
    visitRfqList(rfqListUrl(number))
}
function changeRfqPageSize(number: number) {
    if (
        ![25, 50, 100].includes(number) ||
        busy.value ||
        detailRecord.value ||
        quoteTarget.value
    )
        return
    rfqFilters.per_page = number
    applyRfqFilters()
}
function sortRfqColumn(sort: { key: string; direction: 'asc' | 'desc' }) {
    if (
        sort.key !== 'created_at' ||
        busy.value ||
        detailRecord.value ||
        quoteTarget.value
    )
        return
    rfqFilters.sort = sort.direction === 'asc' ? 'oldest' : 'newest'
    applyRfqFilters()
}
async function refreshSelectedRfq() {
    if (busy.value || !detailRecord.value || detailKind.value !== 'rfqs') return
    const id = detailRecord.value.id
    busy.value = true
    rfqRefreshController = new AbortController()
    try {
        const response = await fetch(
            rfqListUrl() + '&selected_rfq=' + encodeURIComponent(id),
            {
                headers: {
                    'X-Inertia': 'true',
                    'X-Requested-With': 'XMLHttpRequest',
                    Accept: 'application/json',
                },
                signal: rfqRefreshController.signal,
                credentials: 'same-origin',
                cache: 'no-store',
            }
        )
        if (!response.ok) {
            if (response.status === 403 || response.status === 404)
                refreshedRfqAccess.value = { id, permissions: {}, quotes: {} }
            throw new Error(
                'Could not refresh this request. The previous review is still displayed.'
            )
        }
        const data = await response.json()
        if (
            rfqDisposed ||
            section.value !== 'rfqs' ||
            detailRecord.value?.id !== id
        )
            return
        const record = data.props?.selectedRfq
        refreshedRfqAccess.value = {
            id,
            permissions: data.props?.recordCapabilities?.rfqs?.[id] || {},
            quotes: data.props?.recordCapabilities?.quotes || {},
        }
        if (
            data.props?.selectedRfqStatus !== 'available' ||
            !record ||
            record.id !== id
        ) {
            refreshedRfqAccess.value = { id, permissions: {}, quotes: {} }
            throw new Error(
                'This request is unavailable or no longer accessible. No request decision can be submitted from this review.'
            )
        }
        detailRecord.value = record
        selectedRfqVisibility.value =
            data.props.rfqVisibility === 'full' ? 'full' : 'supplier'
        pendingOperation.value = ''
        pendingQuote.value = null
        formError.value = ''
    } catch (error) {
        if (
            rfqDisposed ||
            section.value !== 'rfqs' ||
            detailKind.value !== 'rfqs' ||
            detailRecord.value?.id !== id
        )
            return
        formError.value =
            error instanceof Error
                ? error.message
                : 'Could not refresh this request. The previous review is still displayed.'
    } finally {
        busy.value = false
    }
}
onBeforeUnmount(() => {
    rfqDisposed = true
    rfqRefreshController?.abort()
    if (rfqSearchTimer) clearTimeout(rfqSearchTimer)
})
const productFilters = reactive({
    q: props.productPage?.filters.q || '',
    category: props.productPage?.filters.category || '',
    sort: props.productPage?.filters.sort || 'newest',
    per_page: props.productPage?.per_page || 25,
})
const productFiltersPending = computed(() =>
    props.productPage
        ? productFilters.q !== props.productPage.filters.q ||
          productFilters.category !== props.productPage.filters.category ||
          productFilters.sort !== props.productPage.filters.sort ||
          productFilters.per_page !== props.productPage.per_page
        : false
)

const productListLoading = ref(false)
const productListError = ref('')
let productSearchTimer: ReturnType<typeof setTimeout> | undefined
let productQueuedSearch = false
let productListDisposed = false
let productSubmittedFilters = ''
const productFilterSignature = () => JSON.stringify(productFilters)
onBeforeUnmount(() => {
    productListDisposed = true
    if (productSearchTimer) clearTimeout(productSearchTimer)
})
const productReturnUrl = ref('')
const productFilterCategories = computed(() =>
    [
        ...new Set([
            ...(props.productCategories || []),
            ...(productFilters.category ? [productFilters.category] : []),
        ]),
    ].sort()
)
function productListUrl(
    pageNumber = props.productPage?.page || 1,
    filters = props.productPage?.filters || productFilters,
    perPage = props.productPage?.per_page || productFilters.per_page
) {
    const query = new URLSearchParams({
        section: 'products',
        sort: filters.sort || 'newest',
        page: String(pageNumber),
        per_page: String(perPage),
    })
    if (filters.q) query.set('q', filters.q)
    if (filters.category) query.set('category', filters.category)
    return '/dashboard?' + query.toString()
}
watch(
    () => props.productPage,
    (value) => {
        if (
            value &&
            (!productListLoading.value ||
                productFilterSignature() === productSubmittedFilters)
        )
            Object.assign(productFilters, value.filters, {
                per_page: value.per_page,
            })
    },
    { deep: true, flush: 'sync' }
)
function visitProductList(url: string) {
    if (productListDisposed || section.value !== 'products') return
    if (
        busy.value ||
        productRefreshing.value ||
        productOpen.value ||
        deleteTarget.value
    )
        return
    if (productListLoading.value) {
        productQueuedSearch = true
        return
    }
    if (productSearchTimer) clearTimeout(productSearchTimer)
    productQueuedSearch = false
    productSubmittedFilters = productFilterSignature()
    productListLoading.value = true
    productListError.value = ''
    router.get(
        url,
        {},
        {
            preserveState: true,
            preserveScroll: true,
            onError: () => {
                productListError.value =
                    'The product list could not be loaded. Your previous results are still shown.'
            },
            onNetworkError: () => {
                productListError.value =
                    'Connection lost. Your previous product results are still shown. Try again.'
                return false
            },
            onHttpException: () => {
                productListError.value =
                    'The product list could not be loaded. Your previous results are still shown.'
                return false
            },
            onFinish: () => {
                productListLoading.value = false
                if (productQueuedSearch) {
                    productQueuedSearch = false
                    void nextTick(applyProductFilters)
                }
            },
        }
    )
}
function applyProductFilters() {
    visitProductList(productListUrl(1, productFilters, productFilters.per_page))
}
watch(
    () => productFilters.q,
    (value) => {
        if (productSearchTimer) clearTimeout(productSearchTimer)
        if (section.value !== 'products' || !props.productPage) return
        if (productListLoading.value) {
            if (productFilterSignature() === productSubmittedFilters) return
            productQueuedSearch = true
        } else if (value === props.productPage.filters.q) return
        productSearchTimer = setTimeout(applyProductFilters, value ? 350 : 0)
    }
)

function resetProductFilters() {
    if (
        busy.value ||
        productRefreshing.value ||
        productOpen.value ||
        deleteTarget.value
    )
        return
    if (productSearchTimer) clearTimeout(productSearchTimer)
    Object.assign(productFilters, {
        q: '',
        category: '',
        sort: 'newest',
        per_page: 25,
    })
    applyProductFilters()
}
function changeProductPage(pageNumber: number) {
    if (productFiltersPending.value) return
    if (
        !props.productPage ||
        pageNumber < 1 ||
        pageNumber > props.productPage.last_page
    )
        return
    visitProductList(productListUrl(pageNumber))
}
const search = ref('')
const recordStatus = ref('all')
const accountRole = ref('all')
const accountVerification = ref('all')
const accountRoles = computed(() =>
    [
        ...new Set([
            ...(props.users || [])
                .map((account) => String(account.role || ''))
                .filter(Boolean),
            ...(accountRole.value === 'all' ? [] : [accountRole.value]),
        ]),
    ].sort()
)
const hasActiveFilters = computed(() =>
    section.value === 'orders' && props.orderPage
        ? Boolean(props.orderPage.filters.q || props.orderPage.filters.status)
        : section.value === 'products' && props.productPage
          ? Boolean(
                props.productPage.filters.q ||
                props.productPage.filters.category
            )
          : Boolean(search.value.trim()) ||
            (section.value === 'users'
                ? accountRole.value !== 'all' ||
                  accountVerification.value !== 'all'
                : recordStatus.value !== 'all')
)
function keepReviewing() {
    pendingOperation.value = ''
    pendingQuote.value = null
}
function clearRecordFilters() {
    if (section.value === 'orders' && props.orderPage) {
        resetOrderFilters()
        return
    }
    if (section.value === 'products' && props.productPage) {
        resetProductFilters()
        return
    }
    search.value = ''
    recordStatus.value = 'all'
    accountRole.value = 'all'
    accountVerification.value = 'all'
}
const statusOptions = computed(() =>
    [
        ...new Set([
            ...records.value
                .map((item) => String(item.status || ''))
                .filter(Boolean),
            ...(recordStatus.value === 'all' ? [] : [recordStatus.value]),
        ]),
    ].sort()
)
const pageError = ref('')
const localNotice = ref('')
watch(section, (next, previous) => {
    if (previous === 'rfqs' && next !== 'rfqs') {
        rfqRefreshController?.abort()
        if (quoteTarget.value && quoteDirty.value)
            pendingQuoteDraft.value = captureQuoteDraft()
        quoteTarget.value = null
    }
    if (
        previous === 'orders' &&
        next !== 'orders' &&
        detailKind.value === 'orders'
    ) {
        orderRefreshController?.abort()
        detailRecord.value = null
        pendingOperation.value = ''
        pendingQuote.value = null
    }
    if (orderSearchTimer) clearTimeout(orderSearchTimer)
    orderQueuedSearch = false
    if (
        previous === 'users' &&
        next !== 'users' &&
        detailKind.value === 'users'
    ) {
        userRefreshController?.abort()
        detailRecord.value = null
        pendingOperation.value = ''
        pendingQuote.value = null
    }
    if (userSearchTimer) clearTimeout(userSearchTimer)
    userQueuedSearch = false
    if (
        previous === 'offers' &&
        next !== 'offers' &&
        detailKind.value === 'offers'
    ) {
        offerRefreshController?.abort()
        detailRecord.value = null
        pendingOperation.value = ''
        pendingQuote.value = null
    }
    if (offerSearchTimer) clearTimeout(offerSearchTimer)
    offerQueuedSearch = false
    if (previous === 'rfqs' && next !== 'rfqs' && detailKind.value === 'rfqs') {
        rfqRefreshController?.abort()
        detailRecord.value = null
        pendingOperation.value = ''
        pendingQuote.value = null
    }
    if (rfqSearchTimer) clearTimeout(rfqSearchTimer)
    rfqQueuedSearch = false
    search.value = ''
    recordStatus.value = 'all'
    accountRole.value = 'all'
    accountVerification.value = 'all'
    if (productSearchTimer) clearTimeout(productSearchTimer)
    productQueuedSearch = false
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
    if (section.value === 'rfqs' && props.rfqPage) return rfqList.value
    if (section.value === 'offers' && props.offerPage) return offerList.value
    if (section.value === 'users' && props.userPage) return props.users || []
    if (section.value === 'orders' && props.orderPage) return orderList.value
    if (section.value === 'products' && props.productPage)
        return productList.value
    const query = search.value.trim().toLowerCase()
    const scoped = records.value.filter((item) => {
        if (section.value === 'users') {
            if (accountRole.value !== 'all' && item.role !== accountRole.value)
                return false
            if (
                accountVerification.value !== 'all' &&
                (item.verified ? 'verified' : 'unverified') !==
                    accountVerification.value
            )
                return false
        } else if (
            recordStatus.value !== 'all' &&
            item.status !== recordStatus.value
        )
            return false
        if (!query) return true
        if (
            section.value === 'users' &&
            (query === 'verified' || query === 'unverified')
        )
            return (item.verified ? 'verified' : 'unverified') === query
        return [
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
            'country',
            'status',
            'role',
        ].some((key) =>
            String(item[key] || '')
                .toLowerCase()
                .includes(query)
        )
    })
    return scoped
})

const emptyTitle = computed(
    () =>
        ({
            products: 'Your catalog starts here',
            orders: 'No orders yet',
            rfqs: 'No open sourcing requests',
            offers: 'No buyer offers yet',
            users: 'No accounts to display',
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
const productRefreshing = ref(false)
const freshProduct = ref<RecordData | null>(null)
const productRefreshMessage = ref('')
const pendingProductDraft = ref<RecordData | null>(null)
const pendingProductDiscard = ref(false)
function discardPendingProductDraft() {
    if (busy.value || productRefreshing.value || !pendingProductDiscard.value)
        return
    pendingProductDraft.value = null
    pendingProductDiscard.value = false
    productRefreshMessage.value = ''
    router.remember(null, productDraftKey)
}
async function fetchSelectedProduct(id: string) {
    const url = productListUrl() + '&selected_product=' + encodeURIComponent(id)
    const response = await fetch(url, {
        headers: {
            'X-Inertia': 'true',
            'X-Requested-With': 'XMLHttpRequest',
            Accept: 'application/json',
        },
        credentials: 'same-origin',
        cache: 'no-store',
    })
    if (!response.ok)
        throw new Error(
            'The current saved product could not be checked. Your draft is unchanged.'
        )
    const data = await response.json()
    if (productListDisposed || section.value !== 'products')
        throw new Error(
            'Product lookup stopped because you left the product catalog. Your saved draft remains available.'
        )
    const record = data.props?.selectedProduct
    const permissions = data.props?.recordCapabilities?.products?.[id] || {}
    refreshedProductAccess.value = { id, permissions }
    if (
        data.props?.selectedProductStatus !== 'available' ||
        !record ||
        record.id !== id ||
        permissions.update !== true
    )
        throw new Error(
            'This product is unavailable or you no longer have edit access. Your draft has not been discarded.'
        )
    return record as RecordData
}
async function checkSavedProduct() {
    if (busy.value || productRefreshing.value || !editingProduct.value) return
    productRefreshing.value = true
    freshProduct.value = null
    productRefreshMessage.value = ''
    try {
        freshProduct.value = await fetchSelectedProduct(editingProduct.value.id)
        productRefreshMessage.value =
            'Current saved details loaded for review. Your draft has not changed.'
    } catch (error) {
        productRefreshMessage.value =
            error instanceof Error
                ? error.message
                : 'Could not check the saved product.'
    } finally {
        productRefreshing.value = false
    }
}
function replaceWithSavedProduct() {
    if (busy.value || productRefreshing.value || !freshProduct.value) return
    loadProductForm(freshProduct.value)
    freshProduct.value = null
    productRefreshMessage.value = 'Loaded the current saved version.'
}
async function recoverOffpageProduct() {
    if (busy.value || productRefreshing.value || !pendingProductDraft.value)
        return
    productRefreshing.value = true
    try {
        const record = await fetchSelectedProduct(
            pendingProductDraft.value.productId
        )
        Object.assign(
            productForm,
            emptyProduct(),
            safeProductDraft(pendingProductDraft.value.form)
        )
        productForm.digital_download_url = record.digital_download_url || ''
        editingProduct.value = record
        productSnapshot.value = ''
        productRecovered.value = true
        productOpen.value = true
        productReturnUrl.value = productListUrl()
        pendingProductDraft.value = null
        pendingProductDiscard.value = false
        productRefreshMessage.value = ''
    } catch (error) {
        productRefreshMessage.value =
            error instanceof Error
                ? error.message
                : 'Could not reopen the saved draft.'
    } finally {
        productRefreshing.value = false
    }
}
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
        ? [props.selectedProduct, ...(props.products || [])]
              .filter(Boolean)
              .find((p) => p?.id === restoredProductDraft.productId)
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
        productReturnUrl.value = productListUrl()
    } else if (restoredProductDraft.productId) {
        pendingProductDraft.value = restoredProductDraft
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
    pendingProductDraft.value = null
    freshProduct.value = null
    refreshedProductAccess.value = null
    productRefreshMessage.value = ''
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
    if (pendingProductDraft.value) {
        productRefreshMessage.value =
            'Reopen or explicitly discard the saved draft before starting another.'
        return
    }
    if (
        busy.value ||
        productRefreshing.value ||
        productListLoading.value ||
        productOpen.value
    )
        return
    if (productSearchTimer) clearTimeout(productSearchTimer)
    productQueuedSearch = false
    if (
        product
            ? !recordCan('products', product.id, 'update')
            : !props.capabilities?.createProducts
    )
        return
    productReturnUrl.value = productListUrl()
    freshProduct.value = null
    productRefreshMessage.value = ''
    loadProductForm(product)
}
function loadProductForm(product?: RecordData) {
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
    if (busy.value || productRefreshing.value) return
    if (
        editingProduct.value &&
        !recordCan('products', editingProduct.value.id, 'update')
    ) {
        formError.value =
            'Edit access could not be confirmed. Check the current saved version before saving.'
        return
    }
    validatePriceTiers(event.target as HTMLFormElement)
    if (!validate(event)) return
    const id = editingProduct.value?.id
    const returnUrl = productReturnUrl.value || productListUrl()
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
                    returnUrl,
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
    if (busy.value || productRefreshing.value || productListLoading.value)
        return
    if (productSearchTimer) clearTimeout(productSearchTimer)
    productQueuedSearch = false
    if (!recordCan('products', product.id, 'delete')) return
    productReturnUrl.value = productListUrl()
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
            void nextTick(() =>
                router.get(
                    productReturnUrl.value || productListUrl(),
                    {},
                    { preserveState: true, preserveScroll: true }
                )
            )
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
const quoteDraftKey = `soko.quote-editor.${props.currentUser.id}`
type QuoteDraft = {
    rfqId: string
    form: { unit_price: number; lead_time_days: number; notes: string }
}
const pendingQuoteDraft = ref<QuoteDraft | null>(null)
const quoteDiscardPending = ref(false)
const quoteRecoveryError = ref('')
const restoredQuoteDraft = router.restore(quoteDraftKey) as
    QuoteDraft | undefined
if (restoredQuoteDraft?.rfqId && restoredQuoteDraft.form)
    pendingQuoteDraft.value = restoredQuoteDraft
function captureQuoteDraft(): QuoteDraft {
    return {
        rfqId: String(quoteTarget.value?.id || ''),
        form: {
            unit_price: quoteForm.unit_price,
            lead_time_days: quoteForm.lead_time_days,
            notes: quoteForm.notes,
        },
    }
}
watch(
    [quoteForm, quoteTarget, pendingQuoteDraft],
    () => {
        router.remember(
            quoteTarget.value && quoteDirty.value
                ? captureQuoteDraft()
                : pendingQuoteDraft.value,
            quoteDraftKey
        )
    },
    { deep: true }
)
function discardQuoteDraft() {
    if (busy.value) return
    pendingQuoteDraft.value = null
    quoteDiscardPending.value = false
    quoteRecoveryError.value = ''
    router.remember(null, quoteDraftKey)
}
async function recoverQuoteDraft() {
    if (
        busy.value ||
        rfqListLoading.value ||
        quoteTarget.value ||
        detailRecord.value ||
        !pendingQuoteDraft.value
    )
        return
    const draft = pendingQuoteDraft.value
    busy.value = true
    quoteRecoveryError.value = ''
    rfqRefreshController = new AbortController()
    try {
        const response = await fetch(
            rfqListUrl() + '&selected_rfq=' + encodeURIComponent(draft.rfqId),
            {
                headers: {
                    'X-Inertia': 'true',
                    'X-Requested-With': 'XMLHttpRequest',
                    Accept: 'application/json',
                },
                credentials: 'same-origin',
                cache: 'no-store',
                signal: rfqRefreshController.signal,
            }
        )
        if (!response.ok)
            throw new Error(
                'Could not reopen the quotation. Your draft is still saved.'
            )
        const data = await response.json()
        if (
            rfqDisposed ||
            section.value !== 'rfqs' ||
            pendingQuoteDraft.value?.rfqId !== draft.rfqId
        )
            return
        const record = data.props?.selectedRfq,
            permissions =
                data.props?.recordCapabilities?.rfqs?.[draft.rfqId] || {}
        if (
            data.props?.selectedRfqStatus !== 'available' ||
            record?.id !== draft.rfqId ||
            !permissions.quote
        )
            throw new Error(
                'This request is unavailable or no longer accepts your quotations. Your draft is still saved.'
            )
        refreshedRfqAccess.value = {
            id: record.id,
            permissions,
            quotes: data.props?.recordCapabilities?.quotes || {},
        }
        quoteTarget.value = record
        Object.assign(quoteForm, {
            unit_price: draft.form.unit_price ?? 0,
            lead_time_days: draft.form.lead_time_days ?? 7,
            notes: String(draft.form.notes || ''),
        })
        quoteSnapshot.value = ''
        rfqReturnUrl.value = rfqListUrl()
        pendingQuoteDraft.value = null
        quoteDiscardPending.value = false
        formError.value = ''
        invalidFields.value = []
    } catch (error) {
        if (!rfqDisposed && section.value === 'rfqs')
            quoteRecoveryError.value =
                error instanceof Error
                    ? error.message
                    : 'Could not reopen the quotation.'
    } finally {
        busy.value = false
    }
}
const ownQuote = (rfq: RecordData) =>
    (rfq.quotes || []).find(
        (quote: RecordData) => quote.vendor_id === props.currentUser.id
    )
const openQuote = (rfq: RecordData) => {
    if (
        busy.value ||
        rfqListLoading.value ||
        quoteTarget.value ||
        !recordCan('rfqs', rfq.id, 'quote')
    )
        return
    rfqReturnUrl.value = rfqListUrl()
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
    if (!busy.value) {
        quoteTarget.value = null
        pendingQuoteDraft.value = null
        router.remember(null, quoteDraftKey)
    }
}
const saveQuote = (event: Event) => {
    if (
        busy.value ||
        !quoteTarget.value ||
        !recordCan('rfqs', quoteTarget.value.id, 'quote') ||
        !validate(event)
    )
        return
    const returnUrl = rfqReturnUrl.value
    mutate(
        '/dashboard/rfqs/quote?section=rfqs',
        { rfq_id: quoteTarget.value.id, ...quoteForm },
        'quote',
        () => {
            quoteTarget.value = null
            localNotice.value = 'Quotation sent.'
            if (returnUrl)
                void nextTick(() =>
                    router.get(
                        returnUrl,
                        {},
                        { preserveState: true, preserveScroll: true }
                    )
                )
        },
        true
    )
}
const decideOffer = (offer: RecordData, status: 'accepted' | 'declined') => {
    if (busy.value || !recordCan('offers', offer.id, 'decide')) return
    const returnUrl = offerReturnUrl.value || offerListUrl()
    mutate(
        `/dashboard/offers/${offer.id}/decision?section=offers`,
        { status },
        `offer-${offer.id}-${status}`,
        () => {
            localNotice.value =
                status === 'accepted' ? 'Offer accepted.' : 'Offer declined.'
            detailRecord.value = null
            pendingOperation.value = ''
            void nextTick(() =>
                router.get(
                    returnUrl,
                    {},
                    { preserveState: true, preserveScroll: true }
                )
            )
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
.order-search-toolbar {
    display: flex;
    flex-wrap: wrap;
    align-items: end;
    gap: 0.75rem;
    width: 100%;
}
.order-search-toolbar label {
    display: grid;
    gap: 0.35rem;
    font-size: 0.8125rem;
    flex: 1;
    min-width: min(100%, 16rem);
}
.order-search-toolbar input,
.order-filter-fields select {
    min-height: 44px;
    padding: 0.6rem 0.75rem;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    background: var(--soko-surface);
    color: var(--soko-ink);
}
.order-search-toolbar p {
    margin: 0;
    font-size: 0.8125rem;
    color: var(--soko-muted);
}
.order-filter-fields {
    display: flex;
    gap: 1rem;
    flex-wrap: wrap;
}
.order-filter-fields label {
    display: grid;
    gap: 0.35rem;
    font-size: 0.8125rem;
    flex: 1;
    min-width: min(100%, 10rem);
}

.product-list-toolbar {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    align-items: end;
    margin-bottom: 1.5rem;
}
.product-list-toolbar label {
    display: grid;
    gap: 0.4rem;
    font-size: 0.8125rem;
    font-weight: 650;
    flex: 1;
    min-width: min(100%, 9rem);
}
.product-list-toolbar .product-query {
    flex: 2;
    min-width: min(100%, 16rem);
}
.product-list-toolbar input,
.product-list-toolbar select {
    min-height: 44px;
    padding: 0.6rem 0.75rem;
    width: 100%;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    background: var(--soko-surface);
    color: var(--soko-ink);
}
.product-results-count {
    flex-basis: 100%;
    margin: 0;
    font-size: 0.875rem;
    color: var(--soko-muted);
}
.product-pagination {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    width: 100%;
    flex-wrap: wrap;
}
.product-version-review,
.product-recovery-banner {
    padding: 1rem;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    background: var(--soko-canvas);
    margin-bottom: 1rem;
}
.product-version-review p,
.product-recovery-banner p {
    font-size: 0.875rem;
    line-height: 1.6;
    margin-block: 0.5rem;
}

.workspace-toolbar {
    flex-wrap: wrap;
}
.workspace-search .icon-button {
    min-height: 44px;
    min-width: 44px;
}
.operation-identity {
    min-height: 44px;
    display: inline-flex;
    align-items: center;
    padding: 0.5rem 0;
    border: 0;
    background: transparent;
    color: var(--soko-navy);
    text-align: left;
    font: inherit;
    font-weight: 700;
    text-decoration: underline;
    text-decoration-thickness: 1px;
    text-underline-offset: 0.2em;
    cursor: pointer;
    overflow-wrap: anywhere;
}
.operation-identity:focus-visible {
    outline: 2px solid var(--soko-navy);
    outline-offset: 3px;
    border-radius: var(--soko-radius);
}
.operation-identity:disabled {
    opacity: 0.65;
    cursor: wait;
}
@media (max-width: 640px) {
    .operation-status-filter {
        width: 100%;
        justify-content: space-between;
    }
    .operation-status-filter select {
        flex: 1;
        max-width: 70%;
    }
}

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
