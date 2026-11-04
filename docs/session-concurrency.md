# Session and cart concurrency

The current single-process deployment serializes requests carrying the same session cookie. `SerializeSession` must precede `StartSession` in the web middleware chain so the lock covers both the initial session read and the final session save. This includes GET requests: the session middleware persists them too.

Locks use 1,024 fixed stripes with a process-random hash. Cookie cardinality cannot grow the lock table. Unrelated sessions can occasionally share a stripe and wait for each other. Waiting respects request cancellation and is capped at 10 seconds; an expired wait returns HTTP 503 with Retry-After rather than starting session mutation. The lock remains held for the full lifetime of a request that acquired it, so an in-progress session write cannot overlap another request. Requests without a session cookie create independent sessions and are not serialized.

This is a single-process guarantee, not a distributed lock. Before running multiple application workers or replicas, replace this protection with a shared lock covering the whole session lifecycle, or move cart mutations to transactional storage with revision checks. A shared session file directory alone does not provide this guarantee. Concurrent requests issued during session-ID rotation are not guaranteed to share a lock after the ID changes; the existing login/logout rotation security behavior remains unchanged.

Checkout additionally uses a buyer-bound idempotency token to prevent duplicate order reservations on retries. Session serialization does not replace checkout idempotency or database stock guards.

Regression coverage: `TestConcurrentCartAddsPreserveEveryUpdate` uses a real cookie jar and 24 concurrent HTTP cart additions. Without serialization, additions are lost; with serialization every addition survives.
