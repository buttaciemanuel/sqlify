package ollama

import "fmt"

type PromptSubmissionError struct {
	err     error
	context string
}

func (err PromptSubmissionError) Error() string {
	return fmt.Sprintf("%s (%s)", err.context, err.err.Error())
}

func (err PromptSubmissionError) Unwrap() error {
	return err.err
}

func sendPromptRequestError(err error) PromptSubmissionError {
	return PromptSubmissionError{
		err:     err,
		context: "Unable to send post request to Ollama to submit",
	}
}

func pullModelRequestError(err error) PromptSubmissionError {
	return PromptSubmissionError{
		err:     err,
		context: "Unable to pull model in Ollama",
	}
}

func readResponseError(err error) PromptSubmissionError {
	return PromptSubmissionError{
		err:     err,
		context: "Unable to read response from Ollama",
	}
}

func unpackResponseError(err error) PromptSubmissionError {
	return PromptSubmissionError{
		err:     err,
		context: "Unable to unpack embedding response as JSON from Ollama",
	}
}
