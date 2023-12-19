package middleware

import (
	"net/http"

	"github.com/KrishanBhalla/iter/context"
)

// RequireUser is a middleware layer, dependent
// on the User MW layer, to verify a user is logged in.
type RequireUser struct {
	User
}

var _ Middleware = &RequireUser{}

// ApplyFn assumes that the User mw has already been run, otherwise
// it will not work correctly
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signup", http.StatusFound)
			return
		}
		next(w, r)
	})
}
