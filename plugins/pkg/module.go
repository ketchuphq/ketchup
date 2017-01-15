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

	"github.com/octavore/press/proto/press/packages"
	"github.com/octavore/press/server/config"
	"github.com/octavore/press/util/errors"
)

// Module pkg manages downloading pkg
type Module struct {
	Config *config.Module
}

const (
	// defaultRegistryURL       = "https://ketchuphq.com/registry.json"
	defaultRegistryURL = "http://localhost:8000/registry.json"
	githubPattern      = `https://github.com/{user}/{repo}`
	bitbucketPattern   = `https://bitbucket.org/{user}/{repo}`
	githubAPI          = "https://api.github.com/repos/{user}/{repo}/tags"
	bitbucketAPI       = "https://api.bitbucket.org/2.0/repositories/{user}/{repo}/refs/tags"
	themeDir           = "themes"
)

func (m *Module) Init(c *service.Config) {
	c.AddCommand(&service.Command{
		Keyword:    "themes:get <git url>",
		ShortUsage: "get the theme",
		Usage:      "get the theme",
		Run: func(ctx *service.CommandContext) {
			ctx.RequireExactlyNArgs(1)
			// err := m.FetchRepo("test", ctx.Args[0])
			// if err != nil {
			// 	panic(err)
			// }
		},
	})
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
