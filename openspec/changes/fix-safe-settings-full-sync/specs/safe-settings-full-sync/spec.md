# Spec Delta

## Purpose

Defines the Safe Settings full-sync dependency, runtime compatibility, reporting workaround, and non-destructive label behavior.

## ADDED Requirements

### Requirement: Safe Settings uses the fixed immutable revision
The workflow SHALL use Safe Settings commit `6a8b6ae084987025f6c5de85e3cc6df140f64502`, which contains the upstream rename-race fix.

#### Scenario: Full sync starts
- **WHEN** the workflow checks out Safe Settings
- **THEN** it uses the reviewed immutable commit

### Requirement: The full-sync runner supports Probot 14
The custom full-sync runner SHALL wait for Probot initialization before accessing the logger or loading Safe Settings.

#### Scenario: Probot initializes asynchronously
- **WHEN** full sync starts with Safe Settings 2.1.19
- **THEN** the runner waits for Probot readiness before using the logger or application

### Requirement: The check-suite workaround remains narrow
The full-sync runner SHALL preserve the existing predicate `String(error).includes("Cannot read properties of undefined (reading 'check_suite')")` as the only non-fatal exception match and SHALL fail on every other exception or settings error.

#### Scenario: Known reporting error occurs
- **WHEN** synchronization completes and the known missing `check_suite` context error occurs
- **THEN** the runner records the limitation and completes successfully

#### Scenario: Another error occurs
- **WHEN** the runner receives any different exception or settings error
- **THEN** full sync fails

### Requirement: Label configuration is non-destructive
Safe Settings SHALL preserve labels that are not explicitly declared for management, and the `ci` label color SHALL be represented as the string `5319e7`.

#### Scenario: Repository has an undeclared label
- **WHEN** Safe Settings evaluates managed labels
- **THEN** the undeclared label is not deleted or modified

#### Scenario: CI label is loaded
- **WHEN** Safe Settings loads the `ci` label
- **THEN** its color is the string `5319e7`
