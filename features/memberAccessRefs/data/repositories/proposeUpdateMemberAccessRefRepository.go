package memberaccessrefdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateMemberAccessRefRepository struct {
	proposeUpdateMemberAccessRefTransactionComponent memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefTransactionComponent
	mongoDBTransaction                               mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateMemberAccessRefRepository(
	proposeUpdateMemberAccessRefRepositoryTransactionComponent memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefRepository, error) {
	proposeUpdateMemberAccessRefRepo := &proposeUpdateMemberAccessRefRepository{
		proposeUpdateMemberAccessRefRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateMemberAccessRefRepo,
		"ProposeUpdateMemberAccessRefRepository",
	)

	return proposeUpdateMemberAccessRefRepo, nil
}

func (updateOrgRepo *proposeUpdateMemberAccessRefRepository) SetValidation(
	usecaseComponent memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefUsecaseComponent,
) (bool, error) {
	updateOrgRepo.proposeUpdateMemberAccessRefTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateOrgRepo *proposeUpdateMemberAccessRefRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.proposeUpdateMemberAccessRefTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateMemberAccessRef),
	)
}

func (updateOrgRepo *proposeUpdateMemberAccessRefRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.proposeUpdateMemberAccessRefTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateMemberAccessRef),
	)
}

func (updateOrgRepo *proposeUpdateMemberAccessRefRepository) RunTransaction(
	input *model.InternalUpdateMemberAccessRef,
) (*model.MemberAccessRef, error) {
	output, err := updateOrgRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.MemberAccessRef), err
}
