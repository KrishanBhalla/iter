package models

import (
	"regexp"
	"strings"

	"github.com/KrishanBhalla/iter/hash"
	"github.com/KrishanBhalla/iter/rand"
	"golang.org/x/crypto/bcrypt"
)

var _ UserDB = &userValidator{}

// UserValidator allows for validation of the User object
type UserValidator interface {
	UserDB
}
type userValidator struct {
	UserDB
	hmac       hash.HMAC
	pepper     string
	emailRegex *regexp.Regexp
}

// NewUserValidator initialises a UserService object with an open connection
// to the db.
func NewUserValidator(db UserDB, hmac hash.HMAC, pepper string) UserValidator {

	return &userValidator{
		UserDB:     db,
		hmac:       hmac,
		pepper:     pepper,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-_]+\.[a-z]{2,16}$`),
	}
}

// ByEmail will normalise the email and call
// ByEmail on the subsequent UserDB layer
func (uv userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValFuncs(&user, uv.normaliseEmail)
	if err != nil {
		return nil, err
	}

	return uv.UserDB.ByEmail(user.Email)
}

// ByRemember will hash the remember token and then call
// ByRemember on the subsequent UserDB layer
func (uv userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	err := runUserValFuncs(&user, uv.hmacRememberToken)
	if err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(user.RememberHash)
}

// Create will hash the password and remember token if provided
// and call the UserDB Create method
func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		// Email
		uv.normaliseEmail,
		uv.requireEmail,
		uv.formatEmail,
		uv.emailIsAvailable,
		// Pwd
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		// RememberToken
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRememberToken,
		uv.rememberHashRequired,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

// Update will hash a remember token if provided and call the
// UserDB Update method
func (uv *userValidator) Update(user *User) error {

	err := runUserValFuncs(user,
		// Email
		uv.normaliseEmail,
		uv.requireEmail,
		uv.formatEmail,
		uv.emailIsAvailable,
		// Pwd
		uv.passwordHashRequired,
		// RememberToken
		uv.rememberMinBytes,
		uv.hmacRememberToken,
		uv.rememberHashRequired,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete will delete the user with given ID
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.userIDExists)
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// -------------------------------------------------------------------------------------
// User --------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------

func (uv *userValidator) userIDExists(user *User) error {
	_, err := uv.UserDB.ByID(user.ID)
	if err != nil {
		return ErrUserNotFound
	}
	return nil
}

// requireEmail will ensure the email address provided
// is not empty
func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

// formatEmail will ensure the email address provided
// is of the required form
func (uv *userValidator) formatEmail(user *User) error {
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrInvalidEmail
	}
	return nil
}

// normaliseEmail will remove whitespace from an email
// address and convert it to lowercase.
func (uv *userValidator) normaliseEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

// emailIsAvailable will remove whitespace from an email
// address and convert it to lowercase.
func (uv *userValidator) emailIsAvailable(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err != nil {
		if err == ErrNotFound {
			return nil
		}
		return err

	}
	if user.ID != existing.ID {
		return ErrEmailAlreadyTaken
	}
	return nil
}

// -------------------------------------------------------------------------------------
// Password ----------------------------------------------------------------------------
// -------------------------------------------------------------------------------------

func (uv *userValidator) passwordMinLength(user *User) error {
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordHashRequired
	}
	return nil
}

// bcryptPassword will  hash a user's password with a predefined
// pepper and bcrypt if the Password field is not empty
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + uv.pepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

// -------------------------------------------------------------------------------------
// Remember ----------------------------------------------------------------------------
// -------------------------------------------------------------------------------------

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrRememberTooShort
	}
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberHashRequired
	}
	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	return nil
}

func (uv *userValidator) hmacRememberToken(user *User) error {
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

// -------------------------------------------------------------------------------------
// Types -------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------

// userValFunc defines the standard type we need for validation
type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}
