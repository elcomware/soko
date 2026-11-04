<template>
    <div
        v-if="flash.success || flash.error"
        class="soko-toast-region"
        aria-label="Notifications"
    >
        <div
            v-if="flash.success"
            role="status"
            class="soko-toast soko-toast-success"
            @mouseenter="successToast.pause('hover')"
            @mouseleave="successToast.resume('hover')"
            @focusin="successToast.pause('focus')"
            @focusout="onToastFocusOut"
        >
            <SokoIcon name="check" /><span>{{ flash.success }}</span
            ><button
                class="icon-button"
                type="button"
                aria-label="Dismiss success message"
                @click="dismissedSuccess = flash.success"
            >
                <SokoIcon name="close" />
            </button>
        </div>
        <div
            v-if="flash.error"
            role="alert"
            class="soko-toast soko-toast-error"
        >
            <SokoIcon name="warning" /><span>{{ flash.error }}</span
            ><button
                class="icon-button"
                type="button"
                aria-label="Dismiss error message"
                @click="dismissedError = flash.error"
            >
                <SokoIcon name="close" />
            </button>
        </div>
    </div>
</template>
<script setup lang="ts">
import { computed, ref, watch, onBeforeUnmount, onMounted } from 'vue'
import { usePage } from '@inertiajs/vue3'
import SokoIcon from './SokoIcon.vue'
import { createToastTimer } from '@/design/toastTimer'
const page = usePage()
function onToastFocusOut(event: FocusEvent) {
    if (
        !(event.currentTarget as HTMLElement).contains(
            event.relatedTarget as Node | null
        )
    )
        successToast.resume('focus')
}
const dismissedSuccess = ref('')
const dismissedError = ref('')
const flash = computed(() => {
    const value = (page.props.flash as Record<string, string>) || {}
    return {
        success: value.success === dismissedSuccess.value ? '' : value.success,
        error: value.error === dismissedError.value ? '' : value.error,
    }
})
const successToast = createToastTimer((message) => {
    dismissedSuccess.value = message
})
let mounted = false
onMounted(() => {
    mounted = true
    successToast.set(flash.value.success || '')
})
watch(
    () => flash.value.success,
    (message) => {
        if (mounted) successToast.set(message || '')
    }
)
onBeforeUnmount(() => successToast.dispose())
watch(
    () => page.props.flash,
    () => {
        dismissedSuccess.value = ''
        dismissedError.value = ''
    }
)
</script>
