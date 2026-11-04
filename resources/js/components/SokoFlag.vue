<template>
    <svg
        viewBox="0 0 30 20"
        class="soko-flag"
        role="img"
        :aria-label="label || names[code.toUpperCase()] || code.toUpperCase()"
        :style="{ width: `${size}px` }"
    >
        <template v-if="flag">
            <rect width="30" height="20" :fill="flag.base" />
            <rect
                v-for="(stripe, index) in flag.stripes || []"
                :key="index"
                :x="stripe.x || 0"
                :y="stripe.y || 0"
                :width="stripe.w || 30"
                :height="stripe.h || 20"
                :fill="stripe.color"
            />
            <circle
                v-if="flag.circle"
                :cx="flag.circle.x"
                :cy="flag.circle.y"
                :r="flag.circle.r"
                :fill="flag.circle.color"
            />
            <path
                v-if="flag.star"
                d="m15 5 1.4 3 3.3.4-2.4 2.3.6 3.3-2.9-1.6-2.9 1.6.6-3.3L10.3 8.4l3.3-.4L15 5Z"
                :fill="flag.star"
            />
        </template>
        <template v-else>
            <rect width="30" height="20" fill="var(--soko-canvas)" />
            <text
                x="15"
                y="13"
                text-anchor="middle"
                font-size="8"
                fill="currentColor"
            >
                {{ code.toUpperCase().slice(0, 3) }}
            </text>
        </template>
        <rect
            x=".4"
            y=".4"
            width="29.2"
            height="19.2"
            fill="none"
            stroke="currentColor"
            stroke-opacity=".15"
            stroke-width=".8"
        />
    </svg>
</template>
<script setup lang="ts">
import { computed } from 'vue'
type Stripe = { x?: number; y?: number; w?: number; h?: number; color: string }
type Flag = {
    base: string
    stripes?: Stripe[]
    circle?: { x: number; y: number; r: number; color: string }
    star?: string
}
const props = withDefaults(
    defineProps<{ code: string; label?: string; size?: number }>(),
    { size: 24 }
)
const triV = (a: string, b: string, c: string): Flag => ({
    base: b,
    stripes: [
        { w: 10, color: a },
        { x: 20, w: 10, color: c },
    ],
})
const triH = (a: string, b: string, c: string): Flag => ({
    base: b,
    stripes: [
        { h: 20 / 3, color: a },
        { y: 40 / 3, h: 20 / 3, color: c },
    ],
})
const flags: Record<string, Flag> = {
    FR: triV('#002395', '#fff', '#ed2939'),
    IT: triV('#009246', '#fff', '#ce2b37'),
    IE: triV('#169b62', '#fff', '#ff883e'),
    NG: triV('#008751', '#fff', '#008751'),
    CI: triV('#f77f00', '#fff', '#009e60'),
    BE: triV('#000', '#fdda24', '#ef3340'),
    DE: triH('#000', '#d00', '#ffce00'),
    NL: triH('#ae1c28', '#fff', '#21468b'),
    AT: triH('#ed2939', '#fff', '#ed2939'),
    RU: triH('#fff', '#0039a6', '#d52b1e'),
    PL: { base: '#fff', stripes: [{ y: 10, h: 10, color: '#dc143c' }] },
    ID: { base: '#fff', stripes: [{ h: 10, color: '#ff0000' }] },
    JP: { base: '#fff', circle: { x: 15, y: 10, r: 6, color: '#bc002d' } },
    BD: { base: '#006a4e', circle: { x: 13.5, y: 10, r: 6, color: '#f42a41' } },
    GH: { ...triH('#ce1126', '#fcd116', '#006b3f'), star: '#000' },
    CM: { ...triV('#007a5e', '#ce1126', '#fcd116'), star: '#fcd116' },
    SN: { ...triV('#00853f', '#fdef42', '#e31b23'), star: '#00853f' },
    SE: {
        base: '#006aa7',
        stripes: [
            { x: 9, w: 4, color: '#fecc00' },
            { y: 8, h: 4, color: '#fecc00' },
        ],
    },
    FI: {
        base: '#fff',
        stripes: [
            { x: 8, w: 5, color: '#003580' },
            { y: 7, h: 6, color: '#003580' },
        ],
    },
    DK: {
        base: '#c8102e',
        stripes: [
            { x: 9, w: 3, color: '#fff' },
            { y: 8.5, h: 3, color: '#fff' },
        ],
    },
}
const names: Record<string, string> = {
    FR: 'France',
    IT: 'Italy',
    IE: 'Ireland',
    NG: 'Nigeria',
    CI: 'Côte d’Ivoire',
    BE: 'Belgium',
    DE: 'Germany',
    NL: 'Netherlands',
    AT: 'Austria',
    RU: 'Russia',
    PL: 'Poland',
    ID: 'Indonesia',
    JP: 'Japan',
    BD: 'Bangladesh',
    GH: 'Ghana',
    CM: 'Cameroon',
    SN: 'Senegal',
    SE: 'Sweden',
    FI: 'Finland',
    DK: 'Denmark',
}
const flag = computed(() => flags[props.code.toUpperCase()])
</script>
