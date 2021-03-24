package mongodbaccountdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
)

type personRepoMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewPersonRepoMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbaccountdatasourceinterfaces.PersonRepoMongo, error) {
	basicOperation.SetCollection("persons")
	return &personRepoMongo{
		basicOperation: basicOperation,
	}, nil
}

func (prsnRepoMongo *personRepoMongo) FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error) {
	res, err := prsnRepoMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.Person
	res.Decode(&output)
	return &output, err
}

func (prsnRepoMongo *personRepoMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error) {
	res, err := prsnRepoMongo.basicOperation.FindOne(query, operationOptions)
	var output model.Person
	res.Decode(&output)
	return &output, err
}

func (prsnRepoMongo *personRepoMongo) Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.Person, error) {
	var persons = []*model.Person{}
	cursorDecoder := func(cursor *mongodbcoretypes.CursorObject) (interface{}, error) {
		var person *model.Person
		err := cursor.MongoFindCursor.Decode(person)
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
		return nil, nil
	}

	_, err := prsnRepoMongo.basicOperation.Find(query, cursorDecoder, operationOptions)
	if err != nil {
		return nil, err
	}

	return persons, err
}

func (prsnRepoMongo *personRepoMongo) Create(input *model.CreatePerson, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error) {
	defaultedInput, err := prsnRepoMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
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
		CreatedAt:                   personOutput.CreatedAt,
		UpdatedAt:                   personOutput.UpdatedAt,
	}

	return person, err
}

func (prsnRepoMongo *personRepoMongo) Update(ID interface{}, updateData *model.UpdatePerson, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error) {
	defaultedInput, err := prsnRepoMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := prsnRepoMongo.basicOperation.Update(ID, *defaultedInput.UpdatePerson, operationOptions)
	var output model.Person
	res.Decode(&output)

	return &output, err
}

type setPersonDefaultValuesOutput struct {
	CreatePerson *model.CreatePerson
	UpdatePerson *model.UpdatePerson
}

func (prsnRepoMongo *personRepoMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setPersonDefaultValuesOutput, error) {
	var noOfRecentTransactionToKeep int

	var currentTime = time.Now()
	updateInput := input.(model.UpdatePerson)
	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
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
				UpdatedAt:                   &currentTime,
			},
		}, nil
	}
	createInput := (input).(model.CreatePerson)

	if &createInput.NoOfRecentTransactionToKeep == nil {
		noOfRecentTransactionToKeep = 15
	}

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
			UpdatedAt:                   &currentTime,
		},
	}, nil
}
