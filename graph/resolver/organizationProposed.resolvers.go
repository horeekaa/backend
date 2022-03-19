package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	addresspresentationusecaseinterfaces "github.com/horeekaa/backend/features/addresses/presentation/usecases"
	descriptivephotopresentationusecaseinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/presentation/usecases"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	taggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/taggings/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *organizationProposedResolver) ProfilePhotos(ctx context.Context, obj *model.OrganizationProposed) ([]*model.DescriptivePhoto, error) {
	var getDescriptivePhotoUsecase descriptivephotopresentationusecaseinterfaces.GetDescriptivePhotoUsecase
	container.Make(&getDescriptivePhotoUsecase)

	descriptivePhotos := []*model.DescriptivePhoto{}
	if obj.ProfilePhotos != nil {
		for _, photo := range obj.ProfilePhotos {
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

func (r *organizationProposedResolver) Taggings(ctx context.Context, obj *model.OrganizationProposed) ([]*model.Tagging, error) {
	var getTaggingUsecase taggingpresentationusecaseinterfaces.GetTaggingUsecase
	container.Make(&getTaggingUsecase)

	taggings := []*model.Tagging{}
	if obj.Taggings != nil {
		for _, tagg := range obj.Taggings {
			tagging, err := getTaggingUsecase.Execute(
				&model.TaggingFilterFields{
					ID: &tagg.ID,
				},
			)
			if err != nil {
				return nil, err
			}

			taggings = append(taggings, tagging)
		}
	}
	return taggings, nil
}

func (r *organizationProposedResolver) Addresses(ctx context.Context, obj *model.OrganizationProposed) ([]*model.Address, error) {
	var getAddressUsecase addresspresentationusecaseinterfaces.GetAddressUsecase
	container.Make(&getAddressUsecase)

	addresses := []*model.Address{}
	if obj.Addresses != nil {
		for _, addr := range obj.Addresses {
			address, err := getAddressUsecase.Execute(
				&model.AddressFilterFields{
					ID: &addr.ID,
				},
			)
			if err != nil {
				return nil, err
			}

			addresses = append(addresses, address)
		}
	}
	return addresses, nil
}

func (r *organizationProposedResolver) SubmittingAccount(ctx context.Context, obj *model.OrganizationProposed) (*model.Account, error) {
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

func (r *organizationProposedResolver) RecentApprovingAccount(ctx context.Context, obj *model.OrganizationProposed) (*model.Account, error) {
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

func (r *organizationProposedResolver) RecentLog(ctx context.Context, obj *model.OrganizationProposed) (*model.Logging, error) {
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

// OrganizationProposed returns generated.OrganizationProposedResolver implementation.
func (r *Resolver) OrganizationProposed() generated.OrganizationProposedResolver {
	return &organizationProposedResolver{r}
}

type organizationProposedResolver struct{ *Resolver }
