package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to setup and teardown test config
func setupTestConfig(t *testing.T) (configPath string, cleanup func()) {
	t.Helper()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	configPath = filepath.Join(homeDir, configFileName)
	backupPath := configPath + ".backup"

	// Backup existing config if it exists
	existingConfig, err := os.ReadFile(configPath)
	hasExistingConfig := err == nil
	if hasExistingConfig {
		if err := os.WriteFile(backupPath, existingConfig, 0644); err != nil {
			t.Fatalf("failed to backup existing config: %v", err)
		}
	}

	cleanup = func() {
		if hasExistingConfig {
			os.WriteFile(configPath, existingConfig, 0644)
			os.Remove(backupPath)
		} else {
			os.Remove(configPath)
		}
	}

	return configPath, cleanup
}

func TestRead_ValidConfig(t *testing.T) {
	configPath, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create test config
	testConfig := map[string]string{
		"db_url":            "postgres://localhost:5432/testdb",
		"current_user_name": "testuser",
	}
	configData, err := json.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Test
	config, err := Read()
	if err != nil {
		t.Fatalf("Read() returned unexpected error: %v", err)
	}

	// Verify config was populated correctly
	if config.DbUrl != "postgres://localhost:5432/testdb" {
		t.Errorf("expected DbUrl 'postgres://localhost:5432/testdb', got: %s", config.DbUrl)
	}
	if config.CurrentUser != "testuser" {
		t.Errorf("expected CurrentUser 'testuser', got: %s", config.CurrentUser)
	}
}

func TestRead_NonExistentConfig(t *testing.T) {
	configPath, cleanup := setupTestConfig(t)
	defer cleanup()

	// Ensure config file doesn't exist
	os.Remove(configPath)

	// Test: Read should return error when file doesn't exist
	config, err := Read()

	if err == nil {
		t.Error("expected error when config file doesn't exist, got nil")
	}
	if config != (Config{}) {
		t.Errorf("expected empty Config when file doesn't exist, got: %+v", config)
	}
}

func TestRead_InvalidJSON(t *testing.T) {
	configPath, cleanup := setupTestConfig(t)
	defer cleanup()

	// Write invalid JSON
	if err := os.WriteFile(configPath, []byte("not valid json{"), 0644); err != nil {
		t.Fatalf("failed to write invalid config: %v", err)
	}

	// Test: Read should return error for invalid JSON
	config, err := Read()

	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
	if config != (Config{}) {
		t.Errorf("expected empty Config for invalid JSON, got: %+v", config)
	}
}

func TestSetUser_Success(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create initial config
	config := Config{
		DbUrl:       "postgres://localhost:5432/testdb",
		CurrentUser: "",
	}

	// Test SetUser
	err := config.SetUser("newuser")
	if err != nil {
		t.Fatalf("SetUser() returned unexpected error: %v", err)
	}

	// Verify the user was set in memory
	if config.CurrentUser != "newuser" {
		t.Errorf("expected CurrentUser 'newuser', got: %s", config.CurrentUser)
	}

	// Verify the config was written to file
	readConfig, err := Read()
	if err != nil {
		t.Fatalf("failed to read config after SetUser: %v", err)
	}
	if readConfig.CurrentUser != "newuser" {
		t.Errorf("expected persisted CurrentUser 'newuser', got: %s", readConfig.CurrentUser)
	}
}

func TestSetUser_UpdatesExistingUser(t *testing.T) {
	configPath, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create initial config with existing user
	initialConfig := map[string]string{
		"db_url":            "postgres://localhost:5432/testdb",
		"current_user_name": "olduser",
	}
	configData, _ := json.Marshal(initialConfig)
	os.WriteFile(configPath, configData, 0644)

	// Load and update user
	config, err := Read()
	if err != nil {
		t.Fatalf("failed to read initial config: %v", err)
	}

	err = config.SetUser("updateduser")
	if err != nil {
		t.Fatalf("SetUser() returned unexpected error: %v", err)
	}

	// Verify the update persisted
	readConfig, err := Read()
	if err != nil {
		t.Fatalf("failed to read config after update: %v", err)
	}
	if readConfig.CurrentUser != "updateduser" {
		t.Errorf("expected CurrentUser 'updateduser', got: %s", readConfig.CurrentUser)
	}
	// Verify DbUrl was preserved
	if readConfig.DbUrl != "postgres://localhost:5432/testdb" {
		t.Errorf("expected DbUrl to be preserved, got: %s", readConfig.DbUrl)
	}
}

func TestGetConfigFilePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	expectedPath := filepath.Join(homeDir, configFileName)

	path, err := getConfigFilePath()
	if err != nil {
		t.Errorf("getConfigFilePath() returned unexpected error: %v", err)
	}
	if path != expectedPath {
		t.Errorf("expected path '%s', got: '%s'", expectedPath, path)
	}
}

func TestWrite_CreatesValidJSON(t *testing.T) {
	configPath, cleanup := setupTestConfig(t)
	defer cleanup()

	// Remove any existing config
	os.Remove(configPath)

	config := Config{
		DbUrl:       "postgres://localhost:5432/mydb",
		CurrentUser: "writeuser",
	}

	err := write(config)
	if err != nil {
		t.Fatalf("write() returned unexpected error: %v", err)
	}

	// Verify file was created with correct content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read written config: %v", err)
	}

	var readConfig Config
	if err := json.Unmarshal(data, &readConfig); err != nil {
		t.Fatalf("written config is not valid JSON: %v", err)
	}

	if readConfig.DbUrl != config.DbUrl {
		t.Errorf("expected DbUrl '%s', got: '%s'", config.DbUrl, readConfig.DbUrl)
	}
	if readConfig.CurrentUser != config.CurrentUser {
		t.Errorf("expected CurrentUser '%s', got: '%s'", config.CurrentUser, readConfig.CurrentUser)
	}
}

func TestWrite_OverwritesExistingConfig(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Write initial config
	initialConfig := Config{DbUrl: "old_url", CurrentUser: "old_user"}
	write(initialConfig)

	// Overwrite with new config
	newConfig := Config{DbUrl: "new_url", CurrentUser: "new_user"}
	err := write(newConfig)
	if err != nil {
		t.Fatalf("write() returned unexpected error: %v", err)
	}

	// Verify new config was written
	readConfig, err := Read()
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	if readConfig.DbUrl != "new_url" {
		t.Errorf("expected DbUrl 'new_url', got: '%s'", readConfig.DbUrl)
	}
	if readConfig.CurrentUser != "new_user" {
		t.Errorf("expected CurrentUser 'new_user', got: '%s'", readConfig.CurrentUser)
	}
}
