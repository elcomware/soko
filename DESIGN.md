---
version: alpha
name: Soko
description: A connected procurement marketplace, grounded in the practical world of international trade.
colors:
    primary: '#142f48'
    accent: '#f3ac32'
    background: '#f5f7fa'
    surface: '#ffffff'
    text: '#172d43'
    muted: '#5e6d7d'
    border: '#dfe5eb'
    success: '#217254'
    danger: '#ad3544'
    focus: '#3478b3'
    scene-ink: '#0d283b'
    scene-accent: '#e4bb79'
typography:
    sans:
        fontFamily: 'Manrope Variable, Arial, sans-serif'
        lineHeight: '1.6'
rounded:
    DEFAULT: '8px'
    lg: '12px'
spacing:
    page-max: '1320px'
    page-gutter: '40px'
    mobile-gutter: '20px'
    section-gap: '36px'
components:
    button:
        backgroundColor: '{colors.primary}'
        textColor: '{colors.surface}'
    card:
        backgroundColor: '{colors.surface}'
        textColor: '{colors.text}'
    dialog:
        backgroundColor: '{colors.surface}'
        textColor: '{colors.text}'
---

# Soko design system

## Overview

Soko connects business buyers with industrial equipment, digital products and suppliers. README.md describes four commerce formats: wholesale, direct retail, negotiated offers and digital delivery. This is a product marketplace, not a promotional landing page. Make discovery, comparison and sourcing feel familiar and dependable.

The visual reference is a mature, product-led marketplace: navy ink, amber accents, generous product imagery and precise commerce controls. The signature hero presents a real catalog product, supported by two actual product picks, followed by category browsing and benefit-led discovery cards. Promotional prices appear only when the server supplies a current, genuine promotion. Never invent best-seller rankings, reference prices or buyer activity. Corporate storytelling is not the storefront's primary task.

The interface currently uses English for a global audience. Currency selection changes display currency through the existing useCurrency composable; it does not change settlement policy. No Japanese interface or Japan-specific policy is introduced by this slice.

## Colors

Runtime ownership is Model B: resources/css/app.css owns the shared semantic custom properties; resources/css/storefront.css owns the storefront scene values. `colors.scene-ink` maps to `--soko-scene-ink` and `colors.scene-accent` to `--soko-scene-accent`. This file mirrors them and records intent. `colors.primary` maps to `--soko-navy`, accent to `--soko-amber`, background to `--soko-canvas`, surface to `--soko-surface`, text to `--soko-ink`, muted to `--soko-muted`, border to `--soko-line`, and success/danger/focus to identically named `--soko-*` tokens. Components consume these properties directly.

White is the product surface. Cool grey separates the catalog from its cards. Navy carries primary actions and navigation; amber emphasizes the marketplace entrance and the existing Soko identity. Green is reserved for positive outcomes and verification, red for errors. Selected filters include text and aria-pressed semantics. Do not use color alone for status.

The shared chrome and buying tools retain light surfaces; the immersive hero and sourcing studio use deep marine navy, warm gold display text and editorial photographs. Product detail, cart, account, authentication and supplier/admin workspace now share the light trade-desk system. `resources/css/commerce.css` owns customer journey layouts; workspace and auth styles remain scoped to their page family. Forced colors retain system scrollbars and focus outlines.

## Typography

Manrope Variable is bundled from @fontsource-variable/manrope and served locally, with Arial and sans-serif fallbacks. The display uses restrained 750 weight, tight tracking and a compact line-height; controls and body use the same practical family. Hero titles are 52–90px on desktop and 37–56px on phones; editorial sections use 25–38px. Pricing and transactional controls remain quieter. Dense product metadata is secondary; form fields grow to 16px on phones. Numeric prices use tabular numerals. The product's long names wrap; they are not cut off for decoration.

## Layout

The existing Tailwind spacing rhythm uses 4px steps; the refreshed shell follows it. The container is at most 1320px with 40px desktop and 20px mobile gutters. The catalog has a 210px filter rail and three columns. Below 900px filters become a disclosure; below 640px the header search moves to its own full-width row, the hero keeps its full-width scene, sector cards and sourcing requests stack, and product cards use two columns, reducing to one below 360px to preserve readable prices. The hero measures the shared header to fill the remaining viewport without trapping document scrolling. Buyer requests use semantic cards and definition lists; their data and quotation destination are unchanged. Images reserve aspect ratios and the locally served font avoids remote font dependencies.

## Elevation & Depth

Borders and tonal changes establish groups. Product cards gain only a subtle hover shadow. Hero and sector photographs have dark contrast overlays for readable text; the hero product preview uses real catalog data. Shadow and backdrop blur belong to modal dialogs, not static panels. There are no decorative glow effects.

## Shapes

Standard controls and cards use an 8px radius; hero and dialogs use 12px. The Soko mark keeps a small amber container silhouette. Icons are consistent line drawings from SokoIcon.vue, with 1.6px strokes and text labels for actions.

## Components

Layout.vue owns shared navigation, search, region/currency selectors, flash notices and footer. SokoLogo.vue owns the wordmark; SokoIcon.vue owns the small icon vocabulary. Store/Index.vue owns the catalog composition and existing form payloads. useCatalogNavigation.ts merges committed URL filters instead of dropping sibling filter state. useCurrency remains the owner of price conversion.

SokoDialog.vue owns modal naming, native focus containment, Escape dismissal, focus restoration and scroll locking. Platform-native select popups are explicitly accepted; their geometry remains OS-owned. Buttons have hover, pressed, visible focus and disabled states. Pending form labels keep full-width button geometry. Forms use associated labels, a text error summary and first-invalid focus. Scrollbar baseline is global, with forced-color and reduced-motion paths.

Motion is optional: a slow 24-second background drift and 16-second route tracer give the trade scene a sense of movement. The visible Pause motion control stops ambient motion and spatial hover/reveal effects. useStorefrontMotion.ts keeps that choice across SPA navigation, following the existing useCurrency state pattern. Background motion also pauses offscreen. Manual story changes update the image, copy and matching real product; scenes never advance automatically. Intersection observers apply subtle one-time section translation without hiding content. Reduced-motion preferences disable these effects in both CSS and Vue and are watched for changes. All observers and preference listeners are cleaned up on unmount. The reduced-motion override intentionally takes precedence over component transitions. Source photograph: https://images.unsplash.com/photo-1494412519320-aa613dfb7738 (served locally as public/images/sourcing-harbor.jpg).

## Do's and Don'ts

- Do preserve the existing commerce routes, quantities, price tiers and currency calculations.
- Do show real product data and clear action labels.
- Do keep one dominant action per decision area.
- Don't turn README ambitions or seeded demo metrics into unqualified customer guarantees.
- Don't apply this migration to legacy page interiors without verifying their states.
- Don't introduce browser alert/confirm/prompt calls, fake live activity or forced animation.

Additional editorial image sources, served locally: industry https://images.unsplash.com/photo-1581092160607-ee22621dd758; digital https://images.unsplash.com/photo-1558494949-ef010cbdcc31; agriculture https://images.unsplash.com/photo-1500382017468-9049fed747ef. These set the scene and do not represent supplier endorsements.

## End-to-end marketplace extension · 2026-11-04

The active route family now includes storefront, product details, cart, private account, login/register and the supplier/admin workspace. All use the shared light shell. A product detail splits editorial media and a practical buying panel; mobile places the product visual, buying controls and specifications in reading order. Customer accounts use URL-backed activity tabs. The workspace uses a quiet navy navigation rail and white operational panels.

Cart prices are recalculated on the server. Currency selection is an indicative display conversion; orders record USD amounts. Order submission records a pending request, with no payment collection or automatic digital entitlement. Offer acceptance and quotation awards describe their actual recorded state. Do not render seeded ratings as a real reviews system or imply escrow, delivery guarantees or verification beyond actual account data.

Admin content editing changes the first hero heading, accent, description, announcement and sourcing section. Stored plain text is rendered with Vue interpolation. Default content remains available before the additive content migration. See docs/end-to-end-plan.md for implementation and verification scope.


## Product-first identity revision · 2026-11-04

The current storefront supersedes the older immersive-story section below. `resources/css/shop.css` owns its page-specific semantic styles. Shared merchant tokens in `resources/js/design/brand.ts` map validated brand settings into CSS variables on the application shell. The body/display typeface choices are bundled Manrope or system fonts, not an exclusive custom typeface. The bespoke container/S mark is an SVG component; merchant display name and tagline are data, not baked into artwork.

`SokoIcon` reads 159 distinct original SVG paths from `resources/js/design/icons.ts`, inherits currentColor and configurable `--soko-icon-stroke`, and accepts accessible labels and sizes. `SokoFlag` adds 20 national symbols with fixed flag colors; `SokoCurrencySymbol` adds 23 currency labels. These are scalable vector symbols, not photorealistic illustrations. National colors intentionally do not follow brand tokens. `SokoIconGallery` is the canonical searchable admin preview.

Product cards lead with merchant-supplied short descriptions, falling back to actual description or a neutral exploration prompt. They link to product details for prices, quantity and options. Hero offers show original and current server prices, campaign name and explicit UTC expiry; without a valid promotion they present the product without a discount. The browser removes expired offer presentation using its local clock; server checkout remains authoritative. No auto-rotating carousel or forced motion is introduced.

### User feedback refinement · 2026-11-04

The original amber rounded-square lowercase-s identity is restored and refined, replacing the container/hexagon mark above. Merchant color and typography tokens still control it. Discovery copy is now deliberately bounded (88-character card snippets; 120-character hero), with full content reserved for product details. The main hero is a manually operated slide carousel with arrows, keyboard navigation and horizontal swipe. The two side-pick cards are retained; no auto-advance is introduced.

## Unified operational workspace · 2026-11-04

All `Dashboard/*` pages use `AdminLayout.vue`, consistently in client and SSR resolution. Buyers retain `Layout.vue` for the storefront and their account. Merchants and platform staff share the operational shell; effective permission grants, not role labels, determine navigation and actions. The same Soko logo, bundled typography, brand colors, semantic tokens and `--soko-radius` apply throughout. The difference is arrangement: a persistent navy navigation rail with restrained amber edge, quiet white account header and spacious fluid work area. No storefront search, category controls, shopping cart or storefront footer is inserted into admin screens.

`resources/css/admin.css` owns workspace shell geometry only, consuming the shared brand tokens. At narrow widths the navigation becomes an accessible inline disclosure rather than an overlay, so there is no hidden modal focus trap. Page content retains natural document scrolling, while genuine comparison tables may scroll horizontally. All admin pages inherit one `SharedFlashNotifications.vue` instance; recoverable form errors remain inline as well.
