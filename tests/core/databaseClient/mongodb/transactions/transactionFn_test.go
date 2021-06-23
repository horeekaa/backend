package mongodbcoretransactiontests

import (
	"errors"

	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbtransactionfixtures "github.com/horeekaa/backend/tests/fixtures/core/databaseClient/mongodb/transactions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (mongodbtrxsuite *MongodbTransactionTestSuite) TestTransactionFnOK() {
	mockedSessionCtx := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mongodbtrxsuite.mockedTransactionComponent.
		On("TransactionBody", mock.Anything, mock.Anything).
		Return(mongodbtransactionfixtures.TransactionOutput, nil).
		Once()

	output, err := mongodbtrxsuite.mongodbTransactionUnderTest.TransactionFn(
		new(interface{}),
	)(mockedSessionCtx)

	mongodbtrxsuite.mockedTransactionComponent.AssertExpectations(mongodbtrxsuite.T())

	assert.Equal(mongodbtrxsuite.T(), mongodbtransactionfixtures.TransactionOutput, output)
	assert.Nil(mongodbtrxsuite.T(), err)
}

func (mongodbtrxsuite *MongodbTransactionTestSuite) TestTransactionFnError() {
	mockedSessionCtx := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mongodbtrxsuite.mockedTransactionComponent.
		On("TransactionBody", mock.Anything, mock.Anything).
		Return(nil, errors.New("Some Upstream Error")).
		Once()

	output, err := mongodbtrxsuite.mongodbTransactionUnderTest.TransactionFn(
		new(interface{}),
	)(mockedSessionCtx)

	mongodbtrxsuite.mockedTransactionComponent.AssertExpectations(mongodbtrxsuite.T())

	assert.Nil(mongodbtrxsuite.T(), output)
	assert.Equal(mongodbtrxsuite.T(), errors.New("Some Upstream Error"), err)
}
