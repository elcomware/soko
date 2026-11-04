import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { createToastTimer } from '../../resources/js/design/toastTimer.ts'
function harness() {
    let id = 0
    const pending = new Map(),
        dismissed = []
    const timer = createToastTimer((m) => dismissed.push(m), {
        schedule: (fn) => {
            pending.set(++id, fn)
            return id
        },
        unschedule: (key) => pending.delete(key),
    })
    const fire = () => {
        const callbacks = [...pending.values()]
        pending.clear()
        callbacks.forEach((fn) => fn())
    }
    return { timer, fire, pending, dismissed }
}
test('success dismisses after timer and only current message can expire', () => {
    const h = harness()
    h.timer.set('Saved')
    const stale = [...h.pending.values()][0]
    h.timer.set('Updated')
    stale()
    assert.deepEqual(h.dismissed, [])
    h.fire()
    assert.deepEqual(h.dismissed, ['Updated'])
})
test('hover/focus pause and resume keeps success notice readable', () => {
    const h = harness()
    h.timer.set('Saved')
    h.timer.pause()
    h.fire()
    assert.deepEqual(h.dismissed, [])
    h.timer.resume()
    h.fire()
    assert.deepEqual(h.dismissed, ['Saved'])
})
test('cleanup and empty messages leave no delayed notifications', () => {
    const h = harness()
    h.timer.set('Saved')
    h.timer.dispose()
    h.fire()
    assert.deepEqual(h.dismissed, [])
    h.timer.set('')
    assert.equal(h.pending.size, 0)
})
test('shared toast keeps correct live roles and errors are not timed', () => {
    const source = readFileSync(
        new URL(
            '../../resources/js/components/SharedFlashNotifications.vue',
            import.meta.url
        ),
        'utf8'
    )
    assert.match(source, /role="status"/)
    assert.match(source, /role="alert"/)
    assert.match(source, /Dismiss success message/)
    assert.match(source, /Dismiss error message/)
    assert.equal(source.includes('notice-success'), false)
    assert.match(source, /flash.value.success/)
    assert.equal(
        source.includes('createToastTimer(message => { dismissedError'),
        false
    )
})
test('category controls flank an unboxed icon-and-label rail', () => {
    const source = readFileSync(
        new URL('../../resources/js/Pages/Store/Index.vue', import.meta.url),
        'utf8'
    )
    const prev = source.indexOf('aria-label="Previous categories"'),
        rail = source.indexOf('id="soko-category-rail"'),
        next = source.indexOf('aria-label="Next categories"')
    assert.ok(prev >= 0 && rail >= 0 && next >= 0)
    assert.ok(prev < rail && rail < next)
    assert.ok(source.includes('category-rail-frame'))
    const css = readFileSync(
        new URL('../../resources/css/shop.css', import.meta.url),
        'utf8'
    )
    assert.match(
        css,
        /\.shop-category-tile,\s*\.shop-category-tile\[aria-current=['"]page['"]\]\s*\{[^}]*border:\s*0/
    )
})

test('hover exit cannot resume while keyboard focus remains inside toast', () => {
    const h = harness()
    h.timer.set('Saved')
    h.timer.pause('hover')
    h.timer.pause('focus')
    h.timer.resume('hover')
    h.fire()
    assert.deepEqual(h.dismissed, [])
    h.timer.resume('focus')
    h.fire()
    assert.deepEqual(h.dismissed, ['Saved'])
})
test('cleared notification releases old interaction pause for the next message', () => {
    const h = harness()
    h.timer.set('Saved')
    h.timer.pause('hover')
    h.timer.set('')
    h.timer.set('Updated')
    h.fire()
    assert.deepEqual(h.dismissed, ['Updated'])
})
