package mongodbproductdatasources

import (
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type productDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
	pathIdentity   string
}

func NewProductDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbproductdatasourceinterfaces.ProductDataSourceMongo, error) {
	basicOperation.SetCollection("products")
	return &productDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "ProductDataSource",
	}, nil
}

func (prodDataSourceMongo *productDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (prodDataSourceMongo *productDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Product, error) {
	var output model.Product
	_, err := prodDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (prodDataSourceMongo *productDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Product, error) {
	var output model.Product
	_, err := prodDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (prodDataSourceMongo *productDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Product, error) {
	var products = []*model.Product{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var product model.Product
		if err := cursor.Decode(&product); err != nil {
			return err
		}
		products = append(products, &product)
		return nil
	}
	_, err := prodDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return products, err
}

func (prodDataSourceMongo *productDataSourceMongo) Create(input *model.DatabaseCreateProduct, operationOptions *mongodbcoretypes.OperationOptions) (*model.Product, error) {
	var outputModel model.Product
	_, err := prodDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (prodDataSourceMongo *productDataSourceMongo) Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateProduct, operationOptions *mongodbcoretypes.OperationOptions) (*model.Product, error) {
	existingObject, err := prodDataSourceMongo.FindOne(updateCriteria, operationOptions)
	if err != nil {
		return nil, err
	}
	if existingObject == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			prodDataSourceMongo.pathIdentity,
			nil,
		)
	}

	var output model.Product
	_, err = prodDataSourceMongo.basicOperation.Update(
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
