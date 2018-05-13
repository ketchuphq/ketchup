package templates

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/content/context"
)

func TestRenderPage(t *testing.T) {
	m, stop := setup(false, testTheme)
	defer stop()

	buf := &bytes.Buffer{}
	page := &models.Page{
		Theme:    testTheme.Name,
		Template: testTemplate.Name,
	}
	contents := map[string]interface{}{"content": "hello"}
	m.RenderPage(buf, page, context.NewContext(nil, page, nil, contents))
	assert.Equal(t, "<div>hello</div>", buf.String())
}
