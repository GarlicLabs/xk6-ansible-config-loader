package ansibleConfigLoader

import (
	"testing"
)

func TestGetInventoryVars(t *testing.T) {
	// Positive test case
	t.Run("Valid Inventory", func(t *testing.T) {
		inventory := Inventory{
			Path:          "test/testdata/ansible/prod-inventory.yaml",
			GroupVars:     "test/testdata/ansible/group_vars",
			HostVars:      "test/testdata/ansible/host_vars",
			VaultPassword: "testpassword",
		}
		limits := []string{"group1"}
		_, err := getInventoryVars(inventory, limits)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// Negative test case
	t.Run("Invalid Inventory Path", func(t *testing.T) {
		inventory := Inventory{
			Path: "invalid/path.yml",
		}
		_, err := getInventoryVars(inventory, nil)
		if err == nil {
			t.Error("Expected error for invalid inventory path, got nil")
		}
	})
}

func TestGetAnsibleVaultVars(t *testing.T) {
	// Positive test case
	t.Run("Valid Vault File", func(t *testing.T) {
		vaults := []Vault{
			{
				Path:     "test/testdata/ansible/prod-vault.yaml",
				Password: "S3Cret!",
			},
		}
		_, err := getAnsibleVaultVars(vaults)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// Negative test case
	t.Run("Invalid Vault File Path", func(t *testing.T) {
		vaults := []Vault{
			{
				Path:     "invalid/vault.yml",
				Password: "testpassword",
			},
		}
		_, err := getAnsibleVaultVars(vaults)
		if err == nil {
			t.Error("Expected error for invalid vault file path, got nil")
		}
	})
}

func TestGetAnsibleVars(t *testing.T) {
	// Positive test case
	t.Run("Valid Variable Files", func(t *testing.T) {
		files := []VarFile{
			{
				Path:          "test/testdata/ansible/ansible-vars.yaml",
				VaultPassword: "testpassword",
			},
		}
		_, err := getAnsibleVars(files)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// Negative test case
	t.Run("Invalid Variable File Path", func(t *testing.T) {
		files := []VarFile{
			{
				Path:          "invalid/vars.yml",
				VaultPassword: "testpassword",
			},
		}
		_, err := getAnsibleVars(files)
		if err == nil {
			t.Error("Expected error for invalid variable file path, got nil")
		}
	})
}

func TestMapGroupConfig(t *testing.T) {
	// Positive test case
	t.Run("Valid Group Config", func(t *testing.T) {
		groupConfig := []groupConfig{
			{
				GroupName: "group1",
				Hosts: []hostConfig{
					{
						HostName: "host1",
						HostVars: map[string]string{"var1": "value1"},
					},
				},
				GroupVars: map[string]string{"groupvar1": "value1"},
			},
		}
		result := mapGroupConfig(groupConfig)
		if len(result) != 1 || result[0].GroupName != "group1" {
			t.Errorf("Expected valid group config, got %v", result)
		}
	})

	// Negative test case
	t.Run("Empty Group Config", func(t *testing.T) {
		groupConfig := []groupConfig{}
		result := mapGroupConfig(groupConfig)
		if len(result) != 0 {
			t.Errorf("Expected empty result, got %v", result)
		}
	})
}
