---
id: "UC-IDENTITY-0002"
title: "Initiate sign-in with email"
bounded_context: "Identity"
status: draft # draft | implemented | deprecated
slices:
  - "Authentication"
---

## Intent

Initiate an email-only sign-in by issuing an email challenge to the provided email
address (no password), without disclosing whether the email is registered.

- Scope: boundary validation, rate limiting, challenge issuance, email delivery attempt, non-enumerating response.
- Non-goals: completing verification, creating an authenticated session, account creation UX (handled by sign-up).

## Ownership (Bounded Context)

This use case is owned by exactly one bounded context and is defined within that
bounded context’s domain.

- Bounded context: `Identity`
- Domain language: terms used here must align to this bounded context’s ubiquitous
  language.

## Primary actor

- End user (via UI)

## Preconditions

- Email delivery capability is configured and available to this bounded context.
- Rate-limiting policy for challenge issuance is configured.

## Trigger

- Entry point: user submits the sign-in form with an email address.
- Frequency expectations: user-initiated, bursty (may retry).

## Inputs (boundary types only)

Inputs:

- `InitiateSignInWithEmailCommand` — issues an email challenge for sign-in
  - `email` (required): syntactically valid email; normalized at the boundary
  - `clientRequestId` (optional): idempotency/correlation token for retries

## Outputs (boundary types only)

Outputs:

- `SignInChallengeIssuedResult` — success shape
  - `challengeId` (opaque identifier; may be present or omitted by policy)
  - `deliveryChannel` (e.g., `email`)
  - `message` (generic: “If an account exists, you’ll receive an email”)
- `SignInRejectedError` — failure/error shape
  - `reason` (one of: `InvalidEmail`, `RateLimited`, `EmailDeliveryUnavailable`, `OfflineNotSupported`)

## Main flow

1. Validate inputs at the boundary.
   - Normalize `email` and validate format.
   - Failure: return `SignInRejectedError(InvalidEmail)`.
2. Apply rate-limiting policy (by normalized email and client context).
   - Failure: return `SignInRejectedError(RateLimited)`.
3. Load user state by normalized email via persistence ports.
4. Create a sign-in challenge (magic link or one-time code) with expiry and attempt limits.
   - Link the challenge to the user if the user exists; otherwise create a “non-account-linked” challenge or no-op per policy.
   - Persist challenge creation if applicable.
5. Attempt delivery of the challenge to the email address.
   - Failure: return `SignInRejectedError(EmailDeliveryUnavailable)`.
6. Produce `SignInChallengeIssuedResult` with a non-enumerating message.

## Alternate flows

Alternate flow A: `Email not registered`

- Condition: no user exists for the normalized email.
- Steps that differ:
  - Continue without disclosing non-existence.
  - Challenge creation may be “phantom” or real per policy, but the outward response must remain identical.
- Output:
  - `SignInChallengeIssuedResult` (indistinguishable from registered-email path).

Alternate flow B: `Idempotent retry`

- Condition: duplicate request is received (same `clientRequestId` for the same email within the idempotency window).
- Steps that differ:
  - Avoid issuing multiple emails on retry (policy-dependent).
  - Return a deterministic result consistent with the first attempt.
- Output:
  - Deterministic `SignInChallengeIssuedResult`.

## Side effects

- Local persistence writes:
  - Create challenge record (expiry, attempt limits; linked to user if applicable).
- Network calls:
  - Email delivery provider call to send a magic link / code.

For each side effect:

- Port/protocol used: persistence ports; `EmailDeliveryPort` (conceptual).
- Ordering requirements: persist challenge before sending email; do not produce “issued” before delivery attempt outcome is known.
- Failure handling: if email send fails, return `SignInRejectedError(EmailDeliveryUnavailable)`.

## Idempotency & concurrency

- Idempotency:
  - Is this use case required to be idempotent? yes (to prevent duplicate email sends on retries)
  - If yes, define the idempotency key at the boundary and what constitutes a duplicate.
    - Key: normalized `email` + `clientRequestId` (when provided); otherwise dedupe within a time window.
  - Define what “same result” means:
    - No duplicate externally-visible side effects (email sends) beyond policy and deterministic boundary output.
- Concurrency:
  - What can run concurrently: different emails may proceed in parallel.
  - What must be serialized: optionally per normalized email to enforce rate limiting and challenge policy.
  - Conflict strategy: latest-challenge-wins or bounded parallel challenges (policy-defined).

## Offline behavior

Offline-first is the default. Explicitly specify behavior in each connectivity state.

- Connectivity assumptions:
  - Works offline: no; requires online connectivity to deliver an email challenge.
- When offline:
  - Return `SignInRejectedError(OfflineNotSupported)`.
  - No queued work is created (unless deferred email delivery is explicitly supported and documented).
- Failure & retry policy:
  - Client may retry with backoff; server enforces rate limiting.

## Observability

- Logging:
  - Key structured fields: use case id, correlation id (`clientRequestId`), challenge id
  - Sensitive data rules: never log raw email address, magic link, or one-time code; log only hashed/opaque identifiers.
- Metrics:
  - Issued/rejected counts by `reason`
  - Latency end-to-end and email-delivery latency
  - Rate-limit hit counts
- Tracing:
  - Span includes persistence calls and email delivery call.
- Auditing (if applicable):
  - Record challenge issuance without exposing secrets.

## Test scenarios (Given / When / Then)

### Success cases

- Given a syntactically valid email address
  - When sign-in initiation is triggered with `InitiateSignInWithEmailCommand`
  - Then `SignInChallengeIssuedResult` is returned
  - And an email delivery attempt is made
  - And the output does not reveal whether the email is registered

### Validation & rejection

- Given an invalid email
  - When sign-in initiation is triggered
  - Then `SignInRejectedError(InvalidEmail)`
  - And no side effects occur

### Offline-first cases

- Given the device is offline
  - When sign-in initiation is triggered
  - Then `SignInRejectedError(OfflineNotSupported)`
  - And no side effects occur

### Concurrency & idempotency

- Given the same `clientRequestId` is submitted twice for the same email
  - When both requests are processed
  - Then the outcome is deterministic
  - And duplicate external email sends are avoided per policy

### Failure & recovery

- Given email delivery fails via the email port
  - When sign-in initiation is triggered
  - Then `SignInRejectedError(EmailDeliveryUnavailable)`
