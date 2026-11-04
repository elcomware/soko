// Read-only presentation of server-owned catalog values. Never fabricate sales or reference prices.
export const productValue = (product: Record<string, any>): string => {
    const supplied = String(
        product.value_proposition || product.short_desc || ''
    ).trim()
    if (supplied) return supplied
    const description = String(product.description || '')
        .replace(/<[^>]*>/g, '')
        .trim()
    return (
        description ||
        'Explore this product, available options and supplier details.'
    )
}
export const hasPromotion = (
    product: Record<string, any>,
    now = Date.now()
): boolean => {
    const price = product.pricing
    return Boolean(
        price?.promotion_id &&
        (!price.promotion_ends_at ||
            new Date(price.promotion_ends_at).getTime() > now) &&
        Number.isFinite(Number(price.original_unit_price)) &&
        Number.isFinite(Number(price.unit_price)) &&
        Number(price.unit_price) >= 0 &&
        Number(price.original_unit_price) > Number(price.unit_price)
    )
}
export const savingsLabel = (product: Record<string, any>): string => {
    if (!hasPromotion(product)) return ''
    const percent = Math.floor(
        (1 - product.pricing.unit_price / product.pricing.original_unit_price) *
            100 +
            1e-9
    )
    return percent > 0 ? `Save ${percent}%` : 'Special offer'
}
export const campaignEnd = (date: string): string => {
    const parsed = new Date(date)
    return Number.isNaN(parsed.getTime())
        ? ''
        : new Intl.DateTimeFormat('en', {
              dateStyle: 'medium',
              timeStyle: 'short',
              timeZone: 'UTC',
          }).format(parsed) + ' UTC'
}

// The legacy coffee-sorter seed points to a wheat field, not the listed machine.
// Do not present that known mismatch as product photography.
export const productImage = (product: Record<string, any>): string => {
    const source = String(product.image_url || '')
    return source.includes('photo-1500937386664-56d1dfef3854') ? '' : source
}
export const isIllustrativeImage = (product: Record<string, any>): boolean => {
    try {
        return new URL(productImage(product)).hostname === 'images.unsplash.com'
    } catch {
        return false
    }
}

/** Compact discovery copy only; the product page retains the full merchant description. */
export const productTeaser = (
    product: Record<string, any>,
    limit = 88
): string => {
    const copy = productValue(product).replace(/\s+/g, ' ').trim()
    if (copy.length <= limit) return copy
    const boundary = copy.slice(0, limit - 1).lastIndexOf(' ')
    return (
        copy
            .slice(0, boundary > limit / 2 ? boundary : limit - 1)
            .replace(/[,;:.!?-]+$/, '') + '…'
    )
}
export const carouselIndex = (
    current: number,
    delta: number,
    count: number
): number => (count > 0 ? (((current + delta) % count) + count) % count : 0)
export const swipeDirection = (dx: number, dy: number): number =>
    Math.abs(dx) >= 48 && Math.abs(dx) > Math.abs(dy) * 1.3
        ? dx < 0
            ? 1
            : -1
        : 0

/** Prefer usable imagery, then real offers and merchant picks; preserve input order for ties. */
export const selectHeroProducts = <T extends Record<string, any>>(
    products: T[],
    failedImages: ReadonlySet<string> = new Set(),
    now = Date.now()
): T[] =>
    [...products]
        .sort(
            (a, b) =>
                Number(Boolean(productImage(b)) && !failedImages.has(b.id)) -
                    Number(
                        Boolean(productImage(a)) && !failedImages.has(a.id)
                    ) ||
                Number(hasPromotion(b, now)) - Number(hasPromotion(a, now)) ||
                Number(Boolean(b.featured)) - Number(Boolean(a.featured))
        )
        .slice(0, 3)

export interface ProductBadge {
    label: string
    tone: 'sale' | 'campaign' | 'unavailable' | 'condition' | 'featured'
}
/** Badges reflect catalog state, never invented popularity or scarcity. */
export const productBadges = (
    product: Record<string, any>,
    now = Date.now()
): ProductBadge[] => {
    if (product.type === 'physical' && product.stock === 0)
        return [{ label: 'Sold out', tone: 'unavailable' }]
    const badges: ProductBadge[] = []
    if (hasPromotion(product, now)) {
        const percent = Math.floor(
            (1 -
                product.pricing.unit_price /
                    product.pricing.original_unit_price) *
                100 +
                1e-9
        )
        badges.push({
            label: percent > 0 ? `Save ${percent}%` : 'Special offer',
            tone: 'sale',
        })
        const campaign = String(product.pricing.promotion_name || '').trim()
        if (campaign) badges.push({ label: campaign, tone: 'campaign' })
    } else if (product.featured)
        badges.push({ label: 'Featured', tone: 'featured' })
    if (
        badges.length < 2 &&
        ['open_box', 'refurbished'].includes(product.condition)
    )
        badges.push({
            label:
                product.condition === 'open_box' ? 'Open box' : 'Refurbished',
            tone: 'condition',
        })
    return badges
}

/** Respect manual merchandising order. Automatic placements prefer successfully
 * loadable, non-mismatched imagery while retaining server rank within each group. */
export const orderResolvedHero = <T extends Record<string, any>>(
    products: T[],
    strategy: string,
    failedImages: ReadonlySet<string> = new Set()
): T[] =>
    strategy === 'manual'
        ? products
        : [...products].sort(
              (a, b) =>
                  Number(Boolean(productImage(b)) && !failedImages.has(b.id)) -
                  Number(Boolean(productImage(a)) && !failedImages.has(a.id))
          )
