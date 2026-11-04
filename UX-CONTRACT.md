# Soko storefront UX contract

## Product context

Business buyers browse, compare, negotiate and source products. The current UI locale is English; the market is global. Accessibility intent for changed surfaces is native semantics, visible focus, labeled controls, keyboard access and responsive reflow. Visual values and migration boundaries live in DESIGN.md.

## Business-context sources

| Domain / scope                                         | Source                                                        | Reviewed   |
| ------------------------------------------------------ | ------------------------------------------------------------- | ---------- |
| Commerce formats and regional scope                    | README.md                                                     | 2026-11-03 |
| Catalog filters and mutation destinations              | routes/web.go; app/http/controllers/store_controller.go       | 2026-11-03 |
| Quantity, price and validation semantics               | app/http/requests/; resources/js/composables/useCurrency.ts   | 2026-11-03 |
| Session and access                                     | app/http/middleware/auth_required.go; routes/web.go           | 2026-11-03 |
| Legal, payment settlement, deletion and privacy policy | No maintained policy found; this visual slice introduces none | 2026-11-03 |

## Canonical UI Map

| Capability     | Canonical owner                                 | Source of truth                    | Allowed variants                                    | Verification                               |
| -------------- | ----------------------------------------------- | ---------------------------------- | --------------------------------------------------- | ------------------------------------------ |
| Select/Listbox | Native select in Layout.vue and Store/Index.vue | This contract                      | native, OS-owned popup geometry accepted            | Open popup and keyboard in browser         |
| Form           | Store/Index.vue validation and postForm helper  | Server form requests               | Offer, checkout, RFQ                                | Empty/invalid submit and pending review    |
| Scrollbar      | resources/css/app.css                           | DESIGN.md                          | Table horizontal overflow, dialog vertical overflow | Browser computed styles                    |
| Toast          | Layout.vue flash notices                        | Existing server flash props        | Success status, error alert                         | Browser live regions                       |
| CRUD           | Inertia router and existing controllers         | routes/web.go; store_controller.go | Preserve established destinations                   | Navigation and invalid-path browser checks |

## Dataset navigation

The existing API returns all matches without pagination. This slice preserves that contract and does not claim to introduce server pagination. Search and format/category/condition/region filters are committed in URL parameters. Changing one preserves the others; reset returns to the unfiltered marketplace. Inertia owns visit cancellation and history. Search runs on explicit form submission, not each keystroke, and has an app-owned clear button. Native IME form behavior is retained; avoid custom Enter handlers.

## Flow ledger

| Operation     | Pending                            | Success                                        | Failure recovery                    | Source                           |
| ------------- | ---------------------------------- | ---------------------------------------------- | ----------------------------------- | -------------------------------- |
| Search/filter | Existing Inertia progress bar      | Updated catalog; URL preserves committed state | Existing catalog remains available  | useCatalogNavigation.ts          |
| Offer         | Disable submit, Sending offer      | Existing server flash; close dialog            | Keep values and show error text     | /make-offer; store_controller.go |
| Order         | Disable submit, Placing order      | Pending order redirects to /account            | Keep values and show error text     | /checkout; store_controller.go   |
| RFQ           | Disable submit, Publishing request | Existing redirect to /#rfq and flash           | Keep values and show error text     | /rfqs; store_controller.go       |
| Dialog cancel | Immediate                          | Close; native focus returns to opener          | Values remain in current page state | SokoDialog.vue                   |

## Validation and overlays

Product forms in the refreshed storefront use novalidate and validity metadata without browser validation bubbles. The first invalid field is focused, marked aria-invalid and associated with the error summary. Server validation remains authoritative. Submits are guarded against duplicates; validation/server failures keep entered data. The UI does not add automatic retries for non-idempotent mutations. A transport interruption retains Inertia's existing error handling and requires deliberate retry.

Native dialog supplies the top layer, focus containment, inert background, Escape and focus restoration. SokoDialog adds an accessible title and body scroll lock. Only one sourcing/offer/order modal is expected at a time. Flash notices remain persistent until dismissed or a new navigation replaces them. Authentication, permission and ownership checks belong to the server. No payment provider or escrow integration is represented by this UI.

## Migration and verification

The current extension covers all active product routes: storefront, product detail, cart, account, login/register and role-scoped workspace. Unregistered Home/About/Contact/Feed adapter samples remain source files, not product destinations. Native alert/confirm is replaced in the workspace by the shared dialog. Existing records are preserved; new persistence uses additive migrations.

Verify changed sources with Prettier, vue-tsc, a production Vite build, the existing Go feature test, and the premium strict audit. Review desktop and narrow storefronts, search/no-results/clear, combined filters, currency changes, native popup keyboard behavior, modal Escape/focus and invalid forms. Compare the shared chrome on login and product detail. Actual checkout and mutation success must be exercised only in an explicitly disposable test workflow; this slice's visual review does not claim payment or deployment readiness.

## Product-led discovery · 2026-11-04

Hero selection is an explicit button, not an auto-advancing carousel. It changes the active real product and its image/value statement. Actual active promotions lead the hero; other featured products are merchant-selected, not claimed best sellers. Category tiles commit filters via existing URL navigation. Search and active filters suppress the discovery hero so results remain primary. Product cards route to details; the compact sourcing utility remains available for missing catalog items. Empty, image-failure and no-result recovery are retained.

Native button/link semantics, visible focus and reduced-motion overrides apply. Product imagery is contained rather than cropping the item out of frame. Offers expire visually using a one-second client clock, and the backend re-prices authoritatively; browser clock precision is not a financial guarantee. No fake scarcity, transaction counts or supplier guarantees are rendered.

## Customer and operations contracts · 2026-11-04

- `/cart` renders server-priced session lines; quantity updates and removal are explicit POST operations. Empty carts explain the next action. Pending requests disable duplicate submission. Authentication is required for checkout, offers and RFQs.
- `/account?section=overview|orders|offers|requests` renders only the actor's activity. Request owners can select a supplier quotation or close an open request through confirmation dialogs that explain the resulting state.
- `/dashboard?section=...` is restricted to supplier/admin roles, with records scoped server-side. Catalog editing exposes meaningful stock zero and optional fields; controls remain open with entered values on server errors. Offer decisions and supplier quotations use actual endpoints.
- Admin-only `/dashboard/content` stores bounded plain text for the public hero and sourcing content. Inertia flash feedback reports persistence outcomes.
- Customer forms share `useMarketplaceForm.ts`; storefront and workspace helpers follow the same first-invalid focus, labelled inline error, duplicate guard and preserved-values behavior. Only the active dialog renders its error summary.
- Product sort is committed in the catalog URL and operates on current server-filtered listings. No server pagination, live exchange rate feed, payment collection, shipping integration or password recovery is claimed.

## Unified operations navigation · 2026-11-04

- `AdminLayout.vue` is the canonical shell for every `Dashboard/*` page, including promotions, media, merchandising and settings. `Layout.vue` remains the buyer/storefront shell. Both client and SSR use the same route-family decision.
- Server `WorkspaceNavigation` is the source for capability-aware grouped links. `media.create` and `media.settings.manage` independently admit appropriate media staff. Navigation does not grant access; controllers/services remain authoritative.
- Per-record `recordCapabilities` represent effective ownership and lifecycle gates for product editing/deletion/media, offers and RFQs. Missing entries deny UI actions. General access to a list never implies permission to modify every row.
- Settings has explicit brand, access and payment destinations. Payment administrators select a real permitted merchant account rather than typing an opaque ID. Saved/rejected form redirects retain their panel and merchant context; optimistic versions advance only on confirmed success.
- Native selects retain operating-system keyboard/popup behavior. Shared dialogs own destructive confirmations. SharedFlashNotifications owns non-blocking success/error notification presentation; inline errors retain correction context. No live payment capability is implied by operational navigation.
