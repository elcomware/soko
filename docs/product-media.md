# Product media studio

The catalog editor supports an independent 300-character value proposition, cover image, up to ten gallery images and a featured flag. Omitted value/gallery/featured fields preserve stored values for older clients.

## Uploads

Authenticated accounts require `media.create` for their own media library. Uploads accept PNG/JPEG only, check encoded and decoded size, inspect dimensions before decoding, then re-encode to PNG to strip source metadata. Filenames are random server-generated UUIDs. Files live in private `storage/app/product-media/` and must be backed up with the database; they are not committed to Git. The `/media/{id}` route permits the owner (or an appropriately authorized administrator) to preview an unreferenced draft. An image becomes publicly readable only when referenced by a saved product cover/gallery. Catalog saves verify media ownership. Removing all product references restores draft-only access; no-store response headers prevent caching of private drafts. Never upload confidential material intended to remain private after publication.

Default policy: 5 MB input, minimum 600 × 600, maximum 4096 px per side. Administrators with `media.settings.manage` may set limits within fixed safety bounds, and enable/disable remote links or mockup creation. Remote HTTPS images remain optional and are not downloaded server-side; their quality/dimensions are not verified. Disabling remote links affects new links, not links retained on an existing listing.

The actor-only library is available at `/dashboard/media-library` as well as in the product editor. It shows 24 images per page, supports searching mockup title/brand or image ID, and filters photos versus mockups. All 500 permitted assets per actor are reachable. `media.create` delegates can use the standalone library without product permissions; `media.settings.manage` delegates can reach policy administration without media creation or access to another account’s images. Asset deletion and automatic orphan cleanup are not yet implemented; this avoids silently deleting images used by listings. Operations should monitor media storage. The quota is a lightweight limit, not a substitute for deployment-level rate limiting.

## Mockups

Six deterministic templates: software screen, package box, bottle/container, apparel shirt, book/report and generic equipment. Creators can set brand, title and three colors and upload an authorized raster logo. Server-controlled shapes, escaped text and strict hex colors prevent arbitrary SVG markup or scripts. Every output includes “Illustrative mockup · Actual product may differ.” These are customizable illustrations, not AI-generated images or photographs of an actual product. They do not infer specifications or invent unseen product detail.

Creating an image adds it to the media library and selects it in the draft editor. Saving the listing publishes that cover/gallery selection. Changing or cancelling the listing does not delete its library assets.

## Verification

- `go test ./tests/media` covers image decoding, dimension/type/size-policy validation, safe URLs, mockup escaping/template validation and omitted product fields.
- `go test ./tests/commerce -run TestMedia` covers HTTP permission denial, upload persistence, per-actor library isolation, generated mockup persistence and administrator policy enforcement.
- Vue type checking passed at implementation time. Full final browser validation is part of the integrated storefront delivery.

## Customize a saved mockup

Choose **Customize a copy** on a mockup to reload its template, title, brand and colors. The recipe endpoint is restricted to the owning account with `media.create`; public image visibility does not grant recipe access. Original logo raster data is not returned or retained in reusable settings, so add the logo again if needed. Creating a revision saves a new image and leaves the original and its product references untouched. New copies respect current media policy and quota.
