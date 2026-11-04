# Workspace order pagination

`GET /dashboard?section=orders` loads only the permitted order page and its items. Buyer account queries and other dashboard sections retain their existing behavior.

Query parameters:
- `q`: literal, Unicode lowercase substring matching against order number, customer name, customer email and status; maximum 120 characters.
- `status`: empty, pending, cancelled, completed or shipped.
- `sort`: newest (default) or oldest, with an ID tiebreaker.
- `page`: 1–1,000,000; pages beyond the last result clamp to the last page.
- `per_page`: 1–100, default 25.
- `selected_order`: optional ID, maximum 64 characters, resolved independently of page/filter within the same authorization scope.

`orders` retains the existing record shape. `orderPage` contains page, per_page, total, filtered_total, last_page, from, to and normalized filters. Total counts orders in the complete authorized scope; filtered_total counts matching orders, not order items. Empty pages use from/to zero and last_page one. No unrelated collection totals are returned.

`selectedOrder` is null or the authorized order with the same item filtering as page rows. `selectedOrderStatus` is none, available or unavailable; missing and forbidden records are indistinguishable. Capabilities include both page rows and the selected record. Client drawers should preserve a reviewed snapshot until explicit refresh, and must not infer cancellation permission from a role or list visibility.

`orderVisibility` is full for orders.read/all or supplier otherwise. Supplier scope contains only orders with items whose products belong to the actor. A single batched child query loads only those items; total_amount and digital/physical flags are recomputed from visible items, and digital license/download values are removed. Customer fields and order-level status remain visible as in the existing supplier workflow. Full readers retain the complete order totals and items. Cancellation requires full visibility, pending state and the existing orders.cancel grant; mutation services recheck authorization and status. No order editing, fulfillment or payment workflow is introduced.

Count, search, page, selected-record and child queries share one read transaction on the supported SQLite deployment. Every failure returns an error, never a successful empty result. Invalid parameters produce 400, denied read access 403, and query failures 500. Unrelated sections reset orderPage to null.

Nonempty search scans only scoped, status-filtered header columns and retains at most the requested and final page IDs. It remains O(N), with bounded pagination memory, not indexed full-text search. Empty search uses SQL count/order/offset/limit. Child hydration is batched for at most 100 orders plus one selected order; an individual order may still contain many items. Separate requests do not share a snapshot, so concurrent writes may shift pages. Search does not normalize accents or canonical Unicode sequences.
