package context

import (
	"context"

	"github.com/KrishanBhalla/iter/models"
)

type privateKey string

const (
	userKey privateKey = "user"
)

// WithUser adds a user to the context
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User returns a user from the context
func User(ctx context.Context) *models.User {
	if tmp := ctx.Value(userKey); tmp != nil {
		if user, ok := tmp.(*models.User); ok {
			return user
		}
	}
	return nil
}
