# Governance builder architecture

This repository is not a product/runtime system. It is a **governance source + toolchain** that generates and synchronizes governance documents into target repositories.

## Core concepts

- **Fragments**: Markdown fragments that are assembled into output documents. Fragments live under `Governance/Core/` (shared across profiles) and `Governance/Profiles/<id>/` (profile-specific overlays).
- **Profiles**: Manifests (`Governance/Profiles/<id>/profile.yaml`) that declare which documents, templates, and playbooks a target repo should receive.
- **Templates**: Reusable doc templates (plans, decision records, etc.) that profiles can emit into targets.
- **Playbooks**: Optional, profile-specific operational guidance (e.g., packaging, ports/adapters, platform guidance).

## Generated documents and managed blocks

Target repositories receive generated governance documents that contain:

- **Managed blocks**: sections owned by governance and synchronized deterministically.
- **Local addenda**: a project-owned section that remains editable and is not overwritten by sync.

The synchronization mechanism updates only managed blocks and preserves local addenda and any other content outside managed sections.

## `agent-gov` workflow

`agent-gov` is the vendorable CLI that performs governance operations in a target repository:

- `init`: create or initialize governance docs with managed blocks and addenda.
- `sync`: update managed blocks in-place.
- `verify`: check that managed blocks match the expected governance content.
- `build`: assemble a governance bundle into an output folder (useful for inspection or CI artifacts).

## Configuration

Targets configure governance via `.governance/config.yaml` (schema versioned). The CLI supports running from any working directory by selecting the nearest config when `--config` is omitted.
