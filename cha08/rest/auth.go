package rest

import (
	"net/http"
	"os"
)

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username, password, ok := r.BasicAuth(); ok {
			if username == os.Getenv("USERNAME") && password == os.Getenv("PASSWORD") {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", "Basic realm=\"restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
}
