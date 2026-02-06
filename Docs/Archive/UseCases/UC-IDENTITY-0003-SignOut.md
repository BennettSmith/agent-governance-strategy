---
id: "UC-IDENTITY-0003"
title: "Sign out"
bounded_context: "Identity"
status: draft # draft | implemented | deprecated
slices:
  - "Authentication"
---

## Intent

End the current authenticated session on the device by clearing local authentication
state. Optionally revoke the server-side session when online.

- Scope: local sign-out completion (offline-capable) and best-effort server revocation.
- Non-goals: deleting the user account, revoking sessions on other devices (unless explicitly supported by Identity policy), resetting email address.

## Ownership (Bounded Context)

This use case is owned by exactly one bounded context and is defined within that
bounded context’s domain.

- Bounded context: `Identity`
- Domain language: terms used here must align to this bounded context’s ubiquitous
language.

## Primary actor

- End user (via UI)

## Preconditions

- A local session may or may not be present; this use case must be safe to run when already signed out.
- Local secure storage (e.g., keychain) is accessible, or failures are handled deterministically.

## Trigger

- Entry point: user taps “Sign out” from an account/settings UI surface.
- Frequency expectations: infrequent, user-initiated; may be re-triggered due to UI retries.

## Inputs (boundary types only)

Inputs:

- `SignOutCommand` — signs out the current session on this device
  - `sessionId` (optional): session identifier to revoke/audit; if omitted, uses the “current session”
  - `clientRequestId` (optional): idempotency/correlation token for retries

## Outputs (boundary types only)

Outputs:

- `SignOutCompletedResult` — success shape
  - `signedOut` (always `true`)
- `SignOutRejectedError` — failure/error shape (rare)
  - `reason` (one of: `LocalStateCorrupt`, `SecureStorageUnavailable`)

## Main flow

1. Validate inputs at the boundary.
   - If `sessionId` is provided, ensure it matches the current session (or treat mismatch deterministically per policy).
2. Load current authentication state via persistence/secure-storage ports.
3. Clear local authentication material.
   - Delete access/refresh tokens (or equivalent session credentials).
   - Clear cached identity claims/user identifiers used for authorization decisions.
   - Clear any “pending verification” challenge state associated with the signed-in session (if applicable).
4. Persist the signed-out state locally (so the app behaves as signed out on restart).
5. If online and a server-side session exists, attempt best-effort session revocation via a network port.
   - If revocation fails, do not block local sign-out completion.
6. Produce `SignOutCompletedResult`.

## Alternate flows

Alternate flow A: `Already signed out`

- Condition: no local session exists at the start of execution.
- Steps that differ:
  - No token deletion required (local state is already empty).
  - Network revoke is skipped (no session to revoke) unless policy supports a “revoke by sessionId” call.
- Output:
  - `SignOutCompletedResult`.

Alternate flow B: `Secure storage unavailable`

- Condition: secure storage port fails (e.g., keychain locked/unavailable).
- Steps that differ:
  - Do not partially clear state.
  - Return a deterministic typed failure.
- Output:
  - `SignOutRejectedError(SecureStorageUnavailable)`.

Alternate flow C: `Server revoke fails`

- Condition: network port fails or device is offline.
- Steps that differ:
  - Local sign-out still completes.
  - Optionally queue a revoke attempt for later (only if such queueing is explicitly supported and documented).
- Output:
  - `SignOutCompletedResult`.

## Side effects

Enumerate all effects outside pure computation. This must be complete.

- Local persistence writes:
  - Delete/clear authentication credentials from secure storage.
  - Clear local session markers and cached authorization data.
- Network calls (optional):
  - Revoke the current server-side session (best effort).

For each side effect:

- Port/protocol used: secure-storage port; local persistence port; `SessionRevocationPort` (conceptual).
- Ordering requirements: local clears must happen before returning success; server revoke must not block completion.
- Failure handling and rollback/compensation strategy:
  - If local clear fails, return a typed rejection and leave state unchanged.
  - If server revoke fails, complete locally and optionally queue revoke per documented policy.

## Idempotency & concurrency

- Idempotency:
  - Is this use case required to be idempotent? yes
  - If yes, define the idempotency key at the boundary and what constitutes a duplicate.
    - Key: current local session (or provided `sessionId`) + optional `clientRequestId`.
  - Define what “same result” means:
    - State equality: local state is “signed out” and remains so after retries; no duplicate adverse side effects.
- Concurrency:
  - What can run concurrently: sign-out should be serialized with other session-mutating use cases on this device.
  - What must be serialized: token/session clearing operations.
  - Conflict strategy: single-flight execution; last-write-wins for local “signed out” marker.

## Offline behavior

Offline-first is the default. Explicitly specify behavior in each connectivity state.

- Connectivity assumptions:
  - Works offline: yes; local sign-out must succeed offline.
- When offline:
  - Clear local auth state and return `SignOutCompletedResult`.
  - Any server-side revoke is skipped or queued (only if explicitly supported).
- When transitioning online:
  - If revoke is queued, attempt it with backoff and stop after policy-defined maximum attempts.
- Failure & retry policy:
  - Retries are safe (idempotent); secure storage failures should not be retried in a tight loop.

## Observability

- Logging:
  - Key structured fields: use case id, correlation id (`clientRequestId`), session id (opaque/non-secret identifier only)
  - Sensitive data rules: never log tokens, secrets, or personally identifying information.
- Metrics:
  - Completion count
  - Local-clear failures by `reason`
  - Server-revoke attempts/failures (if applicable)
- Tracing:
  - Span includes local clears and optional revoke call.
- Auditing (if applicable):
  - Record sign-out action (actor = current user/device) without exposing secrets.

## Test scenarios (Given / When / Then)

### Success cases

- Given an authenticated session exists locally
  - When the user triggers `SignOutCommand`
  - Then `SignOutCompletedResult`
  - And local authentication material is cleared
  - And the app is in a signed-out state on restart

### Validation & rejection

- Given secure storage is unavailable
  - When sign-out is triggered
  - Then `SignOutRejectedError(SecureStorageUnavailable)`
  - And no partial clear occurs

### Offline-first cases

- Given the device is offline
  - When sign-out is triggered
  - Then `SignOutCompletedResult`
  - And local authentication material is cleared
  - And no network call is required for completion

### Concurrency & idempotency

- Given sign-out is triggered twice concurrently
  - When both executions run
  - Then the outcome is deterministic
  - And local state ends as signed out
  - And no duplicate adverse side effects occur

### Failure & recovery

- Given server revoke fails
  - When sign-out is triggered while online
  - Then `SignOutCompletedResult`
  - And local authentication material is cleared
  - And the revoke failure is observable (metrics/logs)

