package tagdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateTagRepository struct {
	approveUpdateTagTransactionComponent tagdomainrepositoryinterfaces.ApproveUpdateTagTransactionComponent
	approveDescriptivePhotoComponent     descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent
	tagDataSource                        databasetagdatasourceinterfaces.TagDataSource
	mongoDBTransaction                   mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdateTagRepository(
	approveUpdateTagRepositoryTransactionComponent tagdomainrepositoryinterfaces.ApproveUpdateTagTransactionComponent,
	approveDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent,
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (tagdomainrepositoryinterfaces.ApproveUpdateTagRepository, error) {
	approveUpdatetagRepo := &approveUpdateTagRepository{
		approveUpdateTagRepositoryTransactionComponent,
		approveDescriptivePhotoComponent,
		tagDataSource,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		approveUpdatetagRepo,
		"ApproveUpdateTagRepository",
	)

	return approveUpdatetagRepo, nil
}

func (updateTagRepo *approveUpdateTagRepository) SetValidation(
	usecaseComponent tagdomainrepositoryinterfaces.ApproveUpdateTagUsecaseComponent,
) (bool, error) {
	updateTagRepo.approveUpdateTagTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateTagRepo *approveUpdateTagRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateTagRepo.approveUpdateTagTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateTag),
	)
}

func (updateTagRepo *approveUpdateTagRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	tagToApprove := input.(*model.InternalUpdateTag)
	existingTag, err := updateTagRepo.tagDataSource.GetMongoDataSource().FindByID(
		tagToApprove.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdateTagRepository",
			err,
		)
	}

	if existingTag.ProposedChanges.ProposalStatus == model.EntityProposalStatusProposed {
		if existingTag.ProposedChanges.Photos != nil {
			for _, photo := range existingTag.ProposedChanges.Photos {
				updateDescriptivePhoto := &model.InternalUpdateDescriptivePhoto{
					ID: &photo.ID,
				}
				updateDescriptivePhoto.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*tagToApprove.RecentApprovingAccount)
				updateDescriptivePhoto.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*tagToApprove.ProposalStatus)

				_, err := updateTagRepo.approveDescriptivePhotoComponent.TransactionBody(
					operationOption,
					updateDescriptivePhoto,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/approveUpdateTagRepository",
						err,
					)
				}
			}
		}
	}

	return updateTagRepo.approveUpdateTagTransactionComponent.TransactionBody(
		operationOption,
		tagToApprove,
	)
}

func (updateTagRepo *approveUpdateTagRepository) RunTransaction(
	input *model.InternalUpdateTag,
) (*model.Tag, error) {
	output, err := updateTagRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Tag), err
}
