package serve_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/Skpow1234/Katame/internal/serve"
	"github.com/Skpow1234/Katame/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

func TestReadyzAndRequestID(t *testing.T) {
	base, stop := startServer(t, serve.Options{Listen: "127.0.0.1:0", AuthDisabled: true})
	defer stop()

	req, _ := http.NewRequest(http.MethodGet, base+"/readyz", nil)
	req.Header.Set("X-Request-Id", "test-req-123")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("status %d", resp.StatusCode)
	}
	if got := resp.Header.Get("X-Request-Id"); got != "test-req-123" {
		t.Fatalf("X-Request-Id=%q", got)
	}
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["status"] != "ready" {
		t.Fatalf("body=%v", body)
	}
}

func TestErrorCodeUnauthorized(t *testing.T) {
	base, stop := startServer(t, serve.Options{Listen: "127.0.0.1:0", AuthToken: "secret"})
	defer stop()

	resp, err := http.Get(base + "/v1/version")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("status=%d body=%s", resp.StatusCode, raw)
	}
	if body["error_code"] != "unauthorized" {
		t.Fatalf("error_code=%v", body["error_code"])
	}
}

func TestRateLimit429(t *testing.T) {
	base, stop := startServer(t, serve.Options{
		Listen:         "127.0.0.1:0",
		AuthToken:      "tok",
		RateLimitRPS:   1,
		RateLimitBurst: 1,
	})
	defer stop()

	req, _ := http.NewRequest(http.MethodGet, base+"/v1/version", nil)
	req.Header.Set("Authorization", "Bearer tok")
	resp1, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp1.Body.Close()
	if resp1.StatusCode != 200 {
		t.Fatalf("first status=%d", resp1.StatusCode)
	}

	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	raw, _ := io.ReadAll(resp2.Body)
	if resp2.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d body=%s", resp2.StatusCode, raw)
	}
	var body map[string]any
	_ = json.Unmarshal(raw, &body)
	if body["error_code"] != "rate_limited" {
		t.Fatalf("error_code=%v", body["error_code"])
	}
}

func TestKeygenAndHash(t *testing.T) {
	base, stop := startServer(t, serve.Options{Listen: "127.0.0.1:0", AuthToken: "tok"})
	defer stop()

	kg := postJSON(t, base+"/v1/keygen", "tok", map[string]any{"algo": "ed25519"})
	if kg["_status"] != 200 || kg["ok"] != true {
		t.Fatalf("keygen: %v", kg)
	}
	if kg["private_key_pem"] == nil || kg["public_key_pem"] == nil {
		t.Fatalf("missing keys: %v", kg)
	}

	hs := postJSON(t, base+"/v1/hash", "tok", map[string]any{
		"data_b64": util.B64Encode([]byte("hello-katame")),
		"algo":     "sha256",
	})
	if hs["_status"] != 200 || hs["ok"] != true {
		t.Fatalf("hash: %v", hs)
	}
	if hs["digest_b64"] == nil {
		t.Fatalf("missing digest: %v", hs)
	}
}

func TestOIDCFailure401(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	jwks := map[string]any{
		"keys": []any{
			map[string]any{
				"kty": "RSA",
				"kid": "test",
				"n":   base64.RawURLEncoding.EncodeToString(priv.N.Bytes()),
				"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(priv.E)).Bytes()),
			},
		},
	}
	jwksSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(jwks)
	}))
	defer jwksSrv.Close()

	base, stop := startServer(t, serve.Options{
		Listen: "127.0.0.1:0",
		OIDC: serve.OIDCOptions{
			Issuer:   "https://issuer.example",
			Audience: "katame",
			JWKSURL:  jwksSrv.URL,
		},
	})
	defer stop()

	req, _ := http.NewRequest(http.MethodGet, base+"/v1/version", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "https://issuer.example",
		"aud": "wrong-aud",
		"sub": "user-1",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tok.Header["kid"] = "test"
	signed, err := tok.SignedString(priv)
	if err != nil {
		t.Fatal(err)
	}
	req2, _ := http.NewRequest(http.MethodGet, base+"/v1/version", nil)
	req2.Header.Set("Authorization", "Bearer "+signed)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 401 {
		raw, _ := io.ReadAll(resp2.Body)
		t.Fatalf("expected 401 for bad aud, got %d %s", resp2.StatusCode, raw)
	}
}

func TestUnixSocketMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix sockets")
	}
	dir := t.TempDir()
	sock := filepath.Join(dir, "katame.sock")
	srv, err := serve.NewServer(serve.Options{Listen: "unix:" + sock, AuthDisabled: true})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe(ctx) }()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		st, err := os.Stat(sock)
		if err == nil {
			if st.Mode().Perm()&0o777 != 0o600 {
				t.Fatalf("socket mode=%o want 0600", st.Mode().Perm())
			}
			cancel()
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("socket not created")
}

func TestMaxInFlightProtectSingleOK(t *testing.T) {
	base, stop := startServer(t, serve.Options{
		Listen:             "127.0.0.1:0",
		AuthToken:          "tok",
		MaxInFlightProtect: 1,
	})
	defer stop()

	res := postJSON(t, base+"/v1/protect", "tok", map[string]any{
		"plaintext_b64": util.B64Encode([]byte("x")),
		"key_b64":       util.B64Encode(bytes.Repeat([]byte{1}, 32)),
		"return_bundle": true,
	})
	if res["_status"] == 503 {
		t.Fatalf("single protect should not 503: %v", res)
	}
}
