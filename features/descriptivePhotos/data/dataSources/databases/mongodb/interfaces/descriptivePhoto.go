package mongodbdescriptivePhotodatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DescriptivePhotoDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.DescriptivePhoto, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.DescriptivePhoto, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.DescriptivePhoto, error)
	Create(input *model.InternalCreateDescriptivePhoto, operationOptions *mongodbcoretypes.OperationOptions) (*model.DescriptivePhoto, error)
	Update(ID primitive.ObjectID, updateData *model.InternalUpdateDescriptivePhoto, operationOptions *mongodbcoretypes.OperationOptions) (*model.DescriptivePhoto, error)
	GenerateObjectID() primitive.ObjectID
}
