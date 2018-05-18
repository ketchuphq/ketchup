package pkg

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
	"github.com/stretchr/testify/assert"
)

func TestFetchGithubPackage(t *testing.T) {
	called := false
	m := &Module{
		httpGet: func(url string) (*http.Response, error) {
			called = true
			return &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer([]byte(testGithubResponse)))}, nil
		},
	}
	resp, err := m.FetchTags(&packages.Package{
		VcsUrl: proto.String("https://github.com/ketchuphq/ketchup"),
	})
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.True(t, called)
	assert.Equal(t, RemoteTag{
		Name: "v0.2.0",
		SHA:  "dbf737990b33a980103f9a723d254502a8686886",
		URL:  "https://api.github.com/repos/ketchuphq/ketchup/commits/dbf737990b33a980103f9a723d254502a8686886",
	}, resp[0])
}

func TestFetchBitbucketPackage(t *testing.T) {
	m := &Module{
		httpGet: func(url string) (*http.Response, error) {
			return &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer([]byte(testBitbucketResponse)))}, nil
		},
	}
	resp, err := m.FetchTags(&packages.Package{
		VcsUrl: proto.String("https://bitbucket.org/tortoisehg/thg"),
	})
	assert.NoError(t, err)
	tags := []RemoteTag{{
		Name: "tip",
		SHA:  "83b7e2f14fcec35a739a99a0589e56a4cffa4f6d",
		URL:  "https://bitbucket.org/tortoisehg/thg/commits/tag/tip",
	}, {
		Name: "4.5.3",
		SHA:  "07157a0943be2f031a3cbe2d7a1e0d875ddf7731",
		URL:  "https://bitbucket.org/tortoisehg/thg/commits/tag/4.5.3",
	}, {
		Name: "4.5.2",
		SHA:  "4bbca812fbe6fe5abab4da1201cb25d9a0be59ae",
		URL:  "https://bitbucket.org/tortoisehg/thg/commits/tag/4.5.2",
	}}
	for i, tag := range tags {
		assert.Equalf(t, tag, resp[i], "comparing tag %s", tag.Name)
	}
}
