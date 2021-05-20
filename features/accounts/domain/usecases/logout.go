package accountpresentationusecases

import (
	horeekaacoreerror "github.com/horeekaa/backend/core/errors/errors"
	horeekaacoreerrorenums "github.com/horeekaa/backend/core/errors/errors/enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type logoutUsecase struct {
	getAccountFromAuthDataRepo         accountdomainrepositoryinterfaces.GetAccountFromAuthData
	manageAccountDeviceTokenRepository accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository
}

func NewLogoutUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	manageAccountDeviceTokenRepository accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository,
) (accountpresentationusecaseinterfaces.LogoutUsecase, error) {
	return &logoutUsecase{
		getAccountFromAuthDataRepo,
		manageAccountDeviceTokenRepository,
	}, nil
}

func (logoutUcase *logoutUsecase) validation(input accountpresentationusecasetypes.LogoutUsecaseInput) (*accountpresentationusecasetypes.LogoutUsecaseInput, error) {
	if &input.Context == nil {
		return &accountpresentationusecasetypes.LogoutUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/logoutUsecase",
				nil,
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

	account, err := logoutUcase.getAccountFromAuthDataRepo.Execute(
		accountdomainrepositorytypes.GetAccountFromAuthDataInput{
			Context: validatedInput.Context,
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/logoutUsecase",
			err,
		)
	}
	if account == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			horeekaacoreerrorenums.AuthenticationTokenNotExist,
			401,
			"/logoutUsecase",
			nil,
		)
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
