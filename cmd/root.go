/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/nu12/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var actionsPath string
var cfgFile string
var workflowsOutput string

var log = logging.NewLogger()

var rootCmd = &cobra.Command{
	Use:   "action-docs",
	Short: "Create documentation for github actions and workflows",
	Long:  `Create documentation for github actions and workflows`,
	//Run: func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(actionsCmd)
	rootCmd.AddCommand(workflowsCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.action-docs.yaml)")

	actionsCmd.Flags().StringVarP(&actionsPath, "path", "p", ".", "Path to the directory containing github actions to be scanned")
	workflowsCmd.Flags().StringVarP(&workflowsOutput, "output", "o", ".github/workflows", "Path to place the documentation for workflows")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".action-docs" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".action-docs")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
