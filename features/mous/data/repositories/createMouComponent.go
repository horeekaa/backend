package moudomainrepositories

import (
	"encoding/json"
	"fmt"
	"reflect"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moudomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createMouTransactionComponent struct {
	mouDataSource              databasemoudatasourceinterfaces.MouDataSource
	loggingDataSource          databaseloggingdatasourceinterfaces.LoggingDataSource
	structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility
	partyLoader                moudomainrepositoryutilityinterfaces.PartyLoader
	generatedObjectID          *primitive.ObjectID
}

func NewCreateMouTransactionComponent(
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
	partyLoader moudomainrepositoryutilityinterfaces.PartyLoader,
) (moudomainrepositoryinterfaces.CreateMouTransactionComponent, error) {
	return &createMouTransactionComponent{
		mouDataSource:              mouDataSource,
		loggingDataSource:          loggingDataSource,
		structFieldIteratorUtility: structFieldIteratorUtility,
		partyLoader:                partyLoader,
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
	fieldChanges := []*model.FieldChangeDataInput{}
	createMouTrx.structFieldIteratorUtility.SetIteratingFunc(
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
	createMouTrx.structFieldIteratorUtility.SetPreDeepIterateFunc(
		func(tag interface{}, tagString *interface{}) {
			*tagString = fmt.Sprintf(
				"%v%v.",
				*tagString,
				tag,
			)
		},
	)
	var tagString interface{} = ""
	createMouTrx.structFieldIteratorUtility.IterateStruct(
		*input,
		&tagString,
	)

	generatedObjectID := createMouTrx.GetCurrentObjectID()
	loggingOutput, err := createMouTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "mou",
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
			"/createMou",
			err,
		)
	}

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
