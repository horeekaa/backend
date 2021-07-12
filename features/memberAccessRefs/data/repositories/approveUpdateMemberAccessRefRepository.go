package memberaccessrefdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateMemberAccessRefRepository struct {
	approveUpdateMemberAccessRefTransactionComponent memberaccessrefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefTransactionComponent
	mongoDBTransaction                               mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdateMemberAccessRefRepository(
	approveUpdateMemberAccessRefRepositoryTransactionComponent memberaccessrefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (memberaccessrefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefRepository, error) {
	approveUpdateMemberAccessRefRepo := &approveUpdateMemberAccessRefRepository{
		approveUpdateMemberAccessRefRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		approveUpdateMemberAccessRefRepo,
		"ApproveUpdateMemberAccessRefRepository",
	)

	return approveUpdateMemberAccessRefRepo, nil
}

func (updateOrgRepo *approveUpdateMemberAccessRefRepository) SetValidation(
	usecaseComponent memberaccessrefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefUsecaseComponent,
) (bool, error) {
	updateOrgRepo.approveUpdateMemberAccessRefTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateOrgRepo *approveUpdateMemberAccessRefRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.approveUpdateMemberAccessRefTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateMemberAccessRef),
	)
}

func (updateOrgRepo *approveUpdateMemberAccessRefRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.approveUpdateMemberAccessRefTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateMemberAccessRef),
	)
}

func (updateOrgRepo *approveUpdateMemberAccessRefRepository) RunTransaction(
	input *model.InternalUpdateMemberAccessRef,
) (*model.MemberAccessRef, error) {
	output, err := updateOrgRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.MemberAccessRef), err
}
