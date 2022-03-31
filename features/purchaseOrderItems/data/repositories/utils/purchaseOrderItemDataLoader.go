package purchaseorderitemdomainrepositoryutilities

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/utils"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
)

type purchaseOrderItemLoader struct {
	accountDataSource          databaseaccountdatasourceinterfaces.AccountDataSource
	personDataSource           databaseaccountdatasourceinterfaces.PersonDataSource
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource
	mouItemDataSource          databasemouitemdatasourceinterfaces.MouItemDataSource
	productVariantDataSource   databaseproductvariantdatasourceinterfaces.ProductVariantDataSource
	productDataSource          databaseproductdatasourceinterfaces.ProductDataSource
	tagDataSource              databasetagdatasourceinterfaces.TagDataSource
	taggingDataSource          databasetaggingdatasourceinterfaces.TaggingDataSource
	addressDataSource          databaseaddressdatasourceinterfaces.AddressDataSource
}

func NewPurchaseOrderItemLoader(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
	descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
	mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
	productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
) (purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader, error) {
	return &purchaseOrderItemLoader{
		accountDataSource,
		personDataSource,
		descriptivePhotoDataSource,
		mouItemDataSource,
		productVariantDataSource,
		productDataSource,
		tagDataSource,
		taggingDataSource,
		addressDataSource,
	}, nil
}

func (purcOrderItemLoader *purchaseOrderItemLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	mouItem *model.MouItemForPurchaseOrderItemInput,
	productVariant *model.ProductVariantForPurchaseOrderItemInput,
	delivery *model.DatabasePurchaseOrderItemDelivery,
) (bool, error) {
	mouItemLoadedChan := make(chan bool)
	productVariantLoadedChan := make(chan bool)
	deliveryLoadedChan := make(chan bool)
	errChan := make(chan error)

	go func() {
		if delivery == nil {
			deliveryLoadedChan <- true
			return
		}
		if delivery.ExpectedArrivalDate != nil {
			loc, _ := time.LoadLocation("Asia/Bangkok")
			*delivery.ExpectedArrivalDate = time.Date(
				delivery.ExpectedArrivalDate.Year(),
				delivery.ExpectedArrivalDate.Month(),
				delivery.ExpectedArrivalDate.Day(),
				0, 0, 0, 0,
				loc,
			)
		}

		addressLoadedChan := make(chan bool)
		accountLoadedChan := make(chan bool)

		go func() {
			if delivery.Address == nil {
				addressLoadedChan <- true
				return
			}
			loadedAddress, err := purcOrderItemLoader.addressDataSource.GetMongoDataSource().FindByID(
				delivery.Address.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}
			jsonTemp, _ := json.Marshal(
				map[string]interface{}{
					"Address":  loadedAddress,
					"Status":   model.DeliveryStatusAddressNoted,
					"Courier":  model.AccountForPurchaseOrderItemInput{},
					"PublicID": "",
				},
			)
			json.Unmarshal(jsonTemp, delivery)

			addressLoadedChan <- true
		}()

		go func() {
			if delivery.Courier == nil {
				accountLoadedChan <- true
				return
			}
			account, err := purcOrderItemLoader.accountDataSource.GetMongoDataSource().FindByID(
				*delivery.Courier.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}

			jsonTemp, _ := json.Marshal(account)
			json.Unmarshal(jsonTemp, &delivery.Courier)

			person, err := purcOrderItemLoader.personDataSource.GetMongoDataSource().FindByID(
				delivery.Courier.Person.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}

			jsonTemp, _ = json.Marshal(person)
			json.Unmarshal(jsonTemp, &delivery.Courier.Person)

			accountLoadedChan <- true
		}()

		for i := 0; i < 2; {
			select {
			case _ = <-addressLoadedChan:
				i++
			case _ = <-accountLoadedChan:
				i++
			}
		}

		deliveryLoadedChan <- true
	}()

	go func() {
		if mouItem == nil {
			mouItemLoadedChan <- true
			return
		}

		loadedMouItem, err := purcOrderItemLoader.mouItemDataSource.GetMongoDataSource().FindByID(
			mouItem.ID,
			session,
		)
		if err != nil {
			errChan <- err
			return
		}

		jsonTemp, _ := json.Marshal(loadedMouItem)
		json.Unmarshal(jsonTemp, mouItem)
		mouItemLoadedChan <- true
	}()

	go func() {
		if productVariant == nil {
			productVariantLoadedChan <- true
			return
		}

		loadedProductVariant, err := purcOrderItemLoader.productVariantDataSource.GetMongoDataSource().FindByID(
			productVariant.ID,
			session,
		)
		if err != nil {
			errChan <- err
			return
		}

		jsonTemp, _ := json.Marshal(loadedProductVariant)
		json.Unmarshal(jsonTemp, productVariant)

		descriptivePhotoLoadedChan := make(chan bool)
		productLoadedChan := make(chan bool)

		go func() {
			if productVariant.Photo == nil {
				descriptivePhotoLoadedChan <- true
				return
			}
			loadedDescriptivePhoto, err := purcOrderItemLoader.descriptivePhotoDataSource.GetMongoDataSource().FindByID(
				productVariant.Photo.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}

			jsonTemp, _ := json.Marshal(loadedDescriptivePhoto)
			json.Unmarshal(jsonTemp, &productVariant.Photo)
			descriptivePhotoLoadedChan <- true
		}()

		go func() {
			loadedProduct, err := purcOrderItemLoader.productDataSource.GetMongoDataSource().FindByID(
				productVariant.Product.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}

			jsonTemp, _ := json.Marshal(loadedProduct)
			json.Unmarshal(jsonTemp, &productVariant.Product)

			prodDescriptivePhotoLoadedChan := make(chan bool)
			prodTaggingsLoadedChan := make(chan bool)

			go func() {
				for i := 0; i < len(loadedProduct.Photos); i++ {
					loadedDescriptivePhoto, err := purcOrderItemLoader.descriptivePhotoDataSource.GetMongoDataSource().FindByID(
						productVariant.Product.Photos[i].ID,
						session,
					)
					if err != nil {
						errChan <- err
						return
					}

					jsonTemp, _ := json.Marshal(loadedDescriptivePhoto)
					json.Unmarshal(jsonTemp, &productVariant.Product.Photos[i])
				}
				prodDescriptivePhotoLoadedChan <- true
			}()

			go func() {
				for i := 0; i < len(loadedProduct.Taggings); i++ {
					loadedTagging, err := purcOrderItemLoader.taggingDataSource.GetMongoDataSource().FindByID(
						productVariant.Product.Taggings[i].ID,
						session,
					)
					if err != nil {
						errChan <- err
						return
					}

					jsonTemp, _ := json.Marshal(loadedTagging)
					json.Unmarshal(jsonTemp, &productVariant.Product.Taggings[i])

					loadedTag, err := purcOrderItemLoader.tagDataSource.GetMongoDataSource().FindByID(
						productVariant.Product.Taggings[i].Tag.ID,
						session,
					)
					jsonTemp, _ = json.Marshal(loadedTag)
					json.Unmarshal(jsonTemp, &productVariant.Product.Taggings[i].Tag)
				}
				prodTaggingsLoadedChan <- true
			}()

			for i := 0; i < 2; {
				select {
				case _ = <-prodDescriptivePhotoLoadedChan:
					i++
				case _ = <-prodTaggingsLoadedChan:
					i++
				}
			}
			productLoadedChan <- true
		}()

		for i := 0; i < 2; {
			select {
			case _ = <-descriptivePhotoLoadedChan:
				i++
			case _ = <-productLoadedChan:
				i++
			}
		}

		productVariantLoadedChan <- true
	}()

	for i := 0; i < 3; {
		select {
		case err := <-errChan:
			return false, err
		case _ = <-mouItemLoadedChan:
			i++
		case _ = <-productVariantLoadedChan:
			i++
		case _ = <-deliveryLoadedChan:
			i++
		}
	}
	return true, nil
}
