package accountpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type LoginUsecaseInput struct {
	DeviceToken string
	Context     context.Context
}

type LogoutUsecaseInput struct {
	DeviceToken string
	Context     context.Context
}

type GetPersonDataFromAccountInput struct {
	Context         context.Context
	Account         *model.Account
	ViewProfileMode bool
}

type GetAuthUserAndAttachToCtxInput struct {
	AuthHeader string
	Context    context.Context
}

type GetAccountInput struct {
	Context      context.Context
	FilterFields *model.AccountFilterFields
}
