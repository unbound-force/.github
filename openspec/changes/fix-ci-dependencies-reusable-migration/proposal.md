## Why

`.github/workflows/ci_dependencies.yml` runs `actions/dependency-review-action@da24556b… (v4.7.1)` as a bare inline step with no `continue-on-error` guard. Because the Dependency Graph feature is not enabled on this repo (and there is no plan to enable GitHub Advanced Security here), the action hard-fails on **every** pull request to `main` with "Dependency review is not supported on this repository". This is a 100%-reproducible blocking CI gate, and the inline form also diverges from the org-standard reusable-caller pattern already adopted in this repo's `ci_security.yml` / `ci_scheduled.yml` and tracked for `unbound-force/replicator#38`.

## What Changes

- Replace the inline `dependency-review` job in `ci_dependencies.yml` with a caller to `complytime/org-infra/.github/workflows/reusable_deps_reviewer.yml`, which wraps `dependency-review-action` with `continue-on-error: true` so a missing Dependency Graph no longer blocks PRs; the result is captured as an output and surfaced informally rather than as a hard gate.
- Add a caller to `complytime/org-infra/.github/workflows/reusable_dependabot_reviewer.yml` for Dependabot-authored PRs, matching the canonical pattern in `unbound-force/unbound-force` and `complytime/.github`.
- **Explicitly exclude** the Dependabot auto-approve / auto-merge job from scope for this repo. Because `unbound-force/.github` governs org policy (Peribolos membership, safe-settings, rulesets), auto-approving Dependabot PRs here is an elevated attack surface; it is deferred pending a separate threat-model review.
- Pin both `uses:` references to `complytime/org-infra` at a full commit SHA with a version comment, matching the existing pin-by-SHA convention (`0c784711…926c9864f027ec565fd7c06a382d80f8 # v0.7.1`) already used by `ci_security.yml` and `ci_scheduled.yml`.
- Add `.github/dependabot.yml` for the `github-actions` ecosystem so action SHA updates are proposed automatically. This complements `safe-settings/settings.yml`'s org-wide `enableAutomatedSecurityFixes: true` and keeps action pins auditable.
- Remove the existing inline `actions/dependency-review-action` step and its `actions/checkout` step.

## Capabilities

### New Capabilities
- `ci-dependencies`: The dependency CI workflow (`ci_dependencies.yml`) SHALL delegate dependency review to org-infra reusable workflows rather than running an inline `dependency-review-action` step, and the repo SHALL declare a `github-actions` Dependabot configuration.

### Modified Capabilities
<!-- None — no existing specs under openspec/specs/. -->

## Impact

- `.github/workflows/ci_dependencies.yml` — inline `dependency-review` job (checkout + dependency-review-action steps) replaced by reusable-workflow callers.
- `.github/dependabot.yml` — new file declaring the `github-actions` ecosystem for automated action SHA updates.
- New CI availability coupling to `complytime/org-infra` for the dependency-review callers — identical in nature to the coupling already accepted by `ci_security.yml` and `ci_scheduled.yml`.
- The Dependency Review check name visible in GitHub CI changes from the inline job to the reusable-workflow caller job name(s) (`General`, `Dependabot`).
- No source files, Peribolos config, or safe-settings config are affected.
- Unblocks all future PRs to `main`, which currently hit a hard `dependency-review` failure. Same adoption is tracked in parallel for `unbound-force/replicator#38`.
