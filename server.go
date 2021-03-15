package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"

	configs "github.com/horeekaa/backend/_commons/configs"
	"github.com/horeekaa/backend/graph/generated"
	graphresolver "github.com/horeekaa/backend/graph/resolver"
	mongodbclients "github.com/horeekaa/backend/repositories/databaseClient/mongoDB"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
  
	timeout, _ := strconv.Atoi(configs.GetEnvVariable(configs.DbConfigTimeout))
	repository, _ := mongodbclients.NewMongoClientRef(
		configs.GetEnvVariable(configs.DbConfigURL),
		configs.GetEnvVariable(configs.DbConfigDBName),
		timeout,
	)
	mongodbclients.DatabaseClient = repository

	router := chi.NewRouter()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graphresolver.Resolver{}}))

	router.Handle("/api/v1/graphql", playground.Handler("GraphQL playground", "/api/v1/graphql/query"))
	router.Handle("/api/v1/graphql/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
