<template>
    <div class="auth-page site-container">
        <aside class="auth-scene">
            <img
                src="/public/images/sourcing-harbor.jpg"
                alt="Container ships and cranes at an international harbor"
            />
            <div>
                <p class="eyebrow">Trade starts with a connection</p>
                <h2>Back to business.<br />Open to possibility.</h2>
                <p>
                    Sign in to manage purchases, discover sourcing opportunities
                    and keep your business moving.
                </p>
                <Link href="/"
                    >Explore the marketplace <SokoIcon name="arrow"
                /></Link>
            </div>
        </aside>
        <section class="auth-content" aria-labelledby="login-title">
            <p class="eyebrow">Welcome to Soko</p>
            <h1 id="login-title">Sign in</h1>
            <p class="auth-intro">
                Your marketplace, orders and conversations are waiting.
            </p>
            <form
                class="soko-form auth-form"
                novalidate
                @submit.prevent="submit"
            >
                <p
                    v-if="error"
                    id="login-error"
                    role="alert"
                    class="form-error"
                >
                    {{ error }}
                </p>
                <div class="soko-field">
                    <label for="login-email">Email address</label
                    ><input
                        id="login-email"
                        v-model="form.email"
                        name="email"
                        type="email"
                        required
                        autocomplete="username"
                        placeholder="you@company.com"
                        :aria-invalid="invalidFields.includes('email')"
                        :aria-describedby="error ? 'login-error' : undefined"
                    />
                </div>
                <div class="soko-field">
                    <label for="login-password">Password</label>
                    <div class="auth-password">
                        <input
                            id="login-password"
                            v-model="form.password"
                            name="password"
                            :type="passwordVisible ? 'text' : 'password'"
                            required
                            autocomplete="current-password"
                            :aria-invalid="invalidFields.includes('password')"
                            :aria-describedby="
                                error ? 'login-error' : undefined
                            "
                        /><button
                            type="button"
                            :aria-pressed="passwordVisible"
                            aria-controls="login-password"
                            @click="passwordVisible = !passwordVisible"
                        >
                            {{ passwordVisible ? 'Hide' : 'Show' }}
                        </button>
                    </div>
                </div>
                <button
                    type="submit"
                    class="button button-navy auth-submit"
                    :disabled="form.processing"
                >
                    {{ form.processing ? 'Signing in…' : 'Sign in'
                    }}<SokoIcon name="arrow" />
                </button>
            </form>
            <p class="auth-switch">
                New to Soko? <Link href="/register">Create an account</Link>
            </p>
            <details v-if="demoUsers?.length" class="auth-demos">
                <summary>Explore with a demo account</summary>
                <p>
                    These shared demonstration profiles let you try each
                    workspace. Changes may be visible to other demo users.
                </p>
                <div>
                    <button
                        v-for="demo in demoUsers"
                        :key="demo.email"
                        type="button"
                        :disabled="form.processing"
                        @click="quickLogin(demo)"
                    >
                        <span
                            ><strong>{{ demo.role }}</strong
                            ><small>{{ demo.email }}</small></span
                        ><SokoIcon name="arrow" />
                    </button>
                </div>
            </details>
        </section>
    </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { Link, useForm } from '@inertiajs/vue3'
import SokoIcon from '@/components/SokoIcon.vue'
import '../../../css/workspace.css'
const props = defineProps<{
    demoUsers?: {
        role: string
        email: string
        password: string
        company?: string
        badge?: string
    }[]
}>()
const form = useForm({ email: '', password: '' })
const passwordVisible = ref(false)
const error = ref('')
const invalidFields = ref<string[]>([])
const postLogin = () => {
    if (form.processing) return
    error.value = ''
    form.post('/login', {
        preserveScroll: (response) => response.component === 'Auth/Login',
        onSuccess: (response) => {
            const message = (response.props.flash as Record<string, string>)
                ?.error
            if (message) error.value = message
        },
        onError: (errors) => {
            error.value =
                Object.values(errors).join(' ') ||
                'Sign-in failed. Check your details and try again.'
        },
    })
}
const submit = (event: Event) => {
    if (form.processing) return
    invalidFields.value = []
    error.value = ''
    const element = event.target as HTMLFormElement
    const invalid = Array.from(element.querySelectorAll('input')).filter(
        (input) => !input.checkValidity()
    )
    if (invalid.length) {
        invalidFields.value = invalid.map((input) => input.name)
        error.value =
            invalid[0].name === 'email'
                ? 'Enter a valid email address.'
                : 'Enter your password.'
        invalid[0].focus()
        return
    }
    postLogin()
}
const quickLogin = (demo: { email: string; password: string }) => {
    if (form.processing) return
    form.email = demo.email
    form.password = demo.password
    invalidFields.value = []
    postLogin()
}
</script>
