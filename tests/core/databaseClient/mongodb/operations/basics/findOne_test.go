package mongodbcorebasicoperationtests

import (
	"errors"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"

	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	"github.com/horeekaa/backend/model"
	mongodbbasicoperationfixtures "github.com/horeekaa/backend/tests/fixtures/core/databaseClient/mongodb/operations"
	mongodbcoreoperationwrappermocks "github.com/horeekaa/backend/tests/mocks/core/databaseClient/mongodb/interfaces/wrappers"
)

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindOneWithSessionOK() {
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

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("FlattenMap", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.FindOne(
		map[string]interface{}{},
		&account,
		&mongodbcoretypes.OperationOptions{
			Session: mockedMongoSession,
		},
	)

	mockedSingleResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), true, out)
	assert.Nil(mongodbBscOpSuite.T(), err)
	assert.Equal(mongodbBscOpSuite.T(), mongodbbasicoperationfixtures.BasicOpsSingleResultOutput, account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindOneWithSessionNoDocumentError() {
	mockedSingleResult := &mongodbcoreoperationwrappermocks.MongoSingleResult{}
	mockedMongoSession := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mockedSingleResult.
		On("Decode", mock.Anything).
		Return(mongo.ErrNoDocuments).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("FindOne", mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"), mock.AnythingOfType("primitive.M")).
		Return(mockedSingleResult).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("FlattenMap", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.FindOne(
		map[string]interface{}{},
		&account,
		&mongodbcoretypes.OperationOptions{
			Session: mockedMongoSession,
		},
	)

	mockedSingleResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), true, out)
	assert.Equal(mongodbBscOpSuite.T(),
		mongo.ErrNoDocuments, err)
	assert.Zero(mongodbBscOpSuite.T(), account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindOneWithSessionUpstreamError() {
	mockedSingleResult := &mongodbcoreoperationwrappermocks.MongoSingleResult{}
	mockedMongoSession := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mockedSingleResult.
		On("Decode", mock.Anything).
		Return(errors.New("Some Upstream Error")).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("FindOne", mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"), mock.AnythingOfType("primitive.M")).
		Return(mockedSingleResult).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("FlattenMap", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.FindOne(
		map[string]interface{}{},
		&account,
		&mongodbcoretypes.OperationOptions{
			Session: mockedMongoSession,
		},
	)

	mockedSingleResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), false, out)
	assert.Equal(mongodbBscOpSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fmt.Sprintf("/%s/findOne", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
			nil,
		), err)
	assert.Zero(mongodbBscOpSuite.T(), account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindOneWithoutSessionOK() {
	mockedSingleResult := &mongodbcoreoperationwrappermocks.MongoSingleResult{}

	mockedSingleResult.
		On("Decode", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*model.Account)
			*arg = mongodbbasicoperationfixtures.BasicOpsSingleResultOutput
		}).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("FindOne", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("primitive.M")).
		Return(mockedSingleResult).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("FlattenMap", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.FindOne(
		map[string]interface{}{},
		&account,
		&mongodbcoretypes.OperationOptions{},
	)

	mockedSingleResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), true, out)
	assert.Nil(mongodbBscOpSuite.T(), err)
	assert.Equal(mongodbBscOpSuite.T(), mongodbbasicoperationfixtures.BasicOpsSingleResultOutput, account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindOneWithoutSessionNoDocumentError() {
	mockedSingleResult := &mongodbcoreoperationwrappermocks.MongoSingleResult{}

	mockedSingleResult.
		On("Decode", mock.Anything).
		Return(mongo.ErrNoDocuments).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("FindOne", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("primitive.M")).
		Return(mockedSingleResult).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("FlattenMap", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.FindOne(
		map[string]interface{}{},
		&account,
		&mongodbcoretypes.OperationOptions{},
	)

	mockedSingleResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), true, out)
	assert.Equal(mongodbBscOpSuite.T(),
		mongo.ErrNoDocuments, err)
	assert.Zero(mongodbBscOpSuite.T(), account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindOneWithoutSessionUpstreamError() {
	mockedSingleResult := &mongodbcoreoperationwrappermocks.MongoSingleResult{}

	mockedSingleResult.
		On("Decode", mock.Anything).
		Return(errors.New("Some Upstream Error")).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("FindOne", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("primitive.M")).
		Return(mockedSingleResult).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("FlattenMap", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.FindOne(
		map[string]interface{}{},
		&account,
		&mongodbcoretypes.OperationOptions{},
	)

	mockedSingleResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), false, out)
	assert.Equal(mongodbBscOpSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fmt.Sprintf("/%s/findOne", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
			nil,
		), err)
	assert.Zero(mongodbBscOpSuite.T(), account)
}
