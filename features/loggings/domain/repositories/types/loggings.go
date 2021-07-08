package loggingdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type LogEntityProposalActivityInput struct {
	CollectionName   string
	CreatedByAccount *model.Account
	Activity         model.LoggedActivity
	ProposalStatus   model.EntityProposalStatus
	ExistingObject   *interface{}
	ExistingObjectID *model.ObjectIDOnly
	NewObject        *interface{}
}

type LogEntityApprovalActivityInput struct {
	PreviousLog      *model.Logging
	ApprovingAccount *model.Account
	ApprovalStatus   model.EntityProposalStatus
}
