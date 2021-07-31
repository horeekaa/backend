package mongodbtagdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	mongodbtagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type tagDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewTagDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbtagdatasourceinterfaces.TagDataSourceMongo, error) {
	basicOperation.SetCollection("tags")
	return &tagDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (tagDataSourceMongo *tagDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (tagDataSourceMongo *tagDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tag, error) {
	var output model.Tag
	_, err := tagDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (tagDataSourceMongo *tagDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tag, error) {
	var output model.Tag
	_, err := tagDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (tagDataSourceMongo *tagDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Tag, error) {
	var tags = []*model.Tag{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var tag model.Tag
		if err := cursor.Decode(&tag); err != nil {
			return err
		}
		tags = append(tags, &tag)
		return nil
	}
	_, err := tagDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return tags, err
}

func (tagDataSourceMongo *tagDataSourceMongo) Create(input *model.DatabaseCreateTag, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tag, error) {
	defaultedInput, err := tagDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Tag
	_, err = tagDataSourceMongo.basicOperation.Create(*defaultedInput.CreateTag, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (tagDataSourceMongo *tagDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateTag, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tag, error) {
	updateData.ID = ID
	defaultedInput, err := tagDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.Tag
	_, err = tagDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateTag, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

type setTagDefaultValuesOutput struct {
	CreateTag *model.DatabaseCreateTag
	UpdateTag *model.DatabaseUpdateTag
}

func (tagDataSourceMongo *tagDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setTagDefaultValuesOutput, error) {
	currentTime := time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed
	defaultIsActive := true

	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.DatabaseUpdateTag)
		_, err := tagDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}
		updateInput.UpdatedAt = &currentTime

		return &setTagDefaultValuesOutput{
			UpdateTag: &updateInput,
		}, nil
	}
	createInput := (input).(model.DatabaseCreateTag)
	if createInput.ProposalStatus == nil {
		createInput.ProposalStatus = &defaultProposalStatus
	}
	if createInput.IsActive == nil {
		createInput.IsActive = &defaultIsActive
	}
	if createInput.Photos == nil {
		createInput.Photos = []*model.ObjectIDOnly{}
	}
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime

	return &setTagDefaultValuesOutput{
		CreateTag: &createInput,
	}, nil
}
