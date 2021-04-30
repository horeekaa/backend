package accountpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type LoginUsecaseInput struct {
	AuthHeader  string
	DeviceToken string
	Context     context.Context
}

type LogoutUsecaseInput struct {
	AuthHeader  string
	DeviceToken string
	Context     context.Context
}

type GetPersonDataFromAccountInput struct {
	AuthHeader      string
	Context         context.Context
	Account         *model.Account
	ViewProfileMode bool
}
