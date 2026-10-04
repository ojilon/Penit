package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Layout helpers for the Penit data-root.
//
//   <data_root>/
//     config.json
//     history/
//     projects/
//       <key>/
//         cache.json
//         logs/
//         releases/
//           <version>/

func EnsureRoot(dataRoot string) error {
	dirs := []string{
		dataRoot,
		filepath.Join(dataRoot, "history"),
		filepath.Join(dataRoot, "projects"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", d, err)
		}
	}
	return nil
}

func EnsureProject(dataRoot, key string) error {
	base := filepath.Join(dataRoot, "projects", key)
	dirs := []string{
		base,
		filepath.Join(base, "logs"),
		filepath.Join(base, "releases"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", d, err)
		}
	}
	return nil
}

func ProjectDir(dataRoot, key string) string {
	return filepath.Join(dataRoot, "projects", key)
}

func LogsDir(dataRoot, key string) string {
	return filepath.Join(dataRoot, "projects", key, "logs")
}

func ReleasesDir(dataRoot, key string) string {
	return filepath.Join(dataRoot, "projects", key, "releases")
}

func ReleaseVersionDir(dataRoot, key, version string) string {
	return filepath.Join(dataRoot, "projects", key, "releases", version)
}

// NewBuildLogPath returns a timestamped log file path for a build.
func NewBuildLogPath(dataRoot, key, variant string) string {
	ts := time.Now().Format("20060102-150405")
	name := fmt.Sprintf("build-%s-%s.log", ts, variant)
	return filepath.Join(LogsDir(dataRoot, key), name)
}

// NewLogcatPath returns a timestamped logcat capture path.
func NewLogcatPath(dataRoot, key string) string {
	ts := time.Now().Format("20060102-150405")
	name := fmt.Sprintf("logcat-%s.log", ts)
	return filepath.Join(LogsDir(dataRoot, key), name)
}
