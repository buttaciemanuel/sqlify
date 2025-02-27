package context

type DocumentStore interface {
	Clear() error
	Initialize(subContexts []string) error
	StoreDocument(subContext string, document Document) error
	FetchSimilarDocuments(subContext, queryKey string, limit uint) ([]Document, error)
}
