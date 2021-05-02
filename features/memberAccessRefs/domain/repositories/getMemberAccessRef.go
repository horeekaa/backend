package memberaccessrefdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetMemberAccessRefRepository interface {
	Execute(filterFields *model.MemberAccessRefFilterFields) (*model.MemberAccessRef, error)
}
