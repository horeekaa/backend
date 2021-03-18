package mongodbaccountdatasources

import (
	databaseclient "github.com/horeekaa/backend/core/databaseClient/mongoDB"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	mongodbcoreoperations "github.com/horeekaa/backend/core/databaseClient/mongoDB/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
)

type accountRepoMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewAccountRepoMongo(mongoRepo *databaseclient.MongoRepository) (mongodbaccountdatasourceinterfaces.AccountRepoMongo, error) {
	basicOperation, err := mongodbcoreoperations.NewBasicOperation(
		(*mongoRepo).Client,
		(*mongoRepo.Client.Database((*mongoRepo).DatabaseName)).Collection("accounts"),
		(*mongoRepo).Timeout,
		"accounts",
	)
	if err != nil {
		return nil, err
	}

	return &accountRepoMongo{
		basicOperation: basicOperation,
	}, nil
}

func (accRepoMongo *accountRepoMongo) FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	res, err := accRepoMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.Account
	res.Decode(&output)
	return &output, err
}

func (accRepoMongo *accountRepoMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	res, err := accRepoMongo.basicOperation.FindOne(query, operationOptions)
	var output model.Account
	res.Decode(&output)
	return &output, err
}

func (accRepoMongo *accountRepoMongo) Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.Account, error) {
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

	_, err := accRepoMongo.basicOperation.Find(query, cursorDecoder, operationOptions)
	if err != nil {
		return nil, err
	}

	return accounts, err
}

func (accRepoMongo *accountRepoMongo) Create(input *model.CreateAccount, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	defaultedInput, err := accRepoMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	output, err := accRepoMongo.basicOperation.Create(*defaultedInput.CreateAccount, operationOptions)
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

func (accRepoMongo *accountRepoMongo) Update(ID interface{}, updateData *model.UpdateAccount, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error) {
	defaultedInput, err := accRepoMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := accRepoMongo.basicOperation.Update(ID, *defaultedInput.UpdateAccount, operationOptions)
	var output model.Account
	res.Decode(&output)

	return &output, err
}

type setAccountDefaultValuesOutput struct {
	CreateAccount *model.CreateAccount
	UpdateAccount *model.UpdateAccount
}

func (accRepoMongo *accountRepoMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setAccountDefaultValuesOutput, error) {
	var accountStatus model.AccountStatus
	var accountType model.AccountType

	updateInput := input.(model.UpdateAccount)
	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		existingObject, err := accRepoMongo.FindByID(updateInput.ID, operationOptions)
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
		},
	}, nil
}
