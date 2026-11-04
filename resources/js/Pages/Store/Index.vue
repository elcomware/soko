<template>
    <Head title="Global marketplace" />
    <div
        ref="storefrontElement"
        class="storefront immersive-storefront"
        :class="{ 'motion-disabled': !motionEnabled || prefersReducedMotion }"
    >
        <section
            ref="heroElement"
            class="immersive-hero"
            :class="{
                'ambient-paused':
                    !motionEnabled || prefersReducedMotion || !heroVisible,
            }"
            aria-labelledby="storefront-title"
        >
            <div class="immersive-visual" aria-hidden="true">
                <img
                    v-for="(story, index) in stories"
                    :key="story.id"
                    :src="story.image"
                    alt=""
                    width="1400"
                    height="934"
                    :fetchpriority="index === 0 ? 'high' : 'auto'"
                    class="immersive-backdrop"
                    :class="{ 'is-active': activeStory === index }"
                />
            </div>
            <div class="immersive-shade" aria-hidden="true"></div>
            <svg
                class="trade-route"
                viewBox="0 0 1000 600"
                fill="none"
                aria-hidden="true"
            >
                <path class="route-track" d="M 340 530 Q 510 190 950 170" />
                <path class="route-flow" d="M 340 530 Q 510 190 950 170" />
                <circle cx="340" cy="530" r="5" />
                <circle cx="950" cy="170" r="5" />
            </svg>
            <div class="site-container immersive-inner">
                <div class="immersive-copy">
                    <p class="immersive-eyebrow">
                        <span></span> Soko / {{ currentStory.label }}
                    </p>
                    <div
                        class="story-announcement"
                        aria-live="polite"
                        aria-atomic="true"
                    >
                        <div class="story-copy" :key="currentStory.id">
                            <h1 id="storefront-title">
                                {{ currentStory.headline }}<br /><span>{{
                                    currentStory.accent
                                }}</span>
                            </h1>
                            <p class="immersive-description">
                                {{ currentStory.description }}
                            </p>
                        </div>
                    </div>
                    <div class="immersive-actions">
                        <a
                            v-if="activeStory === 0"
                            href="#catalog"
                            class="button button-amber"
                            >Explore the marketplace <SokoIcon name="arrow"
                        /></a>
                        <Link
                            v-else
                            :href="
                                catalogHref({
                                    format: currentStory.format,
                                    category: 'all',
                                }) + '#catalog'
                            "
                            class="button button-amber"
                            >{{ currentStory.action }} <SokoIcon name="arrow"
                        /></Link>
                        <button
                            class="immersive-secondary"
                            @click="openRfqModal"
                        >
                            Request a quotation <SokoIcon name="arrow" />
                        </button>
                    </div>
                </div>
                <Link
                    v-if="featuredProduct"
                    :href="'/products/' + featuredProduct.id"
                    class="hero-market-card"
                >
                    <span class="market-card-label"
                        ><SokoIcon name="layers" /> In the marketplace</span
                    >
                    <div class="market-card-product">
                        <img
                            :src="featuredProduct.image_url"
                            alt=""
                            width="64"
                            height="64"
                        />
                        <div>
                            <span>{{ featuredProduct.category }}</span>
                            <h2>{{ featuredProduct.title }}</h2>
                        </div>
                    </div>
                    <div class="market-card-bottom">
                        <div>
                            <strong>{{
                                formatPrice(featuredProduct.price)
                            }}</strong
                            ><small
                                >Min. {{ featuredProduct.moq }}
                                {{
                                    featuredProduct.moq === 1 ? 'unit' : 'units'
                                }}
                                · {{ featuredProduct.vendor_country }}</small
                            >
                        </div>
                        <span class="market-card-arrow"
                            ><SokoIcon name="arrow"
                        /></span>
                    </div>
                </Link>
                <div class="hero-bottom-bar">
                    <div
                        class="story-selectors"
                        role="group"
                        aria-label="Explore marketplace stories"
                    >
                        <button
                            v-for="(story, index) in stories"
                            :key="story.id"
                            :aria-pressed="activeStory === index"
                            @click="activeStory = index"
                        >
                            <SokoIcon :name="story.icon" /><span
                                ><strong>{{ story.shortLabel }}</strong
                                ><small>{{ story.caption }}</small></span
                            >
                        </button>
                    </div>
                    <button
                        class="motion-toggle"
                        :aria-pressed="!motionEnabled || prefersReducedMotion"
                        :disabled="prefersReducedMotion"
                        @click="motionEnabled = !motionEnabled"
                    >
                        <SokoIcon
                            :name="
                                motionEnabled && !prefersReducedMotion
                                    ? 'pause'
                                    : 'play'
                            "
                        />{{
                            prefersReducedMotion
                                ? 'Reduced motion'
                                : motionEnabled
                                  ? 'Pause motion'
                                  : 'Enable motion'
                        }}
                    </button>
                </div>
            </div>
            <a
                href="#discover"
                class="hero-scroll-link"
                aria-label="Discover more on Soko"
                ><span>More possibilities below</span><SokoIcon name="down"
            /></a>
        </section>

        <div class="storefront-benefits">
            <div
                class="site-container benefit-strip"
                aria-label="Ways to trade on Soko"
            >
                <div class="benefit">
                    <SokoIcon name="layers" />
                    <div>
                        <strong>Value that scales</strong
                        ><span>Discover wholesale volume pricing</span>
                    </div>
                </div>
                <div class="benefit">
                    <SokoIcon name="globe" />
                    <div>
                        <strong>Your business. Your terms.</strong
                        ><span>Buy directly or start a conversation</span>
                    </div>
                </div>
                <div class="benefit">
                    <SokoIcon name="bolt" />
                    <div>
                        <strong>Tools for your digital business.</strong
                        ><span>Software, licenses and business assets</span>
                    </div>
                </div>
            </div>
        </div>
        <div class="site-container storefront-sections">
            <section
                id="discover"
                class="sector-section reveal-item"
                aria-labelledby="sector-title"
            >
                <div class="section-heading">
                    <div>
                        <p class="eyebrow">
                            Different worlds. Shared ambition.
                        </p>
                        <h2 id="sector-title">Built for your next move.</h2>
                    </div>
                    <a href="#catalog" class="section-explore"
                        >Explore all products <SokoIcon name="arrow"
                    /></a>
                </div>
                <div class="sector-grid">
                    <Link
                        v-for="sector in sectors"
                        :key="sector.title"
                        :href="catalogHref(sector.query) + '#catalog'"
                        class="sector-card"
                        ><img
                            :src="sector.image"
                            alt=""
                            width="800"
                            height="600"
                            loading="lazy" />
                        <div class="sector-card-shade"></div>
                        <div class="sector-card-content">
                            <SokoIcon :name="sector.icon" />
                            <p>{{ sector.caption }}</p>
                            <h3>{{ sector.title }}</h3>
                            <span class="sector-card-action"
                                >Explore sector <SokoIcon name="arrow"
                            /></span></div
                    ></Link>
                </div>
            </section>

            <section
                id="catalog"
                class="catalog-section reveal-item"
                aria-labelledby="catalog-title"
            >
                <div class="section-heading">
                    <div>
                        <p class="eyebrow">Find your next opportunity</p>
                        <h2 id="catalog-title">Products for your business</h2>
                    </div>
                    <span class="catalog-count" role="status"
                        >{{ productList.length }}
                        {{ productList.length === 1 ? 'product' : 'products' }}
                        to explore</span
                    ><button
                        class="button button-outline mobile-filter-button"
                        :aria-expanded="filtersOpen"
                        aria-controls="catalog-filters"
                        @click="filtersOpen = !filtersOpen"
                    >
                        <SokoIcon name="filter" /> Filters
                    </button>
                </div>
                <div class="catalog-layout">
                    <aside
                        id="catalog-filters"
                        class="catalog-filters"
                        :class="{ 'is-expanded': filtersOpen }"
                        aria-label="Product filters"
                    >
                        <div class="filter-heading">
                            <h3>Refine your search</h3>
                            <button class="text-button" @click="resetFilters">
                                Reset all
                            </button>
                        </div>
                        <div class="filter-group">
                            <h4 class="filter-label">Categories</h4>
                            <button
                                v-for="category in categories"
                                :key="category.id"
                                class="category-option"
                                :aria-pressed="activeCategory === category.id"
                                @click="setCategory(category.id)"
                            >
                                <SokoIcon
                                    :name="
                                        categoryIcons[category.id] || 'layers'
                                    "
                                /><span>{{ category.name }}</span>
                            </button>
                        </div>
                        <div class="filter-group">
                            <h4 class="filter-label">Condition</h4>
                            <button
                                v-for="condition in conditions"
                                :key="condition.id"
                                class="condition-option"
                                :aria-pressed="activeCondition === condition.id"
                                @click="setCondition(condition.id)"
                            >
                                <span
                                    class="condition-mark"
                                    aria-hidden="true"
                                ></span
                                >{{ condition.label }}
                            </button>
                        </div>
                        <div class="filter-group">
                            <label for="sourcing-region" class="filter-label"
                                >Sourcing region</label
                            ><select
                                id="sourcing-region"
                                v-model="filterRegion"
                                class="filter-select"
                                @change="onFilterRegionChange"
                            >
                                <option value="global">All regions</option>
                                <option
                                    v-for="region in regions.filter(
                                        (region) => region.id !== 'global'
                                    )"
                                    :key="region.id"
                                    :value="region.id"
                                >
                                    {{ region.name }}
                                </option>
                            </select>
                        </div>
                        <div class="filter-help">
                            <SokoIcon name="document" />
                            <h4>Have something specific in mind?</h4>
                            <p>
                                Share your requirements and let suppliers come
                                to you.
                            </p>
                            <button class="text-button" @click="openRfqModal">
                                Post a sourcing request
                                <SokoIcon name="arrow" />
                            </button>
                        </div>
                    </aside>
                    <div class="catalog-results">
                        <p
                            v-if="
                                formError && !offerModalProduct && !showRfqModal
                            "
                            class="form-error"
                            role="alert"
                        >
                            {{ formError }}
                        </p>
                        <div class="catalog-toolbar">
                            <div
                                class="format-tabs"
                                role="group"
                                aria-label="Product format"
                            >
                                <button
                                    v-for="format in formats"
                                    :key="format.id"
                                    :aria-pressed="activeFormat === format.id"
                                    @click="setFormat(format.id)"
                                >
                                    {{ format.label }}
                                </button>
                            </div>
                            <label class="catalog-sort"
                                ><span class="sr-only">Sort products</span
                                ><select
                                    :value="sortOrder"
                                    @change="
                                        visitCatalog({
                                            sort: (
                                                $event.target as HTMLSelectElement
                                            ).value,
                                        })
                                    "
                                >
                                    <option value="newest">
                                        Newest listings
                                    </option>
                                    <option value="price_asc">
                                        Price: low to high
                                    </option>
                                    <option value="price_desc">
                                        Price: high to low
                                    </option>
                                    <option value="title">Product name</option>
                                </select></label
                            ><span class="result-currency"
                                >Prices in {{ activeCurrency.code }}</span
                            >
                        </div>
                        <div
                            v-if="productList.length === 0"
                            class="empty-state"
                            role="status"
                        >
                            <SokoIcon name="search" />
                            <h3>No matching products</h3>
                            <p>
                                {{
                                    searchQuery
                                        ? `We couldn’t find results for “${searchQuery}”. Try another search or adjust your filters.`
                                        : 'Try a different category, condition or sourcing region.'
                                }}
                            </p>
                            <button
                                class="button button-navy"
                                @click="resetFilters"
                            >
                                Clear search and filters
                            </button>
                        </div>
                        <div v-else class="product-grid">
                            <article
                                v-for="product in productList"
                                :key="product.id"
                                class="product-card"
                            >
                                <Link
                                    :href="'/products/' + product.id"
                                    class="product-image-link"
                                    :class="{
                                        'image-missing': missingImages.has(
                                            product.id
                                        ),
                                    }"
                                    :aria-label="'View ' + product.title"
                                    ><img
                                        :src="product.image_url"
                                        :alt="product.title"
                                        loading="lazy"
                                        width="480"
                                        height="330"
                                        @error="missingImages.add(product.id)"
                                    /><span
                                        class="image-fallback"
                                        aria-hidden="true"
                                        ><SokoIcon
                                            :name="
                                                product.type === 'physical'
                                                    ? 'box'
                                                    : 'cloud'
                                            " /></span
                                    ><span class="product-format">{{
                                        formatLabel(product.format)
                                    }}</span></Link
                                >
                                <div class="product-body">
                                    <p class="product-category">
                                        {{ product.category }}
                                    </p>
                                    <h3>
                                        <Link
                                            :href="'/products/' + product.id"
                                            >{{ product.title }}</Link
                                        >
                                    </h3>
                                    <div class="supplier-line">
                                        <SokoIcon
                                            v-if="product.vendor_verified"
                                            name="shield"
                                        /><span :title="product.vendor_name">{{
                                            product.vendor_name
                                        }}</span
                                        ><span
                                            v-if="product.vendor_verified"
                                            class="sr-only"
                                            >Verified supplier</span
                                        >
                                    </div>
                                    <div class="product-meta">
                                        <span>{{ product.vendor_country }}</span
                                        ><span
                                            v-if="product.condition !== 'new'"
                                            >·
                                            {{
                                                product.condition === 'open_box'
                                                    ? 'Open box'
                                                    : 'Refurbished'
                                            }}</span
                                        >
                                    </div>
                                    <div class="product-price">
                                        {{ formatPrice(product.price) }}
                                        <small
                                            >/
                                            {{
                                                product.type === 'physical'
                                                    ? 'unit'
                                                    : 'license'
                                            }}</small
                                        >
                                    </div>
                                    <p class="product-moq">
                                        Minimum order: {{ product.moq }}
                                        {{
                                            product.moq === 1 ? 'unit' : 'units'
                                        }}
                                    </p>
                                    <dl
                                        v-if="product.price_tiers?.length"
                                        class="price-tiers"
                                        aria-label="Volume pricing"
                                    >
                                        <div
                                            v-for="(
                                                tier, index
                                            ) in product.price_tiers.slice(
                                                0,
                                                3
                                            )"
                                            :key="index"
                                        >
                                            <dt>
                                                {{
                                                    tier.max_qty &&
                                                    tier.max_qty < 999
                                                        ? `${tier.min_qty}–${tier.max_qty}`
                                                        : `${tier.min_qty}+`
                                                }}
                                                units
                                            </dt>
                                            <dd>
                                                {{ formatPrice(tier.price) }}
                                            </dd>
                                        </div>
                                    </dl>
                                    <div class="product-actions">
                                        <Link
                                            :href="'/products/' + product.id"
                                            class="button button-subtle"
                                            >View details
                                            <SokoIcon name="arrow" /></Link
                                        ><button
                                            v-if="product.accepts_offers"
                                            class="button button-navy"
                                            @click="openOfferModal(product)"
                                        >
                                            Make an offer</button
                                        ><button
                                            v-else
                                            class="button button-navy"
                                            :disabled="
                                                submitting ||
                                                product.stock === 0
                                            "
                                            @click="openOrderModal(product)"
                                        >
                                            {{
                                                product.type === 'physical'
                                                    ? 'Add to cart'
                                                    : 'Add to cart'
                                            }}
                                        </button>
                                    </div>
                                </div>
                            </article>
                        </div>
                    </div>
                </div>
            </section>

            <section
                class="sourcing-studio reveal-item"
                aria-labelledby="sourcing-banner-title"
            >
                <div class="sourcing-studio-copy">
                    <p class="eyebrow">From possibility to specification</p>
                    <h2 id="sourcing-banner-title">
                        {{
                            storefrontContent?.sourcing_heading ??
                            'Your next big thing starts with a brief.'
                        }}
                    </h2>
                    <p>
                        {{
                            storefrontContent?.sourcing_description ??
                            'Some needs don’t fit a search box. Share your specifications and give suppliers a clear place to start.'
                        }}
                    </p>
                    <button class="button button-amber" @click="openRfqModal">
                        Create a sourcing request <SokoIcon name="arrow" />
                    </button>
                </div>
                <div class="sourcing-studio-steps">
                    <div>
                        <span>01</span>
                        <div>
                            <h3>Tell your story</h3>
                            <p>What do you need to source?</p>
                        </div>
                        <SokoIcon name="document" />
                    </div>
                    <div>
                        <span>02</span>
                        <div>
                            <h3>Define the details</h3>
                            <p>Quantity, destination and specifications.</p>
                        </div>
                        <SokoIcon name="filter" />
                    </div>
                    <div>
                        <span>03</span>
                        <div>
                            <h3>Open the conversation</h3>
                            <p>Make your request available to suppliers.</p>
                        </div>
                        <SokoIcon name="globe" />
                    </div>
                </div>
            </section>
            <section
                id="rfq"
                class="rfq-section reveal-item"
                aria-labelledby="rfq-title"
            >
                <div class="section-heading">
                    <div>
                        <p class="eyebrow">Opportunity goes both ways</p>
                        <h2 id="rfq-title">What buyers are looking for.</h2>
                        <p>
                            Real sourcing requests. A starting point for your
                            next business connection.
                        </p>
                    </div>
                    <button class="button button-outline" @click="openRfqModal">
                        Post a request <SokoIcon name="arrow" />
                    </button>
                </div>
                <div v-if="rfqs?.length" class="request-grid">
                    <article
                        v-for="request in rfqs"
                        :key="request.id"
                        class="request-card"
                    >
                        <div class="request-card-top">
                            <span
                                ><SokoIcon
                                    :name="
                                        categoryIcons[request.category] ||
                                        'document'
                                    "
                                />{{ request.category }}</span
                            ><small>Buyer request</small>
                        </div>
                        <h3>{{ request.title }}</h3>
                        <p class="request-reference">
                            {{ request.rfq_number }}
                        </p>
                        <dl class="request-facts">
                            <div>
                                <dt>Quantity</dt>
                                <dd>
                                    {{ request.quantity.toLocaleString() }}
                                    <small>units</small>
                                </dd>
                            </div>
                            <div>
                                <dt>Target per unit</dt>
                                <dd>
                                    {{
                                        request.target_price
                                            ? formatPrice(request.target_price)
                                            : 'Not specified'
                                    }}
                                </dd>
                            </div>
                            <div class="request-destination">
                                <dt>Destination</dt>
                                <dd>{{ request.destination_country }}</dd>
                            </div>
                        </dl>
                        <Link href="/dashboard" class="request-quote-link"
                            >Submit a quote
                            <span><SokoIcon name="arrow" /></span
                        ></Link>
                    </article>
                </div>
                <div v-else class="empty-state">
                    <h3>No open sourcing requests</h3>
                    <p>
                        Create a request to share your requirements with
                        suppliers.
                    </p>
                    <button class="button button-navy" @click="openRfqModal">
                        Post a sourcing request
                    </button>
                </div>
            </section>
        </div>

        <SokoDialog
            :open="!!offerModalProduct"
            title="Make an offer"
            @close="closeOfferModal"
        >
            <template v-if="offerModalProduct">
                <div class="dialog-product">
                    <img
                        :src="offerModalProduct.image_url"
                        alt=""
                        width="56"
                        height="56"
                    />
                    <div>
                        <strong>{{ offerModalProduct.title }}</strong>
                        <p>
                            {{ offerModalProduct.vendor_name }} · Listed at
                            {{ formatPrice(offerModalProduct.price) }}
                        </p>
                    </div>
                </div>
                <form
                    class="soko-form"
                    novalidate
                    :aria-busy="submitting"
                    @submit.prevent="submitOfferForm"
                >
                    <p
                        v-if="formError"
                        id="form-feedback"
                        class="form-error"
                        role="alert"
                    >
                        {{ formError }}
                    </p>
                    <div class="form-row">
                        <div class="soko-field">
                            <label for="offer-price"
                                >Your price per unit (USD)</label
                            ><input
                                id="offer-price"
                                v-model.number="offerPrice"
                                v-bind="fieldState('offer-price')"
                                type="number"
                                step="0.01"
                                min="0.01"
                                required
                            />
                        </div>
                        <div class="soko-field">
                            <label for="offer-quantity">Quantity</label
                            ><input
                                id="offer-quantity"
                                v-model.number="offerQty"
                                v-bind="fieldState('offer-quantity')"
                                type="number"
                                min="1"
                                required
                            />
                        </div>
                    </div>
                    <div class="form-row">
                        <div class="soko-field">
                            <label for="offer-name">Your name</label
                            ><input
                                id="offer-name"
                                v-model="offerName"
                                v-bind="fieldState('offer-name')"
                                autocomplete="name"
                                required
                            />
                        </div>
                        <div class="soko-field">
                            <label for="offer-email">Work email</label
                            ><input
                                id="offer-email"
                                v-model="offerEmail"
                                v-bind="fieldState('offer-email')"
                                type="email"
                                autocomplete="email"
                                required
                            />
                        </div>
                    </div>
                    <div class="soko-field">
                        <label for="offer-message"
                            >Message to supplier (optional)</label
                        ><textarea
                            id="offer-message"
                            v-model="offerMessage"
                            class="resize-none"
                            rows="4"
                            placeholder="Delivery destination, preferred terms or questions for the supplier"
                        ></textarea>
                    </div>
                    <dl class="order-summary">
                        <div>
                            <dt>Total offer</dt>
                            <dd>
                                ${{
                                    (offerPrice * offerQty).toLocaleString(
                                        'en-US',
                                        {
                                            minimumFractionDigits: 2,
                                            maximumFractionDigits: 2,
                                        }
                                    )
                                }}
                                USD
                            </dd>
                        </div>
                    </dl>
                    <button
                        type="submit"
                        class="button button-navy button-full"
                        :disabled="submitting"
                    >
                        {{
                            submitting
                                ? 'Sending offer…'
                                : 'Send offer to supplier'
                        }}
                    </button>
                </form>
            </template>
        </SokoDialog>
        <SokoDialog
            :open="showRfqModal"
            title="Create a sourcing request"
            @close="closeRfqModal"
        >
            <form
                class="soko-form"
                novalidate
                :aria-busy="submitting"
                @submit.prevent="submitRfqForm"
            >
                <p class="form-note">
                    Describe what your business needs so suppliers can prepare a
                    relevant quotation.
                </p>
                <p
                    v-if="formError"
                    id="form-feedback"
                    class="form-error"
                    role="alert"
                >
                    {{ formError }}
                </p>
                <div class="soko-field">
                    <label for="rfq-title-input">What are you sourcing?</label
                    ><input
                        id="rfq-title-input"
                        v-model="rfqTitle"
                        v-bind="fieldState('rfq-title-input')"
                        placeholder="e.g. Stainless steel flow control valves"
                        required
                    />
                </div>
                <div class="soko-field">
                    <label for="rfq-category">Product category</label>
                    <select
                        id="rfq-category"
                        v-model="rfqCategory"
                        v-bind="fieldState('rfq-category')"
                        class="filter-select"
                        required
                    >
                        <option value="">Select a category</option>
                        <option
                            v-for="category in categories.filter(
                                (category) => category.id !== 'all'
                            )"
                            :key="category.id"
                            :value="category.id"
                        >
                            {{ category.name }}
                        </option>
                    </select>
                </div>
                <div class="form-row">
                    <div class="soko-field">
                        <label for="rfq-quantity">Quantity</label
                        ><input
                            id="rfq-quantity"
                            v-model.number="rfqQuantity"
                            v-bind="fieldState('rfq-quantity')"
                            type="number"
                            min="1"
                            required
                        />
                    </div>
                    <div class="soko-field">
                        <label for="rfq-price"
                            >Target price per unit (USD, optional)</label
                        ><input
                            id="rfq-price"
                            v-model.number="rfqTargetPrice"
                            v-bind="fieldState('rfq-price')"
                            type="number"
                            step="0.01"
                        />
                    </div>
                </div>
                <div class="form-row">
                    <div class="soko-field">
                        <label for="rfq-destination"
                            >Destination country or port</label
                        ><input
                            id="rfq-destination"
                            v-model="rfqDestination"
                            v-bind="fieldState('rfq-destination')"
                            placeholder="e.g. Mombasa, Kenya"
                            required
                        />
                    </div>
                    <div class="soko-field">
                        <label for="rfq-company">Company name</label
                        ><input
                            id="rfq-company"
                            v-model="rfqCompany"
                            autocomplete="organization"
                        />
                    </div>
                </div>
                <div class="soko-field">
                    <label for="rfq-specifications"
                        >Specifications (optional)</label
                    ><textarea
                        id="rfq-specifications"
                        v-model="rfqSpecs"
                        class="resize-none"
                        rows="4"
                        placeholder="Material, dimensions, certifications and delivery requirements"
                    ></textarea>
                </div>
                <button
                    type="submit"
                    class="button button-navy button-full"
                    :disabled="submitting"
                >
                    {{
                        submitting
                            ? 'Publishing request…'
                            : 'Publish sourcing request'
                    }}
                </button>
            </form>
        </SokoDialog>
    </div>
</template>

<script setup lang="ts">
import { Head, Link, router, usePage } from '@inertiajs/vue3'
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import {
    useCurrency,
    type Currency,
    type MarketRegion,
} from '@/composables/useCurrency'
import { useCatalogNavigation } from '@/composables/useCatalogNavigation'
import { useStorefrontMotion } from '@/composables/useStorefrontMotion'
import SokoIcon from '@/components/SokoIcon.vue'
import SokoDialog from '@/components/SokoDialog.vue'

const props = defineProps<{
    products: any[] | null
    categories: any[]
    currencies: Currency[]
    regions: MarketRegion[]
    rfqs: any[] | null
    activeType: string
    activeFormat: string
    activeRegion: string
    activeCategory: string
    activeCondition: string
    searchQuery: string
    stats: any
    storefrontContent?: Record<string, string>
}>()
const { activeCurrency, formatPrice } = useCurrency(props.currencies)
const { visitCatalog, catalogHref } = useCatalogNavigation()
const page = usePage()
const sortOrder = computed(
    () =>
        new URL(page.url, 'http://soko.local').searchParams.get('sort') ||
        'newest'
)
const productList = computed(() => {
    const products = [...(props.products || [])]
    if (sortOrder.value === 'price_asc')
        products.sort((a, b) => a.price - b.price)
    else if (sortOrder.value === 'price_desc')
        products.sort((a, b) => b.price - a.price)
    else if (sortOrder.value === 'title')
        products.sort((a, b) => a.title.localeCompare(b.title))
    return products
})
const storefrontElement = ref<HTMLElement | null>(null)
const heroElement = ref<HTMLElement | null>(null)
const activeStory = ref(0)
const { motionEnabled } = useStorefrontMotion()
const prefersReducedMotion = ref(false)
const heroVisible = ref(true)
const stories = [
    {
        id: 'source',
        label: 'For a world of opportunity',
        shortLabel: 'Source globally',
        caption: 'Equipment & wholesale',
        icon: 'globe',
        image: '/public/images/sourcing-harbor.jpg',
        headline: 'The world’s supply.',
        accent: 'On your terms.',
        description:
            'From factory floors to your next frontier. Discover equipment, connect with suppliers and build a business without boundaries.',
        format: 'all',
        action: 'Explore the marketplace',
    },
    {
        id: 'digital',
        label: 'For the business you’re becoming',
        shortLabel: 'Build digitally',
        caption: 'Software & business assets',
        icon: 'bolt',
        image: '/public/images/soko-digital.jpg',
        headline: 'Build what’s next.',
        accent: 'Start here.',
        description:
            'The tools behind your next leap. Explore business software, licensed digital assets and solutions ready for your team.',
        format: 'digital',
        action: 'Explore digital products',
    },
    {
        id: 'negotiate',
        label: 'For a better kind of deal',
        shortLabel: 'Make an offer',
        caption: 'Negotiable equipment',
        icon: 'factory',
        image: '/public/images/soko-industry.jpg',
        headline: 'Great equipment.',
        accent: 'Better possibilities.',
        description:
            'A listing is a starting point. Discover equipment, review the details and make an offer that works for your business.',
        format: 'make_offer',
        action: 'Explore negotiable products',
    },
]
const currentStory = computed(() => {
    const story = stories[activeStory.value]
    const c = props.storefrontContent
    return activeStory.value === 0 && c
        ? {
              ...story,
              headline: c.hero_heading ?? story.headline,
              accent: c.hero_accent ?? story.accent,
              description: c.hero_description ?? story.description,
          }
        : story
})
const featuredProduct = computed(() => {
    if (activeStory.value === 1)
        return productList.value.find((product) => product.format === 'digital')
    if (activeStory.value === 2)
        return productList.value.find((product) => product.accepts_offers)
    return (
        productList.value.find((product) => product.format === 'wholesale') ||
        productList.value[0]
    )
})
const sectors = [
    {
        title: 'Industry & machinery',
        caption: 'Built to keep you moving',
        icon: 'factory',
        image: '/public/images/soko-industry.jpg',
        query: { category: 'Industrial & Machinery', format: 'all' },
    },
    {
        title: 'Digital & software',
        caption: 'Tools for your next chapter',
        icon: 'cloud',
        image: '/public/images/soko-digital.jpg',
        query: { category: 'all', format: 'digital' },
    },
    {
        title: 'Agriculture & processing',
        caption: 'Rooted in possibility',
        icon: 'leaf',
        image: '/public/images/soko-agriculture.jpg',
        query: { category: 'Agriculture & Processing', format: 'all' },
    },
]
let revealObserver: IntersectionObserver | undefined
let heroObserver: IntersectionObserver | undefined
let headerObserver: ResizeObserver | undefined
let motionPreference: MediaQueryList | undefined
const syncMotionPreference = () => {
    prefersReducedMotion.value = motionPreference?.matches || false
}
onMounted(() => {
    motionPreference = window.matchMedia('(prefers-reduced-motion: reduce)')
    syncMotionPreference()
    motionPreference.addEventListener('change', syncMotionPreference)
    const header = document.querySelector('.site-header')
    if (header && 'ResizeObserver' in window) {
        headerObserver = new ResizeObserver(() => {
            storefrontElement.value?.style.setProperty(
                '--storefront-header-height',
                `${header.getBoundingClientRect().height}px`
            )
        })
        headerObserver.observe(header)
    }
    if ('IntersectionObserver' in window) {
        revealObserver = new IntersectionObserver(
            (entries) => {
                for (const entry of entries)
                    if (entry.isIntersecting) {
                        entry.target.classList.add('is-revealed')
                        revealObserver?.unobserve(entry.target)
                    }
            },
            { threshold: 0.08 }
        )
        storefrontElement.value?.classList.add('reveal-ready')
        storefrontElement.value
            ?.querySelectorAll('.reveal-item')
            .forEach((element) => revealObserver?.observe(element))
        heroObserver = new IntersectionObserver((entries) => {
            heroVisible.value = entries[0]?.isIntersecting || false
        })
        if (heroElement.value) heroObserver.observe(heroElement.value)
    }
})
onBeforeUnmount(() => {
    revealObserver?.disconnect()
    heroObserver?.disconnect()
    headerObserver?.disconnect()
    motionPreference?.removeEventListener('change', syncMotionPreference)
})

const filterRegion = ref(props.activeRegion || 'global')
watch(
    () => props.activeRegion,
    (region) => {
        filterRegion.value = region || 'global'
    }
)
const filtersOpen = ref(false)
const missingImages = ref(new Set<string>())
const formats = [
    { id: 'all', label: 'All products' },
    { id: 'wholesale', label: 'Wholesale' },
    { id: 'retail', label: 'Ready to buy' },
    { id: 'digital', label: 'Digital' },
    { id: 'make_offer', label: 'Make an offer' },
]
const conditions = [
    { id: 'all', label: 'All conditions' },
    { id: 'new', label: 'New' },
    { id: 'refurbished', label: 'Refurbished' },
    { id: 'open_box', label: 'Open box' },
]
const categoryIcons: Record<string, string> = {
    'Industrial & Machinery': 'factory',
    'Electronics & Telecom': 'chip',
    'Digital Software': 'cloud',
    'Digital Assets': 'layers',
    'Agriculture & Processing': 'leaf',
    'Energy & Utilities': 'bolt',
}
const formatLabels: Record<string, string> = {
    wholesale: 'WHOLESALE',
    retail: 'READY TO BUY',
    digital: 'DIGITAL PRODUCT',
    make_offer: 'OPEN TO OFFERS',
}
const formatLabel = (format: string) => formatLabels[format] || 'MARKETPLACE'
const setFormat = (format: string) => visitCatalog({ format })
const setCategory = (category: string) => visitCatalog({ category })
const setCondition = (condition: string) => visitCatalog({ condition })
const onFilterRegionChange = () => visitCatalog({ region: filterRegion.value })
const resetFilters = () =>
    router.get('/', {}, { preserveState: false, preserveScroll: true })

const offerModalProduct = ref<any | null>(null)
const offerPrice = ref(0)
const offerQty = ref(1)
const offerName = ref('')
const offerEmail = ref('')
const offerMessage = ref('')
const showRfqModal = ref(false)
const rfqTitle = ref('')
const rfqCategory = ref('')
const rfqQuantity = ref(100)
const rfqTargetPrice = ref(50)
const rfqDestination = ref('')
const rfqCompany = ref('')
const rfqSpecs = ref('')
const submitting = ref(false)
const formError = ref('')
const invalidField = ref('')
const clearFormError = () => {
    formError.value = ''
    invalidField.value = ''
}
const fieldState = (id: string) => ({
    'aria-invalid': invalidField.value === id || undefined,
    'aria-describedby': invalidField.value === id ? 'form-feedback' : undefined,
})
const closeOfferModal = () => {
    offerModalProduct.value = null
    clearFormError()
}
const closeRfqModal = () => {
    showRfqModal.value = false
    clearFormError()
}
const openRfqModal = () => {
    if (!(usePage().props.auth as any)?.user) {
        router.visit('/login')
        return
    }
    clearFormError()
    showRfqModal.value = true
}
const openOrderModal = (product: any) => {
    postForm(
        '/cart/items',
        { product_id: product.id, quantity: product.moq || 1 },
        () => {}
    )
}
const openOfferModal = (product: any) => {
    if (!(usePage().props.auth as any)?.user) {
        router.visit('/login')
        return
    }
    clearFormError()
    offerModalProduct.value = product
    offerPrice.value = product.wholesale_price || product.price * 0.85
    offerQty.value = product.moq || 1
}

const validateForm = (event: Event) => {
    clearFormError()
    const form = event.currentTarget as HTMLFormElement
    const invalid = Array.from(
        form.querySelectorAll<
            HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
        >('input, textarea, select')
    ).find((field) => !field.validity.valid)
    if (!invalid) return true
    invalidField.value = invalid.id
    const label =
        invalid.labels?.[0]?.textContent?.trim().toLowerCase() ||
        'a valid value'
    formError.value =
        invalid.id === 'rfq-title-input' && invalid.validity.valueMissing
            ? 'Enter a product title.'
            : invalid.tagName === 'SELECT' && invalid.validity.valueMissing
              ? 'Choose a product category.'
              : invalid.validity.valueMissing
                ? `Enter ${label}.`
                : invalid.validity.typeMismatch
                  ? 'Enter a valid email address.'
                  : `Check ${label} and enter a valid value.`
    invalid.focus()
    return false
}
const postForm = (
    url: string,
    data: Record<string, string | number>,
    close: () => void
) => {
    if (submitting.value) return
    submitting.value = true
    router.post(url, data, {
        preserveScroll: (response) => response.component === 'Store/Index',
        onSuccess: (response) => {
            const error = (response.props.flash as Record<string, string>)
                ?.error
            if (error) formError.value = error
            else close()
        },
        onError: (errors) => {
            formError.value =
                Object.values(errors).join(' ') ||
                'Check your details and try again.'
        },
        onFinish: () => {
            submitting.value = false
        },
    })
}
const submitOfferForm = (event: Event) => {
    if (!offerModalProduct.value || submitting.value || !validateForm(event))
        return
    postForm(
        '/make-offer',
        {
            product_id: offerModalProduct.value.id,
            offered_price: offerPrice.value,
            quantity: offerQty.value,
            name: offerName.value,
            email: offerEmail.value,
            message: offerMessage.value,
        },
        () => {
            offerModalProduct.value = null
        }
    )
}
const submitRfqForm = (event: Event) => {
    if (submitting.value || !validateForm(event)) return
    postForm(
        '/rfqs',
        {
            title: rfqTitle.value,
            category: rfqCategory.value,
            quantity: rfqQuantity.value,
            target_price: rfqTargetPrice.value,
            destination_country: rfqDestination.value,
            buyer_company: rfqCompany.value,
            specifications: rfqSpecs.value,
        },
        () => {
            showRfqModal.value = false
        }
    )
}
</script>
