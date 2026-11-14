package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Skpow1234/Katame/internal/config"
	"github.com/Skpow1234/Katame/internal/keysource"
	"github.com/Skpow1234/Katame/internal/kms"
	"github.com/spf13/cobra"
)

type doctorCheck struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Details string `json:"details,omitempty"`
}

func newDoctorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Run local environment diagnostics",
		Long:  "Check local Katame runtime assumptions: Go/runtime info, configured KMS providers, audit/config env vars, and key-source/HSM availability hints.",
		RunE: func(cmd *cobra.Command, args []string) error {
			checks := runDoctorChecks()
			printer := NewPrinter(flagJSON, flagQuiet)
			if flagJSON {
				return printer.JSON(map[string]any{"ok": doctorOK(checks), "checks": checks})
			}
			for _, c := range checks {
				status := "OK"
				if !c.OK {
					status = "WARN"
				}
				printer.Human("%-5s  %-24s %s", status, c.Name, c.Details)
			}
			if !doctorOK(checks) {
				return fmt.Errorf("doctor found warnings")
			}
			return nil
		},
	}
	return cmd
}

func runDoctorChecks() []doctorCheck {
	checks := []doctorCheck{
		{Name: "runtime", OK: true, Details: runtime.GOOS + "/" + runtime.GOARCH + " " + runtime.Version()},
		{Name: "kms providers", OK: len(kms.Providers()) > 0, Details: fmt.Sprintf("%v", kms.Providers())},
	}
	if p := os.Getenv("KATAME_AUDIT_LOG"); p != "" {
		checks = append(checks, doctorCheck{Name: "audit log env", OK: true, Details: p})
	} else {
		checks = append(checks, doctorCheck{Name: "audit log env", OK: true, Details: "not set (audit disabled unless --audit-log is used)"})
	}
	if p := os.Getenv("KATAME_POLICY"); p != "" {
		checks = append(checks, doctorCheck{Name: "policy env", OK: true, Details: p})
	} else {
		checks = append(checks, doctorCheck{Name: "policy env", OK: true, Details: "not set"})
	}
	checks = append(checks,
		doctorCheck{Name: "key-source portable", OK: true, Details: "file:// env:// b64://"},
	)
	for _, s := range keysource.PlatformSchemes() {
		status := "off"
		if s.Implemented {
			status = "on"
		}
		detail := fmt.Sprintf("%s://=%s", s.Scheme, status)
		if s.Details != "" {
			detail += " (" + s.Details + ")"
		}
		checks = append(checks, doctorCheck{
			Name:    "key-source " + s.Scheme,
			OK:      true,
			Details: detail,
		})
	}
	checks = append(checks, doctorPluginDirCheck())
	checks = append(checks, doctorKMSCredCheck())
	return checks
}

func doctorPluginDirCheck() doctorCheck {
	dir := os.Getenv("KATAME_PLUGIN_DIR")
	if dir == "" {
		if c := config.Get(); c != nil {
			dir = c.PluginDir
		}
	}
	if dir == "" {
		return doctorCheck{Name: "plugin dir", OK: true, Details: "not set"}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return doctorCheck{Name: "plugin dir", OK: false, Details: fmt.Sprintf("%s: %v", dir, err)}
	}
	var execs []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		info, err := e.Info()
		if err != nil {
			continue
		}
		mode := info.Mode()
		if mode&0o111 != 0 || filepath.Ext(name) == ".exe" {
			execs = append(execs, name)
		}
	}
	if len(execs) == 0 {
		return doctorCheck{Name: "plugin dir", OK: true, Details: fmt.Sprintf("%s (no executables found)", dir)}
	}
	return doctorCheck{Name: "plugin dir", OK: true, Details: fmt.Sprintf("%s: %v", dir, execs)}
}

func doctorKMSCredCheck() doctorCheck {
	// Soft probe: report which cloud credential env vars are present (no network).
	var present []string
	for _, k := range []string{
		"AWS_ACCESS_KEY_ID", "AWS_PROFILE", "AWS_REGION",
		"GOOGLE_APPLICATION_CREDENTIALS", "CLOUDSDK_CORE_PROJECT",
		"AZURE_CLIENT_ID", "AZURE_TENANT_ID", "AZURE_CLIENT_SECRET",
	} {
		if os.Getenv(k) != "" {
			present = append(present, k)
		}
	}
	if len(present) == 0 {
		return doctorCheck{Name: "kms credentials", OK: true, Details: "no cloud credential env vars detected (mock still works)"}
	}
	return doctorCheck{Name: "kms credentials", OK: true, Details: fmt.Sprintf("env set: %v", present)}
}

func doctorOK(checks []doctorCheck) bool {
	for _, c := range checks {
		if !c.OK {
			return false
		}
	}
	return true
}
