<template>
    <Teleport to="body">
        <dialog
            ref="dialog"
            class="soko-dialog"
            :aria-labelledby="titleId"
            @cancel.prevent="emit('close')"
            @close="emit('close')"
        >
            <div class="soko-dialog-heading">
                <div>
                    <p class="eyebrow">Soko marketplace</p>
                    <h2 :id="titleId">{{ title }}</h2>
                </div>
                <button
                    type="button"
                    class="icon-button"
                    aria-label="Close dialog"
                    @click="emit('close')"
                >
                    <SokoIcon name="close" />
                </button>
            </div>
            <div v-if="open" class="soko-dialog-body"><slot /></div>
        </dialog>
    </Teleport>
</template>
<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, useId, watch } from 'vue'
import SokoIcon from './SokoIcon.vue'
const props = defineProps<{ open: boolean; title: string }>()
const emit = defineEmits<{ close: [] }>()
const dialog = ref<HTMLDialogElement | null>(null)
const titleId = useId()
let previousOverflow = ''
watch(
    () => props.open,
    async (open) => {
        await nextTick()
        if (open && !dialog.value?.open) {
            previousOverflow = document.body.style.overflow
            dialog.value?.showModal()
            document.body.style.overflow = 'hidden'
        } else if (!open && dialog.value?.open) {
            dialog.value.close()
            document.body.style.overflow = previousOverflow
        }
    },
    { immediate: true }
)
onBeforeUnmount(() => {
    if (dialog.value?.open) {
        dialog.value.close()
        document.body.style.overflow = previousOverflow
    }
})
</script>
