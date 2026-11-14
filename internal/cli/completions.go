package cli

import (
	"github.com/Skpow1234/Katame/internal/crypto"
	"github.com/Skpow1234/Katame/internal/plugin"
	"github.com/spf13/cobra"
)

func completeCipher(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return append([]string{}, crypto.SupportedCiphers...), cobra.ShellCompDirectiveNoFileComp
}

func completeHashAlgo(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return append([]string{}, crypto.SupportedHashAlgos...), cobra.ShellCompDirectiveNoFileComp
}

func completeKDF(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{crypto.KDFArgon2id, crypto.KDFScrypt, crypto.KDFPBKDF2SHA256}, cobra.ShellCompDirectiveNoFileComp
}

func completeKeygenAlgos(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	out := append([]string{}, crypto.SupportedSignAlgos...)
	out = append(out, crypto.SupportedHybridSchemes...)
	reg := plugin.GlobalRegistry()
	out = append(out, reg.SignAlgos()...)
	out = append(out, reg.KEMSchemes()...)
	return out, cobra.ShellCompDirectiveNoFileComp
}

func completeCompress(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{crypto.CompressNone, crypto.CompressGzip, crypto.CompressZstd}, cobra.ShellCompDirectiveNoFileComp
}
