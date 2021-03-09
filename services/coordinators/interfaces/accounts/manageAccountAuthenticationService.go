package accountservicecoordinatorinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type ManageAccountAuthenticationUsecaseComponent interface {
	Validation(authHeader string) (string, error)
}

type ManageAccountAuthenticationService interface {
	RunTransaction(authHeader string) (*model.Account, error)
}
