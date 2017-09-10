package router

import (
	"net/http"
	"path"

	"github.com/go-errors/errors"
	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"
	"github.com/octavore/nagax/router"

	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/server/router/middleware/gzip"
	logger2 "github.com/ketchuphq/ketchup/server/router/middleware/logger"
)

type Module struct {
	*router.Module
	Logger *logger.Module
}

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		m.Module.ErrorHandler = m.errorHandler

		m.Module.Middleware.Set(
			logger2.New(m.Logger.Infof),
			gzip.Default,
		)
		return nil
	}
}

func (m *Module) errorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	m.InternalError(rw, err)
}

func (m *Module) InternalError(rw http.ResponseWriter, err error) {
	if err == router.ErrNotFound {
		m.NotFound(rw)
		return
	}

	switch e := err.(type) {
	case *errors.Error:
		s := e.StackFrames()[0]
		f := path.Base(s.File)
		m.Logger.Errorf("[%s/%s:%d] %v", s.Package, f, s.LineNumber, e.Error())
	default:
		m.Logger.Errorf("router: internal error %v", e)
	}
	err = router.Proto(rw, http.StatusInternalServerError, &api.Error{
		Code:   api.ErrorCode_INTERNAL_SERVER_ERROR.Enum(),
		Detail: proto.String("Internal server error."),
	})
	if err != nil {
		m.Logger.Errorf("router: %v", err)
	}
}

func (m *Module) NotFound(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotFound)
	err := router.Proto(rw, http.StatusNotFound, &api.Error{
		Code:   api.ErrorCode_NOT_FOUND.Enum(),
		Detail: proto.String("Not found."),
	})
	if err != nil {
		m.Logger.Errorf("router: not found error %v", err)
	}
}
