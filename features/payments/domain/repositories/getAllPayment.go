package paymentdomainrepositoryinterfaces

import (
	paymentdomainrepositorytypes "github.com/horeekaa/backend/features/payments/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllPaymentRepository interface {
	Execute(filterFields paymentdomainrepositorytypes.GetAllPaymentInput) ([]*model.Payment, error)
}
