package pipeline

import (
	"fmt"
	"strings"

	"github.com/buttaciemanuel/sqlify/context"
	"github.com/buttaciemanuel/sqlify/datasource"
	"github.com/buttaciemanuel/sqlify/ollama"
	"github.com/buttaciemanuel/sqlify/prompt"
)

type Config struct {
	Context    context.DocumentStore
	Database   datasource.DataSource
	Prompt     prompt.Prompt
	Model      string
	AutoSchema bool
}

type Pipeline struct {
	context  context.DocumentStore
	database datasource.DataSource
	prompt   prompt.Prompt
	model    string
}

func New(config *Config) (*Pipeline, error) {
	if err := config.Context.Initialize([]string{"schema", "example"}); err != nil {
		return nil, err
	}

	if config.AutoSchema {
		tables, err := config.Database.GetTables()

		if err != nil {
			return nil, err
		}

		for _, table := range tables {
			name := table["name"].(string)
			schema, err := config.Database.GetTableSchema(name)

			if err != nil {
				return nil, err
			}

			if err := config.Context.StoreDocument("schema", context.Document{Key: name, Value: createTableDefinition(name, schema)}); err != nil {
				return nil, err
			}
		}
	}

	return &Pipeline{context: config.Context, database: config.Database, prompt: config.Prompt, model: config.Model}, nil
}

func (pipeline *Pipeline) DefineSchema(name string, schema string) error {
	if err := pipeline.context.StoreDocument("schema", context.Document{Key: name, Value: schema}); err != nil {
		return err
	}

	return nil
}

func (pipeline *Pipeline) DefineExample(name string, query string) error {
	if err := pipeline.context.StoreDocument("example", context.Document{Key: name, Value: query}); err != nil {
		return err
	}

	return nil
}

func (pipeline *Pipeline) Execute(query string) ([]map[string]any, error) {
	schemas, err := pipeline.context.FetchSimilarDocuments("schema", query, 1)

	if err != nil {
		return nil, err
	}

	examples, err := pipeline.context.FetchSimilarDocuments("example", query, 1)

	if err != nil {
		return nil, err
	}

	prompt := pipeline.prompt.Build(query, schemas, examples)

	fmt.Println(prompt)

	statement, err := ollama.Generate(pipeline.model, prompt)

	if err != nil {
		return nil, err
	}

	fmt.Println(statement)

	result, err := pipeline.database.RunQuery(statement)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func createTableDefinition(name string, schema []map[string]any) string {
	var builder strings.Builder

	builder.WriteString("CREATE TABLE ")
	builder.WriteString(name)
	builder.WriteString("(")

	if len(schema) > 0 {
		builder.WriteString(schema[0]["column_name"].(string))
		builder.WriteString(" ")
		builder.WriteString(schema[0]["column_type"].(string))
	}

	for _, column := range schema[1:] {
		builder.WriteString(", ")
		builder.WriteString(column["column_name"].(string))
		builder.WriteString(" ")
		builder.WriteString(column["column_type"].(string))
	}

	builder.WriteString(");")

	return builder.String()
}
