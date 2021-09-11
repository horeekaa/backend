package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	moupresentationusecaseinterfaces "github.com/horeekaa/backend/features/mous/presentation/usecases"
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *mouItemResolver) Product(ctx context.Context, obj *model.MouItem) (*model.Product, error) {
	var getProductUsecase productpresentationusecaseinterfaces.GetProductUsecase
	container.Make(&getProductUsecase)
	return getProductUsecase.Execute(
		&model.ProductFilterFields{
			ID: &obj.Product.ID,
		},
	)
}

func (r *mouItemResolver) Mou(ctx context.Context, obj *model.MouItem) (*model.Mou, error) {
	var getMouUsecase moupresentationusecaseinterfaces.GetMouUsecase
	container.Make(&getMouUsecase)
	return getMouUsecase.Execute(
		&model.MouFilterFields{
			ID: &obj.Mou.ID,
		},
	)
}

// MouItem returns generated.MouItemResolver implementation.
func (r *Resolver) MouItem() generated.MouItemResolver { return &mouItemResolver{r} }

type mouItemResolver struct{ *Resolver }
