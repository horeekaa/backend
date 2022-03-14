package mongodbdescriptivePhotodatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbdescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type descriptivePhotoDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
	pathIdentity   string
}

func NewDescriptivePhotoDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbdescriptivephotodatasourceinterfaces.DescriptivePhotoDataSourceMongo, error) {
	basicOperation.SetCollection("descriptivephotos")
	return &descriptivePhotoDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "DescriptivePhotoDataSource",
	}, nil
}

func (descPhotoDataSourceMongo *descriptivePhotoDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
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
	_, err := descPhotoDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.DescriptivePhoto
	_, err = descPhotoDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (descPhotoDataSourceMongo *descriptivePhotoDataSourceMongo) Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateDescriptivePhoto, operationOptions *mongodbcoretypes.OperationOptions) (*model.DescriptivePhoto, error) {
	_, err := descPhotoDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.DescriptivePhoto
	_, err = descPhotoDataSourceMongo.basicOperation.Update(
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

func (descPhotoDataSourceMongo *descriptivePhotoDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.DatabaseUpdateDescriptivePhoto,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	currentTime := time.Now()
	existingObject, err := descPhotoDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			descPhotoDataSourceMongo.pathIdentity,
			nil,
		)
	}
	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}
	input.UpdatedAt = &currentTime

	return true, nil
}

func (descPhotoDataSourceMongo *descriptivePhotoDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.DatabaseCreateDescriptivePhoto,
) (bool, error) {
	currentTime := time.Now()

	input.IsActive = true
	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}
	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime

	return true, nil
}
