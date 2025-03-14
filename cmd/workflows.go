package cmd

import (
	"os"
	"strings"

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
		var toc = &markdown.List{
			Items: []string{},
		}
		var markdownBody = ""
		files, err := helper.ScanPattern(".github/workflows", ".yml", false)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			w := workflow.Parse(file, log)
			markdownBody += w.Markdown()

			link := markdown.Hyperlink{
				Text: file,
				URL:  "#" + strings.Replace(w.Name, " ", "-", -1),
			}
			toc.Add(link.String())
		}

		markdownHeader := (&markdown.Markdown{
			Elements: []markdown.Element{
				markdown.H1("Workflows"),
				markdown.P("Table of contents:"),
				toc,
			},
		}).String()

		if err := os.WriteFile(workflowsOutput+"/README.md", []byte(markdownHeader+markdownBody), 0644); err != nil {
			log.Fatal(err)
		}
	},
}
