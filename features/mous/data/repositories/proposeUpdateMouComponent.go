package moudomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moudomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateMouTransactionComponent struct {
	mouDataSource       databasemoudatasourceinterfaces.MouDataSource
	loggingDataSource   databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility
	partyLoader         moudomainrepositoryutilityinterfaces.PartyLoader
	pathIdentity        string
}

func NewProposeUpdateMouTransactionComponent(
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
	partyLoader moudomainrepositoryutilityinterfaces.PartyLoader,
) (moudomainrepositoryinterfaces.ProposeUpdateMouTransactionComponent, error) {
	return &proposeUpdateMouTransactionComponent{
		mouDataSource:       mouDataSource,
		loggingDataSource:   loggingDataSource,
		mapProcessorUtility: mapProcessorUtility,
		partyLoader:         partyLoader,
		pathIdentity:        "ProposeUpdateMouComponent",
	}, nil
}

func (updateMouTrx *proposeUpdateMouTransactionComponent) PreTransaction(
	input *model.InternalUpdateMou,
) (*model.InternalUpdateMou, error) {
	return input, nil
}

func (updateMouTrx *proposeUpdateMouTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateMou,
) (*model.Mou, error) {
	updateMou := &model.DatabaseUpdateMou{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateMou)

	existingMou, err := updateMouTrx.mouDataSource.GetMongoDataSource().FindByID(
		updateMou.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateMouTrx.pathIdentity,
			err,
		)
	}

	newDocumentJson, _ := json.Marshal(*updateMou)
	oldDocumentJson, _ := json.Marshal(*existingMou)
	loggingOutput, err := updateMouTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Mou",
			Document: &model.ObjectIDOnly{
				ID: &existingMou.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateMou.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateMou.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateMouTrx.pathIdentity,
			err,
		)
	}
	updateMou.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	_, err = updateMouTrx.partyLoader.TransactionBody(
		session,
		input.FirstParty,
		updateMou.FirstParty,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateMouTrx.pathIdentity,
			err,
		)
	}
	_, err = updateMouTrx.partyLoader.TransactionBody(
		session,
		input.SecondParty,
		updateMou.SecondParty,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateMouTrx.pathIdentity,
			err,
		)
	}

	fieldsToUpdateMou := &model.DatabaseUpdateMou{
		ID: updateMou.ID,
	}
	jsonExisting, _ := json.Marshal(existingMou)
	json.Unmarshal(jsonExisting, &fieldsToUpdateMou.ProposedChanges)

	var updateMouMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateMou)
	json.Unmarshal(jsonUpdate, &updateMouMap)

	updateMouTrx.mapProcessorUtility.RemoveNil(updateMouMap)

	jsonUpdate, _ = json.Marshal(updateMouMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateMou.ProposedChanges)

	if updateMou.ProposalStatus != nil {
		fieldsToUpdateMou.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateMou.SubmittingAccount.ID,
		}
		if *updateMou.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdateMou)
		}
	}

	updatedMou, err := updateMouTrx.mouDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateMou.ID,
		},
		fieldsToUpdateMou,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateMouTrx.pathIdentity,
			err,
		)
	}

	return updatedMou, nil
}
