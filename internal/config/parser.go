package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseConfigFile(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("config file not found: %s", path)
        }
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    return ParseConfig(data)
}

func ParseConfig(data []byte) (*Config, error) {
    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
		if err := ValidateConfig(&config); err != nil {
			return nil, fmt.Errorf("config validation failed: %w", err)
		}
    return &config, nil
}