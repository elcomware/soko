# Soko marketplace implementation plan

Preserve the accepted immersive storefront and existing database. Snapshot source and SQLite before changes; use additive migrations only. Keep all work local.

## Verified gaps

Current product details and workspace use a competing dark theme. Workspace statistics read incorrect property names; several controls have no action. Every authenticated role receives global commerce data. Checkout marks orders completed and issues demonstration licenses without payment. Catalog digital URLs and private competing quotations leak into public props.

## Contracts and page checklist

1. Shared shell: light surfaces, stable catalog search/filter navigation, cart and private account entries, role-specific workspace, POST sign-out, accessible flash feedback.
2. Marketplace: retain manual immersive scenes/reduced motion; remove unsupported ratings/payment copy; server-backed editable main hero and sourcing section; cards link to complete detail or cart/offer flows.
3. Product: breadcrumbs, meaningful media, supplier identity, explicit stock/MOQ/volume prices, specifications, quantity-aware cart action, real offer dialog, related listings with honest labels.
4. Cart: server repricing, editable quantities/removal, clear USD accounting versus display conversion, authenticated pending order submission, unavailable-payment disclosure, actionable empty/error states.
5. Account: only own orders/offers/requests, real status and item history, owner quotation award/close, account details. No invented digital entitlement.
6. Workspace: URL-backed sections, real role-scoped summaries and records, complete catalog create/edit/delete, buyer offer decisions, validated RFQ quotation, admin users read-only and persisted storefront content form.
7. Authentication: consistent split composition, labels/autocomplete, pending protection, inline errors, password visibility, plainly labelled demo accounts; no unsupported recovery links.
8. Backend: validate actual session identities, role and ownership authorization, sanitized public DTOs, positive bounded commerce inputs, USD pending transactional checkout/stock reservation, no payment or delivery fabrication.

## Verification

Run Vue typecheck, production assets build, Go tests against isolated SQLite, migration on the backed-up local database, and browser checks of active pages at desktop and phone widths. Check cart, quantity errors, authentication boundaries, modal keyboard behavior, workspace editing, content persistence and pending order visibility using disposable test fixtures where mutations are needed.

## Explicit limits

Payment collection, escrow, external shipping/tracking, email delivery, password recovery and protected digital fulfillment need real provider configuration and are not claimed as implemented. Existing demo records remain preserved.

## Pattern references

Amazon Brand Stores: modular hero/category/product sections and mobile preview. eBay Stores: clear store navigation/filtering/promotions. Alibaba: separate direct buying from requirements-led sourcing. Seller Hub: actionable orders, inventory and listing workflows. Adapt interaction patterns without copying branding or commercial claims.

## Review corrections

Pending orders can be cancelled by their owner or an administrator. A per-item reservation marker ensures finite stock is released exactly once, including digital listings; historic rows default to no reservation. Product updates require the version shown when the editor opened and reject stale saves. Deletion is blocked for listings referenced by order or offer history. Prices below one cent after rounding are rejected. Session rotation explicitly refreshes the cookie, verified with a real HTTP cookie jar.
