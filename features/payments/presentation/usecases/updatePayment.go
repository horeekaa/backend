package paymentpresentationusecaseinterfaces

import (
	paymentpresentationusecasetypes "github.com/horeekaa/backend/features/payments/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdatePaymentUsecase interface {
	Execute(input paymentpresentationusecasetypes.UpdatePaymentUsecaseInput) (*model.Payment, error)
}
