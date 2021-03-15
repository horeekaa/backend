package databaseservicetransactions

import (
	horeekaaexceptiontofailure "github.com/horeekaa/backend/_errors/serviceFailures/_exceptionToFailure"
	mongodbclients "github.com/horeekaa/backend/repositories/databaseClient/mongoDB"
	mongotransactioninterfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/transaction"
	mongotransactions "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/transactions"
	databaseservicetransactioninterfaces "github.com/horeekaa/backend/services/database/interfaces/transaction"
)

type databaseServiceTransaction struct {
	MongoTransaction *mongotransactioninterfaces.MongoRepoTransaction
}

func NewDatabaseServiceTransaction(component databaseservicetransactioninterfaces.TransactionComponent, transactionTitle *string) (databaseservicetransactioninterfaces.DatabaseServiceTransaction, error) {
	mongoTransactionComponent, _ := NewMongoTransactionComponent(&component)

	mongoTransaction, _ := mongotransactions.NewMongoTransaction(
		&mongoTransactionComponent,
		mongodbclients.DatabaseClient,
		transactionTitle,
	)

	return &databaseServiceTransaction{
		MongoTransaction: &mongoTransaction,
	}, nil
}

func (dbSvcTransact *databaseServiceTransaction) RunTransaction(input interface{}) (interface{}, error) {
	result, err := (*dbSvcTransact.MongoTransaction).RunTransaction(input)
	if err != nil {
		return nil, horeekaaexceptiontofailure.ConvertException(
			"/services/runTransaction",
			&err,
		)
	}
	return result, nil
}
