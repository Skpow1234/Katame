package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Skpow1234/Katame/internal/util"
)

func TestHash_BadAlgo_ExitCode(t *testing.T) {
	if os.Getenv("KATAME_EXITCODE_HELPER") == "1" {
		dir := os.TempDir()
		in := filepath.Join(dir, "katame-hash-exitcode.txt")
		_ = os.WriteFile(in, []byte("x"), 0o600)
		root := NewRootCmd()
		root.SetArgs([]string{"hash", "--in", in, "--algo", "not-a-hash"})
		err := root.Execute()
		if err == nil {
			os.Exit(0)
		}
		os.Exit(util.ExitCodeForError(err))
	}

	if testing.Short() {
		t.Skip("subprocess")
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestHash_BadAlgo_ExitCode$", "-test.v=false")
	cmd.Env = append(os.Environ(), "KATAME_EXITCODE_HELPER=1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	ee, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("want ExitError, got %T %v", err, err)
	}
	if ee.ExitCode() != util.ExitUnsupportedVersion {
		t.Fatalf("exit=%d want %d", ee.ExitCode(), util.ExitUnsupportedVersion)
	}
}

func TestMenuCmd_Help(t *testing.T) {
	cmd := NewRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"menu", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("menu help: %v", err)
	}
	if !strings.Contains(buf.String(), "menu") {
		t.Fatalf("unexpected help: %s", buf.String())
	}
}
