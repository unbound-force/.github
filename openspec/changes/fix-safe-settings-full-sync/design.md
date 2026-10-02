# Design

## Context

See `proposal.md` for motivation. The current workflow pins Safe Settings 2.1.18 and generates a custom full-sync runner that assumes Probot 13 initialization. Current label configuration can treat undeclared labels as deletions, and the unquoted `ci` color can be parsed as a non-string YAML value.

## Goals / Non-Goals

**Goals:**

- Upgrade to the immutable Safe Settings revision containing the rename-race fix.
- Make the existing custom runner compatible with Probot 14.
- Preserve the current narrow `check_suite` workaround.
- Restrict scoped sync input to managed repository names.
- Make label synchronization non-destructive and keep the `ci` color as a string.

**Non-Goals:**

- Redesign the workflow, rollout process, or organization governance.
- Add new evidence formats, dependency policy, approval systems, or branch-protection automation.
- Change settings unrelated to labels and full-sync compatibility.

## Decisions

### Pin the reviewed Safe Settings commit

Update the workflow from 2.1.18 to commit `6a8b6ae084987025f6c5de85e3cc6df140f64502`, which contains upstream rename-race fix PR #943. Keep the action immutable by commit SHA.

### Adapt the existing custom runner in place

Keep the workflow's current custom full-sync runner and add the Probot 14 initialization wait before logger access and application loading. This is smaller than introducing a new runner architecture.

Preserve the existing predicate `String(error).includes("Cannot read properties of undefined (reading 'check_suite')")`. Focused tests cover the exact matching error plus similar and unrelated errors that must remain fatal.

### Preserve undeclared labels

Convert label declarations to Safe Settings' additive include/exclude form so configured labels are managed while undeclared labels are preserved. Quote only the `ci` label color as `'5319e7'` for this change.

### Restrict scoped sync input

Accept only unique, plain repository names that appear in the Safe Settings managed-repository allowlist. Reject invalid, unmanaged, and duplicate names before generating the scoped deployment configuration.

## Risks / Trade-offs

- [The dependency changes behavior beyond the race fix] -> Review the pinned diff and Safe Settings dry-run output.
- [The reporting workaround becomes too broad] -> Match only the known `check_suite` failure and test that other errors fail.
- [Old labels remain present] -> Accept preservation as safer than unintended deletion; deliberate deletion requires a separate reviewed change.

## Validation

- Unit-level Go workflow/configuration tests read local files only and cover the immutable SHA, Probot readiness wait, exact workaround predicate, non-matching errors, additive label structure, and the `ci` color string.
- The six specification scenarios have 100% test mapping: each scenario is covered by at least one focused behavioral test. This change introduces no production helpers; CI executes the focused tests through `make test-unit`.
- A reviewed Safe Settings dry run is the integration check for the pinned upstream runtime; automated tests require no GitHub service or network access.
- Run the focused tests, `make safe-settings-validate`, `make test-unit`, and strict OpenSpec validation.
