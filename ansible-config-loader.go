package ansibleConfigLoader

import (
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/ansible-config-loader", new(AnsibleVariables))
}

type EnvironmentVariables struct {
	ConfigPath  string
	Inventories string
	Vaults      string
	Limits      string
}

func (c *AnsibleVariables) GetConfig(configPath string) (AnsibleVariables, error) {
	extensionConfig, err := GetExtensionConfig(configPath)
	if err != nil {
		panic(err)
	}
	ansibleVars, err := GetVars(extensionConfig)
	if err != nil {
		panic(err)
	}
	return ansibleVars, nil
}
