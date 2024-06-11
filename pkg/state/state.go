package state

import (
	"gopkg.in/yaml.v2"
	"os"
)

type State struct {
	Entries []struct {
		Commit string `yaml:"commit"`
		Path   string `yaml:"path"`
	} `yaml:"entries"`
}

func LoadState(path string) (*State, error) {
	var state State
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (state State) SaveState(path string) error {
	data, err := yaml.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
