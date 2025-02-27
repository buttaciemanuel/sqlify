package context

import "fmt"

type ContextRetrievalError struct {
	err     error
	context string
}

func (err ContextRetrievalError) Error() string {
	return fmt.Sprintf("%s (%s)", err.context, err.err.Error())
}

func (err ContextRetrievalError) Unwrap() error {
	return err.err
}

func openDocumentStoreError(err error) ContextRetrievalError {
	return ContextRetrievalError{err: err, context: "Unable to open document store"}
}

func checkDocumentCollectionError(err error) ContextRetrievalError {
	return ContextRetrievalError{err: err, context: "Unable to check document collection"}
}

func createDocumentCollectionError(err error) ContextRetrievalError {
	return ContextRetrievalError{err: err, context: "Unable to create document collection"}
}

func addDocumentError(err error) ContextRetrievalError {
	return ContextRetrievalError{err: err, context: "Unable to add document to collection"}
}

func fetchDocumentsError(err error) ContextRetrievalError {
	return ContextRetrievalError{err: err, context: "Unable to fetch documents from the collection"}
}
