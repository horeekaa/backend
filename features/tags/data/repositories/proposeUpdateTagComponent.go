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

type proposeUpdateTagTransactionComponent struct {
	tagDataSource                    databasetagdatasourceinterfaces.TagDataSource
	taggingDataSource                databasetaggingdatasourceinterfaces.TaggingDataSource
	loggingDataSource                databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility              coreutilityinterfaces.MapProcessorUtility
	proposeUpdateTagUsecaseComponent tagdomainrepositoryinterfaces.ProposeUpdateTagUsecaseComponent
	pathIdentity                     string
}

func NewProposeUpdateTagTransactionComponent(
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (tagdomainrepositoryinterfaces.ProposeUpdateTagTransactionComponent, error) {
	return &proposeUpdateTagTransactionComponent{
		tagDataSource:       tagDataSource,
		taggingDataSource:   taggingDataSource,
		loggingDataSource:   loggingDataSource,
		mapProcessorUtility: mapProcessorUtility,
		pathIdentity:        "ProposeUpdateTagComponent",
	}, nil
}

func (updateTagTrx *proposeUpdateTagTransactionComponent) SetValidation(
	usecaseComponent tagdomainrepositoryinterfaces.ProposeUpdateTagUsecaseComponent,
) (bool, error) {
	updateTagTrx.proposeUpdateTagUsecaseComponent = usecaseComponent
	return true, nil
}

func (updateTagTrx *proposeUpdateTagTransactionComponent) PreTransaction(
	input *model.InternalUpdateTag,
) (*model.InternalUpdateTag, error) {
	if updateTagTrx.proposeUpdateTagUsecaseComponent == nil {
		return input, nil
	}
	return updateTagTrx.proposeUpdateTagUsecaseComponent.Validation(input)
}

func (updateTagTrx *proposeUpdateTagTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateTag,
) (*model.Tag, error) {
	updateTag := &model.DatabaseUpdateTag{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateTag)

	existingTag, err := updateTagTrx.tagDataSource.GetMongoDataSource().FindByID(
		updateTag.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateTagTrx.pathIdentity,
			err,
		)
	}
	newDocumentJson, _ := json.Marshal(*updateTag)
	oldDocumentJson, _ := json.Marshal(*existingTag)
	loggingOutput, err := updateTagTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Tag",
			Document: &model.ObjectIDOnly{
				ID: &existingTag.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateTag.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateTag.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateTagTrx.pathIdentity,
			err,
		)
	}
	updateTag.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	currentTime := time.Now().UTC()
	updateTag.UpdatedAt = &currentTime

	fieldsToUpdateTag := &model.DatabaseUpdateTag{
		ID: updateTag.ID,
	}
	jsonExisting, _ := json.Marshal(existingTag)
	json.Unmarshal(jsonExisting, &fieldsToUpdateTag.ProposedChanges)

	var updateTagMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateTag)
	json.Unmarshal(jsonUpdate, &updateTagMap)

	updateTagTrx.mapProcessorUtility.RemoveNil(updateTagMap)

	jsonUpdate, _ = json.Marshal(updateTagMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateTag.ProposedChanges)

	if updateTag.ProposalStatus != nil {
		fieldsToUpdateTag.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateTag.SubmittingAccount.ID,
		}
		if *updateTag.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdateTag)

			tagForTagging := &model.TagForTaggingInput{}
			jsonTag, _ := json.Marshal(fieldsToUpdateTag.ProposedChanges)
			json.Unmarshal(jsonTag, tagForTagging)

			updateTagTrx.taggingDataSource.GetMongoDataSource().Update(
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

	updatedTag, err := updateTagTrx.tagDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateTag.ID,
		},
		fieldsToUpdateTag,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateTagTrx.pathIdentity,
			err,
		)
	}

	return updatedTag, nil
}
