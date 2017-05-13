package engines

import (
	"html/template"
	"io"

	"github.com/Masterminds/sprig"

	"github.com/ketchuphq/ketchup/server/content/context"
	"github.com/ketchuphq/ketchup/util/errors"
)

var sprigWhitelist = []string{
	"date",
	"dateModify",
	"dateInZone",
	"now",
}

var funcMap = template.FuncMap{}

func init() {
	fm := sprig.GenericFuncMap()
	for _, k := range sprigWhitelist {
		funcMap[k] = fm[k]
	}
}

const EngineTypeHTML = "html"

type htmlEngine struct{}

func (h *htmlEngine) Execute(w io.Writer, templateData string, context *context.EngineContext) (err error) {
	t, err := template.New("temp").Funcs(funcMap).Parse(templateData)
	if err != nil {
		return errors.New("engines: error in html template %q", err)
	}
	return t.Execute(w, context)
}
