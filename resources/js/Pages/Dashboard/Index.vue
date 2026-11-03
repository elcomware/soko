<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 py-8 space-y-8">
    <!-- Dashboard Top Profile Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-3xl bg-slate-900 border border-slate-800 shadow-xl">
      <div class="flex items-center gap-4">
        <img
          :src="currentUser.avatar || 'https://ui-avatars.com/api/?name=' + currentUser.name"
          class="w-16 h-16 rounded-2xl object-cover border-2 border-amber-500/50 shadow-md"
          alt="Avatar"
        />
        <div>
          <div class="flex items-center gap-2">
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
          </div>
          <div class="text-xs text-slate-400 mt-1 flex items-center gap-2">
            <span>{{ currentUser.email }}</span>
            <span>·</span>
            <span>Region: {{ currentUser.country || 'Global' }}</span>
            <span>·</span>
            <span class="text-emerald-400 font-medium">✓ Verified Enterprise Entity</span>
          </div>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <Link
          href="/"
          class="px-4 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs font-semibold text-slate-200 transition-colors"
        >
          View Live eStore
        </Link>
        <button
          @click="showAddModal = true"
          class="px-4 py-2.5 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-bold text-xs shadow-lg shadow-amber-500/20 transition-all flex items-center gap-1.5"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4" />
          </svg>
          <span>List New Product</span>
        </button>
      </div>
    </div>

    <!-- KPI Metric Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="p-5 rounded-2xl bg-slate-900 border border-slate-800 space-y-2">
        <div class="flex items-center justify-between text-xs text-slate-400">
          <span>Gross Volume (GMV)</span>
          <span class="text-emerald-400 font-bold">+18.4%</span>
        </div>
        <div class="text-2xl font-black text-white">${{ (stats.total_gmv || 0).toLocaleString() }}</div>
        <div class="text-[11px] text-slate-500">Processed through Soko Escrow</div>
      </div>

      <div class="p-5 rounded-2xl bg-slate-900 border border-slate-800 space-y-2">
        <div class="flex items-center justify-between text-xs text-slate-400">
          <span>Total Orders</span>
          <span class="text-amber-400 font-bold">Live</span>
        </div>
        <div class="text-2xl font-black text-white">{{ stats.total_orders }}</div>
        <div class="text-[11px] text-slate-500">Wholesale freight & digital assets</div>
      </div>

      <div class="p-5 rounded-2xl bg-slate-900 border border-slate-800 space-y-2">
        <div class="flex items-center justify-between text-xs text-slate-400">
          <span>Digital Deliveries</span>
          <span class="text-cyan-400 font-mono text-[10px]">Zero-Latency</span>
        </div>
        <div class="text-2xl font-black text-cyan-400">{{ stats.instant_deliveries }}</div>
        <div class="text-[11px] text-slate-500">Instant licenses & download tokens</div>
      </div>

      <div class="p-5 rounded-2xl bg-slate-900 border border-slate-800 space-y-2">
        <div class="flex items-center justify-between text-xs text-slate-400">
          <span>Catalog Inventory</span>
          <span class="text-slate-300 font-bold">{{ products.length }} Active</span>
        </div>
        <div class="text-2xl font-black text-white">
          {{ stats.physical_products_count }} <span class="text-xs text-amber-400 font-normal">Physical</span> /
          {{ stats.digital_products_count }} <span class="text-xs text-cyan-400 font-normal">Digital</span>
        </div>
        <div class="text-[11px] text-slate-500">Distributed across global nodes</div>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="border-b border-slate-800 flex items-center gap-2">
      <button
        @click="currentTab = 'products'"
        :class="[
          'px-4 py-3 text-xs font-bold border-b-2 transition-all flex items-center gap-2',
          currentTab === 'products'
            ? 'border-amber-500 text-amber-400'
            : 'border-transparent text-slate-400 hover:text-slate-200'
        ]"
      >
        <span>📦 Product Catalog</span>
        <span class="px-2 py-0.5 rounded-full bg-slate-800 text-[10px] text-slate-300">{{ products.length }}</span>
      </button>

      <button
        @click="currentTab = 'orders'"
        :class="[
          'px-4 py-3 text-xs font-bold border-b-2 transition-all flex items-center gap-2',
          currentTab === 'orders'
            ? 'border-amber-500 text-amber-400'
            : 'border-transparent text-slate-400 hover:text-slate-200'
        ]"
      >
        <span>📜 Orders & Digital Licenses</span>
        <span class="px-2 py-0.5 rounded-full bg-slate-800 text-[10px] text-slate-300">{{ orders.length }}</span>
      </button>

      <button
        v-if="currentUser.role === 'admin'"
        @click="currentTab = 'users'"
        :class="[
          'px-4 py-3 text-xs font-bold border-b-2 transition-all flex items-center gap-2',
          currentTab === 'users'
            ? 'border-amber-500 text-amber-400'
            : 'border-transparent text-slate-400 hover:text-slate-200'
        ]"
      >
        <span>👥 Verified Suppliers & Directory</span>
        <span class="px-2 py-0.5 rounded-full bg-slate-800 text-[10px] text-slate-300">{{ users.length }}</span>
      </button>
    </div>

    <!-- Tab 1: Product Catalog Management -->
    <div v-if="currentTab === 'products'" class="space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-bold text-white">Enterprise Catalog Listings</h3>
        <button
          @click="showAddModal = true"
          class="px-3.5 py-1.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs transition-colors"
        >
          + Add Product
        </button>
      </div>

      <div class="rounded-2xl bg-slate-900 border border-slate-800 overflow-hidden shadow-xl">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-slate-950/80 text-slate-400 uppercase tracking-wider text-[10px] border-b border-slate-800">
              <tr>
                <th class="py-3.5 px-4">Item & SKU</th>
                <th class="py-3.5 px-4">Type</th>
                <th class="py-3.5 px-4">Category</th>
                <th class="py-3.5 px-4">Pricing</th>
                <th class="py-3.5 px-4">Fulfillment Details</th>
                <th class="py-3.5 px-4">Supplier</th>
                <th class="py-3.5 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800/80">
              <tr v-for="product in products" :key="product.id" class="hover:bg-slate-800/40 transition-colors">
                <td class="py-3 px-4">
                  <div class="flex items-center gap-3">
                    <img :src="product.image_url" class="w-10 h-10 rounded-lg object-cover border border-slate-800 shrink-0" alt="Thumb" />
                    <div>
                      <div class="font-bold text-white max-w-xs truncate">{{ product.title }}</div>
                      <div class="font-mono text-[10px] text-slate-500">{{ product.sku }}</div>
                    </div>
                  </div>
                </td>
                <td class="py-3 px-4">
                  <span
                    v-if="product.type !== 'physical'"
                    class="px-2 py-0.5 rounded text-[10px] font-bold uppercase bg-cyan-950 text-cyan-400 border border-cyan-800"
                  >
                    ⚡ Digital
                  </span>
                  <span
                    v-else
                    class="px-2 py-0.5 rounded text-[10px] font-bold uppercase bg-amber-950 text-amber-400 border border-amber-800"
                  >
                    🚢 Freight
                  </span>
                </td>
                <td class="py-3 px-4 text-slate-300">
                  {{ product.category }}
                </td>
                <td class="py-3 px-4">
                  <div class="font-bold text-white">${{ product.price.toFixed(2) }}</div>
                  <div v-if="product.wholesale_price" class="text-[10px] text-amber-400">
                    Tier: ${{ product.wholesale_price.toFixed(2) }}
                  </div>
                </td>
                <td class="py-3 px-4">
                  <div v-if="product.type !== 'physical'" class="text-cyan-400 font-mono text-[11px]">
                    Size: {{ product.digital_file_size || 'Instant' }}
                  </div>
                  <div v-else class="text-slate-300">
                    MOQ: {{ product.moq }} · {{ product.lead_time_days }}d dispatch
                  </div>
                </td>
                <td class="py-3 px-4 text-slate-400">
                  {{ product.vendor_name }}
                </td>
                <td class="py-3 px-4 text-right space-x-2">
                  <Link
                    :href="'/products/' + product.id"
                    class="px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 font-semibold"
                  >
                    View
                  </Link>
                  <button
                    @click="deleteProduct(product.id)"
                    class="px-2.5 py-1 rounded bg-rose-950/60 hover:bg-rose-900 border border-rose-800/40 text-rose-300 font-semibold"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Tab 2: Orders & Digital Key Fulfillment Tracking -->
    <div v-if="currentTab === 'orders'" class="space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-bold text-white">Order History & Digital License Audit</h3>
        <span class="text-xs text-slate-400">Escrow verification & license tokens</span>
      </div>

      <div class="rounded-2xl bg-slate-900 border border-slate-800 overflow-hidden shadow-xl">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-slate-950/80 text-slate-400 uppercase tracking-wider text-[10px] border-b border-slate-800">
              <tr>
                <th class="py-3.5 px-4">Order ID & Date</th>
                <th class="py-3.5 px-4">Customer Details</th>
                <th class="py-3.5 px-4">Items Procured</th>
                <th class="py-3.5 px-4">Total Amount</th>
                <th class="py-3.5 px-4">Status & Escrow</th>
                <th class="py-3.5 px-4">Fulfillment / Digital Keys</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800/80">
              <tr v-for="order in orders" :key="order.id" class="hover:bg-slate-800/40 transition-colors">
                <td class="py-3 px-4">
                  <div class="font-bold text-white font-mono">{{ order.order_number }}</div>
                  <div class="text-[10px] text-slate-500">{{ new Date(order.created_at).toLocaleDateString() }}</div>
                </td>
                <td class="py-3 px-4">
                  <div class="font-bold text-slate-200">{{ order.customer_name }}</div>
                  <div class="text-[10px] text-slate-400">{{ order.customer_email }}</div>
                </td>
                <td class="py-3 px-4">
                  <div v-for="item in order.items" :key="item.product_id" class="space-y-0.5">
                    <span class="text-slate-300 font-medium">{{ item.quantity }}x {{ item.product_title }}</span>
                  </div>
                </td>
                <td class="py-3 px-4 font-black text-amber-400">
                  ${{ order.total_amount.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}
                </td>
                <td class="py-3 px-4">
                  <span
                    :class="[
                      'px-2.5 py-1 rounded-full text-[10px] font-bold uppercase',
                      order.status === 'completed' ? 'bg-emerald-950 text-emerald-400 border border-emerald-800' : 'bg-amber-950 text-amber-400 border border-amber-800'
                    ]"
                  >
                    {{ order.status }}
                  </span>
                </td>
                <td class="py-3 px-4">
                  <div v-for="item in order.items" :key="item.product_id">
                    <div v-if="item.digital_license_key" class="p-2 rounded bg-slate-950 border border-slate-800 space-y-1">
                      <div class="text-[10px] text-cyan-400 font-mono select-all">
                        Key: {{ item.digital_license_key }}
                      </div>
                      <a
                        :href="item.digital_download_url || '#'"
                        target="_blank"
                        class="text-[10px] text-amber-400 hover:underline flex items-center gap-1"
                      >
                        <span>Download Asset Token &rarr;</span>
                      </a>
                    </div>
                    <div v-else class="text-slate-400 text-[11px]">
                      Freight Bill of Lading Generated
                    </div>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Tab 3: Users & Verified Suppliers Directory (Admin Only) -->
    <div v-if="currentTab === 'users' && currentUser.role === 'admin'" class="space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-bold text-white">Verified Enterprise Directory</h3>
        <span class="text-xs text-slate-400">Super Admin Access</span>
      </div>

      <div class="rounded-2xl bg-slate-900 border border-slate-800 overflow-hidden shadow-xl">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-slate-950/80 text-slate-400 uppercase tracking-wider text-[10px] border-b border-slate-800">
              <tr>
                <th class="py-3.5 px-4">Entity & Contact</th>
                <th class="py-3.5 px-4">Role</th>
                <th class="py-3.5 px-4">Legal Company</th>
                <th class="py-3.5 px-4">Jurisdiction</th>
                <th class="py-3.5 px-4">Verification</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800/80">
              <tr v-for="u in users" :key="u.id" class="hover:bg-slate-800/40 transition-colors">
                <td class="py-3 px-4">
                  <div class="font-bold text-white">{{ u.name }}</div>
                  <div class="text-[10px] text-slate-400 font-mono">{{ u.email }}</div>
                </td>
                <td class="py-3 px-4 uppercase font-bold text-[10px] text-amber-400">
                  {{ u.role }}
                </td>
                <td class="py-3 px-4 text-slate-300">
                  {{ u.company_name || 'N/A' }}
                </td>
                <td class="py-3 px-4 text-slate-400">
                  {{ u.country || 'Global' }}
                </td>
                <td class="py-3 px-4">
                  <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-950 text-emerald-400 border border-emerald-800">
                    Verified
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Modal: Add New Enterprise Product -->
    <div
      v-if="showAddModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm"
    >
      <div class="bg-slate-900 border border-slate-700 rounded-3xl max-w-2xl w-full p-6 shadow-2xl space-y-6 relative max-h-[90vh] overflow-y-auto">
        <button
          @click="showAddModal = false"
          class="absolute top-5 right-5 text-slate-400 hover:text-white text-lg font-bold"
        >
          &times;
        </button>

        <div>
          <h2 class="text-xl font-black text-white">List New Product on Soko</h2>
          <p class="text-xs text-slate-400">Configure physical wholesale freight or zero-latency digital software licenses</p>
        </div>

        <form @submit.prevent="submitProduct" class="space-y-4">
          <!-- Type Selector -->
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Product Classification</label>
            <div class="grid grid-cols-2 gap-3">
              <button
                type="button"
                @click="productForm.type = 'digital_software'"
                :class="[
                  'p-3 rounded-xl border text-left transition-all',
                  productForm.type !== 'physical'
                    ? 'bg-cyan-500/10 border-cyan-500/40 text-cyan-300'
                    : 'bg-slate-950 border-slate-800 text-slate-400'
                ]"
              >
                <div class="font-bold text-xs">⚡ Digital Software / Asset / License</div>
                <div class="text-[10px] text-slate-500">Automated key dispatch & instant CDN download</div>
              </button>

              <button
                type="button"
                @click="productForm.type = 'physical'"
                :class="[
                  'p-3 rounded-xl border text-left transition-all',
                  productForm.type === 'physical'
                    ? 'bg-amber-500/10 border-amber-500/40 text-amber-300'
                    : 'bg-slate-950 border-slate-800 text-slate-400'
                ]"
              >
                <div class="font-bold text-xs">🚢 Physical Commodity / Heavy Freight</div>
                <div class="text-[10px] text-slate-500">B2B wholesale volume tiers, freight lead times & MOQ</div>
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
                placeholder="e.g. Soko Quantum Analytics Engine"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white focus:border-amber-500 outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Category</label>
              <select
                v-model="productForm.category"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white focus:border-amber-500 outline-none"
              >
                <option value="Digital Software">Digital Software & SaaS</option>
                <option value="Digital Assets">Digital Assets & AI Models</option>
                <option value="Industrial & Machinery">Industrial & Heavy CNC</option>
                <option value="Electronics & Telecom">Datacenter & Networking</option>
                <option value="Energy & Utilities">Commercial Energy & Storage</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Retail Price ($)</label>
              <input
                v-model.number="productForm.price"
                type="number"
                step="0.01"
                required
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white focus:border-amber-500 outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Wholesale Tier Price ($)</label>
              <input
                v-model.number="productForm.wholesale_price"
                type="number"
                step="0.01"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white focus:border-amber-500 outline-none"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Min. Order (MOQ)</label>
              <input
                v-model.number="productForm.moq"
                type="number"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white focus:border-amber-500 outline-none"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Description</label>
            <textarea
              v-model="productForm.description"
              rows="3"
              placeholder="Enterprise specifications, SLA guarantees, or material certifications..."
              class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2 text-sm text-white focus:border-amber-500 outline-none"
            ></textarea>
          </div>

          <!-- Conditional fields: Digital vs Physical -->
          <div v-if="productForm.type !== 'physical'" class="p-3.5 rounded-2xl bg-slate-950 border border-cyan-500/30 space-y-3">
            <div class="text-xs font-bold text-cyan-400">⚡ Digital Fulfillment Settings</div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block text-[11px] text-slate-400 mb-1">License Key Format</label>
                <input
                  v-model="productForm.digital_license_format"
                  type="text"
                  placeholder="SOKO-ENTERPRISE-XXXX-XXXX"
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-white outline-none"
                />
              </div>
              <div>
                <label class="block text-[11px] text-slate-400 mb-1">File / Asset Size</label>
                <input
                  v-model="productForm.digital_file_size"
                  type="text"
                  placeholder="e.g. 1.2 GB"
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-white outline-none"
                />
              </div>
            </div>
          </div>

          <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-800">
            <button
              type="button"
              @click="showAddModal = false"
              class="px-4 py-2.5 rounded-xl bg-slate-800 text-xs font-semibold text-slate-300 hover:text-white"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="px-5 py-2.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs shadow-lg shadow-amber-500/20"
            >
              Publish to Enterprise Catalog
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Link, router } from '@inertiajs/vue3'
import { ref } from 'vue'

const props = defineProps<{
  currentUser: any
  stats: any
  products: any[]
  orders: any[]
  users: any[]
}>()

const currentTab = ref('products')
const showAddModal = ref(false)

const productForm = ref({
  title: '',
  category: 'Digital Software',
  sub_category: 'Enterprise Suite',
  type: 'digital_software',
  price: 499.00,
  wholesale_price: 399.00,
  moq: 1,
  stock: 100,
  sku: '',
  description: '',
  image_url: 'https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=800&auto=format&fit=crop&q=80',
  digital_license_format: 'SOKO-LIC-XXXX-ENTERPRISE',
  digital_file_size: '650 MB',
  digital_download_url: 'https://cdn.soko.global/releases/enterprise-pkg.zip',
})

const submitProduct = () => {
  router.post('/dashboard/products', productForm.value, {
    onSuccess: () => {
      showAddModal.value = false
    },
  })
}

const deleteProduct = (id: string) => {
  if (confirm('Are you sure you want to remove this product from the enterprise catalog?')) {
    router.post(`/dashboard/products/${id}/delete`)
  }
}
</script>

