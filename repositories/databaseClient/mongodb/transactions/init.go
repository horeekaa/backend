package mongotransaction

import (
	"context"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	horeekaaexception "github.com/horeekaa/backend/_errors/repoExceptions"
	horeekaaexceptionenums "github.com/horeekaa/backend/_errors/repoExceptions/_enums"
	databaseclient "github.com/horeekaa/backend/repositories/databaseClient/mongoDB"
	mongotransactioninterfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/transaction"
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepoTransaction struct {
	Component        *mongotransactioninterfaces.TransactionComponent
	Repository       *databaseclient.MongoRepository
	RetryCounter     int
	TransactionTitle *string
}

func NewMongoTransaction(component *mongotransactioninterfaces.TransactionComponent, repository *databaseclient.MongoRepository, transactionTitle *string) (mongotransactioninterfaces.MongoRepoTransaction, error) {
	return &mongoRepoTransaction{
		Component:        component,
		Repository:       repository,
		RetryCounter:     0,
		TransactionTitle: transactionTitle,
	}, nil
}

func (mongoTrx *mongoRepoTransaction) RunTransaction(input interface{}) (interface{}, error) {
	if mongoTrx.TransactionTitle == nil {
		*mongoTrx.TransactionTitle = strconv.Itoa(
			int(math.Floor(rand.Float64()*900000+100000) + 1),
		)
	}

	preTransactOutput, err := (*mongoTrx.Component).PreTransaction(input)
	if err != nil {
		return nil, err
	}

	session, err := (*mongoTrx).Repository.Client.StartSession()
	if err != nil {
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.DBConnectionFailed,
			"/mongoTransaction/createSession",
			err,
		)
	}
	if err = session.StartTransaction(); err != nil {
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.DBConnectionFailed,
			"/mongoTransaction/startTransaction",
			err,
		)
	}

	transactResult := make(chan interface{})
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration((*mongoTrx).Repository.Timeout)*time.Second)
	defer cancel()
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		result, err := (*mongoTrx.Component).TransactionBody(&mongooperations.OperationOptions{
			Session: &sc,
		}, preTransactOutput)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		if err = session.CommitTransaction(sc); err != nil {
			session.AbortTransaction(sc)

			(*mongoTrx).RetryCounter = (*mongoTrx).RetryCounter + 1
			if (*mongoTrx).RetryCounter < 10 {
				log.Printf("Retrying Transaction %s in 50ms", *mongoTrx.TransactionTitle)
				time.Sleep(50 * time.Millisecond)
				result, err = (*mongoTrx).RunTransaction(input)
				if err != nil {
					return err
				}
				transactResult <- result
				return nil
			}

			return horeekaaexception.NewExceptionObject(
				horeekaaexceptionenums.DBConnectionFailed,
				"/mongoTransaction/commitTransaction",
				err,
			)
		}

		log.Printf("Transaction %s successfully run", *mongoTrx.TransactionTitle)
		transactResult <- result
		return nil
	}); err != nil {
		return nil, err
	}
	session.EndSession(ctx)

	return <-transactResult, nil
}
