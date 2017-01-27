package tls

import (
	"crypto/tls"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

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
		return nil, err
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
			return nil, err
		}
		m.Logger.Infof("found certs %s %s", keyFile, crtFile)
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}

// StartTLSProxy starts the tls proxy
func (m *Module) StartTLSProxy() error {
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
	l, err := tls.Listen("tcp", ":443", tlsConfig)
	if err != nil {
		return err
	}
	server := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      m.Router,
		TLSConfig:    tlsConfig, // needed for http/2
	}
	m.serverStarted = true
	return server.Serve(l)
}
