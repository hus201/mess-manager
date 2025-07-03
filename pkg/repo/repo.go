package repo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"mess/pkg/config"
)

// CloneRepository clones a repository to the appropriate directory
func CloneRepository(repo *config.RepoDefinition, configPath string) error {
	// Get the directory containing the config file
	configDir := filepath.Dir(configPath)
	if configDir == "." {
		configDir, _ = os.Getwd()
	}

	// Create repos directory if it doesn't exist
	reposDir := filepath.Join(configDir, "repos")
	if err := os.MkdirAll(reposDir, 0755); err != nil {
		return fmt.Errorf("failed to create repos directory: %v", err)
	}

	// Target directory for the repository
	targetDir := filepath.Join(reposDir, repo.Name)

	// Check if repository already exists
	if _, err := os.Stat(targetDir); err == nil {
		return fmt.Errorf("repository directory already exists: %s", targetDir)
	}

	// Clone the repository
	fmt.Printf("Cloning repository %s from %s...\n", repo.Name, repo.URL)
	
	// Build git clone command with optional clone_params
	args := []string{"clone"}
	if len(repo.CloneParams) > 0 {
		args = append(args, repo.CloneParams...)
	}
	args = append(args, repo.URL, targetDir)
	
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Clean up partially cloned directory on failure
		os.RemoveAll(targetDir)
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	return nil
}

// IsRepositoryCloned checks if a repository is already cloned
func IsRepositoryCloned(repoName, configPath string) bool {
	// Get the directory containing the config file
	configDir := filepath.Dir(configPath)
	if configDir == "." {
		configDir, _ = os.Getwd()
	}

	// Target directory for the repository
	targetDir := filepath.Join(configDir, "repos", repoName)

	// Check if directory exists and contains .git directory
	if stat, err := os.Stat(targetDir); err == nil && stat.IsDir() {
		gitDir := filepath.Join(targetDir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return true
		}
	}

	return false
}

// GetRepositoryPath returns the path to a cloned repository
func GetRepositoryPath(repoName, configPath string) string {
	configDir := filepath.Dir(configPath)
	if configDir == "." {
		configDir, _ = os.Getwd()
	}
	return filepath.Join(configDir, "repos", repoName)
} 