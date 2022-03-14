package productdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	productdomainrepositorytypes "github.com/horeekaa/backend/features/products/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllProductRepository struct {
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder
	pathIdentity      string
}

func NewGetAllProductRepository(
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (productdomainrepositoryinterfaces.GetAllProductRepository, error) {
	return &getAllProductRepository{
		productDataSource,
		mongoQueryBuilder,
		"GetAllProductRepository",
	}, nil
}

func (getAllProductRepo *getAllProductRepository) Execute(
	input productdomainrepositorytypes.GetAllProductInput,
) ([]*model.Product, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllProductRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	products, err := getAllProductRepo.productDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAllProductRepo.pathIdentity,
			err,
		)
	}

	return products, nil
}
