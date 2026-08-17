// Package server implements the HTTP handler that serves the built frontend.
package server

import (
	"io/fs"
	"net/http"
	"path"
)

// New returns a handler that serves files from dist. Requests for paths that
// don't match a file fall back to index.html so client-side routing works.
func New(dist fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(dist))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := path.Clean(r.URL.Path[1:])
		if name == "" {
			name = "."
		}

		if _, err := fs.Stat(dist, name); err != nil {
			r = r.Clone(r.Context())
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}
