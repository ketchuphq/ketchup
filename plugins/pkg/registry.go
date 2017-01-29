package pkg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/octavore/ketchup/proto/ketchup/packages"
	"github.com/octavore/ketchup/util/errors"
)

// FetchRegistry fetches
func (m *Module) FetchDefaultRegistry() (*packages.Registry, error) {
	return m.FetchRegistry(defaultRegistryURL)
}
func (m *Module) FetchRegistry(registryURL string) (*packages.Registry, error) {
	res, err := http.Get(registryURL)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	repo := &packages.Registry{}
	err = json.Unmarshal(b, &repo)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return repo, nil
}

// press registry daemon should periodically scrape
func getGithubTags(p *packages.Package) {
	// paginate should cache
}

func getBitbucketTags(p *packages.Package) {
	// paginate should cache
}
