package supplyorderitemdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateSupplyOrderItemTransactionComponent struct {
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource
	loggingDataSource         databaseloggingdatasourceinterfaces.LoggingDataSource
	supplyOrderItemLoader     supplyorderitemdomainrepositoryutilityinterfaces.SupplyOrderItemLoader
	mapProcessorUtility       coreutilityinterfaces.MapProcessorUtility
}

func NewProposeUpdateSupplyOrderItemTransactionComponent(
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	supplyOrderItemLoader supplyorderitemdomainrepositoryutilityinterfaces.SupplyOrderItemLoader,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemTransactionComponent, error) {
	return &proposeUpdateSupplyOrderItemTransactionComponent{
		supplyOrderItemDataSource: supplyOrderItemDataSource,
		loggingDataSource:         loggingDataSource,
		supplyOrderItemLoader:     supplyOrderItemLoader,
		mapProcessorUtility:       mapProcessorUtility,
	}, nil
}

func (updateSupplyOrderItemTrx *proposeUpdateSupplyOrderItemTransactionComponent) PreTransaction(
	input *model.InternalUpdateSupplyOrderItem,
) (*model.InternalUpdateSupplyOrderItem, error) {
	return input, nil
}

func (updateSupplyOrderItemTrx *proposeUpdateSupplyOrderItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateSupplyOrderItem *model.InternalUpdateSupplyOrderItem,
) (*model.SupplyOrderItem, error) {
	existingSupplyOrderItem, err := updateSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().FindByID(
		*updateSupplyOrderItem.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateSupplyOrderItemComponent",
			err,
		)
	}

	_, err = updateSupplyOrderItemTrx.supplyOrderItemLoader.TransactionBody(
		session,
		updateSupplyOrderItem.PurchaseOrderToSupply,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateSupplyOrderItemComponent",
			err,
		)
	}

	unitPrice := existingSupplyOrderItem.UnitPrice
	if updateSupplyOrderItem.UnitPrice != nil {
		unitPrice = *updateSupplyOrderItem.UnitPrice
	}

	quantity := existingSupplyOrderItem.QuantityOffered
	if existingSupplyOrderItem.QuantityAccepted > 0 {
		quantity = existingSupplyOrderItem.QuantityAccepted
	}
	if updateSupplyOrderItem.QuantityOffered != nil {
		quantity = *updateSupplyOrderItem.QuantityOffered
	}
	if updateSupplyOrderItem.QuantityAccepted != nil {
		quantity = *updateSupplyOrderItem.QuantityAccepted
	}
	subTotal := quantity * unitPrice
	updateSupplyOrderItem.SubTotal = &subTotal

	newDocumentJson, _ := json.Marshal(*updateSupplyOrderItem)
	oldDocumentJson, _ := json.Marshal(*existingSupplyOrderItem)
	loggingOutput, err := updateSupplyOrderItemTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "SupplyOrderItem",
			Document: &model.ObjectIDOnly{
				ID: &existingSupplyOrderItem.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateSupplyOrderItem.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateSupplyOrderItem.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateSupplyOrderItemComponent",
			err,
		)
	}
	updateSupplyOrderItem.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdatesupplyOrderItem := &model.DatabaseUpdateSupplyOrderItem{
		ID: *updateSupplyOrderItem.ID,
	}
	jsonExisting, _ := json.Marshal(existingSupplyOrderItem)
	json.Unmarshal(jsonExisting, &fieldsToUpdatesupplyOrderItem.ProposedChanges)

	var updatesupplyOrderItemMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateSupplyOrderItem)
	json.Unmarshal(jsonUpdate, &updatesupplyOrderItemMap)

	updateSupplyOrderItemTrx.mapProcessorUtility.RemoveNil(updatesupplyOrderItemMap)

	jsonUpdate, _ = json.Marshal(updatesupplyOrderItemMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdatesupplyOrderItem.ProposedChanges)

	if updateSupplyOrderItem.ProposalStatus != nil {
		fieldsToUpdatesupplyOrderItem.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateSupplyOrderItem.SubmittingAccount.ID,
		}
		if *updateSupplyOrderItem.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdatesupplyOrderItem)
		}
	}

	updatedSupplyOrderItem, err := updateSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdatesupplyOrderItem.ID,
		},
		fieldsToUpdatesupplyOrderItem,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateSupplyOrderItemComponent",
			err,
		)
	}

	return updatedSupplyOrderItem, nil
}
