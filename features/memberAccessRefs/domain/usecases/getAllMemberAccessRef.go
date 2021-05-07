package memberaccessrefpresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories/types"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
	memberaccessrefpresentationusecasetypes "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type getAllMemberAccessRefUsecase struct {
	getAccountFromAuthDataRepo          accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo          accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllMemberAccessRefRepo           memberaccessrefdomainrepositoryinterfaces.GetAllMemberAccessRefRepository
	getAllMemberAccessRefAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewGetAllMemberAccessRefUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.GetAllMemberAccessRefRepository,
) (memberaccessrefpresentationusecaseinterfaces.GetAllMemberAccessRefUsecase, error) {
	return &getAllMemberAccessRefUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllMemberAccessRefRepo,
		&model.MemberAccessRefOptionsInput{
			MemberAccessRefAccesses: &model.MemberAccessRefAccessesInput{
				MemberAccessRefRead: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (getAllMmbAccRefUcase *getAllMemberAccessRefUsecase) validation(input memberaccessrefpresentationusecasetypes.GetAllMemberAccessRefUsecaseInput) (*memberaccessrefpresentationusecasetypes.GetAllMemberAccessRefUsecaseInput, error) {
	if &input.Context == nil {
		return &memberaccessrefpresentationusecasetypes.GetAllMemberAccessRefUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/getAllMemberAccessRefUsecase",
				nil,
			)
	}
	return &input, nil
}

func (getAllMmbAccRefUcase *getAllMemberAccessRefUsecase) Execute(
	input memberaccessrefpresentationusecasetypes.GetAllMemberAccessRefUsecaseInput,
) ([]*model.MemberAccessRef, error) {
	validatedInput, err := getAllMmbAccRefUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := getAllMmbAccRefUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllMemberAccessRefUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/getAllMemberAccessRefUsecase",
			nil,
		)
	}

	_, err = getAllMmbAccRefUcase.getAccountMemberAccessRepo.Execute(
		accountdomainrepositorytypes.GetAccountMemberAccessInput{
			Account:                account,
			MemberAccessRefType:    model.MemberAccessRefTypeOrganizationsBased,
			MemberAccessRefOptions: *getAllMmbAccRefUcase.getAllMemberAccessRefAccessIdentity,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllMemberAccessRefUsecase",
			err,
		)
	}

	memberAccessRefs, err := getAllMmbAccRefUcase.getAllMemberAccessRefRepo.Execute(
		memberaccessrefdomainrepositorytypes.GetAllMemberAccessRefInput{
			FilterFields:  validatedInput.FilterFields,
			PaginationOpt: validatedInput.PaginationOps,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAllMemberAccessRefUsecase",
			err,
		)
	}

	return memberAccessRefs, nil
}
