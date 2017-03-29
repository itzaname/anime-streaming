package controller

import (
	"net/http"
	"strings"
)

// ServeStatic maps static files
func ServeStatic(w http.ResponseWriter, r *http.Request) {
	// Disable listing directories
	if strings.HasSuffix(r.URL.Path, "/") {
		http.Error(w, "Item not found.", 404)
		return
	}
	http.ServeFile(w, r, r.URL.Path[1:])
}
