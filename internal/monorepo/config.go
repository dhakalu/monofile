// parses the mono repo config
package monorepo

import (
	"dhakalu/monofile/internal/parsers"
	"io/fs"

	"github.com/go-yaml/yaml"
)

// Config represents the root of the monorepo configuration file.
type MonofileConfig struct {
	Name     string             `yaml:"name"`
	Defaults Defaults           `yaml:"defaults"`
	Scan     ScanConfig         `yaml:"scan"`
	Projects map[string]Project `yaml:"projects"`
}

// Defaults handles global and language-specific configurations.
type Defaults struct {
	Env       map[string]string             `yaml:"env"`        // Global environment variables
	BuildArgs map[parsers.SourceLang]string `yaml:"build-args"` // Language-specific default CLI arguments (e.g., "dotnet": "--arch x64")
}

// ScanConfig defines the auto-scanning engine behavior.
type ScanConfig struct {
	Enabled *bool               `yaml:"enabled"` // Pointer to bool to distinguish between "missing/nil" (default true) and "false"
	Exclude []string            `yaml:"exclude"` // Glob patterns to ignore during discovery
	Types   map[string]ScanType `yaml:"types"`   // Definitions for automated language scanners
}

// ScanType contains the rules for how the scanner identifies and builds a language family.
type ScanType struct {
	Match string `yaml:"match"` // File glob pattern to search for (e.g., "**/*.csproj")
	Build string `yaml:"build"` // The default build token string (e.g., "dotnet build {path}")
}

// Project defines explicit manual overrides or completely custom projects.
type Project struct {
	Path  string            `yaml:"path,omitempty"`  // Path to the project root or file
	Type  string            `yaml:"type,omitempty"`  // Language type matching ScanConfig.Types keys
	Build string            `yaml:"build,omitempty"` // Custom build command override
	Env   map[string]string `yaml:"env,omitempty"`   // Project-specific environment variables
}

func ParseMonofileConfig(s Scanner) (MonofileConfig, error) {
	cfg := MonofileConfig{}
	raw, err := fs.ReadFile(s.FileSystem, ".mono.yaml")
	if err != nil {
		return MonofileConfig{}, err
	}
	err = yaml.Unmarshal(raw, &cfg)
	return cfg, err
}
