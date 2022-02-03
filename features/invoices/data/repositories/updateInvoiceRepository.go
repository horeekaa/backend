package invoicedomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type updateInvoiceRepository struct {
	updateInvoiceTransactionComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent
	mongoDBTransaction                mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewUpdateInvoiceRepository(
	updateInvoiceRepositoryTransactionComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (invoicedomainrepositoryinterfaces.UpdateInvoiceRepository, error) {
	updateInvoiceRepo := &updateInvoiceRepository{
		updateInvoiceRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		updateInvoiceRepo,
		"updateInvoiceRepository",
	)

	return updateInvoiceRepo, nil
}

func (updateInvoiceRepo *updateInvoiceRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateInvoiceRepo.updateInvoiceTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateInvoice),
	)
}

func (updateInvoiceRepo *updateInvoiceRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	invoiceToCreate := input.(*model.InternalUpdateInvoice)

	return updateInvoiceRepo.updateInvoiceTransactionComponent.TransactionBody(
		operationOption,
		invoiceToCreate,
	)
}

func (updateInvoiceRepo *updateInvoiceRepository) RunTransaction(
	input *model.InternalUpdateInvoice,
) (*model.Invoice, error) {
	output, err := updateInvoiceRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.Invoice), nil
}
