package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/golobby/container/v2"
	descriptivephotopresentationusecaseinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/presentation/usecases"
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *productVariantResolver) Photo(ctx context.Context, obj *model.ProductVariant) (*model.DescriptivePhoto, error) {
	var getDescriptivePhotoUsecase descriptivephotopresentationusecaseinterfaces.GetDescriptivePhotoUsecase
	container.Make(&getDescriptivePhotoUsecase)
	return getDescriptivePhotoUsecase.Execute(
		&model.DescriptivePhotoFilterFields{
			ID: &obj.Photo.ID,
		},
	)
}

func (r *productVariantResolver) Product(ctx context.Context, obj *model.ProductVariant) (*model.Product, error) {
	var getProductUsecase productpresentationusecaseinterfaces.GetProductUsecase
	container.Make(&getProductUsecase)
	return getProductUsecase.Execute(
		&model.ProductFilterFields{
			ID: &obj.Product.ID,
		},
	)
}

// ProductVariant returns generated.ProductVariantResolver implementation.
func (r *Resolver) ProductVariant() generated.ProductVariantResolver {
	return &productVariantResolver{r}
}

type productVariantResolver struct{ *Resolver }
