package content

import (
	"html/template"

	"github.com/octavore/press/proto/press/models"
)

type Content interface {
	Name() string
	Render(string) (interface{}, error)
}

type HTMLContent struct{}

func (h *HTMLContent) Name() string {
	return models.Content_html.String()
}

func (h *HTMLContent) Render(s string) (interface{}, error) {
	return template.HTML(s), nil
}
