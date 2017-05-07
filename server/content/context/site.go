package context

import (
	"github.com/ketchuphq/ketchup/server/content/content"
)

// SiteContext wraps engine context and provides site related functions.
type SiteContext struct {
	*EngineContext
}

// Title is shorthand for accessing site title
func (s *SiteContext) Title() interface{} {
	return s.Data("title")
}

// Data returns rendered site data for the given key
func (s *SiteContext) Data(key string) interface{} {
	data, err := s.backend.GetData(key)
	if err != nil {
		s.logger.Errorf("site data: %v", err)
		return ""
	}
	rendered, err := content.RenderData(data)
	if err != nil {
		s.logger.Errorf("site data: %v", err)
		return ""
	}
	return rendered
}
