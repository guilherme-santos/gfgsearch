package http

import (
	"crypto/subtle"
	"fmt"
	"net/http"
)

// BasicAuthMiddleware is a http middleware to accept only authenticated requests.
// This middleware could be implemented used JWT for example, but to keep the
// code simple and without external dependencies I decided to use BasicAuth.
func BasicAuthMiddleware(username, password string) Middleware {
	checkAuth := func(usr, passwd string) bool {
		return subtle.ConstantTimeCompare([]byte(usr), []byte(username)) == 1 &&
			subtle.ConstantTimeCompare([]byte(passwd), []byte(password)) == 1
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			usr, passwd, ok := req.BasicAuth()
			if !ok || !checkAuth(usr, passwd) {
				w.Header().Set("WWW-Authenticate", `Basic realm="GFG-Search"`)
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, http.StatusText(http.StatusUnauthorized))
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}
