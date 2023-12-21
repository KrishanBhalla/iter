package models

import (
	"github.com/KrishanBhalla/iter/internal/services"
	"github.com/dgraph-io/badger"
)

var _ ContentDB = &contentService{}

// ContentService is a set of methods used to manipulate
// and work with the content model
type ContentService interface {
	ContentDB
}

type contentService struct {
	ContentDB
}

// NewContentService initialises a ContentService object with an open connection
// to the db.
func NewContentService(db *badger.DB, similarityThreshold float64) ContentService {

	cdb := &contentDB{db, similarityThreshold, &services.EmbeddingModel{ModelName: services.ADA002}}
	return &contentService{
		ContentDB: cdb,
	}
}
