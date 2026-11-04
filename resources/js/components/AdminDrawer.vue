<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, useId, watch } from 'vue'
import { router } from '@inertiajs/vue3'
import SokoIcon from './SokoIcon.vue'
const props = withDefaults(
    defineProps<{
        open: boolean
        title: string
        description?: string
        busy?: boolean
        dirty?: boolean
    }>(),
    { busy: false, dirty: false }
)
const emit = defineEmits<{ close: [] }>()
const titleId = useId()
const descriptionId = useId()
const dialog = ref<HTMLDialogElement | null>(null)
const heading = ref<HTMLElement | null>(null)
const keepButton = ref<HTMLButtonElement | null>(null)
const confirming = ref(false)
let opener: HTMLElement | null = null
let returnToEditor: HTMLElement | null = null
let previousOverflow = ''
let locked = false
let pendingVisit: any = null
let allowVisit = false
let removeVisitListener: (() => void) | undefined
onMounted(() => {
    removeVisitListener = router.on('before', (event) => {
        if (allowVisit || !props.open || event.detail.visit.method !== 'get')
            return
        if (props.busy) {
            event.preventDefault()
            return
        }
        if (!props.dirty) return
        event.preventDefault()
        pendingVisit = event.detail.visit
        void requestClose()
    })
})
function unlock() {
    if (!locked) return
    document.body.style.overflow = previousOverflow
    locked = false
    window.removeEventListener('beforeunload', beforeUnload)
}
function beforeUnload(event: BeforeUnloadEvent) {
    if (props.open && props.dirty) {
        event.preventDefault()
        event.returnValue = ''
    }
}
async function requestClose() {
    if (props.busy) return
    if (!props.dirty) {
        emit('close')
        return
    }
    returnToEditor = document.activeElement as HTMLElement
    confirming.value = true
    await nextTick()
    keepButton.value?.focus()
}
async function keepEditing() {
    pendingVisit = null
    confirming.value = false
    await nextTick()
    if (returnToEditor?.isConnected) returnToEditor.focus()
    else heading.value?.focus()
}
function discard() {
    if (props.busy) return
    confirming.value = false
    const visit = pendingVisit
    pendingVisit = null
    emit('close')
    if (visit) {
        allowVisit = true
        try {
            router.visit(visit.url, { ...visit })
        } finally {
            allowVisit = false
        }
    }
}
function escape() {
    if (props.busy) return
    if (confirming.value) void keepEditing()
    else void requestClose()
}
watch(
    () => props.open,
    async (open) => {
        await nextTick()
        if (props.open !== open || !dialog.value) return
        if (open && !dialog.value.open) {
            opener = document.activeElement as HTMLElement
            previousOverflow = document.body.style.overflow
            locked = true
            document.body.style.overflow = 'hidden'
            window.addEventListener('beforeunload', beforeUnload)
            confirming.value = false
            dialog.value.showModal()
            heading.value?.focus()
        } else if (!open && dialog.value.open) {
            dialog.value.close()
            confirming.value = false
            unlock()
            if (opener?.isConnected) opener.focus()
        }
    },
    { immediate: true }
)
onBeforeUnmount(() => {
    removeVisitListener?.()
    dialog.value?.close()
    unlock()
    if (opener?.isConnected) opener.focus()
})
defineExpose({ requestClose })
</script>
<template>
    <dialog
        ref="dialog"
        class="admin-drawer"
        :aria-labelledby="titleId"
        :aria-describedby="description ? descriptionId : undefined"
        :aria-busy="busy"
        @cancel.prevent="escape"
    >
        <div class="admin-drawer-frame">
            <header class="admin-drawer-header">
                <div>
                    <p class="eyebrow">Workspace editor</p>
                    <h2 :id="titleId" ref="heading" tabindex="-1">
                        {{ title }}
                    </h2>
                    <p
                        v-if="description"
                        :id="descriptionId"
                        class="admin-drawer-description"
                    >
                        {{ description }}
                    </p>
                </div>
                <button
                    type="button"
                    class="admin-drawer-close"
                    :disabled="busy"
                    :aria-label="`Close ${title}`"
                    @click="requestClose"
                >
                    <SokoIcon name="close" />
                </button>
            </header>
            <section
                v-if="confirming"
                class="admin-drawer-discard"
                aria-label="Unsaved changes"
            >
                <p role="alert">
                    <strong>Discard unsaved changes?</strong> Your changes have
                    not been saved.
                </p>
                <div>
                    <button
                        ref="keepButton"
                        type="button"
                        class="button button-navy"
                        @click="keepEditing"
                    >
                        Keep editing</button
                    ><button
                        type="button"
                        class="button button-outline"
                        :disabled="busy"
                        @click="discard"
                    >
                        Discard changes
                    </button>
                </div>
            </section>
            <div
                v-if="open"
                class="admin-drawer-body"
                :inert="confirming || undefined"
            >
                <slot :request-close="requestClose" />
            </div>
            <footer
                v-if="open && $slots.footer"
                class="admin-drawer-footer"
                :inert="confirming || undefined"
            >
                <slot name="footer" :request-close="requestClose" />
            </footer>
        </div>
    </dialog>
</template>
<style scoped>
.admin-drawer {
    position: fixed;
    inset: 0 0 0 auto;
    margin: 0;
    width: min(760px, 94vw);
    max-width: 100vw;
    height: 100dvh;
    max-height: 100dvh;
    padding: 0;
    overflow: hidden;
    border: 0;
    border-left: 1px solid var(--soko-line);
    border-radius: var(--soko-radius) 0 0 var(--soko-radius);
    background: var(--soko-surface);
    color: var(--soko-ink);
    font-family: var(--soko-font);
}
.admin-drawer::backdrop {
    background: color-mix(in srgb, var(--soko-ink) 42%, transparent);
}
.admin-drawer-frame {
    display: flex;
    flex-direction: column;
    height: 100%;
    min-height: 0;
}
.admin-drawer-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 20px;
    padding: 24px;
    border-bottom: 1px solid var(--soko-line);
    flex: none;
}
.admin-drawer-header h2 {
    margin: 4px 0;
    font-size: 22px;
    line-height: 1.3;
    font-family: var(--soko-font-display, var(--soko-font));
}
.admin-drawer-description {
    color: var(--soko-muted);
    font-size: 13px;
    line-height: 1.6;
    margin: 8px 0 0;
}
.admin-drawer-close {
    width: 44px;
    height: 44px;
    flex: none;
    display: grid;
    place-items: center;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    background: var(--soko-surface);
    color: inherit;
}
.admin-drawer-close:focus-visible,
.admin-drawer-header h2:focus-visible {
    outline: 2px solid var(--soko-focus);
    outline-offset: 3px;
}
.admin-drawer-body {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    overscroll-behavior: contain;
    padding: 24px;
}
.admin-drawer-footer {
    flex: none;
    border-top: 1px solid var(--soko-line);
    padding: 16px 24px max(16px, env(safe-area-inset-bottom));
    background: var(--soko-surface);
}
.admin-drawer-discard {
    padding: 16px 24px;
    border-bottom: 1px solid var(--soko-line);
    background: var(--soko-canvas);
}
.admin-drawer-discard p {
    font-size: 14px;
    margin: 0 0 12px;
}
.admin-drawer-discard > div {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
}
@media (max-width: 640px) {
    .admin-drawer {
        inset: auto 0 0;
        width: 100%;
        height: 92dvh;
        max-height: 92dvh;
        border-left: 0;
        border-top: 1px solid var(--soko-line);
        border-radius: var(--soko-radius) var(--soko-radius) 0 0;
    }
    .admin-drawer-header {
        padding: 16px;
        gap: 12px;
    }
    .admin-drawer-header h2 {
        font-size: 19px;
    }
    .admin-drawer-body {
        padding: 16px;
    }
    .admin-drawer-footer {
        padding: 12px 16px max(12px, env(safe-area-inset-bottom));
    }
}
@media (prefers-reduced-motion: no-preference) {
    .admin-drawer[open] {
        animation: admin-drawer-enter 160ms ease-out;
    }
    @keyframes admin-drawer-enter {
        from {
            opacity: 0.7;
            transform: translateX(24px);
        }
        to {
            opacity: 1;
            transform: translateX(0);
        }
    }
    @media (max-width: 640px) {
        @keyframes admin-drawer-enter {
            from {
                opacity: 0.7;
                transform: translateY(24px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }
    }
}
</style>
