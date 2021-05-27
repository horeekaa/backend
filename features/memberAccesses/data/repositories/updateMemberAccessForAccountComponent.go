package memberaccessdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type updateMemberAccessForAccountComponent struct {
	memberAccessDataSource    databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	mapProcessorUtility       coreutilityinterfaces.MapProcessorUtility
}

func NewUpdateMemberAccessForAccountTransactionComponent(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (memberaccessdomainrepositoryinterfaces.UpdateMemberAccessForAccountTransactionComponent, error) {
	return &updateMemberAccessForAccountComponent{
		memberAccessDataSource,
		memberAccessRefDataSource,
		mapProcessorUtility,
	}, nil
}

func (updateMmbAccForAccountTrx *updateMemberAccessForAccountComponent) PreTransaction(
	input *model.UpdateMemberAccess,
) (*model.UpdateMemberAccess, error) {
	return input, nil
}

func (updateMmbAccForAccountTrx *updateMemberAccessForAccountComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateMemberAccess *model.UpdateMemberAccess,
) (*memberaccessdomainrepositorytypes.UpdateMemberAccessOutput, error) {
	existingMemberAccess, err := updateMmbAccForAccountTrx.memberAccessDataSource.GetMongoDataSource().FindByID(
		updateMemberAccess.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMemberAccess",
			err,
		)
	}

	queryMap := map[string]interface{}{
		"memberAccessRefType": existingMemberAccess.MemberAccessRefType,
		"proposalStatus":      model.EntityProposalStatusApproved,
	}
	if updateMemberAccess.OrganizationMembershipRole != nil &&
		existingMemberAccess.OrganizationMembershipRole != updateMemberAccess.OrganizationMembershipRole {
		queryMap["organizationMembershipRole"] = *updateMemberAccess.OrganizationMembershipRole
	}
	if updateMemberAccess.Organization != nil {
		queryMap["organizationType"] = *updateMemberAccess.Organization.Type
	}

	if queryMap["organizationMembershipRole"] != nil || queryMap["organizationType"] != nil {
		memberAccessRef, err := updateMmbAccForAccountTrx.memberAccessRefDataSource.GetMongoDataSource().FindOne(
			queryMap,
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateMemberAccess",
				err,
			)
		}
		if memberAccessRef == nil {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.MemberAccessRefNotExist,
				"/updateMemberAccess",
				nil,
			)
		}
		var accessInput model.MemberAccessRefOptionsInput
		jsonTemp, _ := json.Marshal(memberAccessRef.Access)
		json.Unmarshal(jsonTemp, &accessInput)
		updateMemberAccess.Access = &accessInput
		updateMemberAccess.DefaultAccess = &model.ObjectIDOnly{ID: &memberAccessRef.ID}
	}

	if updateMemberAccess.ApprovingAccount != nil &&
		updateMemberAccess.ProposalStatus != nil {
		updatedMemberAccess, err := updateMmbAccForAccountTrx.memberAccessDataSource.GetMongoDataSource().Update(
			existingMemberAccess.ID,
			updateMemberAccess,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateMemberAccess",
				err,
			)
		}

		if existingMemberAccess.PreviousEntity != nil &&
			*updateMemberAccess.ProposalStatus == model.EntityProposalStatusApproved {
			replacedProposalStatus := model.EntityProposalStatusReplaced
			previousMemberAccess, err := updateMmbAccForAccountTrx.memberAccessDataSource.GetMongoDataSource().Update(
				existingMemberAccess.PreviousEntity.ID,
				&model.UpdateMemberAccess{
					ProposalStatus: &replacedProposalStatus,
				},
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/updateMemberAccess",
					err,
				)
			}
			return &memberaccessdomainrepositorytypes.UpdateMemberAccessOutput{
				PreviousMemberAccess: previousMemberAccess,
				UpdatedMemberAccess:  updatedMemberAccess,
			}, nil
		}

		return &memberaccessdomainrepositorytypes.UpdateMemberAccessOutput{
			PreviousMemberAccess: existingMemberAccess,
			UpdatedMemberAccess:  updatedMemberAccess,
		}, nil
	}

	var combinedMemberAccess model.CreateMemberAccess
	ja, _ := json.Marshal(existingMemberAccess)
	json.Unmarshal(ja, &combinedMemberAccess)

	var updateMemberAccessMap map[string]interface{}
	jsonTemp, _ := json.Marshal(updateMemberAccess)
	json.Unmarshal(jsonTemp, &updateMemberAccessMap)

	updateMmbAccForAccountTrx.mapProcessorUtility.RemoveNil(updateMemberAccessMap)

	jb, _ := json.Marshal(updateMemberAccessMap)
	json.Unmarshal(jb, &combinedMemberAccess)
	proposedProposalStatus := model.EntityProposalStatusProposed
	combinedMemberAccess.ProposalStatus = &proposedProposalStatus

	combinedMemberAccess.PreviousEntity = &model.ObjectIDOnly{ID: &existingMemberAccess.ID}

	updatedMemberAccess, err := updateMmbAccForAccountTrx.memberAccessDataSource.GetMongoDataSource().Create(
		&combinedMemberAccess,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMemberAccess",
			err,
		)
	}

	return &memberaccessdomainrepositorytypes.UpdateMemberAccessOutput{
		PreviousMemberAccess: existingMemberAccess,
		UpdatedMemberAccess:  updatedMemberAccess,
	}, nil
}
