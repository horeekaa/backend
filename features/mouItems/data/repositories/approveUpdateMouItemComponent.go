package mouitemdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateMouItemTransactionComponent struct {
	mouItemDataSource   databasemouitemdatasourceinterfaces.MouItemDataSource
	loggingDataSource   databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility
	pathIdentity        string
}

func NewApproveUpdateMouItemTransactionComponent(
	mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (mouitemdomainrepositoryinterfaces.ApproveUpdateMouItemTransactionComponent, error) {
	return &approveUpdateMouItemTransactionComponent{
		mouItemDataSource:   mouItemDataSource,
		loggingDataSource:   loggingDataSource,
		mapProcessorUtility: mapProcessorUtility,
		pathIdentity:        "ApproveUpdateMouItemComponent",
	}, nil
}

func (approveUpdateMouItemTrx *approveUpdateMouItemTransactionComponent) PreTransaction(
	input *model.InternalUpdateMouItem,
) (*model.InternalUpdateMouItem, error) {
	return input, nil
}

func (approveUpdateMouItemTrx *approveUpdateMouItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateMouItem,
) (*model.MouItem, error) {
	updateMouItem := &model.DatabaseUpdateMouItem{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateMouItem)

	existingMouItem, err := approveUpdateMouItemTrx.mouItemDataSource.GetMongoDataSource().FindByID(
		updateMouItem.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveUpdateMouItemTrx.pathIdentity,
			err,
		)
	}
	if existingMouItem.ProposedChanges.ProposalStatus == model.EntityProposalStatusApproved {
		return existingMouItem, nil
	}

	previousLog, err := approveUpdateMouItemTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingMouItem.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveUpdateMouItemTrx.pathIdentity,
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateMouItem.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateMouItem.ProposalStatus,
	}
	jsonTemp, _ = json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveUpdateMouItemTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveUpdateMouItemTrx.pathIdentity,
			err,
		)
	}

	updateMouItem.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdateMouItem := &model.DatabaseUpdateMouItem{
		ID: updateMouItem.ID,
	}
	jsonExisting, _ := json.Marshal(existingMouItem.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdateMouItem.ProposedChanges)

	var updateMouItemMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateMouItem)
	json.Unmarshal(jsonUpdate, &updateMouItemMap)

	approveUpdateMouItemTrx.mapProcessorUtility.RemoveNil(updateMouItemMap)

	jsonUpdate, _ = json.Marshal(updateMouItemMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateMouItem.ProposedChanges)

	if updateMouItem.ProposalStatus != nil {
		if *updateMouItem.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdateMouItem.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdateMouItem)
		}
	}

	updatedMouItem, err := approveUpdateMouItemTrx.mouItemDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateMouItem.ID,
		},
		fieldsToUpdateMouItem,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveUpdateMouItemTrx.pathIdentity,
			err,
		)
	}

	return updatedMouItem, nil
}
