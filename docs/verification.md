# Marketplace and CMS verification — 2026-11-04

## Automated checks

- `go test ./... -count=1`: passed across all packages. The isolated commerce suite covers 24 behaviors and 16 table subcases, including real HTTP cookie-jar login/registration, nested tier creation, ownership, pricing, stock rollback/reservation/cancellation, stale edits, history retention, origin checks and region compatibility.
- `npx vue-tsc --noEmit`: passed for the active customer, authentication and workspace pages.
- `npm run build`: production assets generated successfully (628 modules).
- Prettier checks and `git diff --check`: passed.
- Independent reviews found and closed stale-stock writes, lost supplier history, zero-cent pricing, volume-price regeneration, nested HTTP tier binding and form-feedback defects.
- Premium UI audit: no findings in active product pages; two existing findings remain in the unregistered `Contact.vue` adapter example (validation ownership and textarea styling).

## Browser checks

The browser used the actual compiled Go/Vue application on loopback with a separate SQLite copy and session cookie. Test mutations never touched the local marketplace's business records.

- Catalog sorting commits to the URL and orders actual prices.
- Invalid quantities focus the associated field and expose an inline accessible error.
- Three CNC units use the preserved $44,000 volume price, totaling $132,000, including after a title-only catalog edit.
- Cart changes reprice on the server; navigation to the cart starts at its heading.
- Real login establishes a rotated session. A buyer is redirected from the supplier workspace to their private account.
- CMS rejected non-plain-text content while preserving the administrator's draft, then successfully published valid hero/announcement text to the storefront.
- Catalog editing saves through the real endpoint with a captured version; existing friendly region names select the corresponding canonical market.
- Pending order submission appears in private order history. Cancellation updates status, removes the repeat action and restores reserved stock exactly once.
- Desktop and 375px phone views checked for product, cart, authentication and workspace. No document overflow; product controls remain within their container. Shared navigation and data tables use contained horizontal scrolling where needed.
- Browser console error log was empty during final checks. Native shared dialogs provide focus containment, Escape and inert background; first invalid inputs receive focus.

## Data and limits

Both additive migrations were applied to the backed-up local SQLite database. Existing counts stayed unchanged: 4 users, 6 products, 0 orders, 0 order items, 2 RFQs and 2 offers. SQLite integrity check returned `ok`. Original listing prices and region strings were preserved; compatibility is handled in reads and explicit edits.

Payment collection, escrow, shipment tracking, email/password recovery and protected digital fulfillment are not connected. Pending orders do not assert payment or issue entitlements. Offer acceptance and quotation awards record decisions without creating a paid order. Listings with order/offer history cannot be hard-deleted; set availability to zero instead. Automatic reservation expiry and production-scale pagination/search remain future work.

No deployment or merge is part of this round. The repository retains unregistered adapter example pages but does not advertise them as working product routes.
