package context

import (
	"fmt"

	"github.com/buttaciemanuel/sqlify/embed"
	"github.com/google/uuid"
)

type Document struct {
	Key   string
	Value string
}

func (document Document) String() string {
	return fmt.Sprintf("(`%s`, `%s`)", document.Key, document.Value)
}

func (document Document) UUID() string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(document.String())).String()
}

func (document Document) Embedding(model *embed.Embedder) ([]float32, error) {
	return model.Embed(document.String())
}

func (document Document) Metadata() map[string]any {
	return map[string]any{
		"key":   document.Key,
		"value": document.Value,
	}
}
