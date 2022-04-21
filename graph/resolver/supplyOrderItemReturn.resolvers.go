package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	descriptivephotopresentationusecaseinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *supplyOrderItemReturnResolver) Photos(ctx context.Context, obj *model.SupplyOrderItemReturn) ([]*model.DescriptivePhoto, error) {
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

// SupplyOrderItemReturn returns generated.SupplyOrderItemReturnResolver implementation.
func (r *Resolver) SupplyOrderItemReturn() generated.SupplyOrderItemReturnResolver {
	return &supplyOrderItemReturnResolver{r}
}

type supplyOrderItemReturnResolver struct{ *Resolver }
