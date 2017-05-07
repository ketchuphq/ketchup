package engines

import (
	"html/template"
	"io"

	"github.com/ketchuphq/ketchup/server/content/context"
	"github.com/ketchuphq/ketchup/util/errors"
)

const EngineTypeHTML = "html"

type htmlEngine struct{}

func (h *htmlEngine) Execute(w io.Writer, templateData string, context *context.EngineContext) (err error) {
	t, err := template.New("temp").Parse(templateData)
	if err != nil {
		return errors.New("engines: error in html template %q", err)
	}
	return t.Execute(w, context)
}
