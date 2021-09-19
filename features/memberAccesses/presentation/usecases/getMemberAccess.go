package memberaccesspresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetMemberAccessUsecase interface {
	Execute(input *model.InternalMemberAccessFilterFields) (*model.MemberAccess, error)
}
