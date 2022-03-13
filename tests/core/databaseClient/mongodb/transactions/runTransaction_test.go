package mongodbcoretransactiontests

import (
	"errors"
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbtransactionfixtures "github.com/horeekaa/backend/tests/fixtures/core/databaseClient/mongodb/transactions"
)

func (mongodbtrxsuite *MongodbTransactionTestSuite) TestRunTransactionOK() {
	mongodbtrxsuite.mockedTransactionComponent.
		On("PreTransaction", mock.Anything).
		Return(new(interface{}), nil).
		Once()

	mongodbtrxsuite.mockedMongoClient.
		On("GetDatabaseTimeout").
		Return(20*time.Second, nil).
		Once()

	mongodbtrxsuite.mockedMongoClient.
		On("CreateNewSession").
		Return(mongodbtrxsuite.mockedSession, nil).
		Once()

	mongodbtrxsuite.mockedSession.
		On("EndSession", mock.AnythingOfType("*context.timerCtx")).
		Return().
		Once()

	mongodbtrxsuite.mockedSession.
		On("WithTransaction", mock.AnythingOfType("*context.timerCtx"), mock.Anything).
		Return(mongodbtransactionfixtures.TransactionOutput, nil).
		Once()

	result, err := mongodbtrxsuite.mongodbTransactionUnderTest.RunTransaction(
		new(interface{}),
	)

	mongodbtrxsuite.mockedTransactionComponent.AssertExpectations(mongodbtrxsuite.T())
	mongodbtrxsuite.mockedSession.AssertExpectations(mongodbtrxsuite.T())
	mongodbtrxsuite.mockedMongoClient.AssertExpectations(mongodbtrxsuite.T())

	assert.Equal(mongodbtrxsuite.T(), mongodbtransactionfixtures.TransactionOutput, result)
	assert.Nil(mongodbtrxsuite.T(), err)
}

func (mongodbtrxsuite *MongodbTransactionTestSuite) TestRunTransactionPreTransactionError() {
	mongodbtrxsuite.mockedTransactionComponent.
		On("PreTransaction", mock.Anything).
		Return(nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.CreateObjectFailed,
			fmt.Sprintf("%s.Create", "MyCollection"),
			errors.New("Some Upstream Error"),
		)).
		Once()

	result, err := mongodbtrxsuite.mongodbTransactionUnderTest.RunTransaction(
		new(interface{}),
	)

	mongodbtrxsuite.mockedTransactionComponent.AssertExpectations(mongodbtrxsuite.T())

	assert.Nil(mongodbtrxsuite.T(), result)
	assert.Equal(mongodbtrxsuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.CreateObjectFailed,
			fmt.Sprintf("%s.Create", "MyCollection"),
			errors.New("Some Upstream Error"),
		), err,
	)
}

func (mongodbtrxsuite *MongodbTransactionTestSuite) TestRunTransactionCreateSessionError() {
	mongodbtrxsuite.mockedTransactionComponent.
		On("PreTransaction", mock.Anything).
		Return(new(interface{}), nil).
		Once()

	mongodbtrxsuite.mockedMongoClient.
		On("CreateNewSession").
		Return(nil, errors.New("Some Upstream Error")).
		Once()

	result, err := mongodbtrxsuite.mongodbTransactionUnderTest.RunTransaction(
		new(interface{}),
	)

	mongodbtrxsuite.mockedMongoClient.AssertExpectations(mongodbtrxsuite.T())
	mongodbtrxsuite.mockedTransactionComponent.AssertExpectations(mongodbtrxsuite.T())

	assert.Nil(mongodbtrxsuite.T(), result)
	assert.Equal(mongodbtrxsuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.DBConnectionFailed,
			"MongoDBTransaction",
			errors.New("Some Upstream Error"),
		), err,
	)
}

func (mongodbtrxsuite *MongodbTransactionTestSuite) TestRunTransactionWithTransactionError() {
	mongodbtrxsuite.mockedTransactionComponent.
		On("PreTransaction", mock.Anything).
		Return(new(interface{}), nil).
		Once()

	mongodbtrxsuite.mockedMongoClient.
		On("GetDatabaseTimeout").
		Return(20*time.Second, nil).
		Once()

	mongodbtrxsuite.mockedMongoClient.
		On("CreateNewSession").
		Return(mongodbtrxsuite.mockedSession, nil).
		Once()

	mongodbtrxsuite.mockedSession.
		On("EndSession", mock.AnythingOfType("*context.timerCtx")).
		Return().
		Once()

	mongodbtrxsuite.mockedSession.
		On("WithTransaction", mock.AnythingOfType("*context.timerCtx"), mock.Anything).
		Return(nil, errors.New("Some Upstream Error")).
		Once()

	result, err := mongodbtrxsuite.mongodbTransactionUnderTest.RunTransaction(
		new(interface{}),
	)

	mongodbtrxsuite.mockedMongoClient.AssertExpectations(mongodbtrxsuite.T())
	mongodbtrxsuite.mockedTransactionComponent.AssertExpectations(mongodbtrxsuite.T())

	assert.Nil(mongodbtrxsuite.T(), result)
	assert.Equal(mongodbtrxsuite.T(),
		horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			"MongoDBTransaction",
			errors.New("Some Upstream Error"),
		), err,
	)
}
