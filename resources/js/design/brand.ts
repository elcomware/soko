export interface BrandSettings {
    display_name: string
    tagline: string
    primary_color: string
    accent_color: string
    font_body: string
    font_display: string
    card_footer_price?: boolean
    card_footer_rating?: boolean
    card_footer_availability?: boolean
    card_footer_mode?: 'together' | 'rotate'
    card_footer_interval_seconds?: number
    version?: number
}
export const defaultBrand: BrandSettings = {
    display_name: 'Soko',
    tagline: 'Global commerce',
    primary_color: '#142f48',
    accent_color: '#f3ac32',
    font_body: 'Manrope',
    font_display: 'Manrope',
    card_footer_price: true,
    card_footer_rating: false,
    card_footer_availability: true,
    card_footer_mode: 'together',
    card_footer_interval_seconds: 8,
}
const fonts: Record<string, string> = {
    Manrope: "'Manrope Variable', Arial, sans-serif",
    manrope: "'Manrope Variable', Arial, sans-serif",
    system: 'system-ui, -apple-system, Segoe UI, sans-serif',
}
export const normalizeBrand = (
    value?: Partial<BrandSettings> | null
): BrandSettings => ({
    ...defaultBrand,
    ...value,
    display_name: value?.display_name?.trim() || defaultBrand.display_name,
    tagline: value?.tagline?.trim() || defaultBrand.tagline,
    primary_color: /^#[\da-f]{6}$/i.test(value?.primary_color || '')
        ? value!.primary_color!
        : defaultBrand.primary_color,
    accent_color: /^#[\da-f]{6}$/i.test(value?.accent_color || '')
        ? value!.accent_color!
        : defaultBrand.accent_color,
})
export const brandTokens = (brand: BrandSettings) => ({
    '--soko-navy': brand.primary_color,
    '--soko-amber': brand.accent_color,
    '--soko-font': fonts[brand.font_body] || fonts.Manrope,
    '--soko-font-display': fonts[brand.font_display] || fonts.Manrope,
    '--soko-mark-color': brand.accent_color,
    '--soko-icon-deep-accent': `color-mix(in srgb, ${brand.accent_color} 78%, ${brand.primary_color})`,
    '--soko-amber-hover': `color-mix(in srgb, ${brand.accent_color} 88%, black)`,
})
