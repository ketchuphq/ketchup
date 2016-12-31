package content

import (
	"errors"
	"html/template"
	"io"

	"github.com/russross/blackfriday"

	"github.com/octavore/press/proto/press/models"
)

func (m *Module) render(w io.Writer, page *models.Page) error {
	contents, err := m.createContentMap(page)
	if err != nil {
		return err
	}
	return m.Templates.Render(w, page, contents)
}

type contentMap map[string]interface{}

func (m *Module) createContentMap(page *models.Page) (contentMap, error) {
	contents := map[string]interface{}{}
	var err error
	for _, c := range page.Contents {
		switch c.GetType().(type) {
		case *models.Content_Short:
			contents[c.GetKey()], err = renderTextualContent(c.GetValue(), c.GetShort())
		case *models.Content_Text:
			contents[c.GetKey()], err = renderTextualContent(c.GetValue(), c.GetText())
		case *models.Content_Multiple:
			contents[c.GetKey()] = c.GetValue()
		default:
			m.Logger.Warningf("unknown content type: %s", c.GetType())
		}
		// future: allow custom registere content type renderers
		if err != nil {
			return nil, err
		}
	}
	return contents, err
}

type textualContent interface {
	GetType() models.ContentTextType
}

func renderTextualContent(s string, t textualContent) (interface{}, error) {
	switch t.GetType() {
	case models.ContentTextType_text:
		return s, nil
	case models.ContentTextType_html:
		return template.HTML(s), nil
	case models.ContentTextType_markdown:
		data := blackfriday.MarkdownCommon([]byte(s))
		return template.HTML(string(data)), nil
	}
	return nil, errors.New("unknown content text type: " + t.GetType().String())
}
