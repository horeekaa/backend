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
	createMemberAccessForAccountRepo   memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository
	getAccountMemberAccessRepository   memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	manageAccountDeviceTokenRepository accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository
	loginAccessIdentity                *model.MemberAccessRefOptionsInput
	pathIdentity                       string
}

func NewLoginUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	createAccountFromAuthDataRepo accountdomainrepositoryinterfaces.CreateAccountFromAuthDataRepository,
	createMemberAccessForAccountRepo memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository,
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
		"LoginUsecase",
	}, nil
}

func (loginUsecase *loginUsecase) validation(input accountpresentationusecasetypes.LoginUsecaseInput) (*accountpresentationusecasetypes.LoginUsecaseInput, error) {
	if &input.Context == nil {
		return &accountpresentationusecasetypes.LoginUsecaseInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				loginUsecase.pathIdentity,
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
			loginUcase.pathIdentity,
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
				loginUcase.pathIdentity,
				err,
			)
		}

		_, err = loginUcase.createMemberAccessForAccountRepo.RunTransaction(
			&model.InternalCreateMemberAccess{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				SubmittingAccount:   &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: model.MemberAccessRefTypeAccountsBasics,
				InvitationAccepted:  func(b bool) *bool { return &b }(true),
				ProposalStatus: func(ps model.EntityProposalStatus) *model.EntityProposalStatus {
					return &ps
				}(model.EntityProposalStatusApproved),
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				loginUcase.pathIdentity,
				err,
			)
		}
	}

	memberAccessRefTypeAccountsBasics := model.MemberAccessRefTypeAccountsBasics
	_, err = loginUcase.getAccountMemberAccessRepository.Execute(
		memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
			MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
				Account:             &model.ObjectIDOnly{ID: &account.ID},
				MemberAccessRefType: &memberAccessRefTypeAccountsBasics,
				Access:              loginUcase.loginAccessIdentity,
			},
		},
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			loginUcase.pathIdentity,
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
			loginUcase.pathIdentity,
			err,
		)
	}

	return account, nil
}
