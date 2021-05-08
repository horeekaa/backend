package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
	memberaccessrefpresentationusecasetypes "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *memberAccessRefResolver) SubmittingAccount(ctx context.Context, obj *model.MemberAccessRef) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)
	return getAccountUsecase.Execute(
		&model.AccountFilterFields{
			ID: &obj.SubmittingAccount.ID,
		},
	)
}

func (r *memberAccessRefResolver) ApprovingAccount(ctx context.Context, obj *model.MemberAccessRef) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)
	return getAccountUsecase.Execute(
		&model.AccountFilterFields{
			ID: &obj.ApprovingAccount.ID,
		},
	)
}

func (r *memberAccessRefResolver) PreviousEntity(ctx context.Context, obj *model.MemberAccessRef) (*model.MemberAccessRef, error) {
	var getMemberAccessRefUsecase memberaccessrefpresentationusecaseinterfaces.GetMemberAccessRefUsecase
	container.Make(&getMemberAccessRefUsecase)
	return getMemberAccessRefUsecase.Execute(
		&model.MemberAccessRefFilterFields{
			ID: &obj.PreviousEntity.ID,
		},
	)
}

func (r *memberAccessRefResolver) CorrespondingLog(ctx context.Context, obj *model.MemberAccessRef) (*model.Logging, error) {
	var getLoggingUsecase loggingpresentationusecaseinterfaces.GetLoggingUsecase
	container.Make(&getLoggingUsecase)
	return getLoggingUsecase.Execute(
		&model.LoggingFilterFields{
			ID: &obj.CorrespondingLog.ID,
		},
	)
}

func (r *mutationResolver) CreateMemberAccessRef(ctx context.Context, createMemberAccessRef model.CreateMemberAccessRef) (*model.MemberAccessRef, error) {
	var createMemberAccessRefUsecase memberaccessrefpresentationusecaseinterfaces.CreateMemberAccessRefUsecase
	container.Make(&createMemberAccessRefUsecase)
	return createMemberAccessRefUsecase.Execute(
		memberaccessrefpresentationusecasetypes.CreateMemberAccessRefUsecaseInput{
			Context:               ctx,
			CreateMemberAccessRef: &createMemberAccessRef,
		},
	)
}

func (r *mutationResolver) UpdateMemberAccessRef(ctx context.Context, updateMemberAccessRef model.UpdateMemberAccessRef) (*model.MemberAccessRef, error) {
	var updateMemberAccessRefUsecase memberaccessrefpresentationusecaseinterfaces.UpdateMemberAccessRefUsecase
	container.Make(&updateMemberAccessRefUsecase)
	return updateMemberAccessRefUsecase.Execute(
		memberaccessrefpresentationusecasetypes.UpdateMemberAccessRefUsecaseInput{
			Context:               ctx,
			UpdateMemberAccessRef: &updateMemberAccessRef,
		},
	)
}

func (r *queryResolver) MemberAccessRefs(ctx context.Context, filterFields *model.MemberAccessRefFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.MemberAccessRef, error) {
	var getMemberAccessRefsUsecase memberaccessrefpresentationusecaseinterfaces.GetAllMemberAccessRefUsecase
	container.Make(&getMemberAccessRefsUsecase)
	return getMemberAccessRefsUsecase.Execute(
		memberaccessrefpresentationusecasetypes.GetAllMemberAccessRefUsecaseInput{
			Context:       ctx,
			FilterFields:  filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

// MemberAccessRef returns generated.MemberAccessRefResolver implementation.
func (r *Resolver) MemberAccessRef() generated.MemberAccessRefResolver {
	return &memberAccessRefResolver{r}
}

type memberAccessRefResolver struct{ *Resolver }
