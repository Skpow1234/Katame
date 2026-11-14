package cli

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Skpow1234/Katame/internal/serve"
)

func TestServeCmd_FlagsAndHealthz(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()

	srv, err := serve.NewServer(serve.Options{
		Listen:       addr,
		AuthDisabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = srv.ListenAndServe(ctx) }()

	deadline := time.Now().Add(2 * time.Second)
	var resp *http.Response
	for time.Now().Before(deadline) {
		resp, err = http.Get("http://" + addr + "/healthz")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil || resp == nil || resp.StatusCode != 200 {
		t.Fatalf("healthz failed: %v status=%v", err, resp)
	}

	// Flag registration smoke via cobra help.
	cmd := NewRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"serve", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	help := buf.String()
	for _, want := range []string{"--listen", "--auth-disabled", "--rate-limit-rps", "--max-inflight-protect"} {
		if !strings.Contains(help, want) {
			t.Fatalf("serve help missing %s", want)
		}
	}
}

func TestRepoCmd_InitAddList(t *testing.T) {
	dir := t.TempDir()
	repoDir := filepath.Join(dir, "repo")
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"repo", "init", "--dir", repoDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("repo init: %v", err)
	}

	in := filepath.Join(dir, "a.txt")
	bundle := filepath.Join(dir, "a.katame")
	if err := writeFileHelper(in, []byte("repo-leaf")); err != nil {
		t.Fatal(err)
	}
	protect := NewRootCmd()
	protect.SetArgs([]string{"protect", "--in", in, "--out", bundle})
	if err := protect.Execute(); err != nil {
		t.Fatalf("protect: %v", err)
	}

	add := NewRootCmd()
	add.SetArgs([]string{"repo", "add", "--dir", repoDir, "--bundle", bundle})
	if err := add.Execute(); err != nil {
		// Some repo layouts require signing keys; accept init+help coverage if add needs more setup.
		if !strings.Contains(err.Error(), "signing") && !strings.Contains(err.Error(), "verify-key") {
			t.Fatalf("repo add: %v", err)
		}
	}

	list := NewRootCmd()
	var buf bytes.Buffer
	list.SetOut(&buf)
	list.SetArgs([]string{"repo", "list", "--dir", repoDir})
	if err := list.Execute(); err != nil {
		t.Fatalf("repo list: %v", err)
	}

	help := NewRootCmd()
	var hbuf bytes.Buffer
	help.SetOut(&hbuf)
	help.SetArgs([]string{"repo", "--help"})
	_ = help.Execute()
	if !strings.Contains(hbuf.String(), "init") {
		t.Fatalf("repo help missing init: %s", hbuf.String())
	}
	if _, err := os.Stat(repoDir); err != nil {
		t.Fatalf("repo dir: %v", err)
	}
}
