package organizationdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type updateOrganizationTransactionComponent struct {
	organizationDataSource             databaseorganizationdatasourceinterfaces.OrganizationDataSource
	mapProcessorUtility                coreutilityinterfaces.MapProcessorUtility
	updateOrganizationUsecaseComponent organizationdomainrepositoryinterfaces.UpdateOrganizationUsecaseComponent
}

func NewUpdateOrganizationTransactionComponent(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (organizationdomainrepositoryinterfaces.UpdateOrganizationTransactionComponent, error) {
	return &updateOrganizationTransactionComponent{
		organizationDataSource: organizationDataSource,
		mapProcessorUtility:    mapProcessorUtility,
	}, nil
}

func (updateOrgTrx *updateOrganizationTransactionComponent) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.UpdateOrganizationUsecaseComponent,
) (bool, error) {
	updateOrgTrx.updateOrganizationUsecaseComponent = usecaseComponent
	return true, nil
}

func (updateOrgTrx *updateOrganizationTransactionComponent) PreTransaction(
	input *model.InternalUpdateOrganization,
) (*model.InternalUpdateOrganization, error) {
	if updateOrgTrx.updateOrganizationUsecaseComponent == nil {
		return input, nil
	}
	return updateOrgTrx.updateOrganizationUsecaseComponent.Validation(input)
}

func (updateOrgTrx *updateOrganizationTransactionComponent) TransactionBody(
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

	if updateOrganization.RecentApprovingAccount != nil &&
		updateOrganization.ProposalStatus != nil {
		if existingOrganization.ProposedChanges.ProposalStatus == model.EntityProposalStatusRejected {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.NothingToBeApproved,
				"/updateOrganization",
				nil,
			)
		}

		if *updateOrganization.ProposalStatus == model.EntityProposalStatusApproved {
			jsonTemp, _ := json.Marshal(fieldsToUpdateOrganization.ProposedChanges)
			json.Unmarshal(jsonTemp, fieldsToUpdateOrganization)
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
