# Workspace lists and actions

Use `AdminTable` for operational lists, within the existing `AdminLayout`.

- `caption` is required and names both the native table and keyboard-focusable horizontal scroll region.
- `header` slot contains a header row with `th scope="col"`. The default slot contains body rows; identify records using `th scope="row"`.
- `empty` and `emptyText` provide the actual empty state. `busy` communicates a pending refresh/save; pages must also disable their mutations while saving.
- Optional `toolbar` and `footer` slots accommodate existing filters and real pagination. Do not invent totals or pagination for a truncated backend result.
- Default rows use 13px text, semantic table markup and shared surface, border, radius and typography tokens. Do not wrap it in legacy workspace-table styling.

Use `AdminActionMenu` for secondary row actions, keeping a record's primary view link visible.

- Required props: contextual `label` (for example, `Actions for Linen shirt`) and `items` array of `{ id, label, disabled?, danger? }`.
- Emits `select(id)`. The page retains the existing business handler, permission check and destructive-action confirmation. A menu is not an authorization boundary.
- Omit unauthorized actions entirely. Empty menus render nothing. `disabled` locks the trigger and selection while a mutation is pending.
- Native popover top layer prevents clipping while preserving the branded DOM ancestry. No body teleport or independent palette is used. Older engines without the popover API retain fixed-position rendering.
- Arrow keys, Home/End and first-character navigation move among enabled actions; Escape restores the trigger, Tab continues native navigation, and outside focus/pointer or outer scrolling closes it. Resize closes it before its anchor can move. Event listeners are removed on close/unmount.

Current integrations: product and account lists, order records, and promotion campaign history. Buyer account/storefront layouts remain separate.

## Responsive editing drawers

Use `AdminDrawer` for workspace CRUD and detailed inspection rather than centered form modals. It remains beneath `AdminLayout` in the DOM so native dialog top-layer rendering keeps the current merchant's brand tokens.

Props: `open`, `title`, optional `description`, `busy`, and `dirty`. It emits `close` only after a clean dismissal or explicit discard confirmation. Compute `dirty` from a stable serialized draft snapshot taken when the editor opens; do not compare to changing server props during a pending save. A successful save may set `open=false` directly. Failed submissions keep the drawer open and existing inline errors intact.

The default and optional footer slots expose `requestClose`. Cancel buttons must call that function, not directly close/reset the parent. Footer submit buttons target the inner form's stable `id` through their `form` attribute. Keep actual permission checks and mutation handlers on the owning page.

On desktop the dialog is a right-hand panel; at 640px and below it becomes a bottom sheet. Header/footer remain visible while body content scrolls, with dynamic viewport units, reduced-motion support and safe-area footer padding. Native `showModal` handles focus containment and inert outside content. Opening focuses the title; closing restores the original connected trigger. Escape and header close share the unsaved-change guard, and a pending save blocks dismissal. Dirty refresh/tab-close invokes the browser's native warning. The component intercepts cancellable Inertia GET visits while dirty, using the same discard prompt; pending mutations block these visits. Full-page forms still need their own navigation guard. Browser Back/Forward uses non-cancellable popstate: the product editor stores a whitelisted draft in Inertia history and recovers it on return, rather than pretending to block history navigation. Recovery excludes transient upload files, arbitrary fields and digital delivery URLs, rechecks record editing permission, and clearly asks users to reselect unfinished uploads or re-enter unpublished delivery links. Successful save or explicit discard clears the stored draft; validation errors preserve it.

## Optional data-table capabilities

Existing `header`/default slots continue unchanged. Each page opts in only to supported features:

- `columns`: `{ key, label, sortable?, width?, minWidth?, maxWidth? }[]`. Omit the header slot to use generated accessible headers. Column widths are explicit pixels; the table scrolls horizontally rather than collapsing content. Default width180, lower bound100 and upper640; configured bounds are safely capped64–1200.
- `sort`: `{key,direction:'asc'|'desc'}`; `@sort` emits that shape. The page maps it to supported server query values and owns request/error state. The component never sorts only the current page while claiming the whole dataset is sorted.
- `resizableColumns`: pointer capture or focus a separator and use Left/Right (10px), Shift+Left/Right (40px), Home/End. Cancellation restores the original width; `column-resize` emits `{key,width}`. A custom header slot owns its own header controls.
- `viewportHeight`: preferred pixel maximum, bounded240–900 and capped to75dvh for small screens. `resizableHeight` adds a pointer/Up/Down/Home/End separator. `update:viewportHeight` emits the preference for optional caller persistence. Neither resizing feature touches shared browser storage.
- `filters` slot is a native collapsible filter region. `toolbar` stays visible. `bulk` is optional and omitted entirely unless supported by the page.
- `AdminPagination` takes authoritative `page,perPage,total,lastPage,from,to,busy`; emits `page` and `perPage`. It does not fetch, invent totals, or claim lazy loading. Use filtered total for a filtered list. Pages must preserve filters, reset page when filters change, and return busy/error state.
- `useAdminSelection(() => eligibleLoadedIds, () => busy)` manages only loaded authorized IDs, prunes selection on page/filter/permission change, and exposes selected/all/mixed/toggle/toggleAll/clear. Pair `AdminSelectionCheckbox` with a contextual label. Never interpret selecting a page as selecting every matching record.
- `AdminBulkActions` takes count/actions/busy and emits action/clear only. It performs no mutation. The page must check each selected record's capability, provide any required confirmation, call a real server operation and handle partial failures. No unsupported bulk mutation is exposed in the Orders integration.
- Action menu entries optionally specify `icon` from SokoIcon; readable text remains present. Icons are decorative and do not replace contextual action labels.

Resizing and sorting controls are disabled while busy. The table retains native semantic rows, focusable scrolling, sticky column headings and shared token typography/surfaces/radii. It is pagination, not virtualized or infinite loading; actual lazy loading should be added only alongside a bounded server contract and measured need.
