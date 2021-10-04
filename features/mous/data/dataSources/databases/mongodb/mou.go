package mongodbmoudatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbmoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mouDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewMouDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbmoudatasourceinterfaces.MouDataSourceMongo, error) {
	basicOperation.SetCollection("mous")
	return &mouDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (mouDataSourceMongo *mouDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (mouDataSourceMongo *mouDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Mou, error) {
	var output model.Mou
	_, err := mouDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (mouDataSourceMongo *mouDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Mou, error) {
	var output model.Mou
	_, err := mouDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (mouDataSourceMongo *mouDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Mou, error) {
	var mous = []*model.Mou{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var mou model.Mou
		if err := cursor.Decode(&mou); err != nil {
			return err
		}
		mous = append(mous, &mou)
		return nil
	}
	_, err := mouDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return mous, err
}

func (mouDataSourceMongo *mouDataSourceMongo) Create(input *model.DatabaseCreateMou, operationOptions *mongodbcoretypes.OperationOptions) (*model.Mou, error) {
	_, err := mouDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Mou
	_, err = mouDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (mouDataSourceMongo *mouDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdateMou,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.Mou, error) {
	_, err := mouDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.Mou
	_, err = mouDataSourceMongo.basicOperation.Update(
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

func (mouDataSourceMongo *mouDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.DatabaseUpdateMou,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	currentTime := time.Now()
	existingObject, err := mouDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			"/mouDataSource/update",
			nil,
		)
	}

	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}

func (mouDataSourceMongo *mouDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.DatabaseCreateMou,
) (bool, error) {
	currentTime := time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed
	defaultIsActive := true

	if input.ProposalStatus == nil {
		input.ProposalStatus = &defaultProposalStatus
	}
	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime
	if input.IsActive == nil {
		input.IsActive = &defaultIsActive
	}
	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}
