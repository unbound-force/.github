// SPDX-License-Identifier: Apache-2.0

package config

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func readSafeSettingsWorkflow(t *testing.T) string {
	t.Helper()

	workflow, err := os.ReadFile("../.github/workflows/safe_settings_sync.yml")
	if err != nil {
		t.Fatalf("read Safe Settings workflow: %v", err)
	}
	return string(workflow)
}

func TestSafeSettingsWorkflow_PinsReviewedCommit(t *testing.T) {
	workflow := readSafeSettingsWorkflow(t)

	const pinnedRevision = "ref: '6a8b6ae084987025f6c5de85e3cc6df140f64502'"
	if !strings.Contains(workflow, pinnedRevision) {
		t.Errorf("Safe Settings workflow must contain %q", pinnedRevision)
	}
	if strings.Contains(workflow, "inputs.version") {
		t.Error("Safe Settings workflow must not allow overriding the reviewed revision")
	}
}

func TestSafeSettingsWorkflow_ValidatesScopedRepositories(t *testing.T) {
	scopedSync := extractScopedSync(t, readSafeSettingsWorkflow(t))

	testCases := []struct {
		name             string
		targetRepos      string
		wantExitCode     int
		wantConfig       string
		wantOutput       string
		wantConfigExists bool
	}{
		{
			name:             "managed repository succeeds",
			targetRepos:      "dewey",
			wantConfig:       "    - dewey\n",
			wantConfigExists: true,
		},
		{
			name:         "invalid repository fails",
			targetRepos:  "*",
			wantExitCode: 1,
			wantOutput:   "Invalid repository name: *",
		},
		{
			name:         "unmanaged repository fails",
			targetRepos:  "unmanaged",
			wantExitCode: 1,
			wantOutput:   "Repository is not managed by Safe Settings: unmanaged",
		},
		{
			name:         "duplicate repository fails",
			targetRepos:  "dewey,dewey",
			wantExitCode: 1,
			wantOutput:   "Repository was specified more than once: dewey",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			exitCode, output, configuration, configurationExists := runScopedSync(t, scopedSync, testCase.targetRepos)
			if exitCode != testCase.wantExitCode {
				t.Errorf("scoped sync exit code = %d, want %d; output: %s", exitCode, testCase.wantExitCode, output)
			}
			if testCase.wantOutput != "" && !strings.Contains(output, testCase.wantOutput) {
				t.Errorf("scoped sync output must contain %q; got %q", testCase.wantOutput, output)
			}
			if configurationExists != testCase.wantConfigExists {
				t.Errorf("scoped configuration exists = %t, want %t", configurationExists, testCase.wantConfigExists)
			}
			if testCase.wantConfig != "" && !strings.Contains(configuration, testCase.wantConfig) {
				t.Errorf("scoped configuration must contain %q; got %q", testCase.wantConfig, configuration)
			}
		})
	}
}

func TestSafeSettingsWorkflow_AwaitsProbotReadinessBeforeUse(t *testing.T) {
	workflowText := readSafeSettingsWorkflow(t)
	readiness := strings.Index(workflowText, "await probot.ready()")
	loggerAccess := strings.Index(workflowText, "probot.log.info(`Starting full sync")
	applicationLoad := strings.Index(workflowText, "const app = appFn(probot, {})")
	if readiness == -1 || loggerAccess == -1 || applicationLoad == -1 {
		t.Fatalf("Safe Settings runner must await readiness and then access the logger and load the application")
	}
	if readiness > loggerAccess || readiness > applicationLoad {
		t.Errorf("Safe Settings runner must await Probot readiness before logger access and application loading")
	}
}

func TestSafeSettingsWorkflow_PreservesExactCheckSuitePredicate(t *testing.T) {
	workflow := readSafeSettingsWorkflow(t)
	const predicate = `String(error).includes("Cannot read properties of undefined (reading 'check_suite')")`
	if strings.Count(workflow, predicate) != 1 {
		t.Errorf("Safe Settings runner must contain exactly one %q predicate", predicate)
	}
}

func TestSafeSettingsWorkflow_CheckSuiteFailureHandling(t *testing.T) {
	runner := extractFullSyncRunner(t, readSafeSettingsWorkflow(t))
	testCases := []struct {
		name          string
		exception     string
		settingsError bool
		wantExitCode  int
		wantOutput    string
	}{
		{
			name:         "exact check_suite exception is non-fatal",
			exception:    "Cannot read properties of undefined (reading 'check_suite')",
			wantExitCode: 0,
			wantOutput:   "check_suite reporting skipped (expected in full-sync mode)",
		},
		{
			name:         "similar check_run exception is fatal",
			exception:    "Cannot read properties of undefined (reading 'check_run')",
			wantExitCode: 1,
			wantOutput:   "Unexpected error during full sync",
		},
		{
			name:         "unrelated exception is fatal",
			exception:    "network failure",
			wantExitCode: 1,
			wantOutput:   "Unexpected error during full sync",
		},
		{
			name:          "settings error is fatal",
			settingsError: true,
			wantExitCode:  1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			exitCode, output := runFullSyncRunner(t, runner, testCase.exception, testCase.settingsError)
			if exitCode != testCase.wantExitCode {
				t.Errorf("runner exit code = %d, want %d; output: %s", exitCode, testCase.wantExitCode, output)
			}
			if testCase.wantOutput != "" && !strings.Contains(output, testCase.wantOutput) {
				t.Errorf("runner output must contain %q; got %q", testCase.wantOutput, output)
			}
		})
	}
}

func extractFullSyncRunner(t *testing.T, workflow string) string {
	t.Helper()

	const startMarker = "          cat > full-sync-patched.js << 'ENDPATCH'\n"
	const endMarker = "\n          ENDPATCH"
	start := strings.Index(workflow, startMarker)
	if start == -1 {
		t.Fatal("Safe Settings workflow must define the patched full-sync runner")
	}
	start += len(startMarker)
	end := strings.Index(workflow[start:], endMarker)
	if end == -1 {
		t.Fatal("Safe Settings workflow must terminate the patched full-sync runner")
	}

	lines := strings.Split(workflow[start:start+end], "\n")
	for index := range lines {
		lines[index] = strings.TrimPrefix(lines[index], "          ")
	}
	return strings.Join(lines, "\n")
}

func extractScopedSync(t *testing.T, workflow string) string {
	t.Helper()

	const startMarker = "          echo \"Scoping safe-settings to repos: $TARGET_REPOS\"\n"
	const endMarker = "\n\n      - name: Checkout safe-settings code"
	start := strings.Index(workflow, startMarker)
	if start == -1 {
		t.Fatal("Safe Settings workflow must define scoped sync generation")
	}
	start += len(startMarker)
	end := strings.Index(workflow[start:], endMarker)
	if end == -1 {
		t.Fatal("Safe Settings workflow must terminate scoped sync generation")
	}

	lines := strings.Split(workflow[start:start+end], "\n")
	for index := range lines {
		lines[index] = strings.TrimPrefix(lines[index], "          ")
	}
	return strings.Join(lines, "\n")
}

func runFullSyncRunner(t *testing.T, runner, exception string, settingsError bool) (int, string) {
	t.Helper()

	testDirectory := t.TempDir()
	writeTestFile(t, filepath.Join(testDirectory, "full-sync-patched.js"), runner)
	writeTestFile(t, filepath.Join(testDirectory, "index.js"), `module.exports = () => ({
  syncInstallation: async () => {
    if (process.env.RUNNER_EXCEPTION) throw new Error(process.env.RUNNER_EXCEPTION)
    return { errors: process.env.SETTINGS_ERROR === 'true' ? ['settings failure'] : [] }
  }
})
`)
	writeTestFile(t, filepath.Join(testDirectory, "lib", "env.js"), "exports.FULL_SYNC_NOP = true\n")
	writeTestFile(t, filepath.Join(testDirectory, "node_modules", "probot", "index.js"), `exports.createProbot = () => ({
  ready: async () => {},
  log: { info: (message) => process.stdout.write(message + '\n'), error: () => {} }
})
`)

	command := exec.Command("node", "full-sync-patched.js")
	command.Dir = testDirectory
	command.Env = append(os.Environ(),
		"RUNNER_EXCEPTION="+exception,
		"SETTINGS_ERROR="+map[bool]string{true: "true", false: "false"}[settingsError],
	)
	output, err := command.CombinedOutput()
	if err == nil {
		return 0, string(output)
	}
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		t.Fatalf("run patched full-sync runner: %v", err)
	}
	return exitError.ExitCode(), string(output)
}

func runScopedSync(t *testing.T, scopedSync, targetRepos string) (int, string, string, bool) {
	t.Helper()

	testDirectory := t.TempDir()
	writeTestFile(t, filepath.Join(testDirectory, "safe-settings", "deployment-settings.yml"), "restrictedRepos: {}\n")
	writeTestFile(t, filepath.Join(testDirectory, "bin", "yq"), `#!/usr/bin/env bash
if [[ "$1" == "-r" && "$2" == ".restrictedRepos.include[]" ]]; then
  printf '%s\n' dewey website
else
  printf 'configvalidators: []\n'
fi
`)
	if err := os.Chmod(filepath.Join(testDirectory, "bin", "yq"), 0o755); err != nil {
		t.Fatalf("make mock yq executable: %v", err)
	}

	command := exec.Command("bash", "-c", scopedSync)
	command.Dir = testDirectory
	command.Env = append(os.Environ(),
		"TARGET_REPOS="+targetRepos,
		"RUNNER_TEMP="+testDirectory,
		"PATH="+filepath.Join(testDirectory, "bin")+":"+os.Getenv("PATH"),
	)
	output, err := command.CombinedOutput()
	configurationPath := filepath.Join(testDirectory, "scoped-deployment-settings.yml")
	configuration, readErr := os.ReadFile(configurationPath)
	configurationExists := readErr == nil
	if readErr != nil && !os.IsNotExist(readErr) {
		t.Fatalf("read scoped deployment settings: %v", readErr)
	}
	if err == nil {
		return 0, string(output), string(configuration), configurationExists
	}
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		t.Fatalf("run scoped sync: %v", err)
	}
	return exitError.ExitCode(), string(output), string(configuration), configurationExists
}

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create test directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}
}
