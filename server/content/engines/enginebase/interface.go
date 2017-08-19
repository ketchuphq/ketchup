package enginebase

import (
	"io"

	"github.com/ketchuphq/ketchup/server/content/context"
	"github.com/ketchuphq/ketchup/server/content/templates/store"
)

// Engine has an ExecuteTemplate method which renders data into a template
type Engine interface {
	Execute(w io.Writer, templateName string, context *context.EngineContext) error
}

type EngineFactory func(store.Theme) (Engine, error)
