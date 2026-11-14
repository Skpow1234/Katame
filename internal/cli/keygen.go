package cli

import (
	"fmt"
	"strings"

	"github.com/Skpow1234/Katame/internal/crypto"
	"github.com/Skpow1234/Katame/internal/plugin"
	"github.com/spf13/cobra"
)

func newKeygenCmd() *cobra.Command {
	var (
		outPrefix string
		algo      string
		listAlgos bool
	)

	cmd := &cobra.Command{
		Use:   "keygen",
		Short: "Generate a key pair",
		Long: "Generate a new key pair for signing, verifying, or hybrid encryption.\n\n" +
			"Signing algorithms: ed25519 (default), ed448, ecdsa-p256, ecdsa-p384, rsa-pss-2048, rsa-pss-4096, ml-dsa-65, ml-dsa-87, slh-dsa-128s.\n" +
			"Hybrid encryption: x25519, ecies-p256, rsa-oaep-2048, rsa-oaep-4096, ml-kem-768, ml-kem-1024.\n" +
			"Plugin schemes from KATAME_PLUGIN_DIR (or config plugin_dir) are added to --algo; see docs/plugins.md.",
		RunE: func(cmd *cobra.Command, args []string) error {
			printer := NewPrinter(flagJSON, flagQuiet)

			if listAlgos {
				reg := plugin.GlobalRegistry()
				signing := append([]string{}, crypto.SupportedSignAlgos...)
				hybrid := append([]string{}, crypto.SupportedHybridSchemes...)
				for _, a := range reg.SignAlgos() {
					signing = append(signing, a)
				}
				for _, a := range reg.KEMSchemes() {
					hybrid = append(hybrid, a)
				}
				if flagJSON {
					return printer.JSON(map[string]any{
						"signing": signing,
						"hybrid":  hybrid,
					})
				}
				printer.Human("Signing: %s", strings.Join(signing, ", "))
				printer.Human("Hybrid:  %s", strings.Join(hybrid, ", "))
				return nil
			}

			if outPrefix == "" {
				return fmt.Errorf("--out is required (e.g. --out signing)")
			}

			privPath := outPrefix
			if !strings.HasSuffix(privPath, ".key") {
				privPath = outPrefix + ".key"
			}
			pubPath := strings.TrimSuffix(privPath, ".key") + ".pub"

			var err error
			var purpose string
			reg := plugin.GlobalRegistry()

			// Built-in or plugin: signing first, then hybrid (KEM).
			if crypto.SupportedSignAlgo(algo) {
				privPEM, pubPEM, err := crypto.GenerateSigningKeys(algo)
				if err != nil {
					return fmt.Errorf("generate key pair: %w", err)
				}
				purpose = "signing"
				if err := crypto.SaveKeyPEM(privPath, privPEM, 0o600); err != nil {
					return fmt.Errorf("save private key: %w", err)
				}
				if err := crypto.SaveKeyPEM(pubPath, pubPEM, 0o644); err != nil {
					return fmt.Errorf("save public key: %w", err)
				}
			} else if reg.SignAlgo(algo) != "" {
				err = reg.KeygenSign(algo, outPrefix)
				purpose = "signing"
				if err != nil {
					return fmt.Errorf("generate key pair: %w", err)
				}
			} else if crypto.SupportedHybridScheme(algo) {
				privPEM, pubPEM, err := crypto.GenerateHybridKeys(algo)
				if err != nil {
					return fmt.Errorf("generate key pair: %w", err)
				}
				purpose = "encryption"
				if err := crypto.SaveKeyPEM(privPath, privPEM, 0o600); err != nil {
					return fmt.Errorf("save private key: %w", err)
				}
				if err := crypto.SaveKeyPEM(pubPath, pubPEM, 0o644); err != nil {
					return fmt.Errorf("save public key: %w", err)
				}
			} else if reg.KEMScheme(algo) != "" {
				err = reg.KeygenKEM(algo, outPrefix)
				purpose = "encryption"
				if err != nil {
					return fmt.Errorf("generate key pair: %w", err)
				}
			} else {
				return fmt.Errorf("unsupported algorithm %q; signing: ed25519, ed448, ecdsa-p256, ecdsa-p384, rsa-pss-2048, rsa-pss-4096, ml-dsa-65, ml-dsa-87, slh-dsa-128s; encryption: x25519, ecies-p256, rsa-oaep-2048, rsa-oaep-4096, ml-kem-768, ml-kem-1024 (or plugin schemes)", algo)
			}

			switch printer.Mode {
			case OutputJSON:
				return printer.JSON(map[string]any{
					"private_key": privPath,
					"public_key":  pubPath,
					"algorithm":   algo,
					"purpose":     purpose,
					"format":      "pem",
				})
			default:
				printer.Human("Generated %s %s key pair:", algo, purpose)
				printer.Human("  Private: %s (PEM/PKCS#8)", privPath)
				printer.Human("  Public:  %s (PEM/PKIX)", pubPath)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&outPrefix, "out", "", "output path prefix (e.g. 'signing' produces signing.key + signing.pub)")
	cmd.Flags().StringVar(&algo, "algo", crypto.SignAlgoEd25519, "algorithm (see --help for full list)")
	cmd.Flags().BoolVar(&listAlgos, "list-algos", false, "list built-in and plugin signing/hybrid algorithms and exit")
	_ = cmd.RegisterFlagCompletionFunc("algo", completeKeygenAlgos)

	return cmd
}
