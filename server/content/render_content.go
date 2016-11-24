package content

import (
	"errors"

	"html/template"

	"github.com/octavore/press/proto/press/models"
)

type contentMap map[string]interface{}

var contentRenders = map[string]func(string) (interface{}, error){
	models.ContentType_html.String(): renderHTMLContent,
	// models.ContentType_markdown.String(): renderHTMLContent,
}

func renderPageContents(page *models.Page) (contentMap, error) {
	contents := map[string]interface{}{}
	var err error
	for _, c := range page.Contents {
		renderer := contentRenders[c.GetContentType().String()]
		if renderer == nil {
			return nil, errors.New("unknown renderer")
		}
		contents[c.GetKey()], err = renderer(c.GetValue())
		if err != nil {
			return nil, err
		}
	}
	return contents, err
}

func renderHTMLContent(s string) (interface{}, error) {
	return template.HTML(s), nil
}
