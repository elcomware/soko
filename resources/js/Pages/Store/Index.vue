<template>
    <div
        ref="storefrontElement"
        class="storefront product-first-storefront"
        :class="{ 'has-showcase': showHero || sideProducts.length }"
    >
        <template v-for="sectionKey in sectionOrder" :key="sectionKey">
            <template v-if="sectionKey === 'showcase'"
                ><section
                    v-if="showHero || sideProducts.length"
                    class="shop-hero site-container"
                    data-storefront-section="showcase"
                    :class="{
                        'hero-without-side': !sideProducts.length,
                        'hero-without-main': !showHero,
                    }"
                    :aria-labelledby="showHero ? 'storefront-title' : undefined"
                    :aria-label="
                        !showHero
                            ? merchandising.hero_side?.title ||
                              'Featured product picks'
                            : undefined
                    "
                >
                    <div
                        v-if="showHero"
                        class="shop-hero-main"
                        role="region"
                        aria-roledescription="carousel"
                        aria-label="Featured products"
                        tabindex="0"
                        @keydown="onHeroKeydown"
                        @touchstart.passive="onHeroTouchStart"
                        @touchend="onHeroTouchEnd"
                        @click.capture="onHeroClick"
                    >
                        <Transition name="product-slide" mode="out-in">
                            <div
                                :key="heroProduct.id"
                                class="shop-hero-slide"
                                role="group"
                                aria-roledescription="slide"
                                :aria-label="`${heroIndex + 1} of ${heroProducts.length}: ${heroProduct.title}`"
                            >
                                <div class="shop-hero-copy">
                                    <span class="shop-kicker"
                                        ><SokoIcon
                                            :name="
                                                hasPromotion(heroProduct)
                                                    ? 'tag'
                                                    : 'store'
                                            "
                                        />{{
                                            hasPromotion(heroProduct)
                                                ? heroProduct.pricing
                                                      .promotion_name
                                                : merchandising.hero?.title ||
                                                  'Discover in the shop'
                                        }}</span
                                    >
                                    <h1 id="storefront-title">
                                        {{ heroProduct.title }}
                                    </h1>
                                    <p class="shop-hero-value">
                                        {{ productTeaser(heroProduct, 88) }}
                                    </p>
                                    <div
                                        v-if="hasPromotion(heroProduct)"
                                        class="shop-deal-price"
                                    >
                                        <del>{{
                                            formatPrice(
                                                heroProduct.pricing
                                                    .original_unit_price
                                            )
                                        }}</del>
                                        <strong>{{
                                            formatPrice(
                                                heroProduct.pricing.unit_price
                                            )
                                        }}</strong>
                                        <span class="sale-badge">{{
                                            savingsLabel(heroProduct)
                                        }}</span>
                                        <small
                                            v-if="
                                                heroProduct.pricing
                                                    .promotion_ends_at
                                            "
                                            >Ends
                                            {{
                                                campaignEnd(
                                                    heroProduct.pricing
                                                        .promotion_ends_at
                                                )
                                            }}</small
                                        >
                                    </div>
                                    <p v-else class="shop-hero-context">
                                        {{ heroProduct.category }}
                                        <span>·</span>
                                        {{ formatLabel(heroProduct.format) }}
                                    </p>
                                    <Link
                                        :href="'/products/' + heroProduct.id"
                                        class="button button-navy shop-discover-button"
                                        >Discover this product
                                        <SokoIcon name="arrow"
                                    /></Link>
                                    <a href="#catalog" class="shop-all-link"
                                        >Browse all products</a
                                    >
                                </div>
                                <Link
                                    :href="'/products/' + heroProduct.id"
                                    class="shop-hero-visual"
                                    :aria-label="'Explore ' + heroProduct.title"
                                >
                                    <img
                                        v-if="
                                            productImage(heroProduct) &&
                                            !missingImages.has(heroProduct.id)
                                        "
                                        :src="productImage(heroProduct)"
                                        :alt="heroProduct.title"
                                        width="720"
                                        height="600"
                                        fetchpriority="high"
                                        @error="
                                            missingImages.add(heroProduct.id)
                                        "
                                    />
                                    <div v-else class="shop-image-placeholder">
                                        <SokoIcon
                                            :name="
                                                heroProduct.type === 'physical'
                                                    ? 'box'
                                                    : 'cloud'
                                            "
                                            :size="96"
                                        /><span>Product image unavailable</span>
                                    </div>
                                    <span class="shop-visual-caption">{{
                                        isIllustrativeImage(heroProduct)
                                            ? 'Illustrative stock image · ' +
                                              heroProduct.vendor_name
                                            : heroProduct.vendor_name
                                    }}</span>
                                </Link>
                            </div>
                        </Transition>
                        <div
                            v-if="heroProducts.length > 1"
                            class="shop-carousel-controls"
                        >
                            <button
                                type="button"
                                aria-label="Previous featured product"
                                @click="moveHero(-1)"
                            >
                                <SokoIcon
                                    name="chevron"
                                    class="carousel-previous"
                                />
                            </button>
                            <span aria-live="polite" aria-atomic="true"
                                >{{ heroIndex + 1 }} /
                                {{ heroProducts.length }}</span
                            >
                            <button
                                type="button"
                                aria-label="Next featured product"
                                @click="moveHero(1)"
                            >
                                <SokoIcon name="chevron" />
                            </button>
                        </div>
                    </div>
                    <div
                        v-if="sideProducts.length"
                        class="shop-side-picks"
                        role="region"
                        :aria-label="
                            merchandising.hero_side?.title ||
                            'More product picks'
                        "
                    >
                        <article
                            v-for="product in sideProducts"
                            :key="product.id"
                            class="shop-side-product-card"
                        >
                            <Link
                                :href="'/products/' + product.id"
                                class="side-product-image"
                                :aria-label="'View ' + product.title"
                            >
                                <img
                                    v-if="
                                        productImage(product) &&
                                        !missingImages.has(product.id)
                                    "
                                    :src="productImage(product)"
                                    :alt="product.title"
                                    @error="missingImages.add(product.id)"
                                />
                                <SokoIcon v-else name="box" :size="48" />
                                <small v-if="isIllustrativeImage(product)"
                                    >Illustrative stock image</small
                                >
                            </Link>
                            <div class="side-product-body">
                                <span
                                    v-for="badge in productBadges(
                                        product,
                                        promotionClock
                                    ).slice(0, 1)"
                                    :key="badge.tone"
                                    class="product-state-badge"
                                    :class="'badge-' + badge.tone"
                                    >{{ badge.label }}</span
                                >
                                <h2>
                                    <Link
                                        :href="'/products/' + product.id"
                                        :title="product.title"
                                        >{{ product.title }}</Link
                                    >
                                </h2>
                                <p>{{ productTeaser(product, 44) }}</p>
                                <ProductCardFacts
                                    :product="product"
                                    :settings="footerSettings"
                                    :format-price="formatPrice"
                                    :now="promotionClock"
                                    :paused="
                                        footerPaused || footerReducedMotion
                                    "
                                />
                                <Link
                                    :href="'/products/' + product.id"
                                    class="shop-pick-action"
                                    >Explore product <SokoIcon name="arrow"
                                /></Link>
                            </div>
                        </article>
                    </div></section
            ></template>
            <div
                v-else-if="
                    sectionKey === 'category_navigation' &&
                    merchandising.category_navigation?.enabled !== false
                "
                class="site-container storefront-sections shop-sections storefront-block"
                :data-storefront-section="sectionKey"
            >
                <section
                    v-if="merchandising.category_navigation?.enabled !== false"
                    id="discover"
                    class="shop-categories"
                    aria-labelledby="category-title"
                >
                    <div class="section-heading">
                        <div>
                            <p class="eyebrow">Find what you need</p>
                            <h2 id="category-title">
                                {{
                                    merchandising.category_navigation?.title ||
                                    'Shop by category'
                                }}
                            </h2>
                        </div>
                    </div>
                    <div class="category-rail-frame">
                        <button
                            type="button"
                            class="category-rail-arrow"
                            aria-label="Previous categories"
                            aria-controls="soko-category-rail"
                            :disabled="!categoryBoundary.previous"
                            @click="scrollCategories(-1)"
                        >
                            <SokoIcon
                                name="chevron"
                                class="carousel-previous"
                            />
                        </button>
                        <div
                            :ref="setCategoryScroller"
                            id="soko-category-rail"
                            @scroll.passive="updateCategoryBoundary"
                            @keydown="onCategoryKeydown"
                            class="shop-category-grid"
                            tabindex="0"
                            aria-label="Browse product categories"
                        >
                            <Link
                                v-for="category in categories.filter(
                                    (item) => item.id !== 'all'
                                )"
                                :key="category.id"
                                :href="catalogHref({ category: category.id })"
                                class="shop-category-tile"
                                :aria-current="
                                    activeCategory === category.id
                                        ? 'page'
                                        : undefined
                                "
                            >
                                <span class="shop-category-symbol"
                                    ><SokoCommerceIllustration
                                        :name="
                                            categoryScenes[category.id] ||
                                            'storefront'
                                        "
                                        :size="72" /></span
                                ><span>{{ category.name }}</span>
                            </Link>
                        </div>
                        <button
                            type="button"
                            class="category-rail-arrow"
                            aria-label="Next categories"
                            aria-controls="soko-category-rail"
                            :disabled="!categoryBoundary.next"
                            @click="scrollCategories(1)"
                        >
                            <SokoIcon name="chevron" />
                        </button>
                    </div>
                </section>
            </div>
            <div
                v-else-if="sectionKey === 'catalog'"
                class="site-container storefront-sections shop-sections storefront-block"
                :data-storefront-section="sectionKey"
            >
                <section
                    id="catalog"
                    class="catalog-section reveal-item"
                    aria-labelledby="catalog-title"
                >
                    <div class="section-heading">
                        <div>
                            <p class="eyebrow">Explore the collection</p>
                            <h2 id="catalog-title">
                                {{
                                    searchQuery
                                        ? `Results for “${searchQuery}”`
                                        : activeCategory !== 'all'
                                          ? activeCategory
                                          : 'Find your next essential'
                                }}
                            </h2>
                        </div>
                        <span class="catalog-count" role="status"
                            >{{ productList.length }}
                            {{
                                productList.length === 1
                                    ? 'product'
                                    : 'products'
                            }}
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
                                <button
                                    class="text-button"
                                    @click="resetFilters"
                                >
                                    Reset all
                                </button>
                            </div>
                            <div class="filter-group">
                                <h4 class="filter-label">Categories</h4>
                                <button
                                    v-for="category in categories"
                                    :key="category.id"
                                    class="category-option"
                                    :aria-pressed="
                                        activeCategory === category.id
                                    "
                                    @click="setCategory(category.id)"
                                >
                                    <SokoIcon
                                        :name="
                                            categoryIcons[category.id] ||
                                            'layers'
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
                                    :aria-pressed="
                                        activeCondition === condition.id
                                    "
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
                                <label
                                    for="sourcing-region"
                                    class="filter-label"
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
                                    Share your requirements and let suppliers
                                    come to you.
                                </p>
                                <button
                                    class="text-button"
                                    @click="openRfqModal"
                                >
                                    Post a sourcing request
                                    <SokoIcon name="arrow" />
                                </button>
                            </div>
                        </aside>
                        <div class="catalog-results">
                            <p
                                v-if="
                                    formError &&
                                    !offerModalProduct &&
                                    !showRfqModal
                                "
                                class="form-error"
                                role="alert"
                            >
                                {{ formError }}
                            </p>
                            <button
                                v-if="
                                    footerConfig.mode === 'rotate' &&
                                    !footerReducedMotion
                                "
                                type="button"
                                class="card-details-toggle"
                                :aria-pressed="footerPaused"
                                @click="footerPaused = !footerPaused"
                            >
                                {{
                                    footerPaused
                                        ? 'Resume rotating details'
                                        : 'Pause rotation · show all details'
                                }}
                            </button>
                            <div class="catalog-toolbar">
                                <div
                                    class="format-tabs"
                                    role="group"
                                    aria-label="Product format"
                                >
                                    <button
                                        v-for="format in formats"
                                        :key="format.id"
                                        :aria-pressed="
                                            activeFormat === format.id
                                        "
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
                                        <option value="title">
                                            Product name
                                        </option>
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
                                    <div class="product-card-media">
                                        <Link
                                            :href="'/products/' + product.id"
                                            class="product-image-link"
                                            :class="{
                                                'image-missing':
                                                    !productImage(product) ||
                                                    missingImages.has(
                                                        product.id
                                                    ),
                                            }"
                                            :aria-label="
                                                'View ' + product.title
                                            "
                                            ><img
                                                v-if="productImage(product)"
                                                :src="productImage(product)"
                                                :alt="product.title"
                                                loading="lazy"
                                                width="480"
                                                height="330"
                                                @error="
                                                    missingImages.add(
                                                        product.id
                                                    )
                                                " /><span
                                                class="image-fallback"
                                                aria-hidden="true"
                                                ><SokoIcon
                                                    :name="
                                                        product.type ===
                                                        'physical'
                                                            ? 'box'
                                                            : 'cloud'
                                                    " /></span
                                        ></Link>
                                        <div
                                            v-if="
                                                productBadges(
                                                    product,
                                                    promotionClock
                                                ).length
                                            "
                                            class="product-card-labels"
                                            aria-label="Product labels"
                                        >
                                            <span
                                                v-for="badge in productBadges(
                                                    product,
                                                    promotionClock
                                                )"
                                                :key="badge.tone"
                                                class="product-state-badge"
                                                :class="'badge-' + badge.tone"
                                                :title="badge.label"
                                                >{{ badge.label }}</span
                                            >
                                        </div>
                                        <span
                                            v-if="isIllustrativeImage(product)"
                                            class="product-stock-disclosure"
                                            >Illustrative stock image</span
                                        >
                                    </div>
                                    <div
                                        class="product-body"
                                        :class="{
                                            'has-boundary-badges':
                                                productBadges(
                                                    product,
                                                    promotionClock
                                                ).length,
                                        }"
                                    >
                                        <p class="product-category">
                                            {{ product.category }}
                                        </p>
                                        <h3>
                                            <Link
                                                :href="
                                                    '/products/' + product.id
                                                "
                                                :title="product.title"
                                                >{{ product.title }}</Link
                                            >
                                        </h3>
                                        <p class="product-value">
                                            {{ productTeaser(product, 72) }}
                                        </p>
                                        <div class="product-discovery-meta">
                                            <span>{{
                                                product.vendor_name
                                            }}</span
                                            ><span
                                                v-if="
                                                    product.condition !== 'new'
                                                "
                                                >{{
                                                    product.condition ===
                                                    'open_box'
                                                        ? 'Open box'
                                                        : 'Refurbished'
                                                }}</span
                                            >
                                        </div>
                                        <ProductCardFacts
                                            :product="product"
                                            :settings="footerSettings"
                                            :format-price="formatPrice"
                                            :now="promotionClock"
                                            :paused="
                                                footerPaused ||
                                                footerReducedMotion
                                            "
                                        />
                                        <div class="product-discovery-actions">
                                            <Link
                                                :href="
                                                    '/products/' + product.id
                                                "
                                                class="product-discover-link"
                                                >Explore product
                                                <SokoIcon name="arrow"
                                            /></Link>
                                        </div>
                                    </div>
                                </article>
                            </div>
                        </div>
                    </div>
                </section>
            </div>
            <div
                v-else-if="
                    sectionKey === 'discovery' &&
                    digitalProducts.length > 0 &&
                    merchandising.discovery?.enabled !== false
                "
                class="site-container storefront-sections shop-sections storefront-block"
                :data-storefront-section="sectionKey"
            >
                <section
                    v-if="
                        digitalProducts.length &&
                        merchandising.discovery?.enabled !== false
                    "
                    class="shop-digital-shelf"
                    aria-labelledby="digital-shelf-title"
                >
                    <div class="section-heading">
                        <div>
                            <p class="eyebrow">Software & digital products</p>
                            <h2 id="digital-shelf-title">
                                {{
                                    merchandising.discovery?.title ||
                                    'Equip your digital workspace'
                                }}
                            </h2>
                        </div>
                        <Link
                            :href="catalogHref({ format: 'digital' })"
                            class="section-explore"
                            >Shop digital <SokoIcon name="arrow"
                        /></Link>
                    </div>
                    <div class="shop-digital-grid">
                        <Link
                            v-for="product in digitalProducts"
                            :key="product.id"
                            :href="'/products/' + product.id"
                            class="shop-digital-card"
                            ><img
                                v-if="
                                    productImage(product) &&
                                    !missingImages.has(product.id)
                                "
                                :src="productImage(product)"
                                :alt="product.title"
                                width="400"
                                height="260"
                                loading="lazy"
                                @error="
                                    missingImages.add(product.id)
                                " /><SokoIcon v-else name="cloud" :size="64" />
                            <div>
                                <span>{{ product.category }}</span>
                                <small
                                    v-if="isIllustrativeImage(product)"
                                    class="shop-stock-note"
                                    >Illustrative stock image</small
                                >
                                <h3>{{ product.title }}</h3>
                                <p>{{ productTeaser(product, 72) }}</p>
                                <strong
                                    >Explore product <SokoIcon name="arrow"
                                /></strong></div
                        ></Link>
                    </div>
                </section>
            </div>
            <div
                v-else-if="sectionKey === 'sourcing'"
                class="site-container storefront-sections shop-sections storefront-block"
                :data-storefront-section="sectionKey"
            >
                <section
                    id="rfq"
                    class="shop-sourcing-help"
                    aria-labelledby="sourcing-title"
                >
                    <div class="shop-sourcing-icon">
                        <SokoIcon name="search" :size="32" />
                    </div>
                    <div>
                        <h2 id="sourcing-title">
                            Can’t find the right product?
                        </h2>
                        <p>
                            Request a specific item, quantity or specification
                            from suppliers.
                        </p>
                    </div>
                    <button class="button button-outline" @click="openRfqModal">
                        Request a product <SokoIcon name="arrow" /></button
                    ><Link
                        href="/dashboard?section=rfqs"
                        class="shop-supplier-link"
                        >Supplier? View sourcing requests</Link
                    >
                </section>
            </div>
        </template>

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
import { Link, router, usePage } from '@inertiajs/vue3'
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { normalizeSectionOrder } from '@/design/sectionOrder'
import type { ComponentPublicInstance } from 'vue'
import { categoryRailState, categoryRailTarget } from '@/design/categoryRail'
import {
    useCurrency,
    type Currency,
    type MarketRegion,
} from '@/composables/useCurrency'
import { useCatalogNavigation } from '@/composables/useCatalogNavigation'
import SokoIcon from '@/components/SokoIcon.vue'
import SokoCommerceIllustration from '@/components/SokoCommerceIllustration.vue'
import ProductCardFacts from '@/components/ProductCardFacts.vue'
import { normalizeCardFooter } from '@/design/cardFacts'
import SokoDialog from '@/components/SokoDialog.vue'
import {
    productTeaser,
    productBadges,
    carouselIndex,
    swipeDirection,
    productImage,
    selectHeroProducts,
    orderResolvedHero,
    isIllustrativeImage,
    hasPromotion as isActivePromotion,
    savingsLabel,
    campaignEnd,
} from '@/design/merchandising'

const promotionClock = ref(Date.now())
const footerPaused = ref(false)
const footerReducedMotion = ref(false)
let footerMotionQuery: MediaQueryList | undefined
const syncFooterMotion = () => {
    footerReducedMotion.value = Boolean(footerMotionQuery?.matches)
}
onMounted(() => {
    footerMotionQuery = window.matchMedia('(prefers-reduced-motion: reduce)')
    syncFooterMotion()
    footerMotionQuery.addEventListener('change', syncFooterMotion)
})
onBeforeUnmount(() =>
    footerMotionQuery?.removeEventListener('change', syncFooterMotion)
)

let promotionTimer: ReturnType<typeof setInterval> | undefined
onMounted(() => {
    promotionTimer = setInterval(() => {
        promotionClock.value = Date.now()
    }, 1000)
})
onBeforeUnmount(() => {
    if (promotionTimer) clearInterval(promotionTimer)
})
const hasPromotion = (product: Record<string, any>) =>
    isActivePromotion(product, promotionClock.value)
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
    merchandising?: Record<string, any>
}>()
const { activeCurrency, formatPrice } = useCurrency(props.currencies)
const { visitCatalog, catalogHref } = useCatalogNavigation()
const page = usePage()
const footerSettings = computed(
    () => (page.props.brandSettings || {}) as Record<string, any>
)
const footerConfig = computed(() => normalizeCardFooter(footerSettings.value))
const sortOrder = computed(
    () =>
        new URL(page.url, 'http://soko.local').searchParams.get('sort') ||
        'newest'
)
const productList = computed(() => {
    const products = [...(props.products || [])]
    if (sortOrder.value === 'price_asc')
        products.sort(
            (a, b) =>
                (hasPromotion(a)
                    ? a.pricing.unit_price
                    : (a.pricing?.original_unit_price ?? a.price)) -
                (hasPromotion(b)
                    ? b.pricing.unit_price
                    : (b.pricing?.original_unit_price ?? b.price))
        )
    else if (sortOrder.value === 'price_desc')
        products.sort(
            (a, b) =>
                (hasPromotion(b)
                    ? b.pricing.unit_price
                    : (b.pricing?.original_unit_price ?? b.price)) -
                (hasPromotion(a)
                    ? a.pricing.unit_price
                    : (a.pricing?.original_unit_price ?? a.price))
        )
    else if (sortOrder.value === 'title')
        products.sort((a, b) => a.title.localeCompare(b.title))
    return products
})
const storefrontElement = ref<HTMLElement | null>(null)
const categoryScroller = ref<HTMLElement | null>(null)
const merchandising = computed(() => props.merchandising || {})
const sectionOrder = computed(() =>
    normalizeSectionOrder(merchandising.value.section_order)
)
function setCategoryScroller(
    element: Element | ComponentPublicInstance | null
) {
    categoryScroller.value = element as HTMLElement | null
}
const categoryBoundary = ref({ previous: false, next: false, maximum: 0 })
let categoryObserver: ResizeObserver | undefined
function updateCategoryBoundary() {
    const rail = categoryScroller.value
    categoryBoundary.value = rail
        ? categoryRailState(rail.scrollLeft, rail.scrollWidth, rail.clientWidth)
        : { previous: false, next: false, maximum: 0 }
}
function scrollCategories(direction: number) {
    const rail = categoryScroller.value
    if (!rail) return
    const target = categoryRailTarget(
        direction < 0 ? 'ArrowLeft' : 'ArrowRight',
        rail.scrollLeft,
        rail.scrollWidth,
        rail.clientWidth
    )
    rail.scrollTo({
        left: target || 0,
        behavior: footerReducedMotion.value ? 'auto' : 'smooth',
    })
}
function onCategoryKeydown(event: KeyboardEvent) {
    if (event.target !== event.currentTarget) return
    const rail = categoryScroller.value
    if (!rail) return
    const target = categoryRailTarget(
        event.key,
        rail.scrollLeft,
        rail.scrollWidth,
        rail.clientWidth
    )
    if (target === null) return
    event.preventDefault()
    rail.scrollTo({
        left: target,
        behavior: footerReducedMotion.value ? 'auto' : 'smooth',
    })
}
watch(categoryScroller, async (rail, _previous, onCleanup) => {
    let active = true
    onCleanup(() => {
        active = false
        categoryObserver?.disconnect()
    })
    categoryObserver?.disconnect()
    await nextTick()
    if (!active) return
    if (rail && typeof ResizeObserver !== 'undefined') {
        categoryObserver = new ResizeObserver(updateCategoryBoundary)
        categoryObserver.observe(rail)
    }
    updateCategoryBoundary()
})
watch(
    () => props.categories,
    async () => {
        await nextTick()
        updateCategoryBoundary()
    }
)
onMounted(() => {
    window.addEventListener('resize', updateCategoryBoundary)
    updateCategoryBoundary()
})
onBeforeUnmount(() => {
    categoryObserver?.disconnect()
    window.removeEventListener('resize', updateCategoryBoundary)
})
const isCatalogFiltered = computed(
    () =>
        Boolean(props.searchQuery) ||
        [
            'activeType',
            'activeCategory',
            'activeFormat',
            'activeCondition',
            'activeRegion',
        ].some(
            (key) =>
                !['all', 'global', '', undefined].includes((props as any)[key])
        )
)
const missingImages = ref(new Set<string>())
const heroProducts = computed(() => {
    const placement = merchandising.value.hero
    if (placement?.enabled === false) return []
    if (placement?.products)
        return orderResolvedHero(
            placement.products,
            placement.strategy || 'manual',
            missingImages.value
        )
    return selectHeroProducts(
        productList.value,
        missingImages.value,
        promotionClock.value
    )
})
const showHero = computed(() => heroProducts.value.length > 0)
const heroIndex = ref(0)
function moveHero(delta: number) {
    heroIndex.value = carouselIndex(
        heroIndex.value,
        delta,
        heroProducts.value.length
    )
}
function onHeroKeydown(event: KeyboardEvent) {
    if (event.key === 'ArrowLeft' || event.key === 'ArrowRight') {
        event.preventDefault()
        if ((event.target as HTMLElement)?.closest('.shop-hero-slide'))
            (event.currentTarget as HTMLElement).focus()
        moveHero(event.key === 'ArrowRight' ? 1 : -1)
    } else if (event.key === 'Home' || event.key === 'End') {
        event.preventDefault()
        heroIndex.value =
            event.key === 'Home' ? 0 : heroProducts.value.length - 1
    }
}
let touchStart: { x: number; y: number } | null = null
let suppressClickUntil = 0
function onHeroTouchStart(event: TouchEvent) {
    touchStart =
        event.touches.length === 1
            ? { x: event.touches[0].clientX, y: event.touches[0].clientY }
            : null
}
function onHeroTouchEnd(event: TouchEvent) {
    if (!touchStart || !event.changedTouches.length) return
    const delta = swipeDirection(
        event.changedTouches[0].clientX - touchStart.x,
        event.changedTouches[0].clientY - touchStart.y
    )
    touchStart = null
    if (delta && heroProducts.value.length > 1) {
        if (event.cancelable) event.preventDefault()
        suppressClickUntil = Date.now() + 400
        moveHero(delta)
    }
}
function onHeroClick(event: MouseEvent) {
    if (Date.now() < suppressClickUntil) {
        event.preventDefault()
        event.stopPropagation()
    }
}

watch(heroProducts, () => {
    if (heroIndex.value >= heroProducts.value.length) heroIndex.value = 0
})
const heroProduct = computed(() => heroProducts.value[heroIndex.value])
const sideProducts = computed(() =>
    merchandising.value.hero_side?.enabled === false
        ? []
        : merchandising.value.hero_side?.products ||
          productList.value
              .filter(
                  (product) =>
                      !heroProducts.value.some(
                          (hero: any) => hero.id === product.id
                      )
              )
              .slice(0, 2)
)
const digitalProducts = computed(() =>
    merchandising.value.discovery?.enabled === false
        ? []
        : merchandising.value.discovery?.products ||
          productList.value
              .filter(
                  (product) =>
                      product.type !== 'physical' ||
                      product.format === 'digital'
              )
              .slice(0, 3)
)

const filterRegion = ref(props.activeRegion || 'global')
watch(
    () => props.activeRegion,
    (region) => {
        filterRegion.value = region || 'global'
    }
)
const filtersOpen = ref(false)
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
const categoryScenes: Record<string, string> = {
    'Industrial & Machinery': 'warehouse',
    'Electronics & Telecom': 'analytics',
    'Digital Software': 'digital',
    'Digital Assets': 'customization',
    'Agriculture & Processing': 'agriculture',
    'Energy & Utilities': 'energy',
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
    router.get('/', {}, { preserveState: false, preserveScroll: false })

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
