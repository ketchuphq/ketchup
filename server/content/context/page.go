package context

import "time"

// PageContext wraps page related data
type PageContext struct {
	*EngineContext
}

// Data returns rendered page data for the given key
func (p *PageContext) Data(key string) interface{} {
	return p.contents[key]
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
