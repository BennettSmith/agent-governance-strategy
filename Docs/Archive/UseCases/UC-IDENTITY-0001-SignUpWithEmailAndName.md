---
id: "UC-IDENTITY-0001"
title: "Sign up with email and name"
bounded_context: "Identity"
status: draft # draft | implemented | deprecated
slices:
  - "Authentication"
  - "Onboarding"
---

## Intent

Create a new user account using an email address and name (no password), and begin
email ownership verification via an emailed sign-in challenge.

- Scope: input validation, account creation (or reuse), challenge issuance, and email delivery attempt.
- Non-goals: password creation, completing email verification, establishing an authenticated session.

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
- Rate-limiting policy for sign-up/sign-in challenge issuance is configured.

## Trigger

- Entry point: user submits the sign-up form.
- Frequency expectations: user-initiated, bursty (may retry due to latency).

## Inputs (boundary types only)

Inputs:

- `SignUpWithEmailAndNameCommand` — creates or resumes an account sign-up flow and issues a challenge
  - `email` (required): must be a syntactically valid email; normalized at the boundary
  - `name` (required): must meet display-name rules (length/characters)
  - `clientRequestId` (optional): idempotency/correlation token for retries

## Outputs (boundary types only)

Outputs:

- `SignUpAcceptedResult` — success shape
  - `challengeId` (opaque identifier)
  - `deliveryChannel` (e.g., `email`)
  - `message` (generic user-facing message such as “Check your email”)
- `SignUpRejectedError` — failure/error shape
  - `reason` (one of: `InvalidEmail`, `InvalidName`, `RateLimited`, `EmailDeliveryUnavailable`, `OfflineNotSupported`)

## Main flow

1. Validate inputs at the boundary.
   - Normalize `email` (case/whitespace) and validate format.
   - Validate `name` per Identity rules.
   - Failure: return `SignUpRejectedError(InvalidEmail|InvalidName)`.
2. Apply rate-limiting policy for challenge issuance (by normalized email and client context).
   - Failure: return `SignUpRejectedError(RateLimited)`.
3. Load existing user state by normalized email via persistence ports.
4. If no user exists for the email, create a new user in a “pending verification” state.
   - Persist user creation.
5. Create a sign-in/verification challenge with expiry and attempt limits.
   - Persist challenge creation.
6. Attempt delivery of the challenge to the email address (e.g., magic link or one-time code).
   - Failure: return `SignUpRejectedError(EmailDeliveryUnavailable)` (do not produce an authenticated session).
7. Produce `SignUpAcceptedResult` with an intentionally non-enumerating message.

## Alternate flows

Alternate flow A: `Email already registered`

- Condition: a user already exists for the normalized email.
- Steps that differ:
  - Do not create a new user.
  - Continue by issuing and delivering a sign-in challenge as usual (avoid account enumeration).
- Output:
  - `SignUpAcceptedResult` (indistinguishable from the “new user” path).

Alternate flow B: `Idempotent retry`

- Condition: duplicate request is received (same `clientRequestId` for the same email within the idempotency window).
- Steps that differ:
  - Do not create duplicate users or duplicate challenges.
  - Return the same logical `challengeId` and output semantics, or deterministically re-issue per policy.
- Output:
  - Deterministic `SignUpAcceptedResult`.

## Side effects

- Local persistence writes:
  - Create user (if new) in “pending verification” state.
  - Create challenge record (expiry, attempt limits, linkage to user if applicable).
- Network calls:
  - Email delivery provider call to send a magic link / code.

For each side effect:

- Port/protocol used: persistence ports; `EmailDeliveryPort` (conceptual).
- Ordering requirements: persist user/challenge before sending email; do not emit “accepted” before delivery attempt outcome is known.
- Failure handling: if email send fails, return `SignUpRejectedError(EmailDeliveryUnavailable)`; no authenticated session is created.

## Idempotency & concurrency

- Idempotency:
  - Is this use case required to be idempotent? yes
  - If yes, define the idempotency key at the boundary and what constitutes a duplicate.
    - Key: normalized `email` + `clientRequestId` (when provided); otherwise policy-defined dedupe window by email.
  - Define what “same result” means:
    - State equality (no duplicate users) and no duplicate externally-visible side effects beyond policy (e.g., avoid multiple emails on retry).
- Concurrency:
  - What must be serialized: per normalized email (prevent duplicate account creation).
  - Conflict strategy: serialize/lock per email; deterministic outcome.

## Offline behavior

Offline-first is the default. Explicitly specify behavior in each connectivity state.

- Connectivity assumptions:
  - Works offline: no; requires online connectivity to deliver an email challenge.
- When offline:
  - Return `SignUpRejectedError(OfflineNotSupported)`.
  - No queued work is created (unless the system explicitly supports deferred email delivery, which must be documented).
- Failure & retry policy:
  - Client may retry with exponential backoff; server enforces rate limiting.

## Observability

- Logging:
  - Key structured fields: use case id, correlation id (`clientRequestId`), challenge id
  - Sensitive data rules: never log raw email address, magic link, or one-time code; if needed, log a one-way hash of normalized email.
- Metrics:
  - Success/failure counts by `reason`
  - Latency end-to-end and email-delivery latency
  - Rate-limit hit counts
- Tracing:
  - Span includes persistence calls and email delivery call.
- Auditing (if applicable):
  - Record challenge issuance without exposing secrets.

## Test scenarios (Given / When / Then)

### Success cases

- Given a valid email and a valid name
  - When the user triggers sign-up with `SignUpWithEmailAndNameCommand`
  - Then `SignUpAcceptedResult` is returned
  - And a user exists in “pending verification” state (if previously absent)
  - And a challenge is created
  - And an email delivery attempt is made

### Validation & rejection

- Given an invalid email
  - When sign-up is triggered
  - Then `SignUpRejectedError(InvalidEmail)`
  - And no persistence writes occur
  - And no email is sent

- Given an invalid name
  - When sign-up is triggered
  - Then `SignUpRejectedError(InvalidName)`
  - And no persistence writes occur
  - And no email is sent

### Offline-first cases

- Given the device is offline
  - When sign-up is triggered
  - Then `SignUpRejectedError(OfflineNotSupported)`
  - And no side effects occur

### Concurrency & idempotency

- Given the same `clientRequestId` is submitted twice for the same email
  - When both requests are processed
  - Then the outcome is deterministic
  - And no duplicate user is created
  - And duplicate external email sends are avoided per policy

### Failure & recovery

- Given email delivery fails via the email port
  - When sign-up is triggered
  - Then `SignUpRejectedError(EmailDeliveryUnavailable)`
  - And no authenticated session is created

