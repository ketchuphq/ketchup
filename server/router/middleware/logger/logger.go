package logger

import (
	"net/http"
	"time"
)

func New(log func(format string, args ...interface{})) func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		t := time.Now()
		next(rw, req)
		ua := req.Header.Get("User-Agent")
		log("[%s] %q %v %s", req.Method, req.URL.Path, time.Now().Sub(t), ua)
	}
}
