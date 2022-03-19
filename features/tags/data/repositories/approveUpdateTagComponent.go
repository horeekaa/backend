package tagdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateTagTransactionComponent struct {
	tagDataSource                    databasetagdatasourceinterfaces.TagDataSource
	taggingDataSource                databasetaggingdatasourceinterfaces.TaggingDataSource
	loggingDataSource                databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility              coreutilityinterfaces.MapProcessorUtility
	approveUpdateTagUsecaseComponent tagdomainrepositoryinterfaces.ApproveUpdateTagUsecaseComponent
	pathIdentity                     string
}

func NewApproveUpdateTagTransactionComponent(
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (tagdomainrepositoryinterfaces.ApproveUpdateTagTransactionComponent, error) {
	return &approveUpdateTagTransactionComponent{
		tagDataSource:       tagDataSource,
		taggingDataSource:   taggingDataSource,
		loggingDataSource:   loggingDataSource,
		mapProcessorUtility: mapProcessorUtility,
		pathIdentity:        "ApproveUpdateTagComponent",
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
	input *model.InternalUpdateTag,
) (*model.Tag, error) {
	updateTag := &model.DatabaseUpdateTag{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateTag)

	existingTag, err := approveTagTrx.tagDataSource.GetMongoDataSource().FindByID(
		updateTag.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveTagTrx.pathIdentity,
			err,
		)
	}

	previousLog, err := approveTagTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingTag.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveTagTrx.pathIdentity,
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
	jsonTemp, _ = json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveTagTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveTagTrx.pathIdentity,
			err,
		)
	}

	updateTag.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	currentTime := time.Now()
	updateTag.UpdatedAt = &currentTime

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

			tagForTagging := &model.TagForTaggingInput{}
			json.Unmarshal(jsonUpdate, tagForTagging)

			approveTagTrx.taggingDataSource.GetMongoDataSource().Update(
				map[string]interface{}{
					"tag._id": existingTag.ID,
				},
				&model.DatabaseUpdateTagging{
					Tag: tagForTagging,
				},
				session,
			)
		}
	}

	updatedTag, err := approveTagTrx.tagDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateTag.ID,
		},
		fieldsToUpdateTag,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveTagTrx.pathIdentity,
			err,
		)
	}

	return updatedTag, nil
}
