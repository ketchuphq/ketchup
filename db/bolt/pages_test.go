package bolt

import (
	"testing"
	"time"

	"github.com/ketchuphq/ketchup/proto/ketchup/api"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

func TestPages(t *testing.T) {
	app := &Module{}
	stop := service.New(app).StartForTest()
	defer stop()

	a := &models.Page{
		Uuid:        proto.String("a4f5b3cf-5761-4e94-854d-e28f93777bf8"),
		Title:       proto.String("1st page"),
		Theme:       proto.String("san-marzano"),
		Template:    proto.String("page.html"),
		PublishedAt: proto.Int64(time.Now().Unix()),
		Contents: []*models.Content{
			&models.Content{Key: proto.String("content"), Value: proto.String("hello world")},
		},
		// Contents         []*Content,
		// Metadata: map[string]string{},
		// Tags             []string,
		// Authors          []*Author,
	}
	b := &models.Page{
		Title:    proto.String("about page"),
		Theme:    proto.String("san-marzano"),
		Template: proto.String("about.html"),
	}

	assert.NoError(t, app.UpdatePage(a))
	assert.NoError(t, app.UpdatePage(b))
	assert.NotNil(t, b.Timestamps)
	assert.NotNil(t, b.Uuid)

	expectedPages := []*models.Page{pageWithoutContent(a), b}
	pages, err := app.ListPages(nil)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedPages, pages)
	for _, p := range []*models.Page{a, b} {
		page, err := app.GetPage(p.GetUuid())
		assert.NoError(t, err)
		assert.Equal(t, p, page)
	}

	pages, err = app.ListPages(&api.ListPageRequest{
		Options: &api.ListPageRequest_ListPageOptions{
			Filter: api.ListPageRequest_published.Enum(),
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, []*models.Page{pageWithoutContent(a)}, pages)

	pages, err = app.ListPages(&api.ListPageRequest{
		Options: &api.ListPageRequest_ListPageOptions{
			Filter: api.ListPageRequest_draft.Enum(),
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, []*models.Page{b}, pages)

	err = app.DeletePage(b)
	assert.NoError(t, err)
	page, err := app.GetPage(b.GetUuid())
	expectedErr := ErrNoKey(PAGE_BUCKET + ":" + b.GetUuid())
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, page)
}

func pageWithoutContent(page *models.Page) *models.Page {
	newPage := proto.Clone(page).(*models.Page)
	for _, c := range newPage.Contents {
		c.SetValue(nil)
	}
	return newPage
}
