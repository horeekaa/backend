package productdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getProductRepository struct {
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource
	pathIdentity      string
}

func NewGetProductRepository(
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
) (productdomainrepositoryinterfaces.GetProductRepository, error) {
	return &getProductRepository{
		productDataSource,
		"GetProductRepository",
	}, nil
}

func (getProductRepo *getProductRepository) Execute(filterFields *model.ProductFilterFields) (*model.Product, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	product, err := getProductRepo.productDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getProductRepo.pathIdentity,
			err,
		)
	}

	return product, nil
}
