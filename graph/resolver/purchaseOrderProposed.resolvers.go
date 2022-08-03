package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	purchaseorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases"
	purchaseorderitempresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *purchaseOrderProposedResolver) Items(ctx context.Context, obj *model.PurchaseOrderProposed) ([]*model.PurchaseOrderItem, error) {
	var getAllPurchaseOrderItemUsecase purchaseorderitempresentationusecaseinterfaces.GetAllPurchaseOrderItemUsecase
	container.Make(&getAllPurchaseOrderItemUsecase)

	if obj.Items != nil {
		purchaseOrderItems, err := getAllPurchaseOrderItemUsecase.Execute(
			purchaseorderitempresentationusecasetypes.GetAllPurchaseOrderItemUsecaseInput{
				Context: ctx,
				FilterFields: &model.PurchaseOrderItemFilterFields{
					ID: &model.ObjectIDOnlyFilterField{
						ID: &model.ObjectIDFilterField{
							Operation: model.ObjectIDOperationIn,
							Values: funk.Map(
								obj.Items,
								func(poItem *model.PurchaseOrderItem) *primitive.ObjectID {
									return &poItem.ID
								},
							).([]*primitive.ObjectID),
						},
					},
				},
				PaginationOps: &model.PaginationOptionInput{},
			},
		)
		if err != nil {
			return nil, err
		}
		return purchaseOrderItems, nil
	}
	return []*model.PurchaseOrderItem{}, nil
}

func (r *purchaseOrderProposedResolver) SubmittingAccount(ctx context.Context, obj *model.PurchaseOrderProposed) (*model.Account, error) {
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

func (r *purchaseOrderProposedResolver) RecentApprovingAccount(ctx context.Context, obj *model.PurchaseOrderProposed) (*model.Account, error) {
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

func (r *purchaseOrderProposedResolver) RecentLog(ctx context.Context, obj *model.PurchaseOrderProposed) (*model.Logging, error) {
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

// PurchaseOrderProposed returns generated.PurchaseOrderProposedResolver implementation.
func (r *Resolver) PurchaseOrderProposed() generated.PurchaseOrderProposedResolver {
	return &purchaseOrderProposedResolver{r}
}

type purchaseOrderProposedResolver struct{ *Resolver }
