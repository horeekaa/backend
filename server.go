package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	authenticationmiddlewares "github.com/horeekaa/backend/http/middlewares/authentication"

	masterdependencies "github.com/horeekaa/backend/dependencies"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"

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
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graphresolver.Resolver{}}))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/graphql", func(r chi.Router) {
				r.Get("/", playground.Handler("GraphQL playground", "/api/v1/graphql/query"))
				r.Route("/query", func(r chi.Router) {
					r.Use(authenticationmiddlewares.AuthGateMiddleware)
					r.Handle("/", srv)
				})
			})
		})
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
