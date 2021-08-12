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
	tagpresentationusecaseinterfaces "github.com/horeekaa/backend/features/tags/presentation/usecases"
	tagpresentationusecasetypes "github.com/horeekaa/backend/features/tags/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *mutationResolver) CreateTag(ctx context.Context, createTag model.CreateTag) (*model.Tag, error) {
	var createTagUsecase tagpresentationusecaseinterfaces.CreateTagUsecase
	container.Make(&createTagUsecase)
	return createTagUsecase.Execute(
		tagpresentationusecasetypes.CreateTagUsecaseInput{
			Context:   ctx,
			CreateTag: &createTag,
		},
	)
}

func (r *mutationResolver) UpdateTag(ctx context.Context, updateTag model.UpdateTag) (*model.Tag, error) {
	var updateTagUsecase tagpresentationusecaseinterfaces.UpdateTagUsecase
	container.Make(&updateTagUsecase)
	return updateTagUsecase.Execute(
		tagpresentationusecasetypes.UpdateTagUsecaseInput{
			Context:   ctx,
			UpdateTag: &updateTag,
		},
	)
}

func (r *queryResolver) Tags(ctx context.Context, filterFields model.TagFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.Tag, error) {
	var getAllTagUsecase tagpresentationusecaseinterfaces.GetAllTagUsecase
	container.Make(&getAllTagUsecase)
	return getAllTagUsecase.Execute(
		tagpresentationusecasetypes.GetAllTagUsecaseInput{
			Context:       ctx,
			FilterFields:  &filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

func (r *tagResolver) Photos(ctx context.Context, obj *model.Tag) ([]*model.DescriptivePhoto, error) {
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

func (r *tagResolver) SubmittingAccount(ctx context.Context, obj *model.Tag) (*model.Account, error) {
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

func (r *tagResolver) RecentApprovingAccount(ctx context.Context, obj *model.Tag) (*model.Account, error) {
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

func (r *tagResolver) RecentLog(ctx context.Context, obj *model.Tag) (*model.Logging, error) {
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

// Tag returns generated.TagResolver implementation.
func (r *Resolver) Tag() generated.TagResolver { return &tagResolver{r} }

type tagResolver struct{ *Resolver }
