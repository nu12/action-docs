package cmd

import (
	"os"
	"path/filepath"

	"github.com/nu12/action-docs/internal/action"
	"github.com/nu12/action-docs/internal/helper"
	"github.com/spf13/cobra"
)

var actionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "Generate documentation for github actions",
	Long:  `Generate documentation for github actions`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Scanning actions")

		files, err := helper.ScanPattern(actionsPath, "action.yml", true)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			a := action.Parse(file, log)

			if err := os.WriteFile(filepath.Dir(file)+"/README.md", []byte(a.Markdown()), 0644); err != nil {
				log.Fatal(err)
			}
		}
	},
}
