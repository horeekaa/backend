package mongorepos

import (
	model "github.com/horeekaa/backend/model"
	databaseclient "github.com/horeekaa/backend/repositories/databaseClient/mongoDB"
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
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

func (accRepoMongo *AccountRepoMongo) FindByID(ID interface{}, operationOptions *mongooperations.OperationOptions) (*model.Account, error) {
	object, err := accRepoMongo.basicOperation.FindByID(ID, operationOptions)
	output := (*object).(model.Account)
	return &output, err
}

func (accRepoMongo *AccountRepoMongo) FindOne(query mongooperations.OperationQueryType, operationOptions *mongooperations.OperationOptions) (*model.Account, error) {
	object, err := accRepoMongo.basicOperation.FindOne(query, operationOptions)
	output := (*object).(model.Account)
	return &output, err
}

func (accRepoMongo *AccountRepoMongo) Find(query mongooperations.OperationQueryType, operationOptions *mongooperations.OperationOptions) ([]*model.Account, error) {
	objects, err := accRepoMongo.basicOperation.Find(query, operationOptions)

	var accounts = []*model.Account{}
	for _, obj := range objects {
		account := (*obj).(model.Account)
		accounts = append(accounts, &account)
	}

	return accounts, err
}

func (accRepoMongo *AccountRepoMongo) Create(input *model.CreateAccount, operationOptions *mongooperations.OperationOptions) (*model.Account, error) {
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

func (accRepoMongo *AccountRepoMongo) Update(ID interface{}, updateData *model.UpdateAccount, operationOptions *mongooperations.OperationOptions) (*model.Account, error) {
	defaultedInput, err := accRepoMongo.setDefaultValues(*updateData,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	object, err := accRepoMongo.basicOperation.Update(ID, *defaultedInput.UpdateAccount, operationOptions)
	output := (*object).(model.Account)

	return &output, err
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
