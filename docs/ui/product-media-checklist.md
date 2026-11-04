# Product media editor checklist

Sized as a new component within the existing product editor, not a new app shell.
Reference: Dashboard/Index.vue product form. Reuse soko-field, form-row, workspace-fieldset, button, form-note, form-error; use native labelled inputs and semantic fieldsets. No nested forms.
- Insert ProductMediaStudio in the existing image position, retaining image_url and adding gallery_urls.
- Use current semantic CSS variables and existing buttons; component-scoped layout only.
- Keep all controls keyboard reachable with explicit labels and names.
- Display upload limits and server validation errors in role=alert; operation feedback in role=status.
- Upload PNG/JPEG only; require server decoding/re-encoding and dimension validation.
- Offer URL, image upload and six customizable illustrative mockup templates, never implying a photograph or product specification.
- Keep generated illustration label in the image itself.
- Reuse server-owned image library only for the actor; choose cover and add/remove gallery images.
- Show policy controls only with server-granted permission; enforce the permission again on save.
- Keep busy state explicit, disable mutation controls during operations, and preserve selected image on errors.
- Follow one-column layout on small screens; avoid fixed overflowing widths.
