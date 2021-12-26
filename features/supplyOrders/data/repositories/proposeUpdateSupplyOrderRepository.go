package supplyorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdateSupplyOrderRepository struct {
	supplyOrderDataSource                        databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource
	proposeUpdateSupplyOrderTransactionComponent supplyorderdomainrepositoryinterfaces.ProposeUpdateSupplyOrderTransactionComponent
	createSupplyOrderItemComponent               supplyorderitemdomainrepositoryinterfaces.CreateSupplyOrderItemTransactionComponent
	proposeUpdateSupplyOrderItemComponent        supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemTransactionComponent
	approveUpdateSupplyOrderItemComponent        supplyorderitemdomainrepositoryinterfaces.ApproveUpdateSupplyOrderItemTransactionComponent
	mongoDBTransaction                           mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateSupplyOrderRepository(
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
	proposeUpdateSupplyOrderRepositoryTransactionComponent supplyorderdomainrepositoryinterfaces.ProposeUpdateSupplyOrderTransactionComponent,
	createSupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.CreateSupplyOrderItemTransactionComponent,
	proposeUpdateSupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemTransactionComponent,
	approveUpdateSupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.ApproveUpdateSupplyOrderItemTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (supplyorderdomainrepositoryinterfaces.ProposeUpdateSupplyOrderRepository, error) {
	proposeUpdateSupplyOrderRepo := &proposeUpdateSupplyOrderRepository{
		supplyOrderDataSource,
		proposeUpdateSupplyOrderRepositoryTransactionComponent,
		createSupplyOrderItemComponent,
		proposeUpdateSupplyOrderItemComponent,
		approveUpdateSupplyOrderItemComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateSupplyOrderRepo,
		"ProposeUpdateSupplyOrderRepository",
	)

	return proposeUpdateSupplyOrderRepo, nil
}

func (updateSupplyOrderRepo *proposeUpdateSupplyOrderRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateSupplyOrderRepo.proposeUpdateSupplyOrderTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateSupplyOrder),
	)
}

func (updateSupplyOrderRepo *proposeUpdateSupplyOrderRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	supplyOrderToUpdate := input.(*model.InternalUpdateSupplyOrder)
	existingsupplyOrder, err := updateSupplyOrderRepo.supplyOrderDataSource.GetMongoDataSource().FindByID(
		supplyOrderToUpdate.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateSupplyOrderRepository",
			err,
		)
	}

	if supplyOrderToUpdate.Items != nil {
		savedsupplyOrderItems := existingsupplyOrder.Items
		for _, supplyOrderItemToUpdate := range supplyOrderToUpdate.Items {
			if supplyOrderItemToUpdate.ID != nil {
				if !funk.Contains(
					existingsupplyOrder.Items,
					func(mi *model.SupplyOrderItem) bool {
						return mi.ID == *supplyOrderItemToUpdate.ID
					},
				) {
					continue
				}
				if supplyOrderItemToUpdate.ProposalStatus != nil {
					supplyOrderItemToUpdate.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
						return &m
					}(*supplyOrderToUpdate.SubmittingAccount)

					_, err := updateSupplyOrderRepo.approveUpdateSupplyOrderItemComponent.TransactionBody(
						operationOption,
						supplyOrderItemToUpdate,
					)
					if err != nil {
						return nil, horeekaacoreexceptiontofailure.ConvertException(
							"/proposeUpdateSupplyOrderRepository",
							err,
						)
					}
					continue
				}
				supplyOrderItemToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*supplyOrderToUpdate.ProposalStatus)
				supplyOrderItemToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*supplyOrderToUpdate.SubmittingAccount)
				if *supplyOrderToUpdate.MemberAccess.Organization.Type != model.OrganizationTypePartner {
					supplyOrderItemToUpdate.PartnerAgreed = func(b bool) *bool { return &b }(false)
				}

				_, err := updateSupplyOrderRepo.proposeUpdateSupplyOrderItemComponent.TransactionBody(
					operationOption,
					supplyOrderItemToUpdate,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/proposeUpdateSupplyOrderRepository",
						err,
					)
				}
				continue
			}

			supplyOrderItemToCreate := &model.InternalCreateSupplyOrderItem{}
			jsonTemp, _ := json.Marshal(supplyOrderItemToUpdate)
			json.Unmarshal(jsonTemp, supplyOrderItemToCreate)
			supplyOrderItemToCreate.SupplyOrder = &model.ObjectIDOnly{
				ID: &existingsupplyOrder.ID,
			}
			supplyOrderItemToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*supplyOrderToUpdate.ProposalStatus)
			supplyOrderItemToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*supplyOrderToUpdate.SubmittingAccount)
			if *supplyOrderToUpdate.MemberAccess.Organization.Type == model.OrganizationTypePartner {
				supplyOrderItemToCreate.PartnerAgreed = func(b bool) *bool { return &b }(true)
			}
			for i, descPhoto := range supplyOrderItemToUpdate.Photos {
				supplyOrderItemToCreate.Photos[i].Photo.File = descPhoto.Photo.File
			}

			savedSupplyOrderItem, err := updateSupplyOrderRepo.createSupplyOrderItemComponent.TransactionBody(
				operationOption,
				supplyOrderItemToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/proposeUpdateSupplyOrderRepository",
					err,
				)
			}
			savedsupplyOrderItems = append(savedsupplyOrderItems, savedSupplyOrderItem)
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"Items": savedsupplyOrderItems,
			},
		)
		json.Unmarshal(jsonTemp, supplyOrderToUpdate)
	}

	return updateSupplyOrderRepo.proposeUpdateSupplyOrderTransactionComponent.TransactionBody(
		operationOption,
		supplyOrderToUpdate,
	)
}

func (updateSupplyOrderRepo *proposeUpdateSupplyOrderRepository) RunTransaction(
	input *model.InternalUpdateSupplyOrder,
) (*model.SupplyOrder, error) {
	output, err := updateSupplyOrderRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.SupplyOrder), nil
}
