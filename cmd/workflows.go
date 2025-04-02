package cmd

import (
	"os"

	"github.com/nu12/action-docs/internal/helper"
	"github.com/nu12/action-docs/internal/markdown"
	"github.com/nu12/action-docs/internal/workflow"
	"github.com/spf13/cobra"
)

var workflowsCmd = &cobra.Command{
	Use:   "workflows",
	Short: "Generate documentation for github workflows",
	Long:  `Generate documentation for github workflows`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Scanning workflows")
		var ws = workflow.Workflows{
			Workflows: []workflow.Workflow{},
			Content:   markdown.List{},
		}

		files, err := helper.ScanPattern(".github/workflows", ".yml", false)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			w := workflow.Parse(file, log)
			ws.AddWorkflow(w)
		}

		if err := os.WriteFile(workflowsOutput+"/README.md", []byte(ws.Markdown()), 0644); err != nil {
			log.Fatal(err)
		}
	},
}
