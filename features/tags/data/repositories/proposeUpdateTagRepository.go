package tagdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdateTagRepository struct {
	tagDataSource                          databasetagdatasourceinterfaces.TagDataSource
	createDescriptivePhotoComponent        descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	proposeUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent
	proposeUpdateTagTransactionComponent   tagdomainrepositoryinterfaces.ProposeUpdateTagTransactionComponent
	mongoDBTransaction                     mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateTagRepository(
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	proposeUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent,
	proposeUpdateTagRepositoryTransactionComponent tagdomainrepositoryinterfaces.ProposeUpdateTagTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (tagdomainrepositoryinterfaces.ProposeUpdateTagRepository, error) {
	proposeUpdateTagRepo := &proposeUpdateTagRepository{
		tagDataSource,
		createDescriptivePhotoComponent,
		proposeUpdateDescriptivePhotoComponent,
		proposeUpdateTagRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateTagRepo,
		"ProposeUpdateTagRepository",
	)

	return proposeUpdateTagRepo, nil
}

func (updateTagRepo *proposeUpdateTagRepository) SetValidation(
	usecaseComponent tagdomainrepositoryinterfaces.ProposeUpdateTagUsecaseComponent,
) (bool, error) {
	updateTagRepo.proposeUpdateTagTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateTagRepo *proposeUpdateTagRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateTagRepo.proposeUpdateTagTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateTag),
	)
}

func (updateTagRepo *proposeUpdateTagRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	tagToUpdate := input.(*model.InternalUpdateTag)
	existingTag, err := updateTagRepo.tagDataSource.GetMongoDataSource().FindByID(
		tagToUpdate.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateTagRepository",
			err,
		)
	}

	if tagToUpdate.Photos != nil {
		savedPhotos := existingTag.Photos
		for _, descPhotoToUpdate := range tagToUpdate.Photos {
			if descPhotoToUpdate.ID != nil {
				if !funk.Contains(
					existingTag.Photos,
					func(dp *model.DescriptivePhoto) bool {
						return dp.ID == *descPhotoToUpdate.ID
					},
				) {
					continue
				}
				descPhotoToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*tagToUpdate.ProposalStatus)
				descPhotoToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*tagToUpdate.SubmittingAccount)

				_, err := updateTagRepo.proposeUpdateDescriptivePhotoComponent.TransactionBody(
					operationOption,
					descPhotoToUpdate,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/proposeUpdateTagRepository",
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
			photoToCreate.Category = model.DescriptivePhotoCategoryTag
			photoToCreate.Object = &model.ObjectIDOnly{
				ID: &existingTag.ID,
			}
			photoToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*tagToUpdate.ProposalStatus)
			photoToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*tagToUpdate.SubmittingAccount)

			savedPhoto, err := updateTagRepo.createDescriptivePhotoComponent.TransactionBody(
				operationOption,
				photoToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/proposeUpdateTagRepository",
					err,
				)
			}
			savedPhotos = append(savedPhotos, savedPhoto)
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"Photos": savedPhotos,
			},
		)
		json.Unmarshal(jsonTemp, tagToUpdate)
	}

	return updateTagRepo.proposeUpdateTagTransactionComponent.TransactionBody(
		operationOption,
		tagToUpdate,
	)
}

func (updateTagRepo *proposeUpdateTagRepository) RunTransaction(
	input *model.InternalUpdateTag,
) (*model.Tag, error) {
	output, err := updateTagRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Tag), err
}
