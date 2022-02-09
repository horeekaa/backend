package invoicedomainrepositoryinterfaces

import (
	invoicedomainrepositorytypes "github.com/horeekaa/backend/features/invoices/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllInvoiceRepository interface {
	Execute(filterFields invoicedomainrepositorytypes.GetAllInvoiceInput) ([]*model.Invoice, error)
}
