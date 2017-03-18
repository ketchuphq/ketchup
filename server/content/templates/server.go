package templates

import (
	"mime"
	"net/http"
	"path"
)

// ServeHTTP serves theme assets
func (m *Module) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	asset, err := m.GetAsset(req.URL.Path)
	if err != nil {
		m.Logger.Warning(err)
		m.NotFound(rw, req)
		return
	}
	if asset == nil {
		m.NotFound(rw, req)
		return
	}
	ext := path.Ext(req.URL.Path)
	rw.Header().Add("Content-Type", mime.TypeByExtension(ext))
	rw.Write([]byte(asset.GetData()))
}

// NotFound serves a 404 page from a 404.html asset, if it exists.
func (m *Module) NotFound(rw http.ResponseWriter, req *http.Request) {
	asset, err := m.GetAsset("/404.html")
	if err != nil {
		m.Logger.Warning(err)
		http.NotFound(rw, req)
		return
	}
	if asset == nil {
		http.NotFound(rw, req)
		return
	}
	ext := path.Ext(req.URL.Path)
	rw.Header().Add("Content-Type", mime.TypeByExtension(ext))
	rw.Write([]byte(asset.GetData()))
}
