package mongodbcorebasicoperationtests

import (
	"errors"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	gotest_assert "gotest.tools/assert"

	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	"github.com/horeekaa/backend/model"
	mongodbbasicoperationfixtures "github.com/horeekaa/backend/tests/fixtures/core/databaseClient/mongodb/operations"
	mongodbcoreoperationwrappermocks "github.com/horeekaa/backend/tests/mocks/core/databaseClient/mongodb/interfaces/wrappers"
)

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestUpdateWithSessionOK() {
	mockedSingleResult := &mongodbcoreoperationwrappermocks.MongoSingleResult{}
	mockedMongoSession := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mockedSingleResult.
		On("Decode", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*model.Account)
			*arg = mongodbbasicoperationfixtures.BasicOpsSingleResultOutput
		}).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.On(
		"FindOne",
		mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"),
		mock.AnythingOfType("primitive.M"),
	).
		Return(mockedSingleResult).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.On(
		"UpdateOne",
		mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"),
		mock.AnythingOfType("primitive.M"),
		mock.AnythingOfType("primitive.M"),
	).
		Return(new(interface{}), nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Update(
		map[string]interface{}{},
		map[string]interface{}{},
		&account,
		&mongodbcoretypes.OperationOptions{
			Session: mockedMongoSession,
		},
	)

	mockedSingleResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), true, out)
	assert.Nil(mongodbBscOpSuite.T(), err)

	gotest_assert.DeepEqual(mongodbBscOpSuite.T(), mongodbbasicoperationfixtures.BasicOpsSingleResultOutput, account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestUpdateWithSessionFindByIDError() {
	mockedSingleResult := &mongodbcoreoperationwrappermocks.MongoSingleResult{}
	mockedMongoSession := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mockedSingleResult.
		On("Decode", mock.Anything).
		Return(mongo.ErrNoDocuments).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.On(
		"FindOne",
		mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"),
		mock.AnythingOfType("primitive.M"),
	).
		Return(mockedSingleResult).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.On(
		"UpdateOne",
		mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"),
		mock.AnythingOfType("primitive.M"),
		mock.AnythingOfType("primitive.M"),
	).
		Return(new(interface{}), nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Update(
		map[string]interface{}{},
		map[string]interface{}{},
		&account,
		&mongodbcoretypes.OperationOptions{
			Session: mockedMongoSession,
		},
	)

	mockedSingleResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), false, out)
	gotest_assert.DeepEqual(mongodbBscOpSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.IDNotFound,
			fmt.Sprintf("/%s/update", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
			horeekaacoreexception.NewExceptionObject(
				horeekaacoreexceptionenums.IDNotFound,
				fmt.Sprintf("/%s/findByID", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
				nil,
			),
		), err,
	)

	assert.Zero(mongodbBscOpSuite.T(), account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestUpdateWithSessionUpstreamError() {
	mockedMongoSession := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.On(
		"UpdateOne",
		mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"),
		mock.AnythingOfType("primitive.M"),
		mock.AnythingOfType("primitive.M"),
	).
		Return(nil, errors.New("Some Upstream Error")).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Update(
		map[string]interface{}{},
		map[string]interface{}{},
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
			horeekaacoreexceptionenums.UpdateObjectFailed,
			fmt.Sprintf("/%s/update", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
			errors.New("Some Upstream Error"),
		), err,
	)

	assert.Zero(mongodbBscOpSuite.T(), account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestUpdateWithoutSessionOK() {
	mockedSingleResult := &mongodbcoreoperationwrappermocks.MongoSingleResult{}

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mockedSingleResult.
		On("Decode", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*model.Account)
			*arg = mongodbbasicoperationfixtures.BasicOpsSingleResultOutput
		}).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.On(
		"FindOne",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("primitive.M"),
	).
		Return(mockedSingleResult).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.On(
		"UpdateOne",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("primitive.M"),
		mock.AnythingOfType("primitive.M"),
	).
		Return(new(interface{}), nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Update(
		map[string]interface{}{},
		map[string]interface{}{},
		&account,
		&mongodbcoretypes.OperationOptions{},
	)

	mockedSingleResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), true, out)
	assert.Nil(mongodbBscOpSuite.T(), err)

	gotest_assert.DeepEqual(mongodbBscOpSuite.T(), mongodbbasicoperationfixtures.BasicOpsSingleResultOutput, account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestUpdateWithoutSessionFindByIDError() {
	mockedSingleResult := &mongodbcoreoperationwrappermocks.MongoSingleResult{}

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mockedSingleResult.
		On("Decode", mock.Anything).
		Return(mongo.ErrNoDocuments).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.On(
		"FindOne",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("primitive.M"),
	).
		Return(mockedSingleResult).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.On(
		"UpdateOne",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("primitive.M"),
		mock.AnythingOfType("primitive.M"),
	).
		Return(new(interface{}), nil).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Update(
		map[string]interface{}{},
		map[string]interface{}{},
		&account,
		&mongodbcoretypes.OperationOptions{},
	)

	mockedSingleResult.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), false, out)
	gotest_assert.DeepEqual(mongodbBscOpSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.IDNotFound,
			fmt.Sprintf("/%s/update", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
			horeekaacoreexception.NewExceptionObject(
				horeekaacoreexceptionenums.IDNotFound,
				fmt.Sprintf("/%s/findByID", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
				nil,
			),
		), err,
	)

	assert.Zero(mongodbBscOpSuite.T(), account)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestUpdateWithoutSessionUpstreamError() {
	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.On(
		"UpdateOne",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("primitive.M"),
		mock.AnythingOfType("primitive.M"),
	).
		Return(nil, errors.New("Some Upstream Error")).
		Once()

	var account model.Account
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Update(
		map[string]interface{}{},
		map[string]interface{}{},
		&account,
		&mongodbcoretypes.OperationOptions{},
	)

	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	assert.Equal(mongodbBscOpSuite.T(), false, out)
	assert.Equal(mongodbBscOpSuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpdateObjectFailed,
			fmt.Sprintf("/%s/update", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
			errors.New("Some Upstream Error"),
		), err,
	)

	assert.Zero(mongodbBscOpSuite.T(), account)
}
