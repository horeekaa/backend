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

type logoutUsecase struct {
	manageAccountAuthenticationRepository accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository
	manageAccountDeviceTokenRepository    accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository
}

func NewLogoutUsecase(
	manageAccountAuthenticationRepository accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository,
	manageAccountDeviceTokenRepository accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository,
) (accountpresentationusecaseinterfaces.LogoutUsecase, error) {
	return &logoutUsecase{
		manageAccountAuthenticationRepository,
		manageAccountDeviceTokenRepository,
	}, nil
}

func (logoutUcase *logoutUsecase) validation(input accountpresentationusecasetypes.LogoutUsecaseInput) (*accountpresentationusecasetypes.LogoutUsecaseInput, error) {
	if &input.AuthHeader == nil {
		return &accountpresentationusecasetypes.LogoutUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/logoutUsecase",
				errors.New(horeekaacoreerrorenums.AuthenticationTokenNotExist),
			)
	}
	return &input, nil
}

func (logoutUcase *logoutUsecase) Execute(
	input accountpresentationusecasetypes.LogoutUsecaseInput,
) (*model.Account, error) {
	validatedInput, err := logoutUcase.validation(input)
	if err != nil {
		return nil, err
	}

	account, err := logoutUcase.manageAccountAuthenticationRepository.RunTransaction(
		accountdomainrepositorytypes.ManageAccountAuthenticationInput{
			AuthHeader: validatedInput.AuthHeader,
			Context:    validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/logoutUsecase",
			err,
		)
	}

	if &input.DeviceToken == nil {
		return account, nil
	}
	account, err = logoutUcase.manageAccountDeviceTokenRepository.Execute(
		accountdomainrepositorytypes.ManageAccountDeviceTokenInput{
			Account:                        account,
			DeviceToken:                    validatedInput.DeviceToken,
			ManageAccountDeviceTokenAction: accountdomainrepositorytypes.ManageAccountDeviceTokenActionRemove,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/logoutUsecase",
			err,
		)
	}

	return account, nil
}
