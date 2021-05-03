package main

import (
	"log"
	"net/http"
	"os"

	masterdependencies "github.com/horeekaa/backend/dependencies"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"

	"github.com/horeekaa/backend/graph/generated"
	graphresolver "github.com/horeekaa/backend/graph/resolver"
)

const defaultPort = "8080"

func main() {
	masterBind := &masterdependencies.MasterDependency{}
	masterBind.Bind()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	router := chi.NewRouter()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graphresolver.Resolver{}}))

	router.Handle("/api/v1/graphql", playground.Handler("GraphQL playground", "/api/v1/graphql/query"))
	router.Handle("/api/v1/graphql/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
