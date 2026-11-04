# Workspace request pagination

`GET /dashboard?section=rfqs` loads only scoped requests and their permitted quotations. Other sections and buyer account behavior remain unchanged.

Parameters: `q` (maximum 120 characters), `status` (empty/open/awarded/closed), `sort` (newest default or oldest), `page` (1–1,000,000), `per_page` (1–100, default 25), and optional `selected_rfq` (maximum 64 characters). Sorts include an ID tiebreaker. Out-of-range pages clamp to the last page.

Search uses literal Unicode lowercase substring matching across title, category, request number, buyer company/name, destination country and status. It does not search quotation terms or private supplier notes. Buyer email is not part of the request model.

`rfqs` retains its existing shape. `rfqPage` exposes page, per_page, total authorized requests, filtered_total, last_page, from, to and normalized filters. Empty results have from/to zero and last_page one. Other collection totals are not supplied on this section.

`selectedRfq` resolves independently of current filters/page but within the identical read scope. `selectedRfqStatus` is none, available or unavailable. Missing and forbidden records are indistinguishable. `rfqVisibility` is full or supplier. Other sections clear rfqPage and selected metadata.

Full readers see all requests and quotations. Other authorized workspace actors with rfqs.read or rfqs.quote see open requests plus requests on which they previously quoted. They see only their own quotations, including on selected records. The open-or-previously-quoted condition is grouped before applying every filter, ID selection or count. A request that closes without an actor's quotation becomes unavailable on refresh.

A single batched query loads quotations for at most 100 page requests plus one selected request. Capability maps cover those requests and only visible quotations. Quote, award, close and accept rules remain unchanged. Each quotation retains updated_at for the existing versioned award endpoint; reviewed terms must be explicitly refreshed and reconfirmed after conflict. Pagination creates no new mutation route and does not bypass version validation or parent-request write serialization.

Counts, search, selected records, quotation hydration and capability generation use a single read transaction on SQLite. Invalid filters return 400, missing access 403 and query failures 500 rather than false empty results. Separate requests use fresh snapshots and may reflect concurrent changes.

Nonempty search is a bounded-memory O(N) scan of scoped/status-filtered request headers, not an indexed full-text search. Only requested/final page IDs are retained. Empty search uses SQL count/order/offset/limit. A request with many quotations can still have a large child payload; full quotation pagination is a separate future feature. Search does not add accent removal or canonical Unicode normalization.
