package serve

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/Skpow1234/Katame/internal/audit"
	"github.com/Skpow1234/Katame/internal/crypto"
	"github.com/Skpow1234/Katame/internal/util"
	"github.com/Skpow1234/Katame/pkg/katame"
)

type keygenRequest struct {
	Algo string `json:"algo,omitempty"`
}

func (s *Server) handleKeygen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req keygenRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	algo := req.Algo
	if algo == "" {
		algo = crypto.SignAlgoEd25519
	}
	if err := s.enforcePolicy(audit.OpKeygen, algo, nil); err != nil {
		writePolicyError(w, err)
		return
	}
	priv, pub, err := crypto.GenerateSigningKeys(algo)
	if err != nil {
		s.auditEntry(audit.OpKeygen, algo, "", "", false, err.Error())
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	s.auditEntry(audit.OpKeygen, algo, "", "", true, "")
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":              true,
		"algo":            algo,
		"private_key_pem": string(priv),
		"public_key_pem":  string(pub),
	})
}

type hashRequest struct {
	DataB64   string `json:"data_b64,omitempty"`
	InputPath string `json:"input_path,omitempty"`
	Algo      string `json:"algo,omitempty"`
}

func (s *Server) handleHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req hashRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	algo := req.Algo
	if algo == "" {
		algo = "sha256"
	}
	ref := req.InputPath
	if ref == "" {
		ref = "<memory>"
	}
	if err := s.enforcePolicy(audit.OpHash, ref, nil); err != nil {
		writePolicyError(w, err)
		return
	}

	var digest []byte
	var err error
	switch {
	case req.InputPath != "" && req.DataB64 != "":
		writeError(w, http.StatusBadRequest, "only one of input_path or data_b64 may be set")
		return
	case req.InputPath != "":
		f, errOpen := os.Open(req.InputPath)
		if errOpen != nil {
			writeError(w, http.StatusBadRequest, errOpen.Error())
			return
		}
		defer f.Close()
		digest, err = crypto.HashReader(f, algo)
	case req.DataB64 != "":
		data, decErr := util.B64Decode(req.DataB64)
		if decErr != nil {
			writeError(w, http.StatusBadRequest, "decode data_b64: "+decErr.Error())
			return
		}
		digest, err = crypto.HashReader(bytes.NewReader(data), algo)
	default:
		writeError(w, http.StatusBadRequest, "one of input_path or data_b64 is required")
		return
	}
	if err != nil {
		s.auditEntry(audit.OpHash, ref, "", "", false, err.Error())
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	s.auditEntry(audit.OpHash, ref, "", util.B64Encode(digest), true, "")
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":         true,
		"algo":       algo,
		"digest_b64": util.B64Encode(digest),
		"digest_hex": fmt.Sprintf("%x", digest),
	})
}

type verifyIntegrityRequest struct {
	InputPath      string `json:"input_path,omitempty"`
	BundleB64      string `json:"bundle_b64,omitempty"`
	KeyB64         string `json:"key_b64,omitempty"`
	Password       string `json:"password,omitempty"`
	KMSProvider    string `json:"kms_provider,omitempty"`
	PrivateKeyPEM  string `json:"private_key_pem,omitempty"`
	PrivateKeyPath string `json:"private_key_path,omitempty"`
	AADB64         string `json:"aad_b64,omitempty"`
}

func (s *Server) handleVerifyIntegrity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req verifyIntegrityRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	bundleRef, manifest, err := loadManifest(req.InputPath, req.BundleB64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.enforcePolicy(audit.OpVerifyIntegrity, bundleRef, manifest); err != nil {
		writePolicyError(w, err)
		return
	}

	opts := katame.DecryptOptions{
		InputPath:      req.InputPath,
		Password:       req.Password,
		KMSProvider:    req.KMSProvider,
		PrivateKey:     []byte(req.PrivateKeyPEM),
		PrivateKeyPath: req.PrivateKeyPath,
	}
	if req.BundleB64 != "" {
		b, err := util.B64Decode(req.BundleB64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "decode bundle_b64: "+err.Error())
			return
		}
		opts.InputBytes = b
	}
	if req.KeyB64 != "" {
		k, err := util.B64Decode(req.KeyB64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "decode key_b64: "+err.Error())
			return
		}
		opts.Key = k
	}
	if req.AADB64 != "" {
		a, err := util.B64Decode(req.AADB64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "decode aad_b64: "+err.Error())
			return
		}
		opts.AAD = a
	}
	opts.KMSUnwrap = func(wrapped []byte, keyID string) ([]byte, error) {
		return s.unwrapKMS(req.KMSProvider, wrapped, keyID)
	}

	res, err := katame.Decrypt(opts)
	if err != nil {
		s.auditEntry(audit.OpVerifyIntegrity, bundleRef, "", "", false, err.Error())
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	wantAlgo := "sha256"
	wantB64 := ""
	if manifest != nil {
		if manifest.Plaintext.Algo != "" {
			wantAlgo = manifest.Plaintext.Algo
		}
		wantB64 = manifest.Plaintext.DigestB64
	}
	got, err := crypto.HashReader(bytes.NewReader(res.Plaintext), wantAlgo)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	match := wantB64 != "" && wantB64 == util.B64Encode(got)
	if !match {
		s.auditEntry(audit.OpVerifyIntegrity, bundleRef, "", "", false, "plaintext hash mismatch")
		writeError(w, http.StatusBadRequest, "plaintext hash mismatch")
		return
	}
	s.auditEntry(audit.OpVerifyIntegrity, bundleRef, "", "", true, "")
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":     true,
		"valid":  true,
		"algo":   wantAlgo,
		"digest": wantB64,
	})
}

type attestRequest struct {
	InputPath string `json:"input_path,omitempty"`
	BundleB64 string `json:"bundle_b64,omitempty"`
}

func (s *Server) handleAttest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req attestRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	bundleRef, manifest, err := loadManifest(req.InputPath, req.BundleB64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.enforcePolicy(audit.OpAttest, bundleRef, manifest); err != nil {
		writePolicyError(w, err)
		return
	}

	var digestHex string
	var inputName string
	var inputSize int64
	var hashAlgo, hashB64 string
	switch {
	case req.InputPath != "":
		if st, err := os.Stat(req.InputPath); err == nil {
			inputSize = st.Size()
		}
		f, err := os.Open(req.InputPath)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		sum, err := crypto.HashReader(f, "sha256")
		_ = f.Close()
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		digestHex = fmt.Sprintf("%x", sum)
	case req.BundleB64 != "":
		raw, err := util.B64Decode(req.BundleB64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		sum, err := crypto.HashReader(bytes.NewReader(raw), "sha256")
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		digestHex = fmt.Sprintf("%x", sum)
		inputSize = int64(len(raw))
	}
	if manifest != nil {
		inputName = manifest.Input.Name
		if manifest.Input.Size > 0 {
			inputSize = manifest.Input.Size
		}
		hashAlgo = manifest.Plaintext.Algo
		hashB64 = manifest.Plaintext.DigestB64
	}
	prov := katame.BuildProvenance(bundleRef, digestHex, inputName, inputSize, hashAlgo, hashB64)
	s.auditEntry(audit.OpAttest, bundleRef, "", digestHex, true, "")
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":         true,
		"provenance": prov,
	})
}
