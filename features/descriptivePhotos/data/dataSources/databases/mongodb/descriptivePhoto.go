package mongodbdescriptivePhotodatasources

import (
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
	var outputModel model.DescriptivePhoto
	_, err := descPhotoDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (descPhotoDataSourceMongo *descriptivePhotoDataSourceMongo) Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateDescriptivePhoto, operationOptions *mongodbcoretypes.OperationOptions) (*model.DescriptivePhoto, error) {
	existingObject, err := descPhotoDataSourceMongo.FindOne(updateCriteria, operationOptions)
	if err != nil {
		return nil, err
	}
	if existingObject == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			descPhotoDataSourceMongo.pathIdentity,
			nil,
		)
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
