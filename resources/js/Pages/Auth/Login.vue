<template>
  <div class="max-w-md mx-auto px-4 py-16 space-y-8">
    <div class="text-center space-y-2">
      <div class="inline-flex p-3 rounded-2xl bg-amber-500/10 border border-amber-500/20 text-amber-400 mb-2">
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
        </svg>
      </div>
      <h1 class="text-2xl font-black text-white">Sign In to Soko</h1>
      <p class="text-xs text-slate-400">Access your enterprise procurement, merchant console, or admin controls</p>
    </div>

    <!-- Login Box -->
    <div class="bg-slate-900 border border-slate-800 rounded-3xl p-6 sm:p-8 shadow-2xl space-y-6">
      <form @submit.prevent="submit" class="space-y-4">
        <div>
          <label class="block text-xs font-semibold text-slate-300 mb-1">Corporate Email Address</label>
          <input
            v-model="form.email"
            type="email"
            required
            placeholder="name@company.com"
            class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2.5 text-sm text-white focus:border-amber-500 focus:ring-1 focus:ring-amber-500 outline-none"
          />
        </div>

        <div>
          <div class="flex items-center justify-between mb-1">
            <label class="block text-xs font-semibold text-slate-300">Password</label>
            <span class="text-[11px] text-amber-400 hover:underline cursor-pointer">Forgot?</span>
          </div>
          <input
            v-model="form.password"
            type="password"
            required
            placeholder="••••••••"
            class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2.5 text-sm text-white focus:border-amber-500 focus:ring-1 focus:ring-amber-500 outline-none"
          />
        </div>

        <button
          type="submit"
          :disabled="form.processing"
          class="w-full py-3 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-bold text-sm shadow-lg shadow-amber-500/20 transition-all flex items-center justify-center gap-2 disabled:opacity-50"
        >
          <span>Authenticate Session</span>
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M14 5l7 7m0 0l-7 7m7-7H3" />
          </svg>
        </button>
      </form>

      <!-- Demo Accounts 1-Click Fill -->
      <div class="pt-6 border-t border-slate-800 space-y-3">
        <div class="text-[11px] font-bold uppercase tracking-wider text-slate-400 text-center">
          ⚡ One-Click Demo Profiles
        </div>

        <div class="space-y-2">
          <button
            v-for="demo in demoUsers"
            :key="demo.email"
            @click="quickLogin(demo)"
            type="button"
            class="w-full p-2.5 rounded-xl bg-slate-950/70 hover:bg-slate-800/80 border border-slate-800 text-left transition-colors flex items-center justify-between group"
          >
            <div>
              <div class="font-bold text-xs text-white group-hover:text-amber-400 transition-colors">
                {{ demo.role }}
              </div>
              <div class="text-[10px] text-slate-400 font-mono">{{ demo.email }}</div>
            </div>
            <span class="text-[10px] px-2 py-0.5 rounded bg-slate-800 border border-slate-700 text-amber-400 font-semibold">
              {{ demo.badge }}
            </span>
          </button>
        </div>
      </div>
    </div>

    <!-- Register Link -->
    <div class="text-center text-xs text-slate-400">
      Need an enterprise or supplier account?
      <Link href="/register" class="text-amber-400 font-semibold hover:underline">
        Create one in 2 minutes
      </Link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Link, useForm } from '@inertiajs/vue3'

const props = defineProps<{
  demoUsers: Array<{
    role: string
    email: string
    password: string
    company: string
    badge: string
  }>
}>()

const form = useForm({
  email: '',
  password: '',
})

const submit = () => {
  form.post('/login')
}

const quickLogin = (demo: { email: string; password: string }) => {
  form.email = demo.email
  form.password = demo.password
  form.post('/login')
}
</script>

