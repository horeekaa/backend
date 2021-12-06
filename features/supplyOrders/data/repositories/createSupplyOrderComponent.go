package supplyorderdomainrepositories

import (
	"encoding/json"
	"strings"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createSupplyOrderTransactionComponent struct {
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource
	loggingDataSource     databaseloggingdatasourceinterfaces.LoggingDataSource
	supplyOrderDataLoader supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader
	generatedObjectID     *primitive.ObjectID
}

func NewCreateSupplyOrderTransactionComponent(
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	supplyOrderDataLoader supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader,
) (supplyorderdomainrepositoryinterfaces.CreateSupplyOrderTransactionComponent, error) {
	return &createSupplyOrderTransactionComponent{
		supplyOrderDataSource: supplyOrderDataSource,
		loggingDataSource:     loggingDataSource,
		supplyOrderDataLoader: supplyOrderDataLoader,
	}, nil
}

func (createSupplyOrderTrx *createSupplyOrderTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createSupplyOrderTrx.supplyOrderDataSource.GetMongoDataSource().GenerateObjectID()
	createSupplyOrderTrx.generatedObjectID = &generatedObjectID
	return *createSupplyOrderTrx.generatedObjectID
}

func (createSupplyOrderTrx *createSupplyOrderTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createSupplyOrderTrx.generatedObjectID == nil {
		generatedObjectID := createSupplyOrderTrx.supplyOrderDataSource.GetMongoDataSource().GenerateObjectID()
		createSupplyOrderTrx.generatedObjectID = &generatedObjectID
	}
	return *createSupplyOrderTrx.generatedObjectID
}

func (createSupplyOrderTrx *createSupplyOrderTransactionComponent) PreTransaction(
	input *model.InternalCreateSupplyOrder,
) (*model.InternalCreateSupplyOrder, error) {
	return input, nil
}

func (createSupplyOrderTrx *createSupplyOrderTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	createSupplyOrder *model.InternalCreateSupplyOrder,
) (*model.SupplyOrder, error) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	generatedObjectID := createSupplyOrderTrx.GetCurrentObjectID()

	totalPrice := 0
	for _, item := range createSupplyOrder.Items {
		totalPrice += item.SubTotal
	}
	createSupplyOrder.Total = totalPrice

	createSupplyOrder.Organization = &model.OrganizationForSupplyOrderInput{
		ID: createSupplyOrder.MemberAccess.Organization.ID,
	}
	_, err := createSupplyOrderTrx.supplyOrderDataLoader.TransactionBody(
		session,
		createSupplyOrder.Organization,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createSupplyOrder",
			err,
		)
	}

	newDocumentJson, _ := json.Marshal(*createSupplyOrder)
	loggingOutput, err := createSupplyOrderTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "SupplyOrder",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: createSupplyOrder.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *createSupplyOrder.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createSupplyOrder",
			err,
		)
	}

	createSupplyOrder.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *createSupplyOrder.ProposalStatus == model.EntityProposalStatusApproved {
		createSupplyOrder.RecentApprovingAccount = &model.ObjectIDOnly{ID: createSupplyOrder.SubmittingAccount.ID}
	}

	splittedId := strings.Split(generatedObjectID.Hex(), "")
	supplyOrderToCreate := &model.DatabaseCreateSupplyOrder{
		ID: generatedObjectID,
		PublicID: func(s ...string) string { joinedString := strings.Join(s, "/"); return joinedString }(
			"SO",
			time.Now().In(loc).Format("20060102"),
			strings.ToUpper(
				strings.Join(
					splittedId[len(splittedId)-4:],
					"",
				),
			),
		),
	}
	jsonTemp, _ := json.Marshal(createSupplyOrder)
	json.Unmarshal(jsonTemp, &supplyOrderToCreate)
	json.Unmarshal(jsonTemp, &supplyOrderToCreate.ProposedChanges)

	newsupplyOrder, err := createSupplyOrderTrx.supplyOrderDataSource.GetMongoDataSource().Create(
		supplyOrderToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createSupplyOrder",
			err,
		)
	}
	createSupplyOrderTrx.generatedObjectID = nil

	return newsupplyOrder, nil
}
