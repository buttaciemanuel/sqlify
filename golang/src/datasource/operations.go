package datasource

import (
	"database/sql"
)

func open(driverName, persistentFilePath string) (*sql.DB, error) {
	connection, err := sql.Open(driverName, persistentFilePath)

	if err != nil {
		return nil, openConnectionError(err)
	}

	return connection, nil
}

func execute(connection *sql.DB, statement string) ([]map[string]any, error) {
	rows, err := connection.Query(statement)

	if err != nil {
		return nil, runQueryError(err)
	}

	results := []map[string]any{}
	attributes, _ := rows.Columns()

	for rows.Next() {
		scans := make([]any, len(attributes))
		row := make(map[string]any)

		for i := range scans {
			scans[i] = &scans[i]
		}

		rows.Scan(scans...)

		for i, v := range scans {
			if v != nil {
				row[attributes[i]] = v.(string)
			}
		}

		results = append(results, row)
	}

	return results, nil
}
