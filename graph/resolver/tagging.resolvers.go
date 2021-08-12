package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
	taggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/taggings/presentation/usecases"
	taggingpresentationusecasetypes "github.com/horeekaa/backend/features/taggings/presentation/usecases/types"
	tagpresentationusecaseinterfaces "github.com/horeekaa/backend/features/tags/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *mutationResolver) BulkCreateTagging(ctx context.Context, bulkCreateTagging model.BulkCreateTagging) ([]*model.Tagging, error) {
	var bulkCreateTaggingUsecase taggingpresentationusecaseinterfaces.BulkCreateTaggingUsecase
	container.Make(&bulkCreateTaggingUsecase)
	return bulkCreateTaggingUsecase.Execute(
		taggingpresentationusecasetypes.BulkCreateTaggingUsecaseInput{
			Context:           ctx,
			BulkCreateTagging: &bulkCreateTagging,
		},
	)
}

func (r *mutationResolver) BulkUpdateTagging(ctx context.Context, bulkUpdateTagging model.BulkUpdateTagging) ([]*model.Tagging, error) {
	var bulkUpdateTaggingUsecase taggingpresentationusecaseinterfaces.BulkUpdateTaggingUsecase
	container.Make(&bulkUpdateTaggingUsecase)
	return bulkUpdateTaggingUsecase.Execute(
		taggingpresentationusecasetypes.BulkUpdateTaggingUsecaseInput{
			Context:           ctx,
			BulkUpdateTagging: &bulkUpdateTagging,
		},
	)
}

func (r *queryResolver) Taggings(ctx context.Context, filterFields model.TaggingFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.Tagging, error) {
	var getTaggingsUsecase taggingpresentationusecaseinterfaces.GetAllTaggingUsecase
	container.Make(&getTaggingsUsecase)
	return getTaggingsUsecase.Execute(
		taggingpresentationusecasetypes.GetAllTaggingUsecaseInput{
			Context:       ctx,
			FilterFields:  &filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

func (r *taggingResolver) Tag(ctx context.Context, obj *model.Tagging) (*model.Tag, error) {
	var getTagUsecase tagpresentationusecaseinterfaces.GetTagUsecase
	container.Make(&getTagUsecase)
	return getTagUsecase.Execute(
		&model.TagFilterFields{
			ID: &obj.Tag.ID,
		},
	)
}

func (r *taggingResolver) CorrelatedTag(ctx context.Context, obj *model.Tagging) (*model.Tag, error) {
	var getTagUsecase tagpresentationusecaseinterfaces.GetTagUsecase
	container.Make(&getTagUsecase)

	var filterFields *model.TagFilterFields
	if obj.CorrelatedTag != nil {
		filterFields = &model.TagFilterFields{}
		filterFields.ID = &obj.CorrelatedTag.ID
	}
	return getTagUsecase.Execute(
		filterFields,
	)
}

func (r *taggingResolver) Product(ctx context.Context, obj *model.Tagging) (*model.Product, error) {
	var getProductUsecase productpresentationusecaseinterfaces.GetProductUsecase
	container.Make(&getProductUsecase)

	var filterFields *model.ProductFilterFields
	if obj.Product != nil {
		filterFields = &model.ProductFilterFields{}
		filterFields.ID = &obj.Product.ID
	}
	return getProductUsecase.Execute(
		filterFields,
	)
}

func (r *taggingResolver) Organization(ctx context.Context, obj *model.Tagging) (*model.Organization, error) {
	var getOrganizationUsecase organizationpresentationusecaseinterfaces.GetOrganizationUsecase
	container.Make(&getOrganizationUsecase)

	var filterFields *model.OrganizationFilterFields
	if obj.Organization != nil {
		filterFields = &model.OrganizationFilterFields{}
		filterFields.ID = &obj.Organization.ID
	}
	return getOrganizationUsecase.Execute(
		filterFields,
	)
}

func (r *taggingResolver) SubmittingAccount(ctx context.Context, obj *model.Tagging) (*model.Account, error) {
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

func (r *taggingResolver) RecentApprovingAccount(ctx context.Context, obj *model.Tagging) (*model.Account, error) {
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

func (r *taggingResolver) RecentLog(ctx context.Context, obj *model.Tagging) (*model.Logging, error) {
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

// Tagging returns generated.TaggingResolver implementation.
func (r *Resolver) Tagging() generated.TaggingResolver { return &taggingResolver{r} }

type taggingResolver struct{ *Resolver }
