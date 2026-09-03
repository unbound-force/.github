## 1. Confirm org-infra reusable references

- [x] 1.1 Verify `reusable_deps_reviewer.yml` and `reusable_dependabot_reviewer.yml` exist in `complytime/org-infra` at SHA `0c784711926c9864f027ec565fd7c06a382d80f8` (v0.7.1). Verify by inspecting the org-infra repo at that ref; if either file is absent at that SHA, identify the earliest SHA that contains both and note it for use in tasks 2.x.
- [x] 1.2 If (and only if) the SHA is bumped from v0.7.1 in 1.1, also update the `uses:` pin in `ci_security.yml` and `ci_scheduled.yml` to the same new SHA so all org-infra callers in this repo stay uniform (design.md D2). Verify by grepping for `complytime/org-infra` and confirming every hit uses the identical SHA. If the SHA remains v0.7.1, this task is a no-op — mark it complete.

## 2. Migrate ci_dependencies.yml to reusable callers

- [x] 2.1 Remove the inline `dependency-review` job from `.github/workflows/ci_dependencies.yml` (both the `actions/checkout` and `actions/dependency-review-action` steps). Verify no `uses: actions/dependency-review-action` line remains (`grep` returns nothing).
- [x] 2.2 Add a `call_deps_reviewer` job (name: `General`) calling `complytime/org-infra/.github/workflows/reusable_deps_reviewer.yml` pinned to the confirmed SHA with a version comment. Ensure the inline version comment (`# vX.Y.Z`) matches the actual release tag of the chosen SHA. Verify the `uses:` line includes a full 40-char SHA and a matching version comment.
- [x] 2.3 Add a `call_dependabot_reviewer` job (name: `Dependabot`) calling `complytime/org-infra/.github/workflows/reusable_dependabot_reviewer.yml` pinned to the same SHA with a matching version comment. Verify no auto-approve/auto-merge job is added anywhere in the file.
- [x] 2.4 Preserve the workflow trigger (`pull_request` to `main`) and set a restrictive top-level `permissions:` block — `contents: read`, `issues: none`, `pull-requests: none` — matching the canonical caller and the reusables' own passive permissions (verified: both reusables declare `pull-requests: none` at the pinned SHA; no write scope is needed for the two included callers). Verify the file's `on:` and `permissions:` blocks are correct for the callers.
- [x] 2.5 Add an inline comment in `ci_dependencies.yml` above the caller jobs noting that the dependency review is soft-gated (`continue-on-error: true` is applied inside the reusable) and that this gate should be re-evaluated if the Dependency Graph / GitHub Advanced Security is ever enabled on this repo (design.md Risk 1). Verify the comment is present.

## 3. Add Dependabot configuration

- [x] 3.1 Create `.github/dependabot.yml` with a `version: 2` config and one `updates` entry for `package-ecosystem: "github-actions"`, `directory: "/"`, and a weekly `schedule`. Verify the file contains `package-ecosystem: "github-actions"` and a defined schedule.

## 4. Validate and verify

- [x] 4.1 Run `make sanity` (yamllint + Go tests) and confirm it passes, including YAML validation of the two changed/new files. (Verified equivalents: `make lint`, `go vet ./...`, `go build ./...`, `go test ./...`, and `yamllint` on both changed files all pass with exit 0. Note: literal `make sanity` also runs `git diff --exit-code`, which is expected to fail pre-commit while changes are uncommitted; that gate is satisfied by `/uf.finale` at commit time.)
- [x] 4.2 Open a PR against `main` and confirm the Dependency Review check no longer hard-fails (the `General` caller job runs without blocking the PR), satisfying the `ci-dependencies` spec scenarios. Note that the `call_dependabot_reviewer` job is conditional on the PR author being Dependabot and will NOT fire on a human-authored validation PR; verify its correctness by inspecting the job's `if:` condition in the merged YAML rather than expecting it to run on the test PR. (DEFERRED to PR time / `/uf.finale`: PR creation is not available in this environment. Static verification done in lieu: the caller `uses:` targets `reusable_deps_reviewer.yml`, which was confirmed at the pinned SHA to wrap `dependency-review-action` with `continue-on-error: true` — so the hard-fail is structurally resolved; and `reusable_dependabot_reviewer.yml`'s inner job is gated on `github.event.pull_request.user.login == 'dependabot[bot]'`, confirmed at the pinned SHA.)

<!-- spec-review: passed -->
<!-- code-review: passed -->
