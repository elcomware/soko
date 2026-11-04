import { test } from 'node:test'
import assert from 'node:assert/strict'
import {
    categoryRailState,
    categoryRailTarget,
} from '../../resources/js/design/categoryRail.ts'
test('rail buttons reflect overflow and fractional edges', () => {
    assert.deepEqual(categoryRailState(0, 600, 300), {
        previous: false,
        next: true,
        maximum: 300,
    })
    assert.deepEqual(categoryRailState(299.7, 600, 300), {
        previous: true,
        next: false,
        maximum: 300,
    })
    assert.deepEqual(categoryRailState(0, 200, 300), {
        previous: false,
        next: false,
        maximum: 0,
    })
    assert.equal(categoryRailState(-10, 600, 300).previous, false)
})
test('keyboard moves by viewport and clamps without capturing unrelated keys', () => {
    assert.equal(categoryRailTarget('ArrowRight', 0, 600, 300), 240)
    assert.equal(categoryRailTarget('ArrowRight', 250, 600, 300), 300)
    assert.equal(categoryRailTarget('ArrowLeft', 20, 600, 300), 0)
    assert.equal(categoryRailTarget('Home', 200, 600, 300), 0)
    assert.equal(categoryRailTarget('End', 0, 600, 300), 300)
    assert.equal(categoryRailTarget('Tab', 0, 600, 300), null)
})
