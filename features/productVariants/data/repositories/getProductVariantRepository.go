package productvariantdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getProductVariantRepository struct {
	productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource
}

func NewGetProductVariantRepository(
	productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
) (productvariantdomainrepositoryinterfaces.GetProductVariantRepository, error) {
	return &getProductVariantRepository{
		productVariantDataSource,
	}, nil
}

func (getProductVariantRefRepo *getProductVariantRepository) Execute(filterFields *model.ProductVariantFilterFields) (*model.ProductVariant, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	productVariant, err := getProductVariantRefRepo.productVariantDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getProductVariant",
			err,
		)
	}

	return productVariant, nil
}
