package tls

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xenolf/lego/acme"
)

func TestObtainCert(t *testing.T) {
	module, stop := setup(t)
	defer stop()

	err := module.ObtainCert("test@example.com", "example.com")
	assert.NoError(t, err)

	r, err := module.GetRegistration("example.com", false)
	assert.NoError(t, err)
	assert.Equal(t, &Registration{
		Email:    "test@example.com",
		Domain:   "example.com",
		AgreedOn: "2017-01-01T00:51:26Z",
	}, r)

	cr, err := module.LoadCertResource("example.com")
	assert.NoError(t, err)
	// CertURL is a hack only in tests
	assert.Equal(t, &acme.CertificateResource{Domain: "example.com", CertURL: "x"}, cr)
}

func TestRenew(t *testing.T) {
	domain := "example.com"
	module, stop := setup(t)
	defer stop()

	err := module.ObtainCert("test@example.com", domain)
	assert.NoError(t, err)
	registrationPath, _ := module.getCurrentRegistrationPath(domain)
	assert.Equal(t, "example.com-2017-01-01-v000.json", path.Base(registrationPath))
	cert, err := module.LoadCertResource(domain)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cert.CertURL))

	// test noop renew
	module.renewWithinInterval = time.Minute
	err = module.renewExpiredCerts()
	assert.NoError(t, err)
	registrationPath, _ = module.getCurrentRegistrationPath(domain)
	assert.Equal(t, "example.com-2017-01-01-v000.json", path.Base(registrationPath))
	cert, err = module.LoadCertResource(domain)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cert.CertURL))

	// test valid renew
	module.renewWithinInterval = time.Hour * 24 * 180
	err = module.renewExpiredCerts()
	assert.NoError(t, err)
	registrationPath, _ = module.getCurrentRegistrationPath(domain)
	assert.Equal(t, "example.com-2017-01-01-v000.json", path.Base(registrationPath))
	cert, err = module.LoadCertResource(domain)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cert.CertURL))
}

func TestSaveRegistration(t *testing.T) {
	module, stop := setup(t)
	defer stop()
	p := path.Join(module.Config.DataPath(tlsDir, ""), "example.com-2017-01-01-v000.json")
	_ = ioutil.WriteFile(p, []byte(`{
			"email": "oldadmin@example.com",
			"domain": "example.com"
		}`),
		os.ModePerm,
	)
	err := module.SaveRegistration(&Registration{
		Email:  "admin@example.com",
		Domain: "example.com",
	})
	assert.NoError(t, err)

	expectedPath := path.Join(module.Config.DataPath(tlsDir, ""), "example.com-2017-01-01-v000.json")
	_, err = os.Stat(expectedPath)
	assert.NoError(t, err)

	r, err := module.GetRegistration("example.com", false)
	assert.NoError(t, err)
	assert.Equal(t, &Registration{
		Email:  "admin@example.com",
		Domain: "example.com",
	}, r)
}

func TestGetRegistration(t *testing.T) {
	module, stop := setup(t)
	defer stop()
	p1 := path.Join(module.Config.DataPath(tlsDir, ""), "example.com-2017-01-01-v000.json")
	_ = ioutil.WriteFile(p1, []byte(`{"email": "oldadmin@example.com"}`), os.ModePerm)
	p2 := path.Join(module.Config.DataPath(tlsDir, ""), "example.com-2017-01-01-v001.json")
	_ = ioutil.WriteFile(p2, []byte(`{"email": "admin@example.com"}`), os.ModePerm)

	r, err := module.GetRegistration("fakedomain.com", false)
	if assert.NoError(t, err) {
		assert.Nil(t, r)
	}

	r, err = module.GetRegistration("example.com", false)
	if assert.NoError(t, err) {
		assert.Equal(t, &Registration{Email: "admin@example.com"}, r)
	}
}

func TestHTTPChallenge(t *testing.T) {
	module, stop := setup(t)
	defer stop()

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", challengeBasePath, nil)
	module.ServeHTTP(rw, req)
	assert.Equal(t, http.StatusNotFound, rw.Code)

	assert.NoError(t, module.Present("example.com", "1234", "abcd"))

	// incorrect path
	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", challengeBasePath, nil)
	module.ServeHTTP(rw, req)
	assert.Equal(t, http.StatusNotFound, rw.Code)

	// correct path
	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", acme.HTTP01ChallengePath("1234"), nil)
	module.ServeHTTP(rw, req)
	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "abcd", rw.Body.String())

	// make sure 404 after cleanup
	assert.NoError(t, module.CleanUp("example.com", "1234", "abcd"))
	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", acme.HTTP01ChallengePath("1234"), nil)
	module.ServeHTTP(rw, req)
	assert.Equal(t, http.StatusNotFound, rw.Code)
}
