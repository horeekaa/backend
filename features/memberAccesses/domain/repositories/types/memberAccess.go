package memberaccessdomainrepositorytypes

import (
	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessForAccountInput struct {
	Account                    *model.Account
	MemberAccessRefType        model.MemberAccessRefType
	OrganizationMembershipRole model.OrganizationMembershipRole
	OrganizationType           model.OrganizationType
	Organization               *model.Organization
}

type GetAccountMemberAccessInput struct {
	Account                *model.Account
	MemberAccessRefType    model.MemberAccessRefType
	MemberAccessRefOptions model.MemberAccessRefOptionsInput
}
