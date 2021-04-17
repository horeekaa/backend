package loggingdomainrepositoryinterfaces

import (
	loggingdomainrepositorytypes "github.com/horeekaa/backend/features/loggings/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type LogEntityApprovalActivityUsecaseComponent interface {
	Validation(input loggingdomainrepositorytypes.LogEntityApprovalActivityInput) (loggingdomainrepositorytypes.LogEntityApprovalActivityInput, error)
}

type LogEntityApprovalActivityRepository interface {
	SetValidation(usecaseComponent LogEntityApprovalActivityUsecaseComponent) (bool, error)
	Execute(input loggingdomainrepositorytypes.LogEntityApprovalActivityInput) (*model.Logging, error)
}
