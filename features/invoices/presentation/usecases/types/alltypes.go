package invoicepresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreateInvoiceUsecaseInput struct {
	Context           context.Context
	CreateInvoice     *model.CreateInvoice
	CronAuthenticated bool
}

type UpdateInvoiceUsecaseInput struct {
	Context           context.Context
	UpdateInvoice     *model.UpdateInvoice
	CronAuthenticated bool
}

type GetAllInvoiceUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.InvoiceFilterFields
	PaginationOps *model.PaginationOptionInput
}
