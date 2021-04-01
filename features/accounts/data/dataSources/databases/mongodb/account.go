package mongodbaccountdatasources

import (
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
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

func (accDataSourceMongo *accountDataSourceMongo) FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	res, err := accDataSourceMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.Account
	res.Decode(&output)
	return &output, err
}

func (accDataSourceMongo *accountDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	res, err := accDataSourceMongo.basicOperation.FindOne(query, operationOptions)
	var output model.Account
	res.Decode(&output)
	return &output, err
}

func (accDataSourceMongo *accountDataSourceMongo) Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.Account, error) {
	var accounts = []*model.Account{}
	cursorDecoder := func(cursor *mongodbcoretypes.CursorObject) (interface{}, error) {
		var account *model.Account
		err := cursor.MongoFindCursor.Decode(account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
		return nil, nil
	}

	_, err := accDataSourceMongo.basicOperation.Find(query, cursorDecoder, operationOptions)
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

	account := &model.Account{
		ID:           output.ID,
		Status:       accountOutput.Status,
		StatusReason: accountOutput.StatusReason,
		Type:         accountOutput.Type,
		Person:       accountOutput.Person,
	}

	return account, err
}

func (accDataSourceMongo *accountDataSourceMongo) Update(ID interface{}, updateData *model.UpdateAccount, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	defaultedInput, err := accDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := accDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateAccount, operationOptions)
	var output model.Account
	res.Decode(&output)

	return &output, err
}

type setAccountDefaultValuesOutput struct {
	CreateAccount *model.CreateAccount
	UpdateAccount *model.UpdateAccount
}

func (accDataSourceMongo *accountDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setAccountDefaultValuesOutput, error) {
	var accountStatus model.AccountStatus
	var accountType model.AccountType

	updateInput := input.(model.UpdateAccount)
	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		existingObject, err := accDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}

		if &(*existingObject).Status == nil {
			accountStatus = model.AccountStatusActive
		}
		if &(*existingObject).Type == nil {
			accountType = model.AccountTypePerson
		}

		return &setAccountDefaultValuesOutput{
			UpdateAccount: &model.UpdateAccount{
				ID:           updateInput.ID,
				Status:       &accountStatus,
				StatusReason: updateInput.StatusReason,
				Type:         &accountType,
				Person:       updateInput.Person,
				DeviceTokens: updateInput.DeviceTokens,
			},
		}, nil
	}
	createInput := (input).(model.CreateAccount)

	if &createInput.Status == nil {
		accountStatus = model.AccountStatusActive
	}
	if &createInput.Type == nil {
		accountType = model.AccountTypePerson
	}

	return &setAccountDefaultValuesOutput{
		CreateAccount: &model.CreateAccount{
			Status:       &accountStatus,
			StatusReason: createInput.StatusReason,
			Type:         accountType,
			Person:       createInput.Person,
			DeviceTokens: []*string{},
		},
	}, nil
}
