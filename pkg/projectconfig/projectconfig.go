package projectconfig

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Language string `yaml:"language"`
	Repos    []struct {
		URL    string   `yaml:"url"`
		Paths  []string `yaml:"files"`
		Branch string   `yaml:"branch,omitempty"`
		Force  bool     `yaml:"force,omitempty"`
	} `yaml:"repos"`
}

func DefaultConfig() *Config {
	return &Config{
		Language: "go", // Default programming language
		Repos: []struct {
			URL    string   `yaml:"url"`
			Paths  []string `yaml:"files"`
			Branch string   `yaml:"branch,omitempty"`
			Force  bool     `yaml:"force,omitempty"`
		}{
			{ // Initialize with values from your provided example
				URL:    "https://github.com/B-urb/template-repo.git",
				Paths:  []string{"simple.releaserc.yml", "CODEOWNERS", "renovate.json"}, // Correcting "Files" to "Paths"
				Branch: "main",
				Force:  true,
			},
			{ // Example for additional repo (if needed, can be removed)
				URL:    "https://github.com/B-urb/aaand-action.git",
				Paths:  []string{"python"},
				Branch: "main",
				Force:  true,
			},
		},
	}
}

func SaveConfig(cfg *Config, filename string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
