## Architectural identity (builder repo)

- This repository is a **governance builder**. It produces governance document bundles for other repositories via `agent-gov`.
- The authoritative source for generated governance content lives under `Governance/` (core fragments, profiles, templates, and playbooks).
- Generated governance is applied to target repositories using **managed blocks** (agent-owned) while allowing **local addenda** (project-owned) to coexist without being overwritten.

