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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createTagTransactionComponent struct {
	tagDataSource              databasetagdatasourceinterfaces.TagDataSource
	loggingDataSource          databaseloggingdatasourceinterfaces.LoggingDataSource
	structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility
	createTagUsecaseComponent  tagdomainrepositoryinterfaces.CreateTagUsecaseComponent
	generatedObjectID          *primitive.ObjectID
}

func NewCreateTagTransactionComponent(
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
) (tagdomainrepositoryinterfaces.CreateTagTransactionComponent, error) {
	return &createTagTransactionComponent{
		tagDataSource:              tagDataSource,
		loggingDataSource:          loggingDataSource,
		structFieldIteratorUtility: structFieldIteratorUtility,
	}, nil
}

func (createTagTrx *createTagTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createTagTrx.tagDataSource.GetMongoDataSource().GenerateObjectID()
	createTagTrx.generatedObjectID = &generatedObjectID
	return *createTagTrx.generatedObjectID
}

func (createTagTrx *createTagTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createTagTrx.generatedObjectID == nil {
		generatedObjectID := createTagTrx.tagDataSource.GetMongoDataSource().GenerateObjectID()
		createTagTrx.generatedObjectID = &generatedObjectID
	}
	return *createTagTrx.generatedObjectID
}

func (createTagTrx *createTagTransactionComponent) SetValidation(
	usecaseComponent tagdomainrepositoryinterfaces.CreateTagUsecaseComponent,
) (bool, error) {
	createTagTrx.createTagUsecaseComponent = usecaseComponent
	return true, nil
}

func (createTagTrx *createTagTransactionComponent) PreTransaction(
	input *model.InternalCreateTag,
) (*model.InternalCreateTag, error) {
	if createTagTrx.createTagUsecaseComponent == nil {
		return input, nil
	}
	return createTagTrx.createTagUsecaseComponent.Validation(input)
}

func (createTagTrx *createTagTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateTag,
) (*model.Tag, error) {
	fieldChanges := []*model.FieldChangeDataInput{}
	createTagTrx.structFieldIteratorUtility.SetIteratingFunc(
		func(tag interface{}, field interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v",
				*tagString,
				tag,
			)

			fieldChanges = append(fieldChanges, &model.FieldChangeDataInput{
				Name:     fmt.Sprint(*tagString),
				Type:     reflect.TypeOf(field).Kind().String(),
				NewValue: fmt.Sprint(field),
			})
			*tagString = ""
		},
	)
	createTagTrx.structFieldIteratorUtility.SetPreDeepIterateFunc(
		func(tag interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v.",
				*tagString,
				tag,
			)
		},
	)
	var tagString interface{} = ""
	createTagTrx.structFieldIteratorUtility.IterateStruct(
		*input,
		&tagString,
	)

	generatedObjectID := createTagTrx.GetCurrentObjectID()
	loggingOutput, err := createTagTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Tag",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			FieldChanges: fieldChanges,
			CreatedByAccount: &model.ObjectIDOnly{
				ID: input.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *input.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createTag",
			err,
		)
	}

	input.ID = generatedObjectID
	input.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *input.ProposalStatus == model.EntityProposalStatusApproved {
		input.RecentApprovingAccount = &model.ObjectIDOnly{ID: input.SubmittingAccount.ID}
	}

	tagToCreate := &model.DatabaseCreateTag{}

	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, tagToCreate)
	json.Unmarshal(jsonTemp, &tagToCreate.ProposedChanges)

	newTag, err := createTagTrx.tagDataSource.GetMongoDataSource().Create(
		tagToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createTag",
			err,
		)
	}
	createTagTrx.generatedObjectID = nil

	return newTag, nil
}
