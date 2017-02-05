package logger

import (
	"net/http"
	"time"
)

func New(log func(format string, args ...interface{})) func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		t := time.Now()
		next(rw, req)
		log("[%s] %q %v", req.Method, req.URL.Path, time.Now().Sub(t))
	}
}
