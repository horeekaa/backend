package organizationdomainrepositories

import (
	"encoding/json"
	"fmt"
	"reflect"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateOrganizationTransactionComponent struct {
	organizationDataSource                    databaseorganizationdatasourceinterfaces.OrganizationDataSource
	loggingDataSource                         databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility                       coreutilityinterfaces.MapProcessorUtility
	structComparisonUtility                   coreutilityinterfaces.StructComparisonUtility
	proposeUpdateOrganizationUsecaseComponent organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationUsecaseComponent
}

func NewProposeUpdateOrganizationTransactionComponent(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
	structComparisonUtility coreutilityinterfaces.StructComparisonUtility,
) (organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationTransactionComponent, error) {
	return &proposeUpdateOrganizationTransactionComponent{
		organizationDataSource:  organizationDataSource,
		loggingDataSource:       loggingDataSource,
		mapProcessorUtility:     mapProcessorUtility,
		structComparisonUtility: structComparisonUtility,
	}, nil
}

func (updateOrgTrx *proposeUpdateOrganizationTransactionComponent) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationUsecaseComponent,
) (bool, error) {
	updateOrgTrx.proposeUpdateOrganizationUsecaseComponent = usecaseComponent
	return true, nil
}

func (updateOrgTrx *proposeUpdateOrganizationTransactionComponent) PreTransaction(
	input *model.InternalUpdateOrganization,
) (*model.InternalUpdateOrganization, error) {
	if updateOrgTrx.proposeUpdateOrganizationUsecaseComponent == nil {
		return input, nil
	}
	return updateOrgTrx.proposeUpdateOrganizationUsecaseComponent.Validation(input)
}

func (updateOrgTrx *proposeUpdateOrganizationTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateOrganization *model.InternalUpdateOrganization,
) (*model.Organization, error) {
	existingOrganization, err := updateOrgTrx.organizationDataSource.GetMongoDataSource().FindByID(
		updateOrganization.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateOrganization",
			err,
		)
	}
	fieldChanges := []*model.FieldChangeDataInput{}

	updateOrgTrx.structComparisonUtility.SetComparisonFunc(
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
	updateOrgTrx.structComparisonUtility.SetPreDeepComparisonFunc(
		func(tag interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v.",
				*tagString,
				tag,
			)
		},
	)
	var tagString interface{} = ""
	updateOrgTrx.structComparisonUtility.CompareStructs(
		*updateOrganization,
		*existingOrganization,
		&tagString,
	)

	loggingOutput, err := updateOrgTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Organization",
			Document: &model.ObjectIDOnly{
				ID: &existingOrganization.ID,
			},
			FieldChanges: fieldChanges,
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateOrganization.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateOrganization.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateOrganization",
			err,
		)
	}
	updateOrganization.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdateOrganization := &model.InternalUpdateOrganization{
		ID: updateOrganization.ID,
	}
	jsonExisting, _ := json.Marshal(existingOrganization)
	json.Unmarshal(jsonExisting, &fieldsToUpdateOrganization.ProposedChanges)

	var updateOrganizationMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateOrganization)
	json.Unmarshal(jsonUpdate, &updateOrganizationMap)

	updateOrgTrx.mapProcessorUtility.RemoveNil(updateOrganizationMap)

	jsonUpdate, _ = json.Marshal(updateOrganizationMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateOrganization.ProposedChanges)

	if updateOrganization.ProposalStatus != nil {
		fieldsToUpdateOrganization.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateOrganization.SubmittingAccount.ID,
		}
		if *updateOrganization.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdateOrganization)
		}
	}

	updatedOrganization, err := updateOrgTrx.organizationDataSource.GetMongoDataSource().Update(
		fieldsToUpdateOrganization.ID,
		fieldsToUpdateOrganization,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateOrganization",
			err,
		)
	}

	return updatedOrganization, nil
}
