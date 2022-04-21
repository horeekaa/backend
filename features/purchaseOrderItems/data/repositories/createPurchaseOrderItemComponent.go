package purchaseorderitemdomainrepositories

import (
	"encoding/json"
	"time"

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
	pathIdentity                    string
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
		pathIdentity:                    "CreatePurchaseOrderItemComponent",
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
	input *model.InternalCreatePurchaseOrderItem,
) (*model.PurchaseOrderItem, error) {
	purchaseOrderItemToCreate := &model.DatabaseCreatePurchaseOrderItem{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, purchaseOrderItemToCreate)

	generatedObjectID := createPurchaseOrderItemTrx.GetCurrentObjectID()
	_, err := createPurchaseOrderItemTrx.purchaseOrderItemLoader.TransactionBody(
		session,
		purchaseOrderItemToCreate.MouItem,
		purchaseOrderItemToCreate.ProductVariant,
		purchaseOrderItemToCreate.DeliveryDetail,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createPurchaseOrderItemTrx.pathIdentity,
			err,
		)
	}

	purchaseOrderItemToCreate.UnitPrice = purchaseOrderItemToCreate.ProductVariant.RetailPrice
	if purchaseOrderItemToCreate.MouItem != nil {
		index := funk.IndexOf(
			purchaseOrderItemToCreate.MouItem.AgreedProduct.Variants,
			func(pv *model.InternalAgreedProductVariantInput) bool {
				return pv.ID == purchaseOrderItemToCreate.ProductVariant.ID
			},
		)
		if index > -1 {
			purchaseOrderItemToCreate.UnitPrice = *purchaseOrderItemToCreate.MouItem.AgreedProduct.Variants[index].RetailPrice
		}
	}
	purchaseOrderItemToCreate.SubTotal = purchaseOrderItemToCreate.Quantity * purchaseOrderItemToCreate.UnitPrice
	purchaseOrderItemToCreate.SalesAmount = purchaseOrderItemToCreate.SubTotal

	newDocumentJson, _ := json.Marshal(*purchaseOrderItemToCreate)
	loggingOutput, err := createPurchaseOrderItemTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "PurchaseOrderItem",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: purchaseOrderItemToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *purchaseOrderItemToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createPurchaseOrderItemTrx.pathIdentity,
			err,
		)
	}

	purchaseOrderItemToCreate.ID = generatedObjectID
	purchaseOrderItemToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *purchaseOrderItemToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		purchaseOrderItemToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: purchaseOrderItemToCreate.SubmittingAccount.ID}
	}

	currentTime := time.Now().UTC()
	purchaseOrderItemToCreate.CreatedAt = &currentTime
	purchaseOrderItemToCreate.UpdatedAt = &currentTime

	defaultStatus := model.PurchaseOrderItemStatusPendingConfirmation
	if purchaseOrderItemToCreate.Status == nil {
		purchaseOrderItemToCreate.Status = &defaultStatus
	}

	defaultProposalStatus := model.EntityProposalStatusProposed
	if purchaseOrderItemToCreate.ProposalStatus == nil {
		purchaseOrderItemToCreate.ProposalStatus = &defaultProposalStatus
	}

	jsonTemp, _ = json.Marshal(purchaseOrderItemToCreate)
	json.Unmarshal(jsonTemp, &purchaseOrderItemToCreate.ProposedChanges)

	createdPurchaseOrderItem, err := createPurchaseOrderItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().Create(
		purchaseOrderItemToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createPurchaseOrderItemTrx.pathIdentity,
			err,
		)
	}
	createPurchaseOrderItemTrx.generatedObjectID = nil

	return createdPurchaseOrderItem, nil
}
