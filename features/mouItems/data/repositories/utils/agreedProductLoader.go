package mouitemdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	mouitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories/utils"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type agreedProductLoader struct {
	productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource
	productDataSource        databaseproductdatasourceinterfaces.ProductDataSource
	mapProcessorUtility      coreutilityinterfaces.MapProcessorUtility
}

func NewAgreedProductLoader(
	productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader, error) {
	return &agreedProductLoader{
		productVariantDataSource: productVariantDataSource,
		productDataSource:        productDataSource,
		mapProcessorUtility:      mapProcessorUtility,
	}, nil
}

func (agreedProdLoader *agreedProductLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	product *model.ObjectIDOnly,
	agreedProduct *model.InternalAgreedProductInput,
) (bool, error) {
	agreedProductOutput := model.InternalAgreedProductInput{}
	if agreedProduct != nil {
		agreedProductOutput = *agreedProduct
	}
	existingProduct, err := agreedProdLoader.productDataSource.GetMongoDataSource().FindByID(
		*product.ID,
		session,
	)
	if err != nil {
		return false, horeekaacoreexceptiontofailure.ConvertException(
			"/agreedProductLoader",
			err,
		)
	}

	existingProductJson, _ := json.Marshal(existingProduct)
	json.Unmarshal(existingProductJson, &agreedProductOutput)

	for i := 0; i < len(agreedProductOutput.Variants); i++ {
		loadedVariant, err := agreedProdLoader.productVariantDataSource.GetMongoDataSource().FindByID(
			agreedProductOutput.Variants[i].ID,
			session,
		)
		if err != nil {
			return false, horeekaacoreexceptiontofailure.ConvertException(
				"/agreedProductLoader",
				err,
			)
		}

		loadedVariantJson, _ := json.Marshal(loadedVariant)
		json.Unmarshal(loadedVariantJson, &agreedProductOutput.Variants[i])
	}

	if agreedProduct != nil {
		for _, variant := range agreedProduct.Variants {
			index := funk.IndexOf(
				existingProduct.Variants,
				func(pv *model.ProductVariant) bool {
					return pv.ID == variant.ID
				},
			)
			if index < 0 {
				continue
			}

			var agreedProductVariantMap map[string]interface{}
			agreedProductVariantJson, _ := json.Marshal(variant)
			json.Unmarshal(agreedProductVariantJson, &agreedProductVariantMap)

			agreedProdLoader.mapProcessorUtility.RemoveNil(agreedProductVariantMap)

			agreedProductVariantJson, _ = json.Marshal(agreedProductVariantMap)
			json.Unmarshal(agreedProductVariantJson, &agreedProductOutput.Variants[index])
		}
	}

	agreedProduct = &agreedProductOutput
	return true, nil
}
