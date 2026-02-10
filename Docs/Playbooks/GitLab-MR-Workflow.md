# GitLab merge request workflow (CLI: `glab`)

This playbook provides operational steps for updating and verifying a merge request description using the GitLab CLI.

## Preconditions

- You have `glab` installed and authenticated against the correct GitLab instance.
- Your git remote(s) are configured for the target GitLab project.
- You are on the branch associated with the merge request you want to update.

## Update MR description (markdown)

1. Prepare a markdown description in a file (recommended):

   - `mr.md` (or similar)

2. Update the MR description.

   - If you know the MR IID:

     - `glab mr update <iid> -d "$(cat mr.md)"`

   - If you want to target “the MR for the current branch”, first discover the MR, then update it:

     - `glab mr list --source-branch "$(git branch --show-current)"`
     - Identify the MR IID from the list output
     - `glab mr update <iid> -d "$(cat mr.md)"`

## Verify it was applied

- View the MR and confirm the description content:

  - `glab mr view <iid>`
