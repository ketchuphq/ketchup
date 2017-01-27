package tls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xenolf/lego/acme"

	"regexp"
)

var now = time.Now

func (m *Module) saveCert(cert acme.CertificateResource) error {
	b, err := json.MarshalIndent(cert, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(m.tlsDirPath(cert.Domain+".json"), b, 0600)
	if err != nil {
		return err
	}

	// we generate our own private key, so no need to save cert.PrivateKey
	return ioutil.WriteFile(m.tlsDirPath(cert.Domain+".crt"), cert.Certificate, 0600)
}

// LoadCertResource will return the CertificateResource if it exists on the disk, else nil.
func (m *Module) LoadCertResource(domain string) (*acme.CertificateResource, error) {
	certPath := m.tlsDirPath(domain + ".json")
	_, err := os.Stat(certPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	cert := &acme.CertificateResource{}
	err = json.Unmarshal(b, cert)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// domain-yyyy-mm-dd-###.json where # is an incrementing number
func (m *Module) getNextRegistrationPath(domain string) (string, error) {
	date := now().Format("2006-01-02")
	prefix := fmt.Sprintf("%s-%s-v", domain, date)
	cur, err := m.getCurrentRegistrationPath(domain)
	if err != nil {
		return "", err
	}
	cur = path.Base(cur)
	if !strings.HasPrefix(cur, prefix) {
		return path.Join(m.Config.Config.DataDir, tlsDir, prefix+"000.json"), nil
	}

	// trim [prefix]...[.json] to extract number
	i, err := strconv.Atoi(strings.TrimPrefix(cur, prefix)[0:3])
	if err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%s%03d.json", prefix, i+1)
	return path.Join(m.Config.Config.DataDir, tlsDir, filename), nil
}

func (m *Module) GetAllRegisteredDomains() ([]string, error) {
	// format is domain-yyyy-mm-dd-v###.json
	g := path.Join(m.Config.Config.DataDir, tlsDir, "*-*-*-*-v*.json")
	matches, err := filepath.Glob(g)
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)
	o := []string{}
	seen := map[string]bool{}
	re := regexp.MustCompile(`(.+?)-[0-9]{4}-[0-9]{2}-[0-9]{2}-v[0-9]{3}.json`)
	for _, match := range matches {
		m := path.Base(match)
		s := re.FindStringSubmatch(m)
		if len(s) == 0 {
			continue
		}
		m = s[1]
		if seen[m] {
			continue
		}
		seen[m] = true
		o = append(o, m)
	}
	return o, nil
}

func (m *Module) getCurrentRegistrationPath(domain string) (string, error) {
	g := path.Join(m.Config.Config.DataDir, tlsDir, domain+"-*"+".json")
	matches, err := filepath.Glob(g)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", nil
	}

	sort.Strings(matches)
	return matches[len(matches)-1], nil
}

func (m *Module) SaveRegistration(r *Registration) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	p, err := m.getNextRegistrationPath(r.Domain)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(p, b, 0600)
}

// GetRegistration from the tls dir.
func (m *Module) GetRegistration(domain string, withPrivateKey bool) (*Registration, error) {
	// load user
	r := &Registration{}

	p, err := m.getCurrentRegistrationPath(domain)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, r)
	if err != nil {
		return nil, err
	}
	if !withPrivateKey {
		return r, nil
	}

	return r, nil
}
