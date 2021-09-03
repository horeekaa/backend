package moudomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateMouRepository struct {
	approveUpdateMouTransactionComponent moudomainrepositoryinterfaces.ApproveUpdateMouTransactionComponent
	mongoDBTransaction                   mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdateMouRepository(
	approveUpdateMouTransactionComponent moudomainrepositoryinterfaces.ApproveUpdateMouTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (moudomainrepositoryinterfaces.ApproveUpdateMouRepository, error) {
	approveUpdateMouRepo := &approveUpdateMouRepository{
		approveUpdateMouTransactionComponent,
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
