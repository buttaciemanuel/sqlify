package embed

import "fmt"

type EmbeddingError struct {
	err     error
	context string
}

func (err EmbeddingError) Error() string {
	return fmt.Sprintf("%s (%s)", err.context, err.err.Error())
}

func (err EmbeddingError) Unwrap() error {
	return err.err
}

func createEmbeddingRequestError(err error) EmbeddingError {
	return EmbeddingError{err: err, context: "Unable to create embedding request for Ollama"}
}

func sendEmbeddingRequestError(err error) EmbeddingError {
	return EmbeddingError{err: err, context: "Unable to send embedding request to Ollama"}
}

func readEmbeddingResponseError(err error) EmbeddingError {
	return EmbeddingError{err: err, context: "Unable to read embedding response from Ollama"}
}

func unpackEmbeddingResponseError(err error) EmbeddingError {
	return EmbeddingError{err: err, context: "Unable to unpack embedding response as JSON from Ollama"}
}
