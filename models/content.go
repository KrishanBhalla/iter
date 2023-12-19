package models

import (
	"encoding/json"

	"github.com/KrishanBhalla/iter/internal/helpers"
	"github.com/KrishanBhalla/iter/internal/services"
	"github.com/dgraph-io/badger"
)

// Content implements the content table for
// storing travel content data - this consists of
// the source URL, the content, and the embedding
type Content struct {
	URL       string    `json:"url"`
	Content   string    `json:"content"`
	Embedding []float64 `json:"embedding"`
}

type ContentDB interface {
	ByEmbedding(embedding []float64) ([]Content, error)

	// Methods for altering contents
	Create(content *Content) error
	Update(content *Content) error
	Delete(url string) error
	DbCloser
}

// Define contentDB and ensure it implements ContentDB
var _ ContentDB = &contentDB{}

type contentDB struct {
	db                  *badger.DB
	similarityThreshold float64
	embeddingService    services.EmbeddingService
}

// ByEmbedding finds the closest pieces of content to an embedded query
func (cdb *contentDB) ByEmbedding(embedding []float64) ([]Content, error) {
	var content []Content
	err := getAll(cdb.db, content)
	if err != nil {
		return nil, err
	}

	result := make([]Content, 0)
	for _, c := range content {
		similarity, err := helpers.EmbeddingCosineSimilarity(c.Embedding, embedding)
		if err != nil {
			return nil, err
		}
		if similarity > cdb.similarityThreshold {
			result = append(result, content...)
		}
	}
	return result, nil
}

// Create will create the provided content and backfill data
func (cdb *contentDB) Create(content *Content) error {
	return cdb.Update(content)
}

func (cdb *contentDB) Update(content *Content) error {

	if content.Embedding == nil {
		embedding, err := cdb.embeddingService.GetEmbedding(content.Content)
		if err != nil {
			return err
		}
		content.Embedding = embedding
	}

	err := cdb.db.Update(func(txn *badger.Txn) error {

		contentBytes, err := json.Marshal(&content)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(content.URL), contentBytes)
		return err
	})
	return err
}

func (cdb *contentDB) Delete(url string) error {
	err := cdb.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(url))
		return err
	})
	return err
}

func (cdb *contentDB) CloseDB() error {
	return cdb.db.Close()
}
