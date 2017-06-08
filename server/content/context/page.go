package context

import (
	"sort"
	"time"

	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

// PageContext wraps page related data
type PageContext struct {
	engine *EngineContext
	page   *models.Page
}

// Data returns rendered page data for the given key
func (p *PageContext) Data(key string) interface{} {
	return p.engine.contents[key]
}

// Title is shorthand for accessing page title
func (p *PageContext) Title() string {
	return p.page.GetTitle()
}

// Content is shorthand for accessing page content.
func (p *PageContext) Content() interface{} {
	return p.Data("content")
}

// CreatedAt is shorthand for accessing page created at time
func (p *PageContext) CreatedAt() time.Time {
	createdMillis := p.page.GetTimestamps().GetCreatedAt()
	return time.Unix(createdMillis/1000, 0)
}

// UpdatedAt is shorthand for accessing page updated at time
func (p *PageContext) UpdatedAt() time.Time {
	updatedMillis := p.page.GetTimestamps().GetUpdatedAt()
	return time.Unix(updatedMillis/1000, 0)
}

// PublishedAt is shorthand for accessing page published at time
func (p *PageContext) PublishedAt() time.Time {
	publishedMillis := p.page.GetPublishedAt()
	return time.Unix(publishedMillis/1000, 0)
}

func (p *PageContext) Theme() string {
	return p.page.GetTheme()
}

func (p *PageContext) Template() string {
	return p.page.GetTemplate()
}

// Route returns the first route
func (p *PageContext) Route() string {
	routes, err := p.engine.backend.ListRoutes(&api.ListRouteRequest{
		Options: &api.ListRouteRequest_ListRouteOptions{
			PageUuid: p.page.Uuid,
		},
	})
	if err != nil {
		p.engine.logger.Error(err)
		return ""
	}
	if len(routes) == 0 {
		return ""
	}
	return routes[0].GetPath()
}

type SortablePages []*PageContext

func (p SortablePages) ByPublishedAt() SortablePages {
	sort.Slice(p, func(i, j int) bool {
		return p[i].page.GetPublishedAt() < p[j].page.GetPublishedAt()
	})
	return p
}

func (p SortablePages) Reverse() SortablePages {
	for i := 0; i < len(p)/2; i++ {
		j := len(p) - 1 - i
		p[i], p[j] = p[j], p[i]
	}
	return p
}
