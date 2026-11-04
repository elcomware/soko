import { test } from 'node:test'
import assert from 'node:assert/strict'
import { createDebouncedSearch } from '../../resources/js/design/debouncedSearch.ts'
function harness(initial = '') {
    let time = 0,
        id = 0
    const pending = new Map(),
        visits = []
    const search = createDebouncedSearch({
        initial,
        delay: 350,
        schedule: (fn, delay) => {
            pending.set(++id, { fn, at: time + delay })
            return id
        },
        unschedule: (key) => pending.delete(key),
        visit: (query, request) => visits.push({ query, request }),
    })
    const tick = (ms) => {
        time += ms
        for (const [key, timer] of pending)
            if (timer.at <= time) {
                pending.delete(key)
                timer.fn()
            }
    }
    return { search, visits, tick, pending }
}
test('typing debounces into one request after a short pause', () => {
    const h = harness()
    h.search.update('c')
    h.tick(100)
    h.search.update('co')
    h.tick(100)
    h.search.update('coffee')
    h.tick(349)
    assert.equal(h.visits.length, 0)
    h.tick(1)
    assert.equal(h.visits.length, 1)
    assert.equal(h.visits[0].query, 'coffee')
})
test('clear executes immediately and Enter cancels the pending duplicate', () => {
    const h = harness('coffee')
    h.search.update('')
    assert.equal(h.visits[0].query, '')
    h.search.sync('')
    h.search.update('tea')
    h.search.submit('tea')
    h.tick(500)
    assert.equal(h.visits.length, 2)
    assert.equal(h.visits[1].query, 'tea')
})
test('new input cancels stale visits and stale cancellation tokens immediately cancel', () => {
    const h = harness()
    h.search.update('coffee')
    h.tick(350)
    let cancelled = 0
    h.visits[0].request.setCancel(() => cancelled++)
    h.search.update('tea')
    assert.equal(cancelled, 1)
    assert.equal(h.visits[0].request.isCurrent(), false)
    h.visits[0].request.setCancel(() => cancelled++)
    assert.equal(cancelled, 2)
    h.tick(350)
    assert.equal(h.visits[1].query, 'tea')
    assert.equal(h.visits[1].request.isCurrent(), true)
})
test('navigation/back sync cancels pending edits without another request', () => {
    const h = harness('coffee')
    h.search.update('tea')
    h.search.sync('coffee')
    h.tick(500)
    assert.equal(h.visits.length, 0)
    h.search.update('coffee')
    h.tick(500)
    assert.equal(h.visits.length, 0)
    h.search.update('spare')
    h.search.cancel()
    h.tick(500)
    assert.equal(h.visits.length, 0)
})
