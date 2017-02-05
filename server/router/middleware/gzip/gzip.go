package gzip

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
)

var Default = func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	handle := gziphandler.GzipHandler(next)
	handle.ServeHTTP(rw, req)
}

func New(lvl int) func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	handler := gziphandler.MustNewGzipLevelHandler(lvl)

	return func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		handle := handler(next)
		handle.ServeHTTP(rw, req)
	}
}
