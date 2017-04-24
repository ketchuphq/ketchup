package pkg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"

	"github.com/golang/protobuf/proto"

	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
	"github.com/ketchuphq/ketchup/util/errors"
)

type Registry struct {
	URL      string
	Registry *packages.Registry

	mu sync.RWMutex
}

func (r *Registry) Proto() *packages.Registry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return proto.Clone(r.Registry).(*packages.Registry)
}

func (r *Registry) Sync() error {
	res, err := http.Get(r.URL)
	if err != nil {
		return errors.Wrap(err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err)
	}
	repo := &packages.Registry{}
	err = json.Unmarshal(b, &repo)
	if err != nil {
		return errors.Wrap(err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Registry = repo
	return nil
}

func (r *Registry) Search(name string) (*packages.Package, error) {
	err := r.Sync()
	if err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.Registry.Packages {
		if p.GetName() == name {
			return p, nil
		}
	}
	return nil, nil
}

func (r *Registry) Match(re *regexp.Regexp) ([]*packages.Package, error) {
	err := r.Sync()
	if err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []*packages.Package{}
	for _, p := range r.Registry.Packages {
		if re.MatchString(p.GetName()) {
			out = append(out, p)
		}
	}
	return out, nil
}

// FetchDefaultRegistry fetches the default registry
func (m *Module) Registry(registryURL string) *Registry {
	return &Registry{URL: registryURL}
}

// press registry daemon should periodically scrape
func getGithubTags(p *packages.Package) {
	// paginate should cache
}

func getBitbucketTags(p *packages.Package) {
	// paginate should cache
}
