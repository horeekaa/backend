package organizationdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type approveUpdateOrganizationRepository struct {
	approveUpdateOrganizationTransactionComponent organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationTransactionComponent
	organizationDataSource                        databaseorganizationdatasourceinterfaces.OrganizationDataSource
	approveDescriptivePhotoComponent              descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent
	approveAddressComponent                       addressdomainrepositoryinterfaces.ApproveUpdateAddressTransactionComponent
	bulkApproveUpdateTaggingComponent             taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent
	mongoDBTransaction                            mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                                  string
}

func NewApproveUpdateOrganizationRepository(
	approveUpdateOrganizationRepositoryTransactionComponent organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationTransactionComponent,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	approveDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent,
	approveAddressComponent addressdomainrepositoryinterfaces.ApproveUpdateAddressTransactionComponent,
	bulkApproveUpdateTaggingComponent taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationRepository, error) {
	approveUpdateOrganizationRepo := &approveUpdateOrganizationRepository{
		approveUpdateOrganizationRepositoryTransactionComponent,
		organizationDataSource,
		approveDescriptivePhotoComponent,
		approveAddressComponent,
		bulkApproveUpdateTaggingComponent,
		mongoDBTransaction,
		"ApproveUpdateOrganizationRepository",
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
			approveUpdateOrgRepo.pathIdentity,
			err,
		)
	}

	if existingOrganization.ProposedChanges.ProposalStatus == model.EntityProposalStatusProposed {
		if existingOrganization.ProposedChanges.ProfilePhotos != nil {
			for _, photo := range existingOrganization.ProposedChanges.ProfilePhotos {
				updateDescriptivePhoto := &model.InternalUpdateDescriptivePhoto{
					ID: &photo.ID,
				}
				updateDescriptivePhoto.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*organizationToApprove.RecentApprovingAccount)
				updateDescriptivePhoto.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*organizationToApprove.ProposalStatus)

				_, err := approveUpdateOrgRepo.approveDescriptivePhotoComponent.TransactionBody(
					operationOption,
					updateDescriptivePhoto,
				)
				if err != nil {
					return nil, err
				}
			}
		}

		if existingOrganization.ProposedChanges.Addresses != nil {
			for _, address := range existingOrganization.ProposedChanges.Addresses {
				updateAddress := &model.InternalUpdateAddress{
					ID: &address.ID,
				}
				updateAddress.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*organizationToApprove.RecentApprovingAccount)
				updateAddress.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*organizationToApprove.ProposalStatus)

				_, err := approveUpdateOrgRepo.approveAddressComponent.TransactionBody(
					operationOption,
					updateAddress,
				)
				if err != nil {
					return nil, err
				}
			}
		}

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

			bulkUpdateTagging.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*organizationToApprove.RecentApprovingAccount)
			bulkUpdateTagging.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*organizationToApprove.ProposalStatus)

			_, err := approveUpdateOrgRepo.bulkApproveUpdateTaggingComponent.TransactionBody(
				operationOption,
				bulkUpdateTagging,
			)
			if err != nil {
				return nil, err
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
