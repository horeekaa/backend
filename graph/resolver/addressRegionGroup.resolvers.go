package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	addressregiongrouppresentationusecaseinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases"
	addressregiongrouppresentationusecasetypes "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases/types"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *addressRegionGroupResolver) SubmittingAccount(ctx context.Context, obj *model.AddressRegionGroup) (*model.Account, error) {
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

func (r *addressRegionGroupResolver) RecentApprovingAccount(ctx context.Context, obj *model.AddressRegionGroup) (*model.Account, error) {
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

func (r *addressRegionGroupResolver) RecentLog(ctx context.Context, obj *model.AddressRegionGroup) (*model.Logging, error) {
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

func (r *mutationResolver) CreateAddressRegionGroup(ctx context.Context, createAddressRegionGroup model.CreateAddressRegionGroup) (*model.AddressRegionGroup, error) {
	var createAddressRegionGroupUsecase addressregiongrouppresentationusecaseinterfaces.CreateAddressRegionGroupUsecase
	container.Make(&createAddressRegionGroupUsecase)
	return createAddressRegionGroupUsecase.Execute(
		addressregiongrouppresentationusecasetypes.CreateAddressRegionGroupUsecaseInput{
			Context:                  ctx,
			CreateAddressRegionGroup: &createAddressRegionGroup,
		},
	)
}

func (r *mutationResolver) UpdateAddressRegionGroup(ctx context.Context, updateAddressRegionGroup model.UpdateAddressRegionGroup) (*model.AddressRegionGroup, error) {
	var updateAddressRegionGroupUsecase addressregiongrouppresentationusecaseinterfaces.UpdateAddressRegionGroupUsecase
	container.Make(&updateAddressRegionGroupUsecase)
	return updateAddressRegionGroupUsecase.Execute(
		addressregiongrouppresentationusecasetypes.UpdateAddressRegionGroupUsecaseInput{
			Context:                  ctx,
			UpdateAddressRegionGroup: &updateAddressRegionGroup,
		},
	)
}

func (r *queryResolver) AddressRegionGroups(ctx context.Context, filterFields model.AddressRegionGroupFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.AddressRegionGroup, error) {
	var getAllAddressRegionGroupUsecase addressregiongrouppresentationusecaseinterfaces.GetAllAddressRegionGroupUsecase
	container.Make(&getAllAddressRegionGroupUsecase)
	return getAllAddressRegionGroupUsecase.Execute(
		addressregiongrouppresentationusecasetypes.GetAllAddressRegionGroupUsecaseInput{
			Context:       ctx,
			FilterFields:  &filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

// AddressRegionGroup returns generated.AddressRegionGroupResolver implementation.
func (r *Resolver) AddressRegionGroup() generated.AddressRegionGroupResolver {
	return &addressRegionGroupResolver{r}
}

type addressRegionGroupResolver struct{ *Resolver }
