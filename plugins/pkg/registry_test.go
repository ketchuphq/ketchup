package pkg

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
)

var marshaler = &jsonpb.Marshaler{
	EnumsAsInts: false,
	Indent:      "  ",
}

func TestSync(t *testing.T) {
	rp := &packages.Registry{
		RegistryVersion: proto.String("1.0"),
		RegistryType:    proto.String("themes"),
	}

	didMakeRequest := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		didMakeRequest = true
		marshaler.Marshal(w, rp)
	}))

	registry := &Registry{URL: ts.URL}
	err := registry.Sync()
	assert.NoError(t, err)
	assert.Equal(t, rp, registry.Registry)
	assert.True(t, didMakeRequest)
}

func TestSearch(t *testing.T) {
	aTheme := &packages.Package{Name: proto.String("aTheme")}
	bTheme := &packages.Package{Name: proto.String("bTheme")}
	rp := &packages.Registry{
		RegistryVersion: proto.String("1.0"),
		RegistryType:    proto.String("themes"),
		Packages:        []*packages.Package{aTheme, bTheme},
	}

	didMakeRequest := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		didMakeRequest = true
		marshaler.Marshal(w, rp)
	}))

	registry := &Registry{URL: ts.URL}
	p, err := registry.Search("aTheme")
	assert.NoError(t, err)
	assert.Equal(t, aTheme, p)
	assert.True(t, didMakeRequest)
}

func TestMatch(t *testing.T) {
	aTheme := &packages.Package{Name: proto.String("aTheme")}
	bTheme := &packages.Package{Name: proto.String("bTheme")}
	rp := &packages.Registry{
		RegistryVersion: proto.String("1.0"),
		RegistryType:    proto.String("themes"),
		Packages:        []*packages.Package{aTheme, bTheme},
	}

	didMakeRequest := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		didMakeRequest = true
		marshaler.Marshal(w, rp)
	}))

	registry := &Registry{URL: ts.URL}
	l, err := registry.Match(regexp.MustCompile("(?i)theme"))
	assert.NoError(t, err)
	assert.Equal(t, []*packages.Package{aTheme, bTheme}, l)
	assert.True(t, didMakeRequest)
}
