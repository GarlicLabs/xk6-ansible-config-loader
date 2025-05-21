package ansibleConfigLoader

import (
	"os"

	"github.com/GarlicLabs/xk6-ansible-config-loader/ansible"
	"github.com/dop251/goja"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/ansible-config-loader", new(Module))
}

type Module struct {
	AnsibleVariables ansible.AnsibleVariables
}

type EnvironmentVariables struct {
	ConfigPath  string
	Inventories string
	Vaults      string
	Limits      string
}

func (c *Module) GetConfig(configPath string) (Module, error) {
	extensionConfig, err := ansible.GetExtensionConfig(configPath)
	if err != nil {
		panic(err)
	}
	ansibleVars, err := ansible.GetVars(extensionConfig)
	if err != nil {
		panic(err)
	}
	return Module{AnsibleVariables: ansibleVars}, nil
}

func getEnvVars(rt *goja.Runtime) EnvironmentVariables {
	var envVars EnvironmentVariables
	if val, ok := os.LookupEnv("XK6_CONFIG_PATH"); ok {
		envVars.ConfigPath = val
	}
	if val, ok := os.LookupEnv("XK6_INVENTORIES"); ok {
		envVars.Inventories = val
	}
	if val, ok := os.LookupEnv("XK6_VAULTS"); ok {
		envVars.Vaults = val
	}
	if val, ok := os.LookupEnv("XK6_LIMITS"); ok {
		envVars.Limits = val
	}
	return envVars
}
