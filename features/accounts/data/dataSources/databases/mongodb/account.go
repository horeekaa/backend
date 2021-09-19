package mongodbaccountdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type accountDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewAccountDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbaccountdatasourceinterfaces.AccountDataSourceMongo, error) {
	basicOperation.SetCollection("accounts")
	return &accountDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (accDataSourceMongo *accountDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	var output model.Account
	_, err := accDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (accDataSourceMongo *accountDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	var output model.Account
	_, err := accDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (accDataSourceMongo *accountDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Account, error) {
	var accounts = []*model.Account{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var account model.Account
		if err := cursor.Decode(&account); err != nil {
			return err
		}
		accounts = append(accounts, &account)
		return nil
	}
	_, err := accDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return accounts, err
}

func (accDataSourceMongo *accountDataSourceMongo) Create(input *model.CreateAccount, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	_, err := accDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Account
	_, err = accDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (accDataSourceMongo *accountDataSourceMongo) Update(updateCriteria map[string]interface{}, updateData *model.UpdateAccount, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	_, err := accDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.Account
	_, err = accDataSourceMongo.basicOperation.Update(
		updateCriteria,
		map[string]interface{}{
			"$set": updateData,
		},
		&output,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (accDataSourceMongo *accountDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.UpdateAccount,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	defaultAccountType := model.AccountTypePerson
	defaultLanguage := model.LanguageID
	currentTime := time.Now()

	existingObject, err := accDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			"/accountDataSource/update",
			nil,
		)
	}

	if &(*existingObject).Type == nil {
		input.Type = &defaultAccountType
	}
	if input.Language == nil {
		input.Language = &defaultLanguage
	}
	input.UpdatedAt = &currentTime

	return true, nil
}

func (accDataSourceMongo *accountDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.CreateAccount,
) (bool, error) {
	defaultAccountStatus := model.AccountStatusActive
	defaultAccountType := model.AccountTypePerson
	currentTime := time.Now()

	if input.Status == nil {
		input.Status = &defaultAccountStatus
	}
	if &input.Type == nil {
		input.Type = defaultAccountType
	}
	if input.DeviceTokens == nil {
		input.DeviceTokens = []string{}
	}
	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime

	return true, nil
}
