package cmd

import (
	"os"
	"strings"

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
		files, err := os.ReadDir(".github/workflows")
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			log.Debug(file.Name())

			b, err := os.ReadFile(".github/workflows/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			if !strings.Contains(file.Name(), ".yml") {
				continue
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
