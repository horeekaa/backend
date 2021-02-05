package databaseservicetransactions

import (
	"strconv"

	configs "github.com/horeekaa/backend/_commons/configs"
	horeekaaexceptiontofailure "github.com/horeekaa/backend/_errors/serviceFailures/_exceptionToFailure"
	mongodbclients "github.com/horeekaa/backend/repositories/databaseClient/mongoDB"
	mongotransactioninterfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/transaction"
	mongotransactions "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/transactions"
	databaseservicetransactioninterfaces "github.com/horeekaa/backend/services/database/interfaces/transaction"
)

type databaseServiceTransaction struct {
	MongoTransaction *mongotransactioninterfaces.MongoRepoTransaction
}

func NewDatabaseServiceTransaction(component *databaseservicetransactioninterfaces.TransactionComponent, transactionTitle *string) (databaseservicetransactioninterfaces.DatabaseServiceTransaction, error) {
	mongoTransactionComponent, _ := NewMongoTransactionComponent(component)

	timeout, err := strconv.Atoi(configs.GetEnvVariable(configs.DbConfigTimeout))
	repository, err := mongodbclients.NewMongoClientRef(
		configs.GetEnvVariable(configs.DbConfigURL),
		configs.GetEnvVariable(configs.DbConfigDBName),
		timeout,
	)
	if err != nil {
		return nil, err
	}

	mongoTransaction, _ := mongotransactions.NewMongoTransaction(
		&mongoTransactionComponent,
		repository,
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
