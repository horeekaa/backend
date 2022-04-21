package tagdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createTagTransactionComponent struct {
	tagDataSource             databasetagdatasourceinterfaces.TagDataSource
	loggingDataSource         databaseloggingdatasourceinterfaces.LoggingDataSource
	createTagUsecaseComponent tagdomainrepositoryinterfaces.CreateTagUsecaseComponent
	generatedObjectID         *primitive.ObjectID
	pathIdentity              string
}

func NewCreateTagTransactionComponent(
	tagDataSource databasetagdatasourceinterfaces.TagDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
) (tagdomainrepositoryinterfaces.CreateTagTransactionComponent, error) {
	return &createTagTransactionComponent{
		tagDataSource:     tagDataSource,
		loggingDataSource: loggingDataSource,
		pathIdentity:      "CreateTagComponent",
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
	tagToCreate := &model.DatabaseCreateTag{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, tagToCreate)

	newDocumentJson, _ := json.Marshal(*tagToCreate)
	generatedObjectID := createTagTrx.GetCurrentObjectID()
	loggingOutput, err := createTagTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Tag",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: tagToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *tagToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createTagTrx.pathIdentity,
			err,
		)
	}

	tagToCreate.ID = generatedObjectID
	tagToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *tagToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		tagToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: tagToCreate.SubmittingAccount.ID}
	}

	currentTime := time.Now().UTC()
	tagToCreate.CreatedAt = &currentTime
	tagToCreate.UpdatedAt = &currentTime

	defaultProposalStatus := model.EntityProposalStatusProposed
	if tagToCreate.ProposalStatus == nil {
		tagToCreate.ProposalStatus = &defaultProposalStatus
	}

	defaultIsActive := true
	if tagToCreate.IsActive == nil {
		tagToCreate.IsActive = &defaultIsActive
	}

	if tagToCreate.Photos == nil {
		tagToCreate.Photos = []*model.ObjectIDOnly{}
	}

	jsonTemp, _ = json.Marshal(tagToCreate)
	json.Unmarshal(jsonTemp, &tagToCreate.ProposedChanges)

	newTag, err := createTagTrx.tagDataSource.GetMongoDataSource().Create(
		tagToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createTagTrx.pathIdentity,
			err,
		)
	}
	createTagTrx.generatedObjectID = nil

	return newTag, nil
}
