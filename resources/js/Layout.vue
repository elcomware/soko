<template>
  <div class="min-h-screen flex flex-col bg-slate-950 text-slate-100 selection:bg-amber-500 selection:text-slate-950 font-sans">
    <!-- Top Global Notification & Utility Bar -->
    <header class="border-b border-slate-800 bg-slate-900/90 backdrop-blur-md sticky top-0 z-50">
      <!-- Topmost micro-bar -->
      <div class="bg-slate-950/80 border-b border-slate-800/60 text-xs text-slate-400 py-1.5 px-4 sm:px-6">
        <div class="max-w-7xl mx-auto flex items-center justify-between">
          <div class="flex items-center gap-4">
            <span class="inline-flex items-center gap-1.5 text-emerald-400 font-medium">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-ping"></span>
              Soko Trade Assurance & Escrow Active
            </span>
            <span class="hidden md:inline text-slate-600">|</span>
            <span class="hidden md:inline hover:text-slate-300 transition-colors">
              ⚡ Instant Zero-Latency Digital Delivery & Global Multi-modal Freight
            </span>
          </div>
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2 cursor-pointer hover:text-amber-400 transition-colors">
              <span>🌐 Global / English</span>
              <span>USD ($)</span>
            </div>
            <span class="text-slate-700">|</span>
            <Link v-if="!user" href="/login" class="hover:text-amber-400 transition-colors font-medium">
              Supplier Portal
            </Link>
            <Link v-else href="/dashboard" class="text-amber-400 font-semibold hover:underline">
              {{ user.company_name || user.name }} (Dashboard)
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

        <!-- Search Bar with Product Type Filter -->
        <div class="hidden lg:flex flex-1 max-w-2xl mx-6">
          <div class="relative w-full flex rounded-xl border border-slate-700 bg-slate-900 overflow-hidden shadow-inner focus-within:border-amber-500 focus-within:ring-2 focus-within:ring-amber-500/20 transition-all">
            <select
              v-model="selectedType"
              @change="applySearch"
              class="bg-slate-800/80 text-xs font-semibold text-slate-300 px-3 py-2.5 border-r border-slate-700 outline-none cursor-pointer hover:bg-slate-800 transition-colors"
            >
              <option value="all">All Goods</option>
              <option value="digital">⚡ Digital & Software</option>
              <option value="physical">🚢 Physical & Freight</option>
            </select>
            <input
              v-model="searchQuery"
              @keydown.enter="applySearch"
              type="text"
              placeholder="Search 10,000+ enterprise products, AI models, cloud software, industrial CNC..."
              class="w-full bg-transparent px-4 py-2.5 text-sm text-slate-200 placeholder-slate-500 outline-none"
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

        <!-- Right Action Items -->
        <div class="flex items-center gap-3">
          <!-- RFQ / Wholesale Button -->
          <Link
            href="/#rfq"
            class="hidden sm:inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-semibold bg-slate-800/80 hover:bg-slate-700/80 border border-slate-700 text-slate-200 transition-colors"
          >
            <svg class="w-3.5 h-3.5 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span>Request Quotation</span>
          </Link>

          <!-- Auth Status -->
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
              <div class="hidden md:flex flex-col text-left">
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
              Create Account
            </Link>
          </div>
        </div>
      </div>

      <!-- Category Links Sub-bar -->
      <div class="border-t border-slate-800/80 bg-slate-900/50 px-4 sm:px-6">
        <div class="max-w-7xl mx-auto flex items-center gap-1 sm:gap-2 overflow-x-auto py-2 text-xs scrollbar-none">
          <Link
            href="/"
            :class="[
              'px-3 py-1.5 rounded-lg whitespace-nowrap font-medium transition-colors',
              page.url === '/' ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'
            ]"
          >
            🏬 All Soko Marketplace
          </Link>
          <Link
            href="/?type=digital"
            :class="[
              'px-3 py-1.5 rounded-lg whitespace-nowrap font-medium transition-colors',
              page.url.includes('type=digital') ? 'bg-cyan-500/10 text-cyan-400 border border-cyan-500/20' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'
            ]"
          >
            ⚡ Digital Software & SaaS
          </Link>
          <Link
            href="/?category=Digital+Assets"
            :class="[
              'px-3 py-1.5 rounded-lg whitespace-nowrap font-medium transition-colors',
              page.url.includes('Digital+Assets') ? 'bg-purple-500/10 text-purple-400 border border-purple-500/20' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'
            ]"
          >
            🧠 AI Models & CAD Blueprints
          </Link>
          <Link
            href="/?type=physical"
            :class="[
              'px-3 py-1.5 rounded-lg whitespace-nowrap font-medium transition-colors',
              page.url.includes('type=physical') ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'
            ]"
          >
            🚢 Industrial Machinery & Freight
          </Link>
          <Link
            href="/?category=Electronics+%26+Telecom"
            :class="[
              'px-3 py-1.5 rounded-lg whitespace-nowrap font-medium transition-colors',
              page.url.includes('Electronics') ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'
            ]"
          >
            🔌 Datacenter & Telecom
          </Link>
          <Link
            href="/dashboard"
            :class="[
              'px-3 py-1.5 rounded-lg whitespace-nowrap font-medium transition-colors ml-auto text-amber-400 hover:text-amber-300 font-bold',
              page.url.startsWith('/dashboard') ? 'bg-amber-500/20 text-amber-300 border border-amber-500/30' : ''
            ]"
          >
            📊 Enterprise Dashboard
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

    <!-- Enterprise Footer -->
    <footer class="border-t border-slate-800 bg-slate-900/95 text-slate-400 text-xs mt-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 py-12">
        <div class="grid grid-cols-1 md:grid-cols-5 gap-8 mb-12">
          <!-- Col 1: Brand Info -->
          <div class="md:col-span-2 space-y-4">
            <SokoLogo />
            <p class="text-slate-400 text-xs leading-relaxed max-w-sm">
              Soko is the next-generation global enterprise commerce and digital marketplace engine.
              Seamlessly integrating Alibaba-scale wholesale physical freight with Amazon-scale automated digital goods delivery.
            </p>
            <div class="flex items-center gap-3 text-slate-400 text-xs">
              <span class="px-2 py-1 rounded bg-slate-800 border border-slate-700 font-mono text-amber-400">ISO 27001 Certified</span>
              <span class="px-2 py-1 rounded bg-slate-800 border border-slate-700 font-mono text-emerald-400">Trade Escrow Protected</span>
            </div>
          </div>

          <!-- Col 2: Sourcing Solutions -->
          <div>
            <h4 class="text-slate-200 font-bold uppercase tracking-wider text-[11px] mb-3">Sourcing Solutions</h4>
            <ul class="space-y-2">
              <li><Link href="/?type=physical" class="hover:text-amber-400 transition-colors">B2B Wholesale & Freight</Link></li>
              <li><Link href="/?type=digital" class="hover:text-amber-400 transition-colors">Software & Cloud Licenses</Link></li>
              <li><Link href="/?category=Digital+Assets" class="hover:text-amber-400 transition-colors">AI Models & Weights</Link></li>
              <li><Link href="/#rfq" class="hover:text-amber-400 transition-colors">Request for Quotation (RFQ)</Link></li>
              <li><Link href="/#assurance" class="hover:text-amber-400 transition-colors">Trade Assurance Guarantee</Link></li>
            </ul>
          </div>

          <!-- Col 3: For Suppliers -->
          <div>
            <h4 class="text-slate-200 font-bold uppercase tracking-wider text-[11px] mb-3">For Suppliers & Developers</h4>
            <ul class="space-y-2">
              <li><Link href="/register" class="hover:text-amber-400 transition-colors">Become a Soko Supplier</Link></li>
              <li><Link href="/dashboard" class="hover:text-amber-400 transition-colors">Merchant Dashboard</Link></li>
              <li><Link href="/dashboard" class="hover:text-amber-400 transition-colors">Instant Digital Delivery Setup</Link></li>
              <li><Link href="/dashboard" class="hover:text-amber-400 transition-colors">Cross-Border Escrow Payouts</Link></li>
            </ul>
          </div>

          <!-- Col 4: Platform & Support -->
          <div>
            <h4 class="text-slate-200 font-bold uppercase tracking-wider text-[11px] mb-3">Platform & Architecture</h4>
            <ul class="space-y-2">
              <li class="font-mono text-slate-500">Goravel Framework v1.18</li>
              <li class="font-mono text-slate-500">Inertia.js + Vue 3</li>
              <li class="font-mono text-slate-500">Tailwind CSS v4 Engine</li>
              <li><Link href="/login" class="hover:text-amber-400 transition-colors">Demo Admin Credentials</Link></li>
            </ul>
          </div>
        </div>

        <div class="border-t border-slate-800 pt-6 flex flex-col sm:flex-row items-center justify-between gap-4 text-slate-500">
          <div>
            &copy; 2026 Soko Enterprise Commerce Platform. All rights reserved.
          </div>
          <div class="flex items-center gap-6">
            <span class="hover:text-slate-400 cursor-pointer">Terms of Service</span>
            <span class="hover:text-slate-400 cursor-pointer">Privacy Policy</span>
            <span class="hover:text-slate-400 cursor-pointer">Trade Escrow Agreement</span>
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

const page = usePage()

const user = computed(() => (page.props.auth as any)?.user)
const flash = computed(() => (page.props.flash as any) ?? {})

const searchQuery = ref('')
const selectedType = ref('all')

const applySearch = () => {
  router.get('/', {
    q: searchQuery.value,
    type: selectedType.value !== 'all' ? selectedType.value : undefined,
  }, {
    preserveState: true,
  })
}
</script>
