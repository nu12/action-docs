package cmd

import (
	"os"

	"github.com/nu12/action-docs/internal/action"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var actionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "Generate documentation for github actions",
	Long:  `Generate documentation for github actions`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Scanning actions")
		scanActions(actionsPath)
	},
}

func scanActions(path string) {
	log.Debug("Scanning " + path)
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			if file.Name() == ".git" || file.Name() == ".github" {
				continue
			}
			scanActions(path + "/" + file.Name())
			continue
		}
		if file.Name() != "action.yml" {
			continue
		}

		b, err := os.ReadFile(path + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		a := action.Action{}

		err = yaml.Unmarshal([]byte(b), &a)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(path+"/README.md", []byte(a.Markdown()), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
