<template>
  <div class="min-h-screen flex flex-col bg-slate-950 text-slate-100 selection:bg-amber-500 selection:text-slate-950 font-sans">
    <!-- Top Global Regional & Currency Utility Bar -->
    <header class="border-b border-slate-800 bg-slate-900/90 backdrop-blur-md sticky top-0 z-50">
      <div class="bg-slate-950 border-b border-slate-800/80 text-xs text-slate-400 py-1.5 px-4 sm:px-6">
        <div class="max-w-7xl mx-auto flex items-center justify-between">
          <!-- Left: Trade Assurance & Hybrid Marketplace Specs -->
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center gap-1.5 text-emerald-400 font-semibold">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-ping"></span>
              Soko Escrow: $14.2M Protected
            </span>
            <span class="hidden lg:inline text-slate-700">|</span>
            <span class="hidden lg:inline text-slate-400">
              Alibaba Wholesale + Amazon Fulfillment + eBay Deals
            </span>
          </div>

          <!-- Right: Regional Hub, Currency Switcher & Portal Links -->
          <div class="flex items-center gap-3 sm:gap-4">
            <!-- Region Selector -->
            <div class="relative flex items-center gap-1 text-slate-300">
              <span class="text-xs">Hub:</span>
              <select
                v-model="selectedRegion"
                @change="onRegionChange"
                class="bg-slate-900 text-xs text-amber-400 font-semibold px-2 py-0.5 rounded border border-slate-700 outline-none cursor-pointer"
              >
                <option value="global">🌐 Global All</option>
                <option value="north_america">🇺🇸 North America</option>
                <option value="europe">🇪🇺 Europe & UK</option>
                <option value="east_africa">🇰🇪 East Africa Hub</option>
                <option value="west_africa">🇳🇬 West Africa Hub</option>
                <option value="middle_east">🇦🇪 Middle East</option>
                <option value="asia_pacific">🇨🇳 Asia-Pacific</option>
              </select>
            </div>

            <!-- Currency Switcher -->
            <div class="flex items-center gap-1 text-slate-300">
              <select
                v-model="selectedCurrency"
                class="bg-slate-900 text-xs text-emerald-400 font-mono font-bold px-2 py-0.5 rounded border border-slate-700 outline-none cursor-pointer"
              >
                <option v-for="c in currencies" :key="c.code" :value="c.code">
                  {{ c.code }} ({{ c.symbol }})
                </option>
              </select>
            </div>

            <span class="text-slate-700">|</span>
            <Link v-if="!user" href="/login" class="hover:text-amber-400 transition-colors font-medium">
              Merchant Login
            </Link>
            <Link v-else href="/dashboard" class="text-amber-400 font-semibold hover:underline">
              {{ user.company_name || user.name }}
            </Link>
          </div>
        </div>
      </div>

      <!-- Main Navigation Header -->
      <div class="max-w-7xl mx-auto px-4 sm:px-6 py-3.5 flex items-center justify-between gap-4">
        <!-- Logo -->
        <Link href="/" class="shrink-0 hover:opacity-95 transition-opacity">
          <SokoLogo />
        </Link>

        <!-- Enterprise Search Bar with Multi-Format Filter -->
        <div class="hidden md:flex flex-1 max-w-2xl mx-4">
          <div class="relative w-full flex rounded-xl border border-slate-700 bg-slate-900 overflow-hidden shadow-inner focus-within:border-amber-500 focus-within:ring-2 focus-within:ring-amber-500/20 transition-all">
            <select
              v-model="searchFormat"
              class="bg-slate-800 text-xs font-semibold text-slate-300 px-3 py-2.5 border-r border-slate-700 outline-none cursor-pointer"
            >
              <option value="all">All Formats</option>
              <option value="wholesale">🚢 B2B Wholesale (Alibaba)</option>
              <option value="retail">🛒 Fast Buy (Amazon)</option>
              <option value="digital">⚡ Instant Digital (Licenses)</option>
              <option value="make_offer">🤝 Make Offer (eBay)</option>
            </select>
            <input
              v-model="searchQuery"
              @keydown.enter="applySearch"
              type="text"
              placeholder="Search 10,000+ products, factory machines, cloud software, refurbished datacenter gears..."
              class="w-full bg-transparent px-4 py-2 text-sm text-slate-200 placeholder-slate-500 outline-none"
            />
            <button
              @click="applySearch"
              class="bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold px-5 flex items-center justify-center transition-colors"
              title="Search Soko"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Right Actions -->
        <div class="flex items-center gap-3">
          <!-- B2B RFQ Quick Button -->
          <Link
            href="/#rfq"
            class="hidden sm:inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-bold bg-slate-800 hover:bg-slate-700 border border-slate-700 text-slate-200 transition-colors"
          >
            <span class="text-amber-400">📋</span>
            <span>Post RFQ</span>
          </Link>

          <!-- User Session State -->
          <div v-if="user" class="flex items-center gap-2">
            <Link
              href="/dashboard"
              class="inline-flex items-center gap-2 px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 border border-slate-700 transition-all text-xs"
            >
              <img
                :src="user.avatar || 'https://ui-avatars.com/api/?name=' + user.name"
                class="w-6 h-6 rounded-full object-cover border border-amber-500/50"
                alt="Avatar"
              />
              <div class="hidden xl:flex flex-col text-left">
                <span class="font-bold text-slate-200 leading-none truncate max-w-[120px]">{{ user.name }}</span>
                <span class="text-[10px] text-amber-400 uppercase tracking-wider font-semibold">{{ user.role }}</span>
              </div>
            </Link>
            <Link
              href="/logout"
              method="post"
              as="button"
              class="p-2 text-slate-400 hover:text-rose-400 hover:bg-slate-800 rounded-lg transition-colors"
              title="Sign Out"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
            </Link>
          </div>

          <div v-else class="flex items-center gap-2">
            <Link
              href="/login"
              class="px-3.5 py-2 rounded-xl text-xs font-semibold text-slate-300 hover:text-white hover:bg-slate-800 transition-colors"
            >
              Sign In
            </Link>
            <Link
              href="/register"
              class="px-4 py-2 rounded-xl text-xs font-bold bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 shadow-md shadow-amber-500/20 transition-all"
            >
              Open Account
            </Link>
          </div>
        </div>
      </div>

      <!-- Mega Category & Format Sub-Navigation -->
      <div class="border-t border-slate-800 bg-slate-900/60 px-4 sm:px-6">
        <div class="max-w-7xl mx-auto flex items-center gap-1 sm:gap-2 overflow-x-auto py-2 text-xs scrollbar-none">
          <Link
            href="/"
            :class="[
              'px-3 py-1.5 rounded-lg whitespace-nowrap font-bold transition-colors',
              page.url === '/' ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800'
            ]"
          >
            🏬 Global Soko Center
          </Link>
          <Link
            href="/?format=wholesale"
            :class="[
              'px-3 py-1.5 rounded-lg whitespace-nowrap font-medium transition-colors',
              page.url.includes('format=wholesale') ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800'
            ]"
          >
            🚢 B2B Wholesale Tiers
          </Link>
          <Link
            href="/?format=digital"
            :class="[
              'px-3 py-1.5 rounded-lg whitespace-nowrap font-medium transition-colors',
              page.url.includes('format=digital') ? 'bg-cyan-500/10 text-cyan-400 border border-cyan-500/20' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800'
            ]"
          >
            ⚡ Instant Digital & Software
          </Link>
          <Link
            href="/?format=make_offer"
            :class="[
              'px-3 py-1.5 rounded-lg whitespace-nowrap font-medium transition-colors',
              page.url.includes('format=make_offer') ? 'bg-purple-500/10 text-purple-400 border border-purple-500/20' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800'
            ]"
          >
            🤝 Make an Offer / Deals
          </Link>
          <Link
            href="/?category=Agriculture+%26+Processing"
            :class="[
              'px-3 py-1.5 rounded-lg whitespace-nowrap font-medium transition-colors',
              page.url.includes('Agriculture') ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800'
            ]"
          >
            🌱 Regional Agri-Tech & Milling
          </Link>
          <Link
            href="/#rfq"
            class="px-3 py-1.5 rounded-lg whitespace-nowrap font-medium text-slate-400 hover:text-amber-400 hover:bg-slate-800"
          >
            📋 B2B RFQ Marketplace
          </Link>
          <Link
            href="/dashboard"
            class="px-3 py-1.5 rounded-lg whitespace-nowrap font-bold text-amber-400 hover:text-amber-300 ml-auto bg-amber-500/10 border border-amber-500/20"
          >
            📊 Back-Store Dashboard
          </Link>
        </div>
      </div>
    </header>

    <!-- Flash Alert Banners -->
    <div v-if="flash.success || flash.error" class="max-w-7xl mx-auto w-full px-4 sm:px-6 pt-4">
      <div
        v-if="flash.success"
        class="flex items-center justify-between p-3.5 rounded-xl bg-emerald-950/80 border border-emerald-500/30 text-emerald-300 text-sm shadow-lg shadow-emerald-950/40"
      >
        <div class="flex items-center gap-2.5">
          <svg class="w-5 h-5 text-emerald-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="font-medium">{{ flash.success }}</span>
        </div>
        <button @click="flash.success = ''" class="text-emerald-400 hover:text-white">&times;</button>
      </div>

      <div
        v-if="flash.error"
        class="flex items-center justify-between p-3.5 rounded-xl bg-rose-950/80 border border-rose-500/30 text-rose-300 text-sm shadow-lg shadow-rose-950/40"
      >
        <div class="flex items-center gap-2.5">
          <svg class="w-5 h-5 text-rose-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="font-medium">{{ flash.error }}</span>
        </div>
        <button @click="flash.error = ''" class="text-rose-400 hover:text-white">&times;</button>
      </div>
    </div>

    <!-- Main Content Slot -->
    <main class="flex-1">
      <slot />
    </main>

    <!-- Enterprise Multi-Regional Footer -->
    <footer class="border-t border-slate-800 bg-slate-900/95 text-slate-400 text-xs mt-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 py-12">
        <div class="grid grid-cols-1 md:grid-cols-5 gap-8 mb-12">
          <!-- Col 1: Brand & Regional Network -->
          <div class="md:col-span-2 space-y-4">
            <SokoLogo />
            <p class="text-slate-400 text-xs leading-relaxed max-w-sm">
              Soko combines Alibaba-scale B2B wholesale trade with Amazon-scale automated digital goods delivery and eBay-style marketplace deals. Configurable for any local jurisdiction and currency.
            </p>
            <div class="flex flex-wrap items-center gap-2 text-[11px] text-slate-400">
              <span class="px-2 py-1 rounded bg-slate-800 border border-slate-700 font-mono text-amber-400">ISO 27001</span>
              <span class="px-2 py-1 rounded bg-slate-800 border border-slate-700 font-mono text-emerald-400">Trade Escrow</span>
              <span class="px-2 py-1 rounded bg-slate-800 border border-slate-700 font-mono text-cyan-400">Multi-Currency Engine</span>
            </div>
          </div>

          <!-- Col 2: Marketplace Formats -->
          <div>
            <h4 class="text-slate-200 font-bold uppercase tracking-wider text-[11px] mb-3">Marketplace Formats</h4>
            <ul class="space-y-2">
              <li><Link href="/?format=wholesale" class="hover:text-amber-400 transition-colors">B2B Wholesale Tiers (Alibaba)</Link></li>
              <li><Link href="/?format=retail" class="hover:text-amber-400 transition-colors">Buy It Now / Fast Ship (Amazon)</Link></li>
              <li><Link href="/?format=digital" class="hover:text-amber-400 transition-colors">Instant Software Licenses</Link></li>
              <li><Link href="/?format=make_offer" class="hover:text-amber-400 transition-colors">Make an Offer (eBay)</Link></li>
              <li><Link href="/#rfq" class="hover:text-amber-400 transition-colors">Request for Quotation (RFQ)</Link></li>
            </ul>
          </div>

          <!-- Col 3: Regional Hubs -->
          <div>
            <h4 class="text-slate-200 font-bold uppercase tracking-wider text-[11px] mb-3">Regional Hubs</h4>
            <ul class="space-y-2">
              <li><button @click="switchRegion('north_america', 'USD')" class="hover:text-amber-400 transition-colors">🇺🇸 North America (US/CA)</button></li>
              <li><button @click="switchRegion('europe', 'EUR')" class="hover:text-amber-400 transition-colors">🇪🇺 Europe & UK</button></li>
              <li><button @click="switchRegion('east_africa', 'KES')" class="hover:text-amber-400 transition-colors">🇰🇪 East Africa Hub (M-Pesa)</button></li>
              <li><button @click="switchRegion('west_africa', 'NGN')" class="hover:text-amber-400 transition-colors">🇳🇬 West Africa Hub</button></li>
              <li><button @click="switchRegion('middle_east', 'AED')" class="hover:text-amber-400 transition-colors">🇦🇪 Middle East (GCC)</button></li>
              <li><button @click="switchRegion('asia_pacific', 'CNY')" class="hover:text-amber-400 transition-colors">🇨🇳 Asia-Pacific</button></li>
            </ul>
          </div>

          <!-- Col 4: Platform Architecture -->
          <div>
            <h4 class="text-slate-200 font-bold uppercase tracking-wider text-[11px] mb-3">Back Store Architecture</h4>
            <ul class="space-y-2">
              <li class="font-mono text-slate-500">Goravel Framework v1.18</li>
              <li class="font-mono text-slate-500">Inertia.js + Vue 3</li>
              <li class="font-mono text-slate-500">Tailwind CSS v4 Engine</li>
              <li><Link href="/dashboard" class="hover:text-amber-400 transition-colors">Merchant Back Store</Link></li>
              <li><Link href="/login" class="hover:text-amber-400 transition-colors">Demo Admin Credentials</Link></li>
            </ul>
          </div>
        </div>

        <div class="border-t border-slate-800 pt-6 flex flex-col sm:flex-row items-center justify-between gap-4 text-slate-500">
          <div>
            &copy; 2026 Soko Enterprise Commerce Platform. Tailorable global architecture.
          </div>
          <div class="flex items-center gap-6">
            <span>Trade Escrow Agreement</span>
            <span>Supplier SLA</span>
            <span>Regional Tax Compliance</span>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { Link, usePage, router } from '@inertiajs/vue3'
import { computed, ref } from 'vue'
import SokoLogo from './components/SokoLogo.vue'
import { useCurrency } from '@/composables/useCurrency'

const page = usePage()
const { selectedCurrency, selectedRegion, currencies, setRegion } = useCurrency()

const user = computed(() => (page.props.auth as any)?.user)
const flash = computed(() => (page.props.flash as any) ?? {})

const searchQuery = ref('')
const searchFormat = ref('all')

const applySearch = () => {
  router.get('/', {
    q: searchQuery.value || undefined,
    format: searchFormat.value !== 'all' ? searchFormat.value : undefined,
    region: selectedRegion.value !== 'global' ? selectedRegion.value : undefined,
  }, {
    preserveState: true,
  })
}

const onRegionChange = () => {
  const map: Record<string, string> = {
    global: 'USD',
    north_america: 'USD',
    europe: 'EUR',
    east_africa: 'KES',
    west_africa: 'NGN',
    middle_east: 'AED',
    asia_pacific: 'CNY',
  }
  const def = map[selectedRegion.value] || 'USD'
  setRegion(selectedRegion.value, def)
  applySearch()
}

const switchRegion = (regId: string, curr: string) => {
  setRegion(regId, curr)
  applySearch()
}
</script>
