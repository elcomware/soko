<script setup lang="ts">
import { computed, ref, useId, watch, onBeforeUnmount } from 'vue'
type Column = {
    key: string
    label: string
    sortable?: boolean
    width?: number
    minWidth?: number
    maxWidth?: number
}
const props = withDefaults(
    defineProps<{
        caption: string
        empty?: boolean
        emptyText?: string
        busy?: boolean
        columns?: Column[]
        sort?: { key: string; direction: 'asc' | 'desc' }
        resizableColumns?: boolean
        viewportHeight?: number
        resizableHeight?: boolean
    }>(),
    {
        empty: false,
        emptyText: 'No records to show.',
        busy: false,
        columns: () => [],
        resizableColumns: false,
        resizableHeight: false,
    }
)
const emit = defineEmits<{
    sort: [value: { key: string; direction: 'asc' | 'desc' }]
    'column-resize': [value: { key: string; width: number }]
    'update:viewportHeight': [value: number]
}>()
const captionId = useId()
const regionId = useId()
const widths = ref<Record<string, number>>({})
const clamp = (value: number, min: number, max: number) =>
    Math.max(min, Math.min(max, Number.isFinite(value) ? value : min))
const height = ref(clamp(props.viewportHeight ?? 480, 240, 900))
const bounded = computed(
    () => props.viewportHeight !== undefined || props.resizableHeight
)
const tableWidth = computed(() =>
    props.columns.reduce((sum, col) => sum + columnWidth(col), 0)
)
function limits(col: Column) {
    const min = clamp(col.minWidth ?? 100, 64, 600)
    return { min, max: clamp(col.maxWidth ?? 640, min, 1200) }
}
function columnWidth(col: Column) {
    const { min, max } = limits(col)
    return clamp(widths.value[col.key] ?? col.width ?? 180, min, max)
}
function setColumn(col: Column, value: number, notify = true) {
    const { min, max } = limits(col)
    const width = clamp(value, min, max)
    widths.value = { ...widths.value, [col.key]: width }
    if (notify) emit('column-resize', { key: col.key, width })
}
function setHeight(value: number, notify = true) {
    height.value = clamp(value, 240, 900)
    if (notify) emit('update:viewportHeight', height.value)
}
function sortColumn(col: Column) {
    if (props.busy || !col.sortable) return
    emit('sort', {
        key: col.key,
        direction:
            props.sort?.key === col.key && props.sort.direction === 'asc'
                ? 'desc'
                : 'asc',
    })
}
function keyValue(
    event: KeyboardEvent,
    value: number,
    min: number,
    max: number,
    axis: 'column' | 'height'
) {
    const step = event.shiftKey ? 40 : 10
    const back = axis === 'column' ? 'ArrowLeft' : 'ArrowUp',
        forward = axis === 'column' ? 'ArrowRight' : 'ArrowDown'
    if (![back, forward, 'Home', 'End'].includes(event.key)) return null
    event.preventDefault()
    return event.key === 'Home'
        ? min
        : event.key === 'End'
          ? max
          : event.key === back
            ? value - step
            : value + step
}
function resizeColumnKey(col: Column, event: KeyboardEvent) {
    if (props.busy || !props.resizableColumns) return
    const { min, max } = limits(col)
    const next = keyValue(event, columnWidth(col), min, max, 'column')
    if (next !== null) setColumn(col, next)
}
function resizeHeightKey(event: KeyboardEvent) {
    if (props.busy || !props.resizableHeight) return
    const next = keyValue(event, height.value, 240, 900, 'height')
    if (next !== null) setHeight(next)
}
let drag: {
    key: string | null
    pointerId: number
    start: number
    value: number
    target: HTMLElement
} | null = null
function startResize(key: string | null, event: PointerEvent) {
    if (
        props.busy ||
        event.button !== 0 ||
        event.isPrimary === false ||
        (key === null ? !props.resizableHeight : !props.resizableColumns)
    )
        return
    const col = props.columns.find((c) => c.key === key)
    if (key !== null && !col) return
    event.preventDefault()
    const target = event.currentTarget as HTMLElement
    target.setPointerCapture(event.pointerId)
    drag = {
        key,
        pointerId: event.pointerId,
        start: key === null ? event.clientY : event.clientX,
        value: key === null ? height.value : columnWidth(col!),
        target,
    }
}
function moveResize(event: PointerEvent) {
    if (!drag || event.pointerId !== drag.pointerId) return
    if (props.busy) {
        cancelResize()
        return
    }
    const value =
        drag.value +
        (drag.key === null ? event.clientY : event.clientX) -
        drag.start
    if (drag.key === null) setHeight(value, false)
    else {
        const col = props.columns.find((c) => c.key === drag!.key)
        if (col) setColumn(col, value, false)
    }
}
function releaseDrag() {
    const current = drag
    drag = null
    if (current?.target.hasPointerCapture(current.pointerId))
        current.target.releasePointerCapture(current.pointerId)
}
function endResize(event: PointerEvent) {
    if (!drag || event.pointerId !== drag.pointerId) return
    if (drag.key === null) emit('update:viewportHeight', height.value)
    else {
        const col = props.columns.find((c) => c.key === drag!.key)
        if (col)
            emit('column-resize', { key: col.key, width: columnWidth(col) })
    }
    releaseDrag()
}
function cancelResize(event?: PointerEvent) {
    if (!drag || (event && event.pointerId !== drag.pointerId)) return
    if (drag.key === null) setHeight(drag.value, false)
    else {
        const col = props.columns.find((c) => c.key === drag!.key)
        if (col) setColumn(col, drag.value, false)
    }
    releaseDrag()
}
watch(
    () => props.viewportHeight,
    (value) => {
        if (value !== undefined && !drag) setHeight(value, false)
    }
)
watch(
    () => props.busy,
    (value) => {
        if (value) cancelResize()
    }
)
watch(
    () => [props.resizableColumns, props.resizableHeight],
    () => cancelResize()
)
onBeforeUnmount(cancelResize)
</script>
<template>
    <div class="admin-data-table" :aria-busy="busy">
        <div v-if="$slots.toolbar" class="admin-data-toolbar">
            <slot name="toolbar" />
        </div>
        <details v-if="$slots.filters" class="admin-data-filters">
            <summary>Filters</summary>
            <div><slot name="filters" /></div>
        </details>
        <div
            v-if="$slots.bulk"
            class="admin-data-bulk"
            role="region"
            aria-label="Selected record actions"
        >
            <slot name="bulk" />
        </div>
        <p v-if="empty" class="admin-data-empty" role="status">
            {{ emptyText }}
        </p>
        <div
            v-else
            class="admin-data-scroll"
            :id="regionId"
            :style="
                bounded ? { maxHeight: `min(${height}px, 75dvh)` } : undefined
            "
            role="region"
            :aria-labelledby="captionId"
            tabindex="0"
        >
            <table
                :class="{ 'has-columns': columns.length }"
                :style="
                    columns.length ? { width: `${tableWidth}px` } : undefined
                "
            >
                <colgroup v-if="columns.length">
                    <col
                        v-for="column in columns"
                        :key="column.key"
                        :style="{ width: `${columnWidth(column)}px` }"
                    />
                </colgroup>
                <caption :id="captionId">
                    {{
                        caption
                    }}
                </caption>
                <thead>
                    <slot name="header">
                        <tr>
                            <th
                                v-for="column in columns"
                                :key="column.key"
                                scope="col"
                                :aria-sort="
                                    column.sortable
                                        ? sort?.key === column.key
                                            ? sort.direction === 'asc'
                                                ? 'ascending'
                                                : 'descending'
                                            : 'none'
                                        : undefined
                                "
                            >
                                <button
                                    v-if="column.sortable"
                                    type="button"
                                    class="admin-column-sort"
                                    :disabled="busy"
                                    @click="sortColumn(column)"
                                >
                                    {{ column.label }}
                                    <span aria-hidden="true">{{
                                        sort?.key === column.key
                                            ? sort.direction === 'asc'
                                                ? '↑'
                                                : '↓'
                                            : '↕'
                                    }}</span></button
                                ><span v-else>{{ column.label }}</span>
                                <span
                                    v-if="resizableColumns"
                                    class="admin-column-resize"
                                    role="separator"
                                    aria-orientation="vertical"
                                    :tabindex="busy ? -1 : 0"
                                    :aria-label="`Resize ${column.label} column`"
                                    :aria-disabled="busy"
                                    :aria-valuemin="limits(column).min"
                                    :aria-valuemax="limits(column).max"
                                    :aria-valuenow="columnWidth(column)"
                                    :aria-controls="regionId"
                                    @keydown="resizeColumnKey(column, $event)"
                                    @pointerdown="
                                        startResize(column.key, $event)
                                    "
                                    @pointermove="moveResize"
                                    @pointerup="endResize"
                                    @pointercancel="cancelResize"
                                    @lostpointercapture="cancelResize"
                                />
                            </th>
                        </tr>
                    </slot>
                </thead>
                <tbody>
                    <slot />
                </tbody>
            </table>
        </div>
        <div
            v-if="resizableHeight && !empty"
            class="admin-height-resize"
            role="separator"
            aria-orientation="horizontal"
            :tabindex="busy ? -1 : 0"
            aria-label="Resize table viewport height"
            :aria-disabled="busy"
            aria-valuemin="240"
            aria-valuemax="900"
            :aria-valuenow="height"
            :aria-controls="regionId"
            @keydown="resizeHeightKey"
            @pointerdown="startResize(null, $event)"
            @pointermove="moveResize"
            @pointerup="endResize"
            @pointercancel="cancelResize"
            @lostpointercapture="cancelResize"
        >
            <span aria-hidden="true">⠿</span><span>Table height</span>
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
    overflow-y: auto;
    overscroll-behavior: contain;
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
    position: sticky;
    top: 0;
    z-index: 1;
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
.has-columns {
    table-layout: fixed;
}
.has-columns :deep(td),
.has-columns :deep(tbody th) {
    overflow-wrap: anywhere;
}
.admin-column-sort {
    display: flex;
    align-items: center;
    gap: 8px;
    min-height: 32px;
    max-width: 100%;
    white-space: normal;
    text-align: left;
    background: transparent;
    color: inherit;
    border: 0;
    padding: 0;
    font: inherit;
}
.admin-column-resize {
    position: absolute;
    right: 0;
    top: 0;
    bottom: 0;
    width: 12px;
    cursor: col-resize;
    touch-action: none;
}
.admin-column-resize::after {
    content: '';
    position: absolute;
    inset: 25% 5px;
    background: var(--soko-line);
}
.admin-column-resize:focus-visible,
.admin-height-resize:focus-visible {
    outline: 2px solid var(--soko-focus);
    outline-offset: -2px;
}
.admin-height-resize {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 8px;
    min-height: 32px;
    color: var(--soko-muted);
    font-size: 11px;
    cursor: row-resize;
    touch-action: none;
    user-select: none;
}
.admin-data-filters {
    margin-bottom: 12px;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    padding: 12px 16px;
}
.admin-data-filters summary {
    cursor: pointer;
    font-weight: 650;
    font-size: 13px;
}
.admin-data-filters > div {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    align-items: end;
    padding-top: 12px;
}
.admin-data-bulk {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    align-items: center;
    padding: 12px 16px;
    margin-bottom: 12px;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    background: var(--soko-canvas);
}
</style>
