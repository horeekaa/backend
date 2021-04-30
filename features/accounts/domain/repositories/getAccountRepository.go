package accountdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetAccountRepository interface {
	Execute(filterFields *model.AccountFilterFields) (*model.Account, error)
}
