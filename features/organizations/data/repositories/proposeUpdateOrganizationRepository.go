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
	proposeUpdateDescriptivePhotoComponent        descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent
	createAddressComponent                        addressdomainrepositoryinterfaces.CreateAddressTransactionComponent
	updateAddressComponent                        addressdomainrepositoryinterfaces.ProposeUpdateAddressTransactionComponent
	bulkCreateTaggingComponent                    taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent
	bulkUpdateTaggingComponent                    taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingTransactionComponent
	proposeUpdateOrganizationTransactionComponent organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationTransactionComponent
	mongoDBTransaction                            mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateOrganizationRepository(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	proposeUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent,
	createAddressComponent addressdomainrepositoryinterfaces.CreateAddressTransactionComponent,
	proposeUpdateAddressComponent addressdomainrepositoryinterfaces.ProposeUpdateAddressTransactionComponent,
	bulkCreateTaggingComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent,
	bulkUpdateTaggingComponent taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingTransactionComponent,
	proposeUpdateOrganizationRepositoryTransactionComponent organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationRepository, error) {
	proposeUpdateOrganizationRepo := &proposeUpdateOrganizationRepository{
		organizationDataSource,
		createDescriptivePhotoComponent,
		proposeUpdateDescriptivePhotoComponent,
		createAddressComponent,
		proposeUpdateAddressComponent,
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

				descPhotoToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*organizationToUpdate.ProposalStatus)
				descPhotoToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*organizationToUpdate.SubmittingAccount)

				_, err := updateOrgRepo.proposeUpdateDescriptivePhotoComponent.TransactionBody(
					operationOption,
					descPhotoToUpdate,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/proposeUpdateOrganizationRepository",
						err,
					)
				}

				if descPhotoToUpdate.IsActive != nil {
					if !*descPhotoToUpdate.IsActive {
						index := funk.IndexOf(
							savedPhotos,
							func(dp *model.DescriptivePhoto) bool {
								return dp.ID == *descPhotoToUpdate.ID
							},
						)
						if index > -1 {
							savedPhotos = append(savedPhotos[:index], savedPhotos[index+1:]...)
						}
					}
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
			photoToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*organizationToUpdate.ProposalStatus)
			photoToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*organizationToUpdate.SubmittingAccount)

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
		for _, addressToUpdate := range organizationToUpdate.Addresses {
			if addressToUpdate.ID != nil {
				if !funk.Contains(
					existingOrganization.Addresses,
					func(ad *model.Address) bool {
						return ad.ID == *addressToUpdate.ID
					},
				) {
					continue
				}

				addressToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*organizationToUpdate.ProposalStatus)
				addressToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*organizationToUpdate.SubmittingAccount)

				_, err := updateOrgRepo.updateAddressComponent.TransactionBody(
					operationOption,
					addressToUpdate,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/proposeUpdateOrganizationRepository",
						err,
					)
				}

				if addressToUpdate.IsActive != nil {
					if !*addressToUpdate.IsActive {
						index := funk.IndexOf(
							savedAddresses,
							func(ad *model.Address) bool {
								return ad.ID == *addressToUpdate.ID
							},
						)
						if index > -1 {
							savedAddresses = append(savedAddresses[:index], savedAddresses[index+1:]...)
						}
					}
				}
				continue
			}

			addressToCreate := &model.InternalCreateAddress{}
			jsonTemp, _ := json.Marshal(addressToUpdate)
			json.Unmarshal(jsonTemp, addressToCreate)
			addressToCreate.Type = model.AddressTypeOrganization
			addressToCreate.Object = &model.ObjectIDOnly{
				ID: &existingOrganization.ID,
			}
			addressToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*organizationToUpdate.ProposalStatus)
			addressToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*organizationToUpdate.SubmittingAccount)

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

				bulkUpdateTagging.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*organizationToUpdate.ProposalStatus)
				bulkUpdateTagging.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*organizationToUpdate.SubmittingAccount)

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

				if taggingToUpdate.IsActive != nil {
					if !*taggingToUpdate.IsActive {
						index := funk.IndexOf(
							savedTaggings,
							func(tg *model.Tagging) bool {
								return tg.ID == *taggingToUpdate.ID
							},
						)
						if index > -1 {
							savedTaggings = append(savedTaggings[:index], savedTaggings[index+1:]...)
						}
					}
				}
				continue
			}

			taggingToCreate := &model.InternalCreateTagging{}
			jsonTemp, _ := json.Marshal(taggingToUpdate)
			json.Unmarshal(jsonTemp, taggingToCreate)
			taggingToCreate.Organizations = []*model.ObjectIDOnly{
				{ID: &existingOrganization.ID},
			}
			taggingToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*organizationToUpdate.ProposalStatus)
			taggingToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*organizationToUpdate.SubmittingAccount)

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
