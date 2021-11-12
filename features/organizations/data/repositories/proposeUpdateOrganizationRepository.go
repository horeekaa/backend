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

type proposeUpdateOrganizationRepository struct {
	organizationDataSource                        databaseorganizationdatasourceinterfaces.OrganizationDataSource
	createDescriptivePhotoComponent               descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	updateDescriptivePhotoComponent               descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent
	createAddressComponent                        addressdomainrepositoryinterfaces.CreateAddressTransactionComponent
	updateAddressComponent                        addressdomainrepositoryinterfaces.UpdateAddressTransactionComponent
	bulkCreateTaggingComponent                    taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent
	bulkUpdateTaggingComponent                    taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingTransactionComponent
	proposeUpdateOrganizationTransactionComponent organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationTransactionComponent
	mongoDBTransaction                            mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateOrganizationRepository(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	updateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent,
	createAddressComponent addressdomainrepositoryinterfaces.CreateAddressTransactionComponent,
	updateAddressComponent addressdomainrepositoryinterfaces.UpdateAddressTransactionComponent,
	bulkCreateTaggingComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent,
	bulkUpdateTaggingComponent taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingTransactionComponent,
	proposeUpdateOrganizationRepositoryTransactionComponent organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationRepository, error) {
	proposeUpdateOrganizationRepo := &proposeUpdateOrganizationRepository{
		organizationDataSource,
		createDescriptivePhotoComponent,
		updateDescriptivePhotoComponent,
		createAddressComponent,
		updateAddressComponent,
		bulkCreateTaggingComponent,
		bulkUpdateTaggingComponent,
		proposeUpdateOrganizationRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateOrganizationRepo,
		"ProposeUpdateOrganizationRepository",
	)

	return proposeUpdateOrganizationRepo, nil
}

func (updateOrgRepo *proposeUpdateOrganizationRepository) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationUsecaseComponent,
) (bool, error) {
	updateOrgRepo.proposeUpdateOrganizationTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateOrgRepo *proposeUpdateOrganizationRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.proposeUpdateOrganizationTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateOrganization),
	)
}

func (updateOrgRepo *proposeUpdateOrganizationRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	organizationToUpdate := input.(*model.InternalUpdateOrganization)
	existingOrganization, err := updateOrgRepo.organizationDataSource.GetMongoDataSource().FindByID(
		organizationToUpdate.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateOrganizationRepository",
			err,
		)
	}
	organizationToUpdateProposalStatus := *organizationToUpdate.ProposalStatus
	organizationToUpdateSubmittingAccount := *organizationToUpdate.SubmittingAccount

	if organizationToUpdate.ProfilePhotos != nil {
		savedPhotos := existingOrganization.ProfilePhotos
		for _, descPhotoToUpdate := range organizationToUpdate.ProfilePhotos {
			if descPhotoToUpdate.ID != nil {
				if !funk.Contains(
					existingOrganization.ProfilePhotos,
					func(dp *model.DescriptivePhoto) bool {
						return dp.ID == *descPhotoToUpdate.ID
					},
				) {
					continue
				}

				_, err := updateOrgRepo.updateDescriptivePhotoComponent.TransactionBody(
					operationOption,
					descPhotoToUpdate,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/proposeUpdateOrganizationRepository",
						err,
					)
				}
				continue
			}

			photoToCreate := &model.InternalCreateDescriptivePhoto{}
			jsonTemp, _ := json.Marshal(descPhotoToUpdate)
			json.Unmarshal(jsonTemp, photoToCreate)
			if descPhotoToUpdate.Photo != nil {
				photoToCreate.Photo.File = descPhotoToUpdate.Photo.File
			}
			photoToCreate.Category = model.DescriptivePhotoCategoryOrganizationProfile
			photoToCreate.Object = &model.ObjectIDOnly{
				ID: &existingOrganization.ID,
			}

			savedPhoto, err := updateOrgRepo.createDescriptivePhotoComponent.TransactionBody(
				operationOption,
				photoToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/proposeUpdateOrganizationRepository",
					err,
				)
			}
			savedPhotos = append(savedPhotos, savedPhoto)
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"ProfilePhotos": savedPhotos,
			},
		)
		json.Unmarshal(jsonTemp, organizationToUpdate)
	}

	if organizationToUpdate.Addresses != nil {
		savedAddresses := existingOrganization.Addresses
		for _, address := range organizationToUpdate.Addresses {
			if address.ID != nil {
				if !funk.Contains(
					existingOrganization.Addresses,
					func(ad *model.Address) bool {
						return ad.ID == *address.ID
					},
				) {
					continue
				}

				_, err := updateOrgRepo.updateAddressComponent.TransactionBody(
					operationOption,
					address,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/proposeUpdateOrganizationRepository",
						err,
					)
				}
				continue
			}

			addressToCreate := &model.InternalCreateAddress{}
			jsonTemp, _ := json.Marshal(address)
			json.Unmarshal(jsonTemp, addressToCreate)
			addressToCreate.Type = model.AddressTypeOrganization
			addressToCreate.Object = &model.ObjectIDOnly{
				ID: &existingOrganization.ID,
			}

			savedAddress, err := updateOrgRepo.createAddressComponent.TransactionBody(
				operationOption,
				addressToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/proposeUpdateOrganizationRepository",
					err,
				)
			}
			savedAddresses = append(savedAddresses, savedAddress)
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"Addresses": savedAddresses,
			},
		)
		json.Unmarshal(jsonTemp, organizationToUpdate)
	}

	if organizationToUpdate.Taggings != nil {
		savedTaggings := existingOrganization.Taggings
		for _, taggingToUpdate := range organizationToUpdate.Taggings {
			if taggingToUpdate.ID != nil {
				if !funk.Contains(
					existingOrganization.Taggings,
					func(pv *model.Tagging) bool {
						return pv.ID == *taggingToUpdate.ID
					},
				) {
					continue
				}

				bulkUpdateTagging := &model.InternalBulkUpdateTagging{}
				jsonTemp, _ := json.Marshal(taggingToUpdate)
				json.Unmarshal(jsonTemp, bulkUpdateTagging)
				jsonTemp, _ = json.Marshal(map[string]interface{}{
					"IDs": []interface{}{taggingToUpdate.ID},
				})
				json.Unmarshal(jsonTemp, bulkUpdateTagging)

				bulkUpdateTagging.ProposalStatus = &organizationToUpdateProposalStatus
				bulkUpdateTagging.SubmittingAccount = &organizationToUpdateSubmittingAccount

				_, err := updateOrgRepo.bulkUpdateTaggingComponent.TransactionBody(
					operationOption,
					bulkUpdateTagging,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/proposeUpdateOrganizationRepository",
						err,
					)
				}
				continue
			}

			taggingToCreate := &model.InternalCreateTagging{}
			jsonTemp, _ := json.Marshal(taggingToUpdate)
			json.Unmarshal(jsonTemp, taggingToCreate)
			taggingToCreate.Organizations = []*model.ObjectIDOnly{
				{ID: &existingOrganization.ID},
			}
			taggingToCreate.ProposalStatus = &organizationToUpdateProposalStatus
			taggingToCreate.SubmittingAccount = &organizationToUpdateSubmittingAccount

			savedTagging, err := updateOrgRepo.bulkCreateTaggingComponent.TransactionBody(
				operationOption,
				taggingToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/proposeUpdateOrganizationRepository",
					err,
				)
			}
			savedTaggings = append(savedTaggings, savedTagging...)
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"Taggings": savedTaggings,
			},
		)
		json.Unmarshal(jsonTemp, organizationToUpdate)
	}

	return updateOrgRepo.proposeUpdateOrganizationTransactionComponent.TransactionBody(
		operationOption,
		organizationToUpdate,
	)
}

func (updateOrgRepo *proposeUpdateOrganizationRepository) RunTransaction(
	input *model.InternalUpdateOrganization,
) (*model.Organization, error) {
	output, err := updateOrgRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.Organization), err
}
