package memberaccessrefdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllMemberAccessRefInput struct {
	FilterFields  *model.MemberAccessRefFilterFields
	PaginationOpt *model.PaginationOptionInput
}
