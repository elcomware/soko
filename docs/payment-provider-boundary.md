# Payment provider boundary

Merchant settings are preferences, not activation. This release stores requested provider, methods, country and currency; no credentials are accepted or returned through the settings form. Enabled=true is rejected. The UI explicitly states that payment collection is not activated.

`app/payments.Provider` defines server-side create-session, signature verification and reconciliation operations. `Unconfigured` fails closed. `VerifyStripeSignature` checks exact raw bytes, constant-time HMAC comparison, size limits and a five-minute timestamp tolerance; tests cover tampering, incorrect secret, duplicate timestamps and expired events. This helper is not an activated webhook or complete Stripe adapter.

Before live activation, implement and test the provider's hosted checkout/session API and durable unique `(provider,event_id)` receipts. Verify merchant, order, amount, currency and permitted state transitions inside the same transaction as the event receipt. Treat browser redirects as presentation only, never payment confirmation. Retry ambiguous requests only with stable idempotency keys, then reconcile provider status. Refund and transfer reversal need explicit state models and reconciliation, not a status toggle.

Marketplace settlement is provider/country-specific. Decide platform domicile, seller countries, supported currencies, fees, onboarding and payout responsibility before activating multi-seller payment collection. A merchant choice does not guarantee a supported payment method or provider account. Stripe Connect separate charges/transfers can allocate funds but carry platform fee/dispute responsibility and transfer reversals are separate from refunds. Flutterwave method availability varies by currency and account.

Official references:
- https://docs.stripe.com/connect/separate-charges-and-transfers
- https://docs.stripe.com/webhooks/signature
- https://developer.flutterwave.com/docs/payment-methods
- https://developer.flutterwave.com/docs/split-payments

No live credentials, external transactions, production onboarding, payout or financial account changes were performed.

## Replaceable adapters and regional combinations

`payments.Adapter` extends the execution contract with a stable ID and merchant-specific verified `Capabilities`. Each capability is an explicit country/area/currency/method tuple with collection, refunds and split-settlement flags. Country-wide support uses an empty area; a subdivision-specific capability only matches that subdivision. These are capability declarations from a trusted server-side adapter, not claims accepted from the browser or inferred from provider marketing.

`payments.Route` describes one merchant's selected market and ordered provider allowlist. `payments.Candidates` resolves only matching, collection-capable adapters. Missing adapters, lookup failures, invalid routes and unmatched markets fail closed. Multiple candidates allow a merchant to combine providers across regions or methods without coupling domain services to an SDK. Preference order is not permission to retry payment: an ambiguous result stays pinned to the original provider and idempotency key until reconciliation proves whether money was collected. Never automatically retry a different provider after a timeout.

This release includes a versioned regional route editor under merchant payment settings. Up to 50 country/optional-area/currency/method rules each hold one to three ordered provider preferences. Records use the existing merchant JSON store and compare-and-swap version. `payments.configure/own` permits only that merchant; `/all` administrators can open a target merchant by account ID. Unknown providers, duplicate markets/providers, malformed codes, unsupported contract methods and oversized configurations are rejected. Country-wide routes are fallback only when an exact area route is absent. Legacy single-provider preference fields remain for compatibility.

Routes are draft requests, not verified geographic eligibility: current Stripe contract excludes mobile money, manual transfer accepts only bank transfer, and all country/currency/account combinations still require provider verification. No live adapters are registered. Collection remains disabled and activation requests are rejected, including records with routes. At execution time a future checkout adapter must resolve the persisted rule through `SelectPaymentRoute`, then require `payments.Candidates` using server-verified merchant capability tuples. An unsupported or unavailable capability must never be promoted to live merely because an administrator saved a route. Do not present this contract as activated multi-provider collection. New payment, media and other external-service integrations should follow the same boundary pattern: explicit domain contract, isolated adapter, declared capabilities, safe unavailable state, and contract tests.

## Isolated Stripe test protocol and durable receipt handling

`StripeTestAdapter` now implements Stripe's actual hosted Checkout Session create/retrieve protocol, using a fixed API origin and pinned API version. It requires a test-key prefix, card-only requests, an exact configured merchant and currency, stable caller-provided idempotency keys, and integer minor-unit totals. Requests have a timeout, reject redirects, bound responses, and validate test mode, metadata, amount, currency and the hosted checkout origin. `VerifyEvent` verifies the exact raw body signature before parsing supported Checkout events. A completed session without `payment_status=paid` stays pending. `Reconcile` retrieves the same test session; it never retries another provider.

The adapter is not registered with checkout or merchant routing, exposes no collection capabilities, has no configured credentials, and has no public webhook endpoint. Tests inject an in-memory HTTP transport; no external Stripe request is made. No secrets are collected by the settings UI or committed. It is a protocol implementation tested with fixtures, not proof of successful Stripe account onboarding or live end-to-end payment.

`test_payment_attempts` stores expected test merchant/order/reference/currency/minor-unit amount. `TestPaymentLedger.Handle` verifies the signature then atomically applies a test-only state transition and unique hashed `(provider,event_id)` receipt. Database write locking serializes updates to each attempt across workers. A duplicate with identical canonical data is idempotent; an event ID reused with different data is rejected. Merchant/order/amount/currency mismatches roll back, and delayed pending/failure events cannot undo a successful test payment. Neither raw payloads nor credentials/card data are stored. The ledger does not update production orders, release inventory, deliver digital goods or initiate payout/refund.

Before exposing even a test checkout flow, wire server-authoritative quote/attempt creation, durable pending-attempt recovery after ambiguous session creation, account-specific credentials through an approved secret mechanism, a bounded raw-body webhook endpoint, retry/reconciliation jobs, and deployment/onboarding checks. Before live multi-seller collection, decide merchant of record, direct/destination/separate-charge flow, connected-account ownership, fees, disputes/refunds and transfer reversals, geographic eligibility, payout responsibility and fulfillment state semantics. Existing mixed-vendor orders cannot safely be marked fulfilled from a single provider event without these decisions.

Additional primary references:
- https://docs.stripe.com/api/checkout/sessions/create
- https://docs.stripe.com/checkout/fulfillment
