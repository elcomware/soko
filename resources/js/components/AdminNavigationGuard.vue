<template>
    <SokoDialog
        :open="open"
        :title="busy ? 'Save in progress' : 'Leave unsaved changes?'"
        @close="open = false"
    >
        <p v-if="busy">
            Please wait for the current save to finish before leaving.
        </p>
        <p v-else>
            Your edits have not been saved. Keep editing, or discard them before
            leaving this page.
        </p>
        <div class="admin-guard-actions">
            <button
                type="button"
                class="button button-outline"
                @click="open = false"
            >
                Keep editing
            </button>
            <button
                type="button"
                class="button button-navy"
                :disabled="busy"
                @click="leave"
            >
                Discard and leave
            </button>
        </div>
    </SokoDialog>
</template>
<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { router } from '@inertiajs/vue3'
import SokoDialog from './SokoDialog.vue'
const props = defineProps<{ dirty: boolean; busy?: boolean }>()
const emit = defineEmits<{ discard: [] }>()
const open = ref(false)
let nextVisit: any = null,
    allowed = false,
    removeListener: (() => void) | undefined
function beforeUnload(event: BeforeUnloadEvent) {
    if (props.dirty || props.busy) {
        event.preventDefault()
        event.returnValue = ''
    }
}
onMounted(() => {
    removeListener = router.on('before', (event) => {
        if (
            allowed ||
            (!props.dirty && !props.busy) ||
            event.detail.visit.method !== 'get'
        )
            return
        event.preventDefault()
        nextVisit = event.detail.visit
        open.value = true
    })
    window.addEventListener('beforeunload', beforeUnload)
})
onBeforeUnmount(() => {
    removeListener?.()
    window.removeEventListener('beforeunload', beforeUnload)
})
function leave() {
    const visit = nextVisit
    if (!visit || props.busy) return
    emit('discard')
    open.value = false
    allowed = true
    try {
        router.visit(visit.url, { ...visit })
    } finally {
        allowed = false
        nextVisit = null
    }
}
</script>
<style scoped>
.admin-guard-actions {
    display: flex;
    flex-wrap: wrap;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 24px;
}
</style>
