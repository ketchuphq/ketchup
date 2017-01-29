package engines

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/octavore/ketchup/proto/ketchup/models"
)

func TestRenderTemplate(t *testing.T) {
	type args struct {
		template *models.ThemeTemplate
		contents map[string]interface{}
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
					Engine: proto.String("html"),
					Data:   proto.String("<h1>hello {{.world}}</h1>"),
				},
				map[string]interface{}{
					"world": "earth",
				},
			},
			"<h1>hello earth</h1>",
			false,
		},
	}
	for _, tt := range tests {
		w := &bytes.Buffer{}
		if err := Render(w, tt.args.template, tt.args.contents); (err != nil) != tt.expectedErr {
			t.Errorf("%q. renderTemplate() error = %v, expectedErr %v", tt.name, err, tt.expectedErr)
			continue
		}
		if output := w.String(); output != tt.expectedOutput {
			t.Errorf("%q. renderTemplate() = %v, want %v", tt.name, output, tt.expectedOutput)
		}
	}
}
