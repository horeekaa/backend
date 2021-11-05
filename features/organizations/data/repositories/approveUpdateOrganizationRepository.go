package organizationdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type approveUpdateOrganizationRepository struct {
	approveUpdateOrganizationTransactionComponent organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationTransactionComponent
	organizationDataSource                        databaseorganizationdatasourceinterfaces.OrganizationDataSource
	bulkApproveUpdateTaggingComponent             taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent
	mongoDBTransaction                            mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdateOrganizationRepository(
	approveUpdateOrganizationRepositoryTransactionComponent organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationTransactionComponent,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	bulkApproveUpdateTaggingComponent taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationRepository, error) {
	approveUpdateOrganizationRepo := &approveUpdateOrganizationRepository{
		approveUpdateOrganizationRepositoryTransactionComponent,
		organizationDataSource,
		bulkApproveUpdateTaggingComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		approveUpdateOrganizationRepo,
		"ApproveUpdateOrganizationRepository",
	)

	return approveUpdateOrganizationRepo, nil
}

func (approveUpdateOrgRepo *approveUpdateOrganizationRepository) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationUsecaseComponent,
) (bool, error) {
	approveUpdateOrgRepo.approveUpdateOrganizationTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (approveUpdateOrgRepo *approveUpdateOrganizationRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return approveUpdateOrgRepo.approveUpdateOrganizationTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateOrganization),
	)
}

func (approveUpdateOrgRepo *approveUpdateOrganizationRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	organizationToApprove := input.(*model.InternalUpdateOrganization)
	existingOrganization, err := approveUpdateOrgRepo.organizationDataSource.GetMongoDataSource().FindByID(
		organizationToApprove.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdateOrganizationRepository",
			err,
		)
	}

	if existingOrganization.ProposedChanges.ProposalStatus == model.EntityProposalStatusProposed {
		if existingOrganization.ProposedChanges.Taggings != nil {
			bulkUpdateTagging := &model.InternalBulkUpdateTagging{}
			jsonTemp, _ := json.Marshal(map[string]interface{}{
				"IDs": funk.Map(
					existingOrganization.ProposedChanges.Taggings,
					func(tagging *model.Tagging) interface{} {
						return tagging.ID
					},
				),
			})
			json.Unmarshal(jsonTemp, bulkUpdateTagging)

			bulkUpdateTagging.RecentApprovingAccount = &model.ObjectIDOnly{
				ID: organizationToApprove.RecentApprovingAccount.ID,
			}
			bulkUpdateTagging.ProposalStatus = organizationToApprove.ProposalStatus

			_, err := approveUpdateOrgRepo.bulkApproveUpdateTaggingComponent.TransactionBody(
				operationOption,
				bulkUpdateTagging,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/approveUpdateOrganizationRepository",
					err,
				)
			}
		}
	}

	return approveUpdateOrgRepo.approveUpdateOrganizationTransactionComponent.TransactionBody(
		operationOption,
		organizationToApprove,
	)
}

func (approveUpdateOrgRepo *approveUpdateOrganizationRepository) RunTransaction(
	input *model.InternalUpdateOrganization,
) (*model.Organization, error) {
	output, err := approveUpdateOrgRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.Organization), err
}
