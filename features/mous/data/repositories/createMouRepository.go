package moudomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createMouRepository struct {
	createMouTransactionComponent moudomainrepositoryinterfaces.CreateMouTransactionComponent
	mongoDBTransaction            mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateMouRepository(
	createMouRepositoryTransactionComponent moudomainrepositoryinterfaces.CreateMouTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (moudomainrepositoryinterfaces.CreateMouRepository, error) {
	createMouRepo := &createMouRepository{
		createMouRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createMouRepo,
		"CreateMouRepository",
	)

	return createMouRepo, nil
}

func (createMouRepo *createMouRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createMouRepo.createMouTransactionComponent.PreTransaction(
		input.(*model.InternalCreateMou),
	)
}

func (createMouRepo *createMouRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	mouToCreate := input.(*model.InternalCreateMou)

	return createMouRepo.createMouTransactionComponent.TransactionBody(
		operationOption,
		mouToCreate,
	)
}

func (createMouRepo *createMouRepository) RunTransaction(
	input *model.InternalCreateMou,
) (*model.Mou, error) {
	output, err := createMouRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.Mou), nil
}
