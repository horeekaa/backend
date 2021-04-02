package accountpresentationusecases

import (
	"errors"

	horeekaacoreerror "github.com/horeekaa/backend/core/_errors/usecaseErrors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/_errors/usecaseErrors/_enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/_errors/usecaseErrors/_failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type loginUsecase struct {
	manageAccountAuthenticationRepository accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository
	manageAccountDeviceTokenRepository    accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository
}

func NewLoginUsecase(
	manageAccountAuthenticationRepository accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository,
	manageAccountDeviceTokenRepository accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository,
) (accountpresentationusecaseinterfaces.LoginUsecase, error) {
	return &loginUsecase{
		manageAccountAuthenticationRepository,
		manageAccountDeviceTokenRepository,
	}, nil
}

func (loginUsecase *loginUsecase) Validation(input accountpresentationusecasetypes.LoginUsecaseInput) (accountpresentationusecasetypes.LoginUsecaseInput, error) {
	if &input.AuthHeader == nil {
		return accountpresentationusecasetypes.LoginUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"loginUsecase/",
				errors.New("loginUsecase/"),
			)
	}
	return input, nil
}

func (loginUsecase *loginUsecase) Execute(input accountpresentationusecasetypes.LoginUsecaseInput) (*model.Account, error) {
	account, err := loginUsecase.manageAccountAuthenticationRepository.RunTransaction(
		accountdomainrepositorytypes.ManageAccountAuthenticationInput{
			AuthHeader: input.AuthHeader,
			Context:    input.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"loginUsecase/",
			err,
		)
	}

	if &input.DeviceToken == nil {
		return account, nil
	}
	account, err = loginUsecase.manageAccountDeviceTokenRepository.Execute(
		accountdomainrepositorytypes.ManageAccountDeviceTokenInput{
			Account:                        account,
			DeviceToken:                    input.DeviceToken,
			ManageAccountDeviceTokenAction: accountdomainrepositorytypes.ManageAccountDeviceTokenActionInsert,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"loginUsecase/",
			err,
		)
	}

	return account, nil
}
