package supplyorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateSupplyOrderTransactionComponent struct {
	supplyOrderDataSource     databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource
	loggingDataSource         databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility       coreutilityinterfaces.MapProcessorUtility
	supplyOrderDataLoader     supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader
}

func NewProposeUpdateSupplyOrderTransactionComponent(
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
	supplyOrderDataLoader supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader,
) (supplyorderdomainrepositoryinterfaces.ProposeUpdateSupplyOrderTransactionComponent, error) {
	return &proposeUpdateSupplyOrderTransactionComponent{
		supplyOrderDataSource:     supplyOrderDataSource,
		supplyOrderItemDataSource: supplyOrderItemDataSource,
		loggingDataSource:         loggingDataSource,
		mapProcessorUtility:       mapProcessorUtility,
		supplyOrderDataLoader:     supplyOrderDataLoader,
	}, nil
}

func (updateSupplyOrderTrx *proposeUpdateSupplyOrderTransactionComponent) PreTransaction(
	input *model.InternalUpdateSupplyOrder,
) (*model.InternalUpdateSupplyOrder, error) {
	return input, nil
}

func (updateSupplyOrderTrx *proposeUpdateSupplyOrderTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateSupplyOrder *model.InternalUpdateSupplyOrder,
) (*model.SupplyOrder, error) {
	existingSupplyOrder, err := updateSupplyOrderTrx.supplyOrderDataSource.GetMongoDataSource().FindByID(
		updateSupplyOrder.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateSupplyOrder",
			err,
		)
	}

	_, err = updateSupplyOrderTrx.supplyOrderDataLoader.TransactionBody(
		session,
		updateSupplyOrder.Organization,
		updateSupplyOrder.PickUpDetail.Address,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateSupplyOrder",
			err,
		)
	}

	supplyOrderItems, err := updateSupplyOrderTrx.supplyOrderItemDataSource.GetMongoDataSource().Find(
		map[string]interface{}{
			"supplyOrder._id": existingSupplyOrder.ID,
		},
		&mongodbcoretypes.PaginationOptions{},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateSupplyOrder",
			err,
		)
	}

	totalPrice := 0
	for _, item := range supplyOrderItems {
		totalPrice += item.SubTotal
	}
	updateSupplyOrder.Total = &totalPrice

	newDocumentJson, _ := json.Marshal(*updateSupplyOrder)
	oldDocumentJson, _ := json.Marshal(*existingSupplyOrder)
	loggingOutput, err := updateSupplyOrderTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "SupplyOrder",
			Document: &model.ObjectIDOnly{
				ID: &existingSupplyOrder.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateSupplyOrder.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateSupplyOrder.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateSupplyOrder",
			err,
		)
	}
	updateSupplyOrder.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdatesupplyOrder := &model.DatabaseUpdateSupplyOrder{
		ID: updateSupplyOrder.ID,
	}
	jsonExisting, _ := json.Marshal(existingSupplyOrder)
	json.Unmarshal(jsonExisting, &fieldsToUpdatesupplyOrder.ProposedChanges)

	var updatesupplyOrderMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateSupplyOrder)
	json.Unmarshal(jsonUpdate, &updatesupplyOrderMap)

	updateSupplyOrderTrx.mapProcessorUtility.RemoveNil(updatesupplyOrderMap)

	jsonUpdate, _ = json.Marshal(updatesupplyOrderMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdatesupplyOrder.ProposedChanges)
	jsonUpdate, err = json.Marshal(fieldsToUpdatesupplyOrder.ProposedChanges)

	if updateSupplyOrder.ProposalStatus != nil {
		fieldsToUpdatesupplyOrder.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateSupplyOrder.SubmittingAccount.ID,
		}
		if *updateSupplyOrder.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdatesupplyOrder)
		}
	}

	updatedsupplyOrder, err := updateSupplyOrderTrx.supplyOrderDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdatesupplyOrder.ID,
		},
		fieldsToUpdatesupplyOrder,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateSupplyOrder",
			err,
		)
	}

	return updatedsupplyOrder, nil
}
