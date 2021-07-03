package memberaccessrefdomainrepositoryinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessRefUsecaseComponent interface {
	Validation(input *model.InternalCreateMemberAccessRef) (*model.InternalCreateMemberAccessRef, error)
}

type CreateMemberAccessRefRepository interface {
	SetValidation(usecaseComponent CreateMemberAccessRefUsecaseComponent) (bool, error)
	Execute(input *model.InternalCreateMemberAccessRef) (*model.MemberAccessRef, error)
}
