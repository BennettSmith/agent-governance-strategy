# GitHub pull request workflow (CLI: `gh`)

This playbook provides operational steps for updating and verifying a pull request description using the GitHub CLI.

## Preconditions

- You have `gh` installed and authenticated against the correct GitHub account/org.
- Your git remote(s) are configured for the target GitHub repository.
- You are on the branch associated with the pull request you want to update.

## Update PR description (markdown)

1. Prepare a markdown description in a file (recommended):
   - `pr.md` (or similar)

2. Update the PR body.

- If you know the PR number:
  - `gh pr edit <number> --body-file pr.md`

- If you want to target “the PR for the current branch”:
  - `gh pr edit --body-file pr.md`

## Verify it was applied

- Print the PR body and confirm it matches the expected markdown:
  - `gh pr view --json body -q .body`
