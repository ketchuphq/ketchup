package router

import (
	"errors"
	"log"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/router"
	"github.com/octavore/press/proto/press/api"
)

var ErrNotFound = errors.New("not found")

type Module struct {
	*router.Module
}

func (m *Module) Init(c *service.Config) {
}

func (m *Module) Subrouter(path string) *httprouter.Router {
	r := httprouter.New()
	m.Handle(path, r)
	return r
}

func Proto(rw http.ResponseWriter, pb proto.Message) error {
	rw.Header().Set("Content-Type", "application/json")
	return JSON.Marshal(rw, pb)
}

func InternalError(rw http.ResponseWriter, err error) {
	if err == ErrNotFound {
		NotFound(rw)
		return
	}
	log.Println(err)
	rw.WriteHeader(http.StatusInternalServerError)
	err = Proto(rw, &api.Error{
		Code:   api.ErrorCode_INTERNAL_SERVER_ERROR.Enum(),
		Detail: proto.String("Internal server error."),
	})
	if err != nil {
		log.Println(err)
	}
}

func NotFound(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotFound)
	err := Proto(rw, &api.Error{
		Code:   api.ErrorCode_NOT_FOUND.Enum(),
		Detail: proto.String("Not found."),
	})
	if err != nil {
		log.Println(err)
	}
}

var JSON = &jsonpb.Marshaler{
	EnumsAsInts: false,
	Indent:      "  ",
}
