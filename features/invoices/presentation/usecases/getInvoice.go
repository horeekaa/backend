package invoicepresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetInvoiceUsecase interface {
	Execute(input *model.InvoiceFilterFields) (*model.Invoice, error)
}
