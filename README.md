# Soko — Business marketplace

Soko brings wholesale, retail, negotiated offers and digital product listings into one marketplace. The application uses Go/Goravel 1.18, Inertia, Vue 3, Tailwind 4 and SQLite.

The storefront combines searchable product discovery with quantity-based pricing and requirements-led sourcing. Buyer accounts and supplier/admin workspaces use the same navy, amber and light-surface design system. Design decisions are recorded in [DESIGN.md](DESIGN.md), interaction ownership in [UX-CONTRACT.md](UX-CONTRACT.md), and the current implementation scope in [docs/end-to-end-plan.md](docs/end-to-end-plan.md).

## Main journeys

- **Marketplace**: search, category/format/condition/region filters, URL-backed sorting, product details, supplier identity and volume-price brackets.
- **Cart and orders**: a session cart with server pricing, quantity/stock checks and authenticated pending order submission. Pending orders can be cancelled to release their reserved inventory. Orders are recorded in USD; other currencies are indicative display conversions using configured rates.
- **Buyer account**: private order history, submitted offers, sourcing requests and supplier quotations. Request owners can select a quotation or close a request.
- **Supplier workspace**: owned catalog create/edit/delete, relevant sales/order lines, incoming offers and RFQ quotations. Access and ownership are enforced on the server. Edits reject stale versions; listings with order/offer history cannot be deleted.
- **Administration**: marketplace-wide operational records, read-only user directory and persisted storefront text editing.
- **Authentication**: hashed passwords, validated sessions and role-specific destinations.

## Integration limits

**Submitting an order does not collect payment.** New orders remain pending; suppliers and buyers arrange payment, shipping and digital fulfillment separately. No escrow provider, payment gateway, shipment tracking, email delivery, password recovery or protected digital download service is connected. Offer acceptance and quotation selection record marketplace decisions; they do not charge a customer or create a paid entitlement.

The seed data is demonstration content. Ratings, supplier claims and regional payment-rail definitions in legacy records are not evidence of a live verification, reviews or payment system. The active UI avoids advertising those as implemented capabilities. Do not use the demo credentials in a public deployment.

## Local setup

Use **Go 1.26.1** as declared by `go.mod` and **Node 20.19+ or 22.12+** for the installed Vite tooling. SQLite uses a pure-Go driver.

```bash
git clone https://github.com/elcomware/soko.git
cd soko
go mod download
npm ci
cp .env.example .env
go run . artisan key:generate
```

Configure `.env` for local development:

```ini
APP_NAME=Soko
APP_ENV=local
APP_HOST=127.0.0.1
APP_PORT=8000
APP_URL=http://localhost:8000
DB_CONNECTION=sqlite
DB_DATABASE=database/soko.sqlite
```

Run additive migrations. Seed only a new, disposable development database:

```bash
go run . artisan migrate
# Optional on a new demo database only:
go run . artisan db:seed
npm run build
go run .
```

Open **http://localhost:8000**. Rebuild assets and restart the Go process after production asset changes so the template/manifest cache refreshes. Preserve existing databases; do not use `migrate:fresh` to upgrade a working checkout.

## Demonstration accounts

These accounts are created by the development seeders, with password `secret123`:

| Role          | Email                     | Workspace                                             |
| ------------- | ------------------------- | ----------------------------------------------------- |
| Administrator | `admin@soko.global`       | Marketplace operations and storefront content         |
| Supplier      | `supplier.ea@soko.global` | Owned listings, sales, offers and quotations          |
| Supplier      | `merchant.cn@soko.global` | Owned listings, sales, offers and quotations          |
| Buyer         | `buyer@soko.global`       | Private account, orders, offers and sourcing requests |

## Architecture

- `routes/web.go`: public, authenticated and operational routes.
- `app/http/controllers`: Inertia presentation and mutation responses.
- `app/http/requests`: input rules; actions/services enforce domain and authorization invariants.
- `app/actions`, `app/services`: commerce writes and orchestration.
- `app/queries`: reads, scoped presentation and storefront content.
- `app/models`, `database/migrations`: persistence and additive schema changes.
- `resources/js/Pages/Store`: marketplace, product detail and cart.
- `resources/js/Pages/Account`, `Dashboard`, `Auth`: account, operations and authentication.
- `resources/js/components`: shared logo, icons and native accessible dialog.
- `resources/css`: shared tokens, immersive storefront and customer/workspace page styles.

`Home.vue`, `About.vue`, `Contact.vue` and `Feed.vue` are unregistered adapter examples retained in the source tree; they are not active product destinations.

## Checks

```bash
npx vue-tsc --noEmit
npm run build
go test ./tests/commerce ./app/queries ./tests/feature -count=1
git diff --check
```

Commerce regression tests use a dedicated temporary SQLite database, never the existing local marketplace database. Browser verification should cover desktop and narrow screens, invalid form submission, dialog keyboard behavior, cart updates, ownership restrictions and the customer/supplier journey. Payment and delivery must not be inferred from a successful pending-order test.

## License

[MIT](LICENSE).
