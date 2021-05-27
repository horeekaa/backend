package memberaccessdomainrepositoryinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessForAccountUsecaseComponent interface {
	Validation(input *model.CreateMemberAccess) (*model.CreateMemberAccess, error)
}

type CreateMemberAccessForAccountRepository interface {
	SetValidation(usecaseComponent CreateMemberAccessForAccountUsecaseComponent) (bool, error)
	Execute(input *model.CreateMemberAccess) (*model.MemberAccess, error)
}
