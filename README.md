# SQLIFY

Sqlify is a barebone framework to simulate the functions of a retrieval augmented generation assistant. Once you provide a database schema and optionally a bunch of sample SQL queries, it helps you convert your human text questions to SQL queries.

## Introduction

There are two ways to run Sqlify locally, either as an HTTP server waiting for human text queries, or as a single run application to which ask a question. In either cases, you will need to provide a configuration file that contains the settings for Sqlify to properly set itself up.

```yaml title="/data/config.yaml"
context:
  store: qdrant
  embedder: snowflake-arctic-embed:22m
database:
  autoschema: true
  duckdb:
    filename: /data/database.duckdb
model:
  name: phi4
  url: ollama:11435/api/generate
prompt:
  examples:
    from: /data/examples.sql
```

By default, Sqlify will use a Duckdb database instance to execute the generated queries and fetch data. On the other hand, Qdrant is used as the vector database to store metadata (sample queries and schemas) by similarity for fast retrieval. The `database.autoschema` options allows to read table schemas directly from the database. The following file `examples.sql` contains sample queries in the following format.

```sql title="/data/examples.sql"
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

The dockerized process runs Sqlify as a server.

### Run as a server

Questions can now be asked by firing POST questions to the server.

```sh
curl -X POST -H "Content-Type: application/json" -H "Accept: application/json" http://localhost:3001 -d "{ \"query\": \"<YOUR QUESTION HERE>\" }" | jq .
```
