package organizationdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type updateOrganizationTransactionComponent struct {
	organizationDataSource             databaseorganizationdatasourceinterfaces.OrganizationDataSource
	updateOrganizationUsecaseComponent organizationdomainrepositoryinterfaces.UpdateOrganizationUsecaseComponent
}

func NewUpdateOrganizationTransactionComponent(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
) (organizationdomainrepositoryinterfaces.UpdateOrganizationTransactionComponent, error) {
	return &updateOrganizationTransactionComponent{
		organizationDataSource: organizationDataSource,
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
	jsonExistingOrg, _ := json.Marshal(existingOrganization)
	jsonUpdateOrg, _ := json.Marshal(updateOrganization)
	json.Unmarshal(jsonExistingOrg, fieldsToUpdateOrganization.ProposedChanges)
	json.Unmarshal(jsonUpdateOrg, fieldsToUpdateOrganization.ProposedChanges)

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
