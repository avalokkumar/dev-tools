package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// SDKVersion is the manifest schema version this build understands.
// Increment on breaking changes to the manifest layout.
const SDKVersion = 1

// Manifest declares a plugin's identity and the SDK it expects.
type Manifest struct {
	// SDK is the expected SDK schema version. Must equal SDKVersion.
	SDK int `toml:"sdk"`

	// Name uniquely identifies the plugin. Lowercase, kebab-case.
	Name string `toml:"name"`

	// Version is the plugin's own semver string. Informational.
	Version string `toml:"version"`

	// Description is shown in `devforge plugins list`.
	Description string `toml:"description"`

	// Entrypoint is the relative path (within the plugin directory) of the
	// executable to spawn. If empty, defaults to "plugin" (Unix) or "plugin.exe".
	Entrypoint string `toml:"entrypoint"`

	// Tools enumerates the Tools this plugin contributes. Each Tool may
	// declare multiple Operations registered via the SDK initialize call.
	Tools []ManifestTool `toml:"tools"`
}

// ManifestTool is a documentation hint listed in the manifest. The
// authoritative Operation set comes from the plugin's "register" RPC.
type ManifestTool struct {
	Tool        string `toml:"tool"`
	Description string `toml:"description"`
}

// LoadManifest parses dir/manifest.toml and validates required fields.
func LoadManifest(dir string) (*Manifest, error) {
	path := filepath.Join(dir, "manifest.toml")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("plugin: read manifest %s: %w", path, err)
	}
	var m Manifest
	if err := toml.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("plugin: parse manifest %s: %w", path, err)
	}
	if err := m.Validate(); err != nil {
		return nil, fmt.Errorf("plugin: %s: %w", path, err)
	}
	return &m, nil
}

// Validate enforces required fields and SDK compatibility.
func (m *Manifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("manifest: name required")
	}
	if m.SDK != SDKVersion {
		return fmt.Errorf("manifest: sdk version %d, want %d", m.SDK, SDKVersion)
	}
	return nil
}

// EntrypointPath returns the absolute path to the plugin executable.
func (m *Manifest) EntrypointPath(dir string) string {
	name := m.Entrypoint
	if name == "" {
		name = "plugin"
	}
	return filepath.Join(dir, name)
}
