package tls

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/octavore/nagax/util/errors"
	"github.com/unrolled/secure"
)

var tlsPort = ":443"

func (m *Module) loadTLSConfig() (*tls.Config, error) {
	tlsConfig := &tls.Config{
		NextProtos: []string{"h2"}, // http/2
	}

	// keys are stored in a different place currently
	// extra dot is a hack to skip over session key
	glob := m.tlsDirPath("*.key")
	m.Logger.Infof("loading certs from %s", glob)
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if len(matches) == 0 {
		return nil, nil
	}

	for _, keyFile := range matches {
		crtFile := path.Base(strings.TrimSuffix(keyFile, ".key") + ".crt")
		cert, err := tls.LoadX509KeyPair(m.tlsDirPath(crtFile), keyFile)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, errors.Wrap(err)
		}
		m.Logger.Infof("found certs %s %s", keyFile, crtFile)
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}

// startTLSProxy starts the tls proxy on port 443. SSL is terminated
// and the connection in passed to the Router.
func (m *Module) startTLSProxy() error {
	if m.serverStarted {
		return nil
	}

	tlsConfig, err := m.loadTLSConfig()
	if err != nil {
		return err
	}
	if tlsConfig == nil {
		return nil
	}
	l, err := tls.Listen("tcp", tlsPort, tlsConfig)
	if err != nil {
		return err
	}
	m.server = &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      m.Router.Root,
		TLSConfig:    tlsConfig, // needed for http/2
	}

	// register middleware for https redirection
	https := secure.New(secure.Options{SSLRedirect: true})
	m.Router.Middleware.Prepend(https.HandlerFuncWithNext)

	m.Logger.Info("listening on ", tlsPort)
	m.serverStarted = true
	go m.server.Serve(l)
	return nil
}

// restartTLSProxy restarts the TLS proxy by first shutting
// down the old server and then starting a new one.
func (m *Module) restartTLSProxy() error {
	err := m.stopTLSProxy()
	if err != nil {
		return errors.Wrap(err)
	}
	err = m.startTLSProxy()
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (m *Module) stopTLSProxy() error {
	if !m.serverStarted {
		return nil
	}
	err := m.server.Shutdown(context.Background())
	if err != nil {
		return errors.Wrap(err)
	}
	m.serverStarted = false
	return nil
}
