export interface SearchVisit {
    setCancel: (cancel: () => void) => void
    isCurrent: () => boolean
}
/** Own only search timers/visits; navigation sync cancels both without another request. */
export function createDebouncedSearch(options: {
    initial?: string
    delay?: number
    visit: (query: string, request: SearchVisit) => void
    schedule?: (callback: () => void, delay: number) => any
    unschedule?: (id: any) => void
}) {
    const schedule = options.schedule || setTimeout,
        unschedule = options.unschedule || clearTimeout
    let timer: any,
        cancelVisit: (() => void) | undefined,
        generation = 0,
        committed = (options.initial || '').trim(),
        desired = committed
    function cancel() {
        generation++
        if (timer !== undefined) unschedule(timer)
        timer = undefined
        cancelVisit?.()
        cancelVisit = undefined
    }
    function send() {
        timer = undefined
        if (desired === committed) return
        const id = generation
        options.visit(desired, {
            isCurrent: () => id === generation,
            setCancel: (cancel) => {
                if (id !== generation) cancel()
                else cancelVisit = cancel
            },
        })
    }
    return {
        update(value: string) {
            cancel()
            desired = value.trim()
            if (desired === committed) return
            if (!desired) send()
            else timer = schedule(send, options.delay ?? 350)
        },
        submit(value: string) {
            cancel()
            desired = value.trim()
            send()
        },
        sync(value: string) {
            cancel()
            committed = value.trim()
            desired = committed
        },
        cancel,
    }
}
