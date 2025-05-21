package ansible

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type ExtensionConfig struct {
	Limits      []string `yaml:"limits"`
	Inventories []string `yaml:"inventories"`
	Vaults      []Vault  `yaml:"vaults"`
	Variables   []string `yaml:"variables"`
}

type Vault struct {
	Path     string `yaml:"path"`
	Password string `yaml:"password"`
}

func GetExtensionConfig(configPath string) (ExtensionConfig, error) {
	config, err := getConfigFromFile(configPath)
	if err != nil {
		return ExtensionConfig{}, err
	}
	getConfigFromEnvVars(&config)
	validateExtensionConfig(&config)
	return config, nil
}

func validateExtensionConfig(config *ExtensionConfig) error {
	for _, vault := range config.Vaults {
		if vault.Path != "" && vault.Password == "" {
			return errors.New("vault password required when vault path is set")
		}
	}

	hasInventory := len(config.Inventories) > 0
	hasValidVault := false
	for _, vault := range config.Vaults {
		if vault.Path != "" && vault.Password != "" {
			hasValidVault = true
			break
		}
	}
	hasVariables := len(config.Variables) > 0

	if !hasInventory && !hasValidVault && !hasVariables {
		return errors.New("at least one of inventories, vaults (with path and password), or variables must be set")
	}

	if len(config.Limits) > 0 && !hasInventory {
		return errors.New("inventories must be set if limits are specified")
	}
	if len(config.Inventories) == 0 && len(config.Vaults) == 0 {
		return errors.New("no inventories or vaults defined")
	}
}

func getConfigFromFile(configPath string) (ExtensionConfig, error) {
	if configPath == "" {
		return ExtensionConfig{}, nil
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return ExtensionConfig{}, errors.New("Config file not found with path: " + configPath)
	}
	f, err := os.ReadFile(configPath)
	if err != nil {
		return ExtensionConfig{}, err
	}

	var config ExtensionConfig
	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return ExtensionConfig{}, err
	}

	return config, nil
}

func getConfigFromEnvVars(config *ExtensionConfig) {
	// __ENV.
}
