package cli

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/Skpow1234/Katame/internal/util"
	"github.com/spf13/cobra"
)

func newConvertCmd() *cobra.Command {
	var (
		inFile  string
		outFile string
		force   bool
	)
	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert a legacy Ed25519 key to PKCS#8 / PKIX PEM",
		Long: `Convert Katame's legacy ed25519-priv: / ed25519-pub: key files to
standard PKCS#8 (private) or PKIX (public) PEM. Existing PEM files are rejected
unless they are already the target format (then --force overwrites).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			printer := NewPrinter(flagJSON, flagQuiet)
			if inFile == "" {
				return fmt.Errorf("--in is required")
			}
			if outFile == "" {
				return fmt.Errorf("--out is required")
			}
			raw, err := os.ReadFile(inFile)
			if err != nil {
				return fmt.Errorf("read input: %w", err)
			}
			content := strings.TrimSpace(string(raw))
			var out []byte
			var kind string
			switch {
			case strings.HasPrefix(content, "ed25519-priv:"):
				b64 := strings.TrimPrefix(content, "ed25519-priv:")
				keyRaw, err := util.B64Decode(b64)
				if err != nil {
					return fmt.Errorf("decode legacy private key: %w", err)
				}
				if len(keyRaw) != ed25519.PrivateKeySize {
					return fmt.Errorf("invalid ed25519 private key size: %d", len(keyRaw))
				}
				der, err := x509.MarshalPKCS8PrivateKey(ed25519.PrivateKey(keyRaw))
				if err != nil {
					return fmt.Errorf("marshal PKCS#8: %w", err)
				}
				out = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
				kind = "private-pkcs8"
			case strings.HasPrefix(content, "ed25519-pub:"):
				b64 := strings.TrimPrefix(content, "ed25519-pub:")
				keyRaw, err := util.B64Decode(b64)
				if err != nil {
					return fmt.Errorf("decode legacy public key: %w", err)
				}
				if len(keyRaw) != ed25519.PublicKeySize {
					return fmt.Errorf("invalid ed25519 public key size: %d", len(keyRaw))
				}
				der, err := x509.MarshalPKIXPublicKey(ed25519.PublicKey(keyRaw))
				if err != nil {
					return fmt.Errorf("marshal PKIX: %w", err)
				}
				out = pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
				kind = "public-pkix"
			default:
				return fmt.Errorf("input is not a legacy ed25519-priv:/ed25519-pub: key (already PEM?)")
			}
			if err := refuseIfExists(outFile, force); err != nil {
				return err
			}
			mode := os.FileMode(0o644)
			if kind == "private-pkcs8" {
				mode = 0o600
			}
			if err := os.WriteFile(outFile, out, mode); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			if flagJSON {
				return printer.JSON(map[string]any{
					"ok":   true,
					"in":   inFile,
					"out":  outFile,
					"kind": kind,
				})
			}
			printer.Human("Converted %s → %s (%s)", inFile, outFile, kind)
			return nil
		},
	}
	cmd.Flags().StringVar(&inFile, "in", "", "legacy key file (required)")
	cmd.Flags().StringVar(&outFile, "out", "", "output PEM path (required)")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing output")
	return cmd
}
