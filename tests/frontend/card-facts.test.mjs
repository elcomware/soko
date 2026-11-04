import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import ts from 'typescript'
import { hasPromotion } from '../../resources/js/design/merchandising.ts'
const source = readFileSync(
    new URL('../../resources/js/design/cardFacts.ts', import.meta.url),
    'utf8'
)
const js = ts
    .transpileModule(source, {
        compilerOptions: {
            target: ts.ScriptTarget.ES2022,
            module: ts.ModuleKind.ESNext,
        },
    })
    .outputText.replace(/^import[\s\S]*?;\s*$/gm, '')
    .replace(/export /g, '')
const { normalizeCardFooter, productCardFacts, cardFactIndex } = new Function(
    'hasPromotion',
    js + '; return {normalizeCardFooter,productCardFacts,cardFactIndex}'
)(hasPromotion)
const formatPrice = (n) => `$${n.toFixed(2)}`
const product = {
    type: 'physical',
    stock: 4,
    moq: 2,
    pricing: {
        quantity: 2,
        unit_price: 75,
        original_unit_price: 100,
        promotion_id: 'campaign',
        promotion_ends_at: '2030-01-01T00:00:00Z',
    },
    rating: 5,
    reviews_count: 999,
}
test('admin selections are validated and can hide all footer values', () => {
    assert.equal(
        normalizeCardFooter({ card_footer_interval_seconds: 1 }).interval,
        8
    )
    assert.equal(
        normalizeCardFooter({ card_footer_mode: 'bad' }).mode,
        'together'
    )
    assert.deepEqual(
        productCardFacts(
            product,
            {
                card_footer_price: false,
                card_footer_rating: false,
                card_footer_availability: false,
            },
            formatPrice
        ),
        []
    )
})
test('price uses server quote at its quantity and expired offers revert to original', () => {
    const current = productCardFacts(product, {}, formatPrice, 0)
    assert.equal(current[0].label, '$75.00 / unit')
    assert.match(current[0].description, /at 2 units/)
    const expired = productCardFacts(
        product,
        {},
        formatPrice,
        new Date('2031-01-01').getTime()
    )
    assert.equal(expired[0].label, '$100.00 / unit')
    assert.equal(
        productCardFacts(
            { ...product, pricing: undefined },
            {},
            formatPrice
        ).some((f) => f.key === 'price'),
        false
    )
})
test('legacy demo ratings never render, verified aggregate is required', () => {
    assert.equal(
        productCardFacts(
            product,
            { card_footer_rating: true },
            formatPrice
        ).some((f) => f.key === 'rating'),
        false
    )
    const real = {
        ...product,
        review_summary: {
            source: 'verified_reviews',
            average_rating: 4.2,
            review_count: 7,
        },
    }
    assert.equal(
        productCardFacts(real, { card_footer_rating: true }, formatPrice).find(
            (f) => f.key === 'rating'
        ).label,
        '4.2 / 5 · 7 reviews'
    )
    assert.equal(
        productCardFacts(
            {
                ...real,
                review_summary: { ...real.review_summary, review_count: 0 },
            },
            { card_footer_rating: true },
            formatPrice
        ).some((f) => f.key === 'rating'),
        false
    )
})
test('availability reflects stock without invented urgency and rotation is stable per interval', () => {
    assert.equal(
        productCardFacts(product, {}, formatPrice).find(
            (f) => f.key === 'availability'
        ).label,
        'In stock'
    )
    assert.equal(
        productCardFacts({ ...product, stock: 0 }, {}, formatPrice).find(
            (f) => f.key === 'availability'
        ).label,
        'Sold out'
    )
    assert.equal(cardFactIndex(7999, 8, 3), 0)
    assert.equal(cardFactIndex(8000, 8, 3), 1)
    assert.equal(cardFactIndex(24000, 8, 3), 0)
    assert.equal(cardFactIndex(8000, 8, 0), 0)
})

test('actual footer component pauses rotation and shows all values for reduced-motion control', async () => {
    const { parse, compileScript } = await import('@vue/compiler-sfc')
    const Vue = await import('vue')
    const componentSource = readFileSync(
        new URL(
            '../../resources/js/components/ProductCardFacts.vue',
            import.meta.url
        ),
        'utf8'
    )
    const compiled = compileScript(parse(componentSource).descriptor, {
        id: 'facts-runtime',
    })
    const code = ts
        .transpileModule(compiled.content, {
            compilerOptions: {
                target: ts.ScriptTarget.ES2022,
                module: ts.ModuleKind.ESNext,
            },
        })
        .outputText.replace(/^import[\s\S]*?;\s*$/gm, '')
        .replace('export default', 'return')
    const component = new Function(
        '_defineComponent',
        'computed',
        'ref',
        'normalizeCardFooter',
        'productCardFacts',
        'cardFactIndex',
        'SokoIcon',
        code
    )(
        Vue.defineComponent,
        Vue.computed,
        Vue.ref,
        normalizeCardFooter,
        productCardFacts,
        cardFactIndex,
        {}
    )
    const props = Vue.reactive({
        product,
        settings: { card_footer_mode: 'rotate' },
        formatPrice,
        now: 0,
        paused: false,
    })
    const scope = Vue.effectScope()
    try {
        const state = scope.run(() =>
            component.setup(props, { expose: () => {} })
        )
        assert.equal(state.visibleFacts.value.length, 1)
        assert.equal(state.visibleFacts.value[0].key, 'price')
        state.hold()
        props.now = 8000
        assert.equal(state.visibleFacts.value[0].key, 'price')
        state.release()
        assert.equal(state.visibleFacts.value[0].key, 'availability')
        props.paused = true
        assert.equal(state.visibleFacts.value.length, 2)
        props.settings.card_footer_price = false
        props.settings.card_footer_availability = false
        assert.equal(state.visibleFacts.value.length, 0)
    } finally {
        scope.stop()
    }
})
