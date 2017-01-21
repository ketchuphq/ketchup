package tls

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/keystore"
	"github.com/octavore/nagax/logger"
	"github.com/xenolf/lego/acme"

	"github.com/octavore/press/server/config"
	"github.com/octavore/press/server/router"
	"github.com/octavore/press/util/errors"
)

const (
	defaultCAURL = "https://acme-v01.api.letsencrypt.org/directory"
	// defaultCAURL      = "https://acme-staging.api.letsencrypt.org/directory"
	challengeBasePath = "/.well-known/acme-challenge"
	tlsDir            = "tls"
)

type Module struct {
	Config *config.Module
	Router *router.Module
	Logger *logger.Module

	keystore *keystore.KeyStore
}

func (m *Module) Init(c *service.Config) {
	c.AddCommand(&service.Command{
		Keyword:    "tls:provision <example.com> <my@email.com>",
		ShortUsage: `Provision an ssl cert for the given domain and email`,
		Usage: `Provision an ssl cert for the given domain.
Required params: domain to provision a cert for; contact email for Let's Encrypt.`,
		Flags: []*service.Flag{{Key: "agree"}},
		Run: func(ctx *service.CommandContext) {
			ctx.RequireExactlyNArgs(2)
			if !ctx.Flags["agree"].Present() {
				fmt.Print("Please provide the --agree flag to indicate that you agree to Let's Encrypt's TOS. \n")
				return
			}
			err := m.ObtainCert(ctx.Args[1], ctx.Args[0])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("success!")
		},
	})
	c.Setup = func() error {
		m.keystore = &keystore.KeyStore{Dir: m.Config.Config.DataDir}

		dir := path.Join(m.Config.Config.DataDir, tlsDir)
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0700)
		}
		return err
	}
	c.Start = func() {
		go func() {
			err := m.StartTLSProxy()
			if err != nil {
				m.Logger.Error(errors.Wrap(err))
			}
		}()
	}
}

func (m *Module) tlsDirPath(file string) string {
	return path.Join(m.Config.Config.DataDir, tlsDir, file)
}

// func (m *Module) Renew(r *Registration) error {}

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
	c, err := acme.NewClient(defaultCAURL, r, "")
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

	m.saveCert(cert)
	return errors.Wrap(err)
}
