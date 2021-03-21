package mongodbcoretransactions

import (
	"context"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	horeekaacoreexception "github.com/horeekaa/backend/core/_errors/repoExceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/_errors/repoExceptions/_enums"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	mongodbcoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/init"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepoTransaction struct {
	Component        *mongodbcoretransactioninterfaces.TransactionComponent
	Repository       *mongodbcoreclientinterfaces.MongoClient
	RetryCounter     int
	TransactionTitle *string
}

func NewMongoTransaction(component *mongodbcoretransactioninterfaces.TransactionComponent, repository *mongodbcoreclientinterfaces.MongoClient, transactionTitle *string) (mongodbcoretransactioninterfaces.MongoRepoTransaction, error) {
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

	client, err := (*mongoTrx.Repository).GetClient()
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.DBConnectionFailed,
			"/mongoTransaction/createSession",
			err,
		)
	}

	session, err := client.StartSession()
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.DBConnectionFailed,
			"/mongoTransaction/createSession",
			err,
		)
	}
	if err = session.StartTransaction(); err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.DBConnectionFailed,
			"/mongoTransaction/startTransaction",
			err,
		)
	}

	transactResult := make(chan interface{})
	timeout, err := (*mongoTrx.Repository).GetDatabaseTimeout()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout*time.Second))
	defer cancel()
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		result, err := (*mongoTrx.Component).TransactionBody(&mongodbcoretypes.OperationOptions{
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

			return horeekaacoreexception.NewExceptionObject(
				horeekaacoreexceptionenums.DBConnectionFailed,
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
