package mongodbtransactionfixtures

import (
	"time"

	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var layout = "2006-01-02T15:04:05.000Z"
var str = "2014-11-12T11:45:26.371Z"
var t, err = time.Parse(layout, str)

var TransactionOutput model.MemberAccess = model.MemberAccess{
	ID:                 primitive.ObjectID{byte(255)},
	InvitationAccepted: false,
	Account:            &model.Account{ID: primitive.ObjectID{byte(255)}},
	Status:             "ACTIVE",
	Organization:       &model.Organization{ID: primitive.ObjectID{byte(255)}},
	CreatedAt:          t,
	OrganizationMembershipRole: func(r model.OrganizationMembershipRole) *model.OrganizationMembershipRole {
		return &r
	}(model.OrganizationMembershipRoleManager),
	MemberAccessRefType: model.MemberAccessRefTypeOrganizationsBased,
	Access: &model.MemberAccessRefOptions{
		OrganizationAccesses: &model.OrganizationAccesses{
			OrganizationApproval:  func(b bool) *bool { return &b }(true),
			OrganizationCreate:    func(b bool) *bool { return &b }(true),
			OrganizationReadOwned: func(b bool) *bool { return &b }(true),
			OrganizationUpdate:    func(b bool) *bool { return &b }(true),
		},
	},
	SubmittingAccount: &model.Account{ID: primitive.ObjectID{byte(255)}},
	ProposalStatus:    model.EntityProposalStatusProposed,
}
