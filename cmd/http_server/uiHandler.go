package http_server

import (
	"net/http"
	"strings"
)

// UI Handler with SPA Fallback
func GetUIHandler(uifs http.FileSystem, indexHTML []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If requesting /, redirect to /gateway/
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/gateway/", http.StatusFound)
			return
		}

		// Try to serve static file
		prefix := "/gateway/"
		if strings.HasPrefix(r.URL.Path, prefix) {
			f, err := uifs.Open(r.URL.Path[len(prefix):])
			if err == nil {
				f.Close()
				http.StripPrefix(prefix, http.FileServer(uifs)).ServeHTTP(w, r)
				return
			}
		}

		// Fallback to index.html for SPA routes
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(indexHTML)
	})
}
