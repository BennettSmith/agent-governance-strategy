# Governance sources (builder inputs)

This folder is the **authoritative source** for governance content that will be assembled into target repos by the `agent-gov` CLI.

## Structure

- `Core/`
  - Content that is intended to be shared across all repo types (working agreements, quality gates, branching strategy, etc.).
- `Profiles/`
  - Profile-specific governance for a repo type and architecture style (e.g., mobile Clean Architecture vs Go hexagonal services).
  - Each profile typically has a `profile.yaml` manifest describing which fragments, templates, and playbooks it emits.
- `Templates/`
  - Selectable templates that profiles may emit into target repos (use-case spec templates, bounded context templates, ADR/MADR templates, etc.).

## v1 profiles

- `mobile-clean-ios`
- `backend-go-hex`

## Notes

- Target repos receive *generated* governance docs (single files with managed blocks + local addenda). This folder is **not** copied verbatim into targets.
- Platform/tooling specifics should generally live as **playbooks** under the relevant profile, not in universal core fragments.

