# Workspace product pagination

Only the dashboard products section uses this bounded query path. Other dashboard sections and buyer account queries retain their existing behavior.

## Request and response

`GET /dashboard?section=products` accepts:
- `q`: literal substring search across title, SKU, category and vendor name, at most 120 characters.
- `category`: exact category, at most 128 characters.
- `sort`: `newest` (default), `oldest`, `title_asc` or `title_desc`. Every sort includes an ID tiebreaker.
- `page`: positive integer, at most 1,000,000; requests beyond the last result page clamp to the last page.
- `per_page`: 1–100, default 25.
- `selected_product`: optional product ID, at most 64 characters, resolved independently of current search/category/page but within the same read permission scope.

`products` retains the existing row shape. `productPage` contains page/per_page, total authorized products, filtered_total, last_page, from/to and normalized filters. Empty results have from/to zero. `productCategories` contains distinct nonempty categories from the complete authorized catalog, independent of search. Do not derive total counts or category choices from current-page rows.

`selectedProduct` is a scoped record or null; `selectedProductStatus` is `none`, `available` or `unavailable`. Missing and unauthorized IDs have the same unavailable representation. `recordCapabilities.products` covers returned rows and an available selected record, with the same update/delete/media rules as the existing workspace. A selected record must not silently replace a dirty or reviewed client snapshot.

Invalid filters return 400; missing product-read permission returns 403. Count, category, cursor, row, selection or capability query failures return 500 rather than a successful empty page. Other section responses clear pagination metadata with `productPage: null`.

## Authorization and consistency

Every count, category, search, page and selected-record query begins with products.read own/all scope. Mutations still enforce their own service permissions. Related order/offer history is checked in one bounded query for current page/selection IDs; other workspace collections are not loaded. Counts, search and row hydration run in the same read transaction on the supported SQLite deployment.

The product section does not provide totals for orders, offers, requests or accounts. Its empty unrelated arrays mean not loaded, not a claim that those collections contain zero records. Overview and those sections remain responsible for their own data.

## Search cost and Unicode

SQLite's built-in LOWER is not Unicode-aware. Nonempty search uses a cursor over only authorized, category-filtered ID/title/SKU/category/vendor-name columns and performs Unicode lowercase substring matching in Go. This preserves accented, Greek, Cyrillic and dotted-I case behavior without changing persisted product data. Percent and underscore are literal characters, not SQL wildcards. Search does not add accent removal, transliteration or canonical Unicode normalization.

Search remains O(N) over eligible rows; this is bounded-memory pagination, not an indexed full-text search engine. Only requested-page IDs and the final page's IDs are retained, then at most 100 complete products are hydrated, plus an optional selected record. Empty search uses SQL count/order/offset/limit directly. Very deep offsets and distinct-category counts may still be costly. A maintained normalized search field/index or full-text engine is a future optimization requiring migrations, backfill and coverage of every write path.

The ORM cursor is fully drained, and cursor/scan errors fail the request. A fresh read transaction is used per request; concurrent writes between separate page requests can move records. Stable sorting prevents arbitrary tie order but does not promise a snapshot spanning navigation across multiple requests.
