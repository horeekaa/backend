package memberaccessdomainrepositorytypes

import (
	"github.com/horeekaa/backend/model"
)

type GetAccountMemberAccessInput struct {
	MemberAccessFilterFields *model.MemberAccessFilterFields
	QueryMode                bool
}

type GetAllMemberAccessInput struct {
	FilterFields  *model.MemberAccessFilterFields
	PaginationOpt *model.PaginationOptionInput
}
