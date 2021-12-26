package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	descriptivephotopresentationusecaseinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/presentation/usecases"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	supplyorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases"
	supplyorderitempresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *mutationResolver) UpdateSupplyOrderItemPickUp(ctx context.Context, updateSupplyOrderItemPickUp *model.UpdateSupplyOrderItemPickUpOnly) (*model.SupplyOrderItem, error) {
	var updateSupplyOrderItemPickUpUsecase supplyorderitempresentationusecaseinterfaces.UpdateSupplyOrderItemPickUpUsecase
	container.Make(&updateSupplyOrderItemPickUpUsecase)
	return updateSupplyOrderItemPickUpUsecase.Execute(
		supplyorderitempresentationusecasetypes.UpdateSupplyOrderItemPickUpUsecaseInput{
			Context:                     ctx,
			UpdateSupplyOrderItemPickUp: updateSupplyOrderItemPickUp,
		},
	)
}

func (r *supplyOrderItemResolver) Photos(ctx context.Context, obj *model.SupplyOrderItem) ([]*model.DescriptivePhoto, error) {
	var getDescriptivePhotoUsecase descriptivephotopresentationusecaseinterfaces.GetDescriptivePhotoUsecase
	container.Make(&getDescriptivePhotoUsecase)

	descriptivePhotos := []*model.DescriptivePhoto{}
	if obj.Photos != nil {
		for _, photo := range obj.Photos {
			descriptivePhoto, err := getDescriptivePhotoUsecase.Execute(
				&model.DescriptivePhotoFilterFields{
					ID: &photo.ID,
				},
			)
			if err != nil {
				return nil, err
			}

			descriptivePhotos = append(descriptivePhotos, descriptivePhoto)
		}
	}
	return descriptivePhotos, nil
}

func (r *supplyOrderItemResolver) SubmittingAccount(ctx context.Context, obj *model.SupplyOrderItem) (*model.Account, error) {
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

func (r *supplyOrderItemResolver) RecentApprovingAccount(ctx context.Context, obj *model.SupplyOrderItem) (*model.Account, error) {
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

func (r *supplyOrderItemResolver) RecentLog(ctx context.Context, obj *model.SupplyOrderItem) (*model.Logging, error) {
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

// SupplyOrderItem returns generated.SupplyOrderItemResolver implementation.
func (r *Resolver) SupplyOrderItem() generated.SupplyOrderItemResolver {
	return &supplyOrderItemResolver{r}
}

type supplyOrderItemResolver struct{ *Resolver }
