package productdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	productdomainrepositorytypes "github.com/horeekaa/backend/features/products/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getAllProductRepository struct {
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource
}

func NewGetAllProductRepository(
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
) (productdomainrepositoryinterfaces.GetAllProductRepository, error) {
	return &getAllProductRepository{
		productDataSource,
	}, nil
}

func (getAllMmbAccRefRepo *getAllProductRepository) Execute(
	input productdomainrepositorytypes.GetAllProductInput,
) ([]*model.Product, error) {
	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(input.FilterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	products, err := getAllMmbAccRefRepo.productDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAllproduct",
			err,
		)
	}

	return products, nil
}
