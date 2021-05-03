package organizationdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/serviceFailures/exceptionToFailure"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationdomainrepositorytypes "github.com/horeekaa/backend/features/organizations/domain/repositories/types"
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
	input *model.UpdateOrganization,
) (*model.UpdateOrganization, error) {
	if updateOrgTrx.updateOrganizationUsecaseComponent == nil {
		return input, nil
	}
	return updateOrgTrx.updateOrganizationUsecaseComponent.Validation(input)
}

func (updateOrgTrx *updateOrganizationTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateOrganization *model.UpdateOrganization,
) (*organizationdomainrepositorytypes.UpdateOrganizationOutput, error) {
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

	if updateOrganization.ApprovingAccount != nil &&
		updateOrganization.ProposalStatus != nil {
		updatedOrganization, err := updateOrgTrx.organizationDataSource.GetMongoDataSource().Update(
			existingOrganization.ID,
			updateOrganization,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateOrganization",
				err,
			)
		}

		if &existingOrganization.PreviousEntity.ID != nil &&
			*updateOrganization.ProposalStatus == model.EntityProposalStatusApproved {
			replacedProposalStatus := model.EntityProposalStatusReplaced
			previousOrganization, err := updateOrgTrx.organizationDataSource.GetMongoDataSource().Update(
				existingOrganization.PreviousEntity.ID,
				&model.UpdateOrganization{
					ProposalStatus: &replacedProposalStatus,
				},
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/updateOrganization",
					err,
				)
			}
			return &organizationdomainrepositorytypes.UpdateOrganizationOutput{
				PreviousOrganization: previousOrganization,
				UpdatedOrganization:  updatedOrganization,
			}, nil
		}

		return &organizationdomainrepositorytypes.UpdateOrganizationOutput{
			PreviousOrganization: existingOrganization,
			UpdatedOrganization:  updatedOrganization,
		}, nil
	}

	var combinedOrganization model.CreateOrganization
	ja, _ := json.Marshal(existingOrganization)
	json.Unmarshal(ja, &combinedOrganization)
	jb, _ := json.Marshal(updateOrganization)
	json.Unmarshal(jb, &combinedOrganization)
	proposedProposalStatus := model.EntityProposalStatusProposed
	combinedOrganization.ProposalStatus = &proposedProposalStatus

	combinedOrganization.PreviousEntity.ID = &existingOrganization.ID

	updatedOrganization, err := updateOrgTrx.organizationDataSource.GetMongoDataSource().Create(
		&combinedOrganization,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateOrganization",
			err,
		)
	}

	return &organizationdomainrepositorytypes.UpdateOrganizationOutput{
		PreviousOrganization: existingOrganization,
		UpdatedOrganization:  updatedOrganization,
	}, nil
}
