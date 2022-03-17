package mongodbmoudatasources

import (
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
	pathIdentity   string
}

func NewMouDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbmoudatasourceinterfaces.MouDataSourceMongo, error) {
	basicOperation.SetCollection("mous")
	return &mouDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "MouDataSource",
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
	var outputModel model.Mou
	_, err := mouDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
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
	existingObject, err := mouDataSourceMongo.FindOne(nil, operationOptions)
	if err != nil {
		return nil, err
	}
	if existingObject == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			mouDataSourceMongo.pathIdentity,
			nil,
		)
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
