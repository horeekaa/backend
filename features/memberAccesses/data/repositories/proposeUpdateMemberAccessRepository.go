package memberaccessdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdateMemberAccessRepository struct {
	createNotifComponent                          notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	proposeUpdateMemberAccessTransactionComponent memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessTransactionComponent
	mongoDBTransaction                            mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateMemberAccessRepository(
	createNotifComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	proposeUpdateMemberAccessRepositoryTransactionComponent memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessRepository, error) {
	proposeUpdateMemberAccessRepo := &proposeUpdateMemberAccessRepository{
		createNotifComponent,
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
	memberAccessToUpdate := input.(*model.InternalUpdateMemberAccess)
	updatedMemberAccess, err := updateOrgRepo.proposeUpdateMemberAccessTransactionComponent.TransactionBody(
		operationOption,
		memberAccessToUpdate,
	)
	if err != nil {
		return nil, err
	}

	if updatedMemberAccess.MemberAccessRefType == model.MemberAccessRefTypeOrganizationsBased &&
		updatedMemberAccess.Account.ID.Hex() == memberAccessToUpdate.SubmittingAccount.ID.Hex() &&
		updatedMemberAccess.ProposedChanges.ProposalStatus == model.EntityProposalStatusApproved &&
		(funk.GetOrElse(
			memberAccessToUpdate.InvitationAccepted,
			false,
		)).(bool) {
		notificationToCreate := &model.InternalCreateNotification{
			NotificationCategory: model.NotificationCategoryOrgInvitationAccepted,
			PayloadOptions: &model.PayloadOptionsInput{
				InvitationAcceptedPayload: &model.InvitationAcceptedPayloadInput{
					MemberAccess: &model.MemberAccessForNotifPayloadInput{},
				},
			},
			RecipientAccount: &model.ObjectIDOnly{
				ID: &updatedMemberAccess.SubmittingAccount.ID,
			},
		}

		jsonTemp, _ := json.Marshal(updatedMemberAccess)
		json.Unmarshal(jsonTemp, &notificationToCreate.PayloadOptions.InvitationAcceptedPayload.MemberAccess)

		_, err := updateOrgRepo.createNotifComponent.TransactionBody(
			operationOption,
			notificationToCreate,
		)
		if err != nil {
			return nil, err
		}
	}
	return updatedMemberAccess, nil
}

func (updateOrgRepo *proposeUpdateMemberAccessRepository) RunTransaction(
	input *model.InternalUpdateMemberAccess,
) (*model.MemberAccess, error) {
	output, err := updateOrgRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.MemberAccess), nil
}
