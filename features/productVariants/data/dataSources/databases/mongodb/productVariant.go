package mongodbproductvariantdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	mongodbproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type productVariantDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewProductVariantDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbproductvariantdatasourceinterfaces.ProductVariantDataSourceMongo, error) {
	basicOperation.SetCollection("productvariants")
	return &productVariantDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
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
	defaultedInput, err := prodVarDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.ProductVariant
	_, err = prodVarDataSourceMongo.basicOperation.Create(*defaultedInput.CreateProductvariant, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (prodVarDataSourceMongo *productVariantDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateProductVariant, operationOptions *mongodbcoretypes.OperationOptions) (*model.ProductVariant, error) {
	updateData.ID = ID
	defaultedInput, err := prodVarDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.ProductVariant
	_, err = prodVarDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateProductvariant, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

type setProductVariantDefaultValuesOutput struct {
	CreateProductvariant *model.DatabaseCreateProductVariant
	UpdateProductvariant *model.DatabaseUpdateProductVariant
}

func (prodVarDataSourceMongo *productVariantDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setProductVariantDefaultValuesOutput, error) {
	currentTime := time.Now()
	defaultIsActive := true

	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.DatabaseUpdateProductVariant)
		_, err := prodVarDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}
		updateInput.UpdatedAt = &currentTime

		return &setProductVariantDefaultValuesOutput{
			UpdateProductvariant: &updateInput,
		}, nil
	}
	createInput := (input).(model.DatabaseCreateProductVariant)
	if createInput.IsActive == nil {
		createInput.IsActive = &defaultIsActive
	}
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime

	return &setProductVariantDefaultValuesOutput{
		CreateProductvariant: &createInput,
	}, nil
}
