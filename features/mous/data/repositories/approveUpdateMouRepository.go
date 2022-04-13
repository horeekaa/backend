package moudomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateMouRepository struct {
	approveUpdateMouTransactionComponent moudomainrepositoryinterfaces.ApproveUpdateMouTransactionComponent
	mouDataSource                        databasemoudatasourceinterfaces.MouDataSource
	approveUpdateMouItemComponent        mouitemdomainrepositoryinterfaces.ApproveUpdateMouItemTransactionComponent
	createNotificationComponent          notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	mongoDBTransaction                   mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                         string
}

func NewApproveUpdateMouRepository(
	approveUpdateMouTransactionComponent moudomainrepositoryinterfaces.ApproveUpdateMouTransactionComponent,
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	approveUpdateMouItemComponent mouitemdomainrepositoryinterfaces.ApproveUpdateMouItemTransactionComponent,
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (moudomainrepositoryinterfaces.ApproveUpdateMouRepository, error) {
	approveUpdateMouRepo := &approveUpdateMouRepository{
		approveUpdateMouTransactionComponent,
		mouDataSource,
		approveUpdateMouItemComponent,
		createNotificationComponent,
		mongoDBTransaction,
		"ApproveUpdateMouRepository",
	}

	mongoDBTransaction.SetTransaction(
		approveUpdateMouRepo,
		"ApproveUpdateMouRepository",
	)

	return approveUpdateMouRepo, nil
}

func (approveUpdateMouRepo *approveUpdateMouRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return input, nil
}

func (approveUpdateMouRepo *approveUpdateMouRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	mouToApprove := input.(*model.InternalUpdateMou)
	existingMou, err := approveUpdateMouRepo.mouDataSource.GetMongoDataSource().FindByID(
		mouToApprove.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveUpdateMouRepo.pathIdentity,
			err,
		)
	}

	if existingMou.ProposedChanges.Items != nil {
		for _, item := range existingMou.ProposedChanges.Items {
			updateItem := &model.InternalUpdateMouItem{
				ID: &item.ID,
			}
			updateItem.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*mouToApprove.RecentApprovingAccount)
			updateItem.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*mouToApprove.ProposalStatus)

			_, err := approveUpdateMouRepo.approveUpdateMouItemComponent.TransactionBody(
				operationOption,
				updateItem,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	return approveUpdateMouRepo.approveUpdateMouTransactionComponent.TransactionBody(
		operationOption,
		mouToApprove,
	)
}

func (approveUpdateMouRepo *approveUpdateMouRepository) RunTransaction(
	input *model.InternalUpdateMou,
) (*model.Mou, error) {
	output, err := approveUpdateMouRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}

	approvedMou := output.(*model.Mou)
	go func() {
		notificationToCreate := &model.InternalCreateNotification{
			NotificationCategory: model.NotificationCategoryMouApproval,
			PayloadOptions: &model.PayloadOptionsInput{
				MouPayload: &model.MouPayloadInput{
					Mou: &model.MouForNotifPayloadInput{},
				},
			},
			RecipientAccount: &model.ObjectIDOnly{
				ID: &approvedMou.SecondParty.AccountInCharge.ID,
			},
		}

		jsonTemp, _ := json.Marshal(approvedMou)
		json.Unmarshal(jsonTemp, &notificationToCreate.PayloadOptions.MouPayload.Mou)

		_, err = approveUpdateMouRepo.createNotificationComponent.TransactionBody(
			&mongodbcoretypes.OperationOptions{},
			notificationToCreate,
		)
		if err != nil {
			return
		}
	}()

	return approvedMou, err
}
