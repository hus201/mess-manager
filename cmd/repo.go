package cmd

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"os/exec"

	"github.com/spf13/cobra"
	"mess/pkg/config"
	"mess/pkg/repo"
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo <repo-name> <action>",
	Short: "Manage repositories",
	Long: `Manage repositories in your mess.json file.
Available actions:
  add <repo-url>    - Add a new repository with the given URL
  remove           - Remove a repository (aliases: rm)
  get              - Clone a repository (aliases: clone)
  <git-command>    - Execute git command on the repository`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		repoName := args[0]
		action := args[1]
		remainingArgs := args[2:]

		switch action {
		case "add":
			handleRepoAdd(repoName, remainingArgs)
		case "remove", "rm":
			handleRepoRemove(repoName, remainingArgs)
		case "get", "clone":
			handleRepoGet(repoName, remainingArgs)
		default:
			// Treat as git command
			handleRepoGitCommand(repoName, action, remainingArgs)
		}
	},
}

// handleRepoAdd handles the repo <repo-name> add <repo-url> command
func handleRepoAdd(repoName string, args []string) {
	if len(args) != 1 {
		fmt.Printf("Error: 'repo %s add' requires exactly one URL argument\n", repoName)
		fmt.Printf("Usage: mess repo %s add <repo-url>\n", repoName)
		os.Exit(1)
	}

	repoURL := args[0]

	// Load existing configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Validate that repo name doesn't already exist
	for _, repo := range cfg.Repos {
		if repo.Name == repoName {
			fmt.Printf("Repository '%s' already exists\n", repoName)
			os.Exit(1)
		}
	}

	// Validate that repo URL doesn't already exist
	for _, repo := range cfg.Repos {
		if repo.URL == repoURL {
			fmt.Printf("Repository URL '%s' already exists for repo '%s'\n", repoURL, repo.Name)
			os.Exit(1)
		}
	}

	// Add new repository
	newRepo := config.RepoDefinition{
		Name: repoName,
		URL:  repoURL,
	}
	cfg.Repos = append(cfg.Repos, newRepo)

	// Save configuration
	configPath := configFile
	if configPath == "" {
		configPath = "mess.json"
	}
	if err := config.SaveConfig(cfg, configPath); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully added repository '%s' with URL '%s'\n", repoName, repoURL)
}

// handleRepoRemove handles the repo <repo-name> remove command
func handleRepoRemove(repoName string, args []string) {
	if len(args) > 0 {
		fmt.Printf("Error: 'repo %s remove' takes no additional arguments\n", repoName)
		os.Exit(1)
	}

	// Load existing configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Check if repository exists
	repoIndex := -1
	for i, repo := range cfg.Repos {
		if repo.Name == repoName {
			repoIndex = i
			break
		}
	}

	if repoIndex == -1 {
		fmt.Printf("Repository '%s' not found\n", repoName)
		os.Exit(1)
	}

	// Check if repository is used in any applications
	var appsUsingRepo []string
	for _, app := range cfg.Applications {
		for _, appRepo := range app.Repos {
			if appRepo == repoName {
				appsUsingRepo = append(appsUsingRepo, app.Name)
				break
			}
		}
	}

	// If repository is used in applications, prompt for confirmation
	if len(appsUsingRepo) > 0 {
		fmt.Printf("Repository '%s' is used in the following applications:\n", repoName)
		for _, appName := range appsUsingRepo {
			fmt.Printf("  - %s\n", appName)
		}
		fmt.Print("Do you want to remove it from all applications and delete the repository? (y/N): ")
		
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)
		
		if response != "y" && response != "Y" {
			fmt.Println("Repository removal cancelled.")
			return
		}

		// Remove repository from all applications
		for i := range cfg.Applications {
			var newRepos []string
			for _, appRepo := range cfg.Applications[i].Repos {
				if appRepo != repoName {
					newRepos = append(newRepos, appRepo)
				}
			}
			cfg.Applications[i].Repos = newRepos
		}
	} else {
		// Simple confirmation for repository deletion
		fmt.Printf("Are you sure you want to remove repository '%s'? (y/N): ", repoName)
		
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)
		
		if response != "y" && response != "Y" {
			fmt.Println("Repository removal cancelled.")
			return
		}
	}

	// Remove repository from config
	cfg.Repos = append(cfg.Repos[:repoIndex], cfg.Repos[repoIndex+1:]...)

	// Save configuration
	configPath := configFile
	if configPath == "" {
		configPath = "mess.json"
	}
	if err := config.SaveConfig(cfg, configPath); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully removed repository '%s'\n", repoName)
}

// handleRepoGet handles the repo <repo-name> get command
func handleRepoGet(repoName string, args []string) {
	if len(args) > 0 {
		fmt.Printf("Error: 'repo %s get' takes no additional arguments\n", repoName)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Find repository
	var targetRepo *config.RepoDefinition
	for _, repo := range cfg.Repos {
		if repo.Name == repoName {
			targetRepo = &repo
			break
		}
	}

	if targetRepo == nil {
		fmt.Printf("Repository '%s' not found in configuration\n", repoName)
		os.Exit(1)
	}

	// Get config file directory
	configPath := configFile
	if configPath == "" {
		configPath = "mess.json"
	}

	// Clone repository
	if err := repo.CloneRepository(targetRepo, configPath); err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully cloned repository '%s'\n", repoName)
}

// handleRepoGitCommand handles git command delegation
func handleRepoGitCommand(repoName string, gitCommand string, args []string) {
	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Find repository
	var targetRepo *config.RepoDefinition
	for _, repo := range cfg.Repos {
		if repo.Name == repoName {
			targetRepo = &repo
			break
		}
	}

	if targetRepo == nil {
		fmt.Printf("Repository '%s' not found in configuration\n", repoName)
		os.Exit(1)
	}

	// Get config file directory
	configPath := configFile
	if configPath == "" {
		configPath = "mess.json"
	}

	// Check if repository is cloned
	if !repo.IsRepositoryCloned(repoName, configPath) {
		fmt.Printf("Repository '%s' is not cloned. Run 'mess repo %s get' first\n", repoName, repoName)
		os.Exit(1)
	}

	// Get repository path
	repoPath := repo.GetRepositoryPath(repoName, configPath)

	// Build git command
	gitArgs := []string{gitCommand}
	gitArgs = append(gitArgs, args...)

	fmt.Printf("Executing: git %s in %s\n", strings.Join(gitArgs, " "), repoPath)

	// Execute git command
	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = repoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("Git command failed: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(repoCmd)
} 