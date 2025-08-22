package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RepoPath          string `yaml:"repo_path"`
	GitUser           string `yaml:"git_user"`
	GitEmail          string `yaml:"git_email"`
	Timezone          string `yaml:"timezone"`
	Hour24            int    `yaml:"hour_24"`
	Minute            int    `yaml:"minute"`
	IntensityStrategy string `yaml:"intensity_strategy"`
	IntensityValue    int    `yaml:"intensity_value"`
	PatternFile       string `yaml:"pattern_file"`
}

func defaultPaths() (cfgPath, stateDir string, err error) {
	cfgBase, err := os.UserConfigDir()
	if err != nil { return "", "", err }
	stateBase, err := os.UserHomeDir()
	if err != nil { return "", "", err }
	cfgDir := filepath.Join(cfgBase, "commit-grid-draw")
	stateDir = filepath.Join(stateBase, ".local", "state", "commit-grid-draw")
	_ = os.MkdirAll(cfgDir, 0o755)
	_ = os.MkdirAll(stateDir, 0o755)
	return filepath.Join(cfgDir, "config.yaml"), stateDir, nil
}

func Load() (*Config, string, error) {
	cfgFile, _, err := defaultPaths()
	if err != nil { return nil, "", err }
	b, err := os.ReadFile(cfgFile)
	if err != nil { return nil, "", err }
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil { return nil, "", err }
	return &c, cfgFile, nil
}

func Save(c *Config) (string, error) {
	cfgFile, _, err := defaultPaths()
	if err != nil { return "", err }
	b, err := yaml.Marshal(c)
	if err != nil { return "", err }
	return cfgFile, os.WriteFile(cfgFile, b, 0o644)
}

func Exists() bool {
	cfgFile, _, err := defaultPaths()
	if err != nil { return false }
	_, err = os.Stat(cfgFile)
	return err == nil
}

func Debug(c *Config) string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}
