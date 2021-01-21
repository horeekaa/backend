package mongorepos

import (
	"time"

	model "github.com/horeekaa/backend/model"
	databaseclient "github.com/horeekaa/backend/repositories/databaseClient/mongoDB"
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
)

type PersonRepoMongo struct {
	basicOperation *mongooperations.BasicOperation
}

func NewPersonRepoMongo(mongoRepo *databaseclient.MongoRepository) *PersonRepoMongo {
	return &PersonRepoMongo{
		basicOperation: &mongooperations.BasicOperation{
			Client:         (*mongoRepo).Client,
			CollectionRef:  (*mongoRepo.Client.Database((*mongoRepo).DatabaseName)).Collection("persons"),
			Timeout:        (*mongoRepo).Timeout,
			CollectionName: "persons",
		},
	}
}

func (prsnRepoMongo *PersonRepoMongo) FindByID(id string, operationOptions *mongooperations.OperationOptions) (*model.Person, error) {
	object, err := prsnRepoMongo.basicOperation.FindByID(id, operationOptions)

	return ((*object).(*model.Person)), err
}

func (prsnRepoMongo *PersonRepoMongo) FindOne(query mongooperations.OperationQueryType, operationOptions *mongooperations.OperationOptions) (*model.Person, error) {
	object, err := prsnRepoMongo.basicOperation.FindOne(query, operationOptions)

	return ((*object).(*model.Person)), err
}

func (prsnRepoMongo *PersonRepoMongo) Find(query mongooperations.OperationQueryType, operationOptions *mongooperations.OperationOptions) ([]*model.Person, error) {
	objects, err := prsnRepoMongo.basicOperation.Find(query, operationOptions)

	var persons = []*model.Person{}
	for _, obj := range objects {
		persons = append(persons, (*obj).(*model.Person))
	}

	return persons, err
}

type personCreateOutput struct {
	ID     string
	Object model.Person
}

func (prsnRepoMongo *PersonRepoMongo) Create(input *model.CreateAccount, operationOptions *mongooperations.OperationOptions) (*model.Person, error) {
	defaultedInput, err := prsnRepoMongo.setDefaultValues(*input,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	object, err := prsnRepoMongo.basicOperation.Create(*defaultedInput, operationOptions)
	if err != nil {
		return nil, err
	}

	createOutputObject := (*object).(*personCreateOutput)

	person := &model.Person{
		ID:                          createOutputObject.ID,
		FirstName:                   createOutputObject.Object.FirstName,
		LastName:                    createOutputObject.Object.LastName,
		Gender:                      createOutputObject.Object.Gender,
		PhoneNumber:                 createOutputObject.Object.PhoneNumber,
		Email:                       createOutputObject.Object.Email,
		NoOfRecentTransactionToKeep: createOutputObject.Object.NoOfRecentTransactionToKeep,
	}

	return person, err
}

func (prsnRepoMongo *PersonRepoMongo) Update(ID string, updateData *model.UpdateAccount, operationOptions *mongooperations.OperationOptions) (*model.Person, error) {
	defaultedInput, err := prsnRepoMongo.setDefaultValues(*updateData,
		&defaultValuesOptions{DefaultValuesType: DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	object, err := prsnRepoMongo.basicOperation.Update(ID, *defaultedInput, operationOptions)

	return ((*object).(*model.Person)), err
}

type setPersonDefaultValuesOutput struct {
	CreatePerson *model.CreatePerson
	UpdatePerson *model.UpdatePerson
}

func (prsnRepoMongo *PersonRepoMongo) setDefaultValues(input interface{}, options *defaultValuesOptions, operationOptions *mongooperations.OperationOptions) (*setPersonDefaultValuesOutput, error) {
	var noOfRecentTransactionToKeep int

	updateInput := input.(model.UpdatePerson)
	if (*options).DefaultValuesType == DefaultValuesUpdateType {
		existingObject, err := prsnRepoMongo.FindByID(*updateInput.ID, operationOptions)
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
