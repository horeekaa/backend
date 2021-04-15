package loggingdomainrepositoryinterfaces

import (
	loggingdomainrepositorytypes "github.com/horeekaa/backend/features/loggings/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type LogEntityProposalActivityUsecaseComponent interface {
	Validation(input loggingdomainrepositorytypes.LogEntityProposalActivityInput) (loggingdomainrepositorytypes.LogEntityProposalActivityInput, error)
}

type LogEntityProposalActivityRepository interface {
	SetValidation(usecaseComponent LogEntityProposalActivityUsecaseComponent) (bool, error)
	Execute(input loggingdomainrepositorytypes.LogEntityProposalActivityInput) (*model.Logging, error)
}
