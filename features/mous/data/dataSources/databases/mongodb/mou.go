package mongodbmoudatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
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
	defaultedInput, err := mouDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Mou
	_, err = mouDataSourceMongo.basicOperation.Create(*defaultedInput.CreateMou, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (mouDataSourceMongo *mouDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateMou, operationOptions *mongodbcoretypes.OperationOptions) (*model.Mou, error) {
	updateData.ID = ID
	defaultedInput, err := mouDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.Mou
	_, err = mouDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateMou, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

type setmouDefaultValuesOutput struct {
	CreateMou *model.DatabaseCreateMou
	UpdateMou *model.DatabaseUpdateMou
}

func (mouDataSourceMongo *mouDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setmouDefaultValuesOutput, error) {
	currentTime := time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed
	defaultIsActive := true

	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.DatabaseUpdateMou)
		_, err := mouDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}
		updateInput.UpdatedAt = &currentTime

		return &setmouDefaultValuesOutput{
			UpdateMou: &updateInput,
		}, nil
	}
	createInput := (input).(model.DatabaseCreateMou)
	if createInput.ProposalStatus == nil {
		createInput.ProposalStatus = &defaultProposalStatus
	}
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime
	if createInput.IsActive == nil {
		createInput.IsActive = &defaultIsActive
	}

	return &setmouDefaultValuesOutput{
		CreateMou: &createInput,
	}, nil
}
