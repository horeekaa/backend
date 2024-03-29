package mongodbproductvariantdatasources

import (
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type productVariantDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
	pathIdentity   string
}

func NewProductVariantDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbproductvariantdatasourceinterfaces.ProductVariantDataSourceMongo, error) {
	basicOperation.SetCollection("productvariants")
	return &productVariantDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "ProductVariantDataSource",
	}, nil
}

func (prodVarDataSourceMongo *productVariantDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (prodVarDataSourceMongo *productVariantDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.ProductVariant, error) {
	var output model.ProductVariant
	_, err := prodVarDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (prodVarDataSourceMongo *productVariantDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.ProductVariant, error) {
	var output model.ProductVariant
	_, err := prodVarDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (prodVarDataSourceMongo *productVariantDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.ProductVariant, error) {
	var productVariants = []*model.ProductVariant{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var productvariant model.ProductVariant
		if err := cursor.Decode(&productvariant); err != nil {
			return err
		}
		productVariants = append(productVariants, &productvariant)
		return nil
	}
	_, err := prodVarDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return productVariants, err
}

func (prodVarDataSourceMongo *productVariantDataSourceMongo) Create(input *model.DatabaseCreateProductVariant, operationOptions *mongodbcoretypes.OperationOptions) (*model.ProductVariant, error) {
	var outputModel model.ProductVariant
	_, err := prodVarDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (prodVarDataSourceMongo *productVariantDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdateProductVariant,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.ProductVariant, error) {
	existingObject, err := prodVarDataSourceMongo.FindOne(updateCriteria, operationOptions)
	if err != nil {
		return nil, err
	}
	if existingObject == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			prodVarDataSourceMongo.pathIdentity,
			nil,
		)
	}

	var output model.ProductVariant
	_, err = prodVarDataSourceMongo.basicOperation.Update(
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
