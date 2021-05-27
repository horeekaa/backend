package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
	organizationpresentationusecasetypes "github.com/horeekaa/backend/features/organizations/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *mutationResolver) CreateOrganization(ctx context.Context, createOrganization model.CreateOrganization) (*model.Organization, error) {
	var createOrganizationUsecase organizationpresentationusecaseinterfaces.CreateOrganizationUsecase
	container.Make(&createOrganizationUsecase)
	return createOrganizationUsecase.Execute(
		organizationpresentationusecasetypes.CreateOrganizationUsecaseInput{
			Context:            ctx,
			CreateOrganization: &createOrganization,
		},
	)
}

func (r *mutationResolver) UpdateOrganization(ctx context.Context, updateOrganization model.UpdateOrganization) (*model.Organization, error) {
	var updateOrganizationUsecase organizationpresentationusecaseinterfaces.UpdateOrganizationUsecase
	container.Make(&updateOrganizationUsecase)
	return updateOrganizationUsecase.Execute(
		organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput{
			Context:            ctx,
			UpdateOrganization: &updateOrganization,
		},
	)
}

func (r *organizationResolver) SubmittingAccount(ctx context.Context, obj *model.Organization) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)
	return getAccountUsecase.Execute(
		&model.AccountFilterFields{
			ID: &obj.SubmittingAccount.ID,
		},
	)
}

func (r *organizationResolver) ApprovingAccount(ctx context.Context, obj *model.Organization) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)

	var filterFields *model.AccountFilterFields
	if obj.ApprovingAccount != nil {
		filterFields = &model.AccountFilterFields{}
		filterFields.ID = &obj.ApprovingAccount.ID

	}
	return getAccountUsecase.Execute(
		filterFields,
	)
}

func (r *organizationResolver) PreviousEntity(ctx context.Context, obj *model.Organization) (*model.Organization, error) {
	var getOrganizationUsecase organizationpresentationusecaseinterfaces.GetOrganizationUsecase
	container.Make(&getOrganizationUsecase)

	var filterFields *model.OrganizationFilterFields
	if obj.PreviousEntity != nil {
		filterFields = &model.OrganizationFilterFields{}
		filterFields.ID = &obj.PreviousEntity.ID
	}
	return getOrganizationUsecase.Execute(
		filterFields,
	)
}

func (r *organizationResolver) CorrespondingLog(ctx context.Context, obj *model.Organization) (*model.Logging, error) {
	var getLoggingUsecase loggingpresentationusecaseinterfaces.GetLoggingUsecase
	container.Make(&getLoggingUsecase)

	var filterFields *model.LoggingFilterFields
	if obj.CorrespondingLog != nil {
		filterFields = &model.LoggingFilterFields{}
		filterFields.ID = &obj.CorrespondingLog.ID

	}
	return getLoggingUsecase.Execute(
		filterFields,
	)
}

func (r *queryResolver) Organizations(ctx context.Context, filterFields *model.OrganizationFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.Organization, error) {
	var getOrganizationsUsecase organizationpresentationusecaseinterfaces.GetAllOrganizationUsecase
	container.Make(&getOrganizationsUsecase)
	return getOrganizationsUsecase.Execute(
		organizationpresentationusecasetypes.GetAllOrganizationUsecaseInput{
			Context:       ctx,
			FilterFields:  filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

// Organization returns generated.OrganizationResolver implementation.
func (r *Resolver) Organization() generated.OrganizationResolver { return &organizationResolver{r} }

type organizationResolver struct{ *Resolver }
