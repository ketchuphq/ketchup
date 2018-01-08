package tls

import (
	"crypto"

	"github.com/xenolf/lego/acme"
)

type acmeClient interface {
	SetChallengeProvider(acme.Challenge, acme.ChallengeProvider) error
	Register() (*acme.RegistrationResource, error)
	AgreeToTOS() error
	ObtainCertificate([]string, bool, crypto.PrivateKey, bool) (acme.CertificateResource, map[string]error)
	RenewCertificate(cert acme.CertificateResource, bundle, mustStaple bool) (acme.CertificateResource, error)
}

type legoAcme struct {
	*acme.Client
}

var newAcmeClient = newLegoAcmeClient

func newLegoAcmeClient(user acme.User, cp acme.ChallengeProvider) (acmeClient, error) {
	c, err := acme.NewClient(defaultCAURL, user, "")
	if err != nil {
		return nil, err
	}
	err = c.SetChallengeProvider(acme.HTTP01, cp)
	if err != nil {
		return nil, err
	}
	return legoAcme{c}, nil
}
