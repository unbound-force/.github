# Tasks

## 1. Upgrade Safe Settings Full Sync

- [x] 1.1 Pin `.github/workflows/safe_settings_sync.yml` to Safe Settings commit `6a8b6ae084987025f6c5de85e3cc6df140f64502` and verify a regression test asserts the exact immutable SHA.
- [x] 1.2 Adapt the workflow's existing custom full-sync runner to wait for Probot 14 readiness before logger access and application loading; verify the focused runner/workflow test passes.
- [x] 1.3 Preserve the exact predicate `String(error).includes("Cannot read properties of undefined (reading 'check_suite')")`, add tests proving that exact failure remains non-fatal while similar and unrelated exceptions/settings errors fail, and verify those tests pass.

## 2. Preserve Labels

- [x] 2.1 Quote only the `ci` label color as `'5319e7'` and verify `make safe-settings-validate` passes.
- [x] 2.2 Convert label configuration to additive include/exclude behavior that preserves undeclared labels, add focused regression tests for preservation and the `ci` color string, and verify `make test-unit` passes.

## 3. Validate The Change

- [x] 3.1 Verify all six specification scenarios map to focused local tests, then run those tests, `make safe-settings-validate`, `make test-unit`, and `openspec validate fix-safe-settings-full-sync --strict`; verify every command and the existing coverage ratchet pass.

<!-- spec-review: passed -->

<!-- code-review: passed -->
