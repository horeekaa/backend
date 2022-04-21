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
	pathIdentity          string
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
		pathIdentity:          "CreateSupplyOrderComponent",
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
	input *model.InternalCreateSupplyOrder,
) (*model.SupplyOrder, error) {
	supplyOrderToCreate := &model.DatabaseCreateSupplyOrder{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, supplyOrderToCreate)

	loc, _ := time.LoadLocation("Asia/Bangkok")
	generatedObjectID := createSupplyOrderTrx.GetCurrentObjectID()
	splittedId := strings.Split(generatedObjectID.Hex(), "")
	supplyOrderToCreate.ID = generatedObjectID
	supplyOrderToCreate.PublicID = func(s ...string) string { joinedString := strings.Join(s, "/"); return joinedString }(
		"SO",
		time.Now().In(loc).Format("20060102"),
		strings.ToUpper(
			strings.Join(
				splittedId[len(splittedId)-4:],
				"",
			),
		),
	)

	totalPrice := 0
	for _, item := range input.Items {
		totalPrice += item.SubTotal
	}
	supplyOrderToCreate.Total = totalPrice
	supplyOrderToCreate.FinalSalesAmount = totalPrice

	supplyOrderToCreate.Organization = &model.OrganizationForSupplyOrderInput{
		ID: *input.MemberAccess.Organization.ID,
	}
	_, err := createSupplyOrderTrx.supplyOrderDataLoader.TransactionBody(
		session,
		supplyOrderToCreate.Organization,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createSupplyOrderTrx.pathIdentity,
			err,
		)
	}

	newDocumentJson, _ := json.Marshal(*supplyOrderToCreate)
	loggingOutput, err := createSupplyOrderTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "SupplyOrder",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: supplyOrderToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *supplyOrderToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createSupplyOrderTrx.pathIdentity,
			err,
		)
	}

	supplyOrderToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *supplyOrderToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		supplyOrderToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: supplyOrderToCreate.SubmittingAccount.ID}
	}

	currentTime := time.Now().UTC()
	supplyOrderToCreate.CreatedAt = currentTime
	supplyOrderToCreate.UpdatedAt = currentTime

	defaultProposalStatus := model.EntityProposalStatusProposed
	if supplyOrderToCreate.ProposalStatus == nil {
		supplyOrderToCreate.ProposalStatus = &defaultProposalStatus
	}

	jsonTemp, _ = json.Marshal(supplyOrderToCreate)
	json.Unmarshal(jsonTemp, &supplyOrderToCreate.ProposedChanges)

	newsupplyOrder, err := createSupplyOrderTrx.supplyOrderDataSource.GetMongoDataSource().Create(
		supplyOrderToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createSupplyOrderTrx.pathIdentity,
			err,
		)
	}
	createSupplyOrderTrx.generatedObjectID = nil

	return newsupplyOrder, nil
}
