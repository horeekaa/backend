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
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *paymentProposedResolver) Invoice(ctx context.Context, obj *model.PaymentProposed) (*model.Invoice, error) {
	var getInvoiceUsecase invoicepresentationusecaseinterfaces.GetInvoiceUsecase
	container.Make(&getInvoiceUsecase)

	return getInvoiceUsecase.Execute(
		&model.InvoiceFilterFields{
			ID: &obj.Invoice.ID,
		},
	)
}

func (r *paymentProposedResolver) SubmittingAccount(ctx context.Context, obj *model.PaymentProposed) (*model.Account, error) {
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

func (r *paymentProposedResolver) RecentApprovingAccount(ctx context.Context, obj *model.PaymentProposed) (*model.Account, error) {
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

func (r *paymentProposedResolver) RecentLog(ctx context.Context, obj *model.PaymentProposed) (*model.Logging, error) {
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

// PaymentProposed returns generated.PaymentProposedResolver implementation.
func (r *Resolver) PaymentProposed() generated.PaymentProposedResolver {
	return &paymentProposedResolver{r}
}

type paymentProposedResolver struct{ *Resolver }
