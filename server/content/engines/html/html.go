package html

import (
	"html/template"
	"io"
	"path"
	"strings"

	"github.com/octavore/nagax/util/errors"

	"github.com/ketchuphq/ketchup/server/content/context"
	"github.com/ketchuphq/ketchup/server/content/engines"
	"github.com/ketchuphq/ketchup/server/content/engines/enginebase"
	"github.com/ketchuphq/ketchup/server/content/templates/store"
)

const EngineTypeHTML = "html"

// additional engines can be registered via go plugins
func init() {
	engines.RegisterEngine(EngineTypeHTML, NewHTMLEngine)
}

type htmlEngine struct {
	tmpl *template.Template
}

func (h *htmlEngine) Execute(w io.Writer, templateName string, context *context.EngineContext) error {
	err := h.tmpl.ExecuteTemplate(w, templateName, context)
	if err != nil {
		return errors.New("engines: error in html template %q", err)
	}
	return nil
}

func NewHTMLEngine(theme store.Theme) (enginebase.Engine, error) {
	var tmpl *template.Template
	funcMap := enginebase.FuncMap()
	for _, t := range theme.Proto().GetTemplates() {
		ext := strings.TrimLeft(path.Ext(t.GetName()), ".")
		if t.GetEngine() == "" && ext != EngineTypeHTML {
			continue
		}
		if t.GetEngine() != EngineTypeHTML {
			continue
		}

		t, err := theme.GetTemplate(t.GetName())
		if err != nil {
			return nil, err
		}
		if tmpl == nil {
			tmpl = template.New(t.GetName()).Funcs(funcMap)
		} else {
			tmpl = tmpl.New(t.GetName()).Funcs(funcMap)
		}
		tmpl, err = tmpl.Parse(t.GetData())
		if err != nil {
			return nil, err
		}
	}

	if tmpl == nil {
		return &htmlEngine{}, nil
	}

	return &htmlEngine{tmpl: tmpl}, nil
}
