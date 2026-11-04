<template>
    <section class="workspace-panel media-workspace">
        <header class="workspace-panel-heading media-page-header">
            <div>
                <h1>Media</h1>
                <p class="form-note">
                    Upload photos and create mockups. Unpublished images stay
                    private.
                </p>
            </div>
        </header>
        <ProductMediaStudio
            v-model="preview"
            v-model:gallery="gallery"
            standalone
        />
        <details class="media-visibility">
            <summary>Image visibility and retention</summary>
            <p class="form-note">
                Images become public when saved as a product cover or gallery
                image. Removing an image from a listing does not delete the
                library original. To protect existing listings, permanent
                deletion is not available here.
            </p>
        </details>
    </section>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import ProductMediaStudio from '../../components/ProductMediaStudio.vue'
const preview = ref(''),
    gallery = ref<string[]>([])
</script>

<style scoped>
.media-workspace {
    display: grid;
    gap: 1.5rem;
}
.media-page-header {
    margin-bottom: 0;
}
.media-page-header h1 {
    font-family: var(--soko-font-display);
    font-size: clamp(1.5rem, 3vw, 2rem);
    margin: 0.3rem 0 0.6rem;
}
.media-page-header p {
    max-width: 70ch;
}
.media-visibility {
    padding-top: 1rem;
    border-top: 1px solid var(--soko-border);
    font-size: 0.85rem;
    color: var(--soko-muted);
}
.media-visibility summary {
    cursor: pointer;
    font-weight: 600;
}
.media-visibility p {
    margin: 0.75rem 0 0;
    max-width: 90ch;
}
</style>
