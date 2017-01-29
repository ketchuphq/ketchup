package bolt

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"

	"github.com/octavore/ketchup/proto/ketchup/models"
)

func TestRoutes(t *testing.T) {
	app := &Module{}
	stop := service.New(app).StartForTest()
	defer stop()

	a := &models.Route{
		Uuid: proto.String("a4f5b3cf-5761-4e94-854d-e28f93777bf8"),
		Target: &models.Route_File{
			File: "foo/bar",
		},
	}
	b := &models.Route{
		Uuid: proto.String("bfd3e722-af9f-4344-a8e0-096dffdfc9ac"),
		Target: &models.Route_PageUuid{
			PageUuid: "1234",
		},
	}
	err := app.UpdateRoute(a)
	if err != nil {
		t.Error(err)
	}
	err = app.UpdateRoute(b)
	if err != nil {
		t.Error(err)
	}
	routes, err := app.ListRoutes()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(routes, []*models.Route{a, b}) {
		t.Error(routes)
	}
}
