package memberaccessrefdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateMemberAccessRefTransactionComponent struct {
	memberAccessRefDataSource                    databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	loggingDataSource                            databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility                          coreutilityinterfaces.MapProcessorUtility
	proposeUpdateMemberAccessRefUsecaseComponent memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefUsecaseComponent
	pathIdentity                                 string
}

func NewProposeUpdateMemberAccessRefTransactionComponent(
	MemberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefTransactionComponent, error) {
	return &proposeUpdateMemberAccessRefTransactionComponent{
		memberAccessRefDataSource: MemberAccessRefDataSource,
		loggingDataSource:         loggingDataSource,
		mapProcessorUtility:       mapProcessorUtility,
		pathIdentity:              "ProposeUpdateMemberAccessRefComponent",
	}, nil
}

func (updateMemberAccessRefTrx *proposeUpdateMemberAccessRefTransactionComponent) SetValidation(
	usecaseComponent memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefUsecaseComponent,
) (bool, error) {
	updateMemberAccessRefTrx.proposeUpdateMemberAccessRefUsecaseComponent = usecaseComponent
	return true, nil
}

func (updateMemberAccessRefTrx *proposeUpdateMemberAccessRefTransactionComponent) PreTransaction(
	input *model.InternalUpdateMemberAccessRef,
) (*model.InternalUpdateMemberAccessRef, error) {
	if updateMemberAccessRefTrx.proposeUpdateMemberAccessRefUsecaseComponent == nil {
		return input, nil
	}
	return updateMemberAccessRefTrx.proposeUpdateMemberAccessRefUsecaseComponent.Validation(input)
}

func (updateMemberAccessRefTrx *proposeUpdateMemberAccessRefTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateMemberAccessRef,
) (*model.MemberAccessRef, error) {
	updateMemberAccessRef := &model.DatabaseUpdateMemberAccessRef{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateMemberAccessRef)

	existingMemberAccessRef, err := updateMemberAccessRefTrx.memberAccessRefDataSource.GetMongoDataSource().FindByID(
		updateMemberAccessRef.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateMemberAccessRefTrx.pathIdentity,
			err,
		)
	}

	newDocumentJson, _ := json.Marshal(*updateMemberAccessRef)
	oldDocumentJson, _ := json.Marshal(*existingMemberAccessRef)
	loggingOutput, err := updateMemberAccessRefTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "MemberAccessRef",
			Document: &model.ObjectIDOnly{
				ID: &existingMemberAccessRef.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateMemberAccessRef.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateMemberAccessRef.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateMemberAccessRefTrx.pathIdentity,
			err,
		)
	}
	updateMemberAccessRef.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	var currentTime = time.Now().UTC()
	updateMemberAccessRef.UpdatedAt = &currentTime

	fieldsToUpdateMemberAccessRef := &model.DatabaseUpdateMemberAccessRef{
		ID: updateMemberAccessRef.ID,
	}
	jsonExisting, _ := json.Marshal(existingMemberAccessRef)
	json.Unmarshal(jsonExisting, &fieldsToUpdateMemberAccessRef.ProposedChanges)

	var updateMemberAccessRefMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateMemberAccessRef)
	json.Unmarshal(jsonUpdate, &updateMemberAccessRefMap)

	updateMemberAccessRefTrx.mapProcessorUtility.RemoveNil(updateMemberAccessRefMap)

	jsonUpdate, _ = json.Marshal(updateMemberAccessRefMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateMemberAccessRef.ProposedChanges)

	if updateMemberAccessRef.ProposalStatus != nil {
		fieldsToUpdateMemberAccessRef.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateMemberAccessRef.SubmittingAccount.ID,
		}
		if *updateMemberAccessRef.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdateMemberAccessRef)
		}
	}

	updatedMemberAccessRef, err := updateMemberAccessRefTrx.memberAccessRefDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateMemberAccessRef.ID,
		},
		fieldsToUpdateMemberAccessRef,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateMemberAccessRefTrx.pathIdentity,
			err,
		)
	}

	return updatedMemberAccessRef, nil
}
