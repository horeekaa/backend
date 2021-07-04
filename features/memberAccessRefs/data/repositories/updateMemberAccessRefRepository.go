package memberaccessrefdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type updateMemberAccessRefRepository struct {
	updateMemberAccessRefRepositoryTransactionComponent memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefTransactionComponent
	mongoDBTransaction                                  mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewUpdateMemberAccessRefRepository(
	updateMemberAccessRefRepositoryTransactionComponent memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefRepository, error) {
	updateMmbAccRefRepo := &updateMemberAccessRefRepository{
		updateMemberAccessRefRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		updateMmbAccRefRepo,
		"UpdateMemberAccessRefRepository",
	)

	return updateMmbAccRefRepo, nil
}

func (updateMmbAccRefRepo *updateMemberAccessRefRepository) SetValidation(
	usecaseComponent memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefUsecaseComponent,
) (bool, error) {
	updateMmbAccRefRepo.updateMemberAccessRefRepositoryTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateMmbAccRefRepo *updateMemberAccessRefRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateMmbAccRefRepo.updateMemberAccessRefRepositoryTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateMemberAccessRef),
	)
}

func (updateMmbAccRefRepo *updateMemberAccessRefRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return updateMmbAccRefRepo.updateMemberAccessRefRepositoryTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateMemberAccessRef),
	)
}

func (updateMmbAccRefRepo *updateMemberAccessRefRepository) RunTransaction(
	input *model.InternalUpdateMemberAccessRef,
) (*model.MemberAccessRef, error) {
	output, err := updateMmbAccRefRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.MemberAccessRef), err
}
