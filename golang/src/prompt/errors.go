package prompt

import "fmt"

type PromptConstructionError struct {
	err     error
	context string
}

func (err PromptConstructionError) Error() string {
	return fmt.Sprintf("%s (%s)", err.context, err.err.Error())
}

func (err PromptConstructionError) Unwrap() error {
	return err.err
}

func fetchDefinitionsError(err error) PromptConstructionError {
	return PromptConstructionError{err: err, context: "Unable to fetch definitions from the vector database"}
}

func fetchQueriesError(err error) PromptConstructionError {
	return PromptConstructionError{err: err, context: "Unable to fetch query examples from the vector database"}
}

func readResponseError(err error) PromptConstructionError {
	return PromptConstructionError{err: err, context: "Unable to read response from Ollama"}
}

func unpackResponseError(err error) PromptConstructionError {
	return PromptConstructionError{err: err, context: "Unable to unpack embedding response as JSON from Ollama"}
}
