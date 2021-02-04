package databaseservicetransaction

import (
	mongotransactioninterface "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/transaction"
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
	databaseserviceinterface "github.com/horeekaa/backend/services/database/interfaces/repos"
	databaseservicetransactioninterface "github.com/horeekaa/backend/services/database/interfaces/transaction"
)

type mongoTransactionComponent struct {
	serviceComponent *databaseservicetransactioninterface.TransactionComponent
}

func NewMongoTransactionComponent(serviceComponent *databaseservicetransactioninterface.TransactionComponent) (mongotransactioninterface.TransactionComponent, error) {
	return &mongoTransactionComponent{
		serviceComponent: serviceComponent,
	}, nil
}

func (trxComponent *mongoTransactionComponent) PreTransaction(input interface{}) (interface{}, error) {
	return (*trxComponent.serviceComponent).PreTransaction(input)
}

func (trxComponent *mongoTransactionComponent) TransactionBody(operationOptions *mongooperations.OperationOptions, preOutput interface{}) (interface{}, error) {
	return (*trxComponent.serviceComponent).TransactionBody(
		&databaseserviceinterface.ServiceOptions{
			OperationOptions: operationOptions,
		},
		preOutput,
	)
}
