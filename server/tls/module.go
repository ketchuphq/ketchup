package tls

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/keystore"
	"github.com/octavore/nagax/logger"
	"github.com/octavore/nagax/util/errors"
	"github.com/xenolf/lego/acme"

	"github.com/ketchuphq/ketchup/server/config"
	"github.com/ketchuphq/ketchup/server/router"
)

const (
	defaultCAURL               = "https://acme-v01.api.letsencrypt.org/directory"
	defaultStagingCAURL        = "https://acme-staging.api.letsencrypt.org/directory"
	challengeBasePath          = "/.well-known/acme-challenge/"
	tlsDir                     = "tls"
	defaultRenewWithinInterval = time.Hour * 24 * 14 // renew if cert expires within two weeks
)

type acmeChallenge struct {
	domain, token, keyAuth string
}

type Module struct {
	Config *config.Module
	Router *router.Module
	Logger *logger.Module

	challenge *acmeChallenge
	keystore  *keystore.KeyStore

	server              *http.Server
	serverStarted       bool
	renewWithinInterval time.Duration
}

func (m *Module) Init(c *service.Config) {
	m.registerCommands(c)
	c.Setup = func() error {
		m.renewWithinInterval = defaultRenewWithinInterval
		m.keystore = &keystore.KeyStore{Dir: m.Config.Config.DataDir}

		dir := path.Join(m.Config.Config.DataDir, tlsDir)
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0700)
		}
		m.Router.Root.Handle(challengeBasePath, m)
		return err
	}
	c.Start = func() {
		err := m.startTLSProxy()
		if err != nil {
			m.Logger.Error(errors.Wrap(err))
		}

		go func() {
			for range time.Tick(2 * time.Hour) {
				err := m.renewExpiredCerts()
				if err != nil {
					m.Logger.Error(err)
				}
			}
		}()
	}

	c.Stop = func() {
		err := m.stopTLSProxy()
		if err != nil {
			m.Logger.Error(err)
		}
	}
}

func (m *Module) tlsDirPath(file string) string {
	return path.Join(m.Config.Config.DataDir, tlsDir, file)
}

func (m *Module) renewExpiredCerts() error {
	tlsConfig, err := m.loadTLSConfig()
	if err != nil {
		return err
	}
	for _, cert := range tlsConfig.Certificates {
		x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			return errors.Wrap(err)
		}

		expiration := x509Cert.NotAfter
		nowPlusDelta := now().Add(m.renewWithinInterval)
		if nowPlusDelta.After(expiration) {
			domain := x509Cert.Subject.CommonName
			m.Logger.Infof("expired cert: renewing cert for %s", domain)
			r, err := m.GetRegistration(domain, true)
			if err != nil {
				return errors.Wrap(err)
			}

			err = r.Init(m.keystore)
			if err != nil {
				return errors.Wrap(err)
			}

			acmeCert, err := m.LoadCertResource(domain)
			if err != nil {
				return errors.Wrap(err)
			}

			// set Certificate
			certBytes := &bytes.Buffer{}
			err = pem.Encode(certBytes, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Certificate[0]})
			if err != nil {
				return errors.Wrap(err)
			}
			acmeCert.Certificate = certBytes.Bytes()

			// set private key
			keyFile := path.Join(tlsDir, domain+".key")
			keyBytes, _, err := m.keystore.LoadPrivateKey(keyFile)
			if err != nil {
				return errors.Wrap(err)
			}
			pkBytes := &bytes.Buffer{}
			err = pem.Encode(pkBytes, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes})
			if err != nil {
				return errors.Wrap(err)
			}
			acmeCert.PrivateKey = pkBytes.Bytes()

			c, err := newAcmeClient(r, m)
			if err != nil {
				return errors.Wrap(err)
			}
			newCert, err := c.RenewCertificate(*acmeCert, false, false)
			if err != nil {
				return errors.Wrap(err)
			}
			err = m.saveCert(newCert)
			if err != nil {
				return errors.Wrap(err)
			}
			m.Logger.Infof("successfully renewed cert for %s; restarting server", domain)
			err = m.restartTLSProxy()
			if err != nil {
				return errors.Wrap(err)
			}
			m.Logger.Info("server restarted")
		}
	}
	return nil
}

func (m *Module) ObtainCert(email, domain string) error {
	r, err := m.GetRegistration(domain, true)
	if err != nil {
		return err
	}
	if r == nil {
		r = &Registration{}
	}
	r.Domain = domain
	r.Email = email
	r.AgreedOn = now().Format(time.RFC3339)
	r.Init(m.keystore)
	return m.obtainCert(r)
}

type LetsEncryptError struct{ error }

// ObtainCert obtains a new ssl cert for the given user. Currently uses default
// port 80 and port 443 for challenges.
func (m *Module) obtainCert(r *Registration) error {
	certURL := r.Domain
	// Initialize user and domain
	if certURL == "" {
		return errors.Wrap(fmt.Errorf("no url specified"))
	}
	// hack to URL parse it correctly
	if !strings.HasPrefix(certURL, "https://") && !strings.HasPrefix(certURL, "http://") {
		certURL = "http://" + certURL
	}
	domain, err := url.Parse(certURL)
	if err != nil {
		return err
	}

	if domain.Host == "" {
		return errors.Wrap(fmt.Errorf("no url specified"))
	}

	keyFile := path.Join(tlsDir, domain.Host+".key")
	_, domainKey, err := m.keystore.LoadPrivateKey(keyFile)
	if err != nil {
		return errors.Wrap(err)
	}

	// Initialize the Client
	r.Domain = domain.Host
	c, err := newAcmeClient(r, m)
	if err != nil {
		return errors.Wrap(err)
	}

	registration, err := c.Register()
	if err != nil {
		return errors.Wrap(err)
	}

	r.Registration = registration
	m.SaveRegistration(r)

	err = c.AgreeToTOS()
	if err != nil {
		return errors.Wrap(err)
	}
	cert, errs := c.ObtainCertificate([]string{domain.Host}, true, domainKey, false)
	if len(errs) > 0 {
		lst := []string{}
		for _, e := range errs {
			// todo: check for updated TOS error
			lst = append(lst, e.Error())
		}

		return errors.Wrap(LetsEncryptError{fmt.Errorf(strings.Join(lst, "; "))})
	}

	err = m.saveCert(cert)
	return errors.Wrap(err)
}

func (m *Module) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if m.challenge == nil {
		m.Router.NotFound(rw)
		return
	}
	// The handler validates the HOST header and request type.
	// For validation it then writes the token the server returned with the challenge
	if strings.HasPrefix(req.Host, m.challenge.domain) &&
		req.URL.Path == acme.HTTP01ChallengePath(m.challenge.token) &&
		req.Method == "GET" {
		rw.Header().Add("Content-Type", "text/plain")
		rw.Write([]byte(m.challenge.keyAuth))
		m.challenge = nil
	} else {
		m.Logger.Warningf("Invalid acme challenge for %s", req.Host)
		m.Router.NotFound(rw)
	}
}

// Present implements the acme.ChallengeProvider.Present
func (m *Module) Present(domain, token, keyAuth string) error {
	if m.challenge != nil {
		m.Logger.Warningf("replacing existing challenge for %s with %s", m.challenge.domain, domain)
	}
	m.challenge = &acmeChallenge{domain: domain, token: token, keyAuth: keyAuth}
	return nil
}

// CleanUp implements the acme.ChallengeProvider.CleanUp
func (m *Module) CleanUp(domain, token, keyAuth string) error {
	m.challenge = nil
	return nil
}
