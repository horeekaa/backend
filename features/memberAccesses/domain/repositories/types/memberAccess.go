package memberaccessdomainrepositorytypes

import (
	"github.com/horeekaa/backend/model"
)

type GetAccountMemberAccessInput struct {
	Account                *model.Account
	MemberAccessRefType    model.MemberAccessRefType
	MemberAccessRefOptions model.MemberAccessRefOptionsInput
}

type UpdateMemberAccessOutput struct {
	PreviousMemberAccess *model.MemberAccess
	UpdatedMemberAccess  *model.MemberAccess
}

type GetAllMemberAccessInput struct {
	FilterFields  *model.MemberAccessFilterFields
	PaginationOpt *model.PaginationOptionInput
}
