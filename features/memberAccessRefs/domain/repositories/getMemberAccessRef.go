package memberaccessrefdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetMemberAccessRefRepository interface {
	Execute(filterFields *model.UpdateMemberAccessRef) (*model.MemberAccessRef, error)
}
