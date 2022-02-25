package invoicedomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	invoicedomainrepositorytypes "github.com/horeekaa/backend/features/invoices/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateInvoiceTransactionComponent interface {
	PreTransaction(
		createInvoiceInput *invoicedomainrepositorytypes.CreateInvoiceInput,
	) (*invoicedomainrepositorytypes.CreateInvoiceInput, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createInvoiceInput *invoicedomainrepositorytypes.CreateInvoiceInput,
	) (*model.Invoice, error)

	GenerateNewObjectID() primitive.ObjectID
	GetCurrentObjectID() primitive.ObjectID
}

type CreateInvoiceRepository interface {
	RunTransaction(
		createInvoiceInput *model.InternalCreateInvoice,
	) ([]*model.Invoice, error)
}
