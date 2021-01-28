package mongorepos

import (
	"time"

	model "github.com/horeekaa/backend/model"
	databaseclient "github.com/horeekaa/backend/repositories/databaseClient/mongoDB"
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
)

type personRepoMongo struct {
	basicOperation *mongooperations.BasicOperation
}

func NewPersonRepoMongo(mongoRepo *databaseclient.MongoRepository) *mongorepointerface.PersonRepoMongo {
	return &personRepoMongo{
		basicOperation: &mongooperations.BasicOperation{
			Client:         (*mongoRepo).Client,
			CollectionRef:  (*mongoRepo.Client.Database((*mongoRepo).DatabaseName)).Collection("persons"),
			Timeout:        (*mongoRepo).Timeout,
			CollectionName: "persons",
		},
	}
}

func (prsnRepoMongo *personRepoMongo) FindByID(ID interface{}, operationOptions *mongooperations.OperationOptions) (*model.Person, error) {
	object, err := prsnRepoMongo.basicOperation.FindByID(ID, operationOptions)
	output := (*object).(model.Person)
	return &output, err
}

func (prsnRepoMongo *personRepoMongo) FindOne(query mongooperations.OperationQueryType, operationOptions *mongooperations.OperationOptions) (*model.Person, error) {
	object, err := prsnRepoMongo.basicOperation.FindOne(query, operationOptions)
	output := (*object).(model.Person)
	return &output, err
}

func (prsnRepoMongo *personRepoMongo) Find(query mongooperations.OperationQueryType, operationOptions *mongooperations.OperationOptions) ([]*model.Person, error) {
	objects, err := prsnRepoMongo.basicOperation.Find(query, operationOptions)

	var persons = []*model.Person{}
	for _, obj := range objects {
		person := (*obj).(model.Person)
		persons = append(persons, &person)
	}

	return persons, err
}

func (prsnRepoMongo *personRepoMongo) Create(input *model.CreateAccount, operationOptions *mongooperations.OperationOptions) (*model.Person, error) {
	defaultedInput, err := prsnRepoMongo.setDefaultValues(*input,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	output, err := prsnRepoMongo.basicOperation.Create(*defaultedInput.CreatePerson, operationOptions)
	if err != nil {
		return nil, err
	}

	personOutput := output.Object.(model.Person)

	person := &model.Person{
		ID:                          output.ID,
		FirstName:                   personOutput.FirstName,
		LastName:                    personOutput.LastName,
		Gender:                      personOutput.Gender,
		PhoneNumber:                 personOutput.PhoneNumber,
		Email:                       personOutput.Email,
		NoOfRecentTransactionToKeep: personOutput.NoOfRecentTransactionToKeep,
	}

	return person, err
}

func (prsnRepoMongo *personRepoMongo) Update(ID interface{}, updateData *model.UpdateAccount, operationOptions *mongooperations.OperationOptions) (*model.Person, error) {
	defaultedInput, err := prsnRepoMongo.setDefaultValues(*updateData,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	object, err := prsnRepoMongo.basicOperation.Update(ID, *defaultedInput.UpdatePerson, operationOptions)
	output := (*object).(model.Person)

	return &output, err
}

type setPersonDefaultValuesOutput struct {
	CreatePerson *model.CreatePerson
	UpdatePerson *model.UpdatePerson
}

func (prsnRepoMongo *personRepoMongo) setDefaultValues(input interface{}, options *defaultValuesOptions, operationOptions *mongooperations.OperationOptions) (*setPersonDefaultValuesOutput, error) {
	var noOfRecentTransactionToKeep int

	updateInput := input.(model.UpdatePerson)
	if (*options).DefaultValuesType == DefaultValuesUpdateType {
		existingObject, err := prsnRepoMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}

		if &(*existingObject).NoOfRecentTransactionToKeep == nil {
			noOfRecentTransactionToKeep = 15
		}

		return &setPersonDefaultValuesOutput{
			UpdatePerson: &model.UpdatePerson{
				ID:                          updateInput.ID,
				FirstName:                   updateInput.FirstName,
				LastName:                    updateInput.LastName,
				Gender:                      updateInput.Gender,
				PhoneNumber:                 updateInput.PhoneNumber,
				Email:                       updateInput.Email,
				NoOfRecentTransactionToKeep: &noOfRecentTransactionToKeep,
			},
		}, nil
	}
	createInput := (input).(model.CreatePerson)

	if &createInput.NoOfRecentTransactionToKeep == nil {
		noOfRecentTransactionToKeep = 15
	}

	var currentTime = time.Now()
	return &setPersonDefaultValuesOutput{
		CreatePerson: &model.CreatePerson{
			FirstName:                   createInput.FirstName,
			LastName:                    createInput.LastName,
			Gender:                      createInput.Gender,
			PhoneNumber:                 createInput.PhoneNumber,
			Email:                       createInput.Email,
			DeviceTokens:                []*string{},
			NoOfRecentTransactionToKeep: &noOfRecentTransactionToKeep,
			CreatedAt:                   &currentTime,
		},
	}, nil
}
