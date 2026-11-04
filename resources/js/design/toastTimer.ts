/** Success notices expire gently; focused/hovered notices stay readable. Errors never use this timer. */
export function createToastTimer(
    onDismiss: (message: string) => void,
    options: {
        delay?: number
        schedule?: (fn: () => void, delay: number) => any
        unschedule?: (id: any) => void
    } = {}
) {
    const schedule = options.schedule || setTimeout,
        unschedule = options.unschedule || clearTimeout
    let timer: any,
        message = '',
        generation = 0
    const pauseReasons = new Set<string>()
    function clear() {
        generation++
        if (timer !== undefined) unschedule(timer)
        timer = undefined
    }
    function start() {
        clear()
        if (!message || pauseReasons.size) return
        const current = message,
            id = generation
        timer = schedule(() => {
            if (id === generation) onDismiss(current)
        }, options.delay ?? 8000)
    }
    return {
        set(value: string) {
            if (!value) pauseReasons.clear()
            message = value
            start()
        },
        pause(reason = 'interaction') {
            pauseReasons.add(reason)
            clear()
        },
        resume(reason = 'interaction') {
            pauseReasons.delete(reason)
            start()
        },
        dispose() {
            message = ''
            pauseReasons.clear()
            clear()
        },
    }
}
