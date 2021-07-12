package memberaccessrefdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessRefRepository struct {
	createMemberAccessRefTransactionComponent memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefTransactionComponent
	mongoDBTransaction                        mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateMemberAccessRefRepository(
	createMemberAccessRefRepositoryTransactionComponent memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefRepository, error) {
	createMemberAccessRefRepo := &createMemberAccessRefRepository{
		createMemberAccessRefRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createMemberAccessRefRepo,
		"CreateMemberAccessRefRepository",
	)

	return createMemberAccessRefRepo, nil
}

func (createProdRepo *createMemberAccessRefRepository) SetValidation(
	usecaseComponent memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefUsecaseComponent,
) (bool, error) {
	createProdRepo.createMemberAccessRefTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (createProdRepo *createMemberAccessRefRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createProdRepo.createMemberAccessRefTransactionComponent.PreTransaction(
		input.(*model.InternalCreateMemberAccessRef),
	)
}

func (createProdRepo *createMemberAccessRefRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return createProdRepo.createMemberAccessRefTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalCreateMemberAccessRef),
	)
}

func (createProdRepo *createMemberAccessRefRepository) RunTransaction(
	input *model.InternalCreateMemberAccessRef,
) (*model.MemberAccessRef, error) {
	output, err := createProdRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.MemberAccessRef), err
}
