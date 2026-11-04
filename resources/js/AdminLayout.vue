<template>
    <div class="soko-app admin-shell" :style="brandTokens(brand)">
        <Head :title="`${currentTitle} · ${brand.display_name}`" />
        <a href="#admin-content" class="skip-link">Skip to workspace content</a>
        <aside class="admin-sidebar" aria-label="Workspace navigation">
            <Link href="/dashboard" class="admin-brand"
                ><SokoLogo /><span class="admin-brand-caption"
                    >Workspace</span
                ></Link
            >
            <button
                id="admin-menu-toggle"
                class="admin-menu-toggle"
                type="button"
                :aria-expanded="menuOpen"
                aria-controls="admin-navigation"
                @click="menuOpen = !menuOpen"
            >
                <SokoIcon name="menu" />{{
                    menuOpen ? 'Close navigation' : 'Open navigation'
                }}
            </button>
            <div
                id="admin-navigation"
                class="admin-navigation"
                :class="{ 'is-open': menuOpen }"
            >
                <nav
                    v-for="group in groups"
                    :key="group.label"
                    :aria-label="group.label"
                >
                    <h2>{{ group.label }}</h2>
                    <Link
                        v-for="item in group.items"
                        :key="item.href"
                        :href="item.href"
                        :class="{ 'is-current': active(item.href) }"
                        :aria-current="active(item.href) ? 'page' : undefined"
                        @click="closeNavigation"
                        ><SokoIcon :name="item.icon" /><span>{{
                            item.label
                        }}</span></Link
                    >
                </nav>
                <div class="admin-sidebar-bottom">
                    <Link href="/"
                        ><SokoIcon name="arrow" />View storefront</Link
                    ><Link href="/account"
                        ><SokoIcon name="user" />My purchases</Link
                    >
                </div>
            </div>
        </aside>
        <div class="admin-stage">
            <header class="admin-topbar">
                <nav aria-label="Breadcrumb">
                    <Link href="/dashboard">Workspace</Link
                    ><span aria-hidden="true">/</span
                    ><span aria-current="page">{{ currentTitle }}</span>
                </nav>
                <div class="admin-actor">
                    <span
                        ><strong>{{
                            user.company_name || user.name || 'Workspace member'
                        }}</strong
                        ><small>{{ user.email }}</small></span
                    ><Link
                        href="/logout"
                        method="post"
                        as="button"
                        class="admin-signout"
                        >Sign out</Link
                    >
                </div>
            </header>
            <main id="admin-content" class="admin-content" tabindex="-1">
                <slot />
            </main>
        </div>
        <SharedFlashNotifications />
    </div>
</template>
<script setup lang="ts">
import { computed, ref } from 'vue'
import { Head, Link, usePage } from '@inertiajs/vue3'
import SokoLogo from '@/components/SokoLogo.vue'
import SokoIcon from '@/components/SokoIcon.vue'
import SharedFlashNotifications from '@/components/SharedFlashNotifications.vue'
import { brandTokens, normalizeBrand } from '@/design/brand'
import '../css/workspace.css'
import '../css/admin.css'
type NavItem = { label: string; href: string; icon: string; group: string }
const page = usePage()
const menuOpen = ref(false)
function closeNavigation() {
    if (!menuOpen.value) return
    menuOpen.value = false
    // Keep focus on a visible control when the mobile menu containing the link closes.
    const toggle = document.getElementById('admin-menu-toggle')
    if (toggle?.getClientRects().length) toggle.focus()
}
const brand = computed(() => normalizeBrand(page.props.brandSettings as any))
const user = computed(() => (page.props.auth as any)?.user || {})
const navigation = computed(
    () => (page.props.workspaceNavigation as NavItem[]) || []
)
const groups = computed(() =>
    [...new Set(navigation.value.map((item) => item.group))].map((label) => ({
        label,
        items: navigation.value.filter((item) => item.group === label),
    }))
)
const currentUrl = computed(
    () => new URL(page.url, 'https://workspace.invalid')
)
function active(href: string) {
    const candidate = new URL(href, 'https://workspace.invalid')
    if (candidate.pathname !== currentUrl.value.pathname) return false
    for (const key of ['section', 'panel']) {
        let value = currentUrl.value.searchParams.get(key) || ''
        if (key === 'section' && value === 'overview') value = ''
        if (
            key === 'panel' &&
            !value &&
            currentUrl.value.pathname === '/dashboard/settings'
        ) {
            const first = navigation.value.find(
                (item) =>
                    new URL(item.href, 'https://workspace.invalid').pathname ===
                    '/dashboard/settings'
            )
            value = first
                ? new URL(
                      first.href,
                      'https://workspace.invalid'
                  ).searchParams.get('panel') || ''
                : ''
        }
        if ((candidate.searchParams.get(key) || '') !== value) return false
    }
    return true
}
const currentTitle = computed(
    () =>
        navigation.value.find((item) => active(item.href))?.label ||
        (page.component === 'Dashboard/Settings' ? 'Settings' : 'Workspace')
)
</script>
