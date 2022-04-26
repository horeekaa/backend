package supplyorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateSupplyOrderRepository struct {
	memberAccessDataSource                       databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	supplyOrderDataSource                        databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource
	approveUpdateDescriptivePhotoComponent       descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent
	approveUpdatesupplyOrderItemComponent        supplyorderitemdomainrepositoryinterfaces.ApproveUpdateSupplyOrderItemTransactionComponent
	approveUpdateSupplyOrderTransactionComponent supplyorderdomainrepositoryinterfaces.ApproveUpdateSupplyOrderTransactionComponent
	createNotificationComponent                  notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	mongoDBTransaction                           mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                                 string
}

func NewApproveUpdateSupplyOrderRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
	approveUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent,
	approveUpdatesupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.ApproveUpdateSupplyOrderItemTransactionComponent,
	approveUpdateSupplyOrderTransactionComponent supplyorderdomainrepositoryinterfaces.ApproveUpdateSupplyOrderTransactionComponent,
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (supplyorderdomainrepositoryinterfaces.ApproveUpdateSupplyOrderRepository, error) {
	approveUpdateSupplyOrderRepo := &approveUpdateSupplyOrderRepository{
		memberAccessDataSource,
		supplyOrderDataSource,
		approveUpdateDescriptivePhotoComponent,
		approveUpdatesupplyOrderItemComponent,
		approveUpdateSupplyOrderTransactionComponent,
		createNotificationComponent,
		mongoDBTransaction,
		"ApproveUpdateSupplyOrderRepository",
	}

	mongoDBTransaction.SetTransaction(
		approveUpdateSupplyOrderRepo,
		"ApproveUpdateSupplyOrderRepository",
	)

	return approveUpdateSupplyOrderRepo, nil
}

func (approveUpdateSupplyOrderRepo *approveUpdateSupplyOrderRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return input, nil
}

func (approveUpdateSupplyOrderRepo *approveUpdateSupplyOrderRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	supplyOrderToApprove := input.(*model.InternalUpdateSupplyOrder)
	existingSupplyOrder, err := approveUpdateSupplyOrderRepo.supplyOrderDataSource.GetMongoDataSource().FindByID(
		supplyOrderToApprove.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveUpdateSupplyOrderRepo.pathIdentity,
			err,
		)
	}
	if existingSupplyOrder.ProposedChanges.ProposalStatus == model.EntityProposalStatusProposed {
		if existingSupplyOrder.ProposedChanges.PaymentProofPhoto != nil {
			updateDescriptivePhoto := &model.InternalUpdateDescriptivePhoto{
				ID: &existingSupplyOrder.ProposedChanges.PaymentProofPhoto.ID,
			}
			updateDescriptivePhoto.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*supplyOrderToApprove.RecentApprovingAccount)
			updateDescriptivePhoto.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*supplyOrderToApprove.ProposalStatus)

			_, err := approveUpdateSupplyOrderRepo.approveUpdateDescriptivePhotoComponent.TransactionBody(
				operationOption,
				updateDescriptivePhoto,
			)
			if err != nil {
				return nil, err
			}
		}

		if existingSupplyOrder.ProposedChanges.Items != nil {
			for _, soItem := range existingSupplyOrder.ProposedChanges.Items {
				updateSupplyOrderItem := &model.InternalUpdateSupplyOrderItem{
					ID: &soItem.ID,
				}
				updateSupplyOrderItem.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*supplyOrderToApprove.RecentApprovingAccount)
				updateSupplyOrderItem.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*supplyOrderToApprove.ProposalStatus)

				_, err := approveUpdateSupplyOrderRepo.approveUpdatesupplyOrderItemComponent.TransactionBody(
					operationOption,
					updateSupplyOrderItem,
				)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return approveUpdateSupplyOrderRepo.approveUpdateSupplyOrderTransactionComponent.TransactionBody(
		operationOption,
		supplyOrderToApprove,
	)
}

func (approveUpdateSupplyOrderRepo *approveUpdateSupplyOrderRepository) RunTransaction(
	input *model.InternalUpdateSupplyOrder,
) (*model.SupplyOrder, error) {
	output, err := approveUpdateSupplyOrderRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}

	approvedSupplyOrder := (output).(*model.SupplyOrder)
	go func() {
		memberAccesses, err := approveUpdateSupplyOrderRepo.memberAccessDataSource.GetMongoDataSource().Find(
			map[string]interface{}{
				"memberAccessRefType": model.MemberAccessRefTypeOrganizationsBased,
				"status":              model.MemberAccessStatusActive,
				"proposalStatus":      model.EntityProposalStatusApproved,
				"invitationAccepted":  true,
				"organization._id":    approvedSupplyOrder.Organization.ID,
			},
			&mongodbcoretypes.PaginationOptions{},
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			return
		}

		for _, memberAccess := range memberAccesses {
			notificationToCreate := &model.InternalCreateNotification{
				NotificationCategory: model.NotificationCategorySupplyOrderApproval,
				PayloadOptions: &model.PayloadOptionsInput{
					SupplyOrderPayload: &model.SupplyOrderPayloadInput{
						SupplyOrder: &model.SupplyOrderForNotifPayloadInput{},
					},
				},
				RecipientAccount: &model.ObjectIDOnly{
					ID: &memberAccess.Account.ID,
				},
			}

			jsonTemp, _ := json.Marshal(approvedSupplyOrder)
			json.Unmarshal(jsonTemp, &notificationToCreate.PayloadOptions.SupplyOrderPayload.SupplyOrder)

			_, err = approveUpdateSupplyOrderRepo.createNotificationComponent.TransactionBody(
				&mongodbcoretypes.OperationOptions{},
				notificationToCreate,
			)
			if err != nil {
				return
			}
		}
	}()

	return approvedSupplyOrder, err
}
