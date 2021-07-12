package memberaccessdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateMemberAccessRepository struct {
	approveUpdateMemberAccessTransactionComponent memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessTransactionComponent
	mongoDBTransaction                            mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdateMemberAccessRepository(
	approveUpdateMemberAccessRepositoryTransactionComponent memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessRepository, error) {
	approveUpdateMemberAccessRepo := &approveUpdateMemberAccessRepository{
		approveUpdateMemberAccessRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		approveUpdateMemberAccessRepo,
		"ApproveUpdateMemberAccessRepository",
	)

	return approveUpdateMemberAccessRepo, nil
}

func (approveUpdateMmbAccRepo *approveUpdateMemberAccessRepository) SetValidation(
	usecaseComponent memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessUsecaseComponent,
) (bool, error) {
	approveUpdateMmbAccRepo.approveUpdateMemberAccessTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (approveUpdateMmbAccRepo *approveUpdateMemberAccessRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return approveUpdateMmbAccRepo.approveUpdateMemberAccessTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateMemberAccess),
	)
}

func (approveUpdateMmbAccRepo *approveUpdateMemberAccessRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return approveUpdateMmbAccRepo.approveUpdateMemberAccessTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateMemberAccess),
	)
}

func (approveUpdateMmbAccRepo *approveUpdateMemberAccessRepository) RunTransaction(
	input *model.InternalUpdateMemberAccess,
) (*model.MemberAccess, error) {
	output, err := approveUpdateMmbAccRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.MemberAccess), err
}
