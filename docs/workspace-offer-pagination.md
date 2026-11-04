# Workspace offer pagination

`GET /dashboard?section=offers` accepts q (maximum 120 characters), status (empty/pending/accepted/declined/countered), sort (newest default or oldest), page (1–1,000,000), per_page (1–100, default 25) and selected_offer (optional ID, maximum 64 characters). Stable date/ID ordering and last-page clamping match the product/order/request sections.

`offerPage` provides page, per_page, total, filtered_total, last_page, from, to and normalized filters. Totals cover authorized offers, not the current page. `selectedOffer` is resolved within that same scope independently of filters/page; selectedOfferStatus is none, available or unavailable, with no distinction between missing and forbidden IDs. `offerVisibility` is full or supplier. Other sections reset pagination metadata.

Full readers see all offers. Supplier readers see only offers for products currently owned by them. Buyer/contact/message fields are unchanged within that authorized scope; the query does not expose other suppliers' offers. Search uses Unicode lowercase literal substring matching across product title, buyer name/email/company and status. It never searches unscoped data.

One bounded product-owner lookup calculates decision capabilities for page and selected records, independently of products.read. An offer can be decided only while pending, with an existing product and offers.decide permission for its current owner. Orphan offers visible to full readers have no decision capability.

Decision mutations recheck these conditions under a transaction: lock offer, then referenced product, verify current ownership and update only a pending offer. Competing decisions cannot both succeed, and product transfer/deletion cannot slip between authorization and update. Offer amount/quantity/message are immutable through current application workflows, so no editable-term version contract is added. Buyer account queries remain unchanged.

All pagination queries share a SQLite read transaction; failures abort rather than returning false empty data or permissive capabilities. Invalid filters return 400, missing access 403 and query failures 500. Nonempty search remains a bounded-memory O(N) header scan, not indexed full-text search; empty search uses SQL count/order/offset/limit. Separate page requests are separate snapshots and can reflect concurrent writes.
