package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config
type Config struct {
	Resources    []Resource   `yaml:"resources"`
	EventHandler EventHandler `yaml:"eventHandler"`
}

// Resource contains resources to watch
type Resource struct {
	Name string `yaml:"name"`
}

type EventHandler struct {
	Slack Slack `yaml:"slack"`
}

// Slack configuration
type Slack struct {
	Channel string `yaml:"channel"`
	Token   string `yaml:"token,omitempty"`
}

// New returns new Config
func New() (*Config, error) {
	c := &Config{}
	configPath := os.Getenv("CONFIG_PATH")
	configFile := filepath.Join(configPath, "config.yaml")
	file, err := os.Open(configFile)
	defer file.Close()
	if err != nil {
		return c, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return c, err
	}

	if len(b) != 0 {
		yaml.Unmarshal(b, c)
	}
	return c, nil
}
