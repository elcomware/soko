import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import ts from 'typescript'
import * as Vue from 'vue'
import { parse, compileScript } from '@vue/compiler-sfc'
import { renderToString } from '@vue/server-renderer'
import { iconPaths, iconNames } from '../../resources/js/design/icons.ts'
import {
    commerceIllustrations,
    illustrationNames,
} from '../../resources/js/design/commerceIllustrations.ts'
const appearance = readFileSync(
    new URL('../../resources/js/design/iconAppearance.ts', import.meta.url),
    'utf8'
)
const js = ts
    .transpileModule(appearance, {
        compilerOptions: {
            target: ts.ScriptTarget.ES2022,
            module: ts.ModuleKind.ESNext,
        },
    })
    .outputText.replace(/^import[\s\S]*?;\s*$/gm, '')
    .replace(/export /g, '')
const { iconSurface, dimensionalIconNames } = new Function(
    'iconPaths',
    js + ';return {iconSurface,dimensionalIconNames}'
)(iconPaths)
function component(file, extras = {}) {
    const source = readFileSync(
        new URL('../../resources/js/components/' + file, import.meta.url),
        'utf8'
    )
    const compiled = compileScript(parse(source).descriptor, {
        id: file,
        inlineTemplate: true,
    })
    const ast = ts.createSourceFile(
        file,
        compiled.content,
        ts.ScriptTarget.Latest,
        true,
        ts.ScriptKind.TS
    )
    const env = { ...extras }
    for (const node of ast.statements)
        if (ts.isImportDeclaration(node) && node.moduleSpecifier.text === 'vue')
            for (const binding of node.importClause.namedBindings.elements)
                env[binding.name.text] =
                    Vue[binding.propertyName?.text || binding.name.text]
    const code = ts
        .transpileModule(compiled.content, {
            compilerOptions: {
                target: ts.ScriptTarget.ES2022,
                module: ts.ModuleKind.ESNext,
            },
        })
        .outputText.replace(/^import[\s\S]*?;\s*$/gm, '')
        .replace('export default', 'return')
    return new Function(...Object.keys(env), code)(...Object.values(env))
}
test('all 159 symbols support dimensional rendering; 100+ have semantic filled surfaces', () => {
    assert.equal(dimensionalIconNames.length, iconNames.length)
    assert.ok(iconNames.filter((n) => iconSurface(n)).length >= 100)
    for (const name of iconNames)
        assert.match(iconSurface(name), /^[MmLlHhVvCcSsQqTtAaZz0-9.,\s-]*$/)
})
test('12 unique editorial commerce scenes contain safe, layered vector detail', () => {
    assert.equal(illustrationNames.length, 12)
    assert.equal(
        new Set(
            Object.values(commerceIllustrations).map((l) => JSON.stringify(l))
        ).size,
        12
    )
    for (const layers of Object.values(commerceIllustrations)) {
        assert.ok(layers.length >= 4)
        for (const layer of layers)
            assert.match(layer.d, /^[MmLlHhVvCcSsQqTtAaZz0-9.,\s-]+$/)
    }
})
test('actual dimensional SVG render has unique gradient IDs, labels and safe fallback', async () => {
    const Icon = component('SokoIcon.vue', { iconPaths, iconSurface })
    const html = await renderToString(
        Vue.createSSRApp({
            render: () =>
                Vue.h('div', [
                    Vue.h(Icon, {
                        name: 'cart',
                        size: 48,
                        label: 'Shopping cart',
                    }),
                    Vue.h(Icon, { name: 'unknown', size: 48 }),
                ]),
        })
    )
    const ids = [...html.matchAll(/<linearGradient id="([^"]+)"/g)].map(
        (m) => m[1]
    )
    assert.equal(ids.length, 2)
    assert.equal(new Set(ids).size, 2)
    assert.match(html, /aria-label="Shopping cart"/)
    assert.match(html, /data-icon="help"/)
    assert.match(html, /aria-hidden="true"/)
    const compact = await renderToString(
        Vue.createSSRApp({
            render: () => Vue.h(Icon, { name: 'cart', size: 16 }),
        })
    )
    assert.equal(compact.includes('linearGradient'), false)
})
test('illustrations render token colors with no embedded bitmap or scripts', async () => {
    const Scene = component('SokoCommerceIllustration.vue', {
        commerceIllustrations,
    })
    const html = await renderToString(
        Vue.createSSRApp({
            render: () =>
                Vue.h(Scene, {
                    name: 'warehouse',
                    label: 'Warehouse illustration',
                }),
        })
    )
    assert.match(html, /aria-label="Warehouse illustration"/)
    assert.match(html, /var\(--soko-navy/)
    assert.equal(/<image|<script|foreignObject/.test(html), false)
})

test('six primary storefront categories have distinct relevant scene assignments', () => {
    const source = readFileSync(
        new URL('../../resources/js/Pages/Store/Index.vue', import.meta.url),
        'utf8'
    )
    const body = source.match(
        /const categoryScenes: Record<string, string> = \{([\s\S]*?)\n\}/
    )[1]
    const values = [...body.matchAll(/:\s*'([^']+)'/g)].map((m) => m[1])
    assert.equal(values.length, 6)
    assert.equal(new Set(values).size, 6)
    for (const scene of values) assert.ok(commerceIllustrations[scene])
})
