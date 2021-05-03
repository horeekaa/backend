package memberaccessrefdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type updateMemberAccessRefTransactionComponent struct {
	memberAccessRefDataSource             databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	updateMemberAccessRefUsecaseComponent memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefUsecaseComponent
}

func NewUpdateMemberAccessRefTransactionComponent(
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
) (memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefTransactionComponent, error) {
	return &updateMemberAccessRefTransactionComponent{
		memberAccessRefDataSource: memberAccessRefDataSource,
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

		if &existingMemberAccessRef.PreviousEntity.ID != nil &&
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
	jb, _ := json.Marshal(updateMemberAccessRef)
	json.Unmarshal(jb, &combinedMemberAccessRef)
	proposedProposalStatus := model.EntityProposalStatusProposed
	combinedMemberAccessRef.ProposalStatus = &proposedProposalStatus

	combinedMemberAccessRef.PreviousEntity.ID = &existingMemberAccessRef.ID

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
