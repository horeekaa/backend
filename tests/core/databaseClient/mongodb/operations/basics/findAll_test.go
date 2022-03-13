package mongodbcorebasicoperationtests

import (
	"errors"
	"fmt"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"

	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	model "github.com/horeekaa/backend/model"
	mongodbbasicoperationfixtures "github.com/horeekaa/backend/tests/fixtures/core/databaseClient/mongodb/operations"
	mongodbcoreoperationwrappermocks "github.com/horeekaa/backend/tests/mocks/core/databaseClient/mongodb/interfaces/wrappers"
)

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindAllWithSessionOK() {
	mockedFindCursor := &mongodbcoreoperationwrappermocks.MongoCursor{}
	mockedMongoSession := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mockedFindCursor.
		On("Next", mock.AnythingOfType("*context.timerCtx")).
		Return(true).
		Times(3)
	mockedFindCursor.
		On("Next", mock.AnythingOfType("*context.timerCtx")).
		Return(false).
		Once()
	mockedFindCursor.
		On("Decode", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*model.Account)
			*arg = mongodbbasicoperationfixtures.BasicOpsSingleResultOutput
		}).
		Times(3)

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("Find",
			mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"),
			mock.AnythingOfType("primitive.M"),
			mock.AnythingOfType("*options.FindOptions"),
		).
		Return(mockedFindCursor, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("FlattenMap", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	var accounts []*model.Account
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var account model.Account
		if err := cursor.Decode(&account); err != nil {
			return err
		}
		accounts = append(accounts, &account)
		return nil
	}
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Find(
		map[string]interface{}{},
		&mongodbcoretypes.PaginationOptions{
			QueryLimit: func(i int) *int { return &i }(5),
		},
		appendingFn,
		&mongodbcoretypes.OperationOptions{
			Session: mockedMongoSession,
		},
	)

	mockedFindCursor.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	accountFixtures := []*model.Account{
		&mongodbbasicoperationfixtures.BasicOpsSingleResultOutput,
		&mongodbbasicoperationfixtures.BasicOpsSingleResultOutput,
		&mongodbbasicoperationfixtures.BasicOpsSingleResultOutput,
	}
	assert.Equal(mongodbBscOpSuite.T(), true, out)
	assert.Nil(mongodbBscOpSuite.T(), err)
	assert.Equal(mongodbBscOpSuite.T(), accountFixtures, accounts)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindAllWithSessionDecodingError() {
	mockedFindCursor := &mongodbcoreoperationwrappermocks.MongoCursor{}
	mockedMongoSession := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mockedFindCursor.
		On("Next", mock.AnythingOfType("*context.timerCtx")).
		Return(true).
		Times(2)
	mockedFindCursor.
		On("Decode", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*model.Account)
			*arg = mongodbbasicoperationfixtures.BasicOpsSingleResultOutput
		}).
		Once()
	mockedFindCursor.
		On("Decode", mock.Anything).
		Return(errors.New("Some Upstream Error")).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("Find",
			mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"),
			mock.AnythingOfType("primitive.M"),
			mock.AnythingOfType("*options.FindOptions"),
		).
		Return(mockedFindCursor, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("FlattenMap", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	var accounts []*model.Account
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var account model.Account
		if err := cursor.Decode(&account); err != nil {
			return err
		}
		accounts = append(accounts, &account)
		return nil
	}
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Find(
		map[string]interface{}{},
		&mongodbcoretypes.PaginationOptions{
			QueryLimit: func(i int) *int { return &i }(5),
		},
		appendingFn,
		&mongodbcoretypes.OperationOptions{
			Session: mockedMongoSession,
		},
	)

	mockedFindCursor.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	accountFixtures := []*model.Account{
		&mongodbbasicoperationfixtures.BasicOpsSingleResultOutput,
	}
	assert.Equal(mongodbBscOpSuite.T(), false, out)
	assert.Equal(mongodbBscOpSuite.T(), horeekaacoreexception.NewExceptionObject(
		horeekaacoreexceptionenums.QueryObjectFailed,
		fmt.Sprintf("%s.Find", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
		errors.New("Some Upstream Error"),
	), err)
	assert.Equal(mongodbBscOpSuite.T(), accountFixtures, accounts)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindAllWithSessionUpstreamError() {
	mockedMongoSession := &struct {
		mongodbcorewrapperinterfaces.MongoSessionContext
	}{}

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("Find",
			mock.AnythingOfType("*struct { mongodbcorewrapperinterfaces.MongoSessionContext }"),
			mock.AnythingOfType("primitive.M"),
			mock.AnythingOfType("*options.FindOptions"),
		).
		Return(nil, errors.New("Some Upstream Error")).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("FlattenMap", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	var accounts []*model.Account
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var account model.Account
		if err := cursor.Decode(&account); err != nil {
			return err
		}
		accounts = append(accounts, &account)
		return nil
	}
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Find(
		map[string]interface{}{},
		&mongodbcoretypes.PaginationOptions{
			QueryLimit: func(i int) *int { return &i }(5),
		},
		appendingFn,
		&mongodbcoretypes.OperationOptions{
			Session: mockedMongoSession,
		},
	)

	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	var accountFixtures []*model.Account
	assert.Equal(mongodbBscOpSuite.T(), false, out)
	assert.Equal(mongodbBscOpSuite.T(), horeekaacoreexception.NewExceptionObject(
		horeekaacoreexceptionenums.QueryObjectFailed,
		fmt.Sprintf("%s.Find", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
		errors.New("Some Upstream Error"),
	), err)
	assert.Equal(mongodbBscOpSuite.T(), accountFixtures, accounts)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindAllWithoutSessionOK() {
	mockedFindCursor := &mongodbcoreoperationwrappermocks.MongoCursor{}

	mockedFindCursor.
		On("Next", mock.AnythingOfType("*context.timerCtx")).
		Return(true).
		Times(3)
	mockedFindCursor.
		On("Next", mock.AnythingOfType("*context.timerCtx")).
		Return(false).
		Once()
	mockedFindCursor.
		On("Decode", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*model.Account)
			*arg = mongodbbasicoperationfixtures.BasicOpsSingleResultOutput
		}).
		Times(3)

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("Find",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("primitive.M"),
			mock.AnythingOfType("*options.FindOptions"),
		).
		Return(mockedFindCursor, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("FlattenMap", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	var accounts []*model.Account
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var account model.Account
		if err := cursor.Decode(&account); err != nil {
			return err
		}
		accounts = append(accounts, &account)
		return nil
	}
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Find(
		map[string]interface{}{},
		&mongodbcoretypes.PaginationOptions{
			QueryLimit: func(i int) *int { return &i }(5),
		},
		appendingFn,
		&mongodbcoretypes.OperationOptions{},
	)

	mockedFindCursor.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	accountFixtures := []*model.Account{
		&mongodbbasicoperationfixtures.BasicOpsSingleResultOutput,
		&mongodbbasicoperationfixtures.BasicOpsSingleResultOutput,
		&mongodbbasicoperationfixtures.BasicOpsSingleResultOutput,
	}
	assert.Equal(mongodbBscOpSuite.T(), true, out)
	assert.Nil(mongodbBscOpSuite.T(), err)
	assert.Equal(mongodbBscOpSuite.T(), accountFixtures, accounts)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindAllWithoutSessionDecodingError() {
	mockedFindCursor := &mongodbcoreoperationwrappermocks.MongoCursor{}

	mockedFindCursor.
		On("Next", mock.AnythingOfType("*context.timerCtx")).
		Return(true).
		Times(2)
	mockedFindCursor.
		On("Decode", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*model.Account)
			*arg = mongodbbasicoperationfixtures.BasicOpsSingleResultOutput
		}).
		Once()
	mockedFindCursor.
		On("Decode", mock.Anything).
		Return(errors.New("Some Upstream Error")).
		Once()

	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("Find",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("primitive.M"),
			mock.AnythingOfType("*options.FindOptions"),
		).
		Return(mockedFindCursor, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("FlattenMap", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	var accounts []*model.Account
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var account model.Account
		if err := cursor.Decode(&account); err != nil {
			return err
		}
		accounts = append(accounts, &account)
		return nil
	}
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Find(
		map[string]interface{}{},
		&mongodbcoretypes.PaginationOptions{
			QueryLimit: func(i int) *int { return &i }(5),
		},
		appendingFn,
		&mongodbcoretypes.OperationOptions{},
	)

	mockedFindCursor.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	accountFixtures := []*model.Account{
		&mongodbbasicoperationfixtures.BasicOpsSingleResultOutput,
	}
	assert.Equal(mongodbBscOpSuite.T(), false, out)
	assert.Equal(mongodbBscOpSuite.T(), horeekaacoreexception.NewExceptionObject(
		horeekaacoreexceptionenums.QueryObjectFailed,
		fmt.Sprintf("%s.Find", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
		errors.New("Some Upstream Error"),
	), err)
	assert.Equal(mongodbBscOpSuite.T(), accountFixtures, accounts)
}

func (mongodbBscOpSuite *MongodbBasicOperationTestSuite) TestFindAllWithoutSessionUpstreamError() {
	mongodbBscOpSuite.mockedMongoCollectionRef.
		On("Find",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("primitive.M"),
			mock.AnythingOfType("*options.FindOptions"),
		).
		Return(nil, errors.New("Some Upstream Error")).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("RemoveNil", mock.Anything).
		Return(true, nil).
		Once()

	mongodbBscOpSuite.mockedMapProcessorUtility.
		On("FlattenMap", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	var accounts []*model.Account
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var account model.Account
		if err := cursor.Decode(&account); err != nil {
			return err
		}
		accounts = append(accounts, &account)
		return nil
	}
	out, err := mongodbBscOpSuite.basicOperationUnderTest.Find(
		map[string]interface{}{},
		&mongodbcoretypes.PaginationOptions{
			QueryLimit: func(i int) *int { return &i }(5),
		},
		appendingFn,
		&mongodbcoretypes.OperationOptions{},
	)

	mongodbBscOpSuite.mockedMapProcessorUtility.AssertExpectations(mongodbBscOpSuite.T())
	mongodbBscOpSuite.mockedMongoCollectionRef.AssertExpectations(mongodbBscOpSuite.T())

	var accountFixtures []*model.Account
	assert.Equal(mongodbBscOpSuite.T(), false, out)
	assert.Equal(mongodbBscOpSuite.T(), horeekaacoreexception.NewExceptionObject(
		horeekaacoreexceptionenums.QueryObjectFailed,
		fmt.Sprintf("%s.Find", mongodbBscOpSuite.basicOperationUnderTest.GetCollectionName()),
		errors.New("Some Upstream Error"),
	), err)
	assert.Equal(mongodbBscOpSuite.T(), accountFixtures, accounts)
}
