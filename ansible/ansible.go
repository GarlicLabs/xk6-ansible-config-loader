package ansible

type AnsibleVariables struct {
	GlobalVariables map[string]string
	GlobalConfig    map[string]string
}
type InventoryConfig struct {
	GroupName string
	Hosts     HostConfig
	GroupVars map[string]string
}

type HostConfig struct {
	HostName string
	HostVars map[string]string
}

func GetVars(config ExtensionConfig) (AnsibleVariables, error) {
	// 1. Read ansible inventory and all group vars, put to sorted vars object
	// 2. Take all groups and get their vars -> Sort them to InventoryConfig
	// 3. Take all hosts of groups and get their vars -> Sort them to InventoryConfig
	// 4. Read ansible vault and put to not sorted vars object
	// 5. Read all Variables and put to not sorted vars object
	// 6. Validate: If Vaults oder Variables are set, then GlobalConfig must be set, If Inventory is set, then GlobalConfig must be set

	return AnsibleVariables{}, nil
}
