# Agent checkpoint workflow

This playbook defines the **operational workflow** for agent-run work in repositories that use branch plans (`Docs/Plans/<branch-name>.md`) and checkpoint approval gates.

This playbook is intentionally tool-agnostic. Its purpose is to make the “stop for approval” rule concrete and repeatable.

## Definitions

- **Branch plan**: `Docs/Plans/<branch-name>.md` created from the plan template, with a checklist of checkpoints.
- **Checkpoint**: a small, reviewable unit of work in the branch plan. A checkpoint is “complete” when:
  - its stated outcomes are implemented, and
  - relevant quality gates for that checkpoint have been run (as applicable), and
  - the agent is ready to present the results for review.
- **Stop**: end the agent’s work and wait for explicit human direction/approval before doing additional work.

## Required approval ritual (per checkpoint)

When a checkpoint is complete, the agent must stop and request approval **before proceeding past that checkpoint**, even if no commit has been created yet.

Use the template below.

### Checkpoint approval request template

Include the following sections in the approval request:

- **Summary**: what was accomplished in this checkpoint (1–3 bullets)
- **Changes**: key files/areas touched (short list)
- **Quality gates run**: what was executed and what passed (or what is intentionally deferred)
- **Risks / follow-ups**: anything notable that might affect review or later checkpoints
- **Next checkpoint**: what work will begin next if approved

Then include the required approval phrase:

`APPROVAL REQUEST (checkpoint X): Please approve proceeding past checkpoint X.`

## Human responses (expected)

Humans should respond with one of:

- **Approved**: proceed to the next checkpoint.
- **Changes requested**: the agent should adjust the work within the same checkpoint and re-request approval.
- **Re-scope**: split/merge checkpoints or update the plan scope before proceeding.

## Commit guidance

- If checkpoint approval is granted, the agent may create a checkpoint commit (if appropriate for the repo/workflow), and must keep the branch plan up to date.
- If checkpoint approval is not granted, the agent must not create a checkpoint commit for that checkpoint.
