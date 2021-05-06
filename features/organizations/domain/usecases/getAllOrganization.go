package organizationpresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationdomainrepositorytypes "github.com/horeekaa/backend/features/organizations/domain/repositories/types"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
	organizationpresentationusecasetypes "github.com/horeekaa/backend/features/organizations/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type getAllOrganizationUsecase struct {
	getAccountFromAuthDataRepo       accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo       accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllOrganizationRepo           organizationdomainrepositoryinterfaces.GetAllOrganizationRepository
	getAllOrganizationAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewGetAllOrganizationUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllOrganizationRepo organizationdomainrepositoryinterfaces.GetAllOrganizationRepository,
) (organizationpresentationusecaseinterfaces.GetAllOrganizationUsecase, error) {
	return &getAllOrganizationUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllOrganizationRepo,
		&model.MemberAccessRefOptionsInput{
			OrganizationAccesses: &model.OrganizationAccessesInput{
				OrganizationRead: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (getAllMmbAccRefUcase *getAllOrganizationUsecase) validation(input organizationpresentationusecasetypes.GetAllOrganizationUsecaseInput) (*organizationpresentationusecasetypes.GetAllOrganizationUsecaseInput, error) {
	if &input.User == nil {
		return &organizationpresentationusecasetypes.GetAllOrganizationUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/getAllOrganizationUsecase",
				nil,
			)
	}
	return &input, nil
}

func (getAllMmbAccRefUcase *getAllOrganizationUsecase) Execute(
	input organizationpresentationusecasetypes.GetAllOrganizationUsecaseInput,
) ([]*model.Organization, error) {
	validatedInput, err := getAllMmbAccRefUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllMmbAccRefUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			User:    validatedInput.User,
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllOrganizationUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/getAllOrganizationUsecase",
			nil,
		)
	}

	_, err = getAllMmbAccRefUcase.getAccountMemberAccessRepo.Execute(
		accountdomainrepositorytypes.GetAccountMemberAccessInput{
			Account:                account,
			MemberAccessRefType:    model.MemberAccessRefTypeOrganizationsBased,
			MemberAccessRefOptions: *getAllMmbAccRefUcase.getAllOrganizationAccessIdentity,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllOrganizationUsecase",
			err,
		)
	}

	organizations, err := getAllMmbAccRefUcase.getAllOrganizationRepo.Execute(
		organizationdomainrepositorytypes.GetAllOrganizationInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllOrganizationUsecase",
			err,
		)
	}

	return organizations, nil
}
