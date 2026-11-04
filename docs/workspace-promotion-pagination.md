# Campaign pagination and targeting

`GET /dashboard/promotions` accepts q (maximum 120 characters), status (empty/active/scheduled/ended/disabled), sort (newest default or oldest), page (1–1,000,000), per_page (1–100, default 25) and selected_campaign (optional ID, maximum 64 characters). Sort includes an ID tiebreaker; pages beyond the end clamp to the final page.

`campaignPage` supplies page, per_page, total, filtered_total, last_page, from, to and normalized filters. `campaignCounts` contains total/active/scheduled/ended/disabled over the complete authorized scope, independent of current filters. `selectedCampaign` and selectedCampaignStatus (none/available/unavailable) resolve a full scoped campaign independently of pagination, retaining its scope and updated_at save version. Missing and forbidden IDs have the same response.

All queries use one read transaction and one captured serverTime. Disabled takes precedence. Enabled campaigns are scheduled before starts_at, active at/after starts_at and before ends_at, and ended at/after ends_at. `nextStatusChange` is the earliest future enabled start/end in the authorized scope, or empty. Counts describe the server snapshot; clients should refresh at a boundary when safe, or visibly mark stale results while editing/offline. Client clocks must not silently re-filter a page while retaining old server totals.

Search is literal Unicode lowercase matching over the joined campaign name, displayed merchant name and kind. Displayed merchant names come from the same authorized target projection. Scope always precedes counts, search, selected lookup and boundary calculation: promotions.manage/all sees all campaigns, otherwise only the actor's vendor ID.

The save endpoint retains existing validation and compare-and-swap version protection. Optional return_q, return_status, return_sort, return_page and return_per_page restore the list filters through a validated, URL-encoded fixed /dashboard/promotions redirect. Invalid return filters fall back to the list root; no arbitrary return URL is accepted.

## Searchable targeting

The page sends only distinct authorized merchant choices, not the product catalog. Merchant names use a deterministic maximum nonempty catalog display name per vendor, falling back to the vendor ID. The same names drive campaign search. `targetingApi: true` identifies the supported selector API.

`POST /dashboard/promotions/targets` is a read-only JSON query under existing authentication and same-origin protections. Fields are vendor_id, q, page, per_page and scope (a JSON-encoded PromotionScope string). An actor must have promotions.manage for that exact vendor; own-scope users cannot request a foreign or empty/platform vendor. Page bounds match campaign pagination. Scope payloads are limited to 1 MiB and existing scope arrays to 1,000 entries each. Empty unfinished inclusion scopes are allowed for preview, but save validation remains unchanged.

The response contains:
- products: a page of id/title/category/sub_category/vendor_id/vendor_name, ordered by title and ID.
- productPage: page/per_page/total/filtered_total/last_page/from/to.
- categories: complete scoped category/count facets; subcategories: complete scoped category/sub_category/count facets. Blank categories and parentless/blank subcategories are omitted as selectable choices.
- selectedProducts: the union of supplied inclusion/exclusion product IDs, each with status available/unavailable and product or null. Resolution ignores search/page but respects vendor scope. Missing and unauthorized targets return the same unavailable response without disclosing their names.
- preview: total matching products and up to six samples across the full authorized vendor scope, independent of search/page. Exclusions take precedence over all inclusion modes, matching campaign coverage semantics.

Stored selections are never deleted or changed by the query. Clients retain unavailable product/category/subcategory tokens and require an explicit edit to remove them. Query failures must not replace selected chips or coverage with false empty results; stale/pending preview must be visibly distinguished. Mutations still use the original scope and versioned save flow.

One narrow cursor computes products, facets, selected labels and preview in a single SQLite snapshot. Product memory is bounded by page/selected IDs/six samples rather than catalog size. Search and full preview remain O(N); facet and merchant lists scale with distinct values, and selected hydration is bounded by existing array limits. Category or all-products targeting still includes future matching products when the campaign is evaluated, regardless of what was displayed in the picker. Separate requests have separate snapshots. Invalid filters return 400, missing target access 403, and query failures 500.
