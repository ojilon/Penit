package build

import (
	"fmt"
	"os"
	"path/filepath"
)

// Variant is the Android build variant.
type Variant string

const (
	VariantDebug   Variant = "debug"
	VariantRelease Variant = "release"
)

// ExpectedAPKPaths returns the conventional Gradle output paths for a variant.
func ExpectedAPKPaths(projectRoot string, v Variant) []string {
	base := filepath.Join(projectRoot, "app", "build", "outputs", "apk", string(v))
	switch v {
	case VariantDebug:
		return []string{
			filepath.Join(base, "app-debug.apk"),
		}
	case VariantRelease:
		return []string{
			filepath.Join(base, "app-release-unsigned.apk"),
			filepath.Join(base, "app-release.apk"),
		}
	default:
		return []string{filepath.Join(base, "app-"+string(v)+".apk")}
	}
}

// APKInfo is a found APK with size.
type APKInfo struct {
	Path string
	Size int64
}

// FindAPKs returns existing APKs among the expected paths for a variant.
func FindAPKs(projectRoot string, v Variant) []APKInfo {
	var found []APKInfo
	for _, p := range ExpectedAPKPaths(projectRoot, v) {
		st, err := os.Stat(p)
		if err != nil || st.IsDir() {
			continue
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			abs = p
		}
		found = append(found, APKInfo{Path: abs, Size: st.Size()})
	}
	return found
}

// FormatSize returns a human-readable size.
func FormatSize(n int64) string {
	const (
		kb = 1024
		mb = 1024 * kb
	)
	switch {
	case n >= mb:
		return fmt.Sprintf("%.1f MB", float64(n)/float64(mb))
	case n >= kb:
		return fmt.Sprintf("%.1f KB", float64(n)/float64(kb))
	default:
		return fmt.Sprintf("%d B", n)
	}
}
