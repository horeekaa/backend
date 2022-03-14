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
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllMemberAccessRefUsecase struct {
	getAccountFromAuthDataRepo          accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepo          memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getAllMemberAccessRefRepo           memberaccessrefdomainrepositoryinterfaces.GetAllMemberAccessRefRepository
	getAllMemberAccessRefAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity                        string
}

func NewGetAllMemberAccessRefUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getAllMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.GetAllMemberAccessRefRepository,
) (memberaccessrefpresentationusecaseinterfaces.GetAllMemberAccessRefUsecase, error) {
	return &getAllMemberAccessRefUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepo,
		getAllMemberAccessRefRepo,
		&model.MemberAccessRefOptionsInput{
			MemberAccessRefAccesses: &model.MemberAccessRefAccessesInput{
				MemberAccessRefReadAll: func(b bool) *bool { return &b }(true),
			},
		},
		"GetAllMemberAccessRefUsecase",
	}, nil
}

func (getAllMmbAccRefUcase *getAllMemberAccessRefUsecase) validation(input memberaccessrefpresentationusecasetypes.GetAllMemberAccessRefUsecaseInput) (*memberaccessrefpresentationusecasetypes.GetAllMemberAccessRefUsecaseInput, error) {
	if &input.Context == nil {
		return &memberaccessrefpresentationusecasetypes.GetAllMemberAccessRefUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAllMmbAccRefUcase.pathIdentity,
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
			getAllMmbAccRefUcase.pathIdentity,
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationError,
			getAllMmbAccRefUcase.pathIdentity,
			nil,
		)
	}

	_, err = getAllMmbAccRefUcase.getAccountMemberAccessRepo.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account: &model.ObjectIDOnly{ID: &account.ID},
				Access:  getAllMmbAccRefUcase.getAllMemberAccessRefAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAllMmbAccRefUcase.pathIdentity,
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
			getAllMmbAccRefUcase.pathIdentity,
			err,
		)
	}

	return memberAccessRefs, nil
}
