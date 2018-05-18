package api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xenolf/lego/acme"

	"github.com/ketchuphq/ketchup/server/tls"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/octavore/nagax/users"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Users["user123"] = &models.User{
		Uuid: proto.String("user123"),
	}

	// not logged in
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/user", nil)
	err := te.module.GetUser(rw, req, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNotFound, rw.Code)
	}

	// logged in
	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/user", nil)
	ctx := context.WithValue(context.Background(), users.UserTokenKey{}, "user123")
	err = te.module.GetUser(rw, req.WithContext(ctx), nil)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.JSONEq(t, `{"uuid":"user123"}`, rw.Body.String())
	}

	// logged in but wrong token
	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/user", nil)
	ctx = context.WithValue(context.Background(), users.UserTokenKey{}, "userABC")
	err = te.module.GetUser(rw, req.WithContext(ctx), nil)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNotFound, rw.Code)
	}
}

func TestGetTLS(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Users["user123"] = &models.User{
		Uuid: proto.String("user123"),
	}

	testTLS := &testTLSModule{}
	te.module.tls = testTLS

	testTLS.On("GetAllRegisteredDomains").
		Return([]string{"example.com"}, nil)
	testTLS.On("GetRegistration", "example.com", false).
		Return(&tls.Registration{
			Email:        "me@example.com",
			Domain:       "example.com",
			AgreedOn:     "2018-05-01",
			Registration: &acme.RegistrationResource{TosURL: "http://example.com/tos"},
		}, nil)
	testTLS.On("LoadCertResource", "example.com").
		Return(&acme.CertificateResource{}, nil)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/settings/tls", nil)
	err := te.module.GetTLS(rw, req, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		expected := &api.TLSSettingsResponse{
			TlsEmail:       proto.String("me@example.com"),
			TlsDomain:      proto.String("example.com"),
			AgreedOn:       proto.String("2018-05-01"),
			TermsOfService: proto.String("http://example.com/tos"),
			HasCertificate: proto.Bool(true),
		}
		output := &api.TLSSettingsResponse{}
		assert.NoError(t, jsonpb.Unmarshal(rw.Body, output))
		assert.Equal(t, expected, output)
		testTLS.AssertExpectations(t)
	}
}

func TestEnableTLS(t *testing.T) {
	te := setup()
	defer te.stop()

	testTLS := &testTLSModule{}
	te.module.tls = testTLS
	testTLS.On("ObtainCert", "me@example.com", "example.com").
		Return(nil)
	testTLS.On("GetRegistration", "example.com", false).
		Return(&tls.Registration{
			Email:        "me@example.com",
			Domain:       "example.com",
			AgreedOn:     "2018-05-01",
			Registration: &acme.RegistrationResource{TosURL: "http://example.com/tos"},
		}, nil)
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/settings/tls", bytes.NewBufferString(`{
		"tls_email": "me@example.com",
		"tls_domain": "example.com",
		"agreed": true
	}`))
	err := te.module.EnableTLS(rw, req, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		expected := &api.TLSSettingsResponse{
			TlsEmail:       proto.String("me@example.com"),
			TlsDomain:      proto.String("example.com"),
			AgreedOn:       proto.String("2018-05-01"),
			TermsOfService: proto.String("http://example.com/tos"),
			// HasCertificate: proto.Bool(true),
		}
		output := &api.TLSSettingsResponse{}
		assert.NoError(t, jsonpb.Unmarshal(rw.Body, output))
		assert.Equal(t, expected, output)
		testTLS.AssertExpectations(t)
	}
}
