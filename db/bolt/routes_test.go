package bolt

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

func TestRoutes(t *testing.T) {
	app := &Module{}
	stop := service.New(app).StartForTest()
	defer stop()

	a := &models.Route{
		Uuid: proto.String("a4f5b3cf-5761-4e94-854d-e28f93777bf8"),
		Path: proto.String("test/a"),
		Target: &models.Route_File{
			File: "foo/bar",
		},
	}
	b := &models.Route{
		Uuid: proto.String("bfd3e722-af9f-4344-a8e0-096dffdfc9ac"),
		Path: proto.String("test/b"),
		Target: &models.Route_PageUuid{
			PageUuid: "1234",
		},
	}
	c := &models.Route{
		Path: proto.String("test/c"),
		Target: &models.Route_PageUuid{
			PageUuid: "newroute",
		},
	}
	assert.NoError(t, app.UpdateRoute(a))
	assert.NoError(t, app.UpdateRoute(b))
	assert.NoError(t, app.UpdateRoute(c))
	assert.NotNil(t, c.Uuid)

	routes, err := app.ListRoutes(nil)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []*models.Route{a, b, c}, routes)

	route, err := app.GetRoute(a.GetUuid())
	assert.NoError(t, err)
	assert.Equal(t, a, route)

	route, err = app.GetRoute(b.GetUuid())
	assert.NoError(t, err)
	assert.Equal(t, b, route)

	route, err = app.GetRoute(c.GetUuid())
	assert.NoError(t, err)
	assert.Equal(t, c, route)

	err = app.DeleteRoute(b)
	assert.NoError(t, err)

	route, err = app.GetRoute(b.GetUuid())
	expectedErr := ErrNoKey(ROUTE_BUCKET + ":" + b.GetUuid())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, route)
}
