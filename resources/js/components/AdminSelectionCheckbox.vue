<script setup lang="ts">
const props = defineProps<{
    checked: boolean
    mixed?: boolean
    disabled?: boolean
    label: string
}>()
const emit = defineEmits<{ change: [checked: boolean] }>()
function change(event: Event) {
    if (!props.disabled)
        emit('change', (event.target as HTMLInputElement).checked)
}
</script>
<template>
    <label class="admin-select-record"
        ><input
            type="checkbox"
            :checked="checked"
            :indeterminate="mixed"
            :disabled="disabled"
            :aria-label="label"
            @change="change"
    /></label>
</template>
<style scoped>
.admin-select-record {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 32px;
    min-height: 40px;
}
.admin-select-record input {
    width: 18px;
    height: 18px;
    accent-color: var(--soko-navy);
}
.admin-select-record input:focus-visible {
    outline: 2px solid var(--soko-focus);
    outline-offset: 3px;
}
</style>
