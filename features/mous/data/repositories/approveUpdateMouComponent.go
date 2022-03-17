package moudomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateMouTransactionComponent struct {
	mouDataSource       databasemoudatasourceinterfaces.MouDataSource
	loggingDataSource   databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility
	pathIdentity        string
}

func NewApproveUpdateMouTransactionComponent(
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (moudomainrepositoryinterfaces.ApproveUpdateMouTransactionComponent, error) {
	return &approveUpdateMouTransactionComponent{
		mouDataSource:       mouDataSource,
		loggingDataSource:   loggingDataSource,
		mapProcessorUtility: mapProcessorUtility,
		pathIdentity:        "ApproveUpdateMouComponent",
	}, nil
}

func (approveMouTrx *approveUpdateMouTransactionComponent) PreTransaction(
	input *model.InternalUpdateMou,
) (*model.InternalUpdateMou, error) {
	return input, nil
}

func (approveMouTrx *approveUpdateMouTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateMou,
) (*model.Mou, error) {
	updateMou := &model.DatabaseUpdateMou{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateMou)

	existingMou, err := approveMouTrx.mouDataSource.GetMongoDataSource().FindByID(
		updateMou.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveMouTrx.pathIdentity,
			err,
		)
	}

	previousLog, err := approveMouTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingMou.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveMouTrx.pathIdentity,
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateMou.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateMou.ProposalStatus,
	}
	jsonTemp, _ = json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveMouTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveMouTrx.pathIdentity,
			err,
		)
	}

	updateMou.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	currentTime := time.Now()
	updateMou.UpdatedAt = &currentTime
	updateMou.RemainingCreditLimit = func(i int) *int {
		return &i
	}(
		existingMou.RemainingCreditLimit +
			existingMou.ProposedChanges.CreditLimit - existingMou.CreditLimit,
	)

	fieldsToUpdateMou := &model.DatabaseUpdateMou{
		ID: updateMou.ID,
	}
	jsonExisting, _ := json.Marshal(existingMou.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdateMou.ProposedChanges)

	var updateMouMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateMou)
	json.Unmarshal(jsonUpdate, &updateMouMap)

	approveMouTrx.mapProcessorUtility.RemoveNil(updateMouMap)

	jsonUpdate, _ = json.Marshal(updateMouMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateMou.ProposedChanges)

	if updateMou.ProposalStatus != nil {
		if *updateMou.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdateMou.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdateMou)
		}
	}

	updatedMou, err := approveMouTrx.mouDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateMou.ID,
		},
		fieldsToUpdateMou,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveMouTrx.pathIdentity,
			err,
		)
	}

	return updatedMou, nil
}
