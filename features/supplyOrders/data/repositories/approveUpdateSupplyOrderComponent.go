package supplyorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type approveUpdateSupplyOrderTransactionComponent struct {
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource
	loggingDataSource     databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility   coreutilityinterfaces.MapProcessorUtility
	supplyOrderDataLoader supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader
}

func NewApproveUpdateSupplyOrderTransactionComponent(
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
	supplyOrderDataLoader supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader,
) (supplyorderdomainrepositoryinterfaces.ApproveUpdateSupplyOrderTransactionComponent, error) {
	return &approveUpdateSupplyOrderTransactionComponent{
		supplyOrderDataSource: supplyOrderDataSource,
		loggingDataSource:     loggingDataSource,
		mapProcessorUtility:   mapProcessorUtility,
		supplyOrderDataLoader: supplyOrderDataLoader,
	}, nil
}

func (approvesupplyOrderTrx *approveUpdateSupplyOrderTransactionComponent) PreTransaction(
	input *model.InternalUpdateSupplyOrder,
) (*model.InternalUpdateSupplyOrder, error) {
	return input, nil
}

func (approvesupplyOrderTrx *approveUpdateSupplyOrderTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateSupplyOrder,
) (*model.SupplyOrder, error) {
	updateSupplyOrder := &model.DatabaseUpdateSupplyOrder{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateSupplyOrder)

	existingSupplyOrder, err := approvesupplyOrderTrx.supplyOrderDataSource.GetMongoDataSource().FindByID(
		updateSupplyOrder.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateSupplyOrder",
			err,
		)
	}

	_, err = approvesupplyOrderTrx.supplyOrderDataLoader.TransactionBody(
		session,
		updateSupplyOrder.Organization,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateSupplyOrder",
			err,
		)
	}

	previousLog, err := approvesupplyOrderTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingSupplyOrder.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateSupplyOrder",
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateSupplyOrder.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateSupplyOrder.ProposalStatus,
	}
	jsonTemp, _ = json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approvesupplyOrderTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateSupplyOrder",
			err,
		)
	}

	updateSupplyOrder.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdatesupplyOrder := &model.DatabaseUpdateSupplyOrder{
		ID: updateSupplyOrder.ID,
	}
	jsonExisting, _ := json.Marshal(existingSupplyOrder.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdatesupplyOrder.ProposedChanges)

	var updatesupplyOrderMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateSupplyOrder)
	json.Unmarshal(jsonUpdate, &updatesupplyOrderMap)

	approvesupplyOrderTrx.mapProcessorUtility.RemoveNil(updatesupplyOrderMap)

	jsonUpdate, _ = json.Marshal(updatesupplyOrderMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdatesupplyOrder.ProposedChanges)

	if updateSupplyOrder.ProposalStatus != nil {
		if *updateSupplyOrder.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdatesupplyOrder.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdatesupplyOrder)
		}
	}

	updatedsupplyOrder, err := approvesupplyOrderTrx.supplyOrderDataSource.GetMongoDataSource().Update(
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
