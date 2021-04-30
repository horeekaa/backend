package memberaccessrefdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories/types"
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
	mongoDBTransaction.SetTransaction(
		updateMemberAccessRefRepositoryTransactionComponent,
		"UpdateMemberAccessRefRepository",
	)

	return &updateMemberAccessRefRepository{
		updateMemberAccessRefRepositoryTransactionComponent: updateMemberAccessRefRepositoryTransactionComponent,
		mongoDBTransaction: mongoDBTransaction,
	}, nil
}

func (updateMmbAccRefRepo *updateMemberAccessRefRepository) SetValidation(
	usecaseComponent memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefUsecaseComponent,
) (bool, error) {
	updateMmbAccRefRepo.updateMemberAccessRefRepositoryTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateMmbAccRefRepo *updateMemberAccessRefRepository) RunTransaction(
	input *model.UpdateMemberAccessRef,
) (*memberaccessrefdomainrepositorytypes.UpdateMemberAccessRefOutput, error) {
	output, err := updateMmbAccRefRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*memberaccessrefdomainrepositorytypes.UpdateMemberAccessRefOutput), err
}
