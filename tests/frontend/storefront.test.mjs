import { test } from 'node:test'
import assert from 'node:assert/strict'
import { iconPaths, iconNames } from '../../resources/js/design/icons.ts'
import { normalizeBrand, brandTokens } from '../../resources/js/design/brand.ts'
import {
    hasPromotion,
    productValue,
    savingsLabel,
    campaignEnd,
} from '../../resources/js/design/merchandising.ts'
test('icon registry contains 100+ distinct safe path geometries', () => {
    assert.ok(iconNames.length >= 100)
    assert.equal(new Set(Object.values(iconPaths)).size, iconNames.length)
    for (const path of Object.values(iconPaths))
        assert.match(path, /^[MmLlHhVvCcSsQqTtAaZz0-9.,\s-]+$/)
    assert.ok(iconPaths.help)
})
test('brand settings are validated and fonts have safe fallbacks', () => {
    const brand = normalizeBrand({
        primary_color: 'url(unsafe)',
        accent_color: '#AA1133',
        display_name: '  Example ',
        font_body: 'unknown',
    })
    assert.equal(brand.primary_color, '#142f48')
    assert.equal(brand.display_name, 'Example')
    assert.equal(brandTokens(brand)['--soko-amber'], '#AA1133')
    assert.match(brandTokens(brand)['--soko-font'], /Manrope/)
})
test('only genuine current promotions render savings', () => {
    const p = {
        pricing: {
            promotion_id: '1',
            original_unit_price: 100,
            unit_price: 80,
            promotion_ends_at: '2030-01-01T00:00:00Z',
        },
    }
    assert.equal(hasPromotion(p, 0), true)
    assert.equal(hasPromotion(p, new Date('2031-01-01').getTime()), false)
    assert.equal(
        hasPromotion({ pricing: { ...p.pricing, unit_price: 100 } }),
        false
    )
    assert.equal(
        hasPromotion({ pricing: { ...p.pricing, unit_price: -10 } }),
        false
    )
    assert.equal(
        hasPromotion({ pricing: { ...p.pricing, promotion_id: '' } }),
        false
    )
    assert.equal(savingsLabel(p), 'Save 20%')
    assert.equal(campaignEnd('invalid'), '')
})
test('merchant value text is preferred and fallback never fabricates benefits', () => {
    assert.equal(
        productValue({
            short_desc: 'Keep your team connected',
            description: 'Technical specs',
        }),
        'Keep your team connected'
    )
    assert.equal(
        productValue({ description: '<p>Supplier description</p>' }),
        'Supplier description'
    )
    assert.match(productValue({}), /Explore this product/)
})

test('known mismatched stock photography is not rendered as a product', async () => {
    const { productImage, isIllustrativeImage } =
        await import('../../resources/js/design/merchandising.ts')
    assert.equal(
        productImage({
            image_url:
                'https://images.unsplash.com/photo-1500937386664-56d1dfef3854?w=800',
        }),
        ''
    )
    assert.equal(
        isIllustrativeImage({
            image_url: 'https://images.unsplash.com/photo-other',
        }),
        true
    )
    assert.equal(
        productImage({ image_url: '/uploads/product.png' }),
        '/uploads/product.png'
    )
})

test('storefront teaser is compact, word bounded and leaves detail copy untouched', async () => {
    const { productTeaser, productValue } =
        await import('../../resources/js/design/merchandising.ts')
    const product = {
        short_desc:
            'Designed to help your team create beautiful product packaging with editable templates, reusable layouts, multiple color choices and useful export options.',
    }
    assert.ok(productTeaser(product).length <= 88)
    assert.match(productTeaser(product), /…$/)
    assert.equal(productValue(product), product.short_desc)
    assert.equal(
        productTeaser({ short_desc: 'Create with confidence.' }),
        'Create with confidence.'
    )
})
test('carousel wraps both ways and ignores vertical scrolling or small gestures', async () => {
    const { carouselIndex, swipeDirection } =
        await import('../../resources/js/design/merchandising.ts')
    assert.equal(carouselIndex(0, -1, 3), 2)
    assert.equal(carouselIndex(2, 1, 3), 0)
    assert.equal(carouselIndex(0, 1, 0), 0)
    assert.equal(swipeDirection(-80, 10), 1)
    assert.equal(swipeDirection(80, 10), -1)
    assert.equal(swipeDirection(20, 0), 0)
    assert.equal(swipeDirection(60, 100), 0)
})

test('hero prioritizes usable images and retains honest fallback when none exist', async () => {
    const { selectHeroProducts } =
        await import('../../resources/js/design/merchandising.ts')
    const missing = {
        id: 'missing',
        featured: true,
        image_url:
            'https://images.unsplash.com/photo-1500937386664-56d1dfef3854',
    }
    const usable = { id: 'usable', image_url: '/uploads/product.png' }
    const broken = {
        id: 'broken',
        featured: true,
        image_url: '/uploads/broken.png',
    }
    assert.equal(selectHeroProducts([missing, usable])[0].id, 'usable')
    assert.equal(
        selectHeroProducts([missing, broken, usable], new Set(['broken']))[0]
            .id,
        'usable'
    )
    assert.equal(selectHeroProducts([missing])[0].id, 'missing')
    assert.equal(selectHeroProducts([]).length, 0)
})

test('card badges show only active real campaigns and truthful catalog states', async () => {
    const { productBadges } =
        await import('../../resources/js/design/merchandising.ts')
    const product = {
        type: 'physical',
        stock: 5,
        pricing: {
            promotion_id: 'promo',
            promotion_name: 'Black Friday',
            original_unit_price: 100,
            unit_price: 75,
            promotion_ends_at: '2030-01-01T00:00:00Z',
        },
    }
    assert.deepEqual(
        productBadges(product, 0).map((b) => b.label),
        ['Save 25%', 'Black Friday']
    )
    assert.deepEqual(
        productBadges(product, new Date('2031-01-01').getTime()),
        []
    )
    assert.deepEqual(
        productBadges({ ...product, stock: 0 }, 0).map((b) => b.label),
        ['Sold out']
    )
    assert.deepEqual(productBadges({ type: 'digital', stock: 0 }), [])
    assert.deepEqual(
        productBadges({ featured: true, condition: 'open_box' }).map(
            (b) => b.label
        ),
        ['Featured', 'Open box']
    )
    assert.deepEqual(
        productBadges({ condition: 'refurbished' }).map((b) => b.label),
        ['Refurbished']
    )
    assert.deepEqual(productBadges({ type: 'physical', stock: 1 }), []) // no invented low-stock urgency
})

test('hero and product-card media fill their viewports', async () => {
    const { readFileSync } = await import('node:fs')
    const css = readFileSync(
        new URL('../../resources/css/shop.css', import.meta.url),
        'utf8'
    )
    const heroRules = [...css.matchAll(/\.shop-hero-visual img\s*\{([^}]+)\}/g)]
    assert.ok(heroRules.length)
    assert.match(heroRules.at(-1)[1], /object-fit:\s*cover/)
    assert.match(heroRules.at(-1)[1], /height:\s*100%/)
    const cardRules = [
        ...css.matchAll(
            /\.product-first-storefront \.product-image-link img\s*\{([^}]+)\}/g
        ),
    ]
    assert.match(cardRules.at(-1)[1], /object-fit:\s*cover/)
    assert.match(cardRules.at(-1)[1], /padding:\s*0/)
    assert.match(css, /\.shop-carousel-controls\s*\{[^}]*position:\s*absolute/)
})

test('hero has one control group and card labels use a dedicated boundary container', async () => {
    const { readFileSync } = await import('node:fs')
    const source = readFileSync(
        new URL('../../resources/js/Pages/Store/Index.vue', import.meta.url),
        'utf8'
    )
    assert.equal(
        (source.match(/class="shop-carousel-controls"/g) || []).length,
        1
    )
    assert.equal(source.includes('class="shop-hero-selector"'), false)
    assert.equal(source.includes('class="product-image-badges"'), false)
    assert.equal(source.includes('class="product-card-labels"'), true)
})

test('cards use single-line names and boundary-positioned labels', async () => {
    const { readFileSync } = await import('node:fs')
    const css = readFileSync(
        new URL('../../resources/css/shop.css', import.meta.url),
        'utf8'
    )
    assert.match(
        css,
        /\.product-card-labels\s*\{[^}]*transform:\s*translateY\(50%\)/
    )
    assert.match(
        css,
        /\.product-first-storefront \.product-body h3\s*\{[^}]*white-space:\s*nowrap/
    )
    const source = readFileSync(
        new URL('../../resources/js/Pages/Store/Index.vue', import.meta.url),
        'utf8'
    )
    assert.ok(source.includes('productTeaser(product, 72)'))
    assert.ok(source.includes('class="product-card-media"'))
})

test('resolved automatic hero prefers images but manual pins retain exact order', async () => {
    const { orderResolvedHero } =
        await import('../../resources/js/design/merchandising.ts')
    const missing = {
        id: 'missing',
        image_url:
            'https://images.unsplash.com/photo-1500937386664-56d1dfef3854',
    }
    const good = { id: 'good', image_url: '/uploads/good.png' }
    const later = { id: 'later', image_url: '/uploads/later.png' }
    assert.deepEqual(
        orderResolvedHero([missing, good, later], 'featured').map((p) => p.id),
        ['good', 'later', 'missing']
    )
    assert.deepEqual(
        orderResolvedHero([missing, good, later], 'manual').map((p) => p.id),
        ['missing', 'good', 'later']
    )
    assert.equal(
        orderResolvedHero([good, later], 'newest', new Set(['good']))[0].id,
        'later'
    )
    assert.equal(orderResolvedHero([missing], 'promotions')[0].id, 'missing')
})
