package memberaccessrefdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberAccessRefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberAccessRefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateMemberAccessRefTransactionComponent struct {
	memberAccessRefDataSource                    databasememberAccessRefdatasourceinterfaces.MemberAccessRefDataSource
	loggingDataSource                            databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility                          coreutilityinterfaces.MapProcessorUtility
	approveUpdateMemberAccessRefUsecaseComponent memberAccessRefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefUsecaseComponent
	pathIdentity                                 string
}

func NewApproveUpdateMemberAccessRefTransactionComponent(
	memberAccessRefDataSource databasememberAccessRefdatasourceinterfaces.MemberAccessRefDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (memberAccessRefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefTransactionComponent, error) {
	return &approveUpdateMemberAccessRefTransactionComponent{
		memberAccessRefDataSource: memberAccessRefDataSource,
		loggingDataSource:         loggingDataSource,
		mapProcessorUtility:       mapProcessorUtility,
		pathIdentity:              "ApproveUpdateMemberAccessRefComponent",
	}, nil
}

func (approveProdTrx *approveUpdateMemberAccessRefTransactionComponent) SetValidation(
	usecaseComponent memberAccessRefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefUsecaseComponent,
) (bool, error) {
	approveProdTrx.approveUpdateMemberAccessRefUsecaseComponent = usecaseComponent
	return true, nil
}

func (approveProdTrx *approveUpdateMemberAccessRefTransactionComponent) PreTransaction(
	input *model.InternalUpdateMemberAccessRef,
) (*model.InternalUpdateMemberAccessRef, error) {
	if approveProdTrx.approveUpdateMemberAccessRefUsecaseComponent == nil {
		return input, nil
	}
	return approveProdTrx.approveUpdateMemberAccessRefUsecaseComponent.Validation(input)
}

func (approveProdTrx *approveUpdateMemberAccessRefTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateMemberAccessRef,
) (*model.MemberAccessRef, error) {
	updateMemberAccessRef := &model.DatabaseUpdateMemberAccessRef{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateMemberAccessRef)

	existingmemberAccessRef, err := approveProdTrx.memberAccessRefDataSource.GetMongoDataSource().FindByID(
		updateMemberAccessRef.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveProdTrx.pathIdentity,
			err,
		)
	}

	previousLog, err := approveProdTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingmemberAccessRef.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveProdTrx.pathIdentity,
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateMemberAccessRef.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateMemberAccessRef.ProposalStatus,
	}
	jsonTemp, _ = json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveProdTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveProdTrx.pathIdentity,
			err,
		)
	}

	updateMemberAccessRef.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	var currentTime = time.Now()
	updateMemberAccessRef.UpdatedAt = &currentTime

	fieldsToUpdateMemberAccessRef := &model.DatabaseUpdateMemberAccessRef{
		ID: updateMemberAccessRef.ID,
	}
	jsonExisting, _ := json.Marshal(existingmemberAccessRef.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdateMemberAccessRef.ProposedChanges)

	var updateMemberAccessRefMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateMemberAccessRef)
	json.Unmarshal(jsonUpdate, &updateMemberAccessRefMap)

	approveProdTrx.mapProcessorUtility.RemoveNil(updateMemberAccessRefMap)

	jsonUpdate, _ = json.Marshal(updateMemberAccessRefMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateMemberAccessRef.ProposedChanges)

	if updateMemberAccessRef.ProposalStatus != nil {
		if *updateMemberAccessRef.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdateMemberAccessRef.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdateMemberAccessRef)
		}
	}

	updatedMemberAccessRef, err := approveProdTrx.memberAccessRefDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateMemberAccessRef.ID,
		},
		fieldsToUpdateMemberAccessRef,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveProdTrx.pathIdentity,
			err,
		)
	}

	return updatedMemberAccessRef, nil
}
