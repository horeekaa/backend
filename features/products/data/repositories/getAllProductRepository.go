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
}

func NewGetAllProductRepository(
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (productdomainrepositoryinterfaces.GetAllProductRepository, error) {
	return &getAllProductRepository{
		productDataSource,
		mongoQueryBuilder,
	}, nil
}

func (getAllMmbAccRefRepo *getAllProductRepository) Execute(
	input productdomainrepositorytypes.GetAllProductInput,
) ([]*model.Product, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllMmbAccRefRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

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
