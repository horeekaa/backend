package graphqlroutes

import (
	"context"
	"errors"

	"github.com/go-chi/chi"

	horeekaacorebaseerror "github.com/horeekaa/backend/core/errors/errors/base"
	authenticationmiddlewares "github.com/horeekaa/backend/http/middlewares/authentication"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/horeekaa/backend/graph/generated"
	graphresolver "github.com/horeekaa/backend/graph/resolver"
)

func Route(r chi.Router) {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graphresolver.Resolver{}}))
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)

		var myErr *horeekaacorebaseerror.Error
		if errors.As(e, &myErr) {
			err.Message = myErr.Code
		}

		return err
	})

	r.Get("/", playground.Handler("GraphQL playground", "/api/v1/graphql/query"))
	r.Route("/query", func(r chi.Router) {
		r.Use(authenticationmiddlewares.AuthGateMiddleware)
		r.Handle("/", srv)
	})
}
