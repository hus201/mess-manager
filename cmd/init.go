package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mess/pkg/config"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new mess.json file",
	Long: `Initialize a new mess.json file with sample configuration.
This will create a mess.json file in the current directory with example 
repositories and applications to help you get started.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := "mess.json"
		if configFile != "" {
			configPath = configFile
		}

		// Check if config file already exists
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("Config file already exists: %s\n", configPath)
			fmt.Print("Do you want to overwrite it? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Initialization cancelled.")
				return
			}
		}

		// Create sample configuration
		sampleConfig := config.CreateSampleConfig()

		// Save configuration
		if err := config.SaveConfig(sampleConfig, configPath); err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created %s\n", configPath)
		fmt.Println("You can now customize the configuration for your project.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
} 