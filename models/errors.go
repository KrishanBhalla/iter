package models

import (
	"strings"

	"golang.org/x/text/cases"
)

const (
	// ErrNotFound is returned when a resource cannot be found
	// in the database
	ErrNotFound modelError = "models: resource not found"

	// ErrUserNotFound is returned when a resource cannot be found
	// in the database
	ErrUserNotFound modelError = "models: resource not found"
	// ErrIncorrectPassword is returned when a password does not match the value
	// in the database
	ErrIncorrectPassword modelError = "models: incorrect password provided"
	// ErrPasswordTooShort is returned when a password does not meet
	// length requirements
	ErrPasswordTooShort modelError = `models: password provided was too short.
	 It must be at least 8 characters`
	// ErrPasswordRequired is returned when a password is not provided
	// on Create
	ErrPasswordRequired modelError = `models: password is required`
	// ErrPasswordHashRequired is returned when a password hash is not provided
	ErrPasswordHashRequired modelError = `models: password hash is required`
	// ErrInvalidEmail is returned when a password does not match our
	// email regexp
	ErrInvalidEmail modelError = "models: email address is not valid"
	// ErrEmailRequired is returned when an email is the empty string
	ErrEmailRequired modelError = "models: email address is required"
	// ErrEmailAlreadyTaken is returned on Update or Create if
	// an email address is already in use
	ErrEmailAlreadyTaken modelError = "models: email address is already taken"

	// private errors

	// ErrInvalidID is returned when an invalid ID is provided
	// to a method like Delete
	ErrInvalidID privateError = "models: ID provided was invalid"
	// ErrRememberTooShort is returned when a password does not meet
	// length requirements
	ErrRememberTooShort privateError = `models: remember provided was too short.\n
	It must be at least 32 bytes`
	// ErrRememberHashRequired is returned when a create or update is attempted and
	// a remember token hash is not provided
	ErrRememberHashRequired privateError = `models: remember hash is required`
	// ErrUserIDRequired when user id is not provided on create (gallery)
	ErrUserIDRequired privateError = `models: user id is required`
	// ErrTitleRequired  when title is not provided on create (gallery)
	ErrTitleRequired modelError = `models: title is required`
)

// modelError is a string so that some errors could be made constant
type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	return cases.Title(s).String()
}

// privateError is a string so that some errors could be made constant
type privateError string

func (e privateError) Error() string {
	return string(e)
}
