package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Skpow1234/Katame/internal/policy"
)

func TestBatchProtect_HonorsPolicy(t *testing.T) {
	srcDir := filepath.Join(t.TempDir(), "src")
	encDir := filepath.Join(t.TempDir(), "enc")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("alpha"), 0o600); err != nil {
		t.Fatal(err)
	}

	policyFile := writeYAMLPolicy(t, "deny-batch.yaml", `
version: 1
default: allow
rules:
  - name: no-batch-protect
    action: deny
    reason: "batch-protect disabled for tests"
    when:
      operation: [batch-protect]
`)
	t.Cleanup(func() { policy.SetGlobal(nil) })

	root := NewRootCmd()
	root.SetArgs([]string{
		"--policy", policyFile,
		"batch-protect",
		"--dir", srcDir,
		"--out-dir", encDir,
	})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected batch-protect to be denied by policy")
	}
	if !strings.Contains(err.Error(), "policy denied") {
		t.Fatalf("expected policy denied, got: %v", err)
	}
}
