package cli

import (
	"errors"
	"maps"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Context struct {
		Store    string
		Embedder string
	}
	Database struct {
		Autoschema bool
		Duckdb     struct {
			Filename string
		}
	}
	Model struct {
		Name string
	}
	Prompt struct {
		Schemas struct {
			Title string
			Items []string
			From  string
		}
		Examples struct {
			Title string
			Items []string
			From  string
		}
		Query struct {
			Title string
		}
		Sections []struct {
			Title string
			Body  string
			Items []string
		}
	}
}

func validate(config *Config) error {
	if len(config.Context.Store) == 0 {
		return errors.New("Field context.store cannot be empty")
	}

	if len(config.Context.Store) == 0 {
		return errors.New("Field context.embedder cannot be empty")
	}

	if len(config.Database.Duckdb.Filename) == 0 {
		return errors.New("Field database.duckdb.filename cannot be empty")
	}

	if len(config.Model.Name) == 0 {
		return errors.New("Field model.name cannot be empty")
	}

	if !config.Database.Autoschema && len(config.Prompt.Schemas.Items) == 0 {
		return errors.New(
			"Field config.database.autoschema must be set if you do not provide schemas in config.prompt.schemas.items",
		)
	}

	if len(config.Prompt.Examples.Title) == 0 {
		config.Prompt.Examples.Title = "examples"
	}

	if len(config.Prompt.Examples.Items) > 0 &&
		len(config.Prompt.Examples.From) > 0 {
		return errors.New(
			"Field config.prompt.examples.items and config.prompt.examples.from cannot be set at the same time",
		)
	}

	if len(config.Prompt.Schemas.Title) == 0 {
		config.Prompt.Schemas.Title = "schemas"
	}

	if len(config.Prompt.Query.Title) == 0 {
		config.Prompt.Query.Title = "user"
	}

	sections := map[string]struct {
		Title string
		Body  string
		Items []string
	}{}

	for _, section := range config.Prompt.Sections {
		sections[strings.ToLower(section.Title)] = section
	}

	if task, ok := sections["task"]; !ok {
		sections["task"] = struct {
			Title string
			Body  string
			Items []string
		}{
			Title: "task",
			Body:  "You are a SQL query builder for a specific knowledge base.",
		}
	} else if len(task.Body) == 0 {
		task.Body = "You are a SQL query builder for a specific knowledge base."
	}

	if instruction, ok := sections["instruction"]; !ok {
		sections["instruction"] = struct {
			Title string
			Body  string
			Items []string
		}{
			Title: "instruction",
			Body:  "Your objective is to generate a valid SQL query from the following user input using the given relational table and prior example of queries.",
		}
	} else if len(instruction.Body) == 0 {
		instruction.Body = "Your objective is to generate a valid SQL query from the following user input using the given relational table and prior example of queries."
	}

	if rules, ok := sections["rules"]; !ok {
		sections["rules"] = struct {
			Title string
			Body  string
			Items []string
		}{
			Title: "rules",
			Items: []string{
				"Only use the tables specified before to generate the SQL query.",
				"Only output a SQL statement without explanation.",
			},
		}
	} else if len(rules.Items) == 0 {
		rules.Items = []string{
			"Only use the tables specified before to generate the SQL query.",
			"Only output a SQL statement without explanation.",
		}
	}

	config.Prompt.Sections = slices.Collect(maps.Values(sections))

	return nil
}

func Parse(configFilePath string) (*Config, error) {
	var config Config

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := validate(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func parseSamples(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	results := strings.Split(string(data), "\n\n")
	samples := []string{}

	for _, sample := range results {
		sample = strings.TrimSpace(sample)

		if strings.HasPrefix(sample, "--") && strings.HasSuffix(sample, ";") {
			samples = append(samples, sample)
		}
	}

	return samples, nil
}
