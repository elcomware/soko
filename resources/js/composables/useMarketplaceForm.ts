import { ref } from 'vue'
import { router, usePage } from '@inertiajs/vue3'
export function useMarketplaceForm() {
    const busy = ref(false),
        error = ref(''),
        invalid = ref('')
    function validate(event: Event) {
        error.value = ''
        invalid.value = ''
        const field = Array.from(
            (event.currentTarget as HTMLFormElement).querySelectorAll<
                HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
            >('input,select,textarea')
        ).find((f) => !f.validity.valid)
        if (!field) return true
        invalid.value = field.id
        error.value = `Check ${field.labels?.[0]?.textContent?.trim().toLowerCase() || 'the highlighted field'} and enter a valid value.`
        field.focus()
        return false
    }
    const page = usePage()
    function post(
        url: string,
        data: Record<string, any>,
        success?: () => void
    ) {
        if (busy.value) return
        busy.value = true
        error.value = ''
        const origin = page.component
        router.post(url, data, {
            preserveScroll: (response) => response.component === origin,
            onSuccess: (page) => {
                const message = (page.props.flash as any)?.error
                if (message) error.value = String(message)
                else success?.()
            },
            onError: (errors) => {
                error.value =
                    Object.values(errors).join(' ') ||
                    'Please check your details.'
            },
            onFinish: () => {
                busy.value = false
            },
        })
    }
    const fieldState = (id: string) => ({
        'aria-invalid': invalid.value === id || undefined,
        'aria-describedby':
            invalid.value === id ? 'market-form-error' : undefined,
    })
    return { busy, error, invalid, validate, post, fieldState }
}
