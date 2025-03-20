package workflow

import (
	"github.com/nu12/action-docs/internal/helper"
	"github.com/nu12/action-docs/internal/markdown"
)

type Workflows struct {
	Workflows []Workflow
	Content   markdown.List
}

func (w *Workflows) AddWorkflow(workflow *Workflow) *Workflows {
	w.Workflows = append(w.Workflows, *workflow)

	link := markdown.Hyperlink{
		Text: workflow.Filename,
		URL:  "#" + helper.SanitizeURL(workflow.Name),
	}
	w.Content.Add(link.String())
	return w
}

func (w *Workflows) String() string {
	var body = ""
	for _, workflow := range w.Workflows {
		body += workflow.Markdown()
	}
	header := (&markdown.Markdown{
		Elements: []markdown.Element{
			markdown.H1("Workflows"),
			markdown.P("Table of contents:"),
			&w.Content,
		},
	}).String()
	return header + body
}
