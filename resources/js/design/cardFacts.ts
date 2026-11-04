import { hasPromotion } from './merchandising'
export interface CardFooterSettings {
    card_footer_price?: boolean
    card_footer_rating?: boolean
    card_footer_availability?: boolean
    card_footer_mode?: string
    card_footer_interval_seconds?: number
}
export interface CardFact {
    key: 'price' | 'rating' | 'availability'
    label: string
    description: string
}
export const normalizeCardFooter = (settings: CardFooterSettings = {}) => ({
    price: settings.card_footer_price !== false,
    rating: settings.card_footer_rating === true,
    availability: settings.card_footer_availability !== false,
    mode: settings.card_footer_mode === 'rotate' ? 'rotate' : 'together',
    interval:
        Number.isInteger(settings.card_footer_interval_seconds) &&
        Number(settings.card_footer_interval_seconds) >= 5 &&
        Number(settings.card_footer_interval_seconds) <= 60
            ? Number(settings.card_footer_interval_seconds)
            : 8,
})
export const productCardFacts = (
    product: Record<string, any>,
    settings: CardFooterSettings,
    formatPrice: (value: number) => string,
    now = Date.now()
): CardFact[] => {
    const config = normalizeCardFooter(settings)
    const facts: CardFact[] = []
    // Catalog pricing is server-calculated at its explicit quantity. Do not
    // infer a 'From' price from MOQ, or extrapolate a discount to other tiers.
    const pricing = product.pricing
    const amount = hasPromotion(product, now)
        ? pricing?.unit_price
        : (pricing?.original_unit_price ?? pricing?.unit_price)
    if (
        config.price &&
        pricing &&
        typeof amount === 'number' &&
        Number.isFinite(amount) &&
        amount >= 0
    ) {
        const quantity = Number(pricing.quantity) || Number(product.moq) || 1
        const unit = product.type === 'physical' ? 'unit' : 'license'
        facts.push({
            key: 'price',
            label: `${formatPrice(Number(amount))} / ${unit}`,
            description: `Price per ${unit} at ${quantity} ${quantity === 1 ? unit : unit + 's'}; final price is confirmed in cart.`,
        })
    }
    // Legacy rating/reviews_count are demo fields, not proof of customer reviews.
    const review = product.review_summary
    if (
        config.rating &&
        review?.source === 'verified_reviews' &&
        Number.isInteger(review.review_count) &&
        review.review_count > 0 &&
        Number.isFinite(review.average_rating) &&
        review.average_rating >= 1 &&
        review.average_rating <= 5
    ) {
        facts.push({
            key: 'rating',
            label: `${review.average_rating.toFixed(1)} / 5 · ${review.review_count} reviews`,
            description: 'Average of verified customer reviews.',
        })
    }
    if (config.availability) {
        if (
            product.type === 'physical' &&
            Number.isFinite(product.stock) &&
            product.stock >= 0
        )
            facts.push({
                key: 'availability',
                label: product.stock === 0 ? 'Sold out' : 'In stock',
                description:
                    product.stock === 0
                        ? 'Currently unavailable.'
                        : `${product.stock} units listed in stock; availability is confirmed at checkout.`,
            })
        else if (product.type !== 'physical' && product.format === 'digital')
            facts.push({
                key: 'availability',
                label: 'Digital product',
                description:
                    'Digital product; delivery and license details are on its product page.',
            })
    }
    return facts
}
export const cardFactIndex = (now: number, interval: number, count: number) =>
    count > 0 ? Math.floor(now / (interval * 1000)) % count : 0
