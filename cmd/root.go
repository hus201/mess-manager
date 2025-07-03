package cmd

import (
	"github.com/spf13/cobra"
)

var configFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mess",
	Short: "Mess Manager - A tool to manage messy projects with multiple repos and applications",
	Long: `Mess Manager (mess) is a CLI tool that helps developers manage messy projects 
with multiple repositories and applications. It's useful for:

- Systems with separated frontend and backend applications
- Micro-services backends  
- Modular systems

Use mess.json file to define your project structure and manage repositories and applications.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&configFile, "file", "f", "", "config file path (default is mess.json in current directory)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
} 