package invoicedomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createInvoiceRepository struct {
	createInvoiceTransactionComponent invoicedomainrepositoryinterfaces.CreateInvoiceTransactionComponent
	mongoDBTransaction                mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateInvoiceRepository(
	createInvoiceRepositoryTransactionComponent invoicedomainrepositoryinterfaces.CreateInvoiceTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (invoicedomainrepositoryinterfaces.CreateInvoiceRepository, error) {
	createInvoiceRepo := &createInvoiceRepository{
		createInvoiceRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createInvoiceRepo,
		"CreateInvoiceRepository",
	)

	return createInvoiceRepo, nil
}

func (createInvoiceRepo *createInvoiceRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createInvoiceRepo.createInvoiceTransactionComponent.PreTransaction(
		input.(*model.InternalCreateInvoice),
	)
}

func (createInvoiceRepo *createInvoiceRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	invoiceToCreate := input.(*model.InternalCreateInvoice)

	return createInvoiceRepo.createInvoiceTransactionComponent.TransactionBody(
		operationOption,
		invoiceToCreate,
	)
}

func (createInvoiceRepo *createInvoiceRepository) RunTransaction(
	input *model.InternalCreateInvoice,
) ([]*model.Invoice, error) {
	output, err := createInvoiceRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).([]*model.Invoice), nil
}
