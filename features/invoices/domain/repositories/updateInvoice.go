package invoicedomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type UpdateInvoiceTransactionComponent interface {
	PreTransaction(
		updateInvoiceInput *model.InternalUpdateInvoice,
	) (*model.InternalUpdateInvoice, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateInvoiceInput *model.InternalUpdateInvoice,
	) (*model.Invoice, error)
}

type UpdateInvoiceRepository interface {
	RunTransaction(
		updateInvoiceInput *model.InternalUpdateInvoice,
	) (*model.Invoice, error)
}
