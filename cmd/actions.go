package cmd

import (
	"os"

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

			path := helper.ExtractPath(file)
			if err := os.WriteFile(path+"/README.md", []byte(a.Markdown()), 0644); err != nil {
				log.Fatal(err)
			}
		}
	},
}
