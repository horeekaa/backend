package memberaccesspresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetMemberAccessUsecase interface {
	Execute(input *model.MemberAccessFilterFields) (*model.MemberAccess, error)
}
