package datasource

import "fmt"

type DataSourceError struct {
	err     error
	context string
}

func (err DataSourceError) Error() string {
	return fmt.Sprintf("%s (%s)", err.context, err.err.Error())
}

func (err DataSourceError) Unwrap() error {
	return err.err
}

func openConnectionError(err error) DataSourceError {
	return DataSourceError{err: err, context: "Unable to open connection to the database"}
}

func runQueryError(err error) DataSourceError {
	return DataSourceError{err: err, context: "Unable to run query"}
}
