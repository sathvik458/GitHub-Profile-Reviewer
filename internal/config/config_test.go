package config

import (
	"os"
	"path/filepath"
	"testing"
)

// writeEnvFile puts content in a throwaway .env and returns its path.
// t.TempDir cleans the directory up automatically.
func writeEnvFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	return path
}

// unsetEnv removes key for the duration of the test. t.Setenv registers the
// restore hook, then Unsetenv clears the value so loadEnvFile sees a variable
// that is genuinely absent rather than empty.
func unsetEnv(t *testing.T, key string) {
	t.Helper()

	t.Setenv(key, "")

	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("unset %s: %v", key, err)
	}
}

func TestLoadEnvFileSetsValues(t *testing.T) {
	unsetEnv(t, "GPR_TEST_TOKEN")

	if err := loadEnvFile(writeEnvFile(t, "GPR_TEST_TOKEN=from-file\n")); err != nil {
		t.Fatalf("loadEnvFile() error = %v", err)
	}

	if got := os.Getenv("GPR_TEST_TOKEN"); got != "from-file" {
		t.Fatalf("GPR_TEST_TOKEN = %q, want %q", got, "from-file")
	}
}

// The precedence rule the whole package exists for: a real environment
// variable always beats the file, so containers and CI stay in control.
func TestLoadEnvFileDoesNotOverrideRealEnvironment(t *testing.T) {
	t.Setenv("GPR_TEST_TOKEN", "from-environment")

	if err := loadEnvFile(writeEnvFile(t, "GPR_TEST_TOKEN=from-file\n")); err != nil {
		t.Fatalf("loadEnvFile() error = %v", err)
	}

	if got := os.Getenv("GPR_TEST_TOKEN"); got != "from-environment" {
		t.Fatalf("GPR_TEST_TOKEN = %q, want the real environment value to win", got)
	}
}

func TestLoadEnvFileSkipsCommentsAndBlankLines(t *testing.T) {
	unsetEnv(t, "GPR_TEST_VALUE")
	unsetEnv(t, "GPR_TEST_COMMENTED")

	content := "# GPR_TEST_COMMENTED=should-not-be-set\n\n   \nGPR_TEST_VALUE=set\n"

	if err := loadEnvFile(writeEnvFile(t, content)); err != nil {
		t.Fatalf("loadEnvFile() error = %v", err)
	}

	if _, ok := os.LookupEnv("GPR_TEST_COMMENTED"); ok {
		t.Fatal("commented line was applied")
	}

	if got := os.Getenv("GPR_TEST_VALUE"); got != "set" {
		t.Fatalf("GPR_TEST_VALUE = %q, want %q", got, "set")
	}
}

func TestLoadEnvFileIgnoresMalformedLines(t *testing.T) {
	unsetEnv(t, "GPR_TEST_VALUE")

	content := "LINE_WITH_NO_SEPARATOR\n=value-without-a-key\nGPR_TEST_VALUE=set\n"

	if err := loadEnvFile(writeEnvFile(t, content)); err != nil {
		t.Fatalf("loadEnvFile() error = %v", err)
	}

	if _, ok := os.LookupEnv("LINE_WITH_NO_SEPARATOR"); ok {
		t.Fatal("line without a separator was applied")
	}

	if got := os.Getenv("GPR_TEST_VALUE"); got != "set" {
		t.Fatalf("GPR_TEST_VALUE = %q, want parsing to continue past malformed lines", got)
	}
}

func TestLoadEnvFileTrimsWhitespace(t *testing.T) {
	unsetEnv(t, "GPR_TEST_VALUE")

	if err := loadEnvFile(writeEnvFile(t, "  GPR_TEST_VALUE  =  spaced  \n")); err != nil {
		t.Fatalf("loadEnvFile() error = %v", err)
	}

	if got := os.Getenv("GPR_TEST_VALUE"); got != "spaced" {
		t.Fatalf("GPR_TEST_VALUE = %q, want %q", got, "spaced")
	}
}

// A missing .env is the normal case in a container, so it must not be an error.
func TestLoadEnvFileMissingFileIsNotAnError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.env")

	if err := loadEnvFile(path); err != nil {
		t.Fatalf("loadEnvFile() error = %v, want nil for a missing file", err)
	}
}

func TestLoadReadsGitHubTokenFromEnvironment(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "token-from-environment")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.GitHubToken != "token-from-environment" {
		t.Fatalf("GitHubToken = %q, want %q", cfg.GitHubToken, "token-from-environment")
	}
}
