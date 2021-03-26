package accountdomainrepositoryinterfaces

import (
	accountrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	model "github.com/horeekaa/backend/model"
)

type GetPersonDataFromAccountUsecaseComponent interface {
	Validation(input model.Account) (*model.Account, error)
}

type GetPersonDataFromAccountRepository interface {
	SetValidation(usecaseComponent GetPersonDataFromAccountUsecaseComponent) (bool, error)
	Execute(input model.Account) (*accountrepositorytypes.GetPersonDataByAccountOutput, error)
}
