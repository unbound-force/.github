## 1. Rewrite ci_security.yml

- [x] 1.1 Replace the inline `osv-scanner` job in `.github/workflows/ci_security.yml` with a `call_reusable_vuln_scan` job that calls `complytime/org-infra/.github/workflows/reusable_vuln_scan.yml@0c784711926c9864f027ec565fd7c06a382d80f8 # v0.7.1` with permissions `contents: read`, `actions: read`, `security-events: write`, `packages: write`, `id-token: write`. Verify the job block matches the `unbound-force/unbound-force` canonical.

- [x] 1.2 Replace the inline `scorecards` job with a `call_reusable_security` job that calls `complytime/org-infra/.github/workflows/reusable_security.yml@0c784711926c9864f027ec565fd7c06a382d80f8 # v0.7.1` with permissions `contents: read`, `id-token: write`, `security-events: write`. Verify the job block matches the canonical.

- [x] 1.3 Update the top-level `permissions` block to add `packages: none` alongside the existing `actions: none`, `id-token: none`, `security-events: none` (mirroring `complytime/.github`). Verify no stale permission fields remain.

## 2. Rewrite ci_scheduled.yml

- [x] 2.1 Replace both inline jobs (`osv-scanner` and `scorecards`) in `.github/workflows/ci_scheduled.yml` with a single `call_reusable_scheduled` job that calls `complytime/org-infra/.github/workflows/reusable_scheduled.yml@0c784711926c9864f027ec565fd7c06a382d80f8 # v0.7.1` with permissions `contents: read`, `actions: read`, `security-events: write`, `id-token: write`. Verify the result matches the `complytime/.github` canonical.

## 3. Verify locally

- [x] 3.1 Run `make lint` (or `yamllint .github/workflows/ci_security.yml .github/workflows/ci_scheduled.yml`) and confirm zero lint errors on both edited files.

## 4. Push and confirm CI

- [ ] 4.1 Push the branch. In PR #40, confirm that both `OSV-Scanner` and `OpenSSF Scorecards` checks transition from FAILURE to SUCCESS (or neutral/skipped if the repo lacks OIDC publish rights for Scorecards). Verify the `Dependency Review` check is the only remaining failure (addressed separately by issue #41).

<!-- spec-review: passed -->
