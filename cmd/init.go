package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"mess/pkg/config"
)

var projectName string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new mess.json file",
	Long: `Initialize a new empty mess.json file.
This will create an empty mess.json file in the current directory that you can 
customize for your project by adding repositories and applications.
You can specify a project name using -n/--name flag, otherwise the current directory name will be used.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := "mess.json"
		if configFile != "" {
			configPath = configFile
		}

		// Check if config file already exists - return error if it does (SRS 8.3)
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("Error: Project is already initialized. Config file exists: %s\n", configPath)
			fmt.Println("Use a different directory or remove the existing mess.json file first.")
			os.Exit(1)
		}

		// Determine project name
		var finalProjectName string
		if projectName != "" {
			finalProjectName = projectName
		} else {
			// Use current directory name as fallback
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				os.Exit(1)
			}
			finalProjectName = filepath.Base(currentDir)
		}

		// Create empty configuration with project name
		emptyConfig := &config.MessConfig{
			Name:         finalProjectName,
			Repos:        []config.RepoDefinition{},
			Applications: []config.ApplicationDefinition{},
		}

		// Save configuration
		if err := config.SaveConfig(emptyConfig, configPath); err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created %s with project name '%s'\n", configPath, finalProjectName)
		fmt.Println("You can now customize the configuration for your project.")
		fmt.Println("Use 'mess repo <name> add <url>' to add repositories")
		fmt.Println("Use 'mess app <name> init' to add applications")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&projectName, "name", "n", "", "project name (default is current directory name)")
} 