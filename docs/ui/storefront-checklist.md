# Product-first storefront checklist

Updated: 2026-11-04. This describes the current source and release checks, not a deployment certificate.

## Component ownership

- `Layout.vue`: customer navigation, brand tokens and shared notifications.
- `Store/Index.vue`: filtered catalog, product showcase, category navigation and section rendering.
- `Dashboard/Merchandising.vue`: persisted placement settings and section-order editing.
- `design/brand.ts`: validated merchant appearance and font fallbacks.
- `design/icons.ts`, `SokoIcon.vue` and commerce illustration components: reusable vector artwork.
- `shop.css`: storefront layout and responsive presentation.
- `ProductCardFacts.vue`: configured product facts and optional rotation.

## Implemented behavior

### Product discovery

- Main and side placements use merchandising settings resolved against the already filtered catalog. Side products exclude the full selected hero pool. Empty placements hide rather than fabricate content.
- Featured, newest, genuinely discounted and manual strategies have explicit fallback. Manual pins retain administrator order. Unavailable, missing or out-of-filter pins are omitted.
- Main showcase, category navigation, catalog, digital discovery and sourcing help follow the published section order in the DOM. Hidden sections add no spacer; product catalog cannot be removed by omitting its key.
- Search waits 350 ms after typing, commits immediately on Enter, and clears immediately for empty input. Existing category/format/region filters are retained. New edits and back/forward synchronization cancel stale pending searches.
- Category navigation remains available under filters, supports horizontal swipe/scroll, accurate disabled arrows and active state. ArrowLeft/Right/Home/End operate when the rail itself is focused without intercepting links or Tab. Narrow layouts place controls beneath the full-width rail.

### Media and content hierarchy

- A manually controlled main hero uses previous/count/next, keyboard and swipe navigation without autoplay. Vertical swipes and short gestures do not change slides. Reduced-motion preferences suppress transitions.
- Desktop showcase cards align their allocated height. Smaller screens stack image and copy. Hero and catalog images fill their viewports with cropping; detail pages retain uncropped images for inspection.
- Catalog images dominate concise product names and value excerpts. Full seller descriptions remain available on details. Badges occupy a dedicated image/details boundary row; genuine campaign labels, featured/condition and stock states are not invented.
- Missing or known mismatched imagery receives an honest fallback or disclosure. Illustrations and mockups do not claim to be photographs of the actual product.
- No fabricated seller counts, reviews, popularity, scarcity or discounts are borrowed from visual references.

### Card facts and pricing

- Authorized brand editors can independently enable price, rating and availability; choose simultaneous or timed presentation; and set a 5–60 second interval. Defaults are price and availability enabled, ratings disabled, together, 8 seconds.
- Price uses the server quote at its stated quantity. Expired campaign presentation reverts to the original server price; checkout recalculates authoritative totals.
- Physical availability reflects listed stock. Digital labeling does not promise instant delivery. Ratings stay hidden unless a verified review-summary source exists; seeded rating fields are insufficient.
- Rotation has pause/show-all controls, pauses on hover or focus, and displays all enabled available facts under reduced motion. It does not generate automatic live-region announcements.

### Shared identity and feedback

- The original bag mark, display name, tagline, palette and curated fonts use shared settings. Manrope/system typography is not presented as a custom font family.
- The icon library contains 159 distinct functional vector symbols, plus separate flags/currency symbols and commerce illustrations. Artwork is scalable, accessible and token-driven where appropriate.
- Storefront and administrative layouts share compact fixed snackbars. Success messages use `role=status` and auto-dismiss after 8 seconds unless paused; errors use `role=alert` and persist until dismissal/navigation.
- Notifications do not shift page content or steal focus. Dismiss buttons are labelled, targets are at least 44px, long text wraps within mobile gutters/safe areas, and independent hover/focus state prevents premature dismissal. Cleanup cancels stale timers.

## Implementation boundaries

- Product facts and benefits come from stored seller content, not generated performance claims.
- Merchant settings and promotion displays do not establish payment activation or fulfillment. Orders remain pending until a separately implemented payment workflow is verified.
- The merchandising editor is a constrained section/placement editor, not an arbitrary drag-and-drop page builder, personalized targeting engine or A/B testing system.
- Verified reviews, protected digital downloads and a commissioned font family are not provided.
- Category selection preserves the existing filter model; it is not evidence of a persisted hierarchical taxonomy.

## Release verification checklist

Record the commit, date and result when executing these checks. Implementation descriptions above are not substitutes for testing.

- [ ] `node --test tests/frontend/*.test.mjs`
- [ ] `npx vue-tsc --noEmit`
- [ ] `npm run build`
- [ ] Full Go test suite, including merchandising, promotions and permission boundaries.
- [ ] Desktop plus 375px and 320px rendered checks with long names, populated and empty results.
- [ ] Hero controls, vertical-scroll preservation, keyboard operation and reduced motion.
- [ ] Search debounce, stale cancellation, clearing, filter composition and back/forward navigation.
- [ ] Section reorder publishing, keyboard/reading order and hidden-section behavior.
- [ ] Card facts, expired promotions, unavailable imagery and honest empty states.
- [ ] Snackbar wrap, dismissal, independent hover/focus pauses and route-change cleanup.
- [ ] No horizontal overflow, clipped controls or inaccessible focus after mobile navigation closes.

Automated tests cover catalog setup, filters/search, carousel helpers, icon rendering, product facts, merchandising and notification behavior. Source-level layout assertions cannot establish pixel-level or browser interaction correctness; rendered checks must be recorded separately. Superseded layout experiments and intermediate test counts are intentionally not treated as current release evidence.
