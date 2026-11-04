<template>
    <svg
        viewBox="0 0 96 96"
        :width="size"
        :height="size"
        class="soko-commerce-illustration"
        :role="label ? 'img' : undefined"
        :aria-label="label"
        :aria-hidden="label ? undefined : true"
        fill="none"
    >
        <defs>
            <linearGradient :id="`${id}-paper`" x1="0" y1="0" x2="1" y2="1">
                <stop stop-color="var(--soko-surface, #fff)" />
                <stop offset="1" stop-color="var(--soko-canvas, #f5f7fa)" />
            </linearGradient>
            <linearGradient :id="`${id}-accent`" x1="0" y1="0" x2="1" y2="1">
                <stop stop-color="var(--soko-amber, #f3ac32)" />
                <stop
                    offset="1"
                    stop-color="var(--soko-icon-deep-accent, #c88b26)"
                />
            </linearGradient>
            <linearGradient :id="`${id}-glass`" x1="0" y1="0" x2="1" y2="1">
                <stop stop-color="var(--soko-surface, #fff)" />
                <stop
                    offset="1"
                    stop-color="var(--soko-navy, #142f48)"
                    stop-opacity=".25"
                />
            </linearGradient>
        </defs>
        <ellipse
            cx="48"
            cy="87"
            rx="36"
            ry="5"
            fill="var(--soko-navy)"
            opacity=".09"
        />
        <path
            v-for="(layer, index) in layers"
            :key="index"
            :d="layer.d"
            :fill="layer.detail ? 'none' : tone(layer.tone)"
            stroke="var(--soko-navy, #142f48)"
            :stroke-opacity="layer.detail ? 0.85 : 0.55"
            :stroke-width="layer.detail ? 1.5 : 1.1"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
    </svg>
</template>
<script setup lang="ts">
import { computed, useId } from 'vue'
import {
    commerceIllustrations,
    type IllustrationTone,
} from '@/design/commerceIllustrations'
const props = withDefaults(
    defineProps<{ name: string; size?: number; label?: string }>(),
    { size: 96 }
)
const id = `soko-scene-${useId().replace(/:/g, '')}`
const layers = computed(
    () => commerceIllustrations[props.name] || commerceIllustrations.storefront
)
const tone = (name: IllustrationTone) =>
    name === 'ink'
        ? 'var(--soko-navy)'
        : name === 'shade'
          ? 'color-mix(in srgb, var(--soko-navy) 22%, var(--soko-surface))'
          : `url(#${id}-${name})`
</script>
