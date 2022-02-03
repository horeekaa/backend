package invoicepresentationusecaseinterfaces

import (
	invoicepresentationusecasetypes "github.com/horeekaa/backend/features/invoices/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdateInvoiceUsecase interface {
	Execute(input invoicepresentationusecasetypes.UpdateInvoiceUsecaseInput) (*model.Invoice, error)
}
