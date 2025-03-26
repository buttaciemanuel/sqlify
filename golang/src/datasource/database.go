package datasource

import (
	"database/sql"
	"fmt"

	_ "github.com/marcboeker/go-duckdb"
)

type DataSource interface {
	GetTables() ([]map[string]any, error)
	GetTableSchema(tableName string) ([]map[string]any, error)
	RunQuery(statement string) ([]map[string]any, error)
}

type Database struct {
	connection *sql.DB
}

func Duckdb(persistentFilePath string) (*Database, error) {
	connection, err := open("duckdb", persistentFilePath)

	if err != nil {
		return nil, err
	}

	if _, err := execute(connection, "INSTALL spatial; LOAD spatial;"); err != nil {
		return nil, err
	}

	return &Database{connection}, nil
}

func (database *Database) RunQuery(statement string) ([]map[string]any, error) {
	return execute(database.connection, statement)
}

func (database *Database) GetTables() ([]map[string]any, error) {
	return execute(database.connection, "show tables;")
}

func (database *Database) GetTableSchema(tableName string) ([]map[string]any, error) {
	return execute(database.connection, fmt.Sprintf("show table %s;", tableName))
}
