package tls

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/keystore"
	"github.com/xenolf/lego/acme"

	"github.com/octavore/nagax/logger"

	"github.com/octavore/press/server/config"
	"github.com/octavore/press/server/router"
	"github.com/octavore/press/util/errors"
)

const (
	defaultCAURL = "https://acme-v01.api.letsencrypt.org/directory"
	// defaultCAURL      = "https://acme-staging.api.letsencrypt.org/directory"
	challengeBasePath = "/.well-known/acme-challenge"
	tlsDir            = "tls"
	userFile          = "user.json"
)

type tlsConfig struct {
	TLS struct {
		Email       string
		URL         string
		TLSAddress  int
		HTTPAddress int
	}
}

type Module struct {
	Config *config.Module
	Router *router.Module
	Logger *logger.Module

	keystore *keystore.KeyStore
	config   tlsConfig
}

func (m *Module) Init(c *service.Config) {
	c.AddCommand(&service.Command{
		Keyword: "tls:provision",
		Usage:   "Provision an ssl cert. todo: document required settings in config.json file",
		Flags:   []*service.Flag{{Key: "agree"}},
		Run: func(ctx *service.CommandContext) {
			if !ctx.Flags["agree"].Present() {
				fmt.Print("Please provide the --agree flag to indicate that you agree to Let's Encrypt's TOS. \n")
				return
			}
			user, err := m.GetTLSUser(true)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = m.obtainCert(user)
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
		if err != nil {
			return err
		}

		return m.Config.ReadConfig(&m.config)
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

// func (m *Module) Renew(u *SSLUser) error {}

func (m *Module) handleObtainCert(rw http.ResponseWriter, req *http.Request) {
	user, err := m.GetTLSUser(true)
	if err != nil {
		m.Router.InternalError(rw, err)
		return
	}
	err = m.obtainCert(user)
	if err != nil {
		m.Router.InternalError(rw, err)
	}
}

// ObtainCert obtains a new ssl cert for the given user
func (m *Module) obtainCert(u *SSLUser) error {
	// Initialize user and domain
	if m.config.TLS.URL == "" {
		return errors.Wrap(fmt.Errorf("no url specified"))
	}
	// hack to URL parse it correctly
	if !strings.HasPrefix(m.config.TLS.URL, "https://") && !strings.HasPrefix(m.config.TLS.URL, "http://") {
		m.config.TLS.URL = "http://" + m.config.TLS.URL
	}
	domain, err := url.Parse(m.config.TLS.URL)
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
	c, err := acme.NewClient(defaultCAURL, u, "")
	if err != nil {
		return errors.Wrap(err)
	}

	if m.config.TLS.HTTPAddress != 0 {
		err = c.SetHTTPAddress(fmt.Sprintf(":%d", m.config.TLS.HTTPAddress))
		if err != nil {
			return errors.Wrap(err)
		}
	}

	if m.config.TLS.TLSAddress != 0 {
		err = c.SetTLSAddress(fmt.Sprintf(":%d", m.config.TLS.TLSAddress))
		if err != nil {
			return errors.Wrap(err)
		}
	}

	registration, err := c.Register()
	if err != nil {
		return errors.Wrap(err)
	}

	u.Registration = registration
	m.saveUser(u)

	err = c.AgreeToTOS()
	if err != nil {
		return errors.Wrap(err)
	}

	cert, errs := c.ObtainCertificate([]string{domain.Host}, true, domainKey, false)
	if len(errs) > 0 {
		lst := []string{}
		for domain, e := range errs {
			// todo: check for updated TOS error
			lst = append(lst, fmt.Sprintf("%s: %s", domain, e.Error()))
		}
		return errors.Wrap(fmt.Errorf(strings.Join(lst, "; ")))
	}

	m.saveCert(cert)
	return errors.Wrap(err)
}
