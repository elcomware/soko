# Campaign pagination and targeting

`GET /dashboard/promotions` accepts q (maximum 120 characters), status (empty/active/scheduled/ended/disabled), sort (newest default or oldest), page (1–1,000,000), per_page (1–100, default 25) and selected_campaign (optional ID, maximum 64 characters). Sort includes an ID tiebreaker; pages beyond the end clamp to the final page.

`campaignPage` supplies page, per_page, total, filtered_total, last_page, from, to and normalized filters. `campaignCounts` contains total/active/scheduled/ended/disabled over the complete authorized scope, independent of current filters. `selectedCampaign` and selectedCampaignStatus (none/available/unavailable) resolve a full scoped campaign independently of pagination, retaining its scope and updated_at save version. Missing and forbidden IDs have the same response.

All queries use one read transaction and one captured serverTime. Disabled takes precedence. Enabled campaigns are scheduled before starts_at, active at/after starts_at and before ends_at, and ended at/after ends_at. `nextStatusChange` is the earliest future enabled start/end in the authorized scope, or empty. Counts describe the server snapshot; clients should refresh at a boundary when safe, or visibly mark stale results while editing/offline. Client clocks must not silently re-filter a page while retaining old server totals.

Search is literal Unicode lowercase matching over the joined campaign name, displayed merchant name and kind. Displayed merchant names come from the same authorized target projection. Scope always precedes counts, search, selected lookup and boundary calculation: promotions.manage/all sees all campaigns, otherwise only the actor's vendor ID.

The save endpoint retains existing validation and compare-and-swap version protection. Optional return_q, return_status, return_sort, return_page and return_per_page restore the list filters through a validated, URL-encoded fixed /dashboard/promotions redirect. Invalid return filters fall back to the list root; no arbitrary return URL is accepted.

## Remaining targeting cost

The target picker still receives the complete authorized product choices, with only id/title/category/sub_category/vendor_id/vendor_name selected from SQL. This prevents downloading full product models while preserving saved inclusion/exclusion choices, vendor/category facets and preview accuracy. Target options are not silently limited to campaign page rows or truncated. Their memory and response size remain O(number of authorized products).

A separate future target-search endpoint must preserve complete authorized facets, selected-ID resolution and accurate preview counts before replacing this list. Campaign text search itself remains an O(N), bounded-memory cursor over scoped/status-filtered headers, not indexed full-text search. Empty search uses SQL count/order/offset/limit. Separate requests have separate snapshots. Invalid filters return 400, missing campaign access 403, and any query failure 500 rather than empty successful data.
