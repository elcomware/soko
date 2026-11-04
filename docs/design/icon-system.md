# Soko vector system

## Three intentional scales

1. **Interface, 16–24 px.** 159 original semantic glyphs from `resources/js/design/icons.ts`. High-contrast line mode is retained for search, buttons, forms and status controls, where elaborate details become noise.
2. **Dimensional symbols, 28–64 px.** Every glyph supports `SokoIcon`'s `variant="dimensional"`. 100+ object glyphs have body-shaped filled surfaces; utility glyphs remain open shapes. Brand-aware tonal gradients, fine highlights and offset ink shadows create restrained depth without generic decorative tiles. `variant="auto"` is the default: line below28px, dimensional at larger explicit sizes. A one-unit viewBox inset protects shadows from clipping.
3. **Editorial scenes, 64–160 px.** `SokoCommerceIllustration` has12 original multi-layer scenes: storefront, payments, logistics, warehouse, global trade, analytics, digital goods, agriculture, energy, security, customization and returns. Category tiles use this layer. These illustrations must never masquerade as a merchant's actual product photo.

Colors derive from Soko primary/accent/surface/canvas and icon-specific override tokens: `--soko-icon-light`, `--soko-icon-accent`, `--soko-icon-shade`, `--soko-icon-deep-accent`. The default deep accent is computed from merchant primary/accent. SVG gradient IDs are per-component and SSR-safe. No raw SVG HTML injection, remote images or SVG scripts are involved.

## International vocabulary

-20 country flags retain national colors, with country-code fallback for unsupported regions.
-23 currency symbols support compact text or a dimensional coin-style medallion. These are symbols, not depictions of legal tender.
-12 gallery locale samples show language and region explicitly, including bidirectional isolation. A flag alone never represents a language. Gallery previews do not implement language switching.

## Authoring and accessibility

Use `label` for meaningful standalone symbols; omit for decorative icons beside visible text. Keep SVGs out of the tab order. Unknown names resolve to the help symbol. All actions still need real labeled HTML buttons or links. Shading is decorative, never the sole carrier of meaning.

Gallery: `/dashboard/settings`, “Brand icon library,” authorized by brand settings permission. Users can switch finishes and preview sizes. Current flags/currency sets are deliberately bounded, not a claim of universal territory/currency coverage.

## Verification

`tests/frontend/icon-system.test.mjs` validates159 dimensional names,100+ filled surfaces,12 unique layered illustrations, safe path syntax, actual SSR output, unique gradient IDs, accessible labels, fallback behavior and no embedded bitmap/scripts. Typecheck/build passed. A local raster contact sheet was rendered and inspected at48px symbol/112px scene size; browser integration is checked by the parent.

These are original vector illustrations with depth cues, not photorealistic artwork or a new exclusive typeface. Literal photorealism is unsuitable for tiny24px interface symbols; product photography remains separate.
