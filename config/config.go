package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Website struct {
	Name  string `yaml:"name"  json:"name"`
	Env   string `yaml:"env"  json:"env"`
	Url   string `yaml:"url"   json:"url"`
	Regex string `yaml:"regex" json:"regex"`
}

type Checker struct {
	Interval int       `yaml:"interval"`
	Timeout  int       `yaml:"timeout"`
	Websites []Website `yaml:"websites"`
}

type Config struct {
	Checker  Checker  `yaml:"checker"`
	Telegram Telegram `yaml:"telegram,omitempty" json:"telegram"`
}

type Telegram struct {
	Token  string `yaml:"token" json:"token"`
	ChatId int64  `yaml:"chat_id" json:"chat_id"`
}

func FromEnv(name string) (*Config, error) {
	filename := os.Getenv(name)
	if filename == "" {
		return nil, fmt.Errorf("env variable %s is undefined", name)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %w", filename, err)
	}

	config := new(Config)
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("unable to parse yaml: %w", err)
	}

	return config, nil
}
