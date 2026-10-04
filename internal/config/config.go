package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Project describes one known Android Gradle project.
type Project struct {
	Path             string   `json:"path"`
	AppName          string   `json:"app_name"`
	PackageIDDebug   string   `json:"package_id_debug"`
	PackageIDRelease string   `json:"package_id_release"`
	MainActivity     string   `json:"main_activity"`
	KeystorePath     string   `json:"keystore_path"`
	KeystoreAlias    string   `json:"keystore_alias"`
	GradleUserHome   string   `json:"gradle_user_home,omitempty"`
	ApksignerPath    string   `json:"apksigner_path,omitempty"`
	LogcatFilters    []string `json:"logcat_filters,omitempty"`
	MirrorInRepo     bool     `json:"mirror_in_repo,omitempty"`
}

// Config is the global Penit configuration.
type Config struct {
	DataRoot   string             `json:"data_root"`
	AndroidSDK string             `json:"android_sdk,omitempty"`
	Projects   map[string]Project `json:"projects"`
}

// Default returns a sensible empty config.
func Default() *Config {
	return &Config{
		DataRoot: defaultDataRoot(),
		Projects: map[string]Project{
			"wayer": {
				AppName:          "Wayer",
				PackageIDDebug:   "com.example.wayer.debug",
				PackageIDRelease: "com.example.wayer",
				MainActivity:     "com.example.wayer.core.MainActivity",
				KeystorePath:     "wayer-release.jks",
				KeystoreAlias:    "wayer",
				LogcatFilters:    []string{"WayerNative", "AndroidRuntime"},
			},
			"conductino_android": {
				AppName:       "Conductino-Study",
				KeystorePath:  "conductino-release.jks",
				KeystoreAlias: "conductino",
				LogcatFilters: []string{"AndroidRuntime"},
			},
			"fdroid": {
				AppName: "FDroid",
			},
			"cooda": {
				AppName: "Cooda",
			},
		},
	}
}

func defaultDataRoot() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "PenitData"
	}
	// Prefer a drive-root style on Windows when possible; otherwise home.
	if _, err := os.Stat("D:\\"); err == nil {
		return `D:\PenitData`
	}
	return filepath.Join(home, "PenitData")
}

// Path returns the default config file location under data-root.
func Path(dataRoot string) string {
	return filepath.Join(dataRoot, "config.json")
}

// Load reads config from path. If missing, returns Default() without writing.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Projects == nil {
		cfg.Projects = map[string]Project{}
	}
	return cfg, nil
}

// Save writes config atomically.
func Save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// ResolveProject returns the project entry or an error.
func (c *Config) ResolveProject(key string) (Project, error) {
	p, ok := c.Projects[key]
	if !ok {
		return Project{}, fmt.Errorf("unknown project %q — known: %v", key, c.ProjectKeys())
	}
	if p.Path == "" {
		return Project{}, fmt.Errorf("project %q has no path set — edit config or run: penit projects add %s --path <abs>", key, key)
	}
	return p, nil
}

// ProjectKeys returns sorted keys for display.
func (c *Config) ProjectKeys() []string {
	keys := make([]string, 0, len(c.Projects))
	for k := range c.Projects {
		keys = append(keys, k)
	}
	return keys
}
