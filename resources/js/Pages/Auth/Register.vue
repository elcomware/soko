<template>
    <div class="auth-page auth-register site-container">
        <aside class="auth-scene">
            <img
                src="/public/images/sourcing-harbor.jpg"
                alt="An international container harbor"
            />
            <div>
                <p class="eyebrow">Business without boundaries</p>
                <h2>Your next connection<br />starts here.</h2>
                <p>
                    Find products for your business or bring your catalog to
                    buyers around the world.
                </p>
                <ul>
                    <li>
                        <SokoIcon name="search" />Discover physical and digital
                        products
                    </li>
                    <li>
                        <SokoIcon name="document" />Request quotations and
                        compare offers
                    </li>
                    <li>
                        <SokoIcon name="box" />Manage your trade in one
                        workspace
                    </li>
                </ul>
            </div>
        </aside>
        <section class="auth-content" aria-labelledby="register-title">
            <p class="eyebrow">Make room for possibility</p>
            <h1 id="register-title">Create an account</h1>
            <p class="auth-intro">
                Tell us a little about yourself and how you plan to use Soko.
            </p>
            <form
                class="soko-form auth-form"
                novalidate
                @submit.prevent="submit"
            >
                <p
                    v-if="error"
                    id="register-error"
                    role="alert"
                    class="form-error"
                >
                    {{ error }}
                </p>
                <fieldset class="auth-role">
                    <legend>I want to…</legend>
                    <label
                        :class="{
                            'is-selected': form.account_type === 'buyer',
                        }"
                        ><input
                            v-model="form.account_type"
                            type="radio"
                            name="account_type"
                            value="buyer"
                        /><span
                            ><strong>Buy and source</strong
                            ><small
                                >Find products and request quotations</small
                            ></span
                        ></label
                    ><label
                        :class="{
                            'is-selected': form.account_type === 'supplier',
                        }"
                        ><input
                            v-model="form.account_type"
                            type="radio"
                            name="account_type"
                            value="supplier"
                        /><span
                            ><strong>Sell on Soko</strong
                            ><small
                                >Publish products and connect with buyers</small
                            ></span
                        ></label
                    >
                </fieldset>
                <div class="form-row">
                    <div class="soko-field">
                        <label for="register-name">Full name</label
                        ><input
                            id="register-name"
                            v-model="form.name"
                            name="name"
                            required
                            minlength="2"
                            maxlength="255"
                            autocomplete="name"
                            :aria-invalid="invalidFields.includes('name')"
                            :aria-describedby="
                                error ? 'register-error' : undefined
                            "
                        />
                    </div>
                    <div class="soko-field">
                        <label for="register-email">Email address</label
                        ><input
                            id="register-email"
                            v-model="form.email"
                            name="email"
                            required
                            type="email"
                            autocomplete="email"
                            placeholder="you@company.com"
                            :aria-invalid="invalidFields.includes('email')"
                            :aria-describedby="
                                error ? 'register-error' : undefined
                            "
                        />
                    </div>
                </div>
                <div class="form-row">
                    <div class="soko-field">
                        <label for="register-company"
                            >Company name <span>(optional)</span></label
                        ><input
                            id="register-company"
                            v-model="form.company_name"
                            name="company_name"
                            maxlength="255"
                            autocomplete="organization"
                        />
                    </div>
                    <div class="soko-field">
                        <label for="register-country">Country or region</label
                        ><input
                            id="register-country"
                            v-model="form.country"
                            name="country"
                            required
                            maxlength="64"
                            autocomplete="country-name"
                            placeholder="e.g. Kenya"
                            :aria-invalid="invalidFields.includes('country')"
                            :aria-describedby="
                                error ? 'register-error' : undefined
                            "
                        />
                    </div>
                </div>
                <div class="soko-field">
                    <label for="register-password">Password</label>
                    <div class="auth-password">
                        <input
                            id="register-password"
                            v-model="form.password"
                            name="password"
                            :type="passwordVisible ? 'text' : 'password'"
                            required
                            minlength="8"
                            autocomplete="new-password"
                            :aria-invalid="invalidFields.includes('password')"
                            :aria-describedby="
                                error
                                    ? 'password-guidance register-error'
                                    : 'password-guidance'
                            "
                        /><button
                            type="button"
                            :aria-pressed="passwordVisible"
                            aria-controls="register-password"
                            @click="passwordVisible = !passwordVisible"
                        >
                            {{ passwordVisible ? 'Hide' : 'Show' }}
                        </button>
                    </div>
                    <p id="password-guidance" class="form-note">
                        Use at least 8 characters. A longer, unique password
                        helps protect your account.
                    </p>
                </div>
                <button
                    type="submit"
                    class="button button-navy auth-submit"
                    :disabled="form.processing"
                >
                    {{ form.processing ? 'Creating account…' : 'Create account'
                    }}<SokoIcon name="arrow" />
                </button>
                <p class="form-note">
                    {{
                        form.account_type === 'supplier'
                            ? 'Your supplier workspace will be ready after registration.'
                            : 'Your account keeps orders, offers and sourcing requests together.'
                    }}
                </p>
            </form>
            <p class="auth-switch">
                Already have an account? <Link href="/login">Sign in</Link>
            </p>
        </section>
    </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { Link, useForm } from '@inertiajs/vue3'
import SokoIcon from '@/components/SokoIcon.vue'
import '../../../css/workspace.css'
const form = useForm({
    name: '',
    email: '',
    password: '',
    company_name: '',
    country: '',
    account_type: 'buyer',
})
const passwordVisible = ref(false)
const error = ref('')
const invalidFields = ref<string[]>([])
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
        const first = invalid[0]
        const label =
            element
                .querySelector(`label[for="${first.id}"]`)
                ?.textContent?.trim()
                .toLowerCase() || 'this field'
        error.value = first.validity.valueMissing
            ? `Enter your ${label}.`
            : first.validity.tooShort
              ? `Use at least ${first.minLength} characters for your ${label}.`
              : `Enter a valid ${label}.`
        first.focus()
        return
    }
    form.post('/register', {
        preserveScroll: (response) => response.component === 'Auth/Register',
        onSuccess: (response) => {
            const message = (response.props.flash as Record<string, string>)
                ?.error
            if (message) error.value = message
        },
        onError: (errors) => {
            error.value =
                Object.values(errors).join(' ') ||
                'Your account could not be created. Check your details and try again.'
        },
    })
}
</script>
