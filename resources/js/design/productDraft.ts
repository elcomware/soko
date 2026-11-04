// History recovery intentionally excludes delivery links, upload files and secrets.
const fields = [
    'version',
    'title',
    'category',
    'sub_category',
    'type',
    'format',
    'condition',
    'region',
    'price',
    'wholesale_price',
    'moq',
    'stock',
    'sku',
    'description',
    'image_url',
    'short_desc',
    'featured',
    'digital_license_format',
    'digital_file_size',
    'accepts_offers',
]
export function safeProductDraft(form: Record<string, unknown>) {
    const safe: Record<string, unknown> = {}
    for (const key of fields) {
        const value = form[key]
        if (
            typeof value === 'string' ||
            typeof value === 'boolean' ||
            (typeof value === 'number' && Number.isFinite(value))
        )
            safe[key] = value
    }
    safe.gallery_urls = Array.isArray(form.gallery_urls)
        ? form.gallery_urls.filter((x) => typeof x === 'string')
        : []
    safe.price_tiers = Array.isArray(form.price_tiers)
        ? form.price_tiers
              .filter((t) => t && typeof t === 'object')
              .map((t) => ({
                  min_qty: Number(t.min_qty),
                  max_qty: Number(t.max_qty),
                  price: Number(t.price),
              }))
        : []
    return safe
}
