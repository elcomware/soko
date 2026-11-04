<template>
    <section class="icon-library" aria-labelledby="icon-library-title">
        <div class="section-heading">
            <div>
                <h2 id="icon-library-title">Soko icon library</h2>
                <p>
                    {{ iconNames.length }} commerce symbols in two finishes, 12
                    illustrated scenes, 20 country flags and 23 currency
                    symbols.
                </p>
            </div>
        </div>
        <label class="soko-field"
            >Find a symbol<input
                v-model="query"
                type="search"
                placeholder="Search cart, payment, shipping…"
        /></label>
        <p class="icon-library-note">
            Compact line icons stay crisp in small controls. Dimensional symbols
            add body-shaped shading, highlights and depth at larger sizes.
            Original commerce scenes are designed for category tiles and empty
            states. All are illustrations, not photographs. Flags keep national
            colors.
        </p>
        <div class="icon-preview-controls">
            <label
                >Finish<select v-model="finish">
                    <option value="dimensional">Dimensional</option>
                    <option value="line">Compact line</option>
                </select></label
            >
            <label
                >Preview size<select v-model.number="previewSize">
                    <option :value="24">24 px</option>
                    <option :value="48">48 px</option>
                    <option :value="64">64 px</option>
                </select></label
            >
        </div>
        <div class="icon-library-grid">
            <figure v-for="name in filtered" :key="name">
                <SokoIcon :name="name" :size="previewSize" :variant="finish" />
                <figcaption>{{ name }}</figcaption>
            </figure>
        </div>
        <p v-if="!filtered.length" role="status">
            No symbols match your search.
        </p>
        <h3>Commerce illustrations · 96 px</h3>
        <p class="icon-library-note">
            Original layered vectors for shopping, payments, shipping, global
            trade and industry. Use at 64 px or larger; do not substitute them
            for actual product photos.
        </p>
        <div class="icon-illustration-grid">
            <figure v-for="name in illustrationNames" :key="name">
                <SokoCommerceIllustration :name="name" :size="112" />
                <figcaption>{{ name }}</figcaption>
            </figure>
        </div>
        <h3>Locale labels</h3>
        <p class="icon-library-note">
            Explicit language and region labels. Flags alone do not identify a
            language. These previews do not change the storefront language.
        </p>
        <div class="icon-locale-grid">
            <span v-for="locale in locales" :key="locale.code"
                ><bdi :lang="locale.code">{{ locale.native }}</bdi
                ><small>{{ locale.code }} · {{ locale.region }}</small></span
            >
        </div>
        <h3>Countries</h3>
        <div class="icon-library-grid">
            <figure v-for="code in countries" :key="code">
                <SokoFlag :code="code" :size="36" />
                <figcaption>{{ code }}</figcaption>
            </figure>
        </div>
        <h3>Currencies</h3>
        <div class="icon-library-grid">
            <figure v-for="code in currencies" :key="code">
                <SokoCurrencySymbol
                    :code="code"
                    :variant="finish"
                    style="width: 40px; height: 40px"
                />
                <figcaption>{{ code }}</figcaption>
            </figure>
        </div>
    </section>
</template>
<script setup lang="ts">
import { computed, ref } from 'vue'
import { iconNames } from '@/design/icons'
import SokoIcon from './SokoIcon.vue'
import SokoFlag from './SokoFlag.vue'
import SokoCurrencySymbol from './SokoCurrencySymbol.vue'
import SokoCommerceIllustration from './SokoCommerceIllustration.vue'
import { illustrationNames } from '@/design/commerceIllustrations'
const query = ref('')
const finish = ref<'line' | 'dimensional'>('dimensional')
const previewSize = ref(48)
const locales = [
    { code: 'en-GB', native: 'English', region: 'United Kingdom' },
    { code: 'fr-FR', native: 'Français', region: 'France' },
    { code: 'sw-KE', native: 'Kiswahili', region: 'Kenya' },
    { code: 'ar-AE', native: 'العربية', region: 'United Arab Emirates' },
    { code: 'zh-CN', native: '中文', region: 'China' },
    { code: 'de-DE', native: 'Deutsch', region: 'Germany' },
    { code: 'ja-JP', native: '日本語', region: 'Japan' },
    { code: 'hi-IN', native: 'हिन्दी', region: 'India' },
    { code: 'pt-BR', native: 'Português', region: 'Brazil' },
    { code: 'es-ES', native: 'Español', region: 'Spain' },
    { code: 'ko-KR', native: '한국어', region: 'South Korea' },
    { code: 'id-ID', native: 'Bahasa Indonesia', region: 'Indonesia' },
]
const filtered = computed(() =>
    iconNames.filter((name) => name.includes(query.value.trim().toLowerCase()))
)
const countries = [
    'FR',
    'IT',
    'IE',
    'NG',
    'CI',
    'BE',
    'DE',
    'NL',
    'AT',
    'RU',
    'PL',
    'ID',
    'JP',
    'BD',
    'GH',
    'CM',
    'SN',
    'SE',
    'FI',
    'DK',
]
const currencies = [
    'USD',
    'EUR',
    'GBP',
    'JPY',
    'CNY',
    'INR',
    'NGN',
    'KES',
    'ZAR',
    'GHS',
    'KRW',
    'THB',
    'VND',
    'RUB',
    'TRY',
    'PHP',
    'XOF',
    'XAF',
    'AED',
    'BRL',
    'CAD',
    'AUD',
    'CHF',
]
</script>
<style scoped>
.icon-library h3 {
    margin: 24px 0 12px;
    font-weight: 750;
    font-size: 16px;
}
.icon-library .section-heading p,
.icon-library-note {
    font-size: 12px;
    color: var(--soko-muted);
}
.icon-library-note {
    margin: 16px 0;
}
.icon-library-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: 10px;
}
.icon-library-grid figure {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 20px 8px;
    border: 1px solid var(--soko-line);
    border-radius: var(--soko-radius);
    margin: 0;
    min-height: 95px;
    color: var(--soko-navy);
    background: var(--soko-surface);
}
.icon-library-grid figcaption {
    font-size: 10px;
    overflow-wrap: anywhere;
    color: var(--soko-muted);
}
.icon-preview-controls {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    margin: 16px 0;
}
.icon-preview-controls label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
}
.icon-preview-controls select {
    border: 1px solid var(--soko-line);
    padding: 7px 10px;
    border-radius: 5px;
    background: var(--soko-surface);
}
.icon-illustration-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 12px;
}
.icon-illustration-grid figure {
    margin: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 18px;
    background: var(--soko-canvas);
    border: 1px solid var(--soko-line);
    border-radius: 6px;
}
.icon-illustration-grid figcaption {
    font-size: 11px;
    margin-top: 10px;
    color: var(--soko-muted);
}
.icon-locale-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(170px, 1fr));
    gap: 8px;
}
.icon-locale-grid > span {
    padding: 12px;
    border: 1px solid var(--soko-line);
    border-radius: 5px;
}
.icon-locale-grid small {
    display: block;
    margin-top: 4px;
    color: var(--soko-muted);
    font-size: 10px;
}
</style>
