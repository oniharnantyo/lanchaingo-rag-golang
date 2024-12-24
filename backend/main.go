package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	document2 "github.com/oniharnantyo/lanchaingo-rag-golang/document"
	ollama2 "github.com/oniharnantyo/lanchaingo-rag-golang/ollama"
	"github.com/oniharnantyo/lanchaingo-rag-golang/query"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading env file: %q", err)
	}

	ollamaClient, err := ollama.New(
		ollama.WithModel(os.Getenv("EMBEDDINGS_MODEL")),
		ollama.WithServerURL(os.Getenv("OLLAMA_HOST")),
	)

	ollamaEmbedder := ollama2.NewOllama(ollamaClient)

	ollamaLLM, err := ollama.New(ollama.WithModel(
		os.Getenv("LLM_MODEL")),
		ollama.WithFormat("json"),
	)
	if err != nil {
		log.Fatalf("ollama llm error %q", err)
	}

	//weaviateStore, err := weaviate.New(
	//	weaviate.WithScheme("http"),
	//	weaviate.WithHost(os.Getenv("WEAVIATE_HOST")),
	//	weaviate.WithIndexName("Documents"),
	//	weaviate.WithEmbedder(ollamaEmbedder),
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//mongoClient, err := mongo.Connect(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	//if err != nil {
	//	log.Fatalf("mongo client err: %q", err)
	//}
	//
	//mongoHistoriesColl := mongoClient.Database(os.Getenv("MONGODB_DB")).Collection("histories")
	//
	//mongoStore := mongovector.New(mongoDocumentsColl, ollamaEmbedder)

	pgConncetor, err := pgvector.New(ctx,
		pgvector.WithConnectionURL(os.Getenv("POSTGRES_URI")),
		pgvector.WithEmbeddingTableName("documents"),
		pgvector.WithEmbedder(ollamaEmbedder))
	if err != nil {
		log.Fatalf("postgres err: %q", err)
	}

	retriever := vectorstores.ToRetriever(&pgConncetor, 5, vectorstores.WithEmbedder(ollamaEmbedder))

	documentService := document2.NewService(textsplitter.NewTokenSplitter(), &pgConncetor)
	queryService := query.NewService(
		ollamaLLM, retriever, os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_DB"), "histories")

	documentHandler := document2.NewHandler(documentService)
	queryHandler := query.NewHandler(queryService)

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, "pong")
	})
	e.POST("/documents", documentHandler.AddDocument)
	e.POST("/queries", queryHandler.Query)

	log.Fatal(e.Start(":8000"))
}
