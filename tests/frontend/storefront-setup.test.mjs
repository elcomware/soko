import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { parse, compileScript } from '@vue/compiler-sfc'
import ts from 'typescript'
import * as Vue from 'vue'
import * as merchandising from '../../resources/js/design/merchandising.ts'
import * as categoryRail from '../../resources/js/design/categoryRail.ts'
import * as sectionOrder from '../../resources/js/design/sectionOrder.ts'

// Execute the real compiled SFC setup, including eager Vue watchers. Stub only
// app-boundary services so declaration-order failures cannot hide behind unit tests.
test('storefront setup initializes with populated and empty catalogs', () => {
    const source = readFileSync(
        new URL('../../resources/js/Pages/Store/Index.vue', import.meta.url),
        'utf8'
    )
    const { descriptor } = parse(source)
    const compiled = compileScript(descriptor, {
        id: 'storefront-runtime-test',
    })
    let js = ts.transpileModule(compiled.content, {
        compilerOptions: {
            target: ts.ScriptTarget.ES2022,
            module: ts.ModuleKind.ESNext,
        },
    }).outputText
    js = js
        .replace(/^import[\s\S]*?;\s*$/gm, '')
        .replace('export default', 'return')
    const env = {
        ...Vue,
        ...merchandising,
        ...categoryRail,
        ...sectionOrder,
        _defineComponent: Vue.defineComponent,
        isActivePromotion: merchandising.hasPromotion,
        onMounted: () => {},
        onBeforeUnmount: () => {},
        useCurrency: () => ({
            activeCurrency: Vue.ref({ code: 'USD' }),
            formatPrice: String,
        }),
        useCatalogNavigation: () => ({
            visitCatalog: () => {},
            catalogHref: () => '/',
        }),
        usePage: () => ({ url: '/', props: { auth: { user: null } } }),
        Link: {},
        router: {},
        SokoIcon: {},
        SokoCommerceIllustration: {},
        SokoDialog: {},
        ProductCardFacts: {},
        normalizeCardFooter: () => ({ mode: 'together', interval: 8 }),
    }
    const keys = Object.keys(env).filter(
        (key) => /^[A-Za-z_$][\w$]*$/.test(key) && key !== 'default'
    )
    const component = new Function(...keys, js)(...keys.map((key) => env[key]))
    for (const products of [
        [],
        [
            {
                id: 'a',
                title: 'Product',
                image_url: '/product.png',
                category: 'Other',
                type: 'physical',
            },
        ],
    ]) {
        const scope = Vue.effectScope()
        try {
            const state = scope.run(() =>
                component.setup(
                    {
                        products,
                        categories: [],
                        currencies: [],
                        regions: [],
                        rfqs: [],
                        activeType: 'all',
                        activeFormat: 'wholesale',
                        activeRegion: 'global',
                        activeCategory: 'Other',
                        activeCondition: 'all',
                        searchQuery: '',
                        stats: {},
                    },
                    { expose: () => {} }
                )
            )
            assert.equal(state.heroProducts.value.length, products.length)
            assert.equal(state.missingImages.value.size, 0)
            assert.equal(state.showHero.value, products.length > 0)
            assert.equal(state.sideProducts.value.length, 0)
        } finally {
            scope.stop()
        }
    }
    const sideOnlyScope = Vue.effectScope()
    try {
        const sideOnly = sideOnlyScope.run(() =>
            component.setup(
                {
                    products: [],
                    currencies: [],
                    merchandising: {
                        section_order: ['catalog','category_navigation','showcase','discovery','sourcing'],
                        hero: { enabled: false, products: [] },
                        hero_side: {
                            enabled: true,
                            products: [{ id: 'side-product' }],
                        },
                    },
                },
                { expose: () => {} }
            )
        )
        assert.deepEqual(sideOnly.sectionOrder.value,['catalog','category_navigation','showcase','discovery','sourcing'])
        const rail = {scrollLeft:0,scrollWidth:600,clientWidth:300}
        sideOnly.setCategoryScroller(rail)
        assert.equal(sideOnly.categoryScroller.value.scrollWidth,600)
        sideOnly.setCategoryScroller(null)
        assert.equal(sideOnly.categoryScroller.value,null)
        assert.equal(sideOnly.showHero.value, false)
        assert.equal(sideOnly.sideProducts.value.length, 1)
        assert.match(source, /v-if="showHero \|\| sideProducts.length"/)
        assert.match(source,/v-for="sectionKey in sectionOrder"/)
        assert.match(source,/:ref="setCategoryScroller"/)
    } finally {
        sideOnlyScope.stop()
    }
    const resolvedScope = Vue.effectScope()
    try {
        const resolvedProps = Vue.reactive({
            products: [],
            currencies: [],
            merchandising: {
                hero: {
                    enabled: true,
                    strategy: 'featured',
                    products: [
                        {
                            id: 'missing',
                            image_url:
                                'https://images.unsplash.com/photo-1500937386664-56d1dfef3854',
                        },
                        { id: 'good', image_url: '/uploads/good.png' },
                    ],
                },
            },
        })
        const resolved = resolvedScope.run(() =>
            component.setup(resolvedProps, { expose: () => {} })
        )
        assert.equal(resolved.heroProduct.value.id, 'good')
        resolvedProps.merchandising.hero.strategy = 'manual'
        assert.equal(resolved.heroProducts.value[0].id, 'missing')
    } finally {
        resolvedScope.stop()
    }
    // Calibration: the previous declaration order must fail under this harness.
    const brokenSource = source
        .replace('const missingImages = ref(new Set<string>())', '')
        .replace(
            'const heroProduct = computed(',
            'const missingImages = ref(new Set<string>())\nconst heroProduct = computed('
        )
    const brokenCompiled = compileScript(parse(brokenSource).descriptor, {
        id: 'storefront-tdz-proof',
    })
    const brokenJS = ts
        .transpileModule(brokenCompiled.content, {
            compilerOptions: {
                target: ts.ScriptTarget.ES2022,
                module: ts.ModuleKind.ESNext,
            },
        })
        .outputText.replace(/^import[\s\S]*?;\s*$/gm, '')
        .replace('export default', 'return')
    const brokenComponent = new Function(...keys, brokenJS)(
        ...keys.map((key) => env[key])
    )
    const scope = Vue.effectScope()
    try {
        assert.throws(
            () =>
                scope.run(() =>
                    brokenComponent.setup(
                        { products: [], currencies: [] },
                        { expose: () => {} }
                    )
                ),
            /before initialization/
        )
    } finally {
        scope.stop()
    }
})
