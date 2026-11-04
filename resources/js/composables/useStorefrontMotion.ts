import { ref } from 'vue'

// Keep the visitor’s motion choice across storefront navigation, like currency state.
const motionEnabled = ref(true)

export function useStorefrontMotion() {
    return { motionEnabled }
}
