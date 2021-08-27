package memberaccessdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateMemberAccessRepository struct {
	createNotifComponent                          notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	approveUpdateMemberAccessTransactionComponent memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessTransactionComponent
	mongoDBTransaction                            mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdateMemberAccessRepository(
	createNotifComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	approveUpdateMemberAccessRepositoryTransactionComponent memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessRepository, error) {
	approveUpdateMemberAccessRepo := &approveUpdateMemberAccessRepository{
		createNotifComponent,
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
	memberAccessToUpdate := input.(*model.InternalUpdateMemberAccess)
	updatedMemberAccess, err := approveUpdateMmbAccRepo.approveUpdateMemberAccessTransactionComponent.TransactionBody(
		operationOption,
		memberAccessToUpdate,
	)
	if err != nil {
		return nil, err
	}

	if updatedMemberAccess.MemberAccessRefType == model.MemberAccessRefTypeOrganizationsBased &&
		updatedMemberAccess.Account.ID.Hex() != updatedMemberAccess.SubmittingAccount.ID.Hex() &&
		*memberAccessToUpdate.ProposalStatus == model.EntityProposalStatusApproved &&
		!updatedMemberAccess.InvitationAccepted {
		notificationToCreate := &model.InternalCreateNotification{
			NotificationCategory: model.NotificationCategoryOrgInvitationRequest,
			PayloadOptions: &model.PayloadOptionsInput{
				InvitationRequestPayload: &model.InvitationRequestPayloadInput{
					MemberAccess: &model.MemberAccessForNotifPayloadInput{},
				},
			},
			RecipientAccount: &model.ObjectIDOnly{
				ID: &updatedMemberAccess.Account.ID,
			},
		}

		jsonTemp, _ := json.Marshal(updatedMemberAccess)
		json.Unmarshal(jsonTemp, &notificationToCreate.PayloadOptions.InvitationRequestPayload.MemberAccess)

		_, err := approveUpdateMmbAccRepo.createNotifComponent.TransactionBody(
			operationOption,
			notificationToCreate,
		)
		if err != nil {
			return nil, err
		}
	}
	return updatedMemberAccess, nil
}

func (approveUpdateMmbAccRepo *approveUpdateMemberAccessRepository) RunTransaction(
	input *model.InternalUpdateMemberAccess,
) (*model.MemberAccess, error) {
	output, err := approveUpdateMmbAccRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.MemberAccess), nil
}
