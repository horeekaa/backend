package mouitemdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	mouitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateMouItemTransactionComponent struct {
	mouItemDataSource   databasemouitemdatasourceinterfaces.MouItemDataSource
	loggingDataSource   databaseloggingdatasourceinterfaces.LoggingDataSource
	agreedProductLoader mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility
	pathIdentity        string
}

func NewProposeUpdateMouItemTransactionComponent(
	mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	agreedProductLoader mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (mouitemdomainrepositoryinterfaces.ProposeUpdateMouItemTransactionComponent, error) {
	return &proposeUpdateMouItemTransactionComponent{
		mouItemDataSource:   mouItemDataSource,
		loggingDataSource:   loggingDataSource,
		agreedProductLoader: agreedProductLoader,
		mapProcessorUtility: mapProcessorUtility,
		pathIdentity:        "ProposeUpdateMouItemComponent",
	}, nil
}

func (updateMouItemTrx *proposeUpdateMouItemTransactionComponent) PreTransaction(
	input *model.InternalUpdateMouItem,
) (*model.InternalUpdateMouItem, error) {
	return input, nil
}

func (updateMouItemTrx *proposeUpdateMouItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateMouItem,
) (*model.MouItem, error) {
	updateMouItem := &model.DatabaseUpdateMouItem{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateMouItem)

	existingMouItem, err := updateMouItemTrx.mouItemDataSource.GetMongoDataSource().FindByID(
		updateMouItem.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateMouItemTrx.pathIdentity,
			err,
		)
	}

	if updateMouItem.Product != nil {
		updateMouItemTrx.agreedProductLoader.TransactionBody(
			session,
			updateMouItem.Product,
			updateMouItem.AgreedProduct,
		)
	}

	newDocumentJson, _ := json.Marshal(*updateMouItem)
	oldDocumentJson, _ := json.Marshal(*existingMouItem)
	loggingOutput, err := updateMouItemTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "MouItem",
			Document: &model.ObjectIDOnly{
				ID: &existingMouItem.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateMouItem.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateMouItem.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateMouItemTrx.pathIdentity,
			err,
		)
	}
	updateMouItem.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	currentTime := time.Now()
	updateMouItem.UpdatedAt = &currentTime

	fieldsToUpdateMouItem := &model.DatabaseUpdateMouItem{
		ID: updateMouItem.ID,
	}
	jsonExisting, _ := json.Marshal(existingMouItem)
	json.Unmarshal(jsonExisting, &fieldsToUpdateMouItem.ProposedChanges)

	var updateMouItemMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateMouItem)
	json.Unmarshal(jsonUpdate, &updateMouItemMap)

	updateMouItemTrx.mapProcessorUtility.RemoveNil(updateMouItemMap)

	jsonUpdate, _ = json.Marshal(updateMouItemMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateMouItem.ProposedChanges)

	if updateMouItem.ProposalStatus != nil {
		fieldsToUpdateMouItem.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateMouItem.SubmittingAccount.ID,
		}
		if *updateMouItem.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdateMouItem)
		}
	}

	updatedMouItem, err := updateMouItemTrx.mouItemDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateMouItem.ID,
		},
		fieldsToUpdateMouItem,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateMouItemTrx.pathIdentity,
			err,
		)
	}

	return updatedMouItem, nil
}
