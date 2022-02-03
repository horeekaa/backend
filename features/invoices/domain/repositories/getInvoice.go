package invoicedomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetInvoiceRepository interface {
	Execute(filterFields *model.InvoiceFilterFields) (*model.Invoice, error)
}
