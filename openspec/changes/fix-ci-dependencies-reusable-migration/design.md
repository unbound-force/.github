## Context

See proposal.md — Why. The current `ci_dependencies.yml` runs `actions/dependency-review-action` inline and hard-fails every PR because the Dependency Graph is unavailable. This repo already delegates its other security CI (`ci_security.yml`, `ci_scheduled.yml`) to `complytime/org-infra` reusable workflows pinned at `0c784711926c9864f027ec565fd7c06a382d80f8 # v0.7.1`, so the reusable-caller pattern and the org-infra trust boundary are already established and operationally proven here. The canonical dependency-review callers exist in `unbound-force/unbound-force/.github/workflows/ci_dependencies.yml` and `complytime/.github/.github/workflows/ci_dependencies.yml`. `safe-settings/settings.yml` sets org-wide `enableAutomatedSecurityFixes: true`; there is currently no `.github/dependabot.yml` in this repo.

## Goals / Non-Goals

**Goals:**
- Stop the dependency-review hard-fail on every PR to `main` by delegating to `reusable_deps_reviewer.yml` (which wraps the action with `continue-on-error: true`).
- Align `ci_dependencies.yml` with the org-standard reusable-caller pattern used by the sibling workflows in this repo and across the org.
- Add `.github/dependabot.yml` for the `github-actions` ecosystem so action SHA pins stay current and auditable.

**Non-Goals:**
- Enabling GitHub Advanced Security / the Dependency Graph on this repo (explicitly out of plan).
- Enabling Dependabot auto-approve / auto-merge in this repo (deferred pending a separate threat-model review — see Open Questions).
- Adding Dependabot ecosystems beyond `github-actions` (this repo has no application dependency manifests that Dependabot would act on; `go.mod` exists but is tooling-only and out of scope for this change).
- Changing `ci_security.yml`, `ci_scheduled.yml`, Peribolos config, or safe-settings config.

## Decisions

**D1 — Delegate to org-infra reusables rather than adding `continue-on-error` to the inline step.**
The minimal hotfix would be to add `continue-on-error: true` (or an `if:` guard) to the existing inline `dependency-review-action` step. Rejected in favor of the reusable-caller migration because: (a) it converges on the org-standard pattern already used by the two sibling workflows in this repo, reducing per-repo drift; (b) it matches the parallel adoption tracked in `unbound-force/replicator#38` and the canonical callers in `unbound-force/unbound-force` and `complytime/.github`; (c) the reusable centralizes the soft-gate posture and result-surfacing so it stays consistent across repos.

**D2 — Pin `complytime/org-infra` by full commit SHA with a version comment.**
Both new `uses:` lines will be pinned to `0c784711926c9864f027ec565fd7c06a382d80f8 # v0.7.1` — the exact SHA already used by `ci_security.yml` and `ci_scheduled.yml` in this repo. This matches the repo-wide pin-by-SHA convention and keeps org-infra references uniform. Alternative (floating tag like `@v0.7.1`) rejected: it violates the repo convention and re-introduces the mutable-ref supply-chain risk the SHA pin exists to prevent. Implementation note: before writing, confirm `reusable_deps_reviewer.yml` and `reusable_dependabot_reviewer.yml` exist at that SHA; if v0.7.1 predates those files, bump both org-infra references (all callers in this repo) to the earliest SHA that contains them.

**D3 — Exclude the Dependabot auto-approve/auto-merge job from this repo.**
The canonical pattern includes an auto-approve flow. It is excluded here because `unbound-force/.github` controls org membership (Peribolos) and repository security settings (safe-settings) for the entire org; auto-merging Dependabot PRs into this repo is an elevated attack surface (e.g., a malicious action SHA substitution auto-merged into org policy). The `reusable_dependabot_reviewer.yml` caller is included (review/labeling only), but no auto-approve job is added. Re-enabling is gated on a separate threat-model review.

**D4 — Add `.github/dependabot.yml` scoped to `github-actions` only.**
This is the only ecosystem relevant to this repo's CI, and it is the prerequisite that makes SHA pins maintainable. Weekly schedule, consistent with typical org config. No `open-pull-requests-limit` change beyond defaults is required for this change. `go.mod` exists but is tooling-only and deliberately excluded (see Non-Goals); no other ecosystem entries are added.

**D5 — Permissions posture for the caller jobs.**
Verified against org-infra at the pinned SHA: both `reusable_deps_reviewer.yml` and `reusable_dependabot_reviewer.yml` are passive (review-only, output-emitting) and declare `permissions: { contents: read, issues: none, pull-requests: none }`. The `pull-requests: write` scope in the canonical `ci_dependencies.yml` is granted only at the job level for the comment/auto-approve jobs — which this change deliberately excludes. Therefore the migrated `ci_dependencies.yml` keeps a restrictive top-level block matching the canonical caller: `contents: read`, `issues: none`, `pull-requests: none`. No per-job `permissions:` overrides are needed because the two included callers require no write scopes.

## Risks / Trade-offs

- **Dependency-review signal becomes advisory (`continue-on-error: true`).** If GHAS/Dependency Graph is ever enabled on this repo, real vulnerability findings would no longer block PRs — they would pass silently. → Mitigation: the reusable captures the result as an output and surfaces it (job summary / annotation) rather than discarding it; document the soft-gate posture in the workflow. Re-evaluate the gate if GHAS is ever enabled.
- **CI availability coupling to `complytime/org-infra`.** An org-infra outage or a bad reusable revision could break the dependency CI job. → Mitigation: this coupling is identical to and no worse than the coupling already accepted by `ci_security.yml`/`ci_scheduled.yml`, and the SHA pin prevents unreviewed upstream changes from flowing in.
- **Adding Dependabot PRs where there were none.** Enabling `github-actions` updates will start generating Dependabot PRs. → Mitigation: acceptable and desired (keeps SHAs current); no auto-approve means every PR still requires human review, so there is no unattended-merge risk.

## Migration Plan

1. Replace the `dependency-review` job (checkout + `dependency-review-action` steps) in `ci_dependencies.yml` with `call_deps_reviewer` and `call_dependabot_reviewer` caller jobs pinned per D2.
2. Add `.github/dependabot.yml` with the `github-actions` ecosystem entry.
3. Validate YAML locally (`make sanity` / `yamllint`).
4. Open a PR and confirm the Dependency Review check no longer hard-fails and no inline `dependency-review-action` step remains.
5. Rollback: revert the two files (`ci_dependencies.yml`, `.github/dependabot.yml`) to their prior state; there is no state or data migration to unwind.

## Open Questions

- **Should Dependabot auto-approve ever be enabled for `unbound-force/.github`?** Deferred to a separate security review with an explicit threat model and, if approved, compensating controls (e.g., required human review on top of auto-approve, restricted Dependabot scope). This does not affect the specs, approach, or task breakdown of this change — auto-approve is out of scope here regardless of the eventual answer.
