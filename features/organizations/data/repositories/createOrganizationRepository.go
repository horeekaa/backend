package organizationdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createOrganizationRepository struct {
	createOrganizationTransactionComponent organizationdomainrepositoryinterfaces.CreateOrganizationTransactionComponent
	createDescriptivePhotoComponent        descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	createAddressComponent                 addressdomainrepositoryinterfaces.CreateAddressTransactionComponent
	bulkCreateTaggingComponent             taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent
	mongoDBTransaction                     mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateOrganizationRepository(
	createOrganizationRepositoryTransactionComponent organizationdomainrepositoryinterfaces.CreateOrganizationTransactionComponent,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	createAddressComponent addressdomainrepositoryinterfaces.CreateAddressTransactionComponent,
	bulkCreateTaggingComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (organizationdomainrepositoryinterfaces.CreateOrganizationRepository, error) {
	createOrganizationRepo := &createOrganizationRepository{
		createOrganizationRepositoryTransactionComponent,
		createDescriptivePhotoComponent,
		createAddressComponent,
		bulkCreateTaggingComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createOrganizationRepo,
		"CreateOrganizationRepository",
	)

	return createOrganizationRepo, nil
}

func (createOrgRepo *createOrganizationRepository) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.CreateOrganizationUsecaseComponent,
) (bool, error) {
	createOrgRepo.createOrganizationTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (createOrgRepo *createOrganizationRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createOrgRepo.createOrganizationTransactionComponent.PreTransaction(
		input.(*model.InternalCreateOrganization),
	)
}

func (createOrgRepo *createOrganizationRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	organizationToCreate := input.(*model.InternalCreateOrganization)
	generatedObjectID := createOrgRepo.createOrganizationTransactionComponent.GenerateNewObjectID()
	if organizationToCreate.ProfilePhotos != nil {
		savedPhotos := []*model.InternalCreateDescriptivePhoto{}
		for _, photo := range organizationToCreate.ProfilePhotos {
			photo.Category = model.DescriptivePhotoCategoryOrganizationProfile
			photo.Object = &model.ObjectIDOnly{
				ID: &generatedObjectID,
			}
			createdPhotoOutput, err := createOrgRepo.createDescriptivePhotoComponent.TransactionBody(
				operationOption,
				photo,
			)
			if err != nil {
				return nil, err
			}
			savedPhoto := &model.InternalCreateDescriptivePhoto{}
			jsonTemp, _ := json.Marshal(createdPhotoOutput)
			json.Unmarshal(jsonTemp, savedPhoto)
			savedPhotos = append(savedPhotos, savedPhoto)
		}
		organizationToCreate.ProfilePhotos = savedPhotos
	}

	if organizationToCreate.Addresses != nil {
		savedAddresses := []*model.InternalCreateAddress{}
		for _, address := range organizationToCreate.Addresses {
			address.Type = model.AddressTypeOrganization
			address.Object = &model.ObjectIDOnly{
				ID: &generatedObjectID,
			}
			createdAddressOutput, err := createOrgRepo.createAddressComponent.TransactionBody(
				operationOption,
				address,
			)
			if err != nil {
				return nil, err
			}
			savedAddress := &model.InternalCreateAddress{}
			jsonTemp, _ := json.Marshal(createdAddressOutput)
			json.Unmarshal(jsonTemp, savedAddress)
			savedAddresses = append(savedAddresses, savedAddress)
		}
		organizationToCreate.Addresses = savedAddresses
	}

	if organizationToCreate.Taggings != nil {
		savedTaggings := []*model.InternalCreateTagging{}
		for _, tagging := range organizationToCreate.Taggings {
			tagging.Organizations = []*model.ObjectIDOnly{
				{ID: &generatedObjectID},
			}
			tagging.ProposalStatus = organizationToCreate.ProposalStatus
			tagging.SubmittingAccount = organizationToCreate.SubmittingAccount
			tagging.IgnoreTaggedDocumentCheck = true

			createdTaggingOutput, err := createOrgRepo.bulkCreateTaggingComponent.TransactionBody(
				operationOption,
				tagging,
			)
			if err != nil {
				return nil, err
			}

			savedTagging := &model.InternalCreateTagging{}
			jsonTemp, _ := json.Marshal(createdTaggingOutput[0])
			json.Unmarshal(jsonTemp, savedTagging)
			savedTaggings = append(savedTaggings, savedTagging)
		}
		organizationToCreate.Taggings = savedTaggings
	}

	return createOrgRepo.createOrganizationTransactionComponent.TransactionBody(
		operationOption,
		organizationToCreate,
	)
}

func (createOrgRepo *createOrganizationRepository) RunTransaction(
	input *model.InternalCreateOrganization,
) (*model.Organization, error) {
	output, err := createOrgRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.Organization), err
}
