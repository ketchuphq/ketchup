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
	"github.com/octavore/nagax/util/errors"

	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
	"github.com/ketchuphq/ketchup/server/config"
)

// Module pkg manages downloading pkg
type Module struct {
	Config *config.Module

	httpGet     func(url string) (*http.Response, error)
	RegistryURL string
}

type endpoint struct {
	pattern       string
	baseAPIURL    string
	parseResponse func([]byte) ([]RemoteTag, error)
}

type RemoteTag struct {
	Name string
	SHA  string
	URL  string
}

type githubTag struct {
	Name   string `json:"name"`
	Commit struct {
		SHA string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
}

type bitbucketResponse struct {
	Values []struct {
		Name  string `json:"name"`
		Links struct {
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
		} `json:"links"`
		Target struct {
			Hash string `json:"hash"`
		} `json:"target"`
	} `json:"values"`
}

var github = endpoint{
	pattern:    `https://github.com/(?P<user>[^/]+)/(?P<repo>[^/]+)`,
	baseAPIURL: "https://api.github.com/repos/{user}/{repo}/tags",
	parseResponse: func(data []byte) ([]RemoteTag, error) {
		response := []githubTag{}
		err := json.Unmarshal(data, &response)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		out := []RemoteTag{}
		for _, el := range response {
			out = append(out, RemoteTag{
				Name: el.Name,
				SHA:  el.Commit.SHA,
				URL:  el.Commit.URL,
			})
		}
		return out, nil
	},
}
var bitbucket = endpoint{
	pattern:    `https://bitbucket.org/(?P<user>[^/]+)/(?P<repo>[^/]+)`,
	baseAPIURL: "https://api.bitbucket.org/2.0/repositories/{user}/{repo}/refs/tags?sort=-name",
	parseResponse: func(data []byte) ([]RemoteTag, error) {
		response := &bitbucketResponse{}
		err := json.Unmarshal(data, response)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		out := []RemoteTag{}
		for _, e := range response.Values {
			out = append(out, RemoteTag{
				Name: e.Name,
				SHA:  e.Target.Hash,
				URL:  e.Links.HTML.Href,
			})
		}
		return out, nil
	},
}

func (m *Module) Init(c *service.Config) {
	// todo: keep track of all registries
	c.Setup = func() error {
		m.httpGet = http.Get
		return nil
	}
}

func (e endpoint) fetch(get func(url string) (*http.Response, error), p *packages.Package) ([]RemoteTag, error) {
	re := regexp.MustCompile(e.pattern)
	parts := re.FindStringSubmatch(p.GetVcsUrl())
	if len(parts) == 0 {
		return nil, nil
	}
	url := strings.Replace(e.baseAPIURL, "{user}", parts[1], 1)
	url = strings.Replace(url, "{repo}", parts[2], 1)
	resp, err := get(url)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return e.parseResponse(data)
}

// FetchTags fetches tag information about the given package
// from upstream vcs repo
func (m *Module) FetchTags(p *packages.Package) ([]RemoteTag, error) {
	for _, e := range []endpoint{github, bitbucket} {
		tags, err := e.fetch(m.httpGet, p)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		if tags != nil {
			return tags, nil
		}
	}
	return nil, nil
}
