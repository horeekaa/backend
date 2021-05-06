package accountpresentationusecasetypes

import (
	"context"

	authenticationcoremodels "github.com/horeekaa/backend/core/authentication/models"
	"github.com/horeekaa/backend/model"
)

type LoginUsecaseInput struct {
	User        *authenticationcoremodels.AuthUserWrap
	DeviceToken string
	Context     context.Context
}

type LogoutUsecaseInput struct {
	User        *authenticationcoremodels.AuthUserWrap
	DeviceToken string
	Context     context.Context
}

type GetPersonDataFromAccountInput struct {
	User            *authenticationcoremodels.AuthUserWrap
	Context         context.Context
	Account         *model.Account
	ViewProfileMode bool
}
