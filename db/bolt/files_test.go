package bolt

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

func TestFiles(t *testing.T) {
	app := &Module{}
	stop := service.New(app).StartForTest()
	defer stop()

	a := &models.File{
		Uuid: proto.String("a4f5b3cf-5761-4e94-854d-e28f93777bf8"),
		Name: proto.String("cat.jpg"),
		Url:  proto.String("http://example.com/cat.jpg"),
	}
	b := &models.File{
		Name: proto.String("dog.jpg"),
		Url:  proto.String("http://example.com/dog.jpg"),
	}
	expectedFiles := []*models.File{a, b}

	assert.NoError(t, app.UpdateFile(a))
	assert.NoError(t, app.UpdateFile(b))
	assert.NotNil(t, b.Uuid)

	files, err := app.ListFiles()
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedFiles, files)
	for _, r := range expectedFiles {
		file, err := app.GetFile(r.GetUuid())
		assert.NoError(t, err)
		assert.Equal(t, r, file)

		file, err = app.GetFileByName(r.GetName())
		assert.NoError(t, err)
		assert.Equal(t, r, file)
	}

	err = app.DeleteFile(b)
	assert.NoError(t, err)

	file, err := app.GetFile(b.GetUuid())
	expectedErr := ErrNoKey(FILES_BUCKET + ":" + b.GetUuid())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, file)
}
