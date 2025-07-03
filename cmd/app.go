package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"mess/pkg/app"
	"mess/pkg/config"
)

// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:     "app",
	Aliases: []string{"application"},
	Short:   "Manage applications",
	Long: `Manage applications in your mess.json file.
Usage patterns:
  mess app <application-name> init                          - Create a new application
  mess app <application-name> link <repo-name> [...repo-name] - Link repositories to application  
  mess app <application-name> setup                        - Setup application (clone repos, create symlinks, run setup scripts)
  mess app <application-name> clone                        - Clone application repositories and create symlinks
  mess app <application-name> run <script-name>            - Run a script for the application`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Error: insufficient arguments")
			fmt.Println("Usage:")
			fmt.Println("  mess app <application-name> init")
			fmt.Println("  mess app <application-name> link <repo-name> [...repo-name]")
			fmt.Println("  mess app <application-name> setup")
			fmt.Println("  mess app <application-name> clone")
			fmt.Println("  mess app <application-name> run <script-name>")
			os.Exit(1)
		}

		appName := args[0]
		subCommand := args[1]
		remainingArgs := args[2:]

		switch subCommand {
		case "init":
			handleAppInit(appName, remainingArgs)
		case "link":
			handleAppLink(appName, remainingArgs)
		case "setup":
			handleAppSetup(appName, remainingArgs)
		case "clone":
			handleAppClone(appName, remainingArgs)
		case "run":
			handleAppRun(appName, remainingArgs)
		default:
			fmt.Printf("Error: unknown subcommand '%s'\n", subCommand)
			fmt.Println("Available subcommands: init, link, setup, clone, run")
			os.Exit(1)
		}
	},
}

// handleAppInit handles the app <application-name> init command
func handleAppInit(appName string, args []string) {
	if len(args) > 0 {
		fmt.Printf("Error: 'app %s init' takes no additional arguments\n", appName)
		os.Exit(1)
	}

	// Load existing configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Validate that application name doesn't already exist
	for _, app := range cfg.Applications {
		if app.Name == appName {
			fmt.Printf("Application '%s' already exists\n", appName)
			os.Exit(1)
		}
	}

	// Add new application
	newApp := config.ApplicationDefinition{
		Name:    appName,
		Repos:   []string{},
		Scripts: make(map[string]config.ScriptValue),
		Env:     make(map[string]string),
	}
	cfg.Applications = append(cfg.Applications, newApp)

	// Save configuration
	configPath := configFile
	if configPath == "" {
		configPath = "mess.json"
	}
	if err := config.SaveConfig(cfg, configPath); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created application '%s'\n", appName)
}

// handleAppLink handles the app <application-name> link <repo-name> [...repo-name] command
func handleAppLink(appName string, args []string) {
	if len(args) < 1 {
		fmt.Printf("Error: 'app %s link' requires at least one repository name\n", appName)
		fmt.Printf("Usage: mess app %s link <repo-name> [...repo-name]\n", appName)
		os.Exit(1)
	}

	repoNames := args

	// Load existing configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Find application
	appIndex := -1
	for i, app := range cfg.Applications {
		if app.Name == appName {
			appIndex = i
			break
		}
	}

	if appIndex == -1 {
		fmt.Printf("Application '%s' not found\n", appName)
		fmt.Printf("Available applications:\n")
		for _, app := range cfg.Applications {
			fmt.Printf("  - %s\n", app.Name)
		}
		os.Exit(1)
	}

	// Validate all repositories exist
	var validRepos []string
	for _, repoName := range repoNames {
		repoExists := false
		for _, repo := range cfg.Repos {
			if repo.Name == repoName {
				repoExists = true
				break
			}
		}

		if !repoExists {
			fmt.Printf("Repository '%s' not found\n", repoName)
			fmt.Printf("Available repositories:\n")
			for _, repo := range cfg.Repos {
				fmt.Printf("  - %s\n", repo.Name)
			}
			os.Exit(1)
		}

		// Check if repository is already linked
		alreadyLinked := false
		for _, linkedRepo := range cfg.Applications[appIndex].Repos {
			if linkedRepo == repoName {
				fmt.Printf("Repository '%s' is already linked to application '%s', skipping\n", repoName, appName)
				alreadyLinked = true
				break
			}
		}

		if !alreadyLinked {
			validRepos = append(validRepos, repoName)
		}
	}

	if len(validRepos) == 0 {
		fmt.Println("No new repositories to link")
		return
	}

	// Link repositories to application
	cfg.Applications[appIndex].Repos = append(cfg.Applications[appIndex].Repos, validRepos...)

	// Save configuration
	configPath := configFile
	if configPath == "" {
		configPath = "mess.json"
	}
	if err := config.SaveConfig(cfg, configPath); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	if len(validRepos) == 1 {
		fmt.Printf("Successfully linked repository '%s' to application '%s'\n", validRepos[0], appName)
	} else {
		fmt.Printf("Successfully linked %d repositories to application '%s': %s\n", len(validRepos), appName, strings.Join(validRepos, ", "))
	}
}

// handleAppSetup handles the app <application-name> setup command
func handleAppSetup(appName string, args []string) {
	if len(args) > 0 {
		fmt.Printf("Error: 'app %s setup' takes no additional arguments\n", appName)
		os.Exit(1)
	}

	// Load existing configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Find application
	var targetApp *config.ApplicationDefinition
	for _, app := range cfg.Applications {
		if app.Name == appName {
			targetApp = &app
			break
		}
	}

	if targetApp == nil {
		fmt.Printf("Application '%s' not found\n", appName)
		fmt.Printf("Available applications:\n")
		for _, app := range cfg.Applications {
			fmt.Printf("  - %s\n", app.Name)
		}
		os.Exit(1)
	}

	// Get config file directory
	configPath := configFile
	if configPath == "" {
		configPath = "mess.json"
	}

	// Setup application
	if err := app.SetupApplication(targetApp, cfg, configPath); err != nil {
		fmt.Printf("Error setting up application: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully setup application '%s'\n", appName)
}

// handleAppClone handles the app <application-name> clone command
func handleAppClone(appName string, args []string) {
	if len(args) > 0 {
		fmt.Printf("Error: 'app %s clone' takes no additional arguments\n", appName)
		os.Exit(1)
	}

	// Load existing configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Find application
	var targetApp *config.ApplicationDefinition
	for _, app := range cfg.Applications {
		if app.Name == appName {
			targetApp = &app
			break
		}
	}

	if targetApp == nil {
		fmt.Printf("Application '%s' not found\n", appName)
		fmt.Printf("Available applications:\n")
		for _, app := range cfg.Applications {
			fmt.Printf("  - %s\n", app.Name)
		}
		os.Exit(1)
	}

	// Get config file directory
	configPath := configFile
	if configPath == "" {
		configPath = "mess.json"
	}

	// Clone application repositories
	if err := app.CloneApplication(targetApp, cfg, configPath); err != nil {
		fmt.Printf("Error cloning application: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully cloned application '%s'\n", appName)
}

// handleAppRun handles the app <application-name> run <script-name> command
func handleAppRun(appName string, args []string) {
	if len(args) != 1 {
		fmt.Printf("Error: 'app %s run' requires exactly one script name\n", appName)
		fmt.Printf("Usage: mess app %s run <script-name>\n", appName)
		os.Exit(1)
	}

	scriptName := args[0]

	// Load existing configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Find application
	var targetApp *config.ApplicationDefinition
	for _, app := range cfg.Applications {
		if app.Name == appName {
			targetApp = &app
			break
		}
	}

	if targetApp == nil {
		fmt.Printf("Application '%s' not found\n", appName)
		fmt.Printf("Available applications:\n")
		for _, app := range cfg.Applications {
			fmt.Printf("  - %s\n", app.Name)
		}
		os.Exit(1)
	}

	// Check if script exists
	scriptValue, exists := targetApp.Scripts[scriptName]
	if !exists {
		fmt.Printf("Script '%s' not found in application '%s'\n", scriptName, appName)
		fmt.Printf("Available scripts:\n")
		for name := range targetApp.Scripts {
			fmt.Printf("  - %s\n", name)
		}
		os.Exit(1)
	}

	// Get config file directory
	configPath := configFile
	if configPath == "" {
		configPath = "mess.json"
	}

	// Run script
	if err := app.RunScript(targetApp, scriptName, &scriptValue, configPath); err != nil {
		fmt.Printf("Error running script: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully executed script '%s' for application '%s'\n", scriptName, appName)
}

func init() {
	rootCmd.AddCommand(appCmd)
} 