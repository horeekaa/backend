package tagdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateTagTransactionComponent struct {
	tagDataSource                    databasetagdatasourceinterfaces.TagDataSource
	loggingDataSource                databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility              coreutilityinterfaces.MapProcessorUtility
	approveUpdateTagUsecaseComponent tagdomainrepositoryinterfaces.ApproveUpdateTagUsecaseComponent
}

func NewApproveUpdateTagTransactionComponent(
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (tagdomainrepositoryinterfaces.ApproveUpdateTagTransactionComponent, error) {
	return &approveUpdateTagTransactionComponent{
		tagDataSource:       tagDataSource,
		loggingDataSource:   loggingDataSource,
		mapProcessorUtility: mapProcessorUtility,
	}, nil
}

func (approveTagTrx *approveUpdateTagTransactionComponent) SetValidation(
	usecaseComponent tagdomainrepositoryinterfaces.ApproveUpdateTagUsecaseComponent,
) (bool, error) {
	approveTagTrx.approveUpdateTagUsecaseComponent = usecaseComponent
	return true, nil
}

func (approveTagTrx *approveUpdateTagTransactionComponent) PreTransaction(
	input *model.InternalUpdateTag,
) (*model.InternalUpdateTag, error) {
	if approveTagTrx.approveUpdateTagUsecaseComponent == nil {
		return input, nil
	}
	return approveTagTrx.approveUpdateTagUsecaseComponent.Validation(input)
}

func (approveTagTrx *approveUpdateTagTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateTag *model.InternalUpdateTag,
) (*model.Tag, error) {
	existingTag, err := approveTagTrx.tagDataSource.GetMongoDataSource().FindByID(
		updateTag.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateTag",
			err,
		)
	}

	previousLog, err := approveTagTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingTag.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateTag",
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateTag.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateTag.ProposalStatus,
	}
	jsonTemp, _ := json.Marshal(
		map[string]interface{}{
			"FieldChanges": previousLog.FieldChanges,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveTagTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)

	updateTag.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdateTag := &model.DatabaseUpdateTag{
		ID: updateTag.ID,
	}
	jsonExisting, _ := json.Marshal(existingTag.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdateTag.ProposedChanges)

	var updateTagMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateTag)
	json.Unmarshal(jsonUpdate, &updateTagMap)

	approveTagTrx.mapProcessorUtility.RemoveNil(updateTagMap)

	jsonUpdate, _ = json.Marshal(updateTagMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateTag.ProposedChanges)

	if updateTag.ProposalStatus != nil {
		if *updateTag.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdateTag.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdateTag)
		}
	}

	updatedTag, err := approveTagTrx.tagDataSource.GetMongoDataSource().Update(
		fieldsToUpdateTag.ID,
		fieldsToUpdateTag,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateTag",
			err,
		)
	}

	return updatedTag, nil
}
