package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ojilon/penit/internal/config"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check data-root, project paths, SDK, adb, apksigner, gradlew",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, cfgPath, err := loadConfig()
		if err != nil {
			return err
		}

		ok := true
		check := func(name string, pass bool, detail string) {
			mark := "✓"
			if !pass {
				mark = "✗"
				ok = false
			}
			fmt.Printf("  %s %-18s %s\n", mark, name, detail)
		}

		fmt.Println("Penit doctor")
		fmt.Printf("  config: %s\n", cfgPath)
		fmt.Println()

		// data-root
		info, err := os.Stat(cfg.DataRoot)
		check("data-root", err == nil && info.IsDir(), cfg.DataRoot)

		// projects
		fmt.Println()
		fmt.Println("Projects")
		for _, key := range cfg.ProjectKeys() {
			p := cfg.Projects[key]
			if p.Path == "" {
				check(key, false, "(no path set)")
				continue
			}
			st, err := os.Stat(p.Path)
			if err != nil || !st.IsDir() {
				check(key, false, p.Path+" — missing")
				continue
			}
			gradlew := "gradlew"
			if runtime.GOOS == "windows" {
				gradlew = "gradlew.bat"
			}
			gw := filepath.Join(p.Path, gradlew)
			_, gwErr := os.Stat(gw)
			detail := p.Path
			if gwErr != nil {
				detail += " — no " + gradlew
				check(key, false, detail)
			} else {
				check(key, true, detail)
			}
		}

		// SDK / adb / apksigner
		fmt.Println()
		fmt.Println("Android tools")
		sdk := cfg.AndroidSDK
		if sdk == "" {
			sdk = os.Getenv("ANDROID_HOME")
			if sdk == "" {
				sdk = os.Getenv("ANDROID_SDK_ROOT")
			}
		}
		check("ANDROID_HOME", sdk != "", sdk)

		adbPath, adbErr := exec.LookPath("adb")
		if adbErr != nil && sdk != "" {
			candidate := filepath.Join(sdk, "platform-tools", "adb")
			if runtime.GOOS == "windows" {
				candidate += ".exe"
			}
			if _, err := os.Stat(candidate); err == nil {
				adbPath = candidate
				adbErr = nil
			}
		}
		check("adb", adbErr == nil, adbPath)

		apsPath, apsErr := exec.LookPath("apksigner")
		if apsErr != nil && sdk != "" {
			bt := filepath.Join(sdk, "build-tools")
			if entries, err := os.ReadDir(bt); err == nil {
				var best string
				for _, e := range entries {
					if !e.IsDir() {
						continue
					}
					cand := filepath.Join(bt, e.Name(), "apksigner")
					if runtime.GOOS == "windows" {
						cand += ".bat"
					}
					if _, err := os.Stat(cand); err == nil {
						if best == "" || e.Name() > filepath.Base(filepath.Dir(best)) {
							best = cand
						}
					}
				}
				if best != "" {
					apsPath = best
					apsErr = nil
				}
			}
		}
		check("apksigner", apsErr == nil, apsPath)

		if adbErr == nil {
			out, err := exec.Command(adbPath, "devices").CombinedOutput()
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			deviceCount := 0
			for i, line := range lines {
				if i == 0 {
					continue
				}
				if strings.Contains(line, "device") && !strings.Contains(line, "offline") {
					deviceCount++
				}
			}
			detail := fmt.Sprintf("%d device(s)", deviceCount)
			if err != nil {
				detail = "adb devices failed"
			}
			check("devices", err == nil && deviceCount > 0, detail)
		}

		fmt.Println()
		if !ok {
			fmt.Println("Some checks failed — fix paths / SDK then re-run doctor.")
			return fmt.Errorf("doctor found issues")
		}
		fmt.Println("All critical checks passed.")
		return nil
	},
}

func loadConfig() (*config.Config, string, error) {
	cfg := config.Default()
	if dataRootFlag != "" {
		cfg.DataRoot = dataRootFlag
	}
	cfgPath := configFlag
	if cfgPath == "" {
		cfgPath = config.Path(cfg.DataRoot)
	}
	loaded, err := config.Load(cfgPath)
	if err != nil {
		return nil, cfgPath, err
	}
	if dataRootFlag != "" {
		loaded.DataRoot = dataRootFlag
	}
	return loaded, cfgPath, nil
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
