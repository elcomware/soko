import { onMounted, watch } from 'vue'
import { router } from '@inertiajs/vue3'

// Browser-history recovery only: never include credentials, uploads or secret URLs.
// Callers supply an explicit safe field list and recheck current authorization.
export function useAdminDraftRecovery(options: {
    key: string
    allowed: () => boolean
    read: () => Record<string, unknown>
    apply: (draft: Record<string, unknown>) => void
    dirty: () => boolean
}) {
    let ready = false
    const clone = (value: unknown) => JSON.parse(JSON.stringify(value))
    onMounted(() => {
        if (!options.allowed()) return
        const saved = router.restore(options.key) as Record<
            string,
            unknown
        > | null
        if (saved && typeof saved === 'object' && !Array.isArray(saved)) {
            const current = options.read()
            const draft: Record<string, unknown> = {}
            for (const key of Object.keys(current)) {
                if (Object.prototype.hasOwnProperty.call(saved, key))
                    draft[key] = clone(saved[key])
            }
            options.apply(draft)
        }
        ready = true
    })
    watch(
        () => [options.read(), options.dirty()],
        () => {
            if (!ready || !options.allowed()) return
            router.remember(
                options.dirty() ? clone(options.read()) : null,
                options.key
            )
        },
        { deep: true, flush: 'sync' }
    )
}
