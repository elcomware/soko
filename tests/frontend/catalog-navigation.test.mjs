import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import ts from 'typescript'
test('catalog filters preserve sibling state but reset old page scroll', () => {
    const source = readFileSync(
        new URL(
            '../../resources/js/composables/useCatalogNavigation.ts',
            import.meta.url
        ),
        'utf8'
    )
    const code = ts
        .transpileModule(source, {
            compilerOptions: {
                target: ts.ScriptTarget.ES2022,
                module: ts.ModuleKind.ESNext,
            },
        })
        .outputText.replace(/^import[\s\S]*?;\s*$/gm, '')
        .replace(/export /g, '')
    let sent
    const useCatalogNavigation = new Function(
        'router',
        'usePage',
        code + ';return useCatalogNavigation'
    )({ get: (...args) => (sent = args) }, () => ({
        component: 'Store/Index',
        url: '/?format=wholesale&category=Other',
    }))
    const nav = useCatalogNavigation()
    nav.visitCatalog({ category: 'Electronics & Telecom' })
    assert.equal(sent[1].format, 'wholesale')
    assert.equal(sent[1].category, 'Electronics & Telecom')
    assert.equal(sent[2].preserveScroll, false)
    assert.match(
        nav.catalogHref({ category: 'Agriculture & Processing' }),
        /Agriculture/
    )
    assert.equal(nav.catalogHref({ format: 'all', category: 'all' }), '/')
})
