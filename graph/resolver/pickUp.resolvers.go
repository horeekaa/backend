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

func (r *supplyOrderItemPickUpResolver) Photos(ctx context.Context, obj *model.SupplyOrderItemPickUp) ([]*model.DescriptivePhoto, error) {
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

func (r *supplyOrderItemPickUpResolver) Courier(ctx context.Context, obj *model.SupplyOrderItemPickUp) (*model.AccountForSupplyOrderItem, error) {
	if obj.Courier == nil {
		return nil, nil
	}

	if obj.Courier.ID == nil {
		return nil, nil
	}

	return obj.Courier, nil
}

// SupplyOrderItemPickUp returns generated.SupplyOrderItemPickUpResolver implementation.
func (r *Resolver) SupplyOrderItemPickUp() generated.SupplyOrderItemPickUpResolver {
	return &supplyOrderItemPickUpResolver{r}
}

type supplyOrderItemPickUpResolver struct{ *Resolver }
