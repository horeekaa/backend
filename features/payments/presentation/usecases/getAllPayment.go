package paymentpresentationusecaseinterfaces

import (
	paymentpresentationusecasetypes "github.com/horeekaa/backend/features/payments/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllPaymentUsecase interface {
	Execute(input paymentpresentationusecasetypes.GetAllPaymentUsecaseInput) ([]*model.Payment, error)
}
