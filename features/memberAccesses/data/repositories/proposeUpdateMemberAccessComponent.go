package memberaccessdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateMemberAccessTransactionComponent struct {
	memberAccessDataSource                    databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	loggingDataSource                         databaseloggingdatasourceinterfaces.LoggingDataSource
	organizationDataSource                    databaseorganizationdatasourceinterfaces.OrganizationDataSource
	memberAccessRefDataSource                 databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource
	mapProcessorUtility                       coreutilityinterfaces.MapProcessorUtility
	proposeUpdateMemberAccessUsecaseComponent memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessUsecaseComponent
	pathIdentity                              string
}

func NewProposeUpdateMemberAccessTransactionComponent(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessTransactionComponent, error) {
	return &proposeUpdateMemberAccessTransactionComponent{
		memberAccessDataSource:    memberAccessDataSource,
		loggingDataSource:         loggingDataSource,
		organizationDataSource:    organizationDataSource,
		memberAccessRefDataSource: memberAccessRefDataSource,
		mapProcessorUtility:       mapProcessorUtility,
		pathIdentity:              "ProposeUpdateMemberAccessComponent",
	}, nil
}

func (proposeUpdateMemberAccTrx *proposeUpdateMemberAccessTransactionComponent) SetValidation(
	usecaseComponent memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessUsecaseComponent,
) (bool, error) {
	proposeUpdateMemberAccTrx.proposeUpdateMemberAccessUsecaseComponent = usecaseComponent
	return true, nil
}

func (proposeUpdateMemberAccTrx *proposeUpdateMemberAccessTransactionComponent) PreTransaction(
	input *model.InternalUpdateMemberAccess,
) (*model.InternalUpdateMemberAccess, error) {
	if proposeUpdateMemberAccTrx.proposeUpdateMemberAccessUsecaseComponent == nil {
		return input, nil
	}
	return proposeUpdateMemberAccTrx.proposeUpdateMemberAccessUsecaseComponent.Validation(input)
}

func (proposeUpdateMemberAccTrx *proposeUpdateMemberAccessTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateMemberAccess,
) (*model.MemberAccess, error) {
	updateMemberAccess := &model.DatabaseUpdateMemberAccess{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateMemberAccess)

	existingMemberAccess, err := proposeUpdateMemberAccTrx.memberAccessDataSource.GetMongoDataSource().FindByID(
		updateMemberAccess.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			proposeUpdateMemberAccTrx.pathIdentity,
			err,
		)
	}
	if input.InvitationAccepted != nil {
		if existingMemberAccess.Account.ID.Hex() != input.SubmittingAccount.ID.Hex() ||
			existingMemberAccess.ProposalStatus != model.EntityProposalStatusApproved {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.AcceptInvitationNotAllowed,
				proposeUpdateMemberAccTrx.pathIdentity,
				nil,
			)
		}
	}

	queryMap := map[string]interface{}{
		"memberAccessRefType": existingMemberAccess.MemberAccessRefType,
		"proposalStatus":      model.EntityProposalStatusApproved,
	}
	if updateMemberAccess.OrganizationMembershipRole != nil {
		if *existingMemberAccess.OrganizationMembershipRole != *updateMemberAccess.OrganizationMembershipRole {
			queryMap["organizationMembershipRole"] = *updateMemberAccess.OrganizationMembershipRole
		}
	}
	if updateMemberAccess.Organization != nil {
		orgToUpdate, err := proposeUpdateMemberAccTrx.organizationDataSource.GetMongoDataSource().FindByID(
			*updateMemberAccess.Organization.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				proposeUpdateMemberAccTrx.pathIdentity,
				err,
			)
		}
		queryMap["organizationType"] = orgToUpdate.Type

		jsonTemp, _ := json.Marshal(orgToUpdate)
		json.Unmarshal(jsonTemp, &updateMemberAccess.Organization)
	}

	if queryMap["organizationMembershipRole"] != nil {
		if queryMap["organizationType"] == nil {
			queryMap["organizationType"] = existingMemberAccess.Organization.Type
		}
		memberAccessRef, err := proposeUpdateMemberAccTrx.memberAccessRefDataSource.GetMongoDataSource().FindOne(
			queryMap,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				proposeUpdateMemberAccTrx.pathIdentity,
				err,
			)
		}
		if memberAccessRef == nil {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.MemberAccessRefNotExist,
				proposeUpdateMemberAccTrx.pathIdentity,
				nil,
			)
		}
		jsonTemp, _ := json.Marshal(memberAccessRef.Access)
		json.Unmarshal(jsonTemp, &updateMemberAccess.Access)

		updateMemberAccess.DefaultAccessLatestUpdate = &model.ObjectIDOnly{
			ID: &memberAccessRef.ID,
		}
	}

	newDocumentJson, _ := json.Marshal(*updateMemberAccess)
	oldDocumentJson, _ := json.Marshal(*existingMemberAccess)
	loggingOutput, err := proposeUpdateMemberAccTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "MemberAccess",
			Document: &model.ObjectIDOnly{
				ID: &existingMemberAccess.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateMemberAccess.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateMemberAccess.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			proposeUpdateMemberAccTrx.pathIdentity,
			err,
		)
	}
	updateMemberAccess.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	var currentTime = time.Now()
	updateMemberAccess.UpdatedAt = &currentTime

	fieldsToUpdateMemberAccess := &model.DatabaseUpdateMemberAccess{
		ID: updateMemberAccess.ID,
	}
	jsonExisting, _ := json.Marshal(existingMemberAccess)
	json.Unmarshal(jsonExisting, &fieldsToUpdateMemberAccess.ProposedChanges)

	var updateMemberAccessMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateMemberAccess)
	json.Unmarshal(jsonUpdate, &updateMemberAccessMap)

	proposeUpdateMemberAccTrx.mapProcessorUtility.RemoveNil(updateMemberAccessMap)

	jsonUpdate, _ = json.Marshal(updateMemberAccessMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateMemberAccess.ProposedChanges)

	if updateMemberAccess.ProposalStatus != nil {
		fieldsToUpdateMemberAccess.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateMemberAccess.SubmittingAccount.ID,
		}
		if *updateMemberAccess.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdateMemberAccess)
		}
	}

	updatedMemberAccess, err := proposeUpdateMemberAccTrx.memberAccessDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateMemberAccess.ID,
		},
		fieldsToUpdateMemberAccess,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			proposeUpdateMemberAccTrx.pathIdentity,
			err,
		)
	}

	return updatedMemberAccess, nil
}
