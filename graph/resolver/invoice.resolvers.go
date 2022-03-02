package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	invoicepresentationusecaseinterfaces "github.com/horeekaa/backend/features/invoices/presentation/usecases"
	invoicepresentationusecasetypes "github.com/horeekaa/backend/features/invoices/presentation/usecases/types"
	paymentpresentationusecaseinterfaces "github.com/horeekaa/backend/features/payments/presentation/usecases"
	purchaseorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *invoiceResolver) PurchaseOrders(ctx context.Context, obj *model.Invoice) ([]*model.PurchaseOrder, error) {
	var getPurchaseOrderUsecase purchaseorderpresentationusecaseinterfaces.GetPurchaseOrderUsecase
	container.Make(&getPurchaseOrderUsecase)

	purchaseOrders := []*model.PurchaseOrder{}
	if obj.PurchaseOrders != nil {
		for _, po := range obj.PurchaseOrders {
			purchaseOrder, err := getPurchaseOrderUsecase.Execute(
				&model.PurchaseOrderFilterFields{
					ID: &po.ID,
				},
			)
			if err != nil {
				return nil, err
			}

			purchaseOrders = append(purchaseOrders, purchaseOrder)
		}
	}
	return purchaseOrders, nil
}

func (r *invoiceResolver) Payments(ctx context.Context, obj *model.Invoice) ([]*model.Payment, error) {
	var getPaymentUsecase paymentpresentationusecaseinterfaces.GetPaymentUsecase
	container.Make(&getPaymentUsecase)

	payments := []*model.Payment{}
	if obj.Payments != nil {
		for _, pay := range obj.Payments {
			payment, err := getPaymentUsecase.Execute(
				&model.PaymentFilterFields{
					ID: &pay.ID,
				},
			)
			if err != nil {
				return nil, err
			}

			payments = append(payments, payment)
		}
	}
	return payments, nil
}

func (r *mutationResolver) CreateInvoice(ctx context.Context, createInvoice model.CreateInvoice) ([]*model.Invoice, error) {
	var createInvoiceUsecase invoicepresentationusecaseinterfaces.CreateInvoiceUsecase
	container.Make(&createInvoiceUsecase)
	return createInvoiceUsecase.Execute(
		invoicepresentationusecasetypes.CreateInvoiceUsecaseInput{
			Context:       ctx,
			CreateInvoice: &createInvoice,
		},
	)
}

func (r *mutationResolver) UpdateInvoice(ctx context.Context, updateInvoice model.UpdateInvoice) (*model.Invoice, error) {
	var updateInvoiceUsecase invoicepresentationusecaseinterfaces.UpdateInvoiceUsecase
	container.Make(&updateInvoiceUsecase)
	return updateInvoiceUsecase.Execute(
		invoicepresentationusecasetypes.UpdateInvoiceUsecaseInput{
			Context:       ctx,
			UpdateInvoice: &updateInvoice,
		},
	)
}

func (r *queryResolver) Invoices(ctx context.Context, filterFields model.InvoiceFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.Invoice, error) {
	var getInvoicesUsecase invoicepresentationusecaseinterfaces.GetAllInvoiceUsecase
	container.Make(&getInvoicesUsecase)
	return getInvoicesUsecase.Execute(
		invoicepresentationusecasetypes.GetAllInvoiceUsecaseInput{
			Context:       ctx,
			FilterFields:  &filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

// Invoice returns generated.InvoiceResolver implementation.
func (r *Resolver) Invoice() generated.InvoiceResolver { return &invoiceResolver{r} }

type invoiceResolver struct{ *Resolver }
