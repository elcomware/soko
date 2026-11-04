import { router, usePage } from '@inertiajs/vue3'

/** Keep committed catalog filters together when the shell or storefront changes one. */
export function useCatalogNavigation() {
    const page = usePage()
    const catalogQuery = (changes: Record<string, string | undefined>) => {
        const current = page.url.split('?')[1]?.split('#')[0] || ''
        const query: Record<string, string> =
            page.component === 'Store/Index'
                ? Object.fromEntries(new URLSearchParams(current))
                : {}
        for (const [key, value] of Object.entries(changes)) {
            if (
                !value ||
                value === 'all' ||
                (key === 'region' && value === 'global')
            )
                delete query[key]
            else query[key] = value
        }
        return query
    }
    const catalogHref = (changes: Record<string, string | undefined>) => {
        const query = new URLSearchParams(catalogQuery(changes)).toString()
        return query ? `/?${query}` : '/'
    }
    const visitCatalog = (changes: Record<string, string | undefined>) => {
        router.get('/', catalogQuery(changes), {
            preserveState: true,
            preserveScroll: false,
        })
    }
    return { visitCatalog, catalogHref }
}
