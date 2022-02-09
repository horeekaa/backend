package invoicedomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllInvoiceInput struct {
	FilterFields  *model.InvoiceFilterFields
	PaginationOpt *model.PaginationOptionInput
}
