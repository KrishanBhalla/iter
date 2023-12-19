package helpers

import "fmt"

// Embedding vectors have norm 1, which means the dot product is the cosine similarity
func EmbeddingCosineSimilarity(v1, v2 []float64) (float64, error) {
	if len(v1) != len(v2) {
		return 0., fmt.Errorf("internal: vectors do not have the same length")
	}
	var ab float64
	for i, v := range v1 {
		w := v2[i]
		ab += v * w
	}
	return ab, nil
}
