package mongodbaccountdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type personDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
	pathIdentity   string
}

func NewPersonDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbaccountdatasourceinterfaces.PersonDataSourceMongo, error) {
	basicOperation.SetCollection("persons")
	return &personDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "PersonDataSource",
	}, nil
}

func (prsnDataSourceMongo *personDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error) {
	var output model.Person
	_, err := prsnDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (prsnDataSourceMongo *personDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error) {
	var output model.Person
	_, err := prsnDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (prsnDataSourceMongo *personDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Person, error) {
	var persons = []*model.Person{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var person model.Person
		if err := cursor.Decode(&person); err != nil {
			return err
		}
		persons = append(persons, &person)
		return nil
	}
	_, err := prsnDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return persons, err
}

func (prsnDataSourceMongo *personDataSourceMongo) Create(input *model.CreatePerson, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error) {
	_, err := prsnDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Person
	_, err = prsnDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (prsnDataSourceMongo *personDataSourceMongo) Update(updateCriteria map[string]interface{}, updateData *model.UpdatePerson, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error) {
	_, err := prsnDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.Person
	_, err = prsnDataSourceMongo.basicOperation.Update(
		updateCriteria,
		map[string]interface{}{
			"$set": updateData,
		},
		&output,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (prsnDataSourceMongo *personDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.UpdatePerson,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	var currentTime = time.Now()
	defaultNoOfRecentTransactionToKeep := 15

	existingObject, err := prsnDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			prsnDataSourceMongo.pathIdentity,
			nil,
		)
	}

	if &(*existingObject).NoOfRecentTransactionToKeep == nil {
		input.NoOfRecentTransactionToKeep = &defaultNoOfRecentTransactionToKeep
	}
	input.UpdatedAt = &currentTime

	return true, nil
}

func (prsnDataSourceMongo *personDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.CreatePerson,
) (bool, error) {
	defaultNoOfRecentTransactionToKeep := 15

	var currentTime = time.Now()

	if &input.NoOfRecentTransactionToKeep == nil {
		input.NoOfRecentTransactionToKeep = &defaultNoOfRecentTransactionToKeep
	}
	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime

	return true, nil
}
