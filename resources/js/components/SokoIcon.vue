<template>
    <svg
        :viewBox="dimensional ? '-1 -1 26 26' : '0 0 24 24'"
        fill="none"
        stroke="currentColor"
        :stroke-width="strokeWidth ?? 'var(--soko-icon-stroke, 1.6)'"
        stroke-linecap="round"
        stroke-linejoin="round"
        :aria-hidden="label ? undefined : true"
        :role="label ? 'img' : undefined"
        :aria-label="label || undefined"
        class="soko-icon"
        :style="size ? { width: `${size}px`, height: `${size}px` } : undefined"
        :data-icon="resolvedName"
    >
        <defs v-if="dimensional">
            <linearGradient :id="gradientId" x1="0" y1="0" x2="1" y2="1">
                <stop
                    offset="0"
                    stop-color="var(--soko-icon-light, var(--soko-surface, #fff))"
                />
                <stop
                    offset="0.48"
                    stop-color="var(--soko-icon-accent, var(--soko-amber, #f3ac32))"
                />
                <stop
                    offset="1"
                    stop-color="var(--soko-icon-shade, var(--soko-navy, #142f48))"
                    stop-opacity="0.38"
                />
            </linearGradient>
        </defs>
        <g v-if="dimensional">
            <path
                :d="iconPaths[resolvedName]"
                transform="translate(0 .65)"
                stroke="var(--soko-icon-shade, var(--soko-navy))"
                stroke-opacity=".14"
                :stroke-width="2.2"
            />
            <path
                v-if="surface"
                :d="surface"
                :fill="`url(#${gradientId})`"
                fill-opacity=".85"
                stroke="none"
            />
            <path
                :d="iconPaths[resolvedName]"
                transform="translate(0 -.22)"
                stroke="var(--soko-icon-light, var(--soko-surface))"
                stroke-opacity=".8"
                :stroke-width="2"
            />
        </g>
        <path
            :d="iconPaths[resolvedName]"
            :stroke-width="dimensional ? 1.25 : undefined"
        />
    </svg>
</template>
<script setup lang="ts">
import { computed, useId } from 'vue'
import { iconPaths } from '@/design/icons'
import { iconSurface } from '@/design/iconAppearance'
const props = defineProps<{
    name: string
    label?: string
    size?: number
    strokeWidth?: number
    variant?: 'line' | 'dimensional' | 'auto'
}>()
const resolvedName = computed(() =>
    Object.prototype.hasOwnProperty.call(iconPaths, props.name)
        ? props.name
        : 'help'
)

const gradientId = `soko-icon-${useId().replace(/:/g, '')}`
const dimensional = computed(
    () =>
        props.variant === 'dimensional' ||
        (props.variant !== 'line' && (props.size || 24) >= 28)
)
const surface = computed(() => iconSurface(resolvedName.value))
</script>
