package mongodbcorebasicoperationtests

import (
	"errors"
	"fmt"

	gotest_assert "gotest.tools/assert"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	"github.com/horeekaa/backend/model"
	mongodbbasicoperationfixtures "github.com/horeekaa/backend/tests/fixtures/core/databaseClient/mongodb/operations"
	mongodbcoreoperationwrappermocks "github.com/horeekaa/backend/tests/mocks/core/databaseClient/mongodb/interfaces/wrappers"
)

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestCreateWithSessionOK() {
	mockedInsertOneResult := &mongodbcoreoperationwrappermocks.MongoInsertOneResult{}
	mockedMongoSession := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mockedInsertOneResult.On("GetInsertedID").
		Return(mongodbbasicoperationfixtures.BasicOpsSingleResultOutput.ID).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("InsertOne",
			mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"),
			mock.AnythingOfType("primitive.M"),
		).
		Return(mockedInsertOneResult, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Create(
		mongodbbasicoperationfixtures.BasicOpsSingleResultOutput,
		&account,
		&mongodbcoretypes.OperationOptions{
			Session: mockedMongoSession,
		},
	)

	mockedInsertOneResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), true, out)
	assert.Nil(mongodbBscOpSuite.T(), err)

	gotest_assert.DeepEqual(mongodbBscOpSuite.T(), mongodbbasicoperationfixtures.BasicOpsSingleResultOutput, account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestCreateWithSessionUpstreamError() {
	mockedMongoSession := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("InsertOne",
			mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"),
			mock.AnythingOfType("primitive.M"),
		).
		Return(nil, errors.New("Some Upstream Error")).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Create(
		mongodbbasicoperationfixtures.BasicOpsSingleResultOutput,
		&account,
		&mongodbcoretypes.OperationOptions{
			Session: mockedMongoSession,
		},
	)

	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), false, out)
	assert.Equal(mongodbBscOpSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.CreateObjectFailed,
			fmt.Sprintf("/%s/create", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
			errors.New("Some Upstream Error"),
		), err,
	)

	assert.Zero(mongodbBscOpSuite.T(), account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestCreateWithoutSessionOK() {
	mockedInsertOneResult := &mongodbcoreoperationwrappermocks.MongoInsertOneResult{}

	mockedInsertOneResult.On("GetInsertedID").
		Return(mongodbbasicoperationfixtures.BasicOpsSingleResultOutput.ID).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("InsertOne",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("primitive.M"),
		).
		Return(mockedInsertOneResult, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Create(
		mongodbbasicoperationfixtures.BasicOpsSingleResultOutput,
		&account,
		&mongodbcoretypes.OperationOptions{},
	)

	mockedInsertOneResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), true, out)
	assert.Nil(mongodbBscOpSuite.T(), err)

	gotest_assert.DeepEqual(mongodbBscOpSuite.T(), mongodbbasicoperationfixtures.BasicOpsSingleResultOutput, account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestCreateWithoutSessionUpstreamError() {
	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("InsertOne",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("primitive.M"),
		).
		Return(nil, errors.New("Some Upstream Error")).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Create(
		mongodbbasicoperationfixtures.BasicOpsSingleResultOutput,
		&account,
		&mongodbcoretypes.OperationOptions{},
	)

	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), false, out)
	assert.Equal(mongodbBscOpSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.CreateObjectFailed,
			fmt.Sprintf("/%s/create", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
			errors.New("Some Upstream Error"),
		), err,
	)

	assert.Zero(mongodbBscOpSuite.T(), account)
}
