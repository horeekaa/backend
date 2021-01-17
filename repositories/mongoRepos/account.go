package mongorepos

import (
	model "github.com/horeekaa/backend/model"
	databaseclient "github.com/horeekaa/backend/repositories/databaseClient"
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/operations"
)

type AccountRepoMongo struct {
	basicOperation *mongooperations.BasicOperation
}

func NewAccountRepoMongo(mongoRepo *databaseclient.MongoRepository) *AccountRepoMongo {
	return &AccountRepoMongo{
		basicOperation: &mongooperations.BasicOperation{
			Client:         (*mongoRepo).Client,
			CollectionRef:  (*mongoRepo.Client.Database((*mongoRepo).DatabaseName)).Collection("accounts"),
			Timeout:        (*mongoRepo).Timeout,
			CollectionName: "accounts",
		},
	}
}

func (accRepoMongo *AccountRepoMongo) FindByID(id string, operationOptions *mongooperations.OperationOptions) (*model.Account, error) {
	object, err := accRepoMongo.basicOperation.FindByID(id, operationOptions)

	return ((*object).(*model.Account)), err
}

func (accRepoMongo *AccountRepoMongo) FindOne(query map[string]interface{}, operationOptions *mongooperations.OperationOptions) (*model.Account, error) {
	object, err := accRepoMongo.basicOperation.FindOne(query, operationOptions)

	return ((*object).(*model.Account)), err
}

func (accRepoMongo *AccountRepoMongo) Find(query map[string]interface{}, operationOptions *mongooperations.OperationOptions) ([]*model.Account, error) {
	objects, err := accRepoMongo.basicOperation.Find(query, operationOptions)

	var accounts = []*model.Account{}
	for _, obj := range objects {
		accounts = append(accounts, (*obj).(*model.Account))
	}

	return accounts, err
}

type createOutput struct {
	ID     string
	Object model.Account
}

func (accRepoMongo *AccountRepoMongo) Create(input *model.CreateAccount, operationOptions *mongooperations.OperationOptions) (*model.Account, error) {
	defaultedInput, err := accRepoMongo.setDefaultValues(*input,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	object, err := accRepoMongo.basicOperation.Create(*defaultedInput, operationOptions)
	if err != nil {
		return nil, err
	}

	createOutputObject := (*object).(*createOutput)

	account := &model.Account{
		ID:           createOutputObject.ID,
		Status:       createOutputObject.Object.Status,
		StatusReason: createOutputObject.Object.StatusReason,
		Type:         createOutputObject.Object.Type,
		Person:       createOutputObject.Object.Person,
	}

	return account, err
}

func (accRepoMongo *AccountRepoMongo) Update(ID string, updateData *model.UpdateAccount, operationOptions *mongooperations.OperationOptions) (*model.Account, error) {
	defaultedInput, err := accRepoMongo.setDefaultValues(*updateData,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	object, err := accRepoMongo.basicOperation.Update(ID, *defaultedInput, operationOptions)

	return ((*object).(*model.Account)), err
}

type setDefaultValuesOutput struct {
	CreateAccount *model.CreateAccount
	UpdateAccount *model.UpdateAccount
}

func (accRepoMongo *AccountRepoMongo) setDefaultValues(input interface{}, options *defaultValuesOptions, operationOptions *mongooperations.OperationOptions) (*setDefaultValuesOutput, error) {
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
