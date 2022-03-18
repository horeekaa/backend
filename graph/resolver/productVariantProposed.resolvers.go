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
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *productVariantProposedResolver) Photo(ctx context.Context, obj *model.ProductVariantProposed) (*model.DescriptivePhoto, error) {
	var getDescriptivePhotoUsecase descriptivephotopresentationusecaseinterfaces.GetDescriptivePhotoUsecase
	container.Make(&getDescriptivePhotoUsecase)
	return getDescriptivePhotoUsecase.Execute(
		&model.DescriptivePhotoFilterFields{
			ID: &obj.Photo.ID,
		},
	)
}

func (r *productVariantProposedResolver) Product(ctx context.Context, obj *model.ProductVariantProposed) (*model.Product, error) {
	var getProductUsecase productpresentationusecaseinterfaces.GetProductUsecase
	container.Make(&getProductUsecase)
	return getProductUsecase.Execute(
		&model.ProductFilterFields{
			ID: &obj.Product.ID,
		},
	)
}

func (r *productVariantProposedResolver) SubmittingAccount(ctx context.Context, obj *model.ProductVariantProposed) (*model.Account, error) {
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

func (r *productVariantProposedResolver) RecentApprovingAccount(ctx context.Context, obj *model.ProductVariantProposed) (*model.Account, error) {
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

func (r *productVariantProposedResolver) RecentLog(ctx context.Context, obj *model.ProductVariantProposed) (*model.Logging, error) {
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

// ProductVariantProposed returns generated.ProductVariantProposedResolver implementation.
func (r *Resolver) ProductVariantProposed() generated.ProductVariantProposedResolver {
	return &productVariantProposedResolver{r}
}

type productVariantProposedResolver struct{ *Resolver }
