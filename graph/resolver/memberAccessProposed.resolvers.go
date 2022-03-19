package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *memberAccessProposedResolver) Account(ctx context.Context, obj *model.MemberAccessProposed) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)
	return getAccountUsecase.Execute(
		accountpresentationusecasetypes.GetAccountInput{
			FilterFields: &model.AccountFilterFields{
				ID: &obj.Account.ID,
			},
		},
	)
}

func (r *memberAccessProposedResolver) Organization(ctx context.Context, obj *model.MemberAccessProposed) (*model.Organization, error) {
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

func (r *memberAccessProposedResolver) DefaultAccessLatestUpdate(ctx context.Context, obj *model.MemberAccessProposed) (*model.MemberAccessRef, error) {
	var getMemberAccessRefUsecase memberaccessrefpresentationusecaseinterfaces.GetMemberAccessRefUsecase
	container.Make(&getMemberAccessRefUsecase)
	return getMemberAccessRefUsecase.Execute(
		&model.MemberAccessRefFilterFields{
			ID: &obj.DefaultAccessLatestUpdate.ID,
		},
	)
}

func (r *memberAccessProposedResolver) SubmittingAccount(ctx context.Context, obj *model.MemberAccessProposed) (*model.Account, error) {
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

func (r *memberAccessProposedResolver) RecentApprovingAccount(ctx context.Context, obj *model.MemberAccessProposed) (*model.Account, error) {
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

func (r *memberAccessProposedResolver) RecentLog(ctx context.Context, obj *model.MemberAccessProposed) (*model.Logging, error) {
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

// MemberAccessProposed returns generated.MemberAccessProposedResolver implementation.
func (r *Resolver) MemberAccessProposed() generated.MemberAccessProposedResolver {
	return &memberAccessProposedResolver{r}
}

type memberAccessProposedResolver struct{ *Resolver }
