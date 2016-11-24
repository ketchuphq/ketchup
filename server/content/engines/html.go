package engines

import (
	"fmt"
	"html/template"
	"io"
)

const EngineTypeHTML = "html"

type htmlEngine struct{}

func (h *htmlEngine) Execute(w io.Writer, templateData string, data interface{}) error {
	t, err := template.New("temp").Parse(templateData)
	if err != nil {
		return fmt.Errorf("engines: error in html template %q", err)
	}
	return t.Execute(w, data)
}
