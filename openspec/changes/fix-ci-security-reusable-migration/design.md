## Context

Both `ci_security.yml` and `ci_scheduled.yml` contain two inline jobs with the same dead SHAs:
- `osv-scanner`: calls `google/osv-scanner-action/osv-scanner-action@e5012758...` (v2.0.2) — SHA gone upstream.
- `scorecards`: calls `ossf/scorecard-action@05b42c62...` (v2.4.2) — SHA gone upstream.

Both hard-fail with "unable to find version" — `ci_security.yml` on every PR, `ci_scheduled.yml` on the nightly cron. A previous attempt to fix the OSV SHA in `ci_security.yml` was reverted (`89ea9d7`) because the fix also required switching from the inline step form to OSV's own reusable workflow caller form — a structural change, not just a SHA bump.

The canonical fix pattern is established in `unbound-force/unbound-force` and `complytime/.github`, both of which replace the inline jobs with calls to `complytime/org-infra` reusable workflows. Notably, `complytime/.github` consolidates the scheduled OSV + Scorecard runs into a single `reusable_scheduled.yml` call rather than keeping separate inline jobs.

See `proposal.md` for motivation.

## Goals / Non-Goals

**Goals:**
- Replace both broken inline jobs with org-infra reusable workflow calls.
- Match the structure and permission set used in `unbound-force/unbound-force` (the org canonical).
- Pin the org-infra SHA using the same `@<sha> # <version>` comment convention used everywhere else in this repo.

**Non-Goals:**
- Adding Trivy source scanning (not in scope for this repo; `enable_trivy_source` defaults to `false` in `reusable_vuln_scan.yml`).
- Touching `ci_checks.yml` — its inline Go lint/test jobs are repo-specific and not covered by `reusable_ci.yml`.
- Touching `ci_dependencies.yml` — that is covered separately by issue #41.

## Decisions

### Decision: Pin to the same org-infra SHA as `unbound-force/unbound-force`

The org canonical (`unbound-force/unbound-force`) pins both reusable calls at `0c784711926c9864f027ec565fd7c06a382d80f8 # v0.7.1`. Using the same SHA avoids introducing drift between managed repos and is immediately verifiable against a known-good reference.

**Alternative considered**: Pin to the latest org-infra `main` SHA (`6f6dc6c9...`). Rejected because it has not been vetted as a tagged release and would silently diverge from the org canonical.

### Decision: Mirror the `unbound-force/unbound-force` permissions exactly

The org canonical sets:
- `call_reusable_vuln_scan`: `contents: read`, `actions: read`, `security-events: write`, `packages: write`, `id-token: write`
- `call_reusable_security`: `contents: read`, `id-token: write`, `security-events: write`

The `packages: write` and `id-token: write` on the vuln-scan job are required by the reusable even when trivy image scanning is skipped (the reusable's job-level permissions are inherited). Omitting them causes a permission error at runtime.

**Alternative considered**: Trim to only what the enabled inputs need. Rejected — the org-infra reusable declares its permission needs; callers must satisfy them regardless of which inputs are toggled.

### Decision: Keep the same job names as the current workflow

The current check names (`OSV-Scanner`, `OpenSSF Scorecards`) are set by the reusable workflows' own job `name:` fields and will be preserved automatically. The caller `name:` fields should match so GitHub's required-checks UI stays consistent.

### Decision: Use `reusable_scheduled.yml` for `ci_scheduled.yml` rather than separate reusable calls

`complytime/.github` consolidates the scheduled OSV + Scorecard runs into a single `reusable_scheduled.yml` call. The alternative — calling `reusable_vuln_scan.yml` and `reusable_security.yml` individually from `ci_scheduled.yml` (as `ci_security.yml` does) — would work but diverges from the peer reference. `reusable_scheduled.yml` is purpose-built for the cron use case and requires a simpler permission set than the combined individual calls.

**Alternative considered**: Call `reusable_vuln_scan.yml` + `reusable_security.yml` from `ci_scheduled.yml` the same way `ci_security.yml` does. Rejected — `reusable_scheduled.yml` already encapsulates this combination for the scheduled context and is what the peer repo uses.

## Risks / Trade-offs

- **[Risk] org-infra SHA drifts** → Mitigation: adding `dependabot.yml` for `github-actions` ecosystem (tracked in issue #41 / PR #40) will auto-propose SHA updates via Dependabot PRs.
- **[Risk] Future org-infra reusable changes break callers** → Mitigation: pinning by SHA means no surprise breakage; updates are explicit and reviewable.

## Migration Plan

1. Edit `.github/workflows/ci_security.yml` on branch `opsx/fix-ci-security-reusable-migration` — replace both inline jobs with `reusable_vuln_scan.yml` and `reusable_security.yml` caller jobs.
2. Edit `.github/workflows/ci_scheduled.yml` on the same branch — replace both inline jobs with a single `reusable_scheduled.yml` caller job.
3. Push — CI re-runs automatically on the open PR and the two failing checks should go green. The scheduled workflow fix takes effect on the next nightly cron run after merge.
4. No rollback complexity: if the reusable calls fail, reverting the two file edits restores the previous state (though the previous state also hard-fails, so a rollback is only useful to isolate a new issue).
