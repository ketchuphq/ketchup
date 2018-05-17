package content

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

type renderTextContentCase struct {
	t models.ContentTextType
	c string

	expected interface{}
}

func (r renderTextContentCase) GetType() models.ContentTextType {
	return r.t
}

func TestRenderTextContent(t *testing.T) {
	cases := []renderTextContentCase{
		{
			t:        models.ContentTextType_text,
			c:        "<div>**Hello**</div>",
			expected: "<div>**Hello**</div>",
		},
		{
			t:        models.ContentTextType_html,
			c:        "<div>**Hello**</div>",
			expected: template.HTML("<div>**Hello**</div>"),
		},
		{
			t:        models.ContentTextType_markdown,
			c:        "**hello**",
			expected: template.HTML("<p><strong>hello</strong></p>\n"),
		},
	}
	for _, k := range cases {
		v, err := renderTextualContent(k.c, k)
		if assert.NoError(t, err) {
			assert.Equal(t, v, k.expected)
		}
	}
}

func TestRenderContent(t *testing.T) {
	cases := []struct {
		t           models.IsContent_Type
		c, expected string
	}{
		{
			t:        &models.Content_Short{},
			c:        "<div>**Hello**</div>",
			expected: "<div>**Hello**</div>",
		},
		{
			t:        &models.Content_Text{},
			c:        "<div>**Hello**</div>",
			expected: "<div>**Hello**</div>",
		},
		{
			t:        &models.Content_Multiple{},
			c:        "**hello**",
			expected: "**hello**",
		},
	}
	for _, k := range cases {
		v, err := RenderContent(&models.Content{
			Value: &k.c,
			Type:  k.t,
		})
		if assert.NoError(t, err) {
			assert.Equal(t, k.expected, v)
		}
	}
}

func TestRenderData(t *testing.T) {
	cases := []struct {
		t           models.IsData_Type
		c, expected string
	}{
		{
			t:        &models.Data_Short{},
			c:        "<div>**Hello**</div>",
			expected: "<div>**Hello**</div>",
		},
		{
			t:        &models.Data_Text{},
			c:        "<div>**Hello**</div>",
			expected: "<div>**Hello**</div>",
		},
		{
			t:        &models.Data_Multiple{},
			c:        "**hello**",
			expected: "**hello**",
		},
	}
	for _, k := range cases {
		v, err := RenderData(&models.Data{
			Value: &k.c,
			Type:  k.t,
		})
		if assert.NoError(t, err) {
			assert.Equal(t, k.expected, v)
		}
	}
}
