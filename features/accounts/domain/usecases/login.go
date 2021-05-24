package accountpresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type loginUsecase struct {
	getAccountFromAuthDataRepo         accountdomainrepositoryinterfaces.GetAccountFromAuthData
	createAccountFromAuthDataRepo      accountdomainrepositoryinterfaces.CreateAccountFromAuthDataRepository
	createMemberAccessForAccountRepo   memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountRepository
	getAccountMemberAccessRepository   memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	manageAccountDeviceTokenRepository accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository
	loginAccessIdentity                *model.MemberAccessRefOptionsInput
}

func NewLoginUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	createAccountFromAuthDataRepo accountdomainrepositoryinterfaces.CreateAccountFromAuthDataRepository,
	createMemberAccessForAccountRepo memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountRepository,
	getAccountMemberAccessRepository memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	manageAccountDeviceTokenRepository accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository,
) (accountpresentationusecaseinterfaces.LoginUsecase, error) {
	return &loginUsecase{
		getAccountFromAuthDataRepo,
		createAccountFromAuthDataRepo,
		createMemberAccessForAccountRepo,
		getAccountMemberAccessRepository,
		manageAccountDeviceTokenRepository,
		&model.MemberAccessRefOptionsInput{
			AccountAccesses: &model.AccountAccessesInput{
				AccountLogin: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (loginUsecase *loginUsecase) validation(input accountpresentationusecasetypes.LoginUsecaseInput) (*accountpresentationusecasetypes.LoginUsecaseInput, error) {
	if &input.Context == nil {
		return &accountpresentationusecasetypes.LoginUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/loginUsecase",
				nil,
			)
	}
	return &input, nil
}

func (loginUcase *loginUsecase) Execute(input accountpresentationusecasetypes.LoginUsecaseInput) (*model.Account, error) {
	validatedInput, err := loginUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := loginUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/loginUsecase",
			err,
		)
	}
	if account == nil {
		account, err = loginUcase.createAccountFromAuthDataRepo.RunTransaction(
			accountdomainrepositorytypes.CreateAccountFromAuthDataInput{
				Context: validatedInput.Context,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/loginUsecase",
				err,
			)
		}

		_, err = loginUcase.createMemberAccessForAccountRepo.Execute(
			&model.CreateMemberAccess{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: model.MemberAccessRefTypeAccountsBasics,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/loginUsecase",
				err,
			)
		}
	}

	_, err = loginUcase.getAccountMemberAccessRepository.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			Account:                account,
			MemberAccessRefType:    model.MemberAccessRefTypeAccountsBasics,
			MemberAccessRefOptions: *loginUcase.loginAccessIdentity,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/loginUsecase",
			err,
		)
	}

	account, err = loginUcase.manageAccountDeviceTokenRepository.Execute(
		accountdomainrepositorytypes.ManageAccountDeviceTokenInput{
			Account:                        account,
			DeviceToken:                    validatedInput.DeviceToken,
			ManageAccountDeviceTokenAction: accountdomainrepositorytypes.ManageAccountDeviceTokenActionInsert,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/loginUsecase",
			err,
		)
	}

	return account, nil
}
