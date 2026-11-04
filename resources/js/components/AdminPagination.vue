<script setup lang="ts">
import { computed } from 'vue'
const props = withDefaults(
    defineProps<{
        page: number
        perPage: number
        total: number
        lastPage: number
        from: number
        to: number
        busy?: boolean
        pageSizes?: number[]
    }>(),
    { busy: false, pageSizes: () => [25, 50, 100] }
)
const emit = defineEmits<{ page: [number]; perPage: [number] }>()
const current = computed(() => Math.max(1, Math.trunc(props.page) || 1))
const last = computed(() => Math.max(1, Math.trunc(props.lastPage) || 1))
const sizes = computed(() =>
    [...new Set([...props.pageSizes, props.perPage])]
        .filter((n) => Number.isInteger(n) && n > 0 && n <= 100)
        .sort((a, b) => a - b)
)
function go(page: number) {
    if (
        props.busy ||
        !Number.isInteger(page) ||
        page < 1 ||
        page > last.value ||
        page === current.value
    )
        return
    emit('page', page)
}
function changeSize(event: Event) {
    const value = Number((event.target as HTMLSelectElement).value)
    if (!props.busy && sizes.value.includes(value) && value !== props.perPage)
        emit('perPage', value)
}
</script>
<template>
    <nav
        class="admin-pagination"
        aria-label="Table pagination"
        :aria-busy="busy"
    >
        <p role="status">
            {{
                total
                    ? `${from}–${to} of ${total.toLocaleString()}`
                    : 'No matching records'
            }}
        </p>
        <label
            >Rows per page
            <select :value="perPage" :disabled="busy" @change="changeSize">
                <option v-for="size in sizes" :key="size" :value="size">
                    {{ size }}
                </option>
            </select></label
        >
        <div class="admin-page-controls">
            <button
                type="button"
                :disabled="busy || current <= 1"
                aria-label="First page"
                @click="go(1)"
            >
                First
            </button>
            <button
                type="button"
                :disabled="busy || current <= 1"
                aria-label="Previous page"
                @click="go(current - 1)"
            >
                Previous
            </button>
            <span>Page {{ current }} of {{ last }}</span>
            <button
                type="button"
                :disabled="busy || current >= last"
                aria-label="Next page"
                @click="go(current + 1)"
            >
                Next
            </button>
            <button
                type="button"
                :disabled="busy || current >= last"
                aria-label="Last page"
                @click="go(last)"
            >
                Last
            </button>
        </div>
    </nav>
</template>
<style scoped>
.admin-pagination {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 16px;
    width: 100%;
    font-size: 12px;
    color: var(--soko-muted);
}
.admin-pagination p {
    margin: 0;
    margin-right: auto;
}
.admin-pagination label,
.admin-page-controls {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
}
.admin-pagination button,
.admin-pagination select {
    min-height: 40px;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    background: var(--soko-surface);
    color: var(--soko-ink);
    padding: 8px 12px;
    font: inherit;
}
.admin-pagination button:disabled {
    opacity: 0.45;
}
.admin-pagination button:focus-visible,
.admin-pagination select:focus-visible {
    outline: 2px solid var(--soko-focus);
    outline-offset: 2px;
}
@media (max-width: 640px) {
    .admin-page-controls {
        width: 100%;
        justify-content: space-between;
    }
    .admin-page-controls button {
        min-height: 44px;
    }
    .admin-page-controls button:first-child,
    .admin-page-controls button:last-child {
        display: none;
    }
}
</style>
