package mongodbcoretransactiontests

import (
	"testing"

	"github.com/stretchr/testify/suite"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretransactions "github.com/horeekaa/backend/core/databaseClient/mongodb/transactions"
	mongodbcoreclientmocks "github.com/horeekaa/backend/tests/mocks/core/databaseClient/mongodb/interfaces/init"
	mongodbcoretransactionmocks "github.com/horeekaa/backend/tests/mocks/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcorewrappermocks "github.com/horeekaa/backend/tests/mocks/core/databaseClient/mongodb/interfaces/wrappers"
)

type MongodbTransactionTestSuite struct {
	suite.Suite
	mockedMongoClient          *mongodbcoreclientmocks.MongoClient
	mockedTransactionComponent *mongodbcoretransactionmocks.TransactionComponent
	mockedSession              *mongodbcorewrappermocks.MongoSession

	mongodbTransactionUnderTest mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func (mongodbtrxsuite *MongodbTransactionTestSuite) SetupTest() {
	mongodbtrxsuite.mockedMongoClient = &mongodbcoreclientmocks.MongoClient{}
	mongodbtrxsuite.mockedTransactionComponent = &mongodbcoretransactionmocks.TransactionComponent{}
	mongodbtrxsuite.mockedSession = &mongodbcorewrappermocks.MongoSession{}

	mongoRepoTransaction, _ := mongodbcoretransactions.NewMongoTransaction(
		mongodbtrxsuite.mockedMongoClient,
	)
	mongoRepoTransaction.SetTransaction(
		mongodbtrxsuite.mockedTransactionComponent,
		"MyTransaction",
	)
	mongodbtrxsuite.mongodbTransactionUnderTest = mongoRepoTransaction
}

func TestMongoTransactionSuite(t *testing.T) {
	suite.Run(t, new(MongodbTransactionTestSuite))
}
