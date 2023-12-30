package models

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/KrishanBhalla/iter/internal/helpers"
	"github.com/KrishanBhalla/iter/internal/services"
	"github.com/dgraph-io/badger"
)

// Content implements the content table for
// storing travel content data - this consists of
// the source URL, the content, and the embedding
type Content struct {
	URL       string    `json:"url"`
	Country   string    `json:"country"`
	Location  string    `json:"location"`
	Content   string    `json:"content"`
	Embedding []float64 `json:"embedding"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

type ContentDB interface {
	ByCountryAndSimilarity(country string, embedding []float64) ([]Content, error)
	BySimilarity(embedding []float64) ([]Content, error)
	Countries() ([]string, error)

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

// ByCountryAndSimilarity finds the closest pieces of content to an embedded query
func (cdb *contentDB) ByCountryAndSimilarity(country string, embedding []float64) ([]Content, error) {
	var allContent map[string]Content
	data, err := get(cdb.db, country)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("%s:  %s", "Failed to retrieve content", err.Error())
	}
	fmt.Println(string(data))
	dec := json.NewDecoder(strings.NewReader(string(data)))
	err = dec.Decode(&allContent)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("%s:  %s", "Failed to retrieve unmarshal", err.Error())
	}
	fmt.Println(allContent)

	var content = make([]Content, 0)
	for _, c := range allContent {
		content = append(content, c)
	}

	return cdb.bySimilarity(content, embedding)
}

// Countries finds the closest pieces of content to an embedded query
func (cdb *contentDB) Countries() ([]string, error) {
	keys := make([]string, 0)
	err := keyStrings(cdb.db, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// BySimilarity finds the closest pieces of content to an embedded query
func (cdb *contentDB) BySimilarity(embedding []float64) ([]Content, error) {
	var allContent []map[string]Content
	err := getAll(cdb.db, allContent)
	if err != nil {
		return nil, err
	}

	var content = make([]Content, 0)
	for _, c := range allContent {
		for _, v := range c {
			content = append(content, v)
		}
	}
	return cdb.bySimilarity(content, embedding)
}

func (cdb *contentDB) bySimilarity(content []Content, embedding []float64) ([]Content, error) {
	result := make([]Content, 0)
	for _, c := range content {
		similarity, err := helpers.EmbeddingCosineSimilarity(c.Embedding, embedding)
		if err != nil {
			return nil, err
		}
		fmt.Println(similarity, c.Content)
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

	existingContent := make(map[string]Content)
	data, err := get(cdb.db, content.Country)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	dec := json.NewDecoder(strings.NewReader(string(data)))
	err = dec.Decode(&existingContent)
	if err != nil && err != io.EOF {
		return fmt.Errorf("%s:  %s", "Failed to retrieve unmarshal", err.Error())
	}

	err = cdb.db.Update(func(txn *badger.Txn) error {

		textHash := md5.New()
		io.WriteString(textHash, content.Content)
		existingContent[content.URL+" "+string(textHash.Sum(nil))] = *content
		contentBytes, err := json.Marshal(existingContent)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(content.Country), contentBytes)
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
