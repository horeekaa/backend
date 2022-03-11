package moudomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdateMouRepository struct {
	mouDataSource                        databasemoudatasourceinterfaces.MouDataSource
	proposeUpdateMouTransactionComponent moudomainrepositoryinterfaces.ProposeUpdateMouTransactionComponent
	createMouItemComponent               mouitemdomainrepositoryinterfaces.CreateMouItemTransactionComponent
	proposeUpdateMouItemComponent        mouitemdomainrepositoryinterfaces.ProposeUpdateMouItemTransactionComponent
	mongoDBTransaction                   mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                         string
}

func NewProposeUpdateMouRepository(
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	proposeUpdateMouRepositoryTransactionComponent moudomainrepositoryinterfaces.ProposeUpdateMouTransactionComponent,
	createMouItemComponent mouitemdomainrepositoryinterfaces.CreateMouItemTransactionComponent,
	proposeUpdateMouItemComponent mouitemdomainrepositoryinterfaces.ProposeUpdateMouItemTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (moudomainrepositoryinterfaces.ProposeUpdateMouRepository, error) {
	proposeUpdateMouRepo := &proposeUpdateMouRepository{
		mouDataSource,
		proposeUpdateMouRepositoryTransactionComponent,
		createMouItemComponent,
		proposeUpdateMouItemComponent,
		mongoDBTransaction,
		"ProposeUpdateMouRepository",
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateMouRepo,
		"ProposeUpdateMouRepository",
	)

	return proposeUpdateMouRepo, nil
}

func (updateMouRepo *proposeUpdateMouRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateMouRepo.proposeUpdateMouTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateMou),
	)
}

func (updateMouRepo *proposeUpdateMouRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	mouToUpdate := input.(*model.InternalUpdateMou)
	existingMou, err := updateMouRepo.mouDataSource.GetMongoDataSource().FindByID(
		mouToUpdate.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateMouRepo.pathIdentity,
			err,
		)
	}

	if mouToUpdate.Items != nil {
		savedMouItems := existingMou.Items
		for _, mouItemToUpdate := range mouToUpdate.Items {
			if mouItemToUpdate.ID != nil {
				if !funk.Contains(
					existingMou.Items,
					func(mi *model.MouItem) bool {
						return mi.ID == *mouItemToUpdate.ID
					},
				) {
					continue
				}
				mouItemToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*mouToUpdate.ProposalStatus)
				mouItemToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*mouToUpdate.SubmittingAccount)

				_, err := updateMouRepo.proposeUpdateMouItemComponent.TransactionBody(
					operationOption,
					mouItemToUpdate,
				)
				if err != nil {
					return nil, err
				}

				if mouItemToUpdate.IsActive != nil {
					if !*mouItemToUpdate.IsActive {
						index := funk.IndexOf(
							savedMouItems,
							func(mi *model.MouItem) bool {
								return mi.ID == *mouItemToUpdate.ID
							},
						)
						if index > -1 {
							savedMouItems = append(savedMouItems[:index], savedMouItems[index+1:]...)
						}
					}
				}
				continue
			}

			mouItemToCreate := &model.InternalCreateMouItem{}
			jsonTemp, _ := json.Marshal(mouItemToUpdate)
			json.Unmarshal(jsonTemp, mouItemToCreate)
			mouItemToCreate.Mou = &model.ObjectIDOnly{
				ID: &existingMou.ID,
			}
			mouItemToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*mouToUpdate.ProposalStatus)
			mouItemToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*mouToUpdate.SubmittingAccount)

			savedMouItem, err := updateMouRepo.createMouItemComponent.TransactionBody(
				operationOption,
				mouItemToCreate,
			)
			if err != nil {
				return nil, err
			}
			savedMouItems = append(savedMouItems, savedMouItem)
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"Items": savedMouItems,
			},
		)
		json.Unmarshal(jsonTemp, mouToUpdate)
	}

	return updateMouRepo.proposeUpdateMouTransactionComponent.TransactionBody(
		operationOption,
		mouToUpdate,
	)
}

func (updateMouRepo *proposeUpdateMouRepository) RunTransaction(
	input *model.InternalUpdateMou,
) (*model.Mou, error) {
	output, err := updateMouRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.Mou), nil
}
