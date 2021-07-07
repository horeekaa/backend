package memberaccessrefdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type updateMemberAccessRefTransactionComponent struct {
	memberAccessRefDataSource             databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	mapProcessorUtility                   coreutilityinterfaces.MapProcessorUtility
	updateMemberAccessRefUsecaseComponent memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefUsecaseComponent
}

func NewUpdateMemberAccessRefTransactionComponent(
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefTransactionComponent, error) {
	return &updateMemberAccessRefTransactionComponent{
		memberAccessRefDataSource: memberAccessRefDataSource,
		mapProcessorUtility:       mapProcessorUtility,
	}, nil
}

func (updateMmbAccRefTrx *updateMemberAccessRefTransactionComponent) SetValidation(
	usecaseComponent memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefUsecaseComponent,
) (bool, error) {
	updateMmbAccRefTrx.updateMemberAccessRefUsecaseComponent = usecaseComponent
	return true, nil
}

func (updateMmbAccRefTrx *updateMemberAccessRefTransactionComponent) PreTransaction(
	input *model.InternalUpdateMemberAccessRef,
) (*model.InternalUpdateMemberAccessRef, error) {
	if updateMmbAccRefTrx.updateMemberAccessRefUsecaseComponent == nil {
		return input, nil
	}
	return updateMmbAccRefTrx.updateMemberAccessRefUsecaseComponent.Validation(input)
}

func (updateMmbAccRefTrx *updateMemberAccessRefTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateMemberAccessRef *model.InternalUpdateMemberAccessRef,
) (*model.MemberAccessRef, error) {
	existingMemberAccessRef, err := updateMmbAccRefTrx.memberAccessRefDataSource.GetMongoDataSource().FindByID(
		updateMemberAccessRef.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMemberAccessRef",
			err,
		)
	}
	fieldsToUpdateMemberAccessRef := &model.InternalUpdateMemberAccessRef{
		ID: updateMemberAccessRef.ID,
	}
	jsonExisting, _ := json.Marshal(existingMemberAccessRef)
	json.Unmarshal(jsonExisting, &fieldsToUpdateMemberAccessRef.ProposedChanges)

	var updateMemberAccessRefMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateMemberAccessRef)
	json.Unmarshal(jsonUpdate, &updateMemberAccessRefMap)

	updateMmbAccRefTrx.mapProcessorUtility.RemoveNil(updateMemberAccessRefMap)

	jsonUpdate, _ = json.Marshal(updateMemberAccessRefMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateMemberAccessRef.ProposedChanges)

	if updateMemberAccessRef.RecentApprovingAccount != nil &&
		updateMemberAccessRef.ProposalStatus != nil {
		if existingMemberAccessRef.ProposedChanges.ProposalStatus == model.EntityProposalStatusRejected {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.NothingToBeApproved,
				"/updateMemberAccessRef",
				nil,
			)
		}

		if *updateMemberAccessRef.ProposalStatus == model.EntityProposalStatusApproved {
			jsonTemp, _ := json.Marshal(fieldsToUpdateMemberAccessRef.ProposedChanges)
			json.Unmarshal(jsonTemp, fieldsToUpdateMemberAccessRef)
		}
	}

	updatedMemberAccessRef, err := updateMmbAccRefTrx.memberAccessRefDataSource.GetMongoDataSource().Update(
		fieldsToUpdateMemberAccessRef.ID,
		fieldsToUpdateMemberAccessRef,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMemberAccessRef",
			err,
		)
	}

	return updatedMemberAccessRef, nil
}
