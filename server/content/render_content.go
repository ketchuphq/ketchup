package content

import (
	"fmt"

	"github.com/octavore/press/proto/press/models"
)

type contentMap map[string]interface{}

var renderers = map[string]Content{
	models.Content_html.String(): &HTMLContent{},
}

func renderPageContents(page *models.Page) (contentMap, error) {
	contents := map[string]interface{}{}
	var err error
	for _, c := range page.Contents {
		rt := c.GetContentType().String()
		renderer := renderers[rt]
		if renderer == nil {
			return nil, fmt.Errorf("contents: unknown renderer %q", rt)
		}
		contents[c.GetKey()], err = renderer.Render(c.GetValue())
		if err != nil {
			return nil, err
		}
	}
	return contents, err
}

