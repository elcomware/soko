import { computed, ref, watch } from 'vue'
// The caller supplies only authorized IDs from the loaded page, never a global total.
export function useAdminSelection(
    eligibleKeys: () => string[],
    busy: () => boolean = () => false
) {
    const keys = computed(() => [...new Set(eligibleKeys())])
    const selected = ref<string[]>([])
    const all = computed(
        () =>
            keys.value.length > 0 &&
            keys.value.every((key) => selected.value.includes(key))
    )
    const mixed = computed(() => selected.value.length > 0 && !all.value)
    watch(
        keys,
        (allowed) => {
            selected.value = selected.value.filter((key) =>
                allowed.includes(key)
            )
        },
        { flush: 'sync' }
    )
    function toggle(key: string, checked: boolean) {
        if (busy() || !keys.value.includes(key)) return
        selected.value = checked
            ? [...new Set([...selected.value, key])]
            : selected.value.filter((item) => item !== key)
    }
    function toggleAll(checked: boolean) {
        if (!busy()) selected.value = checked ? [...keys.value] : []
    }
    function clear() {
        if (!busy()) selected.value = []
    }
    return { selected, all, mixed, toggle, toggleAll, clear }
}
