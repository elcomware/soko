<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 py-8 space-y-12">
    <!-- Breadcrumb -->
    <nav class="flex items-center gap-2 text-xs text-slate-400">
      <Link href="/" class="hover:text-amber-400 transition-colors">Marketplace</Link>
      <span>/</span>
      <Link :href="'/?category=' + encodeURIComponent(product.category)" class="hover:text-amber-400 transition-colors">
        {{ product.category }}
      </Link>
      <span>/</span>
      <span class="text-slate-200 font-medium truncate">{{ product.title }}</span>
    </nav>

    <!-- Main Product Showcase -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
      <!-- Left: Visual & Media -->
      <div class="lg:col-span-7 space-y-4">
        <div class="rounded-3xl bg-slate-900 border border-slate-800 overflow-hidden relative shadow-2xl">
          <img
            :src="product.image_url"
            :alt="product.title"
            class="w-full h-96 sm:h-[450px] object-cover"
          />
          <div class="absolute top-4 left-4">
            <span
              v-if="product.type !== 'physical'"
              class="px-3 py-1.5 rounded-full text-xs font-black uppercase tracking-wider bg-cyan-950/90 text-cyan-300 border border-cyan-500/40 backdrop-blur flex items-center gap-1.5 shadow-lg"
            >
              <span class="w-2 h-2 rounded-full bg-cyan-400 animate-pulse"></span>
              Instant Digital Delivery
            </span>
            <span
              v-else
              class="px-3 py-1.5 rounded-full text-xs font-black uppercase tracking-wider bg-amber-950/90 text-amber-300 border border-amber-500/40 backdrop-blur shadow-lg"
            >
              🚢 Freight Cargo · MOQ {{ product.moq }}
            </span>
          </div>
        </div>

        <!-- Supplier Trust Card -->
        <div class="p-5 rounded-2xl bg-slate-900/80 border border-slate-800 flex items-center justify-between gap-4">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 rounded-xl bg-slate-800 flex items-center justify-center font-bold text-amber-400 text-lg border border-slate-700">
              🏢
            </div>
            <div>
              <div class="flex items-center gap-1.5">
                <span class="font-bold text-white text-sm">{{ product.vendor_name }}</span>
                <span class="px-1.5 py-0.5 rounded bg-amber-500/10 text-amber-400 text-[10px] font-bold border border-amber-500/20">
                  Verified Supplier
                </span>
              </div>
              <div class="text-xs text-slate-400 mt-0.5">Location: {{ product.vendor_country }} · 99.4% Fulfillment Rate</div>
            </div>
          </div>
          <button class="px-3.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs font-semibold text-slate-200 transition-colors">
            Contact Supplier
          </button>
        </div>
      </div>

      <!-- Right: Pricing & Order Execution -->
      <div class="lg:col-span-5 space-y-6">
        <div class="rounded-3xl bg-slate-900 border border-slate-800 p-6 space-y-6 shadow-xl">
          <div>
            <div class="text-xs font-mono text-slate-500 mb-1">SKU: {{ product.sku }}</div>
            <h1 class="text-2xl font-black text-white leading-snug">{{ product.title }}</h1>
            <p class="text-xs text-slate-400 mt-2 leading-relaxed">{{ product.description }}</p>
          </div>

          <!-- Pricing Matrix -->
          <div class="p-4 rounded-2xl bg-slate-950/80 border border-slate-800/80 space-y-2">
            <div class="flex items-baseline justify-between">
              <span class="text-xs text-slate-400">Standard Price:</span>
              <span class="text-2xl font-black text-white">
                ${{ product.price.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}
              </span>
            </div>
            <div v-if="product.wholesale_price" class="flex items-baseline justify-between pt-2 border-t border-slate-800 text-amber-400">
              <span class="text-xs">Wholesale Tier (10+ units):</span>
              <span class="text-lg font-bold">
                ${{ product.wholesale_price.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}
              </span>
            </div>
          </div>

          <!-- Specifications Table -->
          <div class="space-y-2.5 text-xs">
            <div class="flex justify-between py-1.5 border-b border-slate-800/60">
              <span class="text-slate-400">Category:</span>
              <span class="text-slate-200 font-semibold">{{ product.category }} ({{ product.sub_category }})</span>
            </div>
            <div class="flex justify-between py-1.5 border-b border-slate-800/60">
              <span class="text-slate-400">Minimum Order (MOQ):</span>
              <span class="text-slate-200 font-semibold">{{ product.moq }} unit(s)</span>
            </div>
            <div v-if="product.type !== 'physical'" class="flex justify-between py-1.5 border-b border-slate-800/60">
              <span class="text-slate-400">Asset File Size:</span>
              <span class="text-cyan-400 font-mono font-semibold">{{ product.digital_file_size || 'Instant Access' }}</span>
            </div>
            <div v-if="product.type !== 'physical'" class="flex justify-between py-1.5 border-b border-slate-800/60">
              <span class="text-slate-400">License Protocol:</span>
              <span class="text-cyan-400 font-mono font-semibold">{{ product.digital_license_format }}</span>
            </div>
            <div v-else class="flex justify-between py-1.5 border-b border-slate-800/60">
              <span class="text-slate-400">Dispatch Lead Time:</span>
              <span class="text-amber-400 font-semibold">{{ product.lead_time_days }} Business Days Freight</span>
            </div>
          </div>

          <!-- Quick Checkout Form -->
          <form @submit.prevent="executeCheckout" class="space-y-4 pt-4 border-t border-slate-800">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">
                Order Quantity
              </label>
              <input
                v-model.number="qty"
                type="number"
                :min="product.moq"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2.5 text-sm text-white focus:border-amber-500 outline-none"
              />
            </div>

            <div class="p-3 rounded-xl bg-slate-950 border border-slate-800 flex justify-between items-center text-sm">
              <span class="text-slate-400">Total Escrow Amount:</span>
              <span class="text-xl font-black text-amber-400">${{ computedTotal.toFixed(2) }}</span>
            </div>

            <button
              type="submit"
              class="w-full py-3.5 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-black text-sm shadow-xl shadow-amber-500/20 transition-all"
            >
              <span v-if="product.type !== 'physical'">⚡ Procure License & Download Instantly</span>
              <span v-else>📦 Place Escrow Wholesale Order</span>
            </button>
          </form>
        </div>
      </div>
    </div>

    <!-- Related Products -->
    <div v-if="related && related.length > 0" class="pt-10 border-t border-slate-800">
      <h3 class="text-xl font-bold text-white mb-6">Related Enterprise Solutions</h3>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-6">
        <Link
          v-for="rel in related"
          :key="rel.id"
          :href="'/products/' + rel.id"
          class="group p-4 rounded-2xl bg-slate-900 border border-slate-800 hover:border-slate-700 transition-all flex items-center gap-4"
        >
          <img :src="rel.image_url" class="w-16 h-16 rounded-xl object-cover shrink-0" alt="rel" />
          <div class="truncate">
            <h4 class="font-bold text-xs text-white group-hover:text-amber-400 transition-colors truncate">{{ rel.title }}</h4>
            <div class="text-[11px] text-slate-400 mt-1">${{ rel.price.toLocaleString() }}</div>
          </div>
        </Link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Link, router } from '@inertiajs/vue3'
import { ref, computed } from 'vue'

const props = defineProps<{
  product: any
  related: any[]
}>()

const qty = ref(props.product.moq || 1)

const computedTotal = computed(() => {
  const price = (qty.value >= 10 && props.product.wholesale_price > 0)
    ? props.product.wholesale_price
    : props.product.price
  return price * (qty.value || 1)
})

const executeCheckout = () => {
  router.post('/checkout', {
    product_id: props.product.id,
    quantity: qty.value,
  })
}
</script>

