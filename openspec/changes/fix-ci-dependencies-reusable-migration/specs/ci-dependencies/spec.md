## Purpose

Defines the dependency CI workflow's obligations for dependency review on pull requests — delegating execution to org-infra reusable workflows rather than running an inline `dependency-review-action` step — and the repository's obligation to declare a `github-actions` Dependabot configuration for automated action SHA updates.

## ADDED Requirements

### Requirement: Dependency review via org-infra reusable workflow
The dependency CI workflow SHALL delegate dependency review to `complytime/org-infra/.github/workflows/reusable_deps_reviewer.yml` rather than running an inline `actions/dependency-review-action` step.

#### Scenario: Dependency review runs on pull_request to main
- **WHEN** a pull request targets the `main` branch
- **THEN** the `reusable_deps_reviewer.yml` caller job is triggered and its Dependency Review check is reported in CI

#### Scenario: Missing Dependency Graph no longer blocks CI
- **WHEN** a pull request is opened against `main` and the Dependency Graph feature is not enabled on the repository
- **THEN** the dependency CI workflow SHALL NOT hard-fail the pull request
- **AND** the dependency review outcome is surfaced informally rather than as a blocking gate

#### Scenario: Inline dependency-review-action step removed
- **WHEN** `ci_dependencies.yml` is read
- **THEN** it SHALL NOT contain an inline `uses: actions/dependency-review-action` step

### Requirement: Dependabot PR review via org-infra reusable workflow
The dependency CI workflow SHALL delegate review of Dependabot-authored pull requests to `complytime/org-infra/.github/workflows/reusable_dependabot_reviewer.yml`.

#### Scenario: Dependabot reviewer caller present
- **WHEN** `ci_dependencies.yml` is read
- **THEN** it SHALL contain a caller job for `reusable_dependabot_reviewer.yml`

### Requirement: Dependabot auto-approve excluded from this repository
The dependency CI workflow SHALL NOT enable Dependabot auto-approve or auto-merge behavior, because this repository governs organization policy (Peribolos membership, safe-settings, rulesets) and auto-approval represents an elevated attack surface pending a separate threat-model review.

#### Scenario: No auto-approve job configured
- **WHEN** `ci_dependencies.yml` is read
- **THEN** it SHALL NOT contain a job that automatically approves or auto-merges Dependabot pull requests

### Requirement: Caller workflows pin org-infra reusable by SHA
The `uses:` reference to each org-infra reusable workflow in `ci_dependencies.yml` SHALL be pinned to a full commit SHA with an inline version comment, following the same pin-by-SHA convention used throughout this repository.

#### Scenario: SHA pin present for each caller
- **WHEN** `ci_dependencies.yml` is read
- **THEN** each `uses: complytime/org-infra/...` line includes a full 40-character SHA and an inline version comment (e.g., `# v0.7.1`)

### Requirement: Dependabot configuration for the github-actions ecosystem
The repository SHALL declare a `.github/dependabot.yml` configuration covering the `github-actions` package ecosystem so that action SHA updates are proposed automatically. As part of this change, no ecosystem entries beyond `github-actions` SHALL be added (the `go.mod` tooling manifest is deliberately out of scope).

#### Scenario: dependabot.yml declares github-actions ecosystem
- **WHEN** `.github/dependabot.yml` is read
- **THEN** it SHALL contain an entry with `package-ecosystem: "github-actions"` and a defined update schedule
- **AND** it SHALL NOT contain any `package-ecosystem` entry other than `github-actions`
