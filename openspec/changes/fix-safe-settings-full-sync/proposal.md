# Proposal

## Why

Safe Settings 2.1.18 has an upstream rename race that can mix repository configuration during full sync. The upgrade also requires a small Probot 14 compatibility change and correction of destructive label behavior observed in dry runs.

## What Changes

- Upgrade Safe Settings to a reviewed immutable commit containing the rename-race fix.
- Adapt the custom full-sync runner in `.github/workflows/safe_settings_sync.yml` for Probot 14 initialization.
- Preserve and validate the existing `check_suite` full-sync workaround.
- Quote the `ci` label color as `'5319e7'`.
- Preserve labels that are not explicitly managed by Safe Settings.

## Capabilities

### New Capabilities

- `safe-settings-full-sync`: Safe Settings full-sync compatibility and non-destructive label configuration.

### Modified Capabilities

None.

## Impact

- `.github/workflows/safe_settings_sync.yml`
- `safe-settings/settings.yml` and repository label overrides
- Focused workflow and Safe Settings configuration tests
- Pinned `github/safe-settings` revision
