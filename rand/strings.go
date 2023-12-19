package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const rememberTokenBytes = 32

// Bytes takes an integer and generates a random byte slice of length n
// this uses the crypto/rand package so is safe to use with
// remember tokens
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	// now contains random bytes
	return b, nil
}

// String will generate a byte slice of size nBytes and then
// return a string that is the base64 URL encoded version
// of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// NBytes returns the number of bytes in a base64 string
func NBytes(base64String string) (int, error) {
	bytes, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return -1, err
	}
	return len(bytes), nil
}

// RememberToken is a helper function designed
// to generate remember tokens of a predetermined byte size
func RememberToken() (string, error) {
	return String(rememberTokenBytes)
}
