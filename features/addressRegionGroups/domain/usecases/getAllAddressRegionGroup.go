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
	pathIdentity                           string
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
		"GetAllAddressRegionGroupUsecase",
	}, nil
}

func (getAllAddressRegionGroupUcase *getAllAddressRegionGroupUsecase) validation(input addressregiongrouppresentationusecasetypes.GetAllAddressRegionGroupUsecaseInput) (*addressregiongrouppresentationusecasetypes.GetAllAddressRegionGroupUsecaseInput, error) {
	if &input.Context == nil {
		return &addressregiongrouppresentationusecasetypes.GetAllAddressRegionGroupUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAllAddressRegionGroupUcase.pathIdentity,
				nil,
			)
	}
	return &input, nil
}

func (getAllAddressRegionGroupUcase *getAllAddressRegionGroupUsecase) Execute(
	input addressregiongrouppresentationusecasetypes.GetAllAddressRegionGroupUsecaseInput,
) ([]*model.AddressRegionGroup, error) {
	validatedInput, err := getAllAddressRegionGroupUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllAddressRegionGroupUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllAddressRegionGroupUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			getAllAddressRegionGroupUcase.pathIdentity,
			nil,
		)
	}

	memberAccessRefTypeOrgBased := model.MemberAccessRefTypeOrganizationsBased
	_, err = getAllAddressRegionGroupUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeOrgBased,
				Access:              getAllAddressRegionGroupUcase.getAllAddressRegionGroupAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllAddressRegionGroupUcase.pathIdentity,
			err,
		)
	}

	addressRegionGroups, err := getAllAddressRegionGroupUcase.getAllAddressRegionGroupRepo.Execute(
		addressregiongroupdomainrepositorytypes.GetAllAddressRegionGroupInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllAddressRegionGroupUcase.pathIdentity,
			err,
		)
	}

	return addressRegionGroups, nil
}
