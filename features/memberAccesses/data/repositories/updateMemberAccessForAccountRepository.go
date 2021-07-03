package memberaccessdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type updateMemberAccessForAccountRepository struct {
	updateMemberAccessRepositoryTransactionComponent memberaccessdomainrepositoryinterfaces.UpdateMemberAccessForAccountTransactionComponent
	mongoDBTransaction                               mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewUpdateMemberAccessForAccountRepository(
	updateMemberAccessRepositoryTransactionComponent memberaccessdomainrepositoryinterfaces.UpdateMemberAccessForAccountTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (memberaccessdomainrepositoryinterfaces.UpdateMemberAccessForAccountRepository, error) {
	updateMmbAccRepo := &updateMemberAccessForAccountRepository{
		updateMemberAccessRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		updateMmbAccRepo,
		"UpdateMemberAccessRepository",
	)

	return updateMmbAccRepo, nil
}

func (updateMmbAccRepo *updateMemberAccessForAccountRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateMmbAccRepo.updateMemberAccessRepositoryTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateMemberAccess),
	)
}

func (updateMmbAccRepo *updateMemberAccessForAccountRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return updateMmbAccRepo.updateMemberAccessRepositoryTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateMemberAccess),
	)
}

func (updateMmbAccRepo *updateMemberAccessForAccountRepository) RunTransaction(
	input *model.InternalUpdateMemberAccess,
) (*memberaccessdomainrepositorytypes.UpdateMemberAccessOutput, error) {
	output, err := updateMmbAccRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*memberaccessdomainrepositorytypes.UpdateMemberAccessOutput), err
}
