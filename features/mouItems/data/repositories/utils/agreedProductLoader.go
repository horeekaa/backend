package mouitemdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories/utils"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type agreedProductLoader struct {
	productVariantDataSource   databaseproductvariantdatasourceinterfaces.ProductVariantDataSource
	productDataSource          databaseproductdatasourceinterfaces.ProductDataSource
	organizationDataSource     databaseorganizationdatasourceinterfaces.OrganizationDataSource
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource
	mapProcessorUtility        coreutilityinterfaces.MapProcessorUtility
	pathIdentity               string
}

func NewAgreedProductLoader(
	productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader, error) {
	return &agreedProductLoader{
		productVariantDataSource:   productVariantDataSource,
		productDataSource:          productDataSource,
		organizationDataSource:     organizationDataSource,
		descriptivePhotoDataSource: descriptivePhotoDataSource,
		mapProcessorUtility:        mapProcessorUtility,
		pathIdentity:               "ProductForMOULoader",
	}, nil
}

func (agreedProdLoader *agreedProductLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	product *model.ObjectIDOnly,
	agreedProduct *model.InternalAgreedProductInput,
	organization *model.OrganizationForMouItemInput,
) (bool, error) {
	agreedProductOutput := model.InternalAgreedProductInput{}
	existingProduct, err := agreedProdLoader.productDataSource.GetMongoDataSource().FindByID(
		*product.ID,
		session,
	)
	if err != nil {
		return false, horeekaacoreexceptiontofailure.ConvertException(
			agreedProdLoader.pathIdentity,
			err,
		)
	}

	existingProductJson, _ := json.Marshal(existingProduct)
	json.Unmarshal(existingProductJson, &agreedProductOutput)

	if agreedProduct != nil {
		var agreedProductMap map[string]interface{}
		agreedProductUpdateJson, _ := json.Marshal(agreedProduct)
		json.Unmarshal(agreedProductUpdateJson, &agreedProductMap)

		agreedProdLoader.mapProcessorUtility.RemoveNil(agreedProductMap)
		delete(agreedProductMap, "variants")

		agreedProductUpdateJson, _ = json.Marshal(agreedProductMap)
		json.Unmarshal(agreedProductUpdateJson, &agreedProductOutput)
	}

	descriptivePhotoLoadedChan := make(chan bool)
	variantsLoadedChan := make(chan bool)
	organizationLoadedChan := make(chan bool)
	errChan := make(chan error)

	go func() {
		for i := 0; i < len(agreedProductOutput.Photos); i++ {
			loadedDescriptivePhoto, err := agreedProdLoader.descriptivePhotoDataSource.GetMongoDataSource().FindByID(
				*agreedProductOutput.Photos[i].ID,
				session,
			)
			if err != nil {
				errChan <- horeekaacoreexceptiontofailure.ConvertException(
					agreedProdLoader.pathIdentity,
					err,
				)
			}

			descriptivePhotoJson, _ := json.Marshal(loadedDescriptivePhoto)
			json.Unmarshal(descriptivePhotoJson, &agreedProductOutput.Photos[i])
		}
		descriptivePhotoLoadedChan <- true
	}()

	go func() {
		for i := 0; i < len(agreedProductOutput.Variants); i++ {
			loadedVariant, err := agreedProdLoader.productVariantDataSource.GetMongoDataSource().FindByID(
				agreedProductOutput.Variants[i].ID,
				session,
			)
			if err != nil {
				errChan <- horeekaacoreexceptiontofailure.ConvertException(
					agreedProdLoader.pathIdentity,
					err,
				)
			}

			if loadedVariant.Photo != nil {
				loadedDescriptivePhoto, err := agreedProdLoader.descriptivePhotoDataSource.GetMongoDataSource().FindByID(
					loadedVariant.Photo.ID,
					session,
				)
				if err != nil {
					errChan <- horeekaacoreexceptiontofailure.ConvertException(
						agreedProdLoader.pathIdentity,
						err,
					)
				}

				loadedDescriptivePhotoJson, _ := json.Marshal(loadedDescriptivePhoto)
				json.Unmarshal(loadedDescriptivePhotoJson, &loadedVariant.Photo)
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

		variantsLoadedChan <- true
	}()

	go func() {
		if organization == nil {
			organizationLoadedChan <- true
			return
		}

		loadedOrganization, err := agreedProdLoader.organizationDataSource.GetMongoDataSource().FindByID(
			organization.ID,
			session,
		)
		if err != nil {
			errChan <- err
			return
		}

		orgJson, _ := json.Marshal(loadedOrganization)
		json.Unmarshal(orgJson, organization)

		organizationLoadedChan <- true
	}()

	for i := 0; i < 3; {
		select {
		case err := <-errChan:
			return false, err
		case _ = <-variantsLoadedChan:
			i++
		case _ = <-descriptivePhotoLoadedChan:
			i++
		case _ = <-organizationLoadedChan:
			i++
		}
	}

	*agreedProduct = agreedProductOutput
	return true, nil
}
