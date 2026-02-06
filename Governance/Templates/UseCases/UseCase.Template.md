---
id: "UC--"
title: "" # use case title/name
bounded_context: "<BoundedContextName>"
status: draft # draft | implemented | deprecated
slices:
  - ""
---

## Intent

State the single behavioral goal this use case exists to achieve, in one or two
sentences. This must be business-meaningful and testable.

- Scope: what is in-scope and explicitly out-of-scope.
- Non-goals: outcomes this use case does not provide.

## Ownership (Bounded Context)

This use case is owned by exactly one bounded context and is defined within that
bounded context’s domain.

- Bounded context: `<BoundedContextName>`
- Domain language: terms used here must align to this bounded context’s ubiquitous
language.

## Primary actor

Who/what initiates the use case.

Examples (choose one and be specific):

- End user (via UI)
- Background task / scheduler
- System event (push, deep link, OS signal)
- External integration (incoming message)

## Preconditions

Conditions that must hold before the use case begins. Do not restate triggers.

- Required local state (e.g., authenticated session present, local store initialized)
- Required permissions/capabilities (e.g., camera permission)
- Required configuration (e.g., feature flag enabled)

## Trigger

The concrete event that causes execution.

- Entry point (UI action / system event / background schedule)
- Frequency expectations (one-shot, repeating, bursty)

## Inputs (boundary types only)

Rules:

- Inputs must be explicit boundary types (DTOs / Commands / Queries), not domain entities.
- Include correlation/idempotency keys if retries or deduplication are required.
- Mark which fields are optional vs required at the boundary.

Inputs:

- `<InputBoundaryType>` — purpose, validation expectations

## Outputs (boundary types only)

Rules:

- Outputs must be explicit boundary types (DTOs / Results), not domain entities.
- Include both success and failure shapes at the boundary (e.g., typed errors /
  failure reasons).

Outputs:

- `<OutputBoundaryType>` — success shape
- `<FailureBoundaryType>` — failure/error shape (if applicable)

## Main flow

Numbered, implementation-driving steps. Each step should be observable and testable.

1. Validate inputs at the boundary (reject invalid inputs deterministically).
2. Load required state via ports (persistence, cache, clock, network, etc.).
3. Apply domain rules/invariants owned by this bounded context.
4. Produce output boundary type(s).
5. Persist state changes via ports (if any).
6. Emit domain events / notifications (if any) via ports.

For each step, specify:

- What state is read/written
- What invariants are enforced
- What ports are called
- What failure conditions exist and how they surface at the boundary

## Alternate flows

List meaningful variations from the main flow. Use the same step-level specificity.

Alternate flow A: `<Name>`

- Condition:
- Steps that differ:
- Output:

## Side effects

Enumerate all effects outside pure computation. This must be complete.

Examples (list only those that apply):

- Local persistence writes (what data, where conceptually, durability expectation)
- Network calls (endpoints/services abstractly, retry expectations)
- Background task scheduling
- Notifications/analytics events
- File system changes
- Secure storage access

For each side effect:

- Port/protocol used
- Ordering requirements relative to other steps
- Failure handling and rollback/compensation strategy

## Idempotency & concurrency

Define correctness under retries, duplicate triggers, and parallel execution.

- Idempotency:
  - Is this use case required to be idempotent? (yes/no)
  - If yes, define the idempotency key at the boundary and what constitutes a duplicate.
  - Define what “same result” means (output equality, state equality, or both).
- Concurrency:
  - What can run concurrently (same actor, different actors)?
  - What must be serialized (per user, per resource, per account, etc.)?
  - Conflict strategy (optimistic concurrency, last-write-wins, merge, reject).

## Offline behavior

Offline-first is the default. Explicitly specify behavior in each connectivity state.

- Connectivity assumptions:
  - Works offline: (yes/no; if no, justify and document the minimum online requirement)
- When offline:
  - What the use case does using local data
  - What is queued for later sync (and the queue boundary type)
  - User-visible outcome (if any) and how “pending” is represented
- When transitioning online:
  - Sync/flush behavior, ordering, and batching
  - Consistency expectations (eventual vs strong, user-visible reconciliation)
  - Conflict detection and resolution rules
- Failure & retry policy:
  - Retry conditions, backoff, maximum attempts, cancellation
  - Poison-message handling (when retries must stop)

## Observability

Define how execution can be understood in production and during testing.

- Logging:
  - Key structured fields (use case id, correlation id, actor id, resource ids)
  - Sensitive data rules (what must never be logged)
- Metrics:
  - Success/failure counts
  - Latency (overall and key steps)
  - Retry/queue depth where applicable
- Tracing:
  - Trace/span boundaries (start/end, external calls)

## Test scenarios (Given / When / Then)

List scenarios that fully specify observable behavior. Avoid implementation details.

### Success cases

- Given `<preconditions and initial state>`
  - And `<additional setup>`
  - When `<trigger with inputs>`
  - Then `<output boundary>`
  - And `<persisted state changes>`
  - And `<emitted events / side effects>`

### Validation & rejection

- Given `<invalid input>`
  - When `<trigger>`
  - Then `<typed failure at boundary>`
  - And `<no side effects>`

