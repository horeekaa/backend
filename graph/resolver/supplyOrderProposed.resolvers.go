package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	paymentpresentationusecaseinterfaces "github.com/horeekaa/backend/features/payments/presentation/usecases"
	supplyorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases"
	supplyorderitempresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *supplyOrderProposedResolver) Payment(ctx context.Context, obj *model.SupplyOrderProposed) (*model.Payment, error) {
	var getPaymentUsecase paymentpresentationusecaseinterfaces.GetPaymentUsecase
	container.Make(&getPaymentUsecase)

	var filterFields *model.PaymentFilterFields
	if obj.Payment != nil {
		filterFields = &model.PaymentFilterFields{}
		filterFields.ID = &obj.Payment.ID
	}
	return getPaymentUsecase.Execute(
		filterFields,
	)
}

func (r *supplyOrderProposedResolver) Items(ctx context.Context, obj *model.SupplyOrderProposed) ([]*model.SupplyOrderItem, error) {
	var getAllSupplyOrderItemUsecase supplyorderitempresentationusecaseinterfaces.GetAllSupplyOrderItemUsecase
	container.Make(&getAllSupplyOrderItemUsecase)

	if obj.Items != nil {
		supplyOrderItems, err := getAllSupplyOrderItemUsecase.Execute(
			supplyorderitempresentationusecasetypes.GetAllSupplyOrderItemUsecaseInput{
				Context: ctx,
				FilterFields: &model.SupplyOrderItemFilterFields{
					ID: &model.ObjectIDFilterField{
						Operation: model.ObjectIDOperationIn,
						Values: funk.Map(
							obj.Items,
							func(soItem *model.SupplyOrderItem) *primitive.ObjectID {
								return &soItem.ID
							},
						).([]*primitive.ObjectID),
					},
				},
				PaginationOps: &model.PaginationOptionInput{
					QueryLimit: func(i int) *int { return &i }(999),
				},
			},
		)
		if err != nil {
			return nil, err
		}
		return supplyOrderItems, nil
	}
	return []*model.SupplyOrderItem{}, nil
}

func (r *supplyOrderProposedResolver) SubmittingAccount(ctx context.Context, obj *model.SupplyOrderProposed) (*model.Account, error) {
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

func (r *supplyOrderProposedResolver) RecentApprovingAccount(ctx context.Context, obj *model.SupplyOrderProposed) (*model.Account, error) {
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

func (r *supplyOrderProposedResolver) RecentLog(ctx context.Context, obj *model.SupplyOrderProposed) (*model.Logging, error) {
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

// SupplyOrderProposed returns generated.SupplyOrderProposedResolver implementation.
func (r *Resolver) SupplyOrderProposed() generated.SupplyOrderProposedResolver {
	return &supplyOrderProposedResolver{r}
}

type supplyOrderProposedResolver struct{ *Resolver }
