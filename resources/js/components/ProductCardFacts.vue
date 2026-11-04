<template>
    <div
        v-if="facts.length"
        class="product-card-facts"
        :class="{ 'facts-rotating': rotating }"
        aria-label="Product details"
        @mouseenter="hold"
        @mouseleave="release"
        @focusin="hold"
        @focusout="release"
    >
        <span
            v-for="fact in visibleFacts"
            :key="fact.key"
            class="product-card-fact"
            :class="`fact-${fact.key}`"
            :title="fact.description"
        >
            <SokoIcon
                :name="
                    fact.key === 'rating'
                        ? 'star'
                        : fact.key === 'availability'
                          ? 'box'
                          : 'tag'
                "
            />
            <span>{{ fact.label }}</span>
            <span class="sr-only">. {{ fact.description }}</span>
        </span>
    </div>
</template>
<script setup lang="ts">
import { computed, ref } from 'vue'
import SokoIcon from './SokoIcon.vue'
import {
    normalizeCardFooter,
    productCardFacts,
    cardFactIndex,
    type CardFooterSettings,
} from '@/design/cardFacts'
const props = defineProps<{
    product: Record<string, any>
    settings: CardFooterSettings
    formatPrice: (value: number) => string
    now: number
    paused: boolean
}>()
const config = computed(() => normalizeCardFooter(props.settings))
const facts = computed(() =>
    productCardFacts(
        props.product,
        props.settings,
        props.formatPrice,
        props.now
    )
)
const rotating = computed(
    () =>
        config.value.mode === 'rotate' &&
        !props.paused &&
        facts.value.length > 1
)
const rotationStartedAt = props.now
const held = ref<number | null>(null)
const currentIndex = computed(() =>
    cardFactIndex(
        Math.max(0, props.now - rotationStartedAt),
        config.value.interval,
        facts.value.length
    )
)
const visibleFacts = computed(() =>
    rotating.value
        ? [facts.value[(held.value ?? currentIndex.value) % facts.value.length]]
        : facts.value
)
function hold() {
    held.value = currentIndex.value
}
function release() {
    held.value = null
}
</script>
