package mongorepos

import (
	model "github.com/horeekaa/backend/model"
	databaseclient "github.com/horeekaa/backend/repositories/databaseClient/mongoDB"
	mongooperationinterfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/operations"
	mongorepointerfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/repos"
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
	mongooperationmodels "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations/models"
)

type accountRepoMongo struct {
	basicOperation mongooperationinterfaces.BasicOperation
}

func NewAccountRepoMongo(mongoRepo *databaseclient.MongoRepository) (mongorepointerfaces.AccountRepoMongo, error) {
	basicOperation, err := mongooperations.NewBasicOperation(
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

func (accRepoMongo *accountRepoMongo) FindByID(ID interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.Account, error) {
	res, err := accRepoMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.Account
	res.Decode(&output)
	return &output, err
}

func (accRepoMongo *accountRepoMongo) FindOne(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.Account, error) {
	res, err := accRepoMongo.basicOperation.FindOne(query, operationOptions)
	var output model.Account
	res.Decode(&output)
	return &output, err
}

func (accRepoMongo *accountRepoMongo) Find(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) ([]*model.Account, error) {
	var accounts = []*model.Account{}
	cursorDecoder := func(cursor *mongooperationmodels.CursorObject) (interface{}, error) {
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

func (accRepoMongo *accountRepoMongo) Create(input *model.CreateAccount, operationOptions *mongooperationmodels.OperationOptions) (*model.Account, error) {
	defaultedInput, err := accRepoMongo.setDefaultValues(*input,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesCreateType},
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

func (accRepoMongo *accountRepoMongo) Update(ID interface{}, updateData *model.UpdateAccount, operationOptions *mongooperationmodels.OperationOptions) (*model.Account, error) {
	defaultedInput, err := accRepoMongo.setDefaultValues(*updateData,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesUpdateType},
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

type setDefaultValuesOutput struct {
	CreateAccount *model.CreateAccount
	UpdateAccount *model.UpdateAccount
}

func (accRepoMongo *accountRepoMongo) setDefaultValues(input interface{}, options *defaultValuesOptions, operationOptions *mongooperationmodels.OperationOptions) (*setDefaultValuesOutput, error) {
	var accountStatus model.AccountStatus
	var accountType model.AccountType

	updateInput := input.(model.UpdateAccount)
	if (*options).DefaultValuesType == DefaultValuesUpdateType {
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

		return &setDefaultValuesOutput{
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

	return &setDefaultValuesOutput{
		CreateAccount: &model.CreateAccount{
			Status:       &accountStatus,
			StatusReason: createInput.StatusReason,
			Type:         accountType,
			Person:       createInput.Person,
		},
	}, nil
}
