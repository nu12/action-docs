package cmd

import (
	"os"

	"github.com/nu12/action-docs/internal/helper"
	"github.com/nu12/action-docs/internal/workflow"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var workflowsCmd = &cobra.Command{
	Use:   "workflows",
	Short: "Generate documentation for github workflows",
	Long:  `Generate documentation for github workflows`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Scanning workflows")
		var workflowOutput = ""
		files, err := helper.ScanPattern(".github/workflows", ".yml", false)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			b, err := os.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}

			w := workflow.Workflow{}
			err = yaml.Unmarshal([]byte(b), &w)
			if err != nil {
				log.Fatal(err)
			}

			workflowOutput += w.Markdown()
		}
		err = os.WriteFile(workflowsOutput+"/README.md", []byte(workflowOutput), 0644)
		if err != nil {
			log.Fatal(err)
		}
	},
}
