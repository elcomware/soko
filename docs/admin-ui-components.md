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
