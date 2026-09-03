## Purpose

Defines the security CI workflows' obligations for OSV vulnerability scanning and OpenSSF Scorecard analysis — both on PRs (`ci_security.yml`) and on schedule (`ci_scheduled.yml`) — delegating execution to org-infra reusable workflows rather than managing action SHAs inline.

## ADDED Requirements

### Requirement: OSV scanning via org-infra reusable workflow
The security CI workflow SHALL delegate OSV dependency scanning to `complytime/org-infra/.github/workflows/reusable_vuln_scan.yml` rather than running an inline `google/osv-scanner-action` step.

#### Scenario: OSV scan runs on pull_request to main
- **WHEN** a pull request targets the `main` branch
- **THEN** the `reusable_vuln_scan.yml` caller job is triggered and the OSV-Scanner check is reported in CI

#### Scenario: OSV scan runs on push to main
- **WHEN** a commit is pushed directly to `main`
- **THEN** the `reusable_vuln_scan.yml` caller job is triggered and the OSV-Scanner check is reported in CI

#### Scenario: Dead OSV SHA no longer blocks CI
- **WHEN** a PR is opened against `main`
- **THEN** the CI workflow SHALL NOT fail with "unable to find version" for `google/osv-scanner-action`

### Requirement: OpenSSF Scorecard analysis via org-infra reusable workflow
The security CI workflow SHALL delegate OpenSSF Scorecard analysis to `complytime/org-infra/.github/workflows/reusable_security.yml` rather than running an inline `ossf/scorecard-action` step.

#### Scenario: Scorecard runs on push to main
- **WHEN** a commit is pushed directly to `main`
- **THEN** the `reusable_security.yml` caller job is triggered and the OpenSSF Scorecards check is reported in CI

#### Scenario: Scorecard runs on schedule
- **WHEN** the scheduled cron trigger fires
- **THEN** the `reusable_security.yml` caller job is triggered

#### Scenario: Dead Scorecard SHA no longer blocks CI
- **WHEN** a PR is opened against `main`
- **THEN** the CI workflow SHALL NOT fail with "unable to find version" for `ossf/scorecard-action`

### Requirement: Scheduled OSV and Scorecard runs via org-infra reusable workflow
The scheduled CI workflow SHALL delegate the nightly OSV scan and Scorecard analysis to `complytime/org-infra/.github/workflows/reusable_scheduled.yml` rather than running inline steps.

#### Scenario: Scheduled workflow runs on cron trigger
- **WHEN** the nightly cron trigger fires
- **THEN** the `reusable_scheduled.yml` caller job is triggered and both OSV-Scanner and OpenSSF Scorecards checks are executed

#### Scenario: Dead SHAs no longer block scheduled runs
- **WHEN** the scheduled cron trigger fires
- **THEN** the workflow SHALL NOT fail with "unable to find version" for either `google/osv-scanner-action` or `ossf/scorecard-action`

### Requirement: Caller workflows pin org-infra reusable by SHA
The `uses:` reference to each org-infra reusable workflow in both `ci_security.yml` and `ci_scheduled.yml` SHALL be pinned to a full commit SHA with a version comment, following the same pin-by-SHA convention used throughout this repo.

#### Scenario: SHA pin present in both workflow files
- **WHEN** `ci_security.yml` or `ci_scheduled.yml` is read
- **THEN** each `uses: complytime/org-infra/...` line includes a full 40-character SHA and an inline version comment (e.g., `# v1.2.3`)
