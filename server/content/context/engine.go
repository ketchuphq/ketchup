package context

import (
	"github.com/octavore/nagax/logger"

	"github.com/ketchuphq/ketchup/db"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

// EngineContext provides functions to templates
type EngineContext struct {
	page     *models.Page
	logger   logger.Logger
	backend  db.Backend
	contents map[string]interface{}
}

// NewContext returns a new EngineContext
func NewContext(l logger.Logger, page *models.Page, backend db.Backend, contents map[string]interface{}) *EngineContext {
	return &EngineContext{
		page:     page,
		logger:   l,
		backend:  backend,
		contents: contents,
	}
}

// Page returns a PageContext
func (e *EngineContext) Page() *PageContext {
	return &PageContext{e, e.page}
}

// Site returns a PageContext
func (e *EngineContext) Site() *SiteContext {
	return &SiteContext{e}
}
