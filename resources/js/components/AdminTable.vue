<script setup lang="ts">
import { useId } from 'vue'
withDefaults(
    defineProps<{
        caption: string
        empty?: boolean
        emptyText?: string
        busy?: boolean
    }>(),
    { empty: false, emptyText: 'No records to show.', busy: false }
)
const captionId = useId()
</script>
<template>
    <div class="admin-data-table" :aria-busy="busy">
        <div v-if="$slots.toolbar" class="admin-data-toolbar">
            <slot name="toolbar" />
        </div>
        <p v-if="empty" class="admin-data-empty" role="status">
            {{ emptyText }}
        </p>
        <div
            v-else
            class="admin-data-scroll"
            role="region"
            :aria-labelledby="captionId"
            tabindex="0"
        >
            <table>
                <caption :id="captionId">
                    {{
                        caption
                    }}
                </caption>
                <thead>
                    <slot name="header" />
                </thead>
                <tbody>
                    <slot />
                </tbody>
            </table>
        </div>
        <div v-if="$slots.footer" class="admin-data-footer">
            <slot name="footer" />
        </div>
    </div>
</template>
<style scoped>
.admin-data-table {
    width: 100%;
    min-width: 0;
    color: var(--soko-ink);
}
.admin-data-scroll {
    overflow-x: auto;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    background: var(--soko-surface);
}
.admin-data-scroll:focus-visible {
    outline: 2px solid var(--soko-navy);
    outline-offset: 3px;
}
table {
    width: 100%;
    border-collapse: collapse;
    font-size: 13px;
    text-align: left;
}
caption {
    text-align: left;
    padding: 12px 16px;
    color: var(--soko-muted);
    font-weight: 650;
    font-size: 12px;
    border-bottom: 1px solid var(--soko-line);
}
:deep(th),
:deep(td) {
    padding: 14px 16px;
    border-bottom: 1px solid var(--soko-line);
    vertical-align: middle;
    line-height: 1.5;
}
:deep(thead th) {
    background: var(--soko-canvas);
    color: var(--soko-muted);
    font-size: 11px;
    font-weight: 750;
    letter-spacing: 0.04em;
    white-space: nowrap;
}
:deep(tbody th) {
    text-align: left;
    font-weight: 650;
}
:deep(tbody tr:last-child > *) {
    border-bottom: 0;
}
:deep(tbody tr:focus-within) {
    background: var(--soko-canvas);
}
.admin-data-toolbar,
.admin-data-footer {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 12px;
    padding: 12px 0;
}
.admin-data-empty {
    padding: 24px;
    border: 1px dashed var(--soko-line);
    border-radius: var(--soko-radius);
    color: var(--soko-muted);
}
</style>
