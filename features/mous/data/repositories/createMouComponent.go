package moudomainrepositories

import (
	"encoding/json"
	"strings"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moudomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createMouTransactionComponent struct {
	mouDataSource     databasemoudatasourceinterfaces.MouDataSource
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource
	partyLoader       moudomainrepositoryutilityinterfaces.PartyLoader
	generatedObjectID *primitive.ObjectID
}

func NewCreateMouTransactionComponent(
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	partyLoader moudomainrepositoryutilityinterfaces.PartyLoader,
) (moudomainrepositoryinterfaces.CreateMouTransactionComponent, error) {
	return &createMouTransactionComponent{
		mouDataSource:     mouDataSource,
		loggingDataSource: loggingDataSource,
		partyLoader:       partyLoader,
	}, nil
}

func (createMouTrx *createMouTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createMouTrx.mouDataSource.GetMongoDataSource().GenerateObjectID()
	createMouTrx.generatedObjectID = &generatedObjectID
	return *createMouTrx.generatedObjectID
}

func (createMouTrx *createMouTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createMouTrx.generatedObjectID == nil {
		generatedObjectID := createMouTrx.mouDataSource.GetMongoDataSource().GenerateObjectID()
		createMouTrx.generatedObjectID = &generatedObjectID
	}
	return *createMouTrx.generatedObjectID
}

func (createMouTrx *createMouTransactionComponent) PreTransaction(
	input *model.InternalCreateMou,
) (*model.InternalCreateMou, error) {
	return input, nil
}

func (createMouTrx *createMouTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateMou,
) (*model.Mou, error) {
	newDocumentJson, _ := json.Marshal(*input)
	generatedObjectID := createMouTrx.GetCurrentObjectID()
	loggingOutput, err := createMouTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Mou",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
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
			"/createMou",
			err,
		)
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	splittedId := strings.Split(generatedObjectID.Hex(), "")
	input.PublicID = func(s ...string) *string { joinedString := strings.Join(s, "/"); return &joinedString }(
		"MOU",
		time.Now().In(loc).Format("20060102"),
		strings.ToUpper(
			strings.Join(
				splittedId[len(splittedId)-4:],
				"",
			),
		),
	)
	input.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *input.ProposalStatus == model.EntityProposalStatusApproved {
		input.RecentApprovingAccount = &model.ObjectIDOnly{ID: input.SubmittingAccount.ID}
	}

	mouToCreate := &model.DatabaseCreateMou{
		ID: generatedObjectID,
	}

	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, mouToCreate)
	_, err = createMouTrx.partyLoader.TransactionBody(
		session,
		input.FirstParty,
		mouToCreate.FirstParty,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMou",
			err,
		)
	}
	_, err = createMouTrx.partyLoader.TransactionBody(
		session,
		input.SecondParty,
		mouToCreate.SecondParty,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMou",
			err,
		)
	}

	jsonTemp, _ = json.Marshal(mouToCreate)
	json.Unmarshal(jsonTemp, &mouToCreate.ProposedChanges)

	newMou, err := createMouTrx.mouDataSource.GetMongoDataSource().Create(
		mouToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMou",
			err,
		)
	}
	createMouTrx.generatedObjectID = nil

	return newMou, nil
}
