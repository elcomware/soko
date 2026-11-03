<template>
  <div class="space-y-12 pb-24">
    <!-- Active Regional Hub Context Strip -->
    <div class="bg-slate-900 border-b border-slate-800 py-2.5 px-4 sm:px-6">
      <div class="max-w-7xl mx-auto flex flex-wrap items-center justify-between gap-3 text-xs">
        <div class="flex items-center gap-2">
          <span class="text-base">{{ currentRegionObj?.flag || '🌐' }}</span>
          <span class="text-slate-400">Current Market:</span>
          <span class="font-bold text-white">{{ currentRegionObj?.name || 'Global International' }}</span>
          <span class="text-slate-600">|</span>
          <span class="text-slate-400">Logistics:</span>
          <span class="text-amber-400 font-medium">{{ currentRegionObj?.local_logistics || 'Global Multimodal Freight' }}</span>
        </div>

        <div class="flex items-center gap-3">
          <span class="text-slate-400">Local Rails:</span>
          <span
            v-for="rail in (currentRegionObj?.payment_rails || ['Trade Escrow', 'Cards', 'Wire'])"
            :key="rail"
            class="px-2 py-0.5 rounded bg-slate-800 text-[11px] text-slate-300 border border-slate-700"
          >
            {{ rail }}
          </span>
        </div>
      </div>
    </div>

    <!-- Hero Multi-Marketplace Section -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6">
      <div class="rounded-3xl bg-gradient-to-br from-slate-900 via-slate-950 to-slate-900 border border-slate-800 p-8 sm:p-12 relative overflow-hidden shadow-2xl">
        <div class="absolute -right-20 -top-20 w-96 h-96 bg-amber-500/10 rounded-full blur-3xl pointer-events-none"></div>
        <div class="absolute -left-20 -bottom-20 w-96 h-96 bg-cyan-500/10 rounded-full blur-3xl pointer-events-none"></div>

        <div class="relative z-10 max-w-3xl space-y-4">
          <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-slate-800/80 border border-slate-700 text-xs font-semibold text-slate-300">
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span>Alibaba Bulk Tiers · Amazon Prime Fulfillment · eBay Make Offer</span>
          </div>

          <h1 class="text-3xl sm:text-5xl font-black text-white tracking-tight leading-tight">
            The Hybrid Enterprise Commerce & Sourcing Engine
          </h1>

          <p class="text-sm sm:text-base text-slate-400 leading-relaxed max-w-2xl">
            Source factory-direct CNC equipment and optical sorting machinery with B2B wholesale volume tiers, procure certified refurbished datacenter network gear with price negotiation, or deploy instant cloud ERP licenses worldwide.
          </p>

          <!-- Format Switcher Pills -->
          <div class="pt-4 flex flex-wrap items-center gap-2">
            <button
              @click="setFormat('all')"
              :class="[
                'px-4 py-2 rounded-xl text-xs font-bold transition-all',
                activeFormat === 'all' ? 'bg-amber-500 text-slate-950 shadow-lg shadow-amber-500/20' : 'bg-slate-800/80 text-slate-300 hover:bg-slate-800'
              ]"
            >
              🏬 All Listings
            </button>
            <button
              @click="setFormat('wholesale')"
              :class="[
                'px-4 py-2 rounded-xl text-xs font-bold transition-all flex items-center gap-1.5',
                activeFormat === 'wholesale' ? 'bg-amber-500 text-slate-950 shadow-lg shadow-amber-500/20' : 'bg-slate-800/80 text-slate-300 hover:bg-slate-800'
              ]"
            >
              <span>🚢 Alibaba Wholesale Tiers</span>
            </button>
            <button
              @click="setFormat('digital')"
              :class="[
                'px-4 py-2 rounded-xl text-xs font-bold transition-all flex items-center gap-1.5',
                activeFormat === 'digital' ? 'bg-cyan-500 text-slate-950 shadow-lg shadow-cyan-500/20' : 'bg-slate-800/80 text-slate-300 hover:bg-slate-800'
              ]"
            >
              <span>⚡ Instant Digital Software</span>
            </button>
            <button
              @click="setFormat('make_offer')"
              :class="[
                'px-4 py-2 rounded-xl text-xs font-bold transition-all flex items-center gap-1.5',
                activeFormat === 'make_offer' ? 'bg-purple-500 text-white shadow-lg shadow-purple-500/20' : 'bg-slate-800/80 text-slate-300 hover:bg-slate-800'
              ]"
            >
              <span>🤝 eBay Deals & Make Offer</span>
            </button>
          </div>
        </div>

        <!-- Metric Badges in Hero -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 mt-10 pt-8 border-t border-slate-800/80 text-xs">
          <div>
            <div class="text-xl font-black text-white">${{ (stats.total_gmv || 48200).toLocaleString() }}</div>
            <div class="text-slate-400">Total GMV Traded</div>
          </div>
          <div>
            <div class="text-xl font-black text-emerald-400">${{ (stats.trade_assurance_reserves || 14250000).toLocaleString() }}</div>
            <div class="text-slate-400">Trade Escrow Covered</div>
          </div>
          <div>
            <div class="text-xl font-black text-cyan-400">Zero Latency</div>
            <div class="text-slate-400">Digital Key Dispatch SLA</div>
          </div>
          <div>
            <div class="text-xl font-black text-amber-400">10 Currencies</div>
            <div class="text-slate-400">Local Currency Converted</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Main Catalog View: Filter Sidebar + Products Grid -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6">
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Left Filter Sidebar -->
        <aside class="space-y-6">
          <div class="p-5 rounded-2xl bg-slate-900 border border-slate-800 space-y-6">
            <div class="flex items-center justify-between pb-3 border-b border-slate-800">
              <h3 class="font-bold text-sm text-white flex items-center gap-1.5">
                <span>⚙️</span>
                <span>Filters & Criteria</span>
              </h3>
              <button @click="resetFilters" class="text-[11px] text-amber-400 hover:underline">Reset</button>
            </div>

            <!-- Categories Filter -->
            <div class="space-y-2">
              <label class="block text-xs font-bold text-slate-300 uppercase tracking-wider text-[10px]">Categories</label>
              <div class="space-y-1">
                <button
                  v-for="cat in categories"
                  :key="cat.id"
                  @click="setCategory(cat.id)"
                  :class="[
                    'w-full text-left px-2.5 py-1.5 rounded-lg text-xs transition-colors flex items-center justify-between',
                    activeCategory === cat.id ? 'bg-amber-500/10 text-amber-400 font-bold' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800'
                  ]"
                >
                  <span class="truncate">{{ cat.name }}</span>
                  <span v-if="cat.count" class="text-[10px] text-slate-600">({{ cat.count }})</span>
                </button>
              </div>
            </div>

            <!-- Condition Filter (eBay style) -->
            <div class="space-y-2 pt-3 border-t border-slate-800">
              <label class="block text-xs font-bold text-slate-300 uppercase tracking-wider text-[10px]">Item Condition</label>
              <div class="flex flex-wrap gap-1.5">
                <button
                  @click="setCondition('all')"
                  :class="[
                    'px-2.5 py-1 rounded text-xs transition-colors',
                    activeCondition === 'all' ? 'bg-slate-700 text-white font-bold' : 'bg-slate-950 text-slate-400 hover:text-white'
                  ]"
                >
                  All
                </button>
                <button
                  @click="setCondition('new')"
                  :class="[
                    'px-2.5 py-1 rounded text-xs transition-colors',
                    activeCondition === 'new' ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30' : 'bg-slate-950 text-slate-400 hover:text-white'
                  ]"
                >
                  Brand New
                </button>
                <button
                  @click="setCondition('refurbished')"
                  :class="[
                    'px-2.5 py-1 rounded text-xs transition-colors',
                    activeCondition === 'refurbished' ? 'bg-cyan-500/20 text-cyan-300 border border-cyan-500/30' : 'bg-slate-950 text-slate-400 hover:text-white'
                  ]"
                >
                  Refurbished
                </button>
              </div>
            </div>

            <!-- Regional Origin Filter -->
            <div class="space-y-2 pt-3 border-t border-slate-800">
              <label class="block text-xs font-bold text-slate-300 uppercase tracking-wider text-[10px]">Sourcing Region</label>
              <select
                v-model="filterRegion"
                @change="onFilterRegionChange"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3 py-2 text-xs text-white outline-none"
              >
                <option value="all">Any Sourcing Origin</option>
                <option value="global">Global International</option>
                <option value="north_america">North America</option>
                <option value="europe">Europe & UK</option>
                <option value="east_africa">East Africa</option>
                <option value="west_africa">West Africa</option>
                <option value="asia_pacific">Asia-Pacific</option>
              </select>
            </div>

            <!-- Trade Assurance Callout -->
            <div class="p-3.5 rounded-xl bg-slate-950 border border-amber-500/20 space-y-1.5 text-xs">
              <div class="font-bold text-amber-400 flex items-center gap-1">
                <span>🛡️</span>
                <span>Trade Assurance</span>
              </div>
              <p class="text-[11px] text-slate-400 leading-relaxed">
                Covers your purchase from payment to final dispatch or instant software key validation.
              </p>
            </div>
          </div>

          <!-- B2B RFQ Promo Box -->
          <div class="p-5 rounded-2xl bg-gradient-to-br from-slate-900 to-amber-950/30 border border-slate-800 space-y-3">
            <span class="text-xs font-bold uppercase tracking-wider text-amber-400">Need Custom Supply?</span>
            <h4 class="font-bold text-white text-sm">Post a Request for Quotation (RFQ)</h4>
            <p class="text-xs text-slate-400">Over 50,000 verified factories and software developers will bid on your custom specs.</p>
            <button
              @click="showRfqModal = true"
              class="w-full py-2 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs shadow-md transition-colors"
            >
              Create Sourcing RFQ &rarr;
            </button>
          </div>
        </aside>

        <!-- Right Products Catalog -->
        <div class="lg:col-span-3 space-y-6">
          <!-- Active filter status bar -->
          <div class="flex items-center justify-between pb-2 border-b border-slate-800 text-xs text-slate-400">
            <div>
              Showing <span class="font-bold text-white">{{ products.length }}</span> items in
              <span class="text-amber-400 font-bold">{{ activeCurrency.code }}</span>
            </div>
            <div class="flex items-center gap-2">
              <span>Sort by:</span>
              <span class="font-semibold text-white">Featured & Verified</span>
            </div>
          </div>

          <!-- Empty State -->
          <div v-if="products.length === 0" class="text-center py-20 bg-slate-900/40 rounded-3xl border border-slate-800">
            <div class="text-3xl mb-2">🔍</div>
            <h3 class="font-bold text-base text-white">No products match your current regional criteria</h3>
            <p class="text-xs text-slate-400 mt-1">Try broadening your format or category filters.</p>
            <button @click="resetFilters" class="mt-4 px-4 py-2 rounded-xl bg-slate-800 text-xs font-bold text-slate-200">
              Clear All Filters
            </button>
          </div>

          <!-- Mature Product Grid -->
          <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
            <div
              v-for="product in products"
              :key="product.id"
              class="rounded-2xl bg-slate-900 border border-slate-800 hover:border-slate-700 transition-all flex flex-col justify-between overflow-hidden shadow-lg hover:shadow-2xl hover:shadow-slate-950/80 group"
            >
              <div>
                <!-- Image Header with Format Badges -->
                <div class="relative h-48 w-full overflow-hidden bg-slate-950">
                  <img
                    :src="product.image_url"
                    :alt="product.title"
                    class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                  />
                  <div class="absolute inset-0 bg-gradient-to-t from-slate-950/90 via-transparent to-transparent"></div>

                  <!-- Format Badge -->
                  <div class="absolute top-3 left-3 flex flex-col gap-1">
                    <span
                      v-if="product.format === 'wholesale'"
                      class="px-2.5 py-0.5 rounded text-[10px] font-black uppercase tracking-wider bg-amber-950/90 text-amber-300 border border-amber-500/40 backdrop-blur"
                    >
                      🚢 B2B Wholesale · MOQ {{ product.moq }}
                    </span>
                    <span
                      v-else-if="product.format === 'digital'"
                      class="px-2.5 py-0.5 rounded text-[10px] font-black uppercase tracking-wider bg-cyan-950/90 text-cyan-300 border border-cyan-500/40 backdrop-blur flex items-center gap-1"
                    >
                      <span class="w-1.5 h-1.5 rounded-full bg-cyan-400 animate-pulse"></span>
                      ⚡ Instant Digital
                    </span>
                    <span
                      v-else-if="product.format === 'make_offer'"
                      class="px-2.5 py-0.5 rounded text-[10px] font-black uppercase tracking-wider bg-purple-950/90 text-purple-300 border border-purple-500/40 backdrop-blur"
                    >
                      🤝 Accepts Offers
                    </span>
                    <span
                      v-else
                      class="px-2.5 py-0.5 rounded text-[10px] font-black uppercase tracking-wider bg-emerald-950/90 text-emerald-300 border border-emerald-500/40 backdrop-blur"
                    >
                      🛒 Fast Buy
                    </span>

                    <span
                      v-if="product.condition === 'refurbished'"
                      class="px-2 py-0.5 rounded text-[9px] font-bold bg-slate-900/90 text-slate-300 border border-slate-700 w-fit"
                    >
                      Certified Refurbished
                    </span>
                  </div>

                  <!-- Rating -->
                  <div class="absolute top-3 right-3 bg-slate-950/80 px-2 py-0.5 rounded text-[10px] font-bold text-amber-400 border border-slate-800 flex items-center gap-1 backdrop-blur">
                    <span>★</span>
                    <span>{{ product.rating }}</span>
                    <span class="text-slate-500">({{ product.reviews_count }})</span>
                  </div>

                  <!-- Category pill -->
                  <div class="absolute bottom-2 left-3 text-[11px] font-medium text-slate-300 truncate max-w-[85%]">
                    {{ product.category }}
                  </div>
                </div>

                <!-- Product Body -->
                <div class="p-4 space-y-3">
                  <h3 class="font-bold text-sm text-white group-hover:text-amber-400 transition-colors line-clamp-2 leading-snug">
                    {{ product.title }}
                  </h3>

                  <!-- Supplier Trust Line -->
                  <div class="p-2 rounded-xl bg-slate-950/60 border border-slate-800 text-[11px] space-y-1">
                    <div class="flex items-center justify-between">
                      <div class="flex items-center gap-1.5 truncate">
                        <span class="text-amber-400 font-bold">✓</span>
                        <span class="font-semibold text-slate-200 truncate">{{ product.vendor_name }}</span>
                      </div>
                      <span class="text-[10px] text-emerald-400 font-mono shrink-0">{{ product.vendor_positive_rate }}</span>
                    </div>
                    <div class="flex items-center justify-between text-[10px] text-slate-400">
                      <span>{{ product.vendor_country }} ({{ product.incoterm || 'FOB/CIF' }})</span>
                      <span class="text-slate-500 font-mono">Escrow: ${{ (product.trade_assurance_cover / 1000).toFixed(0) }}k</span>
                    </div>
                  </div>

                  <!-- Alibaba-Style Volume Price Tiers Preview -->
                  <div v-if="product.price_tiers && product.price_tiers.length > 0" class="space-y-1">
                    <div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Volume Tier Pricing:</div>
                    <div class="grid grid-cols-3 gap-1 text-center bg-slate-950 p-1.5 rounded-lg border border-slate-800 text-[10px]">
                      <div v-for="(tier, idx) in product.price_tiers.slice(0, 3)" :key="idx" class="border-r last:border-0 border-slate-800">
                        <div class="text-slate-500">
                          {{ tier.max_qty && tier.max_qty < 999 ? `${tier.min_qty}-${tier.max_qty}` : `${tier.min_qty}+` }} pcs
                        </div>
                        <div class="font-bold text-amber-400">{{ formatPrice(tier.price) }}</div>
                      </div>
                    </div>
                  </div>
                  <div v-else class="text-xs text-slate-400">
                    <span v-if="product.type !== 'physical'" class="text-cyan-400 font-mono">
                      ⚡ Instant Token · {{ product.digital_file_size || 'Downloadable' }}
                    </span>
                    <span v-else class="text-slate-400">
                      Dispatches in {{ product.lead_time_days }} business days
                    </span>
                  </div>
                </div>
              </div>

              <!-- Product Footer & Actions -->
              <div class="p-4 pt-0 space-y-3">
                <div class="pt-3 border-t border-slate-800 flex items-baseline justify-between">
                  <div>
                    <div class="text-[10px] text-slate-500">Unit Price</div>
                    <div class="text-lg font-black text-white">
                      {{ formatPrice(product.price) }}
                    </div>
                  </div>
                  <div class="text-right">
                    <div class="text-[10px] text-slate-500">Min. Order</div>
                    <div class="text-xs font-bold text-slate-200">{{ product.moq }} unit(s)</div>
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-2">
                  <Link
                    :href="'/products/' + product.id"
                    class="px-2.5 py-2 rounded-xl text-xs font-semibold bg-slate-800 hover:bg-slate-700 text-slate-200 text-center transition-colors truncate"
                  >
                    Specs & Tiers
                  </Link>

                  <!-- Secondary Action: Buy Now vs Make an Offer -->
                  <button
                    v-if="product.accepts_offers"
                    @click="openOfferModal(product)"
                    class="px-2.5 py-2 rounded-xl text-xs font-bold bg-gradient-to-r from-purple-500 to-indigo-600 hover:from-purple-400 hover:to-indigo-500 text-white shadow-md shadow-purple-500/20 transition-all truncate"
                  >
                    🤝 Make Offer
                  </button>
                  <button
                    v-else
                    @click="openOrderModal(product)"
                    class="px-2.5 py-2 rounded-xl text-xs font-bold bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 shadow-md shadow-amber-500/20 transition-all truncate"
                  >
                    {{ product.type !== 'physical' ? '⚡ Procure' : '📦 Order' }}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- B2B RFQ Marketplace Section (Alibaba feature) -->
    <section id="rfq" class="max-w-7xl mx-auto px-4 sm:px-6">
      <div class="rounded-3xl bg-slate-900 border border-slate-800 p-8 space-y-6">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <div class="inline-flex items-center gap-1 text-xs font-bold text-amber-400 uppercase tracking-wider mb-1">
              <span>📋</span>
              <span>Alibaba-Scale B2B RFQ Marketplace</span>
            </div>
            <h2 class="text-2xl font-black text-white">Live Buyer Procurement Demands</h2>
            <p class="text-xs text-slate-400">Suppliers submit custom quotations with shipping Incoterms and delivery timetables.</p>
          </div>
          <button
            @click="showRfqModal = true"
            class="px-5 py-2.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs shadow-md transition-colors shrink-0"
          >
            + Post New RFQ
          </button>
        </div>

        <!-- RFQ Table -->
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-slate-950 text-slate-400 uppercase text-[10px] tracking-wider border-b border-slate-800">
              <tr>
                <th class="py-3 px-4">RFQ Ref & Title</th>
                <th class="py-3 px-4">Category</th>
                <th class="py-3 px-4">Quantity Required</th>
                <th class="py-3 px-4">Target Unit Price</th>
                <th class="py-3 px-4">Destination Port/Region</th>
                <th class="py-3 px-4 text-right">Action</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800">
              <tr v-for="r in rfqs" :key="r.id" class="hover:bg-slate-800/40 transition-colors">
                <td class="py-3.5 px-4">
                  <div class="font-bold text-white">{{ r.title }}</div>
                  <div class="text-[10px] text-slate-500 font-mono">{{ r.rfq_number }} · By {{ r.buyer_name }} ({{ r.buyer_country }})</div>
                </td>
                <td class="py-3.5 px-4 text-slate-300">{{ r.category }}</td>
                <td class="py-3.5 px-4 font-bold text-white">{{ r.quantity.toLocaleString() }} pcs</td>
                <td class="py-3.5 px-4 font-mono text-amber-400 font-bold">{{ formatPrice(r.target_price) }}</td>
                <td class="py-3.5 px-4 text-slate-400">{{ r.destination_country }}</td>
                <td class="py-3.5 px-4 text-right">
                  <Link
                    href="/dashboard"
                    class="px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 font-bold text-xs"
                  >
                    Submit Quotation
                  </Link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </section>

    <!-- Modal: Make an Offer (eBay Style) -->
    <div
      v-if="offerModalProduct"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm"
    >
      <div class="bg-slate-900 border border-slate-700 rounded-3xl max-w-lg w-full p-6 shadow-2xl space-y-6 relative">
        <button @click="offerModalProduct = null" class="absolute top-5 right-5 text-slate-400 hover:text-white text-lg font-bold">
          &times;
        </button>

        <div>
          <div class="text-xs font-bold text-purple-400 uppercase tracking-wider">🤝 Make an Offer (Price Negotiation)</div>
          <h3 class="text-base font-bold text-white mt-1">{{ offerModalProduct.title }}</h3>
          <p class="text-xs text-slate-400">Seller: {{ offerModalProduct.vendor_name }} · Listed at {{ formatPrice(offerModalProduct.price) }}</p>
        </div>

        <form @submit.prevent="submitOfferForm" class="space-y-4">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Your Proposed Price ($ USD)</label>
              <input
                v-model.number="offerPrice"
                type="number"
                step="0.01"
                required
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white focus:border-purple-500 outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Quantity</label>
              <input
                v-model.number="offerQty"
                type="number"
                min="1"
                required
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white focus:border-purple-500 outline-none"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Your Name</label>
              <input
                v-model="offerName"
                type="text"
                required
                placeholder="e.g. John Doe"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Corporate Email</label>
              <input
                v-model="offerEmail"
                type="email"
                required
                placeholder="procurement@corp.com"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Message to Seller (Terms, Destination Port)</label>
            <textarea
              v-model="offerMessage"
              rows="2"
              placeholder="e.g. We require delivery CIF Rotterdam and can pay via Escrow upon bill of lading..."
              class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
            ></textarea>
          </div>

          <div class="p-3 rounded-xl bg-slate-950 border border-purple-500/20 flex justify-between items-center text-xs">
            <span class="text-slate-400">Total Offer Value:</span>
            <span class="font-bold text-purple-400 text-sm">${{ (offerPrice * offerQty).toFixed(2) }} USD</span>
          </div>

          <button
            type="submit"
            class="w-full py-3 rounded-xl bg-gradient-to-r from-purple-500 to-indigo-600 hover:from-purple-400 hover:to-indigo-500 text-white font-bold text-xs shadow-lg transition-all"
          >
            Submit Formal Offer to Merchant
          </button>
        </form>
      </div>
    </div>

    <!-- Modal: Instant Order / Procure Checkout -->
    <div
      v-if="orderModalProduct"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm"
    >
      <div class="bg-slate-900 border border-slate-700 rounded-3xl max-w-lg w-full p-6 shadow-2xl space-y-6 relative">
        <button @click="orderModalProduct = null" class="absolute top-5 right-5 text-slate-400 hover:text-white text-lg font-bold">
          &times;
        </button>

        <div class="flex items-center gap-3">
          <img :src="orderModalProduct.image_url" class="w-14 h-14 rounded-xl object-cover border border-slate-700" alt="prod" />
          <div>
            <div class="text-xs font-bold text-amber-400 uppercase tracking-wider">
              {{ orderModalProduct.type !== 'physical' ? '⚡ Instant Digital Fulfillment' : '🚢 Soko Trade Escrow' }}
            </div>
            <h3 class="text-base font-bold text-white line-clamp-1">{{ orderModalProduct.title }}</h3>
            <div class="text-xs text-slate-400">{{ orderModalProduct.vendor_name }} ({{ orderModalProduct.vendor_country }})</div>
          </div>
        </div>

        <form @submit.prevent="submitCheckout" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">
              Quantity (Min Order: {{ orderModalProduct.moq }} pcs)
            </label>
            <input
              v-model.number="orderQuantity"
              type="number"
              :min="orderModalProduct.moq"
              class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white focus:border-amber-500 outline-none"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Contact Name</label>
              <input
                v-model="orderCustomerName"
                type="text"
                placeholder="Full name"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Work Email</label>
              <input
                v-model="orderCustomerEmail"
                type="email"
                placeholder="Email address"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
          </div>

          <!-- Total Calculation -->
          <div class="p-3.5 rounded-xl bg-slate-950 border border-slate-800 space-y-1.5 text-xs">
            <div class="flex justify-between text-slate-400">
              <span>Unit Price ({{ activeCurrency.code }}):</span>
              <span class="font-mono text-white font-bold">{{ formatPrice(computedUnitPrice) }}</span>
            </div>
            <div class="flex justify-between text-slate-400">
              <span>Settlement Rail:</span>
              <span class="text-emerald-400 font-semibold">{{ currentRegionObj?.payment_rails?.[0] || 'Soko Trade Escrow' }}</span>
            </div>
            <div class="pt-2 border-t border-slate-800 flex justify-between font-bold text-sm text-white">
              <span>Total Amount:</span>
              <span class="text-amber-400 font-black">{{ formatPrice(computedUnitPrice * orderQuantity) }}</span>
            </div>
          </div>

          <button
            type="submit"
            class="w-full py-3.5 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-black text-sm shadow-xl transition-all"
          >
            Confirm & Execute Order
          </button>
        </form>
      </div>
    </div>

    <!-- Modal: Post New RFQ (Alibaba Feature) -->
    <div
      v-if="showRfqModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm"
    >
      <div class="bg-slate-900 border border-slate-700 rounded-3xl max-w-lg w-full p-6 shadow-2xl space-y-6 relative">
        <button @click="showRfqModal = false" class="absolute top-5 right-5 text-slate-400 hover:text-white text-lg font-bold">
          &times;
        </button>

        <div>
          <div class="text-xs font-bold text-amber-400 uppercase tracking-wider">📋 Post Request For Quotation</div>
          <h3 class="text-lg font-black text-white mt-1">Submit Sourcing Brief to Global Suppliers</h3>
        </div>

        <form @submit.prevent="submitRfqForm" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Product Sourcing Title</label>
            <input
              v-model="rfqTitle"
              type="text"
              required
              placeholder="e.g. 2,000 units Stainless Steel Flow Control Valves"
              class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Target Quantity</label>
              <input
                v-model.number="rfqQuantity"
                type="number"
                min="1"
                required
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Target Unit Price ($ USD)</label>
              <input
                v-model.number="rfqTargetPrice"
                type="number"
                step="0.01"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Destination Country / Port</label>
              <input
                v-model="rfqDestination"
                type="text"
                required
                placeholder="e.g. Kenya (Mombasa) or USA"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Company Name</label>
              <input
                v-model="rfqCompany"
                type="text"
                placeholder="e.g. Megacorp Logistics"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Detailed Technical Specifications</label>
            <textarea
              v-model="rfqSpecs"
              rows="2"
              placeholder="Material grade, tolerances, certifications required (CE, ISO, RoHS)..."
              class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
            ></textarea>
          </div>

          <button
            type="submit"
            class="w-full py-3 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs shadow-lg transition-all"
          >
            Publish RFQ to Soko Network
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Link, router } from '@inertiajs/vue3'
import { ref, computed } from 'vue'
import { useCurrency, type Currency, type MarketRegion } from '@/composables/useCurrency'

const props = defineProps<{
  products: any[]
  categories: any[]
  currencies: Currency[]
  regions: MarketRegion[]
  rfqs: any[]
  activeType: string
  activeFormat: string
  activeRegion: string
  activeCategory: string
  activeCondition: string
  searchQuery: string
  stats: any
}>()

const { selectedRegion, activeCurrency, formatPrice } = useCurrency(props.currencies)

const currentRegionObj = computed(() => {
  return props.regions.find(r => r.id === props.activeRegion) || props.regions[0]
})

const filterRegion = ref(props.activeRegion || 'global')

// Modals
const orderModalProduct = ref<any | null>(null)
const orderQuantity = ref(1)
const orderCustomerName = ref('')
const orderCustomerEmail = ref('')

const offerModalProduct = ref<any | null>(null)
const offerPrice = ref(0)
const offerQty = ref(1)
const offerName = ref('')
const offerEmail = ref('')
const offerMessage = ref('')

const showRfqModal = ref(false)
const rfqTitle = ref('')
const rfqQuantity = ref(100)
const rfqTargetPrice = ref(50)
const rfqDestination = ref('')
const rfqCompany = ref('')
const rfqSpecs = ref('')

const setFormat = (fmt: string) => {
  router.get('/', {
    format: fmt !== 'all' ? fmt : undefined,
    region: props.activeRegion !== 'global' ? props.activeRegion : undefined,
    category: props.activeCategory !== 'all' ? props.activeCategory : undefined,
  }, { preserveState: true })
}

const setCategory = (catId: string) => {
  router.get('/', {
    category: catId !== 'all' ? catId : undefined,
    format: props.activeFormat !== 'all' ? props.activeFormat : undefined,
    region: props.activeRegion !== 'global' ? props.activeRegion : undefined,
  }, { preserveState: true })
}

const setCondition = (cond: string) => {
  router.get('/', {
    condition: cond !== 'all' ? cond : undefined,
    format: props.activeFormat !== 'all' ? props.activeFormat : undefined,
    region: props.activeRegion !== 'global' ? props.activeRegion : undefined,
  }, { preserveState: true })
}

const onFilterRegionChange = () => {
  router.get('/', {
    region: filterRegion.value !== 'all' ? filterRegion.value : undefined,
    format: props.activeFormat !== 'all' ? props.activeFormat : undefined,
    category: props.activeCategory !== 'all' ? props.activeCategory : undefined,
  }, { preserveState: true })
}

const resetFilters = () => {
  router.get('/', {}, { preserveState: false })
}

// Order handling
const openOrderModal = (product: any) => {
  orderModalProduct.value = product
  orderQuantity.value = product.moq || 1
}

const computedUnitPrice = computed(() => {
  if (!orderModalProduct.value) return 0
  const p = orderModalProduct.value
  if (p.price_tiers && p.price_tiers.length > 0) {
    for (const tier of p.price_tiers) {
      if (orderQuantity.value >= tier.min_qty && (!tier.max_qty || orderQuantity.value <= tier.max_qty)) {
        return tier.price
      }
    }
  }
  if (orderQuantity.value >= 10 && p.wholesale_price > 0) {
    return p.wholesale_price
  }
  return p.price
})

const submitCheckout = () => {
  if (!orderModalProduct.value) return
  router.post('/checkout', {
    product_id: orderModalProduct.value.id,
    quantity: orderQuantity.value,
    customer_name: orderCustomerName.value,
    customer_email: orderCustomerEmail.value,
    currency: activeCurrency.value.code,
  }, {
    onFinish: () => {
      orderModalProduct.value = null
    },
  })
}

// Offer handling
const openOfferModal = (product: any) => {
  offerModalProduct.value = product
  offerPrice.value = product.wholesale_price || (product.price * 0.85)
  offerQty.value = product.moq || 1
}

const submitOfferForm = () => {
  if (!offerModalProduct.value) return
  router.post('/make-offer', {
    product_id: offerModalProduct.value.id,
    offered_price: offerPrice.value,
    quantity: offerQty.value,
    name: offerName.value,
    email: offerEmail.value,
    message: offerMessage.value,
  }, {
    onFinish: () => {
      offerModalProduct.value = null
    },
  })
}

// RFQ handling
const submitRfqForm = () => {
  router.post('/rfqs', {
    title: rfqTitle.value,
    quantity: rfqQuantity.value,
    target_price: rfqTargetPrice.value,
    destination_country: rfqDestination.value,
    buyer_company: rfqCompany.value,
    specifications: rfqSpecs.value,
  }, {
    onFinish: () => {
      showRfqModal.value = false
    },
  })
}
</script>
