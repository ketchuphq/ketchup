package gzip

import (
	"compress/gzip"
	"net/http"
	"path"

	"github.com/NYTimes/gziphandler"
)

var Default = New(gzip.DefaultCompression)

func New(lvl int) func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	handler := gziphandler.MustNewGzipLevelHandler(lvl)

	return func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		// todo: whitelist gzippable data
		if path.Ext(req.URL.Path) == "" {
			next(rw, req)
			return
		}
		handle := handler(next)
		handle.ServeHTTP(rw, req)
	}
}
