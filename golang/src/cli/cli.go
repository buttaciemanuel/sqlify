package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"

	"github.com/buttaciemanuel/sqlify/context"
	"github.com/buttaciemanuel/sqlify/datasource"
	"github.com/buttaciemanuel/sqlify/embed"
	"github.com/buttaciemanuel/sqlify/pipeline"
	"github.com/buttaciemanuel/sqlify/prompt"
)

func Run() error {
	subcommands := []*cobra.Command{
		{
			Use:   "serve",
			Short: "Run sqlify as a server to query information to",
			RunE:  serve,
		},
		{
			Use:   "query [query]",
			Short: "Ask a direct question to sqlify",
			Args:  cobra.ExactArgs(1),
			RunE:  query,
		},
	}

	subcommands[0].PersistentFlags().
		String("configuration", "config.yaml", "Path to the configuration file")
	subcommands[0].PersistentFlags().
		Uint("port", 3000, "Port number for the server")

	subcommands[1].PersistentFlags().
		String("configuration", "config.yaml", "Path to the configuration file")

	command := &cobra.Command{
		Use:   "sqlify [question]",
		Short: "Sqlify is a cli tool to extract relevant information from a database through natural language queries",
		Long:  "Sqlify connects to your database, and allows for the generation of sql statement using your provided large language model.",
		Args:  cobra.ExactArgs(1),
	}

	command.AddCommand(subcommands...)

	if err := command.Execute(); err != nil {
		return err
	}

	return nil
}

func query(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString("configuration")
	if err != nil {
		return err
	}

	config, err := configure(path)
	if err != nil {
		return err
	}

	assistant, err := pipeline.New(config)
	if err != nil {
		return err
	}

	results, err := assistant.Execute(args[0])
	if err != nil {
		return err
	}

	output := fmt.Sprintf("%v", results)

	if len(output) > 1024 {
		fmt.Printf("%s...\n", output[:1024])
	} else {
		fmt.Printf("%s\n", output)
	}

	return nil
}

func serve(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString("configuration")
	if err != nil {
		return err
	}

	config, err := configure(path)
	if err != nil {
		return err
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Query string `json:"query"`
		}

		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		assistant, err := pipeline.New(config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("%s", body.Query)

		results, err := assistant.Execute(body.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	})

	port, err := cmd.Flags().GetUint("port")
	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%v", port), router)
}

func configure(path string) (*pipeline.Config, error) {
	config, err := Parse(path)
	if err != nil {
		return nil, err
	}

	var embeddingModel embed.Embedder

	if config.Context.Embedder == embed.SnowflakeArticEmbed.Name() {
		embeddingModel = embed.SnowflakeArticEmbed
	} else if config.Context.Embedder == embed.MxbaiEmbedLarge.Name() {
		embeddingModel = embed.MxbaiEmbedLarge
	}

	var modelContext context.DocumentStore

	if config.Context.Store == "qdrant" {
		connection, err := context.Qdrant(embeddingModel)
		if err != nil {
			return nil, err
		}

		modelContext = connection
	}

	// if err := modelContext.Clear(); err != nil {
	// 	return nil, err
	// }

	var modelDatabase *datasource.Database

	if len(config.Database.Duckdb.Filename) > 0 {
		connection, err := datasource.Duckdb(config.Database.Duckdb.Filename)
		if err != nil {
			return nil, err
		}

		modelDatabase = connection
	}

	modelPrompt := prompt.Prompt{
		Sections: []prompt.TextSection{},
		Query: prompt.TextSection{
			Title: config.Prompt.Query.Title,
		},
		Schemas: prompt.CodeSection{
			Title: config.Prompt.Schemas.Title,
		},
		Examples: prompt.CodeSection{
			Title: config.Prompt.Examples.Title,
		},
	}

	for _, schema := range config.Prompt.Schemas.Items {
		modelPrompt.Schemas.Samples = append(
			modelPrompt.Schemas.Samples,
			context.Document{
				Value: strings.TrimSpace(schema),
			},
		)
	}

	if len(config.Prompt.Examples.From) > 0 {
		samples, err := parseSamples(config.Prompt.Examples.From)
		if err != nil {
			return nil, err
		}

		config.Prompt.Examples.Items = samples
	}

	for _, example := range config.Prompt.Examples.Items {
		modelPrompt.Examples.Samples = append(
			modelPrompt.Examples.Samples,
			context.Document{
				Value: strings.TrimSpace(example),
			},
		)
	}

	for _, section := range config.Prompt.Sections {
		modelPrompt.Sections = append(modelPrompt.Sections, prompt.TextSection{
			Title: section.Title,
			Body:  section.Body,
			Items: section.Items,
		})
	}

	return &pipeline.Config{
		Context:    modelContext,
		Database:   modelDatabase,
		Prompt:     modelPrompt,
		Model:      config.Model,
		AutoSchema: config.Database.Autoschema,
	}, nil
}
