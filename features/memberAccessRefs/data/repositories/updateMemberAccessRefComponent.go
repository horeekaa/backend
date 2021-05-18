package memberaccessrefdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories/types"
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
	input *model.UpdateMemberAccessRef,
) (*model.UpdateMemberAccessRef, error) {
	if updateMmbAccRefTrx.updateMemberAccessRefUsecaseComponent == nil {
		return input, nil
	}
	return updateMmbAccRefTrx.updateMemberAccessRefUsecaseComponent.Validation(input)
}

func (updateMmbAccRefTrx *updateMemberAccessRefTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateMemberAccessRef *model.UpdateMemberAccessRef,
) (*memberaccessrefdomainrepositorytypes.UpdateMemberAccessRefOutput, error) {
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

	if updateMemberAccessRef.ApprovingAccount != nil &&
		updateMemberAccessRef.ProposalStatus != nil {
		updatedMemberAccessRef, err := updateMmbAccRefTrx.memberAccessRefDataSource.GetMongoDataSource().Update(
			existingMemberAccessRef.ID,
			updateMemberAccessRef,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateMemberAccessRef",
				err,
			)
		}

		if existingMemberAccessRef.PreviousEntity != nil &&
			*updateMemberAccessRef.ProposalStatus == model.EntityProposalStatusApproved {
			replacedProposalStatus := model.EntityProposalStatusReplaced
			previousMemberAccessRef, err := updateMmbAccRefTrx.memberAccessRefDataSource.GetMongoDataSource().Update(
				existingMemberAccessRef.PreviousEntity.ID,
				&model.UpdateMemberAccessRef{
					ProposalStatus: &replacedProposalStatus,
				},
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/updateMemberAccessRef",
					err,
				)
			}
			return &memberaccessrefdomainrepositorytypes.UpdateMemberAccessRefOutput{
				PreviousMemberAccessRef: previousMemberAccessRef,
				UpdatedMemberAccessRef:  updatedMemberAccessRef,
			}, nil
		}

		return &memberaccessrefdomainrepositorytypes.UpdateMemberAccessRefOutput{
			PreviousMemberAccessRef: existingMemberAccessRef,
			UpdatedMemberAccessRef:  updatedMemberAccessRef,
		}, nil
	}

	var combinedMemberAccessRef model.CreateMemberAccessRef
	ja, _ := json.Marshal(existingMemberAccessRef)
	json.Unmarshal(ja, &combinedMemberAccessRef)

	var updateMemberAccessRefMap map[string]interface{}
	jsonTemp, _ := json.Marshal(updateMemberAccessRef)
	json.Unmarshal(jsonTemp, &updateMemberAccessRefMap)

	updateMmbAccRefTrx.mapProcessorUtility.RemoveNil(updateMemberAccessRefMap)

	jb, _ := json.Marshal(updateMemberAccessRefMap)
	json.Unmarshal(jb, &combinedMemberAccessRef)
	proposedProposalStatus := model.EntityProposalStatusProposed
	combinedMemberAccessRef.ProposalStatus = &proposedProposalStatus

	combinedMemberAccessRef.PreviousEntity = &model.ObjectIDOnly{ID: &existingMemberAccessRef.ID}

	updatedMemberAccessRef, err := updateMmbAccRefTrx.memberAccessRefDataSource.GetMongoDataSource().Create(
		&combinedMemberAccessRef,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMemberAccessRef",
			err,
		)
	}

	return &memberaccessrefdomainrepositorytypes.UpdateMemberAccessRefOutput{
		PreviousMemberAccessRef: existingMemberAccessRef,
		UpdatedMemberAccessRef:  updatedMemberAccessRef,
	}, nil
}
