package engines

import (
	"errors"
	"io"

	"github.com/ketchuphq/ketchup/server/content/context"
	"github.com/ketchuphq/ketchup/server/content/engines/enginebase"
	"github.com/ketchuphq/ketchup/server/content/templates/store"
)

// Engines contain maps extensions to template rendering engines
var engineFactories = map[string]enginebase.EngineFactory{}

// RegisterEngine allows external go plugins to register new engines
func RegisterEngine(name string, f enginebase.EngineFactory) {
	engineFactories[name] = f
}

// Render the given template using the engine it specifies and the given vars.
// Output is written to w.
func Render(w io.Writer, theme store.Theme, templateName string, context *context.EngineContext) error {
	// todo: remove this extra GetTemplate (another one in the factory)
	t, err := theme.GetTemplate(templateName)
	if err != nil {
		return err
	}

	// todo: cache engine for theme
	engineFactory := engineFactories[t.GetEngine()]
	if engineFactory == nil {
		return errors.New("unknown template engine " + t.GetEngine())
	}
	engine, err := engineFactory(theme)
	if err != nil {
		return err
	}
	if engine == nil {
		return errors.New("unknown template engine " + t.GetEngine())
	}
	return engine.Execute(w, templateName, context)
}
