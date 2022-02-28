package invoicedomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type updateInvoiceTransactionComponent struct {
	invoiceDataSource       databaseinvoicedatasourceinterfaces.InvoiceDataSource
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	paymentDataSource       databasepaymentdatasourceinterfaces.PaymentDataSource
}

func NewUpdateInvoiceTransactionComponent(
	invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
) (invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent, error) {
	return &updateInvoiceTransactionComponent{
		invoiceDataSource:       invoiceDataSource,
		purchaseOrderDataSource: purchaseOrderDataSource,
		paymentDataSource:       paymentDataSource,
	}, nil
}

func (updateInvoiceTrx *updateInvoiceTransactionComponent) PreTransaction(
	input *model.InternalUpdateInvoice,
) (*model.InternalUpdateInvoice, error) {
	return input, nil
}

func (updateInvoiceTrx *updateInvoiceTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateInvoiceInput *model.InternalUpdateInvoice,
) (*model.Invoice, error) {
	invoiceToUpdate := &model.DatabaseUpdateInvoice{}
	jsonTemp, _ := json.Marshal(updateInvoiceInput)
	json.Unmarshal(jsonTemp, invoiceToUpdate)

	existingInvoice, err := updateInvoiceTrx.invoiceDataSource.GetMongoDataSource().FindByID(
		invoiceToUpdate.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateInvoice",
			err,
		)
	}

	purchaseOrders := []*model.PurchaseOrder{}
	if invoiceToUpdate.PaymentDueDate != nil {
		dateOnly := time.Date(
			invoiceToUpdate.PaymentDueDate.Year(),
			invoiceToUpdate.PaymentDueDate.Month(),
			invoiceToUpdate.PaymentDueDate.Day()+7,
			0, 0, 0, 0,
			invoiceToUpdate.PaymentDueDate.Location(),
		)
		invoiceToUpdate.PaymentDueDate = &dateOnly
		query := map[string]interface{}{
			"status":         model.PurchaseOrderStatusWaitingForInvoice,
			"paymentDueDate": dateOnly,
		}

		purchaseOrders, err = updateInvoiceTrx.purchaseOrderDataSource.GetMongoDataSource().Find(
			query,
			&mongodbcoretypes.PaginationOptions{},
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateInvoice",
				err,
			)
		}
	}

	if invoiceToUpdate.StartInvoiceDate != nil && invoiceToUpdate.EndInvoiceDate != nil {
		invoiceToUpdate.StartInvoiceDate = func(t time.Time) *time.Time { return &t }(
			time.Date(
				invoiceToUpdate.StartInvoiceDate.Year(),
				invoiceToUpdate.StartInvoiceDate.Month(),
				invoiceToUpdate.StartInvoiceDate.Day(),
				0, 0, 0, 0,
				invoiceToUpdate.StartInvoiceDate.Location(),
			),
		)
		invoiceToUpdate.EndInvoiceDate = func(t time.Time) *time.Time { return &t }(
			time.Date(
				invoiceToUpdate.EndInvoiceDate.Year(),
				invoiceToUpdate.EndInvoiceDate.Month(),
				invoiceToUpdate.EndInvoiceDate.Day(),
				0, 0, 0, 0,
				invoiceToUpdate.EndInvoiceDate.Location(),
			),
		)

		query := map[string]interface{}{
			"status": model.PurchaseOrderStatusWaitingForInvoice,
		}
		query["$and"] = []map[string]interface{}{
			{
				"paymentDueDate": map[string]interface{}{
					"$gte": invoiceToUpdate.StartInvoiceDate,
				},
			},
			{
				"paymentDueDate": map[string]interface{}{
					"$lte": invoiceToUpdate.EndInvoiceDate,
				},
			},
		}
		invoiceToUpdate.PaymentDueDate = invoiceToUpdate.EndInvoiceDate

		purchaseOrders, err = updateInvoiceTrx.purchaseOrderDataSource.GetMongoDataSource().Find(
			query,
			&mongodbcoretypes.PaginationOptions{},
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateInvoice",
				err,
			)
		}
	}

	if len(purchaseOrders) > 0 {
		_, err = updateInvoiceTrx.purchaseOrderDataSource.GetMongoDataSource().UpdateAll(
			map[string]interface{}{
				"_id": map[string]interface{}{
					"$in": funk.Map(
						existingInvoice.PurchaseOrders,
						func(po *model.PurchaseOrder) interface{} {
							return po.ID
						},
					),
				},
			},
			&model.DatabaseUpdatePurchaseOrder{
				Status: func(m model.PurchaseOrderStatus) *model.PurchaseOrderStatus {
					return &m
				}(model.PurchaseOrderStatusWaitingForInvoice),
			},
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateInvoice",
				err,
			)
		}

		jsonTemp, _ = json.Marshal(map[string]interface{}{
			"PurchaseOrders": purchaseOrders,
		})
		json.Unmarshal(jsonTemp, invoiceToUpdate)

		totalPrice := 0
		for _, item := range purchaseOrders {
			totalPrice += item.FinalSalesAmount
		}
		invoiceToUpdate.TotalValue = &totalPrice

		totalDiscounted := existingInvoice.TotalDiscounted
		if invoiceToUpdate.TotalDiscounted != nil {
			totalDiscounted = *invoiceToUpdate.TotalDiscounted
		}

		discountInPercent := existingInvoice.DiscountInPercent
		if invoiceToUpdate.DiscountInPercent != nil {
			discountInPercent = *invoiceToUpdate.DiscountInPercent
		}

		if discountInPercent > 0 {
			totalDiscounted = (discountInPercent / 100.0) * totalPrice
		}
		invoiceToUpdate.TotalDiscounted = &totalDiscounted
		invoiceToUpdate.TotalPayable = func(i int) *int { return &i }(totalPrice - totalDiscounted)
	}

	totalPaidAmount := 0
	if len(invoiceToUpdate.Payments) > 0 {
		invoiceToUpdate.Payments = append(
			funk.Map(
				existingInvoice.Payments,
				func(m *model.Payment) *model.ObjectIDOnly {
					return &model.ObjectIDOnly{
						ID: &m.ID,
					}
				},
			).([]*model.ObjectIDOnly),
			invoiceToUpdate.Payments...,
		)

		payments, err := updateInvoiceTrx.paymentDataSource.GetMongoDataSource().Find(
			map[string]interface{}{
				"_id": map[string]interface{}{
					"$in": funk.Map(
						invoiceToUpdate.Payments,
						func(pyt *model.ObjectIDOnly) interface{} {
							return pyt.ID
						},
					),
				},
			},
			&mongodbcoretypes.PaginationOptions{},
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updateInvoice",
				err,
			)
		}

		for _, payment := range payments {
			if payment.ProposalStatus != model.EntityProposalStatusApproved {
				continue
			}
			totalPaidAmount += payment.Amount
		}
		invoiceToUpdate.TotalPaidAmount = &totalPaidAmount
	}

	updatedInvoice, err := updateInvoiceTrx.invoiceDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": invoiceToUpdate.ID,
		},
		invoiceToUpdate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateInvoice",
			err,
		)
	}

	updatePO := &model.DatabaseUpdatePurchaseOrder{
		Status: func(m model.PurchaseOrderStatus) *model.PurchaseOrderStatus {
			return &m
		}(model.PurchaseOrderStatusInvoiced),
	}
	if totalPaidAmount >= existingInvoice.TotalPayable {
		updatePO.Status = func(m model.PurchaseOrderStatus) *model.PurchaseOrderStatus {
			return &m
		}(model.PurchaseOrderStatusPaid)
		invoiceToUpdate.Status = func(m model.InvoiceStatus) *model.InvoiceStatus {
			return &m
		}(model.InvoiceStatusPaid)
	}
	jsonInv, _ := json.Marshal(updatedInvoice)
	json.Unmarshal(jsonInv, &updatePO.Invoice)
	_, err = updateInvoiceTrx.purchaseOrderDataSource.GetMongoDataSource().UpdateAll(
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
			"/updateInvoice",
			err,
		)
	}

	return updatedInvoice, nil
}
