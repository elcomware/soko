import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import { reactive, effectScope, watch } from 'vue'
const source = readFileSync(
    new URL('../../resources/js/design/adminDraftRecovery.ts', import.meta.url),
    'utf8'
)
    .replace(/^import .*$/gm, '')
    .replace('export function', 'function')
function harness(history = new Map()) {
    const mounted = [],
        scope = effectScope()
    const context = {
        watch,
        onMounted: (fn) => mounted.push(fn),
        router: {
            restore: (key) => history.get(key),
            remember: (value, key) => history.set(key, value),
        },
    }
    vm.runInNewContext(
        ts.transpileModule(
            source + '\nglobalThis.recover = useAdminDraftRecovery',
            {
                compilerOptions: {
                    target: ts.ScriptTarget.ES2022,
                    module: ts.ModuleKind.None,
                },
            }
        ).outputText,
        context
    )
    return {
        history,
        bind: (options) => scope.run(() => context.recover(options)),
        mount: () => mounted.forEach((fn) => fn()),
        stop: () => scope.stop(),
    }
}
test('history recovery preserves unsaved snapshot and concurrency version, excludes unexpected keys', () => {
    const key = 'soko:brand-draft:actor'
    const h = harness(
        new Map([
            [
                key,
                {
                    display_name: 'Unsaved',
                    version: 1,
                    secret: 'must-not-restore',
                },
            ],
        ])
    )
    const form = reactive({ display_name: 'Current', version: 2 })
    let baseline = JSON.stringify(form)
    h.bind({
        key,
        allowed: () => true,
        read: () => ({ ...form }),
        apply: (draft) => Object.assign(form, draft),
        dirty: () => JSON.stringify(form) !== baseline,
    })
    h.mount()
    assert.equal(form.display_name, 'Unsaved')
    assert.equal(form.version, 1)
    assert.equal('secret' in form, false)
    form.display_name = 'Latest draft'
    assert.equal(h.history.get(key).display_name, 'Latest draft')
    form.display_name = 'Current'
    form.version = 2
    assert.equal(h.history.get(key), null)
    h.stop()
})
test('recovery rechecks authorization and isolates actor/merchant keys', () => {
    const originalKey = 'soko:payment-draft:actor:merchant-a'
    const history = new Map([[originalKey, { country: 'KE' }]])
    const h = harness(history)
    const form = reactive({ country: 'US' })
    let allowed = false
    h.bind({
        key: originalKey,
        allowed: () => allowed,
        read: () => ({ ...form }),
        apply: (draft) => Object.assign(form, draft),
        dirty: () => true,
    })
    h.bind({
        key: 'soko:payment-draft:actor:merchant-b',
        allowed: () => true,
        read: () => ({ ...form }),
        apply: (draft) => Object.assign(form, draft),
        dirty: () => true,
    })
    h.mount()
    assert.equal(form.country, 'US')
    form.country = 'GB'
    assert.equal(history.get(originalKey).country, 'KE')
    assert.equal(
        history.get('soko:payment-draft:actor:merchant-b').country,
        'GB'
    )
    h.stop()
})

test('page navigation cannot discard or leave while a save is pending', async () => {
    const { ref } = await import('vue')
    const component = readFileSync(
        new URL(
            '../../resources/js/components/AdminNavigationGuard.vue',
            import.meta.url
        ),
        'utf8'
    )
    const script = component
        .match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
        .replace(/^import .*$/gm, '')
    const props = { dirty: true, busy: true },
        visits = [],
        emitted = []
    let before
    const context = {
        ref,
        defineProps: () => props,
        defineEmits: () => (event) => emitted.push(event),
        onMounted: (fn) => fn(),
        onBeforeUnmount: () => {},
        window: { addEventListener: () => {}, removeEventListener: () => {} },
        router: {
            on: (_, callback) => {
                before = callback
                return () => {}
            },
            visit: (...args) => visits.push(args),
        },
    }
    vm.runInNewContext(
        ts.transpileModule(script + '\nglobalThis.leave = leave', {
            compilerOptions: {
                target: ts.ScriptTarget.ES2022,
                module: ts.ModuleKind.None,
            },
        }).outputText,
        context
    )
    let prevented = false
    before({
        detail: { visit: { method: 'get', url: '/dashboard' } },
        preventDefault: () => {
            prevented = true
        },
    })
    context.leave()
    assert.equal(prevented, true)
    assert.equal(visits.length, 0)
    assert.equal(emitted.length, 0)
    props.busy = false
    context.leave()
    assert.equal(visits.length, 1)
    assert.deepEqual(emitted, ['discard'])
})
