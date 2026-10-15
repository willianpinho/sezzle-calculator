package httpapi

import (
	_ "embed"
	"net/http"
)

// openAPIContract is the domain-owned public API contract served verbatim.
//
//go:embed openapi.json
var openAPIContract []byte

func (a API) openAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(openAPIContract)
}

func (a API) docsRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	http.Redirect(w, r, "/docs/", http.StatusTemporaryRedirect)
}
