package middleware

import (
	"net/http"
	"strings"

	"github.com/KrishanBhalla/iter/context"
	"github.com/KrishanBhalla/iter/models"
)

// User wraps models.UserService with middleware
// It adds a user to the request context going forwards
type User struct {
	models.UserService
}

var _ Middleware = &User{}

// Apply middleware to a http.Handler
func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn to a http.HandlerFunc
func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// if the user is requesting a static asset or image, we don't need to look up the current user.
		if strings.HasPrefix(path, "/assets/") {
			next(w, r)
			return
		}
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}
