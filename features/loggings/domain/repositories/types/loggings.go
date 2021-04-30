package loggingdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type LogEntityProposalActivityInput struct {
	CollectionName   string
	CreatedByAccount *model.Account
	CreatorInitial   string
	Activity         model.LoggedActivity
	ProposalStatus   model.EntityProposalStatus
	ExistingObject   *interface{}
	ExistingObjectID *string
	NewObject        *interface{}
}

type LogEntityApprovalActivityInput struct {
	PreviousLog      *model.Logging
	ApprovingAccount *model.Account
	ApproverInitial  string
	ApprovalStatus   model.EntityProposalStatus
}
