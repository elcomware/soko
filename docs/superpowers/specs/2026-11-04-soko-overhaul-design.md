# Soko storefront and commerce management

Updated: 2026-11-04
Status: Product requirements and implementation boundaries for the current source tree. Implementation does not imply deployment or live payment availability.

## Product intent

Soko is a product-first marketplace: shoppers discover relevant products, understand seller-provided value, inspect details, and place an order or request an offer. The visual system retains navy, warm amber and light surfaces, with configurable brand identity. Administration must operate on persisted records with server-enforced permissions rather than demonstrate controls that have no backend.

The application remains Go/Goravel, Vue 3, Inertia, Tailwind and SQLite. Schema upgrades are additive; existing marketplace data must be preserved. Pending orders are not paid orders. Payment collection, shipping and protected digital fulfillment are separate capabilities.

## Implemented in the current source

### Shopping and merchandising

- Product-led main showcase and side placements, category navigation, catalog, digital discovery and compact sourcing help. Products resolve against the current filtered catalog; unavailable or excluded manual selections cannot bypass filters.
- A dedicated merchandising editor controls announcement text, placement enablement, titles, category restrictions, selection strategy, ordered product pins, limits and fallback. Supported strategies are featured, newest, discounted and manual.
- Five sections can be reordered with accessible Move up/Move down controls. The published order determines DOM, reading and keyboard order. A complete validated section list prevents arbitrary component injection or accidental catalog removal.
- Search commits after a short debounce; Enter commits immediately and clearing removes the search query while retaining other filters. Navigation cancels stale pending searches.
- The hero is manually controlled, with keyboard and swipe support and no autoplay. Product cards prioritize imagery and concise value text; detail pages retain the full description and inspectable gallery.
- Configurable card facts include price and availability. Ratings require a verified review source, which is not currently provided; legacy demo ratings are not represented as genuine reviews.

See [storefront merchandising controls](../../storefront-merchandising.md) and the [storefront checklist](../../ui/storefront-checklist.md).

### Promotions and pricing

- Optional enabled campaigns support percentage or fixed USD reductions, explicit start/end instants, merchant scope, category/subcategory/product inclusion and exclusion. Exclusions win and end time is exclusive.
- The pricing function applies the best eligible single campaign to the applicable quantity tier. Discounts do not stack. Calculations use cents and half-up rounding; existing persisted monetary fields remain floating-point representations.
- Expiration is checked during server pricing rather than depending on a scheduler. Detail quotes, cart and checkout use server pricing. Order lines preserve original unit price, discount and selected campaign identity.
- Quantities, stock and price are validated on submission. An order remains pending and creates no paid digital entitlement.

### Cart and orders

- Quantity changes save automatically with debounce and explicit saving/error behavior. Checkout waits for pending client mutations. Server pricing and inventory validation remain authoritative.
- Buyer-bound checkout tokens prevent repeated reservations on retries. Cancelling an eligible pending order releases reserved stock once.
- Same-session request serialization protects cart updates in a single application process. Multi-worker operation needs shared locking or transactional cart storage; see [session concurrency](../../session-concurrency.md).

### Administration and access

- A dedicated administrative shell separates operational navigation from the storefront. Products, orders, offers/RFQs, promotions, merchandising, media and settings expose actions based on effective permissions.
- Grants use named actions and `own` or `all` scope. Role defaults are templates; an explicit profile replaces them. Queries and mutation services check ownership and current grants, including after revocation.
- Dashboard record capabilities also reflect record state, ownership and deletion-history constraints. Frontend visibility never replaces backend authorization.
- Access changes are audited and constrain delegation and removal of the last administrator. Versioned settings reject stale writes.
- Branding controls display name, tagline, validated primary/accent colors, curated body/display font choices and card facts. Shared tokens carry those choices into the active interfaces.

### Identity, icons and media

- An original bag-shaped Soko mark and shared wordmark use merchant brand settings. Typography uses licensed Manrope/system choices; this is not a commissioned custom font family.
- A typed SVG registry provides 159 distinct functional symbols, with separate flag and currency collections and commerce illustrations. Symbols support accessible labels and predictable fallbacks. They are scalable vector artwork, not photographs.
- Product text, cover image, ordered gallery and featured state are independently editable. Uploads are validated PNG/JPEG, re-encoded to strip source metadata and stored with actor ownership.
- Default media policy is 5 MB input, minimum 600 × 600 and maximum 4096 pixels per side, within administrator-configurable safety bounds. Remote HTTPS images are references, not server-fetched imports.
- Six deterministic mockup templates accept controlled text/colors and raster logos. Outputs are explicitly illustrative, not AI-generated product photography. Draft assets remain private until referenced by a published product.

See [product media](../../product-media.md) for storage, privacy and operational limits.

### Payments: preparation only

- Merchant settings persist provider/method/country/currency preferences and versioned regional routing rules with ordered provider choices. Preferences do not activate collection; enabling live collection is rejected.
- Provider-neutral contracts and capability matching fail closed. A Stripe test protocol adapter creates/retrieves Checkout Sessions and verifies signed events under test-only constraints.
- An isolated test ledger validates merchant/order/reference/currency/amount and deduplicates event receipts. It does not update production orders or trigger fulfillment.
- Adapter tests use an in-memory transport. There are no configured payment credentials, registered live checkout adapters or public payment webhook route. Successful fixture tests are not provider-sandbox or live verification.

See [payment provider boundary](../../payment-provider-boundary.md).

## Requirements retained for future work

The following requirements remain part of the product direction, not claims of delivered functionality:

- **Catalog depth:** persisted category hierarchy with descendant semantics and cycle prevention; structured benefit/specification fields; explicitly modeled variants with their own price and stock. Do not imply that these exist through display copy alone.
- **Promotion lifecycle:** richer draft/scheduled/paused/expired/archive views, explicit affected-product and overlap previews, optional business-approved priority rules, complete minor-unit persistence migration, and an explicit changed-price review/confirmation flow. Stacking and applying discounts to negotiated prices require separate defined rules.
- **Publishing:** complete brand/content draft preview, publish history and rollback; broader administrative audit coverage, search/filter/pagination workflows and explicit resource-scoped grants beyond own/all where needed.
- **Branding:** centralized favicon management, additional curated font roles and contrast validation. A complete original font family is a separate design and licensing decision.
- **Media:** alternative raster formats, original/derivative storage, focal-point/crop tools, richer alt-text and usage management, artwork placement/scaling, undo/reset and editable recipe round trips. Asset deletion/orphan cleanup and deployment-level abuse controls are still needed. A future URL importer must defend against SSRF, redirects to private hosts and oversized responses.
- **Live payments:** secure account-specific credential configuration, server-authoritative attempt creation, hosted checkout integration, bounded signed webhooks, durable recovery/reconciliation and account-specific sandbox verification. Refunds, disputes and transfer reversals need separate state models and idempotency.
- **Multi-seller funds flow:** choose direct merchant payments versus platform collection/payout, merchant of record, fees, onboarding, geographic eligibility, refunds/disputes and payout responsibility. Multiple providers must not be presented as one atomic payment. Never retry another provider after an ambiguous result before reconciliation.
- **Fulfillment:** shipping/tracking, verified reviews and protected digital delivery remain separate integrations.

No future feature should fabricate availability, sales popularity, scarcity, reviews, certifications, discounts or fulfillment guarantees.

## Acceptance and release checks

For each release, record actual results rather than treating the following matrix as evidence that checks passed:

1. Campaign inclusion/exclusion, time boundaries, overlaps, tiers, rounding and immutable order prices.
2. Rapid cart edits, removal during saving, validation/retry, checkout retries, inventory restoration and navigation recovery.
3. Permission action/scope/owner combinations, direct unauthorized requests, revocation, delegation and administrator retention.
4. Settings conflicts, brand-token propagation, keyboard operation, accessible labels and contrast.
5. Media type/size/dimension checks, malformed content, draft privacy, publication and ownership.
6. Payment signature/idempotency/mismatch tests and explicit distinctions between fixture, sandbox and live verification.
7. Typecheck, frontend tests, full Go suite, production build and diff checks on the final source.
8. Rendered desktop, 375px and 320px checks; keyboard, reduced-motion, empty/error states and browser history. Static or source-level checks do not establish rendered correctness.

Preserve the existing database, back up uploaded media with its records, and disclose any failed or unrun verification before release.
