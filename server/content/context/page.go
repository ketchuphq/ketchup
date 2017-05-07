package context

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
