package invoicedomainrepositories

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createInvoiceTransactionComponent struct {
	invoiceDataSource       databaseinvoicedatasourceinterfaces.InvoiceDataSource
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	generatedObjectID       *primitive.ObjectID
}

func NewCreateInvoiceTransactionComponent(
	invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
) (invoicedomainrepositoryinterfaces.CreateInvoiceTransactionComponent, error) {
	return &createInvoiceTransactionComponent{
		invoiceDataSource:       invoiceDataSource,
		purchaseOrderDataSource: purchaseOrderDataSource,
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
	input *model.InternalCreateInvoice,
) (*model.InternalCreateInvoice, error) {
	return input, nil
}

func (createInvoiceTrx *createInvoiceTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateInvoice,
) ([]*model.Invoice, error) {
	invoiceToCreate := &model.DatabaseCreateInvoice{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, invoiceToCreate)

	currentTime := time.Now()
	if invoiceToCreate.PaymentDueDate == nil {
		futureDateOnly := time.Date(
			currentTime.Year(),
			currentTime.Month(),
			currentTime.Day()+7,
			0, 0, 0, 0,
			currentTime.Location(),
		)
		invoiceToCreate.PaymentDueDate = &futureDateOnly
	} else {
		dateOnly := time.Date(
			invoiceToCreate.PaymentDueDate.Year(),
			invoiceToCreate.PaymentDueDate.Month(),
			invoiceToCreate.PaymentDueDate.Day(),
			0, 0, 0, 0,
			invoiceToCreate.PaymentDueDate.Location(),
		)
		invoiceToCreate.PaymentDueDate = &dateOnly
	}

	query := map[string]interface{}{
		"status": model.PurchaseOrderStatusWaitingForInvoice,
		"paymentDueDate": map[string]interface{}{
			"$lte": invoiceToCreate.PaymentDueDate,
		},
	}
	if invoiceToCreate.StartInvoiceDate != nil && invoiceToCreate.EndInvoiceDate != nil {
		invoiceToCreate.StartInvoiceDate = func(t time.Time) *time.Time { return &t }(
			time.Date(
				invoiceToCreate.StartInvoiceDate.Year(),
				invoiceToCreate.StartInvoiceDate.Month(),
				invoiceToCreate.StartInvoiceDate.Day(),
				0, 0, 0, 0,
				invoiceToCreate.StartInvoiceDate.Location(),
			),
		)
		invoiceToCreate.EndInvoiceDate = func(t time.Time) *time.Time { return &t }(
			time.Date(
				invoiceToCreate.EndInvoiceDate.Year(),
				invoiceToCreate.EndInvoiceDate.Month(),
				invoiceToCreate.EndInvoiceDate.Day(),
				0, 0, 0, 0,
				invoiceToCreate.EndInvoiceDate.Location(),
			),
		)
		delete(query, "paymentDueDate")
		query["$and"] = []map[string]interface{}{
			{
				"paymentDueDate": map[string]interface{}{
					"$gte": invoiceToCreate.StartInvoiceDate,
				},
			},
			{
				"paymentDueDate": map[string]interface{}{
					"$lte": invoiceToCreate.EndInvoiceDate,
				},
			},
		}
		invoiceToCreate.PaymentDueDate = invoiceToCreate.EndInvoiceDate
	}

	purchaseOrders, err := createInvoiceTrx.purchaseOrderDataSource.GetMongoDataSource().Find(
		query,
		&mongodbcoretypes.PaginationOptions{},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createInvoice",
			err,
		)
	}

	groupedPurchaseOrderByOrganization := map[string]map[string][]*model.PurchaseOrder{}
	mouForInvoices := map[string]*model.MouForPurchaseOrder{}
	organizationForInvoices := map[string]*model.OrganizationForPurchaseOrder{}
	for _, po := range purchaseOrders {
		orgStringID := po.Organization.ID.Hex()
		organizationForInvoices[orgStringID] = po.Organization

		mouStringID := "NONE"
		if po.Mou != nil {
			mouStringID = po.Mou.ID.Hex()
			mouForInvoices[mouStringID] = po.Mou
		}

		if groupedPurchaseOrderByOrganization[orgStringID][mouStringID] == nil {
			groupedPurchaseOrderByOrganization[orgStringID][mouStringID] = []*model.PurchaseOrder{}
		}
		groupedPurchaseOrderByOrganization[orgStringID][mouStringID] = append(
			groupedPurchaseOrderByOrganization[orgStringID][mouStringID],
			po,
		)
	}

	invoicesCreated := []*model.Invoice{}
	for _, orgKey := range reflect.ValueOf(groupedPurchaseOrderByOrganization).MapKeys() {
		for _, mouKey := range reflect.ValueOf(groupedPurchaseOrderByOrganization[orgKey.String()]).MapKeys() {
			loc, _ := time.LoadLocation("Asia/Bangkok")
			generatedObjectID := createInvoiceTrx.GetCurrentObjectID()
			splittedId := strings.Split(generatedObjectID.Hex(), "")
			invoiceToCreate.ID = generatedObjectID
			invoiceToCreate.PublicID = func(s ...string) string { joinedString := strings.Join(s, "/"); return joinedString }(
				"INV",
				currentTime.In(loc).Format("20060102"),
				strings.ToUpper(
					strings.Join(
						splittedId[len(splittedId)-4:],
						"",
					),
				),
			)
			jsonOrgForInv, _ := json.Marshal(organizationForInvoices[orgKey.String()])
			json.Unmarshal(jsonOrgForInv, &invoiceToCreate.Organization)

			invoiceToCreate.Mou = nil
			if mouForInvoices[mouKey.String()] != nil {
				jsonMouForInv, _ := json.Marshal(mouForInvoices[mouKey.String()])
				json.Unmarshal(jsonMouForInv, &invoiceToCreate.Mou)
			}

			purchaseOrders := groupedPurchaseOrderByOrganization[orgKey.String()][mouKey.String()]

			jsonTemp, _ = json.Marshal(map[string]interface{}{
				"PurchaseOrders": purchaseOrders,
			})
			json.Unmarshal(jsonTemp, invoiceToCreate)

			totalPrice := 0
			for _, item := range purchaseOrders {
				totalPrice += item.FinalSalesAmount
			}
			invoiceToCreate.TotalValue = totalPrice

			totalDiscounted := invoiceToCreate.TotalDiscounted
			if invoiceToCreate.DiscountInPercent > 0 {
				totalDiscounted = (invoiceToCreate.DiscountInPercent / 100.0) * totalPrice
			}
			invoiceToCreate.TotalDiscounted = totalDiscounted
			invoiceToCreate.TotalPayable = totalPrice - totalDiscounted
			invoiceToCreate.PaymentDueDate = purchaseOrders[0].PaymentDueDate

			newInvoice, err := createInvoiceTrx.invoiceDataSource.GetMongoDataSource().Create(
				invoiceToCreate,
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/createInvoice",
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
							purchaseOrders,
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
					"/createInvoice",
					err,
				)
			}
			createInvoiceTrx.generatedObjectID = nil

			invoicesCreated = append(invoicesCreated, newInvoice)
		}
	}

	return invoicesCreated, nil
}
