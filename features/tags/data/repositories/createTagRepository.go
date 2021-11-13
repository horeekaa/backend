package tagdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createTagRepository struct {
	createTagTransactionComponent   tagdomainrepositoryinterfaces.CreateTagTransactionComponent
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	mongoDBTransaction              mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateTagRepository(
	createTagRepositoryTransactionComponent tagdomainrepositoryinterfaces.CreateTagTransactionComponent,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (tagdomainrepositoryinterfaces.CreateTagRepository, error) {
	createTagRepo := &createTagRepository{
		createTagRepositoryTransactionComponent,
		createDescriptivePhotoComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createTagRepo,
		"CreateTagRepository",
	)

	return createTagRepo, nil
}

func (createTagRepo *createTagRepository) SetValidation(
	usecaseComponent tagdomainrepositoryinterfaces.CreateTagUsecaseComponent,
) (bool, error) {
	createTagRepo.createTagTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (createTagRepo *createTagRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createTagRepo.createTagTransactionComponent.PreTransaction(
		input.(*model.InternalCreateTag),
	)
}

func (createTagRepo *createTagRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	tagToCreate := input.(*model.InternalCreateTag)
	generatedObjectID := createTagRepo.createTagTransactionComponent.GenerateNewObjectID()
	if tagToCreate.Photos != nil {
		savedPhotos := []*model.InternalCreateDescriptivePhoto{}
		for _, photo := range tagToCreate.Photos {
			photo.Category = model.DescriptivePhotoCategoryTag
			photo.Object = &model.ObjectIDOnly{
				ID: &generatedObjectID,
			}
			photo.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*tagToCreate.ProposalStatus)
			photo.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*tagToCreate.SubmittingAccount)
			createdPhotoOutput, err := createTagRepo.createDescriptivePhotoComponent.TransactionBody(
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
		tagToCreate.Photos = savedPhotos
	}

	return createTagRepo.createTagTransactionComponent.TransactionBody(
		operationOption,
		tagToCreate,
	)
}

func (createTagRepo *createTagRepository) RunTransaction(
	input *model.InternalCreateTag,
) (*model.Tag, error) {
	output, err := createTagRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Tag), err
}
