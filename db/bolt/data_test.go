package bolt

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

func TestData(t *testing.T) {
	app := &Module{}
	stop := service.New(app).StartForTest()
	defer stop()

	a := &models.Data{
		Key:   proto.String("cat"),
		Value: proto.String("Sheriff"),
	}
	b := &models.Data{
		Key:   proto.String("dog"),
		Value: proto.String("Doge"),
	}
	expectedData := []*models.Data{a, b}

	assert.NoError(t, app.UpdateData(a))
	assert.NoError(t, app.UpdateData(b))
	assert.NotNil(t, b.Timestamps)
	assert.Nil(t, b.Uuid) // data is kept by key
	assert.EqualError(t, app.UpdateData(&models.Data{
		Value: proto.String("Doge"),
	}), "bolt: cannot update data without key")

	data, err := app.ListData()
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedData, data)
	for _, d := range expectedData {
		data, err := app.GetData(d.GetKey())
		assert.NoError(t, err)
		assert.Equal(t, d, data)
	}

	assert.NoError(t, app.UpdateDataBatch([]*models.Data{
		&models.Data{
			Key:   proto.String("dog"),
			Value: proto.String("Cowboy"),
		},
	}))
	datum, err := app.GetData(a.GetKey())
	assert.NoError(t, err)
	assert.Equal(t, datum, a)

	datum, err = app.GetData(b.GetKey())
	assert.NoError(t, err)
	assert.Equal(t, datum.GetValue(), "Cowboy")

	err = app.DeleteData(b)
	assert.NoError(t, err)
	datum, err = app.GetData(b.GetUuid())
	expectedErr := ErrNoKey(DATA_BUCKET + ":" + b.GetUuid())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, datum)
}
