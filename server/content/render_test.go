package content

import (
	"encoding/json"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

func TestCreateContentMap(t *testing.T) {
	page := &models.Page{
		Contents: []*models.Content{
			{
				Key:   proto.String("k0"),
				Value: proto.String("v0"),
				Type: &models.Content_Short{
					Short: nil,
				},
			},
		},
	}
	m := &Module{}
	cm, err := m.CreateContentMap(page)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"k0":"v0","title":""}` {
		t.Error(string(b))
	}
}
