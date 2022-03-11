package mongodbproductdatasources

import (
	"time"

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
	_, err := prodDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Product
	_, err = prodDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (prodDataSourceMongo *productDataSourceMongo) Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateProduct, operationOptions *mongodbcoretypes.OperationOptions) (*model.Product, error) {
	_, err := prodDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
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

func (prodDataSourceMongo *productDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.DatabaseUpdateProduct,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	currentTime := time.Now()
	existingObject, err := prodDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			prodDataSourceMongo.pathIdentity,
			nil,
		)
	}

	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}

func (prodDataSourceMongo *productDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.DatabaseCreateProduct,
) (bool, error) {
	currentTime := time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed
	defaultIsActive := true

	if input.ProposalStatus == nil {
		input.ProposalStatus = &defaultProposalStatus
	}
	if input.IsActive == nil {
		input.IsActive = &defaultIsActive
	}
	if input.Photos == nil {
		input.Photos = []*model.ObjectIDOnly{}
	}
	if input.Variants == nil {
		input.Variants = []*model.ObjectIDOnly{}
	}
	if input.Taggings == nil {
		input.Taggings = []*model.ObjectIDOnly{}
	}
	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime
	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}
