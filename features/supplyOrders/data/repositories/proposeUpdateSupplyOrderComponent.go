package supplyorderdomainrepositories

import (
	"encoding/json"
	"time"

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
	pathIdentity              string
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
		pathIdentity:              "ProposeUpdateSupplyOrderComponent",
	}, nil
}

func (updateSupplyOrderTrx *proposeUpdateSupplyOrderTransactionComponent) PreTransaction(
	input *model.InternalUpdateSupplyOrder,
) (*model.InternalUpdateSupplyOrder, error) {
	return input, nil
}

func (updateSupplyOrderTrx *proposeUpdateSupplyOrderTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateSupplyOrder,
) (*model.SupplyOrder, error) {
	updateSupplyOrder := &model.DatabaseUpdateSupplyOrder{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateSupplyOrder)

	existingSupplyOrder, err := updateSupplyOrderTrx.supplyOrderDataSource.GetMongoDataSource().FindByID(
		updateSupplyOrder.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateSupplyOrderTrx.pathIdentity,
			err,
		)
	}

	_, err = updateSupplyOrderTrx.supplyOrderDataLoader.TransactionBody(
		session,
		updateSupplyOrder.Organization,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateSupplyOrderTrx.pathIdentity,
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
			updateSupplyOrderTrx.pathIdentity,
			err,
		)
	}

	totalPrice := 0
	totalReturn := 0
	for _, item := range supplyOrderItems {
		if !item.PartnerAgreed {
			continue
		}
		totalPrice += item.SubTotal
		if item.SupplyOrderItemReturn != nil {
			totalReturn += item.SupplyOrderItemReturn.SubTotal
		}
	}
	updateSupplyOrder.Total = &totalPrice
	updateSupplyOrder.TotalReturn = &totalReturn
	totalSales := totalPrice - totalReturn

	updateSupplyOrder.FinalSalesAmount = &totalSales

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
			updateSupplyOrderTrx.pathIdentity,
			err,
		)
	}
	updateSupplyOrder.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	currentTime := time.Now()
	updateSupplyOrder.UpdatedAt = &currentTime

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
			updateSupplyOrderTrx.pathIdentity,
			err,
		)
	}

	return updatedsupplyOrder, nil
}
