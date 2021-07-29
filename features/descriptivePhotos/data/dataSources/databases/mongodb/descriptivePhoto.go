package mongodbdescriptivePhotodatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	mongodbdescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type descriptivePhotoDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewDescriptivePhotoDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbdescriptivephotodatasourceinterfaces.DescriptivePhotoDataSourceMongo, error) {
	basicOperation.SetCollection("descriptivephotos")
	return &descriptivePhotoDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (descPhotoDataSourceMongo *descriptivePhotoDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.DescriptivePhoto, error) {
	var output model.DescriptivePhoto
	_, err := descPhotoDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (descPhotoDataSourceMongo *descriptivePhotoDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.DescriptivePhoto, error) {
	var output model.DescriptivePhoto
	_, err := descPhotoDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (descPhotoDataSourceMongo *descriptivePhotoDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.DescriptivePhoto, error) {
	var descriptivePhotos = []*model.DescriptivePhoto{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var descriptivePhoto model.DescriptivePhoto
		if err := cursor.Decode(&descriptivePhoto); err != nil {
			return err
		}
		descriptivePhotos = append(descriptivePhotos, &descriptivePhoto)
		return nil
	}
	_, err := descPhotoDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return descriptivePhotos, err
}

func (descPhotoDataSourceMongo *descriptivePhotoDataSourceMongo) Create(input *model.DatabaseCreateDescriptivePhoto, operationOptions *mongodbcoretypes.OperationOptions) (*model.DescriptivePhoto, error) {
	defaultedInput, err := descPhotoDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.DescriptivePhoto
	_, err = descPhotoDataSourceMongo.basicOperation.Create(*defaultedInput.CreateDescriptivePhoto, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (descPhotoDataSourceMongo *descriptivePhotoDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateDescriptivePhoto, operationOptions *mongodbcoretypes.OperationOptions) (*model.DescriptivePhoto, error) {
	updateData.ID = ID
	defaultedInput, err := descPhotoDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.DescriptivePhoto
	_, err = descPhotoDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateDescriptivePhoto, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

type setDescriptivePhotoDefaultValuesOutput struct {
	CreateDescriptivePhoto *model.DatabaseCreateDescriptivePhoto
	UpdateDescriptivePhoto *model.DatabaseUpdateDescriptivePhoto
}

func (descPhotoDataSourceMongo *descriptivePhotoDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setDescriptivePhotoDefaultValuesOutput, error) {
	currentTime := time.Now()

	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.DatabaseUpdateDescriptivePhoto)
		_, err := descPhotoDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}
		updateInput.UpdatedAt = &currentTime

		return &setDescriptivePhotoDefaultValuesOutput{
			UpdateDescriptivePhoto: &updateInput,
		}, nil
	}
	createInput := (input).(model.DatabaseCreateDescriptivePhoto)
	createInput.IsActive = true
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime

	return &setDescriptivePhotoDefaultValuesOutput{
		CreateDescriptivePhoto: &createInput,
	}, nil
}
