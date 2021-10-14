package addressregiongrouppresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	addressregiongroupdomainrepositorytypes "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories/types"
	addressregiongrouppresentationusecaseinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases"
	addressregiongrouppresentationusecasetypes "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllAddressRegionGroupUsecase struct {
	getAccountFromAuthDataRepo             accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo             memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllAddressRegionGroupRepo           addressregiongroupdomainrepositoryinterfaces.GetAllAddressRegionGroupRepository
	getAllAddressRegionGroupAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewGetAllAddressRegionGroupUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllAddressRegionGroupRepo addressregiongroupdomainrepositoryinterfaces.GetAllAddressRegionGroupRepository,
) (addressregiongrouppresentationusecaseinterfaces.GetAllAddressRegionGroupUsecase, error) {
	return &getAllAddressRegionGroupUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllAddressRegionGroupRepo,
		&model.MemberAccessRefOptionsInput{
			AddressRegionGroupAccesses: &model.AddressRegionGroupAccessesInput{
				AddressRegionGroupReadAll: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (getAlladdressRegionGroupUcase *getAllAddressRegionGroupUsecase) validation(input addressregiongrouppresentationusecasetypes.GetAllAddressRegionGroupUsecaseInput) (*addressregiongrouppresentationusecasetypes.GetAllAddressRegionGroupUsecaseInput, error) {
	if &input.Context == nil {
		return &addressregiongrouppresentationusecasetypes.GetAllAddressRegionGroupUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/getAllAddressRegionGroupUsecase",
				nil,
			)
	}
	return &input, nil
}

func (getAlladdressRegionGroupUcase *getAllAddressRegionGroupUsecase) Execute(
	input addressregiongrouppresentationusecasetypes.GetAllAddressRegionGroupUsecaseInput,
) ([]*model.AddressRegionGroup, error) {
	validatedInput, err := getAlladdressRegionGroupUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAlladdressRegionGroupUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllAddressRegionGroupUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/getAllAddressRegionGroupUsecase",
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	_, err = getAlladdressRegionGroupUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              getAlladdressRegionGroupUcase.getAllAddressRegionGroupAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllAddressRegionGroupUsecase",
			err,
		)
	}

	addressRegionGroups, err := getAlladdressRegionGroupUcase.getAllAddressRegionGroupRepo.Execute(
		addressregiongroupdomainrepositorytypes.GetAllAddressRegionGroupInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllAddressRegionGroupUsecase",
			err,
		)
	}

	return addressRegionGroups, nil
}
