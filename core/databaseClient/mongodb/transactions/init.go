package mongodbcoretransactions

import (
	"context"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	mongodbcoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/init"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
)

type mongoRepoTransaction struct {
	component        mongodbcoretransactioninterfaces.TransactionComponent
	mongoClient      mongodbcoreclientinterfaces.MongoClient
	retryCounter     int
	transactionTitle string
}

func NewMongoTransaction(mongoClient mongodbcoreclientinterfaces.MongoClient) (mongodbcoretransactioninterfaces.MongoRepoTransaction, error) {
	return &mongoRepoTransaction{
		mongoClient:  mongoClient,
		retryCounter: 0,
	}, nil
}

func (mongoTrx *mongoRepoTransaction) SetTransaction(component mongodbcoretransactioninterfaces.TransactionComponent, transactionTitle string) bool {
	mongoTrx.component = component
	mongoTrx.transactionTitle = transactionTitle
	return true
}

func (mongoTrx *mongoRepoTransaction) RunTransaction(input interface{}) (interface{}, error) {
	if &mongoTrx.transactionTitle == nil {
		mongoTrx.transactionTitle = strconv.Itoa(
			int(math.Floor(rand.Float64()*900000+100000) + 1),
		)
	}

	preTransactOutput, err := mongoTrx.component.PreTransaction(input)
	if err != nil {
		return nil, err
	}

	session, err := mongoTrx.mongoClient.CreateNewSession()
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

	timeout, err := mongoTrx.mongoClient.GetDatabaseTimeout()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout*time.Second))
	defer cancel()
	defer session.EndSession(ctx)
	sc, _ := mongoTrx.mongoClient.CreateNewSessionContext(ctx, session)
	result, err := mongoTrx.component.TransactionBody(&mongodbcoretypes.OperationOptions{
		Session: sc,
	}, preTransactOutput)
	if err != nil {
		session.AbortTransaction(sc)
		return nil, err
	}

	if err = session.CommitTransaction(sc); err != nil {
		session.AbortTransaction(sc)

		mongoTrx.retryCounter = mongoTrx.retryCounter + 1
		if mongoTrx.retryCounter < 10 {
			log.Printf("Retrying Transaction %s in 50ms", mongoTrx.transactionTitle)
			time.Sleep(50 * time.Millisecond)
			result, err = (*mongoTrx).RunTransaction(input)
			if err != nil {
				return nil, err
			}
			return result, nil
		}

		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.DBConnectionFailed,
			"/mongoTransaction/commitTransaction",
			err,
		)
	}
	log.Printf("Transaction %s successfully run", mongoTrx.transactionTitle)

	return result, nil
}
