package embed

var (
	SnowflakeArticEmbed Embedder = Embedder{name: "snowflake-arctic-embed:22m", embedSize: 384}
	MxbaiEmbedLarge     Embedder = Embedder{name: "mxbai-embed-large", embedSize: 1024}
)
