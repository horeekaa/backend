package memberaccessrefpresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetMemberAccessRefUsecase interface {
	Execute(input *model.MemberAccessRefFilterFields) (*model.MemberAccessRef, error)
}
