<template>
  <div class="space-y-16 pb-20">
    <!-- Hero Section -->
    <section class="relative overflow-hidden bg-gradient-to-b from-slate-900 via-slate-950 to-slate-950 border-b border-slate-800/80 pt-10 pb-20">
      <!-- Ambient Glow Effects -->
      <div class="absolute top-0 left-1/4 w-96 h-96 bg-amber-500/10 rounded-full blur-3xl pointer-events-none"></div>
      <div class="absolute top-10 right-1/4 w-96 h-96 bg-cyan-500/10 rounded-full blur-3xl pointer-events-none"></div>

      <div class="max-w-7xl mx-auto px-4 sm:px-6 relative z-10">
        <!-- Pill Tag -->
        <div class="flex items-center justify-center">
          <div class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-slate-800/80 border border-slate-700/80 text-xs font-semibold text-slate-300 shadow-sm backdrop-blur">
            <span class="flex h-2 w-2 rounded-full bg-amber-400"></span>
            <span>Enterprise B2B Wholesale & Automated Digital Goods</span>
            <span class="text-slate-500">·</span>
            <span class="text-cyan-400 font-mono">v1.18 Ready</span>
          </div>
        </div>

        <!-- Main Headline -->
        <div class="text-center max-w-4xl mx-auto mt-6 space-y-4">
          <h1 class="text-4xl sm:text-6xl font-black tracking-tight text-white leading-tight">
            The Global Enterprise Commerce &
            <span class="bg-gradient-to-r from-amber-400 via-yellow-300 to-amber-500 bg-clip-text text-transparent">Digital Marketplace</span>
          </h1>
          <p class="text-base sm:text-xl text-slate-400 font-normal leading-relaxed max-w-2xl mx-auto">
            Trade industrial machinery and physical commodities with B2B escrow, or procure enterprise software licenses, AI models, and digital assets with instant zero-latency delivery.
          </p>
        </div>

        <!-- Hero Search Box -->
        <div class="max-w-3xl mx-auto mt-8">
          <div class="p-2 rounded-2xl bg-slate-900/90 border border-slate-700/90 shadow-2xl backdrop-blur flex flex-col sm:flex-row gap-2">
            <!-- Filter Type Picker -->
            <div class="flex items-center gap-1 bg-slate-950/60 p-1 rounded-xl border border-slate-800 shrink-0">
              <button
                @click="setTypeFilter('all')"
                :class="[
                  'px-3 py-2 rounded-lg text-xs font-bold transition-all',
                  activeType === 'all' ? 'bg-amber-500 text-slate-950 shadow-md' : 'text-slate-400 hover:text-white'
                ]"
              >
                All Goods
              </button>
              <button
                @click="setTypeFilter('digital')"
                :class="[
                  'px-3 py-2 rounded-lg text-xs font-bold transition-all flex items-center gap-1',
                  activeType === 'digital' ? 'bg-cyan-500 text-slate-950 shadow-md' : 'text-slate-400 hover:text-cyan-400'
                ]"
              >
                ⚡ Digital
              </button>
              <button
                @click="setTypeFilter('physical')"
                :class="[
                  'px-3 py-2 rounded-lg text-xs font-bold transition-all flex items-center gap-1',
                  activeType === 'physical' ? 'bg-amber-500 text-slate-950 shadow-md' : 'text-slate-400 hover:text-amber-400'
                ]"
              >
                🚢 Physical
              </button>
            </div>

            <!-- Search Field -->
            <div class="flex-1 flex items-center px-3 bg-slate-950/40 rounded-xl border border-slate-800">
              <svg class="w-5 h-5 text-slate-500 mr-2 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <input
                v-model="searchQueryInput"
                @keydown.enter="triggerSearch"
                type="text"
                placeholder="Search products by SKU, name, or category (e.g., ERP, Laser, Fiber)..."
                class="w-full bg-transparent text-sm text-white placeholder-slate-500 outline-none py-2.5"
              />
            </div>

            <button
              @click="triggerSearch"
              class="px-6 py-2.5 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-black text-sm shadow-lg shadow-amber-500/20 transition-all shrink-0 flex items-center justify-center gap-1.5"
            >
              <span>Explore</span>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M14 5l7 7m0 0l-7 7m7-7H3" />
              </svg>
            </button>
          </div>

          <!-- Quick Tags -->
          <div class="flex items-center gap-2 mt-3 text-xs text-slate-500 flex-wrap justify-center">
            <span class="font-medium text-slate-400">Trending Sourcing:</span>
            <button @click="quickSearch('ERP')" class="hover:text-amber-400 bg-slate-900/60 px-2.5 py-1 rounded-md border border-slate-800 transition-colors">Enterprise ERP</button>
            <button @click="quickSearch('AI')" class="hover:text-cyan-400 bg-slate-900/60 px-2.5 py-1 rounded-md border border-slate-800 transition-colors">AI Routing Models</button>
            <button @click="quickSearch('Laser')" class="hover:text-amber-400 bg-slate-900/60 px-2.5 py-1 rounded-md border border-slate-800 transition-colors">3000W CNC Fiber Laser</button>
            <button @click="quickSearch('Ethernet')" class="hover:text-amber-400 bg-slate-900/60 px-2.5 py-1 rounded-md border border-slate-800 transition-colors">Cat6A Spools</button>
          </div>
        </div>

        <!-- Metric KPI Cards -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 max-w-5xl mx-auto mt-14">
          <div class="p-4 rounded-2xl bg-slate-900/80 border border-slate-800/80 backdrop-blur flex items-center gap-3.5">
            <div class="w-10 h-10 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-400 shrink-0">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <div class="text-xl font-black text-white">${{ (stats.total_gmv || 48200).toLocaleString() }}</div>
              <div class="text-[11px] text-slate-400">Total GMV Traded</div>
            </div>
          </div>

          <div class="p-4 rounded-2xl bg-slate-900/80 border border-slate-800/80 backdrop-blur flex items-center gap-3.5">
            <div class="w-10 h-10 rounded-xl bg-cyan-500/10 border border-cyan-500/20 flex items-center justify-center text-cyan-400 shrink-0">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <div>
              <div class="text-xl font-black text-white">Instant Fulfillment</div>
              <div class="text-[11px] text-slate-400">Zero-Latency Digital SLA</div>
            </div>
          </div>

          <div class="p-4 rounded-2xl bg-slate-900/80 border border-slate-800/80 backdrop-blur flex items-center gap-3.5">
            <div class="w-10 h-10 rounded-xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400 shrink-0">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
            <div>
              <div class="text-xl font-black text-white">100% Escrow</div>
              <div class="text-[11px] text-slate-400">Trade Assurance Protected</div>
            </div>
          </div>

          <div class="p-4 rounded-2xl bg-slate-900/80 border border-slate-800/80 backdrop-blur flex items-center gap-3.5">
            <div class="w-10 h-10 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400 shrink-0">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <div class="text-xl font-black text-white">180+ Countries</div>
              <div class="text-[11px] text-slate-400">Global Multimodal Freight</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Category Carousel / Quick Navigation Strip -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h2 class="text-xl font-bold text-white tracking-tight">Enterprise Sourcing Sectors</h2>
          <p class="text-xs text-slate-400">Browse verified industrial suppliers and certified digital software developers</p>
        </div>
        <div class="text-xs font-semibold text-amber-400 hover:underline cursor-pointer" @click="resetFilters">
          View All Categories &rarr;
        </div>
      </div>

      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
        <div
          v-for="cat in categories"
          :key="cat.id"
          @click="setCategoryFilter(cat.id)"
          :class="[
            'p-4 rounded-xl border text-left cursor-pointer transition-all hover:scale-[1.02]',
            activeCategory === cat.id
              ? 'bg-amber-500/10 border-amber-500/40 text-amber-300'
              : 'bg-slate-900/60 border-slate-800 text-slate-300 hover:border-slate-700 hover:bg-slate-900'
          ]"
        >
          <div class="w-8 h-8 rounded-lg mb-3 flex items-center justify-center font-bold text-lg"
               :class="cat.type === 'digital' ? 'bg-cyan-500/20 text-cyan-400' : 'bg-amber-500/20 text-amber-400'">
            <span v-if="cat.type === 'digital'">⚡</span>
            <span v-else>🚢</span>
          </div>
          <div class="font-bold text-xs truncate">{{ cat.name }}</div>
          <div class="text-[10px] text-slate-500 mt-1 uppercase font-semibold">
            {{ cat.type === 'digital' ? 'Instant Access' : 'B2B Freight' }}
          </div>
        </div>
      </div>
    </section>

    <!-- Main Product Catalog Grid -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6">
      <!-- Section Header & Filter Tabs -->
      <div class="flex flex-col md:flex-row md:items-end justify-between gap-4 pb-6 border-b border-slate-800">
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-2xl font-black text-white tracking-tight">Enterprise Marketplace Showcase</h2>
            <span class="text-xs px-2.5 py-0.5 rounded-full bg-slate-800 border border-slate-700 text-slate-300 font-mono">
              {{ products.length }} available
            </span>
          </div>
          <p class="text-xs text-slate-400 mt-1">
            Physical commodities with wholesale volume tiers & instant downloadable digital solutions
          </p>
        </div>

        <!-- Filter Segment -->
        <div class="flex items-center gap-2 flex-wrap">
          <button
            @click="setTypeFilter('all')"
            :class="[
              'px-3.5 py-1.5 rounded-xl text-xs font-bold transition-all',
              activeType === 'all' ? 'bg-slate-800 text-white border border-slate-600' : 'text-slate-400 hover:text-white'
            ]"
          >
            All Items ({{ products.length }})
          </button>
          <button
            @click="setTypeFilter('digital')"
            :class="[
              'px-3.5 py-1.5 rounded-xl text-xs font-bold transition-all flex items-center gap-1.5',
              activeType === 'digital' ? 'bg-cyan-500/20 text-cyan-300 border border-cyan-500/40' : 'text-slate-400 hover:text-cyan-400'
            ]"
          >
            <span class="w-2 h-2 rounded-full bg-cyan-400"></span>
            <span>Digital Goods & SaaS</span>
          </button>
          <button
            @click="setTypeFilter('physical')"
            :class="[
              'px-3.5 py-1.5 rounded-xl text-xs font-bold transition-all flex items-center gap-1.5',
              activeType === 'physical' ? 'bg-amber-500/20 text-amber-300 border border-amber-500/40' : 'text-slate-400 hover:text-amber-400'
            ]"
          >
            <span class="w-2 h-2 rounded-full bg-amber-400"></span>
            <span>Physical Freight & Machinery</span>
          </button>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="products.length === 0" class="text-center py-20 bg-slate-900/40 rounded-3xl border border-slate-800 mt-6">
        <div class="w-16 h-16 rounded-2xl bg-slate-800 flex items-center justify-center mx-auto text-slate-500 text-2xl mb-4">
          🔍
        </div>
        <h3 class="text-lg font-bold text-white">No products found</h3>
        <p class="text-xs text-slate-400 mt-1 max-w-sm mx-auto">
          Try clearing your search query or selecting a different category filter.
        </p>
        <button
          @click="resetFilters"
          class="mt-4 px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs font-semibold text-slate-200 transition-colors"
        >
          Reset All Filters
        </button>
      </div>

      <!-- Product Cards Grid -->
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mt-8">
        <div
          v-for="product in products"
          :key="product.id"
          class="group rounded-2xl bg-slate-900/90 border border-slate-800/80 overflow-hidden hover:border-slate-700 transition-all hover:shadow-2xl hover:shadow-slate-950 flex flex-col justify-between"
        >
          <!-- Card Top: Image + Badges -->
          <div>
            <div class="relative h-48 w-full overflow-hidden bg-slate-950">
              <img
                :src="product.image_url"
                :alt="product.title"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-slate-950/80 via-transparent to-transparent"></div>

              <!-- Product Type Badge -->
              <div class="absolute top-3 left-3">
                <span
                  v-if="product.type !== 'physical'"
                  class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-wider bg-cyan-950/90 text-cyan-300 border border-cyan-500/40 backdrop-blur"
                >
                  <span class="w-1.5 h-1.5 rounded-full bg-cyan-400 animate-pulse"></span>
                  Instant Digital
                </span>
                <span
                  v-else
                  class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-wider bg-amber-950/90 text-amber-300 border border-amber-500/40 backdrop-blur"
                >
                  <span>🚢 Freight · MOQ {{ product.moq }}</span>
                </span>
              </div>

              <!-- Rating Badge -->
              <div class="absolute top-3 right-3 bg-slate-950/80 border border-slate-800 px-2 py-0.5 rounded-lg text-[10px] font-bold text-amber-400 flex items-center gap-1 backdrop-blur">
                <span>★</span>
                <span>{{ product.rating }}</span>
                <span class="text-slate-500">({{ product.reviews_count }})</span>
              </div>

              <!-- Category Pill -->
              <div class="absolute bottom-2 left-3 text-[11px] font-medium text-slate-300 truncate max-w-[80%]">
                {{ product.category }} · {{ product.sub_category }}
              </div>
            </div>

            <!-- Card Body -->
            <div class="p-4 space-y-3">
              <!-- Title -->
              <h3 class="font-bold text-sm text-slate-100 group-hover:text-amber-400 transition-colors line-clamp-2 leading-snug">
                {{ product.title }}
              </h3>

              <!-- Supplier Information -->
              <div class="flex items-center gap-1.5 text-[11px] text-slate-400">
                <svg class="w-3.5 h-3.5 text-amber-500 shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span class="truncate font-medium">{{ product.vendor_name }}</span>
                <span class="text-slate-600">·</span>
                <span class="shrink-0 text-slate-500">{{ product.vendor_country }}</span>
              </div>

              <!-- Specifications Bar -->
              <div class="p-2.5 rounded-xl bg-slate-950/60 border border-slate-800/80 text-[11px] space-y-1">
                <div v-if="product.type !== 'physical'" class="flex justify-between items-center text-slate-400">
                  <span>Delivery Format:</span>
                  <span class="text-cyan-400 font-mono text-[10px]">{{ product.digital_file_size || 'Instant Download' }}</span>
                </div>
                <div v-else class="flex justify-between items-center text-slate-400">
                  <span>Lead Time:</span>
                  <span class="text-amber-400 font-medium">{{ product.lead_time_days }} days dispatch</span>
                </div>
                <div class="flex justify-between items-center text-slate-400">
                  <span>SKU:</span>
                  <span class="font-mono text-slate-400 text-[10px]">{{ product.sku }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Card Footer: Pricing & Action Button -->
          <div class="p-4 pt-0 space-y-3">
            <div class="pt-3 border-t border-slate-800/80 flex items-baseline justify-between">
              <div>
                <div class="text-xs text-slate-500">Unit Price</div>
                <div class="text-lg font-black text-white">
                  ${{ product.price.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}
                </div>
              </div>
              <div v-if="product.wholesale_price && product.wholesale_price < product.price" class="text-right">
                <div class="text-[10px] text-amber-400 font-semibold uppercase">Wholesale Tier</div>
                <div class="text-xs font-bold text-amber-300">
                  ${{ product.wholesale_price.toLocaleString() }} <span class="text-[10px] text-slate-400">/ 10+</span>
                </div>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="grid grid-cols-2 gap-2">
              <Link
                :href="'/products/' + product.id"
                class="px-3 py-2 rounded-xl text-xs font-semibold bg-slate-800 hover:bg-slate-700 text-slate-200 text-center transition-colors"
              >
                Specifications
              </Link>
              <button
                @click="openOrderModal(product)"
                :class="[
                  'px-3 py-2 rounded-xl text-xs font-bold text-slate-950 text-center shadow-md transition-all flex items-center justify-center gap-1',
                  product.type !== 'physical'
                    ? 'bg-gradient-to-r from-cyan-400 to-cyan-500 hover:from-cyan-300 hover:to-cyan-400 shadow-cyan-500/20'
                    : 'bg-gradient-to-r from-amber-400 to-amber-500 hover:from-amber-300 hover:to-amber-400 shadow-amber-500/20'
                ]"
              >
                <span v-if="product.type !== 'physical'">⚡ Procure</span>
                <span v-else>📦 Order</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Enterprise B2B Features Grid -->
    <section id="assurance" class="max-w-7xl mx-auto px-4 sm:px-6">
      <div class="rounded-3xl bg-gradient-to-br from-slate-900 to-slate-950 border border-slate-800 p-8 sm:p-12 shadow-2xl relative overflow-hidden">
        <div class="absolute -right-20 -bottom-20 w-80 h-80 bg-amber-500/10 rounded-full blur-3xl pointer-events-none"></div>

        <div class="max-w-2xl">
          <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-amber-500/10 text-amber-400 text-xs font-bold border border-amber-500/20 mb-4">
            Alibaba Escrow + Amazon Instant Delivery Architecture
          </div>
          <h2 class="text-3xl font-black text-white tracking-tight">
            Built for High-Volume B2B & Zero-Friction Digital Licensing
          </h2>
          <p class="text-sm text-slate-400 mt-2 leading-relaxed">
            Whether you are ordering 500 spools of datacenter fiber cable, importing heavy CNC manufacturing equipment, or deploying enterprise ERP licenses to 1,000 employees globally, Soko streamlines procurement under a single unified platform.
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-10">
          <div class="p-5 rounded-2xl bg-slate-950/80 border border-slate-800 space-y-2.5">
            <div class="w-10 h-10 rounded-xl bg-amber-500/10 text-amber-400 flex items-center justify-center font-bold text-lg">
              🛡️
            </div>
            <h4 class="font-bold text-white text-sm">Trade Assurance & Escrow</h4>
            <p class="text-xs text-slate-400 leading-relaxed">
              Funds are held safely in escrow until maritime/air freight bills of lading are verified, or digital licenses are authenticated.
            </p>
          </div>

          <div class="p-5 rounded-2xl bg-slate-950/80 border border-slate-800 space-y-2.5">
            <div class="w-10 h-10 rounded-xl bg-cyan-500/10 text-cyan-400 flex items-center justify-center font-bold text-lg">
              ⚡
            </div>
            <h4 class="font-bold text-white text-sm">Automated Digital Fulfillment</h4>
            <p class="text-xs text-slate-400 leading-relaxed">
              Instant generation of software license keys, CDN download tokens, and API credentials directly upon checkout.
            </p>
          </div>

          <div class="p-5 rounded-2xl bg-slate-950/80 border border-slate-800 space-y-2.5">
            <div class="w-10 h-10 rounded-xl bg-emerald-500/10 text-emerald-400 flex items-center justify-center font-bold text-lg">
              📊
            </div>
            <h4 class="font-bold text-white text-sm">Tier-1 Merchant Dashboard</h4>
            <p class="text-xs text-slate-400 leading-relaxed">
              Suppliers and system admins manage listings, audit orders, monitor digital download telemetry, and manage enterprise payouts.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Request for Quotation (RFQ) Section -->
    <section id="rfq" class="max-w-7xl mx-auto px-4 sm:px-6">
      <div class="rounded-3xl bg-slate-900 border border-slate-800 p-8 flex flex-col lg:flex-row items-center justify-between gap-8">
        <div class="space-y-2 max-w-xl">
          <div class="text-amber-400 text-xs font-bold uppercase tracking-wider">Enterprise RFQ Portal</div>
          <h3 class="text-2xl font-black text-white">Can't find the exact product or custom digital software specification?</h3>
          <p class="text-xs text-slate-400 leading-relaxed">
            Submit a Request for Quotation (RFQ). Over 50,000 verified industrial suppliers and certified software engineers will respond within 24 hours with tailored pricing and delivery schedules.
          </p>
        </div>

        <div class="flex items-center gap-3 shrink-0">
          <Link
            href="/register"
            class="px-6 py-3 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs shadow-lg shadow-amber-500/20 transition-all"
          >
            Submit an Enterprise RFQ
          </Link>
          <Link
            href="/login"
            class="px-5 py-3 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 font-semibold text-xs transition-colors"
          >
            Access Demo Accounts
          </Link>
        </div>
      </div>
    </section>

    <!-- Interactive Instant Order / Procurement Modal -->
    <div
      v-if="activeModalProduct"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm"
    >
      <div class="bg-slate-900 border border-slate-700 rounded-3xl max-w-lg w-full p-6 shadow-2xl space-y-6 relative">
        <button
          @click="activeModalProduct = null"
          class="absolute top-5 right-5 text-slate-400 hover:text-white text-lg font-bold"
        >
          &times;
        </button>

        <!-- Modal Header -->
        <div class="flex items-center gap-3">
          <img
            :src="activeModalProduct.image_url"
            class="w-14 h-14 rounded-xl object-cover border border-slate-700"
            alt="Product"
          />
          <div>
            <div class="text-xs font-bold text-amber-400 uppercase tracking-wider">
              {{ activeModalProduct.type !== 'physical' ? '⚡ Instant Digital Fulfillment' : '🚢 Physical Freight Delivery' }}
            </div>
            <h3 class="text-base font-bold text-white line-clamp-1">{{ activeModalProduct.title }}</h3>
            <div class="text-xs text-slate-400">Vendor: {{ activeModalProduct.vendor_name }} ({{ activeModalProduct.vendor_country }})</div>
          </div>
        </div>

        <!-- Order Configurator Form -->
        <form @submit.prevent="submitCheckout" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">
              Quantity / Seat Licenses (Min: {{ activeModalProduct.moq }})
            </label>
            <input
              v-model.number="orderQuantity"
              type="number"
              :min="activeModalProduct.moq"
              class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white focus:border-amber-500 outline-none"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Procurement Contact Name</label>
              <input
                v-model="orderCustomerName"
                type="text"
                placeholder="e.g. John Doe"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white focus:border-amber-500 outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Corporate Email</label>
              <input
                v-model="orderCustomerEmail"
                type="email"
                placeholder="procurement@corp.com"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white focus:border-amber-500 outline-none"
              />
            </div>
          </div>

          <!-- Price Calculation Breakdown -->
          <div class="p-3.5 rounded-xl bg-slate-950/70 border border-slate-800 space-y-1.5 text-xs">
            <div class="flex justify-between text-slate-400">
              <span>Unit Rate:</span>
              <span class="font-mono text-slate-200">${{ currentUnitPrice.toFixed(2) }}</span>
            </div>
            <div v-if="orderQuantity >= 10 && activeModalProduct.wholesale_price" class="flex justify-between text-amber-400">
              <span>Wholesale Volume Discount:</span>
              <span>Applied</span>
            </div>
            <div class="flex justify-between text-slate-400">
              <span>Fulfillment Protocol:</span>
              <span v-if="activeModalProduct.type !== 'physical'" class="text-cyan-400 font-semibold">
                Instant License Key + Asset Download
              </span>
              <span v-else class="text-amber-400 font-semibold">
                {{ activeModalProduct.lead_time_days }} Days Freight Dispatch
              </span>
            </div>
            <div class="pt-2 border-t border-slate-800 flex justify-between font-bold text-sm text-white">
              <span>Total Estimated Amount:</span>
              <span class="text-amber-400">${{ totalOrderAmount.toFixed(2) }}</span>
            </div>
          </div>

          <!-- Checkout Submit Button -->
          <button
            type="submit"
            :disabled="isSubmitting"
            class="w-full py-3 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-black text-sm shadow-lg shadow-amber-500/20 transition-all flex items-center justify-center gap-2 disabled:opacity-50"
          >
            <span v-if="isSubmitting">Processing Escrow Order...</span>
            <span v-else>Confirm & Execute Enterprise Order</span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Link, router } from '@inertiajs/vue3'
import { ref, computed } from 'vue'

const props = defineProps<{
  products: any[]
  categories: any[]
  activeType: string
  activeCategory: string
  searchQuery: string
  stats: any
}>()

const searchQueryInput = ref(props.searchQuery || '')
const activeModalProduct = ref<any | null>(null)
const orderQuantity = ref(1)
const orderCustomerName = ref('')
const orderCustomerEmail = ref('')
const isSubmitting = ref(false)

const triggerSearch = () => {
  router.get('/', {
    q: searchQueryInput.value,
    type: props.activeType !== 'all' ? props.activeType : undefined,
    category: props.activeCategory !== 'all' ? props.activeCategory : undefined,
  }, {
    preserveState: true,
  })
}

const quickSearch = (q: string) => {
  searchQueryInput.value = q
  triggerSearch()
}

const setTypeFilter = (type: string) => {
  router.get('/', {
    type: type !== 'all' ? type : undefined,
    category: props.activeCategory !== 'all' ? props.activeCategory : undefined,
    q: searchQueryInput.value || undefined,
  }, {
    preserveState: true,
  })
}

const setCategoryFilter = (catId: string) => {
  router.get('/', {
    category: catId !== 'all' ? catId : undefined,
    type: props.activeType !== 'all' ? props.activeType : undefined,
    q: searchQueryInput.value || undefined,
  }, {
    preserveState: true,
  })
}

const resetFilters = () => {
  searchQueryInput.value = ''
  router.get('/', {}, { preserveState: false })
}

const openOrderModal = (product: any) => {
  activeModalProduct.value = product
  orderQuantity.value = product.moq || 1
}

const currentUnitPrice = computed(() => {
  if (!activeModalProduct.value) return 0
  if (orderQuantity.value >= 10 && activeModalProduct.value.wholesale_price > 0) {
    return activeModalProduct.value.wholesale_price
  }
  return activeModalProduct.value.price
})

const totalOrderAmount = computed(() => {
  return currentUnitPrice.value * (orderQuantity.value || 1)
})

const submitCheckout = () => {
  if (!activeModalProduct.value) return
  isSubmitting.value = true

  router.post('/checkout', {
    product_id: activeModalProduct.value.id,
    quantity: orderQuantity.value,
    customer_name: orderCustomerName.value,
    customer_email: orderCustomerEmail.value,
  }, {
    onFinish: () => {
      isSubmitting.value = false
      activeModalProduct.value = null
    },
  })
}
</script>

