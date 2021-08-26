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

	if createdMemberAccess.MemberAccessRefType == model.MemberAccessRefTypeOrganizationsBased &&
		createdMemberAccess.Account.ID.Hex() != createdMemberAccess.SubmittingAccount.ID.Hex() &&
		createdMemberAccess.ProposalStatus == model.EntityProposalStatusApproved {
		notificationToCreate := &model.InternalCreateNotification{
			NotificationCategory: model.NotificationCategoryOrgInvitationRequest,
			PayloadOptions: &model.PayloadOptionsInput{
				InvitationRequestPayload: &model.InvitationRequestPayloadInput{
					MemberAccess: &model.MemberAccessForNotifPayloadInput{},
				},
			},
		}

		jsonTemp, _ := json.Marshal(createdMemberAccess)
		json.Unmarshal(jsonTemp, &notificationToCreate.PayloadOptions.InvitationRequestPayload.MemberAccess)
		_, err := createProdRepo.createNotifComponent.TransactionBody(
			operationOption,
			notificationToCreate,
		)
		if err != nil {
			return nil, err
		}
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
	return (output).(*model.MemberAccess), nil
}
