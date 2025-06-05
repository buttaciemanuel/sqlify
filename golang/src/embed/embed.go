package embed

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Embedder struct {
	name      string
	embedSize uint
}

func (embedder Embedder) Name() string {
	return embedder.name
}

func (embedder Embedder) EmbeddingLength() uint {
	return embedder.embedSize
}

func (embedder Embedder) Embed(text string) ([]float32, error) {
	if embeddings, err := embedder.EmbedMany([]string{text}); err != nil {
		return nil, err
	} else {
		return embeddings[0], nil
	}
}

func (embedder Embedder) EmbedMany(texts []string) ([][]float32, error) {
	parameters, err := json.Marshal(map[string]any{
		"model": embedder.name,
		"input": texts,
	})
	if err != nil {
		return nil, createEmbeddingRequestError(err)
	}

	response, err := http.Post(
		"http://ollama:11434/api/embed",
		"application/json",
		bytes.NewBuffer(parameters),
	)
	if err != nil {
		return nil, sendEmbeddingRequestError(err)
	}

	defer response.Body.Close()

	var jsonBody struct {
		Model           string      `json:"model"`
		Embeddings      [][]float32 `json:"embeddings"`
		TotalDuration   int         `json:"total_duration"`
		LoadDuration    int         `json:"load_duration"`
		PromptEvalCount int         `json:"prompt_eval_count"`
	}

	rawBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, readEmbeddingResponseError(err)
	}

	if err := json.Unmarshal(rawBody, &jsonBody); err != nil {
		return nil, unpackEmbeddingResponseError(err)
	}

	return jsonBody.Embeddings, nil
}
