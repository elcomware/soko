<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 py-8 space-y-8">
    <!-- Back Store Top Console Bar -->
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-6 p-6 rounded-3xl bg-slate-900 border border-slate-800 shadow-2xl">
      <div class="flex items-center gap-4">
        <img
          :src="currentUser.avatar || 'https://ui-avatars.com/api/?name=' + currentUser.name"
          class="w-16 h-16 rounded-2xl object-cover border-2 border-amber-500/50 shadow-md shrink-0"
          alt="Avatar"
        />
        <div>
          <div class="flex items-center gap-2 flex-wrap">
            <h1 class="text-2xl font-black text-white">{{ currentUser.company_name || currentUser.name }}</h1>
            <span
              :class="[
                'text-[10px] font-black px-2 py-0.5 rounded uppercase tracking-wider',
                currentUser.role === 'admin'
                  ? 'bg-rose-500/20 text-rose-400 border border-rose-500/30'
                  : currentUser.role === 'merchant'
                  ? 'bg-amber-500/20 text-amber-400 border border-amber-500/30'
                  : 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/30'
              ]"
            >
              {{ currentUser.role }} console
            </span>
            <span class="px-2 py-0.5 rounded bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 text-[10px] font-bold">
              ✓ Verified Global Merchant
            </span>
          </div>
          <div class="text-xs text-slate-400 mt-1 flex flex-wrap items-center gap-2">
            <span>{{ currentUser.email }}</span>
            <span>·</span>
            <span>Region: {{ currentUser.country || 'Global' }}</span>
            <span>·</span>
            <span class="text-amber-400 font-mono">Trade Escrow Limit: $2,500,000 USD</span>
          </div>
        </div>
      </div>

      <!-- Quick Action Buttons -->
      <div class="flex items-center gap-3 shrink-0 flex-wrap">
        <Link
          href="/"
          class="px-4 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs font-semibold text-slate-200 transition-colors"
        >
          View eStore Front
        </Link>
        <button
          @click="showAddModal = true"
          class="px-4 py-2.5 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-bold text-xs shadow-lg shadow-amber-500/20 transition-all flex items-center gap-1.5"
        >
          <span>+ Create Enterprise Listing</span>
        </button>
      </div>
    </div>

    <!-- KPI Metric Cards (Hybrid Alibaba + Amazon + eBay) -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="p-5 rounded-2xl bg-slate-900 border border-slate-800 space-y-2">
        <div class="flex items-center justify-between text-xs text-slate-400">
          <span>Gross Merchandise Value (GMV)</span>
          <span class="text-emerald-400 font-bold">+24.8%</span>
        </div>
        <div class="text-2xl font-black text-white">{{ formatPrice(stats.total_gmv || 48200) }}</div>
        <div class="text-[11px] text-slate-500">Processed via Soko Escrow</div>
      </div>

      <div class="p-5 rounded-2xl bg-slate-900 border border-slate-800 space-y-2">
        <div class="flex items-center justify-between text-xs text-slate-400">
          <span>Orders & Shipments</span>
          <span class="text-amber-400 font-bold">100% Fulfilled</span>
        </div>
        <div class="text-2xl font-black text-white">{{ stats.total_orders }} Active</div>
        <div class="text-[11px] text-slate-500">Wholesale freight & instant digital</div>
      </div>

      <div class="p-5 rounded-2xl bg-slate-900 border border-slate-800 space-y-2">
        <div class="flex items-center justify-between text-xs text-slate-400">
          <span>Active B2B RFQs & Offers</span>
          <span class="text-purple-400 font-bold">Open Bidding</span>
        </div>
        <div class="text-2xl font-black text-purple-400">
          {{ stats.active_rfqs_count }} RFQs · {{ stats.active_offers_count }} Offers
        </div>
        <div class="text-[11px] text-slate-500">B2B buyer procurement requests</div>
      </div>

      <div class="p-5 rounded-2xl bg-slate-900 border border-slate-800 space-y-2">
        <div class="flex items-center justify-between text-xs text-slate-400">
          <span>Escrow Capital Reserve</span>
          <span class="text-cyan-400 font-mono text-[10px]">Tier-1 Trust</span>
        </div>
        <div class="text-2xl font-black text-cyan-400">{{ formatPrice(stats.trade_assurance_reserves || 14250000) }}</div>
        <div class="text-[11px] text-slate-500">Backing supplier delivery SLAs</div>
      </div>
    </div>

    <!-- Back Store Navigation Tabs -->
    <div class="border-b border-slate-800 flex items-center gap-2 overflow-x-auto text-xs font-bold scrollbar-none">
      <button
        @click="currentTab = 'catalog'"
        :class="[
          'px-4 py-3 border-b-2 transition-all flex items-center gap-2 whitespace-nowrap',
          currentTab === 'catalog' ? 'border-amber-500 text-amber-400' : 'border-transparent text-slate-400 hover:text-slate-200'
        ]"
      >
        <span>📦 Enterprise Catalog</span>
        <span class="px-2 py-0.5 rounded-full bg-slate-800 text-[10px] text-slate-300">{{ products.length }}</span>
      </button>

      <button
        @click="currentTab = 'orders'"
        :class="[
          'px-4 py-3 border-b-2 transition-all flex items-center gap-2 whitespace-nowrap',
          currentTab === 'orders' ? 'border-amber-500 text-amber-400' : 'border-transparent text-slate-400 hover:text-slate-200'
        ]"
      >
        <span>📜 Orders & Digital Licenses</span>
        <span class="px-2 py-0.5 rounded-full bg-slate-800 text-[10px] text-slate-300">{{ orders.length }}</span>
      </button>

      <button
        @click="currentTab = 'rfqs'"
        :class="[
          'px-4 py-3 border-b-2 transition-all flex items-center gap-2 whitespace-nowrap',
          currentTab === 'rfqs' ? 'border-amber-500 text-amber-400' : 'border-transparent text-slate-400 hover:text-slate-200'
        ]"
      >
        <span>📋 B2B RFQ Marketplace</span>
        <span class="px-2 py-0.5 rounded-full bg-slate-800 text-[10px] text-amber-400">{{ rfqs.length }}</span>
      </button>

      <button
        @click="currentTab = 'offers'"
        :class="[
          'px-4 py-3 border-b-2 transition-all flex items-center gap-2 whitespace-nowrap',
          currentTab === 'offers' ? 'border-amber-500 text-amber-400' : 'border-transparent text-slate-400 hover:text-slate-200'
        ]"
      >
        <span>🤝 Inbound Buyer Offers</span>
        <span class="px-2 py-0.5 rounded-full bg-slate-800 text-[10px] text-purple-400">{{ offers.length }}</span>
      </button>

      <button
        v-if="currentUser.role === 'admin'"
        @click="currentTab = 'users'"
        :class="[
          'px-4 py-3 border-b-2 transition-all flex items-center gap-2 whitespace-nowrap',
          currentTab === 'users' ? 'border-amber-500 text-amber-400' : 'border-transparent text-slate-400 hover:text-slate-200'
        ]"
      >
        <span>👥 Verified Suppliers & Hubs</span>
        <span class="px-2 py-0.5 rounded-full bg-slate-800 text-[10px] text-slate-300">{{ users.length }}</span>
      </button>
    </div>

    <!-- Tab 1: Enterprise Catalog -->
    <div v-if="currentTab === 'catalog'" class="space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="text-base font-bold text-white">Active Product Listings Across Regional Hubs</h3>
        <button
          @click="showAddModal = true"
          class="px-3.5 py-1.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs"
        >
          + Add Listing
        </button>
      </div>

      <div class="rounded-2xl bg-slate-900 border border-slate-800 overflow-hidden shadow-xl">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-slate-950 text-slate-400 uppercase text-[10px] tracking-wider border-b border-slate-800">
              <tr>
                <th class="py-3 px-4">Item & SKU</th>
                <th class="py-3 px-4">Format / Model</th>
                <th class="py-3 px-4">Category</th>
                <th class="py-3 px-4">Volume Tiers / Pricing</th>
                <th class="py-3 px-4">Region & Origin</th>
                <th class="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800">
              <tr v-for="product in products" :key="product.id" class="hover:bg-slate-800/40 transition-colors">
                <td class="py-3.5 px-4">
                  <div class="flex items-center gap-3">
                    <img :src="product.image_url" class="w-10 h-10 rounded-lg object-cover border border-slate-800 shrink-0" alt="Thumb" />
                    <div>
                      <div class="font-bold text-white max-w-xs truncate">{{ product.title }}</div>
                      <div class="font-mono text-[10px] text-slate-500">{{ product.sku }} · {{ product.condition }}</div>
                    </div>
                  </div>
                </td>
                <td class="py-3.5 px-4">
                  <span
                    v-if="product.format === 'wholesale'"
                    class="px-2 py-0.5 rounded text-[10px] font-bold uppercase bg-amber-950 text-amber-400 border border-amber-800"
                  >
                    🚢 Wholesale Tiers
                  </span>
                  <span
                    v-else-if="product.format === 'digital'"
                    class="px-2 py-0.5 rounded text-[10px] font-bold uppercase bg-cyan-950 text-cyan-400 border border-cyan-800"
                  >
                    ⚡ Instant Digital
                  </span>
                  <span
                    v-else-if="product.format === 'make_offer'"
                    class="px-2 py-0.5 rounded text-[10px] font-bold uppercase bg-purple-950 text-purple-400 border border-purple-800"
                  >
                    🤝 Make Offer
                  </span>
                  <span
                    v-else
                    class="px-2 py-0.5 rounded text-[10px] font-bold uppercase bg-emerald-950 text-emerald-400 border border-emerald-800"
                  >
                    🛒 Fast Buy
                  </span>
                </td>
                <td class="py-3.5 px-4 text-slate-300">{{ product.category }}</td>
                <td class="py-3.5 px-4">
                  <div class="font-bold text-white">{{ formatPrice(product.price) }}</div>
                  <div v-if="product.wholesale_price" class="text-[10px] text-amber-400">
                    Bulk: {{ formatPrice(product.wholesale_price) }}
                  </div>
                </td>
                <td class="py-3.5 px-4">
                  <div class="text-slate-300 font-semibold">{{ product.vendor_region || 'Global' }}</div>
                  <div class="text-[10px] text-slate-500">{{ product.vendor_country }} ({{ product.incoterm || 'FOB' }})</div>
                </td>
                <td class="py-3.5 px-4 text-right space-x-2">
                  <Link :href="'/products/' + product.id" class="px-2 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300">
                    View
                  </Link>
                  <button @click="deleteProduct(product.id)" class="px-2 py-1 rounded bg-rose-950 hover:bg-rose-900 text-rose-300">
                    Delete
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Tab 2: Orders & Digital Licenses -->
    <div v-if="currentTab === 'orders'" class="space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="text-base font-bold text-white">Orders & Automated Digital License Issuance Vault</h3>
        <span class="text-xs text-slate-400">Escrow verification & license tokens</span>
      </div>

      <div class="rounded-2xl bg-slate-900 border border-slate-800 overflow-hidden shadow-xl">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-slate-950 text-slate-400 uppercase text-[10px] tracking-wider border-b border-slate-800">
              <tr>
                <th class="py-3 px-4">Order ID & Date</th>
                <th class="py-3 px-4">Buyer & Rail</th>
                <th class="py-3 px-4">Procured Items</th>
                <th class="py-3 px-4">Total Amount</th>
                <th class="py-3 px-4">Escrow Status</th>
                <th class="py-3 px-4">Fulfillment / Issued Keys</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800">
              <tr v-for="order in orders" :key="order.id" class="hover:bg-slate-800/40 transition-colors">
                <td class="py-3.5 px-4">
                  <div class="font-bold text-white font-mono">{{ order.order_number }}</div>
                  <div class="text-[10px] text-slate-500">{{ new Date(order.created_at).toLocaleDateString() }}</div>
                </td>
                <td class="py-3.5 px-4">
                  <div class="font-bold text-slate-200">{{ order.customer_name }}</div>
                  <div class="text-[10px] text-amber-400 font-semibold">{{ order.payment_method }}</div>
                </td>
                <td class="py-3.5 px-4">
                  <div v-for="item in order.items" :key="item.product_id" class="text-slate-300">
                    {{ item.quantity }}x {{ item.product_title }}
                  </div>
                </td>
                <td class="py-3.5 px-4 font-black text-amber-400">
                  {{ formatPrice(order.total_amount) }}
                </td>
                <td class="py-3.5 px-4">
                  <span class="px-2 py-0.5 rounded text-[10px] font-bold uppercase bg-emerald-950 text-emerald-400 border border-emerald-800">
                    {{ order.status }}
                  </span>
                </td>
                <td class="py-3.5 px-4">
                  <div v-for="item in order.items" :key="item.product_id">
                    <div v-if="item.digital_license_key" class="p-2 rounded bg-slate-950 border border-slate-800 space-y-1">
                      <div class="text-[10px] text-cyan-400 font-mono select-all flex items-center justify-between">
                        <span>Key: {{ item.digital_license_key }}</span>
                        <button @click="copyKey(item.digital_license_key)" class="text-slate-400 hover:text-white text-[9px] underline">
                          Copy
                        </button>
                      </div>
                      <a :href="item.digital_download_url || '#'" target="_blank" class="text-[10px] text-amber-400 hover:underline">
                        Download CDN Asset Token &rarr;
                      </a>
                    </div>
                    <div v-else class="text-slate-400 text-[11px]">
                      Bill of Lading / Freight Dispatch Waybill Issued
                    </div>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Tab 3: B2B RFQ Marketplace Bidding -->
    <div v-if="currentTab === 'rfqs'" class="space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-base font-bold text-white">B2B Request for Quotation (RFQ) Tender Board</h3>
          <p class="text-xs text-slate-400">Merchants can submit bids with custom lead times & Incoterms</p>
        </div>
      </div>

      <div class="rounded-2xl bg-slate-900 border border-slate-800 overflow-hidden shadow-xl">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-slate-950 text-slate-400 uppercase text-[10px] tracking-wider border-b border-slate-800">
              <tr>
                <th class="py-3 px-4">RFQ Ref & Brief</th>
                <th class="py-3 px-4">Buyer Company</th>
                <th class="py-3 px-4">Quantity</th>
                <th class="py-3 px-4">Target Unit Price</th>
                <th class="py-3 px-4">Destination</th>
                <th class="py-3 px-4 text-right">Submit Quotation</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800">
              <tr v-for="r in rfqs" :key="r.id" class="hover:bg-slate-800/40 transition-colors">
                <td class="py-3.5 px-4">
                  <div class="font-bold text-white">{{ r.title }}</div>
                  <div class="text-[10px] text-slate-500 font-mono">{{ r.rfq_number }}</div>
                  <div class="text-[11px] text-slate-400 mt-1 max-w-sm">{{ r.specifications }}</div>
                </td>
                <td class="py-3.5 px-4">
                  <div class="font-bold text-slate-200">{{ r.buyer_company || r.buyer_name }}</div>
                  <div class="text-[10px] text-slate-400">{{ r.buyer_country }}</div>
                </td>
                <td class="py-3.5 px-4 font-bold text-white">{{ r.quantity.toLocaleString() }} pcs</td>
                <td class="py-3.5 px-4 font-bold text-amber-400 font-mono">{{ formatPrice(r.target_price) }}</td>
                <td class="py-3.5 px-4 text-slate-400">{{ r.destination_country }}</td>
                <td class="py-3.5 px-4 text-right">
                  <button
                    @click="openQuoteModal(r)"
                    class="px-3 py-1.5 rounded-lg bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs"
                  >
                    Bid on RFQ
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Tab 4: Inbound Buyer Offers (eBay Feature) -->
    <div v-if="currentTab === 'offers'" class="space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-base font-bold text-white">Inbound Buyer Offers & Price Negotiations</h3>
          <p class="text-xs text-slate-400">Buyers proposing custom prices on Make an Offer items</p>
        </div>
      </div>

      <div class="rounded-2xl bg-slate-900 border border-slate-800 overflow-hidden shadow-xl">
        <div v-if="offers.length === 0" class="text-center py-12 text-xs text-slate-400">
          No active buyer offers currently submitted.
        </div>
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-slate-950 text-slate-400 uppercase text-[10px] tracking-wider border-b border-slate-800">
              <tr>
                <th class="py-3 px-4">Product</th>
                <th class="py-3 px-4">Buyer Entity</th>
                <th class="py-3 px-4">Offered Price / Qty</th>
                <th class="py-3 px-4">Total Value</th>
                <th class="py-3 px-4">Buyer Note</th>
                <th class="py-3 px-4 text-right">Decision</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800">
              <tr v-for="o in offers" :key="o.id" class="hover:bg-slate-800/40 transition-colors">
                <td class="py-3.5 px-4 font-bold text-white">{{ o.product_title }}</td>
                <td class="py-3.5 px-4">
                  <div class="text-slate-200">{{ o.buyer_name }}</div>
                  <div class="text-[10px] text-slate-500">{{ o.buyer_email }}</div>
                </td>
                <td class="py-3.5 px-4">
                  <div class="font-bold text-purple-400 font-mono">{{ formatPrice(o.offered_price) }} / unit</div>
                  <div class="text-[10px] text-slate-400">Qty: {{ o.quantity }}</div>
                </td>
                <td class="py-3.5 px-4 font-black text-amber-400 font-mono">
                  {{ formatPrice(o.total_offered) }}
                </td>
                <td class="py-3.5 px-4 text-slate-400 text-[11px] max-w-xs truncate">{{ o.message || 'No extra note' }}</td>
                <td class="py-3.5 px-4 text-right space-x-2">
                  <button class="px-2.5 py-1 rounded bg-emerald-950 text-emerald-400 border border-emerald-800 font-bold">
                    Accept
                  </button>
                  <button class="px-2.5 py-1 rounded bg-rose-950 text-rose-400 border border-rose-800 font-bold">
                    Decline
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Tab 5: Verified Directory (Admin) -->
    <div v-if="currentTab === 'users' && currentUser.role === 'admin'" class="space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="text-base font-bold text-white">Global Merchant Entities & Regional Compliance</h3>
        <span class="text-xs text-slate-400">Super Admin Access</span>
      </div>

      <div class="rounded-2xl bg-slate-900 border border-slate-800 overflow-hidden shadow-xl">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-slate-950 text-slate-400 uppercase text-[10px] tracking-wider border-b border-slate-800">
              <tr>
                <th class="py-3 px-4">User & Email</th>
                <th class="py-3 px-4">Role</th>
                <th class="py-3 px-4">Company Name</th>
                <th class="py-3 px-4">Jurisdiction</th>
                <th class="py-3 px-4">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800">
              <tr v-for="u in users" :key="u.id" class="hover:bg-slate-800/40 transition-colors">
                <td class="py-3.5 px-4">
                  <div class="font-bold text-white">{{ u.name }}</div>
                  <div class="text-[10px] text-slate-400 font-mono">{{ u.email }}</div>
                </td>
                <td class="py-3.5 px-4 uppercase font-bold text-[10px] text-amber-400">{{ u.role }}</td>
                <td class="py-3.5 px-4 text-slate-300">{{ u.company_name || 'N/A' }}</td>
                <td class="py-3.5 px-4 text-slate-400">{{ u.country || 'Global' }}</td>
                <td class="py-3.5 px-4">
                  <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-950 text-emerald-400 border border-emerald-800">
                    Verified Gold
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Modal: Create Enterprise Listing (Hybrid Alibaba + Amazon + eBay) -->
    <div
      v-if="showAddModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm"
    >
      <div class="bg-slate-900 border border-slate-700 rounded-3xl max-w-2xl w-full p-6 shadow-2xl space-y-6 relative max-h-[90vh] overflow-y-auto">
        <button @click="showAddModal = false" class="absolute top-5 right-5 text-slate-400 hover:text-white text-lg font-bold">
          &times;
        </button>

        <div>
          <h2 class="text-xl font-black text-white">Create Enterprise Marketplace Listing</h2>
          <p class="text-xs text-slate-400">Configure Alibaba-style B2B volume tiers, Amazon-style fast buy, or instant digital licenses</p>
        </div>

        <form @submit.prevent="submitProduct" class="space-y-4">
          <!-- Format Selector -->
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Marketplace Model & Format</label>
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
              <button
                type="button"
                @click="productForm.format = 'wholesale'"
                :class="[
                  'p-2.5 rounded-xl border text-left transition-all text-xs font-bold',
                  productForm.format === 'wholesale' ? 'bg-amber-500/10 border-amber-500/40 text-amber-300' : 'bg-slate-950 border-slate-800 text-slate-400'
                ]"
              >
                🚢 B2B Wholesale
              </button>
              <button
                type="button"
                @click="productForm.format = 'retail'"
                :class="[
                  'p-2.5 rounded-xl border text-left transition-all text-xs font-bold',
                  productForm.format === 'retail' ? 'bg-emerald-500/10 border-emerald-500/40 text-emerald-300' : 'bg-slate-950 border-slate-800 text-slate-400'
                ]"
              >
                🛒 Fast Buy
              </button>
              <button
                type="button"
                @click="productForm.format = 'digital'; productForm.type = 'digital_software'"
                :class="[
                  'p-2.5 rounded-xl border text-left transition-all text-xs font-bold',
                  productForm.format === 'digital' ? 'bg-cyan-500/10 border-cyan-500/40 text-cyan-300' : 'bg-slate-950 border-slate-800 text-slate-400'
                ]"
              >
                ⚡ Instant Digital
              </button>
              <button
                type="button"
                @click="productForm.format = 'make_offer'; productForm.accepts_offers = true"
                :class="[
                  'p-2.5 rounded-xl border text-left transition-all text-xs font-bold',
                  productForm.format === 'make_offer' ? 'bg-purple-500/10 border-purple-500/40 text-purple-300' : 'bg-slate-950 border-slate-800 text-slate-400'
                ]"
              >
                🤝 Make Offer
              </button>
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Product Title</label>
              <input
                v-model="productForm.title"
                type="text"
                required
                placeholder="e.g. Industrial 5-Axis CNC Mill"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Category</label>
              <select
                v-model="productForm.category"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              >
                <option value="Industrial & Machinery">Industrial & Heavy CNC</option>
                <option value="Electronics & Telecom">Datacenter & Networking</option>
                <option value="Digital Software">Digital Software & SaaS</option>
                <option value="Digital Assets">AI Models & Weights</option>
                <option value="Agriculture & Processing">Agri-Tech & Sorting Mills</option>
                <option value="Energy & Utilities">Commercial Energy & BESS</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Base Price ($ USD)</label>
              <input
                v-model.number="productForm.price"
                type="number"
                step="0.01"
                required
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Bulk Wholesale Price ($ USD)</label>
              <input
                v-model.number="productForm.wholesale_price"
                type="number"
                step="0.01"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Target Region</label>
              <select
                v-model="productForm.region"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              >
                <option value="global">Global International</option>
                <option value="north_america">North America</option>
                <option value="europe">Europe & UK</option>
                <option value="east_africa">East Africa Hub</option>
                <option value="west_africa">West Africa Hub</option>
                <option value="asia_pacific">Asia-Pacific</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Condition</label>
              <select
                v-model="productForm.condition"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              >
                <option value="new">Brand New</option>
                <option value="refurbished">Certified Refurbished</option>
                <option value="open_box">Open Box / Surplus</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Min. Order (MOQ)</label>
              <input
                v-model.number="productForm.moq"
                type="number"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Description & Specifications</label>
            <textarea
              v-model="productForm.description"
              rows="2"
              placeholder="Detailed specifications, factory compliance, warranty..."
              class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
            ></textarea>
          </div>

          <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-800">
            <button type="button" @click="showAddModal = false" class="px-4 py-2 rounded-xl bg-slate-800 text-xs text-slate-300">
              Cancel
            </button>
            <button type="submit" class="px-5 py-2.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs shadow-lg">
              Publish to Marketplace
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal: Submit Quotation on RFQ (Supplier) -->
    <div
      v-if="activeRfqForQuote"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm"
    >
      <div class="bg-slate-900 border border-slate-700 rounded-3xl max-w-lg w-full p-6 shadow-2xl space-y-6 relative">
        <button @click="activeRfqForQuote = null" class="absolute top-5 right-5 text-slate-400 hover:text-white text-lg font-bold">
          &times;
        </button>

        <div>
          <div class="text-xs font-bold text-amber-400 uppercase tracking-wider">📋 Submit Supplier Quotation</div>
          <h3 class="text-base font-bold text-white mt-1">{{ activeRfqForQuote.title }}</h3>
          <p class="text-xs text-slate-400">Target Qty: {{ activeRfqForQuote.quantity }} · Target: ${{ activeRfqForQuote.target_price }}</p>
        </div>

        <form @submit.prevent="submitRfqQuoteAction" class="space-y-4">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Your Quoted Unit Price ($ USD)</label>
              <input
                v-model.number="quoteUnitPrice"
                type="number"
                step="0.01"
                required
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Production Lead Time (Days)</label>
              <input
                v-model.number="quoteLeadTime"
                type="number"
                required
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Quotation Notes & Incoterms</label>
            <textarea
              v-model="quoteNotes"
              rows="2"
              placeholder="e.g. Includes CIF shipping insurance and CE certificates..."
              class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white outline-none"
            ></textarea>
          </div>

          <button
            type="submit"
            class="w-full py-3 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs shadow-lg transition-all"
          >
            Submit Official Quotation
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Link, router } from '@inertiajs/vue3'
import { ref } from 'vue'
import { useCurrency, type Currency } from '@/composables/useCurrency'

const props = defineProps<{
  currentUser: any
  stats: any
  products: any[]
  orders: any[]
  users: any[]
  rfqs: any[]
  offers: any[]
  currencies?: Currency[]
}>()

const { formatPrice } = useCurrency(props.currencies)

const currentTab = ref('catalog')
const showAddModal = ref(false)

const productForm = ref({
  title: '',
  category: 'Industrial & Machinery',
  sub_category: 'Precision Hardware',
  type: 'physical',
  format: 'wholesale',
  condition: 'new',
  region: 'global',
  price: 980.00,
  wholesale_price: 780.00,
  moq: 5,
  stock: 120,
  sku: '',
  description: '',
  image_url: 'https://images.unsplash.com/photo-1581092160607-ee22621dd758?w=800&auto=format&fit=crop&q=80',
  digital_license_format: 'SOKO-LIC-XXXX-XXXX',
  digital_file_size: '420 MB',
  digital_download_url: 'https://cdn.soko.global/releases/asset.zip',
  accepts_offers: true,
})

const activeRfqForQuote = ref<any | null>(null)
const quoteUnitPrice = ref(0)
const quoteLeadTime = ref(7)
const quoteNotes = ref('')

const submitProduct = () => {
  router.post('/dashboard/products', productForm.value, {
    onSuccess: () => {
      showAddModal.value = false
    },
  })
}

const deleteProduct = (id: string) => {
  if (confirm('Delete this listing from the enterprise catalog?')) {
    router.post(`/dashboard/products/${id}/delete`)
  }
}

const copyKey = (key: string) => {
  navigator.clipboard.writeText(key)
  alert('Digital License Key copied to clipboard: ' + key)
}

const openQuoteModal = (rfq: any) => {
  activeRfqForQuote.value = rfq
  quoteUnitPrice.value = rfq.target_price || 50
}

const submitRfqQuoteAction = () => {
  if (!activeRfqForQuote.value) return
  router.post('/dashboard/rfqs/quote', {
    rfq_id: activeRfqForQuote.value.id,
    unit_price: quoteUnitPrice.value,
    lead_time_days: quoteLeadTime.value,
    notes: quoteNotes.value,
  }, {
    onFinish: () => {
      activeRfqForQuote.value = null
    },
  })
}
</script>
