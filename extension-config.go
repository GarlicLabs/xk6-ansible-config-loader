package ansibleConfigLoader

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type ExtensionConfig struct {
	Limits    []string  `yaml:"limits"`
	Inventory Inventory `yaml:"inventory"`
	Vaults    []Vault   `yaml:"vaults"`
	VarFiles  []VarFile `yaml:"var_files"`
}

type VarFile struct {
	Path          string `yaml:"path"`
	VaultPassword string `yaml:"vault_password"`
}

type Vault struct {
	Path     string `yaml:"path"`
	Password string `yaml:"password"`
}

type Inventory struct {
	Path          string `yaml:"path"`
	GroupVars     string `yaml:"group_vars"`
	HostVars      string `yaml:"host_vars"`
	VaultPassword string `yaml:"vault_password"`
}

func GetExtensionConfig(configPath string) (ExtensionConfig, error) {
	config, err := getConfigFromFile(configPath)
	if err != nil {
		return ExtensionConfig{}, err
	}
	validateExtensionConfig(&config)
	return config, nil
}

func validateExtensionConfig(config *ExtensionConfig) error {
	for _, vault := range config.Vaults {
		if vault.Path != "" && vault.Password == "" {
			return errors.New("vault password required when vault path is set")
		}
	}

	isInventory := config.Inventory != Inventory{}
	hasValidVault := false
	for _, vault := range config.Vaults {
		if vault.Path != "" && vault.Password != "" {
			hasValidVault = true
			break
		}
	}
	hasVariables := len(config.VarFiles) > 0

	if !isInventory && !hasValidVault && !hasVariables {
		return errors.New("at least one of inventories, vaults (with path and password), or variables must be set")
	}

	if len(config.Limits) > 0 && !isInventory {
		return errors.New("inventories must be set if limits are specified")
	}
	return nil
}

func getConfigFromFile(configPath string) (ExtensionConfig, error) {
	if configPath == "" {
		return ExtensionConfig{}, errors.New("Config path is empty")
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
