package configparser

import (
	"os"

	yaml "github.com/goccy/go-yaml"
)

type LatencyParserConfig struct {
	Global       GlobalConfig        `yaml:"global"`
	Destinations []DestinationConfig `yaml:"destinations"`
}

type GlobalConfig struct {
	DNS             GlobalDnsConfig `yaml:"dns"`
	MetricsEndpoint string          `yaml:"metrics_endpoint"`
}

type GlobalDnsConfig struct {
	Resolver string `yaml:"resolver"`
}

type DestinationConfig struct {
	Name        string `yaml:"name"`
	Endpoint    string `yaml:"endpoint"`
	Method      string `yaml:"method"`
	MetricsSlug string `yaml:"metrics_slug"`
}

func ConfigParser() (*LatencyParserConfig, error) {
	var rawConfig []byte
	var err error
	var config *LatencyParserConfig = &LatencyParserConfig{}
	if configPath := os.Getenv("LATENCY_PARSER_CONFIG_PATH"); configPath != "" {
		rawConfig, err = os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
	} else {
		rawConfig, err = os.ReadFile("/var/latency-parser/config.yml")
		if err != nil {
			return nil, err
		}
	}

	if err = yaml.Unmarshal(rawConfig, config); err != nil {
		return nil, err
	}

	return config, nil
}
