package moudomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateMouRepository struct {
	approveUpdateMouTransactionComponent moudomainrepositoryinterfaces.ApproveUpdateMouTransactionComponent
	mouDataSource                        databasemoudatasourceinterfaces.MouDataSource
	approveUpdateMouItemComponent        mouitemdomainrepositoryinterfaces.ApproveUpdateMouItemTransactionComponent
	mongoDBTransaction                   mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdateMouRepository(
	approveUpdateMouTransactionComponent moudomainrepositoryinterfaces.ApproveUpdateMouTransactionComponent,
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	approveUpdateMouItemComponent mouitemdomainrepositoryinterfaces.ApproveUpdateMouItemTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (moudomainrepositoryinterfaces.ApproveUpdateMouRepository, error) {
	approveUpdateMouRepo := &approveUpdateMouRepository{
		approveUpdateMouTransactionComponent,
		mouDataSource,
		approveUpdateMouItemComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		approveUpdateMouRepo,
		"ApproveUpdateMouRepository",
	)

	return approveUpdateMouRepo, nil
}

func (approveUpdateMouRepo *approveUpdateMouRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return input, nil
}

func (approveUpdateMouRepo *approveUpdateMouRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	mouToApprove := input.(*model.InternalUpdateMou)
	existingMou, err := approveUpdateMouRepo.mouDataSource.GetMongoDataSource().FindByID(
		mouToApprove.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdateMouRepository",
			err,
		)
	}

	if existingMou.ProposedChanges.Items != nil {
		for _, item := range existingMou.ProposedChanges.Items {
			updateItem := &model.InternalUpdateMouItem{
				ID: &item.ID,
			}
			updateItem.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*mouToApprove.RecentApprovingAccount)
			updateItem.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*mouToApprove.ProposalStatus)

			_, err := approveUpdateMouRepo.approveUpdateMouItemComponent.TransactionBody(
				operationOption,
				updateItem,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/approveUpdateMouRepository",
					err,
				)
			}
		}
	}

	return approveUpdateMouRepo.approveUpdateMouTransactionComponent.TransactionBody(
		operationOption,
		mouToApprove,
	)
}

func (approveUpdateMouRepo *approveUpdateMouRepository) RunTransaction(
	input *model.InternalUpdateMou,
) (*model.Mou, error) {
	output, err := approveUpdateMouRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.Mou), err
}
