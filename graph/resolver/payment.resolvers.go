package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	invoicepresentationusecaseinterfaces "github.com/horeekaa/backend/features/invoices/presentation/usecases"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	paymentpresentationusecaseinterfaces "github.com/horeekaa/backend/features/payments/presentation/usecases"
	paymentpresentationusecasetypes "github.com/horeekaa/backend/features/payments/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *mutationResolver) CreatePayment(ctx context.Context, createPayment model.CreatePayment) (*model.Payment, error) {
	var createPaymentUsecase paymentpresentationusecaseinterfaces.CreatePaymentUsecase
	container.Make(&createPaymentUsecase)
	return createPaymentUsecase.Execute(
		paymentpresentationusecasetypes.CreatePaymentUsecaseInput{
			Context:       ctx,
			CreatePayment: &createPayment,
		},
	)
}

func (r *mutationResolver) UpdatePayment(ctx context.Context, updatePayment model.UpdatePayment) (*model.Payment, error) {
	var updatePaymentUsecase paymentpresentationusecaseinterfaces.UpdatePaymentUsecase
	container.Make(&updatePaymentUsecase)
	return updatePaymentUsecase.Execute(
		paymentpresentationusecasetypes.UpdatePaymentUsecaseInput{
			Context:       ctx,
			UpdatePayment: &updatePayment,
		},
	)
}

func (r *paymentResolver) Invoice(ctx context.Context, obj *model.Payment) (*model.Invoice, error) {
	var getInvoiceUsecase invoicepresentationusecaseinterfaces.GetInvoiceUsecase
	container.Make(&getInvoiceUsecase)

	return getInvoiceUsecase.Execute(
		&model.InvoiceFilterFields{
			ID: &obj.Invoice.ID,
		},
	)
}

func (r *paymentResolver) SubmittingAccount(ctx context.Context, obj *model.Payment) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)
	return getAccountUsecase.Execute(
		accountpresentationusecasetypes.GetAccountInput{
			FilterFields: &model.AccountFilterFields{
				ID: &obj.SubmittingAccount.ID,
			},
		},
	)
}

func (r *paymentResolver) RecentApprovingAccount(ctx context.Context, obj *model.Payment) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)

	var filterFields *model.AccountFilterFields
	if obj.RecentApprovingAccount != nil {
		filterFields = &model.AccountFilterFields{}
		filterFields.ID = &obj.RecentApprovingAccount.ID
	}
	return getAccountUsecase.Execute(
		accountpresentationusecasetypes.GetAccountInput{
			FilterFields: filterFields,
		},
	)
}

func (r *paymentResolver) RecentLog(ctx context.Context, obj *model.Payment) (*model.Logging, error) {
	var getLoggingUsecase loggingpresentationusecaseinterfaces.GetLoggingUsecase
	container.Make(&getLoggingUsecase)

	var filterFields *model.LoggingFilterFields
	if obj.RecentLog != nil {
		filterFields = &model.LoggingFilterFields{}
		filterFields.ID = &obj.RecentLog.ID
	}
	return getLoggingUsecase.Execute(
		filterFields,
	)
}

func (r *queryResolver) Payments(ctx context.Context, filterFields model.PaymentFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.Payment, error) {
	var getPaymentsUsecase paymentpresentationusecaseinterfaces.GetAllPaymentUsecase
	container.Make(&getPaymentsUsecase)
	return getPaymentsUsecase.Execute(
		paymentpresentationusecasetypes.GetAllPaymentUsecaseInput{
			Context:       ctx,
			FilterFields:  &filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

// Payment returns generated.PaymentResolver implementation.
func (r *Resolver) Payment() generated.PaymentResolver { return &paymentResolver{r} }

type paymentResolver struct{ *Resolver }
