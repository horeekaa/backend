package mouitemdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	mouitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createMouItemTransactionComponent struct {
	mouItemDataSource   databasemouitemdatasourceinterfaces.MouItemDataSource
	agreedProductLoader mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader
	generatedObjectID   *primitive.ObjectID
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
	agreedProductLoader mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader,
) (mouitemdomainrepositoryinterfaces.CreateMouItemTransactionComponent, error) {
	return &createMouItemTransactionComponent{
		mouItemDataSource:   mouItemDataSource,
		agreedProductLoader: agreedProductLoader,
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
	mouItemToCreate.ID = createMouItemTrx.GetCurrentObjectID()

	createMouItemTrx.agreedProductLoader.TransactionBody(
		session,
		mouItemToCreate,
	)

	createdVariant, err := createMouItemTrx.mouItemDataSource.GetMongoDataSource().Create(
		mouItemToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createMouItem",
			err,
		)
	}
	createMouItemTrx.generatedObjectID = nil

	return createdVariant, nil
}
