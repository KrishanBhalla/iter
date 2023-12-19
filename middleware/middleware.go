package middleware

import "net/http"

// Middleware is a generic interface for a middleware
// object. It's purpose is to wrap a handler
// or a handler func with the middleware layer
// Apply should call on ApplyFn by passing in
// the next.ServeHTTP method which is itself a
// http.HandlerFunc
type Middleware interface {
	Apply(next http.Handler) http.HandlerFunc
	ApplyFn(next http.HandlerFunc) http.HandlerFunc
}
