package mongodbcoretransactions

import (
	"context"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	mongodbcoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/init"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
)

type mongoRepoTransaction struct {
	component        mongodbcoretransactioninterfaces.TransactionComponent
	mongoClient      mongodbcoreclientinterfaces.MongoClient
	transactionTitle string
}

func NewMongoTransaction(mongoClient mongodbcoreclientinterfaces.MongoClient) (mongodbcoretransactioninterfaces.MongoRepoTransaction, error) {
	return &mongoRepoTransaction{
		mongoClient: mongoClient,
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

	timeout, err := mongoTrx.mongoClient.GetDatabaseTimeout()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout*time.Second))
	defer cancel()
	defer session.EndSession(ctx)

	// session WithTransaction automatically start and commit or abort the session
	result, err := session.WithTransaction(ctx, mongoTrx.TransactionFn(preTransactOutput))
	if err != nil {
		return nil, err
	}
	log.Printf("Transaction %s successfully run", mongoTrx.transactionTitle)

	return result, nil
}

func (mongoTrx *mongoRepoTransaction) TransactionFn(
	preTransactOutput interface{},
) func(sessCtx mongo.SessionContext) (interface{}, error) {
	return func(sessCtx mongo.SessionContext) (interface{}, error) {
		result, err := mongoTrx.component.TransactionBody(&mongodbcoretypes.OperationOptions{
			Session: sessCtx,
		}, preTransactOutput)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
}
