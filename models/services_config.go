package models

import (
	"github.com/dgraph-io/badger"
)

// ServicesConfig allows for dynamic adding of services
type ServicesConfig func(*Services) error

// withBadger initiates a badger db
func withBadger(dbPath string) (*badger.DB, error) {
	opt := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// WithUser returns a ServicesConfig object that sets a user
func WithUser(hmacKey, pepper string) ServicesConfig {
	return func(s *Services) error {
		userDB, err := withBadger("users")
		if err != nil {
			return err
		}
		s.User = NewUserService(userDB, hmacKey, pepper)
		return nil
	}
}

// WithContent returns a ServicesConfig object that sets content
func WithContent(similarityThreshold float64) ServicesConfig {
	return func(s *Services) error {
		contentDb, err := withBadger("content")
		if err != nil {
			return err
		}
		s.Content = NewContentService(contentDb, similarityThreshold)
		return nil
	}
}
