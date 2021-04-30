package memberaccessrefdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type UpdateMemberAccessRefOutput struct {
	PreviousMemberAccessRef *model.MemberAccessRef
	UpdatedMemberAccessRef  *model.MemberAccessRef
}

type GetAllMemberAccessRefInput struct {
	FilterFields  *model.MemberAccessRefFilterFields
	PaginationOpt *model.PaginationOptionInput
}
