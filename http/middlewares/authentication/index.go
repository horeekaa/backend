package authenticationmiddlewares

import (
	"net/http"

	"github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
)

func AuthGateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var getAuthUserUsecase accountpresentationusecaseinterfaces.GetAuthUserAndAttachToCtxUsecase
			container.Make(&getAuthUserUsecase)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			ctxOutput, err := getAuthUserUsecase.Execute(
				accountpresentationusecasetypes.GetAuthUserAndAttachToCtxInput{
					AuthHeader: authHeader,
					Context:    r.Context(),
				},
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctxOutput))
		},
	)
}
