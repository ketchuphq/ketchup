package engines_test

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/content/context"
	"github.com/ketchuphq/ketchup/server/content/engines"
	_ "github.com/ketchuphq/ketchup/server/content/engines/html"
	"github.com/ketchuphq/ketchup/server/content/templates/dummystore"
)

func TestRenderTemplate(t *testing.T) {
	type args struct {
		template *models.ThemeTemplate
		contents *context.EngineContext
	}
	tests := []struct {
		name           string
		args           args
		expectedOutput string
		expectedErr    bool
	}{
		{
			"test render html template",
			args{
				&models.ThemeTemplate{
					Name:   proto.String("index.html"),
					Engine: proto.String("html"),
					Data:   proto.String(`<h1>hello {{.Page.Data "world"}}</h1>`),
				},
				context.NewContext(
					nil, nil, nil,
					map[string]interface{}{
						"world": "earth",
					},
				),
			},
			"<h1>hello earth</h1>",
			false,
		},
	}
	for _, tt := range tests {
		w := &bytes.Buffer{}
		tmpl := tt.args.template
		theme := &dummy.Theme{
			Theme: &models.Theme{
				Templates: map[string]*models.ThemeTemplate{tmpl.GetName(): tmpl},
			},
		}
		if err := engines.Render(w, theme, tmpl.GetName(), tt.args.contents); (err != nil) != tt.expectedErr {
			t.Errorf("%q. renderTemplate() error = %v, expectedErr %v", tt.name, err, tt.expectedErr)
			continue
		}
		if output := w.String(); output != tt.expectedOutput {
			t.Errorf("%q. renderTemplate() = %v, want %v", tt.name, output, tt.expectedOutput)
		}
	}
}
