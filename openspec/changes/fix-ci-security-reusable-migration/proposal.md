## Why

`ci_security.yml` and `ci_scheduled.yml` both pin two action SHAs that no longer resolve: `google/osv-scanner-action@e5012758...` (v2.0.2) and `ossf/scorecard-action@05b42c62...` (v2.4.2). The dead SHAs cause a hard CI failure on every PR — including PR #40 — with "unable to find version", and will fail the nightly cron job identically. The upstream actions have moved on and those SHAs are gone. Adopting the org-infra reusable workflows (`reusable_vuln_scan.yml`, `reusable_security.yml`, and `reusable_scheduled.yml`) fixes both files by delegating to SHAs that org-infra actively maintains, and aligns this repo with the org-standard pattern already used in `unbound-force/unbound-force` and `complytime/.github`.

## What Changes

- Replace the inline `osv-scanner` and `scorecards` jobs in `ci_security.yml` with calls to `reusable_vuln_scan.yml` and `reusable_security.yml` from `complytime/org-infra` — this also switches from the old inline step form to the OSV-maintained reusable workflow caller form (v2.5.1).
- Replace the inline `osv-scanner` and `scorecards` jobs in `ci_scheduled.yml` with a single call to `reusable_scheduled.yml` from `complytime/org-infra`, which handles the scheduled OSV + Scorecard runs together (matching the `complytime/.github` pattern).
- Both files pin org-infra at the same SHA already used by `unbound-force/unbound-force` (`0c784711... # v0.7.1`).

## Capabilities

### New Capabilities

- `ci-security`: The security CI workflows (`ci_security.yml` and `ci_scheduled.yml`) must delegate OSV scanning and OpenSSF Scorecard analysis to org-infra reusable workflows rather than running inline steps with self-managed action SHAs.

### Modified Capabilities

*(none — no existing specs)*

## Impact

- `.github/workflows/ci_security.yml` — inline OSV + Scorecard jobs replaced.
- `.github/workflows/ci_scheduled.yml` — inline OSV + Scorecard jobs replaced.
- No other workflows, configs, or source files are affected.
- The `OSV-Scanner` and `OpenSSF Scorecards` check names visible in GitHub CI remain the same (defined by the reusable workflows).
- Unblocks PR #40 (`issue-39-enrich-org-labels`) which currently has two hard CI failures from `ci_security.yml`; also prevents the same failures from hitting the nightly cron via `ci_scheduled.yml`.
