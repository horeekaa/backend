package memberaccessdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createMemberAccessRepository struct {
	createNotifComponent                   notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	createMemberAccessTransactionComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessTransactionComponent
	mongoDBTransaction                     mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateMemberAccessRepository(
	createNotifComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	createMemberAccessRepositoryTransactionComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository, error) {
	createMemberAccessRepo := &createMemberAccessRepository{
		createNotifComponent,
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
	memberAccessToCreate := input.(*model.InternalCreateMemberAccess)
	createdMemberAccess, err := createProdRepo.createMemberAccessTransactionComponent.TransactionBody(
		operationOption,
		memberAccessToCreate,
	)
	if err != nil {
		return nil, err
	}

	return createdMemberAccess, nil
}

func (createProdRepo *createMemberAccessRepository) RunTransaction(
	input *model.InternalCreateMemberAccess,
) (*model.MemberAccess, error) {
	output, err := createProdRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}

	createdMemberAccess := (output).(*model.MemberAccess)
	go func() {
		if createdMemberAccess.MemberAccessRefType == model.MemberAccessRefTypeOrganizationsBased &&
			createdMemberAccess.ProposalStatus == model.EntityProposalStatusApproved {
			notificationToCreate := &model.InternalCreateNotification{
				NotificationCategory: model.NotificationCategoryMemberAccessInvitationRequest,
				PayloadOptions: &model.PayloadOptionsInput{
					MemberAccessInvitationPayload: &model.MemberAccessInvitationPayloadInput{
						MemberAccess: &model.MemberAccessForNotifPayloadInput{},
					},
				},
				RecipientAccount: &model.ObjectIDOnly{
					ID: &createdMemberAccess.Account.ID,
				},
			}

			jsonTemp, _ := json.Marshal(createdMemberAccess)
			json.Unmarshal(jsonTemp, &notificationToCreate.PayloadOptions.MemberAccessInvitationPayload.MemberAccess)
			_, err := createProdRepo.createNotifComponent.TransactionBody(
				&mongodbcoretypes.OperationOptions{},
				notificationToCreate,
			)
			if err != nil {
				return
			}
		}
	}()
	return createdMemberAccess, nil
}
