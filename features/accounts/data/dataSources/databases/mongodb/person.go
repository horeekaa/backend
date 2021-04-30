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

type personDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewPersonDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbaccountdatasourceinterfaces.PersonDataSourceMongo, error) {
	basicOperation.SetCollection("persons")
	return &personDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (prsnDataSourceMongo *personDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error) {
	res, err := prsnDataSourceMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.Person
	res.Decode(&output)
	return &output, err
}

func (prsnDataSourceMongo *personDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error) {
	res, err := prsnDataSourceMongo.basicOperation.FindOne(query, operationOptions)
	var output model.Person
	res.Decode(&output)
	return &output, err
}

func (prsnDataSourceMongo *personDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Person, error) {
	var persons = []*model.Person{}
	cursorDecoder := func(cursor *mongo.Cursor) (interface{}, error) {
		var person *model.Person
		err := cursor.Decode(person)
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
		return nil, nil
	}

	_, err := prsnDataSourceMongo.basicOperation.Find(query, paginationOpts, cursorDecoder, operationOptions)
	if err != nil {
		return nil, err
	}

	return persons, err
}

func (prsnDataSourceMongo *personDataSourceMongo) Create(input *model.CreatePerson, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error) {
	defaultedInput, err := prsnDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	output, err := prsnDataSourceMongo.basicOperation.Create(*defaultedInput.CreatePerson, operationOptions)
	if err != nil {
		return nil, err
	}

	personOutput := output.Object.(model.Person)
	personOutput.ID = output.ID

	return &personOutput, err
}

func (prsnDataSourceMongo *personDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.UpdatePerson, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error) {
	defaultedInput, err := prsnDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := prsnDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdatePerson, operationOptions)
	var output model.Person
	res.Decode(&output)

	return &output, err
}

type setPersonDefaultValuesOutput struct {
	CreatePerson *model.CreatePerson
	UpdatePerson *model.UpdatePerson
}

func (prsnDataSourceMongo *personDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setPersonDefaultValuesOutput, error) {
	defaultNoOfRecentTransactionToKeep := 15

	var currentTime = time.Now()
	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.UpdatePerson)
		existingObject, err := prsnDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}

		if &(*existingObject).NoOfRecentTransactionToKeep == nil {
			updateInput.NoOfRecentTransactionToKeep = &defaultNoOfRecentTransactionToKeep
		}
		updateInput.UpdatedAt = &currentTime

		return &setPersonDefaultValuesOutput{
			UpdatePerson: &updateInput,
		}, nil
	}
	createInput := (input).(model.CreatePerson)

	if &createInput.NoOfRecentTransactionToKeep == nil {
		createInput.NoOfRecentTransactionToKeep = &defaultNoOfRecentTransactionToKeep
	}
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime

	return &setPersonDefaultValuesOutput{
		CreatePerson: &createInput,
	}, nil
}
