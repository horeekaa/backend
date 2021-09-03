package moudomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateMouRepository struct {
	proposeUpdateMouTransactionComponent moudomainrepositoryinterfaces.ProposeUpdateMouTransactionComponent
	mongoDBTransaction                   mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateMouRepository(
	proposeUpdateMouRepositoryTransactionComponent moudomainrepositoryinterfaces.ProposeUpdateMouTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (moudomainrepositoryinterfaces.ProposeUpdateMouRepository, error) {
	proposeUpdateMouRepo := &proposeUpdateMouRepository{
		proposeUpdateMouRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateMouRepo,
		"ProposeUpdateMouRepository",
	)

	return proposeUpdateMouRepo, nil
}

func (updateMouRepo *proposeUpdateMouRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateMouRepo.proposeUpdateMouTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateMou),
	)
}

func (updateMouRepo *proposeUpdateMouRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	mouToUpdate := input.(*model.InternalUpdateMou)

	return updateMouRepo.proposeUpdateMouTransactionComponent.TransactionBody(
		operationOption,
		mouToUpdate,
	)
}

func (updateMouRepo *proposeUpdateMouRepository) RunTransaction(
	input *model.InternalUpdateMou,
) (*model.Mou, error) {
	output, err := updateMouRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.Mou), nil
}
