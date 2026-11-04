import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { parse, compileTemplate } from '@vue/compiler-sfc'
for (const page of ['Index', 'Promotions']) test(`${page} admin template compiles in production mode`, () => {
    const filename = `resources/js/Pages/Dashboard/${page}.vue`
    const { descriptor, errors } = parse(readFileSync(filename, 'utf8'), { filename })
    assert.deepEqual(errors, [])
    const result = compileTemplate({ source: descriptor.template.content, filename, id: page, isProd: true })
    assert.deepEqual(result.errors, [])
})
