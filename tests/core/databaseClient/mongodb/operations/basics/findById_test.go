package mongodbcoreoperationtests

import (
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	mongodbbasicoperationfixtures "github.com/horeekaa/backend/tests/fixtures/core/databaseClient/mongodb/operations"
	mongodbcoreoperationwrappermocks "github.com/horeekaa/backend/tests/mocks/core/databaseClient/mongodb/interfaces/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindByIdOK() {
	mockedSingleResult := &mongodbcoreoperationwrappermocks.MongoSingleResult{}
	mockedMongoSession := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mockedSingleResult.
		On("Decode", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*model.Account)
			*arg = mongodbbasicoperationfixtures.BasicOpsSingleResultOutput
		}).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("FindOne", mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"), mock.AnythingOfType("primitive.M")).
		Return(mockedSingleResult).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.FindByID(
		mongodbbasicoperationfixtures.BasicOpsSingleResultOutput.ID,
		&account,
		&mongodbcoretypes.OperationOptions{
			Session: mockedMongoSession,
		},
	)

	mockedSingleResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())
	assert.Equal(mongodbBscOpSuite.T(), true, out)
	assert.Nil(mongodbBscOpSuite.T(), err)
	assert.Equal(mongodbBscOpSuite.T(), mongodbbasicoperationfixtures.BasicOpsSingleResultOutput, account)
}
