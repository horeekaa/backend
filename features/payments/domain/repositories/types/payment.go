package paymentdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllPaymentInput struct {
	FilterFields  *model.PaymentFilterFields
	PaginationOpt *model.PaginationOptionInput
}
