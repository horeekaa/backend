package paymentdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetPaymentRepository interface {
	Execute(filterFields *model.PaymentFilterFields) (*model.Payment, error)
}
