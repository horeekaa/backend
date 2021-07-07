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
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
)

type updateMemberAccessForAccountComponent struct {
	memberAccessDataSource    databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	organizationDataSource    databaseorganizationdatasourceinterfaces.OrganizationDataSource
	mapProcessorUtility       coreutilityinterfaces.MapProcessorUtility
}

func NewUpdateMemberAccessForAccountTransactionComponent(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (memberaccessdomainrepositoryinterfaces.UpdateMemberAccessForAccountTransactionComponent, error) {
	return &updateMemberAccessForAccountComponent{
		memberAccessDataSource,
		memberAccessRefDataSource,
		organizationDataSource,
		mapProcessorUtility,
	}, nil
}

func (updateMmbAccForAccountTrx *updateMemberAccessForAccountComponent) PreTransaction(
	input *model.InternalUpdateMemberAccess,
) (*model.InternalUpdateMemberAccess, error) {
	return input, nil
}

func (updateMmbAccForAccountTrx *updateMemberAccessForAccountComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateMemberAccess *model.InternalUpdateMemberAccess,
) (*model.MemberAccess, error) {
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

		orgToUpdate, err := updateMmbAccForAccountTrx.organizationDataSource.GetMongoDataSource().FindByID(
			updateMemberAccess.Organization.ID,
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateMemberAccess",
				err,
			)
		}

		jsonTemp, _ := json.Marshal(orgToUpdate)
		json.Unmarshal(jsonTemp, &updateMemberAccess.Organization)
		json.Unmarshal(jsonTemp, &updateMemberAccess.OrganizationLatestUpdate)
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
		jsonTemp, _ := json.Marshal(memberAccessRef.Access)
		json.Unmarshal(jsonTemp, &updateMemberAccess.Access)

		jsonTemp, _ = json.Marshal(memberAccessRef)
		json.Unmarshal(jsonTemp, &updateMemberAccess.DefaultAccess)

		updateMemberAccess.DefaultAccessLatestUpdate = &model.ObjectIDOnly{
			ID: &memberAccessRef.ID,
		}
	}

	fieldsToUpdateMemberAccess := &model.InternalUpdateMemberAccess{
		ID: updateMemberAccess.ID,
	}
	jsonExisting, _ := json.Marshal(existingMemberAccess)
	json.Unmarshal(jsonExisting, &fieldsToUpdateMemberAccess.ProposedChanges)

	var updateMemberAccessMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateMemberAccess)
	json.Unmarshal(jsonUpdate, &updateMemberAccessMap)

	updateMmbAccForAccountTrx.mapProcessorUtility.RemoveNil(updateMemberAccessMap)

	jsonUpdate, _ = json.Marshal(updateMemberAccessMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateMemberAccess.ProposedChanges)

	if updateMemberAccess.RecentApprovingAccount != nil &&
		updateMemberAccess.ProposalStatus != nil {
		if existingMemberAccess.ProposedChanges.ProposalStatus == model.EntityProposalStatusRejected {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.NothingToBeApproved,
				"/updateMemberAccess",
				nil,
			)
		}

		if *updateMemberAccess.ProposalStatus == model.EntityProposalStatusApproved {
			jsonTemp, _ := json.Marshal(fieldsToUpdateMemberAccess.ProposedChanges)
			json.Unmarshal(jsonTemp, fieldsToUpdateMemberAccess)
		}
	}

	updatedMemberAccess, err := updateMmbAccForAccountTrx.memberAccessDataSource.GetMongoDataSource().Update(
		fieldsToUpdateMemberAccess.ID,
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
