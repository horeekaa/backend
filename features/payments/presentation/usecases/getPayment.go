package paymentpresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetPaymentUsecase interface {
	Execute(input *model.PaymentFilterFields) (*model.Payment, error)
}
