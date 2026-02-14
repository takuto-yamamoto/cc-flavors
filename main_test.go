package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

var testBinary string

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "cc-flavors-test")
	if err != nil {
		panic(err)
	}
	testBinary = filepath.Join(tmpDir, "cc-flavors")
	cmd := exec.Command("go", "build", "-o", testBinary, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	code := m.Run()
	if err := os.RemoveAll(tmpDir); err != nil {
		panic(err)
	}
	os.Exit(code)
}

func runCLI(t *testing.T, input string, args ...string) (string, string) {
	t.Helper()

	cmd := exec.Command(testBinary, args...)
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v (stderr: %s)", err, stderr.String())
	}
	return stdout.String(), stderr.String()
}

func runCLIExpectError(t *testing.T, input string, args ...string) (string, string) {
	t.Helper()

	cmd := exec.Command(testBinary, args...)
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err == nil {
		t.Fatalf("expected command failure, got success")
	}
	return stdout.String(), stderr.String()
}

func assertBoxTable(t *testing.T, output string) {
	t.Helper()
	lines := []string{}
	for _, line := range strings.Split(output, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) < 3 {
		t.Fatalf("expected table output, got:\n%s", output)
	}
	width := utf8.RuneCountInString(lines[0])
	for _, line := range lines {
		if utf8.RuneCountInString(line) != width {
			t.Fatalf("table lines have inconsistent widths:\n%s", output)
		}
	}
	if !strings.HasPrefix(lines[0], "╭") || !strings.HasSuffix(lines[0], "╮") {
		t.Fatalf("missing top border:\n%s", output)
	}
}

func TestSummaryEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "events.sqlite")

	stdout, _ := runCLI(t, "", "summary", "--db", dbPath)
	if stdout != "No flavor texts found yet.\n" {
		t.Fatalf("unexpected output: %q", stdout)
	}
}

func TestDefaultRunsSummary(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "events.sqlite")

	stdout, _ := runCLI(t, "", "--db", dbPath)
	if stdout != "No flavor texts found yet.\n" {
		t.Fatalf("unexpected output: %q", stdout)
	}
}

func TestIngestAndSummary(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "events.sqlite")

	input := strings.Join([]string{
		"Thinking... Moonwalking...",
		"Refactoring... Thinking...",
	}, "\n")
	_, _ = runCLI(t, input, "ingest", "--db", dbPath)

	stdout, _ := runCLI(t, "", "summary", "--db", dbPath)
	if !strings.Contains(stdout, "Your Flavors") {
		t.Fatalf("expected title in output, got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "Thinking") || !strings.Contains(stdout, "Moonwalking") || !strings.Contains(stdout, "Refactoring") {
		t.Fatalf("expected flavors in output, got:\n%s", stdout)
	}
	assertBoxTable(t, stdout)
}

func TestSummarySinceFilters(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "events.sqlite")

	input := strings.Join([]string{
		"Thinking...",
		"Thinking...",
	}, "\n")
	_, _ = runCLI(t, input, "ingest", "--db", dbPath)

	stdout, _ := runCLI(t, "", "summary", "--db", dbPath, "--since", "2099-01-01T00:00:00Z")
	if stdout != "No flavor texts found yet.\n" {
		t.Fatalf("unexpected output: %q", stdout)
	}
}

func TestSummarySinceValidation(t *testing.T) {
	_, stderr := runCLIExpectError(t, "", "summary", "--since", "not-a-time")
	if !strings.Contains(stderr, "invalid --since") {
		t.Fatalf("expected since validation error, got: %s", stderr)
	}
}

func TestClearConfirmation(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "events.sqlite")

	_, _ = runCLI(t, "Thinking...\n", "ingest", "--db", dbPath)
	_, _ = runCLI(t, "n\n", "clear", "--db", dbPath)
	stdout, _ := runCLI(t, "", "summary", "--db", dbPath)
	if !strings.Contains(stdout, "Thinking") {
		t.Fatalf("expected data to remain, got: %s", stdout)
	}

	_, _ = runCLI(t, "", "clear", "--db", dbPath, "--yes")
	stdout, _ = runCLI(t, "", "summary", "--db", dbPath)
	if stdout != "No flavor texts found yet.\n" {
		t.Fatalf("expected empty after clear, got: %s", stdout)
	}
}
