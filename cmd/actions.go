package cmd

import (
	"os"

	"github.com/nu12/action-docs/internal/action"
	"github.com/nu12/action-docs/internal/helper"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

			b, err := os.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}

			a := action.Action{}

			err = yaml.Unmarshal([]byte(b), &a)
			if err != nil {
				log.Fatal(err)
			}

			path := helper.ExtractPath(file)
			err = os.WriteFile(path+"/README.md", []byte(a.Markdown()), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
