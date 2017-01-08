package tls

import (
	"crypto"
	"crypto/rsa"

	"github.com/xenolf/lego/acme"
)

type SSLUser struct {
	Email        string                     `json:"email"`
	Registration *acme.RegistrationResource `json:"registration"`

	key *rsa.PrivateKey
}

func (s *SSLUser) GetEmail() string {
	return s.Email
}

func (s *SSLUser) GetRegistration() *acme.RegistrationResource {
	return s.Registration
}

func (s *SSLUser) GetPrivateKey() crypto.PrivateKey {
	return s.key
}
