package purchaseorderitemdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createPurchaseOrderItemTransactionComponent struct {
	purchaseOrderItemDataSource     databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource
	loggingDataSource               databaseloggingdatasourceinterfaces.LoggingDataSource
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	purchaseOrderItemLoader         purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader
	generatedObjectID               *primitive.ObjectID
}

func NewCreatePurchaseOrderItemTransactionComponent(
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	purchaseOrderItemLoader purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader,
) (purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent, error) {
	return &createPurchaseOrderItemTransactionComponent{
		purchaseOrderItemDataSource:     purchaseOrderItemDataSource,
		loggingDataSource:               loggingDataSource,
		createDescriptivePhotoComponent: createDescriptivePhotoComponent,
		purchaseOrderItemLoader:         purchaseOrderItemLoader,
	}, nil
}

func (createPurchaseOrderItemTrx *createPurchaseOrderItemTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createPurchaseOrderItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().GenerateObjectID()
	createPurchaseOrderItemTrx.generatedObjectID = &generatedObjectID
	return *createPurchaseOrderItemTrx.generatedObjectID
}

func (createPurchaseOrderItemTrx *createPurchaseOrderItemTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createPurchaseOrderItemTrx.generatedObjectID == nil {
		generatedObjectID := createPurchaseOrderItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().GenerateObjectID()
		createPurchaseOrderItemTrx.generatedObjectID = &generatedObjectID
	}
	return *createPurchaseOrderItemTrx.generatedObjectID
}

func (createPurchaseOrderItemTrx *createPurchaseOrderItemTransactionComponent) PreTransaction(
	createPurchaseOrderItemcreatePurchaseOrderItem *model.InternalCreatePurchaseOrderItem,
) (*model.InternalCreatePurchaseOrderItem, error) {
	return createPurchaseOrderItemcreatePurchaseOrderItem, nil
}

func (createPurchaseOrderItemTrx *createPurchaseOrderItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	createPurchaseOrderItem *model.InternalCreatePurchaseOrderItem,
) (*model.PurchaseOrderItem, error) {
	generatedObjectID := createPurchaseOrderItemTrx.GetCurrentObjectID()
	for i, photoToCreate := range createPurchaseOrderItem.DeliveryDetail.Photos {
		photoToCreate.Category = model.DescriptivePhotoCategoryPurchaseOrderItem
		photoToCreate.Object = &model.ObjectIDOnly{
			ID: &generatedObjectID,
		}
		photoToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
			return &s
		}(*photoToCreate.ProposalStatus)
		photoToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
			return &m
		}(*photoToCreate.SubmittingAccount)
		descriptivePhoto, err := createPurchaseOrderItemTrx.createDescriptivePhotoComponent.TransactionBody(
			&mongodbcoretypes.OperationOptions{},
			photoToCreate,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/createPurchaseOrderItemComponent",
				err,
			)
		}

		jsonTemp, _ := json.Marshal(descriptivePhoto)
		json.Unmarshal(jsonTemp, &createPurchaseOrderItem.DeliveryDetail.Photos[i])
	}
	_, err := createPurchaseOrderItemTrx.purchaseOrderItemLoader.TransactionBody(
		session,
		createPurchaseOrderItem.MouItem,
		createPurchaseOrderItem.ProductVariant,
		createPurchaseOrderItem.DeliveryDetail.Address,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createPurchaseOrderItemComponent",
			err,
		)
	}

	createPurchaseOrderItem.UnitPrice = createPurchaseOrderItem.ProductVariant.RetailPrice
	if createPurchaseOrderItem.MouItem != nil {
		index := funk.IndexOf(
			createPurchaseOrderItem.MouItem.AgreedProduct.Variants,
			func(pv *model.InternalAgreedProductVariantInput) bool {
				return pv.ID == createPurchaseOrderItem.ProductVariant.ID
			},
		)
		if index > -1 {
			createPurchaseOrderItem.UnitPrice = *createPurchaseOrderItem.MouItem.AgreedProduct.Variants[index].RetailPrice
		}
	}
	createPurchaseOrderItem.SubTotal = func(i int) *int { return &i }(createPurchaseOrderItem.Quantity * createPurchaseOrderItem.UnitPrice)

	newDocumentJson, _ := json.Marshal(*createPurchaseOrderItem)
	loggingOutput, err := createPurchaseOrderItemTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "PurchaseOrderItem",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: createPurchaseOrderItem.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *createPurchaseOrderItem.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createPurchaseOrderItemComponent",
			err,
		)
	}

	createPurchaseOrderItem.ID = &generatedObjectID
	createPurchaseOrderItem.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *createPurchaseOrderItem.ProposalStatus == model.EntityProposalStatusApproved {
		createPurchaseOrderItem.RecentApprovingAccount = &model.ObjectIDOnly{ID: createPurchaseOrderItem.SubmittingAccount.ID}
	}

	purchaseOrderItemToCreate := &model.DatabaseCreatePurchaseOrderItem{}
	jsonTemp, _ := json.Marshal(createPurchaseOrderItem)
	json.Unmarshal(jsonTemp, purchaseOrderItemToCreate)
	json.Unmarshal(jsonTemp, &purchaseOrderItemToCreate.ProposedChanges)

	createdPurchaseOrderItem, err := createPurchaseOrderItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().Create(
		purchaseOrderItemToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createPurchaseOrderItem",
			err,
		)
	}
	createPurchaseOrderItemTrx.generatedObjectID = nil

	return createdPurchaseOrderItem, nil
}
