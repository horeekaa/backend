package invoicepresentationusecaseinterfaces

import (
	invoicepresentationusecasetypes "github.com/horeekaa/backend/features/invoices/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllInvoiceUsecase interface {
	Execute(input invoicepresentationusecasetypes.GetAllInvoiceUsecaseInput) ([]*model.Invoice, error)
}
