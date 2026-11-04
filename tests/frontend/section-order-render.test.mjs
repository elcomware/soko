import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { createRequire } from 'node:module'
import { pathToFileURL, fileURLToPath } from 'node:url'
const root = fileURLToPath(new URL('../..', import.meta.url))
const require = createRequire(root + '/package.json')
const { parse, compileScript } = require('@vue/compiler-sfc')
const ts = require('typescript')
const Vue = require('vue')
const SSR = require('@vue/server-renderer')
const merchandising = await import(
    pathToFileURL(root + '/resources/js/design/merchandising.ts')
)
const sections = await import(
    pathToFileURL(root + '/resources/js/design/sectionOrder.ts')
)
const rail = await import(
    pathToFileURL(root + '/resources/js/design/categoryRail.ts')
)
const { descriptor } = parse(
    readFileSync(root + '/resources/js/Pages/Store/Index.vue', 'utf8')
)
const compiled = compileScript(descriptor, {
    id: 'section-order-render',
    inlineTemplate: true,
    templateOptions: { ssr: true },
})
let js = ts.transpileModule(compiled.content, {
    compilerOptions: {
        target: ts.ScriptTarget.ES2022,
        module: ts.ModuleKind.ESNext,
    },
}).outputText
const blank = Vue.defineComponent({ render: () => null })
const env = {
    ...Vue,
    ...merchandising,
    ...sections,
    ...rail,
    isActivePromotion: merchandising.hasPromotion,
    useCurrency: () => ({
        activeCurrency: Vue.ref({ code: 'USD' }),
        formatPrice: String,
    }),
    useCatalogNavigation: () => ({
        visitCatalog: () => {},
        catalogHref: () => '/',
    }),
    usePage: () => ({ url: '/', props: { auth: { user: null } } }),
    Link: Vue.defineComponent({
        setup:
            (_p, { slots }) =>
            () =>
                Vue.h('a', slots.default?.()),
    }),
    router: {},
    SokoIcon: blank,
    SokoCommerceIllustration: blank,
    SokoDialog: blank,
    ProductCardFacts: blank,
    normalizeCardFooter: () => ({ mode: 'together', interval: 8 }),
}
js = js
    .replace(
        /^import\s*\{([^{}]*?)\}\s*from\s*["'](vue|vue\/server-renderer)["'];?\s*$/gm,
        (_all, names, module) => {
            for (const name of names.split(',')) {
                const [original, alias] = name.trim().split(/\s+as\s+/)
                env[alias || original] = (module === 'vue' ? Vue : SSR)[
                    original
                ]
            }
            return ''
        }
    )
    .replace(/^import[\s\S]*?;\s*$/gm, '')
    .replace('export default', 'return')
const keys = Object.keys(env).filter(
    (k) => /^[A-Za-z_$][\w$]*$/.test(k) && k !== 'default'
)
const component = new Function(...keys, js)(...keys.map((k) => env[k]))
async function rendered(order, { hidden = false } = {}) {
    const product = {
        id: 'p',
        title: 'Product',
        image_url: '/product.png',
        category: 'Other',
        type: 'physical',
        price: 10,
        stock: 5,
        moq: 1,
    }
    const props = {
        products: [product],
        categories: [],
        currencies: [],
        regions: [],
        rfqs: [],
        activeType: 'all',
        activeFormat: 'all',
        activeRegion: 'global',
        activeCategory: 'all',
        activeCondition: 'all',
        searchQuery: '',
        stats: {},
        merchandising: {
            section_order: order,
            hero: { enabled: !hidden, products: [product] },
            hero_side: { enabled: false, products: [] },
            discovery: {
                enabled: !hidden,
                products: [{ ...product, id: 'digital', type: 'digital' }],
            },
            category_navigation: { enabled: !hidden, title: 'Categories' },
        },
    }
    const html = await SSR.renderToString(Vue.createSSRApp(component, props))
    return [...html.matchAll(/data-storefront-section="([^"]+)"/g)].map(
        (m) => m[1]
    )
}
test('real compiled storefront renders default section sequence', async () =>
    assert.deepEqual(await rendered(undefined), [
        ...sections.storefrontSectionKeys,
    ]))
test('real compiled storefront renders published custom sequence', async () => {
    const order = [
        'sourcing',
        'catalog',
        'discovery',
        'category_navigation',
        'showcase',
    ]
    assert.deepEqual(await rendered(order), order)
})
test('hidden sections leave no rendered section wrappers', async () =>
    assert.deepEqual(
        await rendered(
            [
                'sourcing',
                'discovery',
                'showcase',
                'catalog',
                'category_navigation',
            ],
            { hidden: true }
        ),
        ['sourcing', 'catalog']
    ))
