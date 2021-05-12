package mongodbaccountdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
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
	res, err := accDataSourceMongo.basicOperation.FindByID(ID, operationOptions)
	if err != nil {
		return nil, err
	}

	var output model.Account
	res.Decode(&output)
	return &output, nil
}

func (accDataSourceMongo *accountDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	res, err := accDataSourceMongo.basicOperation.FindOne(query, operationOptions)
	if err != nil {
		return nil, err
	}

	var output model.Account
	err = res.Decode(&output)
	if err == mongo.ErrNoDocuments {
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
	cursorDecoder := func(cursor *mongo.Cursor) (interface{}, error) {
		var account *model.Account
		err := cursor.Decode(account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
		return nil, nil
	}

	_, err := accDataSourceMongo.basicOperation.Find(query, paginationOpts, cursorDecoder, operationOptions)
	if err != nil {
		return nil, err
	}

	return accounts, err
}

func (accDataSourceMongo *accountDataSourceMongo) Create(input *model.CreateAccount, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	defaultedInput, err := accDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	output, err := accDataSourceMongo.basicOperation.Create(*defaultedInput.CreateAccount, operationOptions)
	if err != nil {
		return nil, err
	}

	accountOutput := output.Object.(model.Account)
	accountOutput.ID = output.ID

	return &accountOutput, err
}

func (accDataSourceMongo *accountDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.UpdateAccount, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	defaultedInput, err := accDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := accDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateAccount, operationOptions)
	if err != nil {
		return nil, err
	}

	var output model.Account
	res.Decode(&output)

	return &output, nil
}

type setAccountDefaultValuesOutput struct {
	CreateAccount *model.CreateAccount
	UpdateAccount *model.UpdateAccount
}

func (accDataSourceMongo *accountDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setAccountDefaultValuesOutput, error) {
	defaultAccountStatus := model.AccountStatusActive
	defaultAccountType := model.AccountTypePerson
	currentTime := time.Now()

	updateInput := input.(model.UpdateAccount)
	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		existingObject, err := accDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}

		if &(*existingObject).Status == nil {
			updateInput.Status = &defaultAccountStatus
		}
		if &(*existingObject).Type == nil {
			updateInput.Type = &defaultAccountType
		}
		updateInput.UpdatedAt = &currentTime

		return &setAccountDefaultValuesOutput{
			UpdateAccount: &updateInput,
		}, nil
	}
	createInput := (input).(model.CreateAccount)

	if &createInput.Status == nil {
		createInput.Status = &defaultAccountStatus
	}
	if &createInput.Type == nil {
		createInput.Type = defaultAccountType
	}
	if createInput.DeviceTokens == nil {
		createInput.DeviceTokens = []*string{}
	}
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime

	return &setAccountDefaultValuesOutput{
		CreateAccount: &createInput,
	}, nil
}
