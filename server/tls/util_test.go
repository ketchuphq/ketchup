package tls

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/octavore/naga/service"
	"github.com/xenolf/lego/acme"
)

const testDomain = "example.com"
const testTime = "2017-01-01T00:51:26Z"
const testPath = testDomain + "-2017-01-01-v000.json"

func init() {
	newAcmeClient = newTestAcmeClient
	now = testNow
	tlsPort = "localhost:9443"
}

// testAcmeClient replaces the actual acme client for tests
type testAcmeClient struct {
}

func (c *testAcmeClient) SetChallengeProvider(acme.Challenge, acme.ChallengeProvider) error {
	return nil
}

func (c *testAcmeClient) Register() (*acme.RegistrationResource, error) {
	return nil, nil
}

func (c *testAcmeClient) AgreeToTOS() error {
	return nil
}

func newCert(domain string, privKey crypto.PrivateKey, certURL string) (acme.CertificateResource, error) {
	var pubKey interface{}
	if k, ok := privKey.(*rsa.PrivateKey); ok {
		pubKey = &k.PublicKey
	} else {
		panic("unsupported key")
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1234),
		Subject: pkix.Name{
			Country:    []string{"US"},
			CommonName: domain,
		},
		NotBefore: now(),
		NotAfter:  now().AddDate(0, 0, 1),
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, pubKey, privKey)
	if err != nil {
		panic(err)
	}
	buf := &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if err != nil {
		panic(err)
	}
	// certURL is a hack to store some metadata that the test can check
	return acme.CertificateResource{
		Domain:      domain,
		CertURL:     certURL,
		Certificate: buf.Bytes(),
	}, nil
}

func (c *testAcmeClient) ObtainCertificate(domains []string, bundle bool, privKey crypto.PrivateKey, mustStaple bool) (acme.CertificateResource, map[string]error) {
	cert, _ := newCert(domains[0], privKey, "x")
	return cert, nil
}

func (c *testAcmeClient) RenewCertificate(cert acme.CertificateResource, bundle bool, mustStaple bool) (acme.CertificateResource, error) {
	pemBlock, _ := pem.Decode(cert.PrivateKey)
	privKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return acme.CertificateResource{}, err
	}
	return newCert(cert.Domain, privKey, cert.CertURL+"x")
}

func newTestAcmeClient(user acme.User, cp acme.ChallengeProvider) (acmeClient, error) {
	return &testAcmeClient{}, nil
}

func testNow() time.Time {
	t, _ := time.Parse(time.RFC3339, testTime)
	return t
}

func setup(t *testing.T) (*Module, func()) {
	m := &Module{}
	svc := service.New(m)
	stop := svc.StartForTest()
	return m, stop
}
