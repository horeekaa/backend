package invoicedomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateInvoiceTransactionComponent interface {
	PreTransaction(
		createInvoiceInput *model.InternalCreateInvoice,
	) (*model.InternalCreateInvoice, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createInvoiceInput *model.InternalCreateInvoice,
	) ([]*model.Invoice, error)

	GenerateNewObjectID() primitive.ObjectID
	GetCurrentObjectID() primitive.ObjectID
}

type CreateInvoiceRepository interface {
	RunTransaction(
		createInvoiceInput *model.InternalCreateInvoice,
	) ([]*model.Invoice, error)
}
