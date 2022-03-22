package mouitemdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	mouitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createMouItemTransactionComponent struct {
	mouItemDataSource   databasemouitemdatasourceinterfaces.MouItemDataSource
	loggingDataSource   databaseloggingdatasourceinterfaces.LoggingDataSource
	agreedProductLoader mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader
	generatedObjectID   *primitive.ObjectID
	pathIdentity        string
}

func (createMouItemTrx *createMouItemTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createMouItemTrx.mouItemDataSource.GetMongoDataSource().GenerateObjectID()
	createMouItemTrx.generatedObjectID = &generatedObjectID
	return *createMouItemTrx.generatedObjectID
}

func (createMouItemTrx *createMouItemTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createMouItemTrx.generatedObjectID == nil {
		generatedObjectID := createMouItemTrx.mouItemDataSource.GetMongoDataSource().GenerateObjectID()
		createMouItemTrx.generatedObjectID = &generatedObjectID
	}
	return *createMouItemTrx.generatedObjectID
}

func NewCreateMouItemTransactionComponent(
	mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	agreedProductLoader mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader,
) (mouitemdomainrepositoryinterfaces.CreateMouItemTransactionComponent, error) {
	return &createMouItemTransactionComponent{
		mouItemDataSource:   mouItemDataSource,
		loggingDataSource:   loggingDataSource,
		agreedProductLoader: agreedProductLoader,
		pathIdentity:        "CreateMouItemComponent",
	}, nil
}

func (createMouItemTrx *createMouItemTransactionComponent) PreTransaction(
	createmouItemInput *model.InternalCreateMouItem,
) (*model.InternalCreateMouItem, error) {
	return createmouItemInput, nil
}

func (createMouItemTrx *createMouItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateMouItem,
) (*model.MouItem, error) {
	mouItemToCreate := &model.DatabaseCreateMouItem{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, mouItemToCreate)

	createMouItemTrx.agreedProductLoader.TransactionBody(
		session,
		mouItemToCreate.Product,
		mouItemToCreate.AgreedProduct,
		mouItemToCreate.Organization,
	)

	newDocumentJson, _ := json.Marshal(*mouItemToCreate)
	generatedObjectID := createMouItemTrx.GetCurrentObjectID()
	loggingOutput, err := createMouItemTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "MouItem",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: mouItemToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *mouItemToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createMouItemTrx.pathIdentity,
			err,
		)
	}

	mouItemToCreate.ID = generatedObjectID
	mouItemToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *mouItemToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		mouItemToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: mouItemToCreate.SubmittingAccount.ID}
	}

	defaultIsActive := true
	if mouItemToCreate.IsActive == nil {
		mouItemToCreate.IsActive = &defaultIsActive
	}

	currentTime := time.Now()
	mouItemToCreate.CreatedAt = &currentTime
	mouItemToCreate.UpdatedAt = &currentTime

	defaultProposalStatus := model.EntityProposalStatusProposed
	if mouItemToCreate.ProposalStatus == nil {
		mouItemToCreate.ProposalStatus = &defaultProposalStatus
	}

	jsonTemp, _ = json.Marshal(mouItemToCreate)
	json.Unmarshal(jsonTemp, &mouItemToCreate.ProposedChanges)

	createdVariant, err := createMouItemTrx.mouItemDataSource.GetMongoDataSource().Create(
		mouItemToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createMouItemTrx.pathIdentity,
			err,
		)
	}
	createMouItemTrx.generatedObjectID = nil

	return createdVariant, nil
}
