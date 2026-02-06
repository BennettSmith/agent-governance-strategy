# Identity — Bounded Context

## Purpose

Provide identity lifecycle and authentication-state management for an end user across:

- Creating (or reusing) a user account based on email + name.
- Initiating email-only sign-in in a non-enumerating way.
- Signing out by clearing local authentication state and optionally revoking a server-side session best-effort.

## Scope

### In Scope

- Boundary validation and normalization of user-provided email and name.
- Creating a user record if one does not already exist for a normalized email.
- Loading user state by normalized email.
- Enforcing rate-limiting policy for sign-up / sign-in initiation.
- Maintaining local authentication state for “current session on this device”, including sign-out.
- Best-effort server-side session revocation when online (does not block sign-out).
- Idempotency semantics (via optional `clientRequestId`) for retry-safe behavior.
- Offline-first behavior as specified by the use cases.

### Out of Scope

- Modeling “challenges”, magic links, one-time codes, or any verification-token lifecycle.
- Completing email ownership verification / completing sign-in (establishing an authenticated session).
- Passwords and password-based authentication.
- Account deletion.
- Global sign-out (“sign out on all devices”).
- Email content, templating, and deliverability concerns beyond “attempt delivery and observe success/failure”.

## Ubiquitous Language

- **User** — A person account identified by a unique normalized email within Identity.
- **Email** — The user-provided email address.
- **Normalized Email** — The email after boundary normalization (case/whitespace) used for identity lookups and policy enforcement.
- **Display Name** — The user-provided name subject to Identity rules (length/characters).
- **Sign-Up** — Creating (or reusing) a user record and initiating email-only sign-in initiation, without creating an authenticated session.
- **Sign-In Initiation** — Starting email-only sign-in in a non-enumerating way; does not create an authenticated session.
- **Non-enumerating response** — A response that does not disclose whether an email is registered.
- **Session** — The current authenticated session on this device, represented by local authentication material (tokens/credentials) and associated cached authorization data.
- **Sign-Out** — Clearing local authentication material and persisting a signed-out local state; may attempt best-effort server session revocation.
- **Client Request ID** — Optional idempotency/correlation token supplied by the client for retry-safe operations.
- **Rate limiting** — Policy that may reject sign-up / sign-in initiation requests for a normalized email + client context.
- **Secure storage** — Local secret storage used to store/delete session credentials (e.g., keychain).

## Domain Model (High Level)

### Entities

- **User**
  - Identity: `UserId` (opaque identifier)
  - Responsibility:
    - Own the user’s normalized email and display name.
    - Represent existence of an account for an email, including any registration state required by current use cases.
  - Lifecycle Notes:
    - Created on sign-up if no user exists for normalized email.
    - Loaded by normalized email for sign-up (reuse path) and sign-in initiation.
- **Session (Device Session)**
  - Identity: `SessionId` (opaque identifier; may be absent when signed out)
  - Responsibility:
    - Represent “current session on this device” sufficient to clear local authentication material deterministically on sign-out.
  - Lifecycle Notes:
    - Creation/establishment of an authenticated session is out of scope for the current use cases.
    - Sign-out transitions session to a signed-out local state.

### Value Types

- **EmailAddress**
  - Meaning: User-provided email address in normalized form for lookup/policy decisions.
  - Constraints:
    - Must be syntactically valid.
    - Must be normalized at the boundary (case/whitespace) consistently.
- **DisplayName**
  - Meaning: User-facing name captured at sign-up.
  - Constraints:
    - Must meet Identity “display-name rules” (length/characters), enforced at the boundary.
- **ClientRequestId**
  - Meaning: Optional idempotency/correlation token for retry-safe use case execution.
  - Constraints:
    - Opaque; compared for equality within an idempotency window (policy-defined).
- **SessionCredentials**
  - Meaning: Local authentication material required to act as signed in (e.g., access/refresh tokens or equivalent).
  - Constraints:
    - Stored in secure storage when present.
    - Must be fully removed on successful sign-out.
- **AuthState**
  - Meaning: The locally persisted authentication state for the app (signed-in vs signed-out) including any cached authorization data that must be cleared on sign-out.
  - Constraints:
    - Signed-out state must not retain session credentials or cached authorization data.

### Aggregates

- **User (Aggregate Root)**
  - Members: `User`, `EmailAddress`, `DisplayName`
  - Invariants Enforced:
    - A normalized email corresponds to at most one user within Identity.
    - Sign-up must not create duplicate users for the same normalized email (serialize/lock per email).
- **Session (Aggregate Root)**
  - Members: `Session`, `SessionCredentials`, `AuthState`
  - Invariants Enforced:
    - After sign-out completes, the local session state is signed out (no credentials/cached authorization data remain).
    - Sign-out is safe and deterministic when already signed out.

### Domain Events (Optional)

- **UserAccountCreated**
  - Occurs When: A new user record is created for a normalized email during sign-up.
  - Key Data: `userId`, (hashed or opaque) email identifier
- **SignedOut**
  - Occurs When: Local authentication material is cleared and signed-out state is persisted.
  - Key Data: `sessionId` (if available), correlation (`clientRequestId` if provided)

## Invariants

- Inputs are validated at the boundary:
  - Email must be normalized and syntactically valid.
  - Display name must meet Identity rules.
- Sign-up and sign-in initiation enforce rate limiting by normalized email and client context.
- Sign-up and sign-in initiation responses must be non-enumerating with respect to email registration.
- Idempotent retries (when `clientRequestId` is provided) must not create duplicate user records and should avoid duplicate externally visible side effects beyond policy.
- Offline behavior:
  - Sign-up and sign-in initiation reject when offline (no queued work is created unless explicitly supported and documented elsewhere).
  - Sign-out must succeed offline (local-only completion).
- Sensitive data handling:
  - Raw email addresses and any secrets must not be logged; only hashed/opaque identifiers are permitted.
  - Tokens/credentials must never be logged.

## Boundary Rules

- Domain models are internal to this bounded context.
- Domain entities and value types must not cross use-case boundaries.
- Interaction occurs only via application-layer use cases.
- Cross-context interaction uses identifiers, boundary types, or events—not shared models.
- External dependencies used by this context (conceptual ports):
  - Email delivery capability (used by sign-up and sign-in initiation).
  - Server-side session revocation capability (best-effort on sign-out when online).

## Mapping Notes

- Use case boundary types (commands/results/errors) map to/from domain models in the application layer.
- Persistence and external API representations map to/from domain models in the infrastructure layer.
- Mapping logic must be explicit, especially for:
  - Email normalization (boundary responsibility).
  - Non-enumerating output semantics (application responsibility).
  - Secure storage failures (typed rejections for sign-out).

## Offline-First Considerations

- Sign-up (`UC-IDENTITY-0001`) and sign-in initiation (`UC-IDENTITY-0002`) require online connectivity due to email delivery; when offline they return typed rejections (`OfflineNotSupported`).
- Sign-out (`UC-IDENTITY-0003`) must complete offline by clearing local state; any server revoke is skipped or queued only if explicitly supported and documented.
- Invariants enforced locally must not require server round-trips (e.g., “signed out means no local credentials remain”).

## Open Questions

- What are the explicit “display-name rules” (length/character set) and where are they specified?
- Should a `challengeId` ever be omitted from the boundary result (as suggested by policy), and what is the consumer expectation?
- For sign-out, what is the policy when a provided `sessionId` does not match the current session (reject vs treat as already signed out)?
- Is there an explicit policy for queuing a best-effort server revoke when offline, and what are the maximum attempts/backoff?

## Related Use Cases

- `UC-IDENTITY-0001` — Sign up with email and name
- `UC-IDENTITY-0002` — Initiate sign-in with email
- `UC-IDENTITY-0003` — Sign out

## Change Log

- 2026-02-05 — Initial version

