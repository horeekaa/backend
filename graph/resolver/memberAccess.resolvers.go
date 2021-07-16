package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
	memberaccesspresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases"
	memberaccesspresentationusecasetypes "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases/types"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *memberAccessResolver) Account(ctx context.Context, obj *model.MemberAccess) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)
	return getAccountUsecase.Execute(
		&model.AccountFilterFields{
			ID: &obj.Account.ID,
		},
	)
}

func (r *memberAccessResolver) OrganizationLatestUpdate(ctx context.Context, obj *model.MemberAccess) (*model.Organization, error) {
	var getOrganizationUsecase organizationpresentationusecaseinterfaces.GetOrganizationUsecase
	container.Make(&getOrganizationUsecase)

	var filterFields *model.OrganizationFilterFields
	if obj.OrganizationLatestUpdate != nil {
		filterFields = &model.OrganizationFilterFields{}
		filterFields.ID = &obj.OrganizationLatestUpdate.ID
	}
	return getOrganizationUsecase.Execute(
		filterFields,
	)
}

func (r *memberAccessResolver) DefaultAccessLatestUpdate(ctx context.Context, obj *model.MemberAccess) (*model.MemberAccessRef, error) {
	var getMemberAccessRefUsecase memberaccessrefpresentationusecaseinterfaces.GetMemberAccessRefUsecase
	container.Make(&getMemberAccessRefUsecase)
	return getMemberAccessRefUsecase.Execute(
		&model.MemberAccessRefFilterFields{
			ID: &obj.DefaultAccess.ID,
		},
	)
}

func (r *memberAccessResolver) SubmittingAccount(ctx context.Context, obj *model.MemberAccess) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)
	return getAccountUsecase.Execute(
		&model.AccountFilterFields{
			ID: &obj.SubmittingAccount.ID,
		},
	)
}

func (r *memberAccessResolver) RecentApprovingAccount(ctx context.Context, obj *model.MemberAccess) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)

	var filterFields *model.AccountFilterFields
	if obj.RecentApprovingAccount != nil {
		filterFields = &model.AccountFilterFields{}
		filterFields.ID = &obj.RecentApprovingAccount.ID
	}
	return getAccountUsecase.Execute(
		filterFields,
	)
}

func (r *memberAccessResolver) RecentLog(ctx context.Context, obj *model.MemberAccess) (*model.Logging, error) {
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

func (r *mutationResolver) CreateMemberAccess(ctx context.Context, createMemberAccess model.CreateMemberAccess) (*model.MemberAccess, error) {
	var createMemberAccessUsecase memberaccesspresentationusecaseinterfaces.CreateMemberAccessUsecase
	container.Make(&createMemberAccessUsecase)
	return createMemberAccessUsecase.Execute(
		memberaccesspresentationusecasetypes.CreateMemberAccessUsecaseInput{
			Context:            ctx,
			CreateMemberAccess: &createMemberAccess,
		},
	)
}

func (r *mutationResolver) UpdateMemberAccess(ctx context.Context, updateMemberAccess model.UpdateMemberAccess) (*model.MemberAccess, error) {
	var updateMemberAccessUsecase memberaccesspresentationusecaseinterfaces.UpdateMemberAccessUsecase
	container.Make(&updateMemberAccessUsecase)
	return updateMemberAccessUsecase.Execute(
		memberaccesspresentationusecasetypes.UpdateMemberAccessUsecaseInput{
			Context:            ctx,
			UpdateMemberAccess: &updateMemberAccess,
		},
	)
}

func (r *queryResolver) MemberAccesses(ctx context.Context, filterFields model.MemberAccessFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.MemberAccess, error) {
	var getMemberAccesssUsecase memberaccesspresentationusecaseinterfaces.GetAllMemberAccessUsecase
	container.Make(&getMemberAccesssUsecase)
	return getMemberAccesssUsecase.Execute(
		memberaccesspresentationusecasetypes.GetAllMemberAccessUsecaseInput{
			Context:       ctx,
			FilterFields:  &filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

// MemberAccess returns generated.MemberAccessResolver implementation.
func (r *Resolver) MemberAccess() generated.MemberAccessResolver { return &memberAccessResolver{r} }

type memberAccessResolver struct{ *Resolver }
