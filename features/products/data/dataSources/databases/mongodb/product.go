package mongodbproductdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	mongodbproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type productDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewProductDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbproductdatasourceinterfaces.ProductDataSourceMongo, error) {
	basicOperation.SetCollection("products")
	return &productDataSourceMongo{
		basicOperation: basicOperation,
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
	defaultedInput, err := prodDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Product
	_, err = prodDataSourceMongo.basicOperation.Create(*defaultedInput.Createproduct, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (prodDataSourceMongo *productDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateProduct, operationOptions *mongodbcoretypes.OperationOptions) (*model.Product, error) {
	updateData.ID = ID
	defaultedInput, err := prodDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.Product
	_, err = prodDataSourceMongo.basicOperation.Update(ID, *defaultedInput.Updateproduct, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

type setProductDefaultValuesOutput struct {
	Createproduct *model.DatabaseCreateProduct
	Updateproduct *model.DatabaseUpdateProduct
}

func (prodDataSourceMongo *productDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setProductDefaultValuesOutput, error) {
	currentTime := time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed
	defaultIsActive := true

	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.DatabaseUpdateProduct)
		_, err := prodDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}
		updateInput.UpdatedAt = &currentTime

		return &setProductDefaultValuesOutput{
			Updateproduct: &updateInput,
		}, nil
	}
	createInput := (input).(model.DatabaseCreateProduct)
	if createInput.ProposalStatus == nil {
		createInput.ProposalStatus = &defaultProposalStatus
	}
	if createInput.IsActive == nil {
		createInput.IsActive = &defaultIsActive
	}
	if createInput.Photos == nil {
		createInput.Photos = []*model.ObjectIDOnly{}
	}
	if createInput.Variants == nil {
		createInput.Variants = []*model.ObjectIDOnly{}
	}
	if createInput.Taggings == nil {
		createInput.Taggings = []*model.ObjectIDOnly{}
	}
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime

	return &setProductDefaultValuesOutput{
		Createproduct: &createInput,
	}, nil
}
