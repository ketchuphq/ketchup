package tls

import (
	"crypto"
	"crypto/rsa"
	"path"

	"github.com/octavore/nagax/keystore"
	"github.com/xenolf/lego/acme"

	"github.com/ketchuphq/ketchup/util/errors"
)

type Registration struct {
	Email        string                     `json:"email"`
	AgreedOn     string                     `json:"agreed_on"`
	Domain       string                     `json:"domain"`
	Registration *acme.RegistrationResource `json:"registration"`

	key *rsa.PrivateKey
}

func (r *Registration) Init(ks *keystore.KeyStore) (err error) {
	keyFile := path.Join(tlsDir, r.GetEmail()+".key")
	_, r.key, err = ks.LoadPrivateKey(keyFile)
	return errors.Wrap(err)
}

func (r *Registration) GetEmail() string {
	return r.Email
}

func (r *Registration) GetPrivateKey() crypto.PrivateKey {
	return r.key
}

func (r *Registration) GetRegistration() *acme.RegistrationResource {
	return r.Registration
}
