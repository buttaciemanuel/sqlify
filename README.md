# SQLIFY

...

## Embedding model

Install the Snowflake model within the Ollama suite to create text embeddings.

```sh
ollama pull snowflake-arctic-embed:22m
```

## Qdrant

Run Qdrant vector database as Dockerized container.

```sh
docker run -p 6333:6333 -p 6334:6334 -v "$(pwd)/qdrant_storage:/qdrant/storage:z" qdrant/qdrant
```
