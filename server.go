package main

import (
	"log"
	"net/http"
	"os"

	graphqlroutes "github.com/horeekaa/backend/http/routes/graphql"
	scheduledjobroutes "github.com/horeekaa/backend/http/routes/scheduledJob"

	"github.com/joho/godotenv"

	masterdependencies "github.com/horeekaa/backend/dependencies"

	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Cannot load .env file!")
	}

	masterBind := &masterdependencies.MasterDependency{}
	masterBind.Bind()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/graphql", graphqlroutes.Route)
		})
	})

	router.Route("/scheduledJob", func(r chi.Router) {
		r.Route("/v1", scheduledjobroutes.Route)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
