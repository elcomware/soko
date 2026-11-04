<template>
    <svg
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.6"
        stroke-linecap="round"
        stroke-linejoin="round"
        aria-hidden="true"
        class="soko-icon"
    >
        <path :d="paths[name] || paths.box" />
    </svg>
</template>
<script setup lang="ts">
defineProps<{ name: string }>()
const paths: Record<string, string> = {
    cart: 'M3 4h2l3 12h10l3-9H6 M10 20h.01 M17 20h.01',
    pause: 'M8 5v14 M16 5v14',
    play: 'm8 4 12 8-12 8V4',
    down: 'M12 4v16 m-6-6 6 6 6-6',
    search: 'm21 21-5-5 M19 11a8 8 0 1 1-16 0 8 8 0 0 1 16 0',
    arrow: 'M4 12h16 m-6-6 6 6-6 6',
    chevron: 'm9 5 7 7-7 7',
    globe: 'M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0 M3 12h18 M12 3c5 5 5 13 0 18 M12 3c-5 5-5 13 0 18',
    shield: 'm12 3 8 3v6c0 5-8 9-8 9s-8-4-8-9V6l8-3 m-4 9 3 3 5-5',
    box: 'm12 3 9 5-9 5-9-5 9-5 M3 8v9l9 5 9-5V8 M12 13v9 M7.5 5.5l9 5',
    bolt: 'm13 2-9 12h7l-1 8 10-13h-8l1-7',
    layers: 'm12 3 10 5-10 5L2 8l10-5 M2 12l10 5 10-5 M2 16l10 5 10-5',
    check: 'm5 12 4 4L19 6',
    close: 'm6 6 12 12 M6 18 18 6',
    filter: 'M4 7h16 M4 17h16 M8 4v6 M16 14v6',
    document: 'M14 3H5v18h14V8l-5-5 M14 3v6h5 M8 13h8 M8 17h5',
    factory: 'M3 21V9l6 4V7l7 5V3h4v18H3 M6 17h1 M11 17h1 M16 17h1',
    cloud: 'M7 18h11a4 4 0 0 0 0-8 6 6 0 0 0-11-2 5 5 0 0 0 0 10',
    chip: 'M7 7h10v10H7V7 M10 10h4v4h-4v-4 M9 3v4 M15 3v4 M9 17v4 M15 17v4 M3 9h4 M3 15h4 M17 9h4 M17 15h4',
    leaf: 'M20 3C8 3 3 7 5 14s15 6 15-11 M4 21 15 10',
    star: 'm12 3 2.8 5.7 6.2.9-4.5 4.4 1.1 6.2-5.6-3-5.6 3 1.1-6.2L3 9.6l6.2-.9L12 3',
    user: 'M16 7a4 4 0 1 1-8 0 4 4 0 0 1 8 0 M4 21v-2a8 8 0 0 1 16 0v2',
    logout: 'M9 4H4v16h5 M9 12h12 m-5-5 5 5-5 5',
    menu: 'M4 6h16 M4 12h16 M4 18h16',
}
</script>
