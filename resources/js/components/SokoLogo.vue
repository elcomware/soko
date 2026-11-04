<template>
    <span
        class="soko-brand"
        :aria-label="`${brand.display_name} — ${brand.tagline}`"
    >
        <svg
            class="soko-bag-mark"
            viewBox="0 0 48 52"
            fill="none"
            aria-hidden="true"
        >
            <path
                d="M17 15V11a7 7 0 0 1 14 0v4"
                stroke="var(--soko-navy)"
                stroke-width="3.2"
                stroke-linecap="round"
            />
            <path
                d="M9 14h30l3 30a4 4 0 0 1-4 4H10a4 4 0 0 1-4-4l3-30Z"
                fill="var(--soko-amber)"
            />
            <path
                d="M17 14v5m14-5v5"
                stroke="var(--soko-navy)"
                stroke-width="2.5"
                stroke-linecap="round"
            />
            <text
                x="24"
                y="41"
                text-anchor="middle"
                fill="var(--soko-navy)"
                font-size="33"
                font-weight="800"
                font-family="var(--soko-font-display)"
            >
                s
            </text>
        </svg>
        <span class="soko-brand-name"
            >{{ brand.display_name.toLowerCase()
            }}<span class="soko-brand-dot">.</span
            ><small>{{ brand.tagline }}</small></span
        >
    </span>
</template>
<script setup lang="ts">
import { computed } from 'vue'
import { usePage } from '@inertiajs/vue3'
import { normalizeBrand, type BrandSettings } from '@/design/brand'
const page = usePage()
const brand = computed(() =>
    normalizeBrand(page.props.brandSettings as Partial<BrandSettings>)
)
</script>
