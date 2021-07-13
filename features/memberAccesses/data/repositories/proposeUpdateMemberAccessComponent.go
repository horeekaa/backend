package memberaccessdomainrepositories

import (
	"encoding/json"
	"fmt"
	"reflect"

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
	structComparisonUtility                   coreutilityinterfaces.StructComparisonUtility
	proposeUpdateMemberAccessUsecaseComponent memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessUsecaseComponent
}

func NewProposeUpdateMemberAccessTransactionComponent(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
	structComparisonUtility coreutilityinterfaces.StructComparisonUtility,
) (memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessTransactionComponent, error) {
	return &proposeUpdateMemberAccessTransactionComponent{
		memberAccessDataSource:    memberAccessDataSource,
		loggingDataSource:         loggingDataSource,
		organizationDataSource:    organizationDataSource,
		memberAccessRefDataSource: memberAccessRefDataSource,
		mapProcessorUtility:       mapProcessorUtility,
		structComparisonUtility:   structComparisonUtility,
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
	updateMemberAccess *model.InternalUpdateMemberAccess,
) (*model.MemberAccess, error) {
	existingMemberAccess, err := proposeUpdateMemberAccTrx.memberAccessDataSource.GetMongoDataSource().FindByID(
		updateMemberAccess.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateMemberAccess",
			err,
		)
	}
	fieldChanges := []*model.FieldChangeDataInput{}

	proposeUpdateMemberAccTrx.structComparisonUtility.SetComparisonFunc(
		func(tag interface{}, field1 interface{}, field2 interface{}, tagString *interface{}) {
			if field1 == field2 {
				return
			}
			*tagString = fmt.Sprintf(
				"%v%v",
				*tagString,
				tag,
			)

			fieldChanges = append(fieldChanges, &model.FieldChangeDataInput{
				Name:     fmt.Sprint(*tagString),
				Type:     reflect.TypeOf(field1).Kind().String(),
				OldValue: fmt.Sprint(field2),
				NewValue: fmt.Sprint(field1),
			})
			*tagString = ""
		},
	)
	proposeUpdateMemberAccTrx.structComparisonUtility.SetPreDeepComparisonFunc(
		func(tag interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v.",
				*tagString,
				tag,
			)
		},
	)
	var tagString interface{} = ""
	proposeUpdateMemberAccTrx.structComparisonUtility.CompareStructs(
		*updateMemberAccess,
		*existingMemberAccess,
		&tagString,
	)

	queryMap := map[string]interface{}{
		"memberAccessRefType": existingMemberAccess.MemberAccessRefType,
		"proposalStatus":      model.EntityProposalStatusApproved,
	}
	if updateMemberAccess.OrganizationMembershipRole != nil &&
		existingMemberAccess.OrganizationMembershipRole != updateMemberAccess.OrganizationMembershipRole {
		queryMap["organizationMembershipRole"] = *updateMemberAccess.OrganizationMembershipRole
	}
	if updateMemberAccess.Organization != nil {
		orgToUpdate, err := proposeUpdateMemberAccTrx.organizationDataSource.GetMongoDataSource().FindByID(
			updateMemberAccess.Organization.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateMemberAccess",
				err,
			)
		}
		queryMap["organizationType"] = orgToUpdate.Type

		jsonTemp, _ := json.Marshal(orgToUpdate)
		json.Unmarshal(jsonTemp, &updateMemberAccess.Organization)
		json.Unmarshal(jsonTemp, &updateMemberAccess.OrganizationLatestUpdate)
	}

	if queryMap["organizationMembershipRole"] != nil || queryMap["organizationType"] != nil {
		memberAccessRef, err := proposeUpdateMemberAccTrx.memberAccessRefDataSource.GetMongoDataSource().FindOne(
			queryMap,
			session,
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

	loggingOutput, err := proposeUpdateMemberAccTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "MemberAccess",
			Document: &model.ObjectIDOnly{
				ID: &existingMemberAccess.ID,
			},
			FieldChanges: fieldChanges,
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
			"/updateMemberAccess",
			err,
		)
	}
	updateMemberAccess.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdateMemberAccess := &model.InternalUpdateMemberAccess{
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
		if *updateMemberAccess.ProposalStatus == model.EntityProposalStatusApproved {
			jsonTemp, _ := json.Marshal(fieldsToUpdateMemberAccess.ProposedChanges)
			json.Unmarshal(jsonTemp, fieldsToUpdateMemberAccess)
		}
	}

	updatedMemberAccess, err := proposeUpdateMemberAccTrx.memberAccessDataSource.GetMongoDataSource().Update(
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
