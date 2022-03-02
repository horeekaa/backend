package paymentpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreatePaymentUsecaseInput struct {
	Context       context.Context
	CreatePayment *model.CreatePayment
}

type UpdatePaymentUsecaseInput struct {
	Context       context.Context
	UpdatePayment *model.UpdatePayment
}

type GetAllPaymentUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.PaymentFilterFields
	PaginationOps *model.PaginationOptionInput
}
