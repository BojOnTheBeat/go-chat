package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

// Remember: implementing the http.Handler interface means implementing ServeHTTP
// So this implementation means our authHandler struct is now a http.Handler
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")

	if err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		// some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// success - call the next handler
	h.next.ServeHTTP(w, r)
}

// MustAuth creates an authHandler that wraps any other http.Handler
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles the third party login process.
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")

	if len(segs) < 4 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not sure what to do with this url")
	}
	action := segs[2]
	provider := segs[3]

	switch action {
	case "login":
		log.Println("TODO handle login for", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
