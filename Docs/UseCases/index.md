# Use-Case Catalog (Authoritative)

In this system, **use cases are the sole source of application behavior**. Any
behavior that matters to users or integrations must be represented as a documented
use case, owned by exactly one bounded context.

## What a use case is (in this repository)

A use case is a bounded, testable unit of application behavior that:

- Lives in exactly one **bounded context** (implemented as its own module/package; platform specifics are documented in playbooks)
- Accepts explicit **input boundary types** (no domain entities)
- Produces explicit **output boundary types** (no domain entities)
- Enforces domain invariants owned by its bounded context
- Is **offline-first by default** unless explicitly documented otherwise

Feature modules (vertical slices) **compose** use cases to deliver user outcomes,
but **do not define behavior**.

## This catalog is the source of truth

This catalog is the **single authoritative list of system behavior**. If it is
not in this catalog (with a linked spec), it is not considered defined behavior.

## Catalog

| ID | Title | Bounded context | Status | Slices | Spec |
|---|---|---|---|---|---|
| `UC-IDENTITY-0001` | Sign up with email and name | `Identity` | draft | Authentication, Onboarding | [`UC-IDENTITY-0001-SignUpWithEmailAndName.md`](./UC-IDENTITY-0001-SignUpWithEmailAndName.md) |
| `UC-IDENTITY-0002` | Initiate sign-in with email | `Identity` | draft | Authentication | [`UC-IDENTITY-0002-InitiateSignInWithEmail.md`](./UC-IDENTITY-0002-InitiateSignInWithEmail.md) |
| `UC-IDENTITY-0003` | Sign out | `Identity` | draft | Authentication | [`UC-IDENTITY-0003-SignOut.md`](./UC-IDENTITY-0003-SignOut.md) |

## Canonical use-case spec template

All use cases must be documented using the canonical spec template:

- [`UseCase.Template.md`](./UseCase.Template.md)

