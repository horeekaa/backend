package supplyorderitemdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createSupplyOrderItemTransactionComponent struct {
	supplyOrderItemDataSource       databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource
	loggingDataSource               databaseloggingdatasourceinterfaces.LoggingDataSource
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	supplyOrderItemLoader           supplyorderitemdomainrepositoryutilityinterfaces.SupplyOrderItemLoader
	generatedObjectID               *primitive.ObjectID
	pathIdentity                    string
}

func NewCreateSupplyOrderItemTransactionComponent(
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	supplyOrderItemLoader supplyorderitemdomainrepositoryutilityinterfaces.SupplyOrderItemLoader,
) (supplyorderitemdomainrepositoryinterfaces.CreateSupplyOrderItemTransactionComponent, error) {
	return &createSupplyOrderItemTransactionComponent{
		supplyOrderItemDataSource:       supplyOrderItemDataSource,
		loggingDataSource:               loggingDataSource,
		createDescriptivePhotoComponent: createDescriptivePhotoComponent,
		supplyOrderItemLoader:           supplyOrderItemLoader,
		pathIdentity:                    "CreateSupplyOrderItemComponent",
	}, nil
}

func (createSupplyOrderItemTrx *createSupplyOrderItemTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().GenerateObjectID()
	createSupplyOrderItemTrx.generatedObjectID = &generatedObjectID
	return *createSupplyOrderItemTrx.generatedObjectID
}

func (createSupplyOrderItemTrx *createSupplyOrderItemTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createSupplyOrderItemTrx.generatedObjectID == nil {
		generatedObjectID := createSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().GenerateObjectID()
		createSupplyOrderItemTrx.generatedObjectID = &generatedObjectID
	}
	return *createSupplyOrderItemTrx.generatedObjectID
}

func (createSupplyOrderItemTrx *createSupplyOrderItemTransactionComponent) PreTransaction(
	createsupplyOrderItemcreatesupplyOrderItem *model.InternalCreateSupplyOrderItem,
) (*model.InternalCreateSupplyOrderItem, error) {
	return createsupplyOrderItemcreatesupplyOrderItem, nil
}

func (createSupplyOrderItemTrx *createSupplyOrderItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateSupplyOrderItem,
) (*model.SupplyOrderItem, error) {
	supplyOrderItemToCreate := &model.DatabaseCreateSupplyOrderItem{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, supplyOrderItemToCreate)

	generatedObjectID := createSupplyOrderItemTrx.GetCurrentObjectID()
	for i, photo := range input.Photos {
		photoToCreate := &model.InternalCreateDescriptivePhoto{}

		jsonTemp, _ := json.Marshal(photo)
		json.Unmarshal(jsonTemp, &photoToCreate)
		photoToCreate.Photo.File = photo.Photo.File
		photoToCreate.Category = model.DescriptivePhotoCategorySupplyOrderItemOnPickup
		photoToCreate.Object = &model.ObjectIDOnly{
			ID: &generatedObjectID,
		}
		photoToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
			return &s
		}(*photoToCreate.ProposalStatus)
		photoToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
			return &m
		}(*photoToCreate.SubmittingAccount)
		descriptivePhoto, err := createSupplyOrderItemTrx.createDescriptivePhotoComponent.TransactionBody(
			&mongodbcoretypes.OperationOptions{},
			photoToCreate,
		)
		if err != nil {
			return nil, err
		}

		jsonTemp, _ = json.Marshal(descriptivePhoto)
		json.Unmarshal(jsonTemp, &supplyOrderItemToCreate.Photos[i])
	}
	_, err := createSupplyOrderItemTrx.supplyOrderItemLoader.TransactionBody(
		session,
		supplyOrderItemToCreate.PurchaseOrderToSupply,
		supplyOrderItemToCreate.PickUpDetail,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createSupplyOrderItemTrx.pathIdentity,
			err,
		)
	}

	quantity := supplyOrderItemToCreate.QuantityOffered
	if supplyOrderItemToCreate.QuantityAccepted > 0 {
		quantity = supplyOrderItemToCreate.QuantityAccepted
	}
	supplyOrderItemToCreate.SubTotal = quantity * supplyOrderItemToCreate.UnitPrice
	supplyOrderItemToCreate.SalesAmount = supplyOrderItemToCreate.SubTotal

	newDocumentJson, _ := json.Marshal(*supplyOrderItemToCreate)
	loggingOutput, err := createSupplyOrderItemTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "SupplyOrderItem",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: supplyOrderItemToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *supplyOrderItemToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createSupplyOrderItemTrx.pathIdentity,
			err,
		)
	}

	supplyOrderItemToCreate.ID = generatedObjectID
	supplyOrderItemToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *supplyOrderItemToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		supplyOrderItemToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: supplyOrderItemToCreate.SubmittingAccount.ID}
	}

	currentTime := time.Now()
	supplyOrderItemToCreate.CreatedAt = &currentTime
	supplyOrderItemToCreate.UpdatedAt = &currentTime

	defaultProposalStatus := model.EntityProposalStatusProposed
	if supplyOrderItemToCreate.ProposalStatus == nil {
		supplyOrderItemToCreate.ProposalStatus = &defaultProposalStatus
	}

	defaultSupplyOrderItemStatus := model.SupplyOrderItemStatusAwaitingAcceptance
	if supplyOrderItemToCreate.Status == nil {
		supplyOrderItemToCreate.Status = &defaultSupplyOrderItemStatus
	}

	jsonTemp, _ = json.Marshal(supplyOrderItemToCreate)
	json.Unmarshal(jsonTemp, &supplyOrderItemToCreate.ProposedChanges)

	createdsupplyOrderItem, err := createSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().Create(
		supplyOrderItemToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createSupplyOrderItemTrx.pathIdentity,
			err,
		)
	}
	createSupplyOrderItemTrx.generatedObjectID = nil

	return createdsupplyOrderItem, nil
}
