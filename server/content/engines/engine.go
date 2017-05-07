package engines

import (
	"errors"
	"io"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/content/context"
)

// Engines contain maps extensions to template rendering engines
var engines = map[string]Engine{}

// Engine has an ExecuteTemplate method which renders data into a template
type Engine interface {
	Execute(w io.Writer, templateData string, context *context.EngineContext) error
}

// additional engines can be registered via go plugins
func init() {
	RegisterEngine(EngineTypeHTML, &htmlEngine{})
}

// RegisterEngine allows external go plugins to register new engines
func RegisterEngine(name string, e Engine) {
	engines[name] = e
}

// Render the given template using the engine it specifies and the given vars.
// Output is written to w.
func Render(w io.Writer, template *models.ThemeTemplate, context *context.EngineContext) error {
	engine := engines[template.GetEngine()]
	if engine == nil {
		return errors.New("unknown template engine " + template.GetEngine())
	}
	return engine.Execute(w, template.GetData(), context)
}
