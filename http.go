package htpasswd

import (
	"fmt"
	"net/http"
)

// BasicAuthMiddleware implements a simple middleware handler for adding basic http auth to a route.
func BasicAuthMiddleware(realm string, htpasswd *Htpasswd) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				user, pass, ok := r.BasicAuth()
				if !ok || !htpasswd.Match(user, pass) {
					w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}
