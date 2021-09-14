package memberaccessdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateMemberAccessTransactionComponent struct {
	memberAccessDataSource                    databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	loggingDataSource                         databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility                       coreutilityinterfaces.MapProcessorUtility
	approveUpdateMemberAccessUsecaseComponent memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessUsecaseComponent
}

func NewApproveUpdateMemberAccessTransactionComponent(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessTransactionComponent, error) {
	return &approveUpdateMemberAccessTransactionComponent{
		memberAccessDataSource: memberAccessDataSource,
		loggingDataSource:      loggingDataSource,
		mapProcessorUtility:    mapProcessorUtility,
	}, nil
}

func (approveProdTrx *approveUpdateMemberAccessTransactionComponent) SetValidation(
	usecaseComponent memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessUsecaseComponent,
) (bool, error) {
	approveProdTrx.approveUpdateMemberAccessUsecaseComponent = usecaseComponent
	return true, nil
}

func (approveProdTrx *approveUpdateMemberAccessTransactionComponent) PreTransaction(
	input *model.InternalUpdateMemberAccess,
) (*model.InternalUpdateMemberAccess, error) {
	if approveProdTrx.approveUpdateMemberAccessUsecaseComponent == nil {
		return input, nil
	}
	return approveProdTrx.approveUpdateMemberAccessUsecaseComponent.Validation(input)
}

func (approveProdTrx *approveUpdateMemberAccessTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateMemberAccess *model.InternalUpdateMemberAccess,
) (*model.MemberAccess, error) {
	existingMemberAccess, err := approveProdTrx.memberAccessDataSource.GetMongoDataSource().FindByID(
		updateMemberAccess.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMemberAccess",
			err,
		)
	}

	previousLog, err := approveProdTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingMemberAccess.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMemberAccess",
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateMemberAccess.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateMemberAccess.ProposalStatus,
	}
	jsonTemp, _ := json.Marshal(
		map[string]interface{}{
			"FieldChanges": previousLog.FieldChanges,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveProdTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)

	updateMemberAccess.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdateMemberAccess := &model.InternalUpdateMemberAccess{
		ID: updateMemberAccess.ID,
	}
	jsonExisting, _ := json.Marshal(existingMemberAccess.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdateMemberAccess.ProposedChanges)

	var updateMemberAccessMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateMemberAccess)
	json.Unmarshal(jsonUpdate, &updateMemberAccessMap)

	approveProdTrx.mapProcessorUtility.RemoveNil(updateMemberAccessMap)

	jsonUpdate, _ = json.Marshal(updateMemberAccessMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateMemberAccess.ProposedChanges)

	if updateMemberAccess.ProposalStatus != nil {
		if *updateMemberAccess.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdateMemberAccess.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdateMemberAccess)
		}
	}

	updatedMemberAccess, err := approveProdTrx.memberAccessDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateMemberAccess.ID,
		},
		fieldsToUpdateMemberAccess,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMemberAccess",
			err,
		)
	}

	return updatedMemberAccess, nil
}
