<template>
    <svg
        :viewBox="variant === 'dimensional' ? '0 0 32 32' : '0 0 32 24'"
        class="soko-currency-symbol"
        role="img"
        :aria-label="label || code.toUpperCase()"
    >
        <template v-if="variant === 'dimensional'">
            <defs>
                <linearGradient :id="id" x1="0" y1="0" x2="1" y2="1">
                    <stop stop-color="var(--soko-surface)" />
                    <stop offset=".45" stop-color="var(--soko-amber)" />
                    <stop
                        offset="1"
                        stop-color="var(--soko-icon-deep-accent, var(--soko-amber))"
                    />
                </linearGradient>
            </defs>
            <circle
                cx="16"
                cy="17"
                r="14"
                fill="var(--soko-navy)"
                opacity=".15"
            />
            <circle
                cx="16"
                cy="15"
                r="14"
                :fill="`url(#${id})`"
                stroke="var(--soko-navy)"
                stroke-width=".7"
            />
            <circle
                cx="16"
                cy="15"
                r="11.3"
                fill="none"
                stroke="var(--soko-navy)"
                stroke-opacity=".3"
                stroke-width=".5"
            />
        </template>
        <text
            x="16"
            :y="variant === 'dimensional' ? 21 : 18"
            text-anchor="middle"
            :font-size="
                symbol.length > 2 ? 10 : variant === 'dimensional' ? 17 : 20
            "
            fill="currentColor"
            font-family="var(--soko-font)"
            font-weight="650"
        >
            {{ symbol }}
        </text>
    </svg>
</template>
<script setup lang="ts">
import { computed, useId } from 'vue'
const props = defineProps<{
    code: string
    label?: string
    variant?: 'line' | 'dimensional'
}>()
const id = `soko-currency-${useId().replace(/:/g, '')}`
const symbols: Record<string, string> = {
    USD: '$',
    EUR: '€',
    GBP: '£',
    JPY: '¥',
    CNY: '¥',
    INR: '₹',
    NGN: '₦',
    KES: 'KSh',
    ZAR: 'R',
    GHS: '₵',
    KRW: '₩',
    THB: '฿',
    VND: '₫',
    RUB: '₽',
    TRY: '₺',
    PHP: '₱',
    XOF: 'CFA',
    XAF: 'FCFA',
    AED: 'AED',
    BRL: 'R$',
    CAD: 'C$',
    AUD: 'A$',
    CHF: 'Fr',
}
const symbol = computed(
    () => symbols[props.code.toUpperCase()] || props.code.toUpperCase()
)
</script>
