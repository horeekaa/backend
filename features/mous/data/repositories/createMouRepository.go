package moudomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createMouRepository struct {
	createMouTransactionComponent moudomainrepositoryinterfaces.CreateMouTransactionComponent
	createMouItemComponent        mouitemdomainrepositoryinterfaces.CreateMouItemTransactionComponent
	mongoDBTransaction            mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateMouRepository(
	createMouRepositoryTransactionComponent moudomainrepositoryinterfaces.CreateMouTransactionComponent,
	createMouItemComponent mouitemdomainrepositoryinterfaces.CreateMouItemTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (moudomainrepositoryinterfaces.CreateMouRepository, error) {
	createMouRepo := &createMouRepository{
		createMouRepositoryTransactionComponent,
		createMouItemComponent,
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
	generatedObjectID := createMouRepo.createMouTransactionComponent.GenerateNewObjectID()

	if mouToCreate.Items != nil {
		savedMouItems := []*model.InternalCreateMouItem{}
		for _, mouItem := range mouToCreate.Items {
			mouItem.Mou = &model.ObjectIDOnly{
				ID: &generatedObjectID,
			}
			mouItem.Organization = &model.OrganizationForMouItemInput{
				ID: *mouToCreate.SecondParty.Organization.ID,
			}
			mouItem.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*mouToCreate.ProposalStatus)
			mouItem.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*mouToCreate.SubmittingAccount)
			createdMouItemOutput, err := createMouRepo.createMouItemComponent.TransactionBody(
				operationOption,
				mouItem,
			)
			if err != nil {
				return nil, err
			}
			savedMouItem := &model.InternalCreateMouItem{}
			jsonTemp, _ := json.Marshal(createdMouItemOutput)
			json.Unmarshal(jsonTemp, savedMouItem)
			savedMouItems = append(savedMouItems, savedMouItem)
		}
		mouToCreate.Items = savedMouItems
	}

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
