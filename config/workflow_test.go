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
	workflow := readSafeSettingsWorkflow(t)

	for _, required := range []string{
		"allowed_repos",
		`^[A-Za-z0-9][A-Za-z0-9._-]*$`,
		"Repository is not managed by Safe Settings",
		"Repository was specified more than once",
	} {
		if !strings.Contains(workflow, required) {
			t.Errorf("Safe Settings workflow must validate scoped repositories with %q", required)
		}
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

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create test directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}
}
