package memberaccessdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateMemberAccessRepository struct {
	proposeUpdateMemberAccessTransactionComponent memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessTransactionComponent
	mongoDBTransaction                            mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateMemberAccessRepository(
	proposeUpdateMemberAccessRepositoryTransactionComponent memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessRepository, error) {
	proposeUpdateMemberAccessRepo := &proposeUpdateMemberAccessRepository{
		proposeUpdateMemberAccessRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateMemberAccessRepo,
		"ProposeUpdateMemberAccessRepository",
	)

	return proposeUpdateMemberAccessRepo, nil
}

func (updateOrgRepo *proposeUpdateMemberAccessRepository) SetValidation(
	usecaseComponent memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessUsecaseComponent,
) (bool, error) {
	updateOrgRepo.proposeUpdateMemberAccessTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateOrgRepo *proposeUpdateMemberAccessRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.proposeUpdateMemberAccessTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateMemberAccess),
	)
}

func (updateOrgRepo *proposeUpdateMemberAccessRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.proposeUpdateMemberAccessTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateMemberAccess),
	)
}

func (updateOrgRepo *proposeUpdateMemberAccessRepository) RunTransaction(
	input *model.InternalUpdateMemberAccess,
) (*model.MemberAccess, error) {
	output, err := updateOrgRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.MemberAccess), err
}
