package models

import (
	"github.com/dgraph-io/badger"
	"gorm.io/gorm"
)

// first will return the first matching record from a db lookup
// 1) If the record is found, we will return a nil error
// 2) If the record is not found, we will return ErrNotFound
// 3) If any other error occurs, we will return and error with more information
// about what went wrong. This may not be generated by the errors package
func first(db *badger.DB, dst interface{}) error {

	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// find will return all matchings record from a db lookup
// 1) If records are found, we will return a nil error
// 2) If records are not found, we will return ErrNotFound
// 3) If any other error occurs, we will return and error with more information
// about what went wrong. This may not be generated by the errors package
func find(db *badger.DB, dst interface{}) error {

	err := db.Find(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
