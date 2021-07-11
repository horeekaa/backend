package memberaccessdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessRepository struct {
	createMemberAccessTransactionComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessTransactionComponent
	mongoDBTransaction                     mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateMemberAccessRepository(
	createMemberAccessRepositoryTransactionComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository, error) {
	createMemberAccessRepo := &createMemberAccessRepository{
		createMemberAccessRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createMemberAccessRepo,
		"CreateMemberAccessRepository",
	)

	return createMemberAccessRepo, nil
}

func (createProdRepo *createMemberAccessRepository) SetValidation(
	usecaseComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessUsecaseComponent,
) (bool, error) {
	createProdRepo.createMemberAccessTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (createProdRepo *createMemberAccessRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createProdRepo.createMemberAccessTransactionComponent.PreTransaction(
		input.(*model.InternalCreateMemberAccess),
	)
}

func (createProdRepo *createMemberAccessRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return createProdRepo.createMemberAccessTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalCreateMemberAccess),
	)
}

func (createProdRepo *createMemberAccessRepository) RunTransaction(
	input *model.InternalCreateMemberAccess,
) (*model.MemberAccess, error) {
	output, err := createProdRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.MemberAccess), err
}
