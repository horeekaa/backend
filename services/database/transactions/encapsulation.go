package databaseservicetransactions

import (
	mongotransactioninterfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/transaction"
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
	databaseservicetransactioninterfaces "github.com/horeekaa/backend/services/database/interfaces/transaction"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type mongoTransactionComponent struct {
	serviceComponent *databaseservicetransactioninterfaces.TransactionComponent
}

func NewMongoTransactionComponent(serviceComponent *databaseservicetransactioninterfaces.TransactionComponent) (mongotransactioninterfaces.TransactionComponent, error) {
	return &mongoTransactionComponent{
		serviceComponent: serviceComponent,
	}, nil
}

func (trxComponent *mongoTransactionComponent) PreTransaction(input interface{}) (interface{}, error) {
	return (*trxComponent.serviceComponent).PreTransaction(input)
}

func (trxComponent *mongoTransactionComponent) TransactionBody(operationOptions *mongooperations.OperationOptions, preOutput interface{}) (interface{}, error) {
	return (*trxComponent.serviceComponent).TransactionBody(
		&databaseserviceoperations.ServiceOptions{
			OperationOptions: operationOptions,
		},
		preOutput,
	)
}
