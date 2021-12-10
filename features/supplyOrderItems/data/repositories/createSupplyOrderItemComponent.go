package supplyorderitemdomainrepositories

import (
	"encoding/json"

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
	createSupplyOrderItem *model.InternalCreateSupplyOrderItem,
) (*model.SupplyOrderItem, error) {
	generatedObjectID := createSupplyOrderItemTrx.GetCurrentObjectID()
	for i, photo := range createSupplyOrderItem.Photos {
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
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/createSupplyOrderItemComponent",
				err,
			)
		}

		jsonTemp, _ = json.Marshal(descriptivePhoto)
		json.Unmarshal(jsonTemp, &createSupplyOrderItem.Photos[i])
	}
	for i, photo := range createSupplyOrderItem.PickUpDetail.Photos {
		photoToCreate := &model.InternalCreateDescriptivePhoto{}

		jsonTemp, _ := json.Marshal(photo)
		json.Unmarshal(jsonTemp, &photoToCreate)
		photoToCreate.Photo.File = photo.Photo.File
		photoToCreate.Category = model.DescriptivePhotoCategorySupplyOrderItem
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
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/createSupplyOrderItemComponent",
				err,
			)
		}

		jsonTemp, _ = json.Marshal(descriptivePhoto)
		json.Unmarshal(jsonTemp, &createSupplyOrderItem.PickUpDetail.Photos[i])
	}
	_, err := createSupplyOrderItemTrx.supplyOrderItemLoader.TransactionBody(
		session,
		createSupplyOrderItem.PurchaseOrderToSupply,
		createSupplyOrderItem.PickUpDetail,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createSupplyOrderItemComponent",
			err,
		)
	}

	quantity := createSupplyOrderItem.QuantityOffered
	if createSupplyOrderItem.QuantityAccepted != nil {
		quantity = *createSupplyOrderItem.QuantityAccepted
	}
	createSupplyOrderItem.SubTotal = quantity * createSupplyOrderItem.UnitPrice

	newDocumentJson, _ := json.Marshal(*createSupplyOrderItem)
	loggingOutput, err := createSupplyOrderItemTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "SupplyOrderItem",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: createSupplyOrderItem.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *createSupplyOrderItem.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createSupplyOrderItemComponent",
			err,
		)
	}

	createSupplyOrderItem.ID = &generatedObjectID
	createSupplyOrderItem.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *createSupplyOrderItem.ProposalStatus == model.EntityProposalStatusApproved {
		createSupplyOrderItem.RecentApprovingAccount = &model.ObjectIDOnly{ID: createSupplyOrderItem.SubmittingAccount.ID}
	}

	supplyOrderItemToCreate := &model.DatabaseCreateSupplyOrderItem{}
	jsonTemp, _ := json.Marshal(createSupplyOrderItem)
	json.Unmarshal(jsonTemp, supplyOrderItemToCreate)
	json.Unmarshal(jsonTemp, &supplyOrderItemToCreate.ProposedChanges)
	supplyOrderItemToCreate.QuantityAccepted = 0

	createdsupplyOrderItem, err := createSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().Create(
		supplyOrderItemToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createSupplyOrderItemComponent",
			err,
		)
	}
	createSupplyOrderItemTrx.generatedObjectID = nil

	return createdsupplyOrderItem, nil
}
