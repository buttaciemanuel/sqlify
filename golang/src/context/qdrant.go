package context

import (
	"context"

	"github.com/buttaciemanuel/sqlify/embed"
	"github.com/qdrant/go-client/qdrant"
)

type QdrantStore struct {
	client   *qdrant.Client
	embedder embed.Embedder
}

func Qdrant(embedder embed.Embedder) (*QdrantStore, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		return nil, openDocumentStoreError(err)
	}

	return &QdrantStore{client, embedder}, nil
}

func (dataContext *QdrantStore) Initialize(subContexts []string) error {
	for _, subContext := range subContexts {
		exists, err := dataContext.client.CollectionExists(context.Background(), subContext)

		if err != nil {
			return checkDocumentCollectionError(err)
		}

		if exists {
			continue
		}

		if err := dataContext.client.CreateCollection(context.Background(), &qdrant.CreateCollection{
			CollectionName: subContext,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     uint64(dataContext.embedder.EmbeddingLength()),
				Distance: qdrant.Distance_Cosine,
			}),
		}); err != nil {
			return createDocumentCollectionError(err)
		}
	}

	return nil
}

func (dataContext *QdrantStore) ContainsDocument(subContext string, document Document) (bool, error) {
	results, err := dataContext.client.Get(context.Background(), &qdrant.GetPoints{
		CollectionName: subContext,
		Ids:            []*qdrant.PointId{qdrant.NewIDUUID(document.UUID())},
	})

	if err != nil {
		return false, fetchDocumentsError(err)
	}

	return len(results) > 0, nil
}

func (dataContext *QdrantStore) Clear() error {
	collectionNames, err := dataContext.client.ListCollections(context.Background())

	if err != nil {
		return checkDocumentCollectionError(err)
	}

	for _, collectionName := range collectionNames {
		dataContext.client.DeleteCollection(context.Background(), collectionName)
	}

	return nil
}

func (dataContext *QdrantStore) StoreDocument(subContext string, document Document) error {
	embedding, err := document.Embedding(&dataContext.embedder)

	if err != nil {
		return err
	}

	if _, err := dataContext.client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: subContext,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDUUID(document.UUID()),
				Vectors: qdrant.NewVectors(embedding...),
				Payload: qdrant.NewValueMap(document.Metadata()),
			},
		},
	}); err != nil {
		return addDocumentError(err)
	}

	return nil
}

func (dataContext *QdrantStore) FetchSimilarDocuments(subContext, queryKey string, limit uint) ([]Document, error) {
	embedding, err := dataContext.embedder.Embed(queryKey)

	if err != nil {
		return nil, err
	}

	results, err := dataContext.client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: subContext,
		Query:          qdrant.NewQuery(embedding...),
		WithPayload:    qdrant.NewWithPayload(true),
		Limit:          qdrant.PtrOf(uint64(limit)),
	})

	if err != nil {
		return nil, fetchDocumentsError(err)
	}

	documents := []Document{}

	for _, result := range results {
		documents = append(documents, Document{
			Key:   result.Payload["key"].GetStringValue(),
			Value: result.Payload["value"].GetStringValue(),
		})
	}

	return documents, nil
}
