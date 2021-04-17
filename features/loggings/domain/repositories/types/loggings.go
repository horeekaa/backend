package loggingdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type LogEntityProposalActivityInput struct {
	CollectionName   string
	ChangedBy        model.Person
	Activity         model.LoggedActivity
	ProposalStatus   model.EntityProposalStatus
	ExistingObject   *interface{}
	ExistingObjectID *string
	NewObject        *interface{}
}

type LogEntityApprovalActivityInput struct {
	PreviousLog    *model.Logging
	ApprovedBy     *model.Person
	ApprovalStatus model.EntityProposalStatus
}
