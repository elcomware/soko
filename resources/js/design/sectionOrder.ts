/** Closed layout vocabulary: no executable markup or hidden catalog via configuration. */
export const storefrontSectionKeys = [
    'showcase',
    'category_navigation',
    'catalog',
    'discovery',
    'sourcing',
] as const
export type StorefrontSectionKey = (typeof storefrontSectionKeys)[number]
export const normalizeSectionOrder = (
    value: unknown
): StorefrontSectionKey[] => {
    if (
        !Array.isArray(value) ||
        value.length !== storefrontSectionKeys.length ||
        new Set(value).size !== storefrontSectionKeys.length ||
        !value.every((key) => storefrontSectionKeys.includes(key))
    )
        return [...storefrontSectionKeys]
    return [...value]
}
