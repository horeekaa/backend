package memberaccessrefdomainrepositoryinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessRefUsecaseComponent interface {
	Validation(input *model.CreateMemberAccessRef) (*model.CreateMemberAccessRef, error)
}

type CreateMemberAccessRefRepository interface {
	SetValidation(usecaseComponent CreateMemberAccessRefUsecaseComponent) (bool, error)
	Execute(input *model.CreateMemberAccessRef) (*model.MemberAccessRef, error)
}
