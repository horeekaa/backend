package purchaseorderdomainrepositories

import (
	"encoding/json"
	"strings"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createPurchaseOrderTransactionComponent struct {
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	loggingDataSource       databaseloggingdatasourceinterfaces.LoggingDataSource
	mouDataSource           databasemoudatasourceinterfaces.MouDataSource
	purchaseOrderDataLoader purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader
	generatedObjectID       *primitive.ObjectID
}

func NewCreatePurchaseOrderTransactionComponent(
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	purchaseOrderDataLoader purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader,
) (purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderTransactionComponent, error) {
	return &createPurchaseOrderTransactionComponent{
		purchaseOrderDataSource: purchaseOrderDataSource,
		loggingDataSource:       loggingDataSource,
		mouDataSource:           mouDataSource,
		purchaseOrderDataLoader: purchaseOrderDataLoader,
	}, nil
}

func (createPurchaseOrderTrx *createPurchaseOrderTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createPurchaseOrderTrx.purchaseOrderDataSource.GetMongoDataSource().GenerateObjectID()
	createPurchaseOrderTrx.generatedObjectID = &generatedObjectID
	return *createPurchaseOrderTrx.generatedObjectID
}

func (createPurchaseOrderTrx *createPurchaseOrderTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createPurchaseOrderTrx.generatedObjectID == nil {
		generatedObjectID := createPurchaseOrderTrx.purchaseOrderDataSource.GetMongoDataSource().GenerateObjectID()
		createPurchaseOrderTrx.generatedObjectID = &generatedObjectID
	}
	return *createPurchaseOrderTrx.generatedObjectID
}

func (createPurchaseOrderTrx *createPurchaseOrderTransactionComponent) PreTransaction(
	input *model.InternalCreatePurchaseOrder,
) (*model.InternalCreatePurchaseOrder, error) {
	purchaseOrder, err := createPurchaseOrderTrx.purchaseOrderDataSource.GetMongoDataSource().FindOne(
		map[string]interface{}{
			"status": map[string]interface{}{
				"$in": [...]model.PurchaseOrderStatus{
					model.PurchaseOrderStatusWaitingForInvoice,
					model.PurchaseOrderStatusInvoiced,
					model.PurchaseOrderStatusOpen,
					model.PurchaseOrderStatusProcessed,
				},
			},
			"type":             model.PurchaseOrderTypeRetail,
			"organization._id": input.MemberAccess.Organization.ID,
		},
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createPurchaseOrder",
			err,
		)
	}
	if purchaseOrder != nil {
		return nil, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.POSalesAmountExceedCreditLimit,
			"/createPurchaseOrder",
			nil,
		)
	}
	return input, nil
}

func (createPurchaseOrderTrx *createPurchaseOrderTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreatePurchaseOrder,
) (*model.PurchaseOrder, error) {
	purchaseOrderToCreate := &model.DatabaseCreatePurchaseOrder{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, purchaseOrderToCreate)

	loc, _ := time.LoadLocation("Asia/Bangkok")
	generatedObjectID := createPurchaseOrderTrx.GetCurrentObjectID()
	splittedId := strings.Split(generatedObjectID.Hex(), "")
	purchaseOrderToCreate.ID = generatedObjectID
	purchaseOrderToCreate.PublicID = func(s ...string) string { joinedString := strings.Join(s, "/"); return joinedString }(
		"PO",
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
		totalPrice += *item.SubTotal
	}
	purchaseOrderToCreate.Total = totalPrice
	purchaseOrderToCreate.FinalSalesAmount = purchaseOrderToCreate.Total

	purchaseOrderToCreate.Organization = &model.OrganizationForPurchaseOrderInput{
		ID: *input.MemberAccess.Organization.ID,
	}
	_, err := createPurchaseOrderTrx.purchaseOrderDataLoader.TransactionBody(
		session,
		purchaseOrderToCreate.Mou,
		purchaseOrderToCreate.Organization,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createPurchaseOrder",
			err,
		)
	}
	if purchaseOrderToCreate.Mou != nil {
		if purchaseOrderToCreate.Total < *purchaseOrderToCreate.Mou.MinimumOrderValueBeforeDelivery {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.POMinimumOrderValueHasNotMet,
				"/createPurchaseOrder",
				nil,
			)
		}

		*purchaseOrderToCreate.Mou.RemainingCreditLimit -= purchaseOrderToCreate.FinalSalesAmount
		if *purchaseOrderToCreate.Mou.RemainingCreditLimit < 0 {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.POSalesAmountExceedCreditLimit,
				"/createPurchaseOrder",
				nil,
			)
		}

		_, err = createPurchaseOrderTrx.mouDataSource.GetMongoDataSource().Update(
			map[string]interface{}{
				"_id": purchaseOrderToCreate.Mou.ID,
			},
			&model.DatabaseUpdateMou{
				RemainingCreditLimit: purchaseOrderToCreate.Mou.RemainingCreditLimit,
			},
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/createPurchaseOrder",
				err,
			)
		}
	}

	newDocumentJson, _ := json.Marshal(*purchaseOrderToCreate)
	loggingOutput, err := createPurchaseOrderTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "PurchaseOrder",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: purchaseOrderToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *purchaseOrderToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createPurchaseOrder",
			err,
		)
	}

	purchaseOrderToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *purchaseOrderToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		purchaseOrderToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: purchaseOrderToCreate.SubmittingAccount.ID}
	}

	jsonTemp, _ = json.Marshal(purchaseOrderToCreate)
	json.Unmarshal(jsonTemp, &purchaseOrderToCreate.ProposedChanges)

	newPurchaseOrder, err := createPurchaseOrderTrx.purchaseOrderDataSource.GetMongoDataSource().Create(
		purchaseOrderToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createPurchaseOrder",
			err,
		)
	}
	createPurchaseOrderTrx.generatedObjectID = nil

	return newPurchaseOrder, nil
}
