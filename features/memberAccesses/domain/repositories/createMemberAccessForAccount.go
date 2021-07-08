package memberaccessdomainrepositoryinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessForAccountUsecaseComponent interface {
	Validation(input *model.InternalCreateMemberAccess) (*model.InternalCreateMemberAccess, error)
}

type CreateMemberAccessForAccountRepository interface {
	SetValidation(usecaseComponent CreateMemberAccessForAccountUsecaseComponent) (bool, error)
	Execute(input *model.InternalCreateMemberAccess) (*model.MemberAccess, error)
}
