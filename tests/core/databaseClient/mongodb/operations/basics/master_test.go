package mongodbcoreoperationtests

import (
	"testing"
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcoreoperations "github.com/horeekaa/backend/core/databaseClient/mongodb/operations"
	mongodbcoreclientmocks "github.com/horeekaa/backend/tests/mocks/core/databaseClient/mongodb/interfaces/init"
	mongodbcoreoperationwrappermocks "github.com/horeekaa/backend/tests/mocks/core/databaseClient/mongodb/interfaces/wrappers"
	coreutilitymocks "github.com/horeekaa/backend/tests/mocks/core/utilities/interfaces"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MongodbBasicOperationTestSuite struct {
	suite.Suite
	mockedMongoClient         *mongodbcoreclientmocks.MongoClient
	mockedMongoCollectionRef  *mongodbcoreoperationwrappermocks.MongoCollectionRef
	mockedMapProcessorUtility *coreutilitymocks.MapProcessorUtility

	basicOperationUnderTest mongodbcoreoperationinterfaces.BasicOperation
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) SetupSuite() {
	mongodbBscOpSuite.mockedMongoClient = &mongodbcoreclientmocks.MongoClient{}
	mongodbBscOpSuite.mockedMapProcessorUtility = &coreutilitymocks.MapProcessorUtility{}
	mongodbBscOpSuite.mockedMongoCollectionRef = &mongodbcoreoperationwrappermocks.MongoCollectionRef{}

	mongodbBscOpSuite.mockedMongoClient.
		On("GetDatabaseTimeout").
		Return(20*time.Second, nil).
		Once()

	mongodbBscOpSuite.mockedMongoClient.
		On("GetCollectionRef", mock.AnythingOfType("string")).
		Return(mongodbBscOpSuite.mockedMongoCollectionRef, nil).
		Once()

	basicOpsUT, _ := mongodbcoreoperations.NewBasicOperation(
		mongodbBscOpSuite.mockedMongoClient,
		mongodbBscOpSuite.mockedMapProcessorUtility,
	)
	basicOpsUT.SetCollection("MyCollection")
	mongodbBscOpSuite.basicOperationUnderTest = basicOpsUT
}

func TestBasicOperationSuite(t *testing.T) {
	suite.Run(t, new(MongodbBasicOperationTestSuite))
}
