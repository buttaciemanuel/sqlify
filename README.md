# SQLIFY

Sqlify is a barebone framework to simulate the functions of a retrieval augmented generation assistant. Once you provide a database schema and optionally a bunch of sample SQL queries, it helps you convert your human text questions to SQL queries.

## Introduction

There are two ways to run Sqlify locally, either as an HTTP server waiting for human text queries, or as a single run application to which ask a question. In either cases, you will need to provide a configuration file that contains the settings for Sqlify to properly set itself up.

```yaml title="config.yaml"
context:
  store: qdrant
  embedder: snowflake-arctic-embed:22m
database:
  autoschema: true
  duckdb:
    filename: /your/path/to/database.duckdb
model:
  name: phi4
  url: localhost:11435/api/generate
prompt:
  examples:
    from: examples.sql
```

By default, Sqlify will use a Duckdb database instance to execute the generated queries and fetch data. On the other hand, Qdrant is used as the vector database to store metadata (sample queries and schemas) by similarity for fast retrieval. The `database.autoschema` options allows to read table schemas directly from the database. The following file `examples.sql` contains sample queries in the following format.

```sql title="examples.sql"
-- <YOUR QUESTION>
SELECT ...;

-- <ANOTHER QUESTION OF YOURS>
SELECT ...;
```

## Run

Luckily, we can launch the Qdrant vector store and Ollama client as a docker container.

```sh
docker compose up -d
```

Place yourself in the source directory.

```sh
cd golang/src
```

### Run as a command line application

To ask a question to Sqlify you can do as follows.

```sh
go run . query --configuration path/to/your/config.yaml <YOUR QUESTION HERE>
```

Running such a command with your questions will yield an output containing

1. The constructed prompt for the large language model as an XML prompt with similar questions and related table schemas.
2. The generated SQL statement corresponding to the question you typed in.
3. The data results (if the statement is correct and any) given by executing such a query on the Duckdb database instance.

### Run as a server

Otherwise, you can run Sqlify as a server.

```sh
go run . serve --configuration path/to/your/config.yaml --port 3001
```

Questions can now be asked by firing POST questions to the server.

```sh
curl -X POST -H "Content-Type: application/json" -H "Accept: application/json" http://localhost:3001 -d "{ \"query\": \"<YOUR QUESTION HERE>\" }" | jq .
```
