package content

import (
	"fmt"
	"html/template"

	"github.com/octavore/nagax/util/errors"
	"github.com/russross/blackfriday"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

type ErrUnknownContentType struct {
	msg string
}

func (e ErrUnknownContentType) Error() string {
	return e.msg
}

func RenderContent(c *models.Content) (interface{}, error) {
	// future: allow custom registered content type renderers
	switch c.GetType().(type) {
	case *models.Content_Short:
		return renderTextualContent(c.GetValue(), c.GetShort())
	case *models.Content_Text:
		return renderTextualContent(c.GetValue(), c.GetText())
	case *models.Content_Multiple:
		return c.GetValue(), nil
	}
	return nil, errors.Wrap(ErrUnknownContentType{
		fmt.Sprintf("unknown content type: %v", c),
	})
}

func RenderData(c *models.Data) (interface{}, error) {
	// ugh not DRY
	switch c.GetType().(type) {
	case *models.Data_Short:
		return renderTextualContent(c.GetValue(), c.GetShort())
	case *models.Data_Text:
		return renderTextualContent(c.GetValue(), c.GetText())
	case *models.Data_Multiple:
		return c.GetValue(), nil
	}
	return nil, errors.Wrap(ErrUnknownContentType{
		fmt.Sprintf("unknown content type: %v", c),
	})
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
	return nil, errors.Wrap(ErrUnknownContentType{
		fmt.Sprintf("unknown content text type: %s", t.GetType()),
	})
}
