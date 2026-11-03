<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 py-8 space-y-12">
    <!-- Breadcrumb & Region Tag -->
    <div class="flex items-center justify-between text-xs text-slate-400">
      <nav class="flex items-center gap-2">
        <Link href="/" class="hover:text-amber-400 transition-colors">Marketplace</Link>
        <span>/</span>
        <Link :href="'/?category=' + encodeURIComponent(product.category)" class="hover:text-amber-400 transition-colors">
          {{ product.category }}
        </Link>
        <span>/</span>
        <span class="text-slate-200 font-medium truncate max-w-xs">{{ product.title }}</span>
      </nav>
      <div class="hidden sm:flex items-center gap-2 text-slate-400">
        <span>Sourcing Origin:</span>
        <span class="font-bold text-white">{{ product.vendor_country }} ({{ product.incoterm || 'FOB' }})</span>
      </div>
    </div>

    <!-- Main Product Showcase -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
      <!-- Left: Visual & Media -->
      <div class="lg:col-span-7 space-y-4">
        <div class="rounded-3xl bg-slate-900 border border-slate-800 overflow-hidden relative shadow-2xl">
          <img :src="product.image_url" :alt="product.title" class="w-full h-96 sm:h-[480px] object-cover" />
          <div class="absolute top-4 left-4 flex flex-col gap-1.5">
            <span
              v-if="product.type !== 'physical'"
              class="px-3 py-1.5 rounded-full text-xs font-black uppercase tracking-wider bg-cyan-950/90 text-cyan-300 border border-cyan-500/40 backdrop-blur flex items-center gap-1.5 shadow-lg"
            >
              <span class="w-2 h-2 rounded-full bg-cyan-400 animate-pulse"></span>
              ⚡ Instant Zero-Latency Digital License
            </span>
            <span
              v-else
              class="px-3 py-1.5 rounded-full text-xs font-black uppercase tracking-wider bg-amber-950/90 text-amber-300 border border-amber-500/40 backdrop-blur shadow-lg"
            >
              🚢 B2B Freight · MOQ {{ product.moq }} pcs
            </span>
            <span
              v-if="product.condition === 'refurbished'"
              class="px-2.5 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider bg-slate-900/90 text-slate-300 border border-slate-700 w-fit"
            >
              Certified Refurbished (1-Yr Warranty)
            </span>
          </div>
        </div>

        <!-- Supplier Credibility Card (Alibaba Style) -->
        <div class="p-6 rounded-2xl bg-slate-900 border border-slate-800 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div class="flex items-center gap-4">
            <div class="w-14 h-14 rounded-2xl bg-slate-800 flex items-center justify-center font-bold text-amber-400 text-xl border border-slate-700 shrink-0">
              🏢
            </div>
            <div>
              <div class="flex items-center gap-2">
                <span class="font-bold text-white text-base">{{ product.vendor_name }}</span>
                <span class="px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 text-[10px] font-bold border border-amber-500/20">
                  {{ product.vendor_badge || 'Gold Verified' }}
                </span>
              </div>
              <div class="text-xs text-slate-400 mt-1 flex items-center gap-2">
                <span>{{ product.vendor_country }}</span>
                <span>·</span>
                <span class="text-emerald-400 font-semibold">{{ product.vendor_positive_rate || '99.5%' }} Positive Feedback</span>
                <span>·</span>
                <span class="text-amber-400">${{ (product.trade_assurance_cover / 1000).toFixed(0) }}k Escrow Protection</span>
              </div>
            </div>
          </div>
          <button class="px-4 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs font-bold text-slate-200 transition-colors shrink-0">
            Audit Inspection Report
          </button>
        </div>
      </div>

      <!-- Right: Pricing, Volume Tiers & Order Engine -->
      <div class="lg:col-span-5 space-y-6">
        <div class="rounded-3xl bg-slate-900 border border-slate-800 p-6 space-y-6 shadow-xl">
          <div>
            <div class="text-xs font-mono text-slate-500 mb-1">SKU: {{ product.sku }} · In Stock: {{ product.stock < 0 ? 'Unlimited Digital' : product.stock + ' units' }}</div>
            <h1 class="text-2xl font-black text-white leading-snug">{{ product.title }}</h1>
            <p class="text-xs text-slate-400 mt-2 leading-relaxed">{{ product.description }}</p>
          </div>

          <!-- Alibaba-Style Tiered Price Bracket Matrix -->
          <div v-if="product.price_tiers && product.price_tiers.length > 0" class="space-y-2">
            <div class="text-xs font-bold text-amber-400 uppercase tracking-wider">B2B Volume Tier Pricing ({{ activeCurrency.code }}):</div>
            <div class="grid grid-cols-3 gap-2 bg-slate-950 p-3 rounded-2xl border border-slate-800 text-center">
              <div
                v-for="(tier, idx) in product.price_tiers"
                :key="idx"
                :class="[
                  'p-2 rounded-xl border transition-all',
                  qty >= tier.min_qty && (!tier.max_qty || qty <= tier.max_qty)
                    ? 'bg-amber-500/10 border-amber-500/40 text-amber-300'
                    : 'bg-slate-900/60 border-slate-800 text-slate-400'
                ]"
              >
                <div class="text-[11px] text-slate-400">
                  {{ tier.max_qty && tier.max_qty < 999 ? `${tier.min_qty}-${tier.max_qty}` : `${tier.min_qty}+` }} pcs
                </div>
                <div class="font-black text-sm text-white mt-1">{{ formatPrice(tier.price) }}</div>
              </div>
            </div>
          </div>
          <div v-else class="p-4 rounded-2xl bg-slate-950 border border-slate-800 flex justify-between items-baseline">
            <span class="text-xs text-slate-400">Unit Price:</span>
            <span class="text-2xl font-black text-white">{{ formatPrice(product.price) }}</span>
          </div>

          <!-- Specifications Table -->
          <div class="space-y-2 text-xs">
            <div class="flex justify-between py-1.5 border-b border-slate-800/80">
              <span class="text-slate-400">Format & Condition:</span>
              <span class="text-slate-200 font-semibold uppercase">{{ product.format }} · {{ product.condition }}</span>
            </div>
            <div class="flex justify-between py-1.5 border-b border-slate-800/80">
              <span class="text-slate-400">Minimum Order (MOQ):</span>
              <span class="text-slate-200 font-semibold">{{ product.moq }} unit(s)</span>
            </div>
            <div v-if="product.type !== 'physical'" class="flex justify-between py-1.5 border-b border-slate-800/80">
              <span class="text-slate-400">Fulfillment Token:</span>
              <span class="text-cyan-400 font-mono font-semibold">{{ product.digital_file_size || 'Instant Download' }} ({{ product.digital_license_format }})</span>
            </div>
            <div v-else class="flex justify-between py-1.5 border-b border-slate-800/80">
              <span class="text-slate-400">Shipping Lead Time:</span>
              <span class="text-amber-400 font-semibold">{{ product.lead_time_days }} Business Days ({{ product.incoterm }})</span>
            </div>
          </div>

          <!-- Order Configurator -->
          <form @submit.prevent="executeCheckout" class="space-y-4 pt-4 border-t border-slate-800">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">
                Order Quantity
              </label>
              <input
                v-model.number="qty"
                type="number"
                :min="product.moq"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2.5 text-sm text-white focus:border-amber-500 outline-none font-bold"
              />
            </div>

            <div class="p-3.5 rounded-xl bg-slate-950 border border-slate-800 flex justify-between items-center text-sm">
              <span class="text-slate-400 text-xs">Total Estimated:</span>
              <span class="text-2xl font-black text-amber-400">{{ formatPrice(computedTotal) }}</span>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <button
                type="submit"
                class="w-full py-3.5 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-black text-xs shadow-xl transition-all"
              >
                {{ product.type !== 'physical' ? '⚡ Instant Procure' : '📦 Place Escrow Order' }}
              </button>

              <button
                v-if="product.accepts_offers"
                type="button"
                @click="showOffer = true"
                class="w-full py-3.5 rounded-xl bg-gradient-to-r from-purple-500 to-indigo-600 hover:from-purple-400 hover:to-indigo-500 text-white font-bold text-xs shadow-xl transition-all"
              >
                🤝 Make an Offer
              </button>
              <Link
                v-else
                href="/#rfq"
                class="w-full py-3.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-bold text-xs text-center transition-all flex items-center justify-center"
              >
                📋 Request Custom Sample
              </Link>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Related Products -->
    <div v-if="related && related.length > 0" class="pt-10 border-t border-slate-800">
      <h3 class="text-xl font-bold text-white mb-6">Frequently Procured Together</h3>
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
            <div class="text-[11px] text-amber-400 font-mono mt-1">{{ formatPrice(rel.price) }}</div>
          </div>
        </Link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Link, router } from '@inertiajs/vue3'
import { ref, computed } from 'vue'
import { useCurrency, type Currency } from '@/composables/useCurrency'

const props = defineProps<{
  product: any
  related: any[]
  currencies?: Currency[]
}>()

const { activeCurrency, formatPrice } = useCurrency(props.currencies)

const qty = ref(props.product.moq || 1)
const showOffer = ref(false)

const computedUnitPrice = computed(() => {
  const p = props.product
  if (p.price_tiers && p.price_tiers.length > 0) {
    for (const tier of p.price_tiers) {
      if (qty.value >= tier.min_qty && (!tier.max_qty || qty.value <= tier.max_qty)) {
        return tier.price
      }
    }
  }
  if (qty.value >= 10 && p.wholesale_price > 0) {
    return p.wholesale_price
  }
  return p.price
})

const computedTotal = computed(() => {
  return computedUnitPrice.value * (qty.value || 1)
})

const executeCheckout = () => {
  router.post('/checkout', {
    product_id: props.product.id,
    quantity: qty.value,
    currency: activeCurrency.value.code,
  })
}
</script>
