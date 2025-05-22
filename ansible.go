package ansibleConfigLoader

import (
	"errors"
	"io/fs"
	"maps"
	"os"
	"path/filepath"

	vault "github.com/sosedoff/ansible-vault-go"

	"gopkg.in/yaml.v2"
)

// Original YAML parsing structs (internal, for unmarshalling)
type inventoryYAML struct {
	All struct {
		Children map[string]struct {
			Hosts map[string]map[string]interface{} `yaml:"hosts"`
			Vars  map[string]string                 `yaml:"vars"`
		} `yaml:"children"`
	} `yaml:"all"`
}

type groupConfig struct {
	GroupName string
	Path      string
	Hosts     []hostConfig
	GroupVars map[string]string
}

type GroupConfig struct {
	GroupName string
	Hosts     []HostConfig
	GroupVars map[string]string
}

type InventoryConfig struct {
	InventoryPath string
	Groups        []groupConfig
}

type hostConfig struct {
	HostName string
	Path     string
	HostVars map[string]string
}

type HostConfig struct {
	HostName string
	HostVars map[string]string
}
type AnsibleVariables struct {
	GroupConfig  []GroupConfig
	GlobalConfig map[string]interface{}
}

type AllGroup struct {
	Children map[string]Group `yaml:"children"`
	Hosts    map[string]Host  `yaml:"hosts"`
}

type Group struct {
	Hosts map[string]Host `yaml:"hosts"`
}

type Host struct {
	AnsibleHost string                 `yaml:"ansible_host,omitempty"`
	Vars        map[string]interface{} `yaml:",inline"`
}

func GetVars(config ExtensionConfig) (AnsibleVariables, error) {
	invConf, err := getInventoryVars(config.Inventory, config.Limits)
	if err != nil {
		return AnsibleVariables{}, err
	}
	vaults, err := getAnsibleVaultVars(config.Vaults)
	if err != nil {
		return AnsibleVariables{}, err
	}
	vars, err := getAnsibleVars(config.VarFiles)
	if err != nil {
		return AnsibleVariables{}, err
	}
	groupConf := mapGroupConfig(invConf.Groups)
	maps.Copy(vars, vaults)
	ansibleVars := AnsibleVariables{
		GroupConfig:  groupConf,
		GlobalConfig: vars,
	}
	return ansibleVars, nil
}

func mapGroupConfig(groupConfig []groupConfig) []GroupConfig {
	var groupConfigs []GroupConfig
	for _, g := range groupConfig {
		var hosts []HostConfig
		for _, h := range g.Hosts {
			hosts = append(hosts, HostConfig{
				HostName: h.HostName,
				HostVars: h.HostVars,
			})
		}
		groupConfigs = append(groupConfigs, GroupConfig{
			GroupName: g.GroupName,
			Hosts:     hosts,
			GroupVars: g.GroupVars,
		})

	}
	return groupConfigs
}

func getAnsibleVaultVars(vaults []Vault) (map[string]interface{}, error) {
	if len(vaults) == 0 {
		return nil, nil
	}
	var result map[string]interface{}
	for _, v := range vaults {
		str, err := vault.DecryptFile(v.Path, v.Password)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal([]byte(str), &result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func getAnsibleVars(files []VarFile) (map[string]interface{}, error) {
	if len(files) == 0 {
		return nil, nil
	}
	result := make(map[string]interface{})
	for _, f := range files {
		if _, err := os.Stat(f.Path); os.IsNotExist(err) {
			return nil, errors.New("File not found with path: " + f.Path)
		}
		file, err := os.ReadFile(f.Path)
		if err != nil {
			return nil, err
		}
		var vars map[string]interface{}
		err = yaml.Unmarshal(file, &vars)
		if err != nil {
			return nil, err
		}
		vars, err = decryptVars(vars, f.VaultPassword)
		if err != nil {
			return nil, err
		}
		maps.Copy(result, vars)
	}
	return result, nil
}

func decryptVars(vars map[string]interface{}, password string) (map[string]interface{}, error) {
	for k, v := range vars {
		if strVal, ok := v.(string); ok {
			if len(strVal) > 12 && strVal[:12] == "$ANSIBLE_VAULT" {
				vaultStr, err := vault.Decrypt(strVal, password)
				if err != nil {
					return nil, err
				}
				vars[k] = vaultStr
			}
		}
	}
	return vars, nil
}

func getInventoryVars(inv Inventory, limits []string) (InventoryConfig, error) {
	if inv == (Inventory{}) {
		return InventoryConfig{}, nil
	}
	groups, err := getGroups(inv)
	if err != nil {
		return InventoryConfig{}, err
	}
	groups = filterGroups(groups, limits)
	inventoryConfig := InventoryConfig{
		InventoryPath: inv.Path,
		Groups:        groups,
	}
	inventoryConfig, err = getGroupVars(inventoryConfig, inv.VaultPassword)
	if err != nil {
		return InventoryConfig{}, err
	}
	inventoryConfig, err = getHostVars(inventoryConfig, inv.VaultPassword)
	if err != nil {
		return InventoryConfig{}, err
	}
	return inventoryConfig, nil
}

func filterGroups(groupConf []groupConfig, limits []string) []groupConfig {
	excludeMap := make(map[string]bool)
	for _, name := range limits {
		excludeMap[name] = true
	}

	var result []groupConfig
	for _, g := range groupConf {
		if !excludeMap[g.GroupName] {
			result = append(result, g)
		}
	}
	return result
}

func getHostVars(invConf InventoryConfig, password string) (InventoryConfig, error) {
	result := invConf
	for _, g := range invConf.Groups {
		for _, h := range g.Hosts {
			files, err := getAllFiles(h.Path, h.HostName)
			if err != nil {
				continue
			}
			var vars map[string]interface{}
			for _, file := range files {
				f, err := os.ReadFile(file)
				if err != nil {
					continue
				}
				err = yaml.Unmarshal(f, &vars)
				if err != nil {
					return invConf, err
				}
				vars, err = decryptVars(vars, password)
				if err != nil {
					return invConf, err
				}
			}
			h.HostVars = make(map[string]string)
			for k, v := range vars {
				if strVal, ok := v.(string); ok {
					h.HostVars[k] = strVal
				}
			}
			g.Hosts = append(g.Hosts, h)
		}
		result.Groups = append(result.Groups, g)
	}
	return result, nil
}

func getGroupVars(invConf InventoryConfig, password string) (InventoryConfig, error) {
	result := invConf
	for _, g := range invConf.Groups {
		files, err := getAllFiles(g.Path, g.GroupName)
		if err != nil {
			continue
		}
		var vars map[string]interface{}
		for _, file := range files {
			f, err := os.ReadFile(file)
			if err != nil {
				return invConf, err
			}
			err = yaml.Unmarshal(f, &vars)
			if err != nil {
				return invConf, err
			}
		}
		vars, err = decryptVars(vars, password)
		if err != nil {
			return invConf, err
		}
		g.GroupVars = make(map[string]string)
		for k, v := range vars {
			if strVal, ok := v.(string); ok {
				g.GroupVars[k] = strVal
			}
		}
		result.Groups = append(result.Groups, g)
	}
	return result, nil
}

func getAllFiles(path string, lastFolderName string) ([]string, error) {
	fullPath := filepath.Join(path, lastFolderName)
	var files []string
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, errors.New("Path does not exist: " + fullPath)
	}
	err := filepath.WalkDir(fullPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Only add files, not directories
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func getGroups(inv Inventory) ([]groupConfig, error) {
	if inv.Path == "" {
		return nil, errors.New("Yaml file not found with path: " + inv.Path)
	}
	f, err := os.ReadFile(inv.Path)
	if err != nil {
		return nil, err
	}

	var raw inventoryYAML
	err = yaml.Unmarshal(f, &raw)
	if err != nil {
		return nil, err
	}

	var configs []groupConfig

	for groupName, groupData := range raw.All.Children {
		config := groupConfig{
			GroupName: groupName,
			Path:      inv.GroupVars,
			GroupVars: groupData.Vars,
		}

		for hostName, hostVarsRaw := range groupData.Hosts {
			hostVars := make(map[string]string)
			for k, v := range hostVarsRaw {
				if strVal, ok := v.(string); ok {
					hostVars[k] = strVal
				}
			}

			hostConfig := hostConfig{
				HostName: hostName,
				Path:     inv.HostVars,
				HostVars: hostVars,
			}
			config.Hosts = append(config.Hosts, hostConfig)
		}

		configs = append(configs, config)
	}
	allGroup := groupConfig{
		GroupName: "all",
		Path:      inv.GroupVars,
	}
	configs = append(configs, allGroup)

	return configs, nil
}
