package main

import (
	"fmt"
	"os"

	"github.com/buttaciemanuel/sqlify/context"
	"github.com/buttaciemanuel/sqlify/datasource"
	"github.com/buttaciemanuel/sqlify/embed"
	"github.com/buttaciemanuel/sqlify/pipeline"
	"github.com/buttaciemanuel/sqlify/prompt"
)

// func main() {
// 	store, err := context.NewQdrantStore()

// 	if err != nil {
// 		panic(err)
// 	}

// 	store.Clear()

// 	if err := store.Initialize([]string{"schema", "example"}); err != nil {
// 		panic(err)
// 	}

// 	// ---

// 	tableBuilders := []*context.TableBuilder{
// 		context.NewTable(
// 			"phone_purchases",
// 			"This table stores the purchases of mobile phones performed by each client along with date and price",
// 		).AddAttribute(
// 			"model",
// 			"VARCHAR",
// 			"Mobile phone or smarphone model",
// 		).AddAttribute(
// 			"client",
// 			"VARCHAR",
// 			"Buyer full name",
// 		).AddAttribute(
// 			"date",
// 			"DATETIME",
// 			"Date and time of the purchase",
// 		).AddAttribute(
// 			"price",
// 			"REAL",
// 			"Price in USD dollars of the purchased phone",
// 		),
// 		context.NewTable(
// 			"gym_subscriptions",
// 			"This table stores the gym subscriptions performed by each client along with date and price",
// 		).AddAttribute(
// 			"gym",
// 			"VARCHAR",
// 			"Name (or branding) of the gym",
// 		).AddAttribute(
// 			"client",
// 			"VARCHAR",
// 			"Gym customer full name",
// 		).AddAttribute(
// 			"date",
// 			"DATETIME",
// 			"Date and time of the purchase of the subscription",
// 		).AddAttribute(
// 			"price",
// 			"REAL",
// 			"Price in USD dollars of the subscription plan",
// 		),
// 	}

// 	for _, tableBuilder := range tableBuilders {
// 		if err := store.StoreDocument("schema", tableBuilder.Build()); err != nil {
// 			panic(err)
// 		}
// 	}

// 	// ---

// 	queryBuilders := []*context.QueryBuilder{
// 		context.NewQuery(
// 			"SELECT DISTINCT(model) FROM phone_purchases WHERE date < '2025-01-01';",
// 			"List of mobile phones sold before January 2025",
// 		),
// 		context.NewQuery(
// 			"SELECT AVG(price) FROM gym_subscriptions WHERE datetime >= '2025-06-21' AND datetime < '2025-09-21';",
// 			"Show the average gym price during summer 2025",
// 		),
// 	}

// 	for _, queryBuilder := range queryBuilders {
// 		if err := store.StoreDocument("example", queryBuilder.Build()); err != nil {
// 			panic(err)
// 		}
// 	}

// 	// ---

// 	generator := prompt.NewDefaultPrompt(1)
// 	modelPrompt, err := generator.Build(os.Args[1], store)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(modelPrompt)

// 	// ---

// 	generatedStatement, err := ollama.Generate("llama3.2", modelPrompt)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("\n%s\n", generatedStatement)

// 	// ---

// 	database, err := datasource.DuckDb("../../database.duckdb")

// 	if err != nil {
// 		panic(err)
// 	}

// 	queryResult, err := database.RunQuery(generatedStatement)

// 	if err != nil {
// 		panic(err)
// 	}

// 	jsonQueryResult, err := json.Marshal(queryResult)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("\n%s\n", string(jsonQueryResult))

// 	// ---

// 	tablesResult, err := database.GetTables()

// 	if err != nil {
// 		panic(err)
// 	}

// 	jsonTablesResult, err := json.Marshal(tablesResult)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("\n%s\n", string(jsonTablesResult))

// 	// ---

// 	tableResult, err := database.GetTableSchema("phone_purchases")

// 	if err != nil {
// 		panic(err)
// 	}

// 	jsonTableResult, err := json.Marshal(tableResult)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("\n%s\n", string(jsonTableResult))
// }

func main() {
	context, err := context.NewQdrantStore(embed.SnowflakeArticEmbed)

	if err != nil {
		panic(err)
	}

	if err := context.Clear(); err != nil {
		panic(err)
	}

	database, err := datasource.DuckDb("../../database.duckdb")

	if err != nil {
		panic(err)
	}

	assistant, err := pipeline.New(&pipeline.Config{
		Context:  context,
		Database: database,
		Prompt: prompt.Prompt{
			Sections: []prompt.TextSection{
				{
					Title: "Task",
					Body:  "You are a SQL query builder for a specific knowledge base.",
				},
				{
					Title: "Rules",
					Items: []string{
						"Only use the tables specified before to generate the SQL query.",
						"Only output a SQL statement without explanation.",
					},
				},
				{
					Title: "Instruction",
					Body:  "Your objective is to generate a valid SQL query from the following user input using the given relational table and prior example of queries.",
				},
			},
			Query: prompt.TextSection{
				Title: "User",
			},
			Schemas: prompt.CodeSection{
				Title: "Schemas",
			},
			Examples: prompt.CodeSection{
				Title: "Examples",
			},
		},
		Model:      "llama3.2",
		AutoSchema: true,
	})

	if err != nil {
		panic(err)
	}

	results, err := assistant.Execute(os.Args[1])

	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", results)
}
