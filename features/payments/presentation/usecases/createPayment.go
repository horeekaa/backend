package paymentpresentationusecaseinterfaces

import (
	paymentpresentationusecasetypes "github.com/horeekaa/backend/features/payments/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type CreatePaymentUsecase interface {
	Execute(input paymentpresentationusecasetypes.CreatePaymentUsecaseInput) (*model.Payment, error)
}
