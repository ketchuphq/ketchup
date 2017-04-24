// pkg is responsible for listing and downloading
// themes and plugins from the internet.
// Themes and plugins are stored in the same registry, but
// get installed to different places? Maybe linked
// to themes dir if active?

package pkg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/octavore/naga/service"

	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
	"github.com/ketchuphq/ketchup/server/config"
	"github.com/ketchuphq/ketchup/util/errors"
)

// Module pkg manages downloading pkg
type Module struct {
	Config *config.Module

	RegistryURL string
}

const (
	githubPattern    = `https://github.com/{user}/{repo}`
	bitbucketPattern = `https://bitbucket.org/{user}/{repo}`
	githubAPI        = "https://api.github.com/repos/{user}/{repo}/tags"
	bitbucketAPI     = "https://api.bitbucket.org/2.0/repositories/{user}/{repo}/refs/tags"
)

func (m *Module) Init(c *service.Config) {
	// todo: keep track of all registries
}

// FetchPackage fetches information about the given package
// from upstream vcs repo
func (m *Module) FetchPackage(p *packages.Package) ([]*packages.Registry, error) {
	matches := regexp.MustCompile(`^https://([^/])/([^/])/([^/])$`).FindStringSubmatch(p.GetVcsUrl())
	if len(matches) != 3 {
		panic("die")
	}
	host, user, repo := matches[0], matches[1], matches[2]

	tagURL := ""
	if host == "github.com" {
		tagURL = githubAPI
	} else {
		tagURL = bitbucketAPI
	}
	tagURL = strings.Replace(tagURL, "{user}", user, 1)
	tagURL = strings.Replace(tagURL, "{repo}", repo, 2)

	res, err := http.Get(tagURL)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	registryList := []*packages.Registry{}
	err = json.Unmarshal(b, &registryList)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return nil, nil
}
