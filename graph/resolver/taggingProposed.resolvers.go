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
	tagpresentationusecaseinterfaces "github.com/horeekaa/backend/features/tags/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *taggingProposedResolver) Tag(ctx context.Context, obj *model.TaggingProposed) (*model.Tag, error) {
	var getTagUsecase tagpresentationusecaseinterfaces.GetTagUsecase
	container.Make(&getTagUsecase)
	return getTagUsecase.Execute(
		&model.TagFilterFields{
			ID: &obj.Tag.ID,
		},
	)
}

func (r *taggingProposedResolver) CorrelatedTag(ctx context.Context, obj *model.TaggingProposed) (*model.Tag, error) {
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

func (r *taggingProposedResolver) Product(ctx context.Context, obj *model.TaggingProposed) (*model.Product, error) {
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

func (r *taggingProposedResolver) Organization(ctx context.Context, obj *model.TaggingProposed) (*model.Organization, error) {
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

func (r *taggingProposedResolver) SubmittingAccount(ctx context.Context, obj *model.TaggingProposed) (*model.Account, error) {
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

func (r *taggingProposedResolver) RecentApprovingAccount(ctx context.Context, obj *model.TaggingProposed) (*model.Account, error) {
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

func (r *taggingProposedResolver) RecentLog(ctx context.Context, obj *model.TaggingProposed) (*model.Logging, error) {
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

// TaggingProposed returns generated.TaggingProposedResolver implementation.
func (r *Resolver) TaggingProposed() generated.TaggingProposedResolver {
	return &taggingProposedResolver{r}
}

type taggingProposedResolver struct{ *Resolver }
