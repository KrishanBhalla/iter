package models

import (
	"github.com/KrishanBhalla/iter/hash"
	"github.com/dgraph-io/badger"
	"golang.org/x/crypto/bcrypt"
)

var _ UserDB = &userService{}

// UserService is a set of methods used to manipulate
// and work with the user model
type UserService interface {
	// Authenticate will verify the provided email address and
	// password are correct. If they are correct, the user corresponding
	// to that email will be returned. Otherwise you will recieve either:
	// ErrNotFound, ErrInvalidPassword, or another error if something
	// goes wrong
	Authenticate(email, password string) (*User, error)
	UserDB
}

type userService struct {
	UserDB
	pepper string
}

// NewUserService initialises a UserService object with an open connection
// to the db.
func NewUserService(db *badger.DB, hmacKey, pepper string) UserService {

	ug := &userDB{db}
	hmac := hash.NewHMAC(hmacKey)
	uv := NewUserValidator(ug, hmac, pepper)
	return &userService{
		UserDB: uv,
		pepper: pepper,
	}
}

// Authenticate can be used to authenticate a user with
// the provided email address and password.
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+us.pepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrIncorrectPassword
		default:
			return nil, err
		}

	}
	return foundUser, nil
}
