# Workspace account directory pagination

`GET /dashboard?section=users` is read-only and requires users.read/all. A users.read/own grant does not allow listing or selected-record access. Buyer account routes remain unchanged, and no account mutation is introduced.

Parameters: q (maximum 120 characters), role (empty/admin/merchant/buyer), verification (empty/verified/unverified), sort (newest default or oldest), page (1–1,000,000), per_page (1–100, default 25) and selected_user (optional ID, maximum 64 characters). Date sorting uses ID as a stable tiebreaker; beyond-last pages clamp to the last page.

`userPage` contains page, per_page, total, filtered_total, last_page, from, to and normalized filters. Total is the complete authorized directory count, not page length. `selectedUser` resolves independently of filters/page under the same strict access gate; selectedUserStatus is none, available or unavailable. Missing selected users return null. Permission loss fails the request with 403. Unrelated sections clear userPage and selected metadata.

Search performs literal Unicode lowercase substring matching across name, email, company name, country and role. Exact case-insensitive verified/unverified searches use the verification boolean as before; an explicit verification filter also applies, so conflicting query/filter values return zero. Role options reflect supported roles rather than the current page.

SQL projections explicitly exclude password and internal permission fields. Existing public account-directory fields are retained. Page results, count and selected record use a single read snapshot; storage/query errors return 500, malformed filters 400, rather than false empty results. The directory exposes no per-record mutations. Access-management links remain controlled by their existing separate permissions.

Nonempty text search remains a bounded-memory O(N) scan over the authorized directory; it is not indexed full-text search. Empty or exact verification queries use SQL counts and pagination. Only requested/final page IDs are retained for text search. Separate requests use fresh snapshots and may reflect deletions or other concurrent changes. No accent removal or canonical Unicode normalization is added.
