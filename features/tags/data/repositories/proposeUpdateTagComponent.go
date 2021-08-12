package tagdomainrepositories

import (
	"encoding/json"
	"fmt"
	"reflect"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateTagTransactionComponent struct {
	tagDataSource                    databasetagdatasourceinterfaces.TagDataSource
	loggingDataSource                databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility              coreutilityinterfaces.MapProcessorUtility
	structComparisonUtility          coreutilityinterfaces.StructComparisonUtility
	proposeUpdateTagUsecaseComponent tagdomainrepositoryinterfaces.ProposeUpdateTagUsecaseComponent
}

func NewProposeUpdateTagTransactionComponent(
	TagDataSource databasetagdatasourceinterfaces.TagDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
	structComparisonUtility coreutilityinterfaces.StructComparisonUtility,
) (tagdomainrepositoryinterfaces.ProposeUpdateTagTransactionComponent, error) {
	return &proposeUpdateTagTransactionComponent{
		tagDataSource:           TagDataSource,
		loggingDataSource:       loggingDataSource,
		mapProcessorUtility:     mapProcessorUtility,
		structComparisonUtility: structComparisonUtility,
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
	updateTag *model.InternalUpdateTag,
) (*model.Tag, error) {
	existingTag, err := updateTagTrx.tagDataSource.GetMongoDataSource().FindByID(
		updateTag.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateTag",
			err,
		)
	}
	fieldChanges := []*model.FieldChangeDataInput{}

	updateTagTrx.structComparisonUtility.SetComparisonFunc(
		func(tag interface{}, field1 interface{}, field2 interface{}, tagString *interface{}) {
			if field1 == field2 {
				return
			}
			*tagString = fmt.Sprintf(
				"%v%v",
				*tagString,
				tag,
			)

			fieldChanges = append(fieldChanges, &model.FieldChangeDataInput{
				Name:     fmt.Sprint(*tagString),
				Type:     reflect.TypeOf(field1).Kind().String(),
				OldValue: fmt.Sprint(field2),
				NewValue: fmt.Sprint(field1),
			})
			*tagString = ""
		},
	)
	updateTagTrx.structComparisonUtility.SetPreDeepComparisonFunc(
		func(tag interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v.",
				*tagString,
				tag,
			)
		},
	)
	var tagString interface{} = ""
	updateTagTrx.structComparisonUtility.CompareStructs(
		*updateTag,
		*existingTag,
		&tagString,
	)

	loggingOutput, err := updateTagTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Tag",
			Document: &model.ObjectIDOnly{
				ID: &existingTag.ID,
			},
			FieldChanges: fieldChanges,
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
			"/updateTag",
			err,
		)
	}
	updateTag.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

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
		}
	}

	updatedTag, err := updateTagTrx.tagDataSource.GetMongoDataSource().Update(
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
