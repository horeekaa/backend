package invoicepresentationusecaseinterfaces

import (
	invoicepresentationusecasetypes "github.com/horeekaa/backend/features/invoices/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type CreateInvoiceUsecase interface {
	Execute(input invoicepresentationusecasetypes.CreateInvoiceUsecaseInput) ([]*model.Invoice, error)
}
