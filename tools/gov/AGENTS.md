# AGENTS (embedded: tools/gov)

This folder is an **embedded target repo** (the `agent-gov` CLI). When working under `tools/gov/**`, the governing documents are the key-three located in this folder:

1. `Non-Negotiables.md`
2. `Architecture.md`
3. `Constitution.md`

### Precedence (conflicts)

- `Non-Negotiables.md` overrides all other documents
- `Architecture.md` overrides `Constitution.md` on matters of system shape
- `Constitution.md` guides behavior when other documents are silent
- If anything is ambiguous, stop and ask for human direction

## Notes

- These embedded key-three docs are generated/synced via `tools/gov/.governance/config.yaml`.

