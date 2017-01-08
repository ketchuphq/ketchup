package tls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xenolf/lego/acme"

	"github.com/octavore/press/util/errors"
)

func (m *Module) saveCert(cert acme.CertificateResource) error {
	b, err := json.MarshalIndent(cert, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(m.tlsDirPath(cert.Domain+".json"), b, 0600)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(m.tlsDirPath(cert.Domain+".crt"), cert.Certificate, 0600)
	if err != nil {
		return err
	}

	// we generate our own private key, so no need to save
	// err = ioutil.WriteFile(m.tlsDirPath(cert.Domain+".key"), cert.PrivateKey, 0600)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (m *Module) loadCert() {
	// todo: load acme.CertificateResource from disk
}

func (m *Module) saveUser(u *SSLUser) error {
	b, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(m.tlsDirPath(userFile), b, 0600)
}

func (m *Module) loadUser() (*SSLUser, error) {
	// load user
	u := &SSLUser{}
	b, err := ioutil.ReadFile(m.tlsDirPath(userFile))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if !os.IsNotExist(err) {
		err = json.Unmarshal(b, u)
		if err != nil {
			return nil, err
		}
	}
	if u.Email == "" {
		// todo: detect if changed
		u.Email = m.config.TLS.Email
	}

	// load key
	if m.config.TLS.URL == "" {
		return nil, errors.Wrap(fmt.Errorf("no ssl url"))
	}
	_, u.key, err = m.keystore.LoadPrivateKey(u.Email + ".key")
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return u, nil
}
