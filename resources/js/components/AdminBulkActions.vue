<script setup lang="ts">
import AdminActionMenu from './AdminActionMenu.vue'
const props = defineProps<{
    count: number
    busy?: boolean
    actions: {
        id: string
        label: string
        icon?: string
        disabled?: boolean
        danger?: boolean
    }[]
}>()
const emit = defineEmits<{ action: [id: string]; clear: [] }>()
function select(id: string) {
    if (
        !props.busy &&
        props.count > 0 &&
        props.actions.some((action) => action.id === id && !action.disabled)
    )
        emit('action', id)
}
</script>
<template>
    <div v-if="count > 0" class="admin-bulk-actions">
        <p role="status">{{ count }} selected on this page</p>
        <AdminActionMenu
            label="Actions for selected records"
            :items="actions"
            :disabled="busy"
            @select="select"
        />
        <button
            type="button"
            class="button button-outline"
            :disabled="busy"
            @click="!busy && emit('clear')"
        >
            Clear selection
        </button>
    </div>
</template>
<style scoped>
.admin-bulk-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
    font-size: 13px;
}
.admin-bulk-actions p {
    margin: 0;
}
</style>
