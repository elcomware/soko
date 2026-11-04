<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, ref, useId, watch } from 'vue'
const props = defineProps<{
    label: string
    items: { id: string; label: string; disabled?: boolean; danger?: boolean }[]
    disabled?: boolean
}>()
const emit = defineEmits<{ select: [id: string] }>()
const menuId = useId()
const trigger = ref<HTMLButtonElement | null>(null)
const menu = ref<HTMLElement | null>(null)
const open = ref(false)
const position = ref({ top: '0px', left: '0px' })
const available = computed(() => props.items.filter((item) => !item.disabled))
function close(restore = false) {
    if (!open.value) return
    open.value = false
    menu.value?.hidePopover?.()
    document.removeEventListener('pointerdown', outside, true)
    document.removeEventListener('focusin', outside, true)
    window.removeEventListener('resize', viewportChanged)
    window.removeEventListener('scroll', scrolled, true)
    if (restore) trigger.value?.focus()
}
function outside(event: Event) {
    const target = event.target as Node
    if (!menu.value?.contains(target) && !trigger.value?.contains(target))
        close()
}
function viewportChanged() {
    close(true)
}
function scrolled(event: Event) {
    if (event.target === window || !menu.value?.contains(event.target as Node))
        close(true)
}
function buttons() {
    return Array.from(
        menu.value?.querySelectorAll<HTMLButtonElement>(
            'button:not(:disabled)'
        ) || []
    )
}
async function show(last = false) {
    if (props.disabled || !available.value.length) return
    open.value = true
    await nextTick()
    if (!open.value || !menu.value || !trigger.value) return
    // Native top layer avoids clipping without detaching inherited brand tokens.
    menu.value.showPopover?.()
    const rect = trigger.value.getBoundingClientRect()
    const width = menu.value.offsetWidth
    const height = menu.value.offsetHeight
    position.value = {
        left:
            Math.max(
                8,
                Math.min(rect.right - width, window.innerWidth - width - 8)
            ) + 'px',
        top:
            Math.max(
                8,
                rect.bottom + height + 8 <= window.innerHeight
                    ? rect.bottom + 4
                    : rect.top - height - 4
            ) + 'px',
    }
    await nextTick()
    if (!open.value) return
    const options = buttons()
    options[last ? options.length - 1 : 0]?.focus()
    document.addEventListener('pointerdown', outside, true)
    document.addEventListener('focusin', outside, true)
    window.addEventListener('resize', viewportChanged)
    window.addEventListener('scroll', scrolled, true)
}
function toggle() {
    if (open.value) close(true)
    else void show()
}
function select(id: string) {
    if (props.disabled || !available.value.some((item) => item.id === id))
        return
    close(true)
    emit('select', id)
}
function keydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
        event.preventDefault()
        event.stopPropagation()
        close(true)
        return
    }
    if (event.key === 'Tab') {
        close(true)
        return
    }
    const options = buttons()
    if (!options.length) return
    const index = options.indexOf(document.activeElement as HTMLButtonElement)
    let target = -1
    if (event.key === 'ArrowDown') target = (index + 1) % options.length
    if (event.key === 'ArrowUp')
        target = (index - 1 + options.length) % options.length
    if (event.key === 'Home') target = 0
    if (event.key === 'End') target = options.length - 1
    if (
        event.key.length === 1 &&
        !event.ctrlKey &&
        !event.altKey &&
        !event.metaKey
    ) {
        for (let offset = 1; offset <= options.length; offset++) {
            const next = (index + offset) % options.length
            if (
                options[next].textContent
                    ?.trim()
                    .toLowerCase()
                    .startsWith(event.key.toLowerCase())
            ) {
                target = next
                break
            }
        }
    }
    if (target >= 0) {
        event.preventDefault()
        options[target]?.focus()
    }
}
watch(
    () => [props.disabled, props.items],
    () => {
        if (open.value) close(true)
    },
    { deep: true }
)
onBeforeUnmount(() => close())
</script>
<template>
    <span v-if="items.length" class="admin-action-menu">
        <button
            ref="trigger"
            type="button"
            class="admin-action-trigger"
            :aria-label="label"
            aria-haspopup="menu"
            :aria-expanded="open"
            :aria-controls="menuId"
            :disabled="disabled || !available.length"
            @click="toggle"
            @keydown.down.prevent="show()"
            @keydown.up.prevent="show(true)"
            @keydown.esc.prevent="close(true)"
        >
            <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="currentColor"
                aria-hidden="true"
            >
                <circle cx="5" cy="12" r="1.8" />
                <circle cx="12" cy="12" r="1.8" />
                <circle cx="19" cy="12" r="1.8" />
            </svg>
        </button>
        <div
            v-show="open"
            :id="menuId"
            ref="menu"
            popover="manual"
            role="menu"
            :aria-label="label"
            class="admin-action-popup"
            :style="position"
            @keydown="keydown"
        >
            <button
                v-for="item in items"
                :key="item.id"
                type="button"
                role="menuitem"
                tabindex="-1"
                :disabled="item.disabled"
                :class="{ 'is-danger': item.danger }"
                @click="select(item.id)"
            >
                {{ item.label }}
            </button>
        </div>
    </span>
</template>
<style scoped>
.admin-action-menu {
    display: inline-flex;
}
.admin-action-trigger {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 44px;
    height: 44px;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    background: var(--soko-surface);
    color: var(--soko-ink);
    cursor: pointer;
}
.admin-action-trigger:hover {
    background: var(--soko-canvas);
}
.admin-action-trigger:disabled,
.admin-action-popup button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}
.admin-action-trigger:focus-visible,
.admin-action-popup button:focus-visible {
    outline: 2px solid var(--soko-navy);
    outline-offset: 2px;
}
.admin-action-popup {
    position: fixed;
    inset: auto;
    margin: 0;
    z-index: 1000;
    min-width: 180px;
    max-width: calc(100vw - 16px);
    max-height: calc(100vh - 16px);
    overflow-y: auto;
    padding: 6px;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    background: var(--soko-surface);
    color: var(--soko-ink);
    font-family: var(--soko-font);
    box-shadow: 0 8px 24px color-mix(in srgb, var(--soko-ink) 16%, transparent);
}
.admin-action-popup button {
    display: block;
    width: 100%;
    border: 0;
    border-radius: var(--soko-radius);
    background: transparent;
    color: inherit;
    font: inherit;
    font-size: 13px;
    text-align: left;
    padding: 10px 12px;
    cursor: pointer;
}
.admin-action-popup button:hover,
.admin-action-popup button:focus {
    background: var(--soko-canvas);
}
.admin-action-popup button.is-danger {
    color: var(--soko-danger, #a32121);
}
</style>
