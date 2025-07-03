package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"mess/pkg/config"
	"mess/pkg/repo"
)

// SetupApplication sets up an application by cloning repos and creating symlinks
func SetupApplication(app *config.ApplicationDefinition, cfg *config.MessConfig, configPath string) error {
	// Get the directory containing the config file
	configDir := filepath.Dir(configPath)
	if configDir == "." {
		configDir, _ = os.Getwd()
	}

	// Get applications directory - check MESS_APPLICATION_ROOT environment variable first
	var appsDir string
	if appRoot := os.Getenv("MESS_APPLICATION_ROOT"); appRoot != "" {
		appsDir = appRoot
	} else {
		appsDir = filepath.Join(configDir, "applications")
	}
	
	// Create applications directory if it doesn't exist
	if err := os.MkdirAll(appsDir, 0755); err != nil {
		return fmt.Errorf("failed to create applications directory: %v", err)
	}

	// Create application directory
	appDir := filepath.Join(appsDir, app.Name)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create application directory: %v", err)
	}

	// Execute pre-setup script if defined
	if app.PreSetup != "" {
		fmt.Printf("Executing pre-setup script for application '%s'...\n", app.Name)
		if err := runSingleCommand(app.PreSetup, appsDir, app.Env); err != nil {
			return fmt.Errorf("pre-setup script failed: %v", err)
		}
	}

	// Get repository definitions for the application
	var reposToProcess []config.RepoDefinition
	for _, repoName := range app.Repos {
		for _, repo := range cfg.Repos {
			if repo.Name == repoName {
				reposToProcess = append(reposToProcess, repo)
				break
			}
		}
	}

	// Step 1: Clone any missing repositories
	fmt.Printf("Checking repositories for application '%s'...\n", app.Name)
	for _, repoToProcess := range reposToProcess {
		if !repo.IsRepositoryCloned(repoToProcess.Name, configPath) {
			fmt.Printf("Repository '%s' not found, cloning...\n", repoToProcess.Name)
			if err := repo.CloneRepository(&repoToProcess, configPath); err != nil {
				return fmt.Errorf("failed to clone repository %s: %v", repoToProcess.Name, err)
			}
		} else {
			fmt.Printf("Repository '%s' already cloned\n", repoToProcess.Name)
		}
	}

	// Step 2: Create symbolic links
	fmt.Printf("Creating symbolic links for application '%s'...\n", app.Name)
	for _, repoName := range app.Repos {
		sourcePath := repo.GetRepositoryPath(repoName, configPath)
		linkPath := filepath.Join(appDir, repoName)

		// Remove existing symlink or directory if it exists
		if _, err := os.Lstat(linkPath); err == nil {
			if err := os.Remove(linkPath); err != nil {
				return fmt.Errorf("failed to remove existing link/directory %s: %v", linkPath, err)
			}
		}

		// Create symbolic link
		if err := os.Symlink(sourcePath, linkPath); err != nil {
			return fmt.Errorf("failed to create symbolic link from %s to %s: %v", sourcePath, linkPath, err)
		}
		fmt.Printf("Created symbolic link: %s -> %s\n", linkPath, sourcePath)
	}

	// Execute post-setup script if defined
	if app.PostSetup != "" {
		fmt.Printf("Executing post-setup script for application '%s'...\n", app.Name)
		if err := runSingleCommand(app.PostSetup, appsDir, app.Env); err != nil {
			return fmt.Errorf("post-setup script failed: %v", err)
		}
	}

	return nil
}

// CloneApplication clones application repositories and creates symlinks without setup scripts
func CloneApplication(app *config.ApplicationDefinition, cfg *config.MessConfig, configPath string) error {
	// Get the directory containing the config file
	configDir := filepath.Dir(configPath)
	if configDir == "." {
		configDir, _ = os.Getwd()
	}

	// Get applications directory - check MESS_APPLICATION_ROOT environment variable first
	var appsDir string
	if appRoot := os.Getenv("MESS_APPLICATION_ROOT"); appRoot != "" {
		appsDir = appRoot
	} else {
		appsDir = filepath.Join(configDir, "applications")
	}
	
	// Create applications directory if it doesn't exist
	if err := os.MkdirAll(appsDir, 0755); err != nil {
		return fmt.Errorf("failed to create applications directory: %v", err)
	}

	// Create application directory
	appDir := filepath.Join(appsDir, app.Name)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create application directory: %v", err)
	}

	// Get repository definitions for the application
	var reposToProcess []config.RepoDefinition
	for _, repoName := range app.Repos {
		for _, repo := range cfg.Repos {
			if repo.Name == repoName {
				reposToProcess = append(reposToProcess, repo)
				break
			}
		}
	}

	// Step 1: Clone any missing repositories
	fmt.Printf("Checking repositories for application '%s'...\n", app.Name)
	for _, repoToProcess := range reposToProcess {
		if !repo.IsRepositoryCloned(repoToProcess.Name, configPath) {
			fmt.Printf("Repository '%s' not found, cloning...\n", repoToProcess.Name)
			if err := repo.CloneRepository(&repoToProcess, configPath); err != nil {
				return fmt.Errorf("failed to clone repository %s: %v", repoToProcess.Name, err)
			}
		} else {
			fmt.Printf("Repository '%s' already cloned\n", repoToProcess.Name)
		}
	}

	// Step 2: Create symbolic links
	fmt.Printf("Creating symbolic links for application '%s'...\n", app.Name)
	for _, repoName := range app.Repos {
		sourcePath := repo.GetRepositoryPath(repoName, configPath)
		linkPath := filepath.Join(appDir, repoName)

		// Remove existing symlink or directory if it exists
		if _, err := os.Lstat(linkPath); err == nil {
			if err := os.Remove(linkPath); err != nil {
				return fmt.Errorf("failed to remove existing link/directory %s: %v", linkPath, err)
			}
		}

		// Create symbolic link
		if err := os.Symlink(sourcePath, linkPath); err != nil {
			return fmt.Errorf("failed to create symbolic link from %s to %s: %v", sourcePath, linkPath, err)
		}
		fmt.Printf("Created symbolic link: %s -> %s\n", linkPath, sourcePath)
	}

	return nil
}

// RunScript runs a script for an application
func RunScript(app *config.ApplicationDefinition, scriptName string, scriptValue *config.ScriptValue, configPath string) error {
	// Get the directory containing the config file
	configDir := filepath.Dir(configPath)
	if configDir == "." {
		configDir, _ = os.Getwd()
	}

	// Get applications directory - check MESS_APPLICATION_ROOT environment variable first
	var appsDir string
	if appRoot := os.Getenv("MESS_APPLICATION_ROOT"); appRoot != "" {
		appsDir = appRoot
	} else {
		appsDir = filepath.Join(configDir, "applications")
	}

	// Application directory
	appDir := filepath.Join(appsDir, app.Name)

	// Check if application directory exists
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		return fmt.Errorf("application directory not found: %s. Run 'mess app setup %s' first", appDir, app.Name)
	}

	// Parse script value
	if scriptValue.IsArray {
		// Array of commands
		return runMultipleCommands(scriptValue.Multiple, appDir, app.Env)
	} else {
		// Single command
		return runSingleCommand(scriptValue.Single, appDir, app.Env)
	}
}

// runSingleCommand executes a single command
func runSingleCommand(command, workingDir string, env map[string]string) error {
	fmt.Printf("Executing: %s\n", command)
	
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	// Set environment variables
	if env != nil && len(env) > 0 {
		// Start with current environment
		cmd.Env = os.Environ()
		// Add application-specific environment variables
		for key, value := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return cmd.Run()
}

// runMultipleCommands executes multiple commands in parallel
func runMultipleCommands(commands []string, workingDir string, env map[string]string) error {
	fmt.Printf("Executing %d commands in parallel...\n", len(commands))
	
	var wg sync.WaitGroup
	errChan := make(chan error, len(commands))

	for i, command := range commands {
		wg.Add(1)
		go func(idx int, cmd string) {
			defer wg.Done()
			
			fmt.Printf("[%d] Executing: %s\n", idx+1, cmd)
			
			execCmd := exec.Command("sh", "-c", cmd)
			execCmd.Dir = workingDir
			// For parallel execution, we might want to prefix output
			// but for simplicity, we'll let them write to stdout/stderr directly
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr
			
			// Set environment variables
			if env != nil && len(env) > 0 {
				// Start with current environment
				execCmd.Env = os.Environ()
				// Add application-specific environment variables
				for key, value := range env {
					execCmd.Env = append(execCmd.Env, fmt.Sprintf("%s=%s", key, value))
				}
			}
			
			if err := execCmd.Run(); err != nil {
				errChan <- fmt.Errorf("command [%d] failed: %s - %v", idx+1, cmd, err)
			}
		}(i, command)
	}

	// Wait for all commands to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		fmt.Printf("Some commands failed:\n")
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
		return fmt.Errorf("%d out of %d commands failed", len(errors), len(commands))
	}

	return nil
} 