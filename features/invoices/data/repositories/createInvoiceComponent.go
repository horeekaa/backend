package invoicedomainrepositories

import (
	"encoding/json"
	"strings"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	invoicedomainrepositorytypes "github.com/horeekaa/backend/features/invoices/domain/repositories/types"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createInvoiceTransactionComponent struct {
	invoiceDataSource       databaseinvoicedatasourceinterfaces.InvoiceDataSource
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	generatedObjectID       *primitive.ObjectID
	pathIdentity            string
}

func NewCreateInvoiceTransactionComponent(
	invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
) (invoicedomainrepositoryinterfaces.CreateInvoiceTransactionComponent, error) {
	return &createInvoiceTransactionComponent{
		invoiceDataSource:       invoiceDataSource,
		purchaseOrderDataSource: purchaseOrderDataSource,
		pathIdentity:            "CreateInvoiceComponent",
	}, nil
}

func (createInvoiceTrx *createInvoiceTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createInvoiceTrx.invoiceDataSource.GetMongoDataSource().GenerateObjectID()
	createInvoiceTrx.generatedObjectID = &generatedObjectID
	return *createInvoiceTrx.generatedObjectID
}

func (createInvoiceTrx *createInvoiceTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createInvoiceTrx.generatedObjectID == nil {
		generatedObjectID := createInvoiceTrx.invoiceDataSource.GetMongoDataSource().GenerateObjectID()
		createInvoiceTrx.generatedObjectID = &generatedObjectID
	}
	return *createInvoiceTrx.generatedObjectID
}

func (createInvoiceTrx *createInvoiceTransactionComponent) PreTransaction(
	input *invoicedomainrepositorytypes.CreateInvoiceInput,
) (*invoicedomainrepositorytypes.CreateInvoiceInput, error) {
	return input, nil
}

func (createInvoiceTrx *createInvoiceTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *invoicedomainrepositorytypes.CreateInvoiceInput,
) (*model.Invoice, error) {
	invoiceToCreate := &model.DatabaseCreateInvoice{}
	jsonTemp, _ := json.Marshal(input.CreateInvoiceInput)
	json.Unmarshal(jsonTemp, invoiceToCreate)

	loc, _ := time.LoadLocation("Asia/Bangkok")
	generatedObjectID := createInvoiceTrx.GetCurrentObjectID()
	splittedId := strings.Split(generatedObjectID.Hex(), "")
	invoiceToCreate.ID = generatedObjectID
	invoiceToCreate.PublicID = func(s ...string) string { joinedString := strings.Join(s, "/"); return joinedString }(
		"INV",
		time.Now().In(loc).Format("20060102"),
		strings.ToUpper(
			strings.Join(
				splittedId[len(splittedId)-4:],
				"",
			),
		),
	)
	jsonOrgForInv, _ := json.Marshal(input.PurchaseOrdersToInvoice[0].Organization)
	json.Unmarshal(jsonOrgForInv, &invoiceToCreate.Organization)

	invoiceToCreate.Mou = nil
	if input.PurchaseOrdersToInvoice[0].Mou != nil {
		jsonMouForInv, _ := json.Marshal(input.PurchaseOrdersToInvoice[0].Mou)
		json.Unmarshal(jsonMouForInv, &invoiceToCreate.Mou)
	}

	jsonTemp, _ = json.Marshal(map[string]interface{}{
		"PurchaseOrders": input.PurchaseOrdersToInvoice,
	})
	json.Unmarshal(jsonTemp, invoiceToCreate)

	totalPrice := 0
	for _, item := range input.PurchaseOrdersToInvoice {
		totalPrice += item.FinalSalesAmount
	}
	invoiceToCreate.TotalValue = totalPrice

	totalDiscounted := invoiceToCreate.TotalDiscounted
	if invoiceToCreate.DiscountInPercent > 0 {
		totalDiscounted = (invoiceToCreate.DiscountInPercent / 100.0) * totalPrice
	}
	invoiceToCreate.TotalDiscounted = totalDiscounted
	invoiceToCreate.TotalPayable = totalPrice - totalDiscounted
	invoiceToCreate.PaymentDueDate = input.PurchaseOrdersToInvoice[0].PaymentDueDate

	newInvoice, err := createInvoiceTrx.invoiceDataSource.GetMongoDataSource().Create(
		invoiceToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createInvoiceTrx.pathIdentity,
			err,
		)
	}

	updatePO := &model.DatabaseUpdatePurchaseOrder{
		Status: func(m model.PurchaseOrderStatus) *model.PurchaseOrderStatus {
			return &m
		}(model.PurchaseOrderStatusInvoiced),
	}
	jsonInv, _ := json.Marshal(newInvoice)
	json.Unmarshal(jsonInv, &updatePO.Invoice)
	_, err = createInvoiceTrx.purchaseOrderDataSource.GetMongoDataSource().UpdateAll(
		map[string]interface{}{
			"_id": map[string]interface{}{
				"$in": funk.Map(
					input.PurchaseOrdersToInvoice,
					func(po *model.PurchaseOrder) interface{} {
						return po.ID
					},
				),
			},
		},
		updatePO,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createInvoiceTrx.pathIdentity,
			err,
		)
	}
	createInvoiceTrx.generatedObjectID = nil

	return newInvoice, nil
}
