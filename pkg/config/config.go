package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// MessConfig represents the main configuration structure
type MessConfig struct {
	Name         string                         `json:"name"`
	Repos        []RepoDefinition              `json:"repos"`
	Applications []ApplicationDefinition       `json:"applications"`
}

// RepoDefinition represents a repository definition
type RepoDefinition struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	CloneParams []string `json:"clone_params,omitempty"`
}

// ApplicationDefinition represents an application definition
type ApplicationDefinition struct {
	Name      string                     `json:"name"`
	Repos     []string                   `json:"repos"`
	Scripts   map[string]ScriptValue     `json:"scripts"`
	Env       map[string]string          `json:"env,omitempty"`
	PreSetup  string                     `json:"pre-setup,omitempty"`
	PostSetup string                     `json:"post-setup,omitempty"`
}

// ScriptValue represents a script that can be either a string or array of strings
type ScriptValue struct {
	Single   string
	Multiple []string
	IsArray  bool
}

// UnmarshalJSON implements custom JSON unmarshaling for ScriptValue
func (sv *ScriptValue) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		sv.Single = str
		sv.IsArray = false
		return nil
	}

	// Try to unmarshal as array of strings
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		sv.Multiple = arr
		sv.IsArray = true
		return nil
	}

	return fmt.Errorf("script value must be either a string or array of strings")
}

// MarshalJSON implements custom JSON marshaling for ScriptValue
func (sv ScriptValue) MarshalJSON() ([]byte, error) {
	if sv.IsArray {
		return json.Marshal(sv.Multiple)
	}
	return json.Marshal(sv.Single)
}

// LoadConfig loads and validates the mess.json file
func LoadConfig(configPath string) (*MessConfig, error) {
	// If no path provided, try to find mess.json in current directory
	if configPath == "" {
		configPath = "mess.json"
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	// Read the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse JSON
	var config MessConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Validate configuration
	if err := ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the specified file
func SaveConfig(config *MessConfig, configPath string) error {
	// If no path provided, use mess.json in current directory
	if configPath == "" {
		configPath = "mess.json"
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// ValidateConfig validates the configuration structure
func ValidateConfig(config *MessConfig) error {
	if config.Name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Validate repos
	repoNames := make(map[string]bool)
	for _, repo := range config.Repos {
		if repo.Name == "" {
			return fmt.Errorf("repo name cannot be empty")
		}
		if repo.URL == "" {
			return fmt.Errorf("repo URL cannot be empty for repo: %s", repo.Name)
		}
		if repoNames[repo.Name] {
			return fmt.Errorf("duplicate repo name: %s", repo.Name)
		}
		repoNames[repo.Name] = true
	}

	// Validate applications
	appNames := make(map[string]bool)
	for _, app := range config.Applications {
		if app.Name == "" {
			return fmt.Errorf("application name cannot be empty")
		}
		if appNames[app.Name] {
			return fmt.Errorf("duplicate application name: %s", app.Name)
		}
		appNames[app.Name] = true

		// Validate that all referenced repos exist
		for _, repoName := range app.Repos {
			if !repoNames[repoName] {
				return fmt.Errorf("application %s references non-existent repo: %s", app.Name, repoName)
			}
		}
	}

	return nil
}

// CreateSampleConfig creates a sample configuration file
func CreateSampleConfig() *MessConfig {
	return &MessConfig{
		Name: "sample-project",
		Repos: []RepoDefinition{
			{
				Name: "frontend",
				URL:  "https://github.com/example/frontend.git",
			},
			{
				Name:        "backend",
				URL:         "https://github.com/example/backend.git",
				CloneParams: []string{"--depth", "1"},
			},
		},
		Applications: []ApplicationDefinition{
			{
				Name:  "web-app",
				Repos: []string{"frontend", "backend"},
				Scripts: map[string]ScriptValue{
					"start": {Single: "npm start", IsArray: false},
					"build": {Multiple: []string{"npm run build:frontend", "npm run build:backend"}, IsArray: true},
					"test":  {Single: "npm test", IsArray: false},
				},
				Env: map[string]string{
					"NODE_ENV":     "development",
					"API_URL":      "http://localhost:3000",
					"DATABASE_URL": "postgres://localhost:5432/myapp",
				},
				PreSetup:  "echo 'Setting up web application...'",
				PostSetup: "echo 'Web application setup completed!'",
			},
		},
	}
} 