# Go CLI structure (recommended)

This playbook describes the preferred structure for Go-based CLI tools in repos that adopt the `backend-go-hex` profile.

## Goals

- Keep CLI entrypoints thin and testable.
- Keep domain logic (or governance logic) in `internal/` packages, not in `cmd/`.
- Make behavior deterministic and automation-friendly (exit codes, stderr/stdout discipline, reproducible flags/config).

## Recommended layout

- `cmd/<tool>/main.go`
  - parse flags
  - load config
  - call a testable runner (e.g., `cli.Run(args, stdout, stderr)`), then `os.Exit(code)`
- `internal/cli/`
  - command dispatch
  - config resolution (including “run from wherever you are” ergonomics)
  - printing usage/errors
- `internal/<domain>/`
  - the real behavior (builder/sync/verify, ports/adapters, etc.)

Rule of thumb: `cmd/` should mostly be wiring; `internal/` is where logic and tests live.

## Testability patterns

- Prefer a `Run(args []string, stdout, stderr io.Writer) int` entrypoint.
- Keep `main()` minimal and move exiting behavior to a small wrapper (e.g., `mainWithExit`), so tests can exercise control flow without terminating the test runner.
- For process execution, use an injectable function variable (e.g., `execCommandContext`) so tests can stub it.

## Output and exit codes

- Use **stdout** for normal output and **stderr** for errors/usage.
- Return stable exit codes:
  - `0` success
  - `1` operation failed (e.g., verify failed)
  - `2` usage/config errors (bad flags, invalid config)

## Config and working directory

- Prefer `.governance/config.yaml` (or equivalent) as the default config location.
- Support `--config PATH` overrides.
- If `--config` is omitted, search upward from the current working directory for the nearest config so the tool can be run from subdirectories.
- When running repo-mutating operations, use the config location to determine the effective repo root (don’t assume CWD is the repo root).
