package api

import (
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/octavore/nagax/router"
	"github.com/octavore/nagax/users"

	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/server/tls"
	"github.com/ketchuphq/ketchup/util/errors"
)

func (m *Module) GetUser(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	userUUID, ok := req.Context().Value(users.UserTokenKey{}).(string)
	if !ok {
		m.Router.EmptyJSON(rw, http.StatusNotFound)
		return nil
	}
	user, err := m.DB.GetUser(userUUID)
	if err != nil {
		return err
	}
	user.HashedPassword = nil
	user.Token = nil
	return router.ProtoOK(rw, user)
}

func (m *Module) GetTLS(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	domains, err := m.TLS.GetAllRegisteredDomains()
	if err != nil {
		return err
	}
	if len(domains) == 0 {
		return m.Router.EmptyJSON(rw, http.StatusOK)
	}
	domain := domains[0]
	r, err := m.TLS.GetRegistration(domain, false)
	if err != nil {
		return err
	}
	tosURL := ""
	if r.Registration != nil {
		tosURL = r.Registration.TosURL
	}
	crt, err := m.TLS.LoadCertResource(domain)
	if err != nil {
		m.Logger.Error(err)
	}
	res := &api.TLSSettingsResponse{
		TlsEmail:       &r.Email,
		TlsDomain:      &r.Domain,
		AgreedOn:       &r.AgreedOn,
		TermsOfService: &tosURL,
		HasCertificate: proto.Bool(crt != nil),
	}
	return router.ProtoOK(rw, res)
}

func (m *Module) EnableTLS(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	rpb := &api.EnableTLSRequest{}
	err := jsonpb.Unmarshal(req.Body, rpb)
	if err != nil {
		return errors.Wrap(err)
	}

	err = m.TLS.ObtainCert(rpb.GetTlsEmail(), rpb.GetTlsDomain())
	if err != nil {
		if errors.IsType(err, tls.LetsEncryptError{}) {
			return m.Router.SimpleError(rw, http.StatusBadRequest, err)
		}
		return errors.Wrap(err)
	}

	r, err := m.TLS.GetRegistration(rpb.GetTlsDomain(), false)
	if err != nil {
		return err
	}
	tosURL := ""
	if r.Registration != nil {
		tosURL = r.Registration.TosURL
	}
	res := &api.TLSSettingsResponse{
		TlsEmail:       &r.Email,
		TlsDomain:      &r.Domain,
		AgreedOn:       &r.AgreedOn,
		TermsOfService: &tosURL,
	}
	return router.ProtoOK(rw, res)
}

func (m *Module) Logout(rw http.ResponseWriter, req *http.Request, _ router.Params) error {
	m.Auth.Auth.Logout(rw, req)
	return nil
}
