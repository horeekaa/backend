package memberaccessdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateMemberAccessRepository struct {
	memberAccessDataSource                        databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	createNotifComponent                          notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	approveUpdateMemberAccessTransactionComponent memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessTransactionComponent
	mongoDBTransaction                            mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                                  string
}

func NewApproveUpdateMemberAccessRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	createNotifComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	approveUpdateMemberAccessRepositoryTransactionComponent memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessRepository, error) {
	approveUpdateMemberAccessRepo := &approveUpdateMemberAccessRepository{
		memberAccessDataSource,
		createNotifComponent,
		approveUpdateMemberAccessRepositoryTransactionComponent,
		mongoDBTransaction,
		"ApproveUpdateMemberAccessRepository",
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
	return approveUpdateMmbAccRepo.approveUpdateMemberAccessTransactionComponent.TransactionBody(
		operationOption,
		memberAccessToUpdate,
	)
}

func (approveUpdateMmbAccRepo *approveUpdateMemberAccessRepository) RunTransaction(
	input *model.InternalUpdateMemberAccess,
) (*model.MemberAccess, error) {
	existingMemberAccess, err := approveUpdateMmbAccRepo.memberAccessDataSource.GetMongoDataSource().FindByID(
		input.ID,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveUpdateMmbAccRepo.pathIdentity,
			err,
		)
	}

	output, err := approveUpdateMmbAccRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}

	updatedMemberAccess := (output).(*model.MemberAccess)
	go func() {
		if updatedMemberAccess.MemberAccessRefType == model.MemberAccessRefTypeOrganizationsBased &&
			*input.ProposalStatus == model.EntityProposalStatusApproved &&
			!existingMemberAccess.InvitationAccepted {
			notificationToCreate := &model.InternalCreateNotification{
				NotificationCategory: model.NotificationCategoryMemberAccessInvitationRequest,
				PayloadOptions: &model.PayloadOptionsInput{
					MemberAccessInvitationPayload: &model.MemberAccessInvitationPayloadInput{
						MemberAccess: &model.MemberAccessForNotifPayloadInput{},
					},
				},
				RecipientAccount: &model.ObjectIDOnly{
					ID: &updatedMemberAccess.Account.ID,
				},
			}

			jsonTemp, _ := json.Marshal(updatedMemberAccess)
			json.Unmarshal(jsonTemp, &notificationToCreate.PayloadOptions.MemberAccessInvitationPayload.MemberAccess)

			_, err := approveUpdateMmbAccRepo.createNotifComponent.TransactionBody(
				&mongodbcoretypes.OperationOptions{},
				notificationToCreate,
			)
			if err != nil {
				return
			}
		}
	}()
	return updatedMemberAccess, nil
}
