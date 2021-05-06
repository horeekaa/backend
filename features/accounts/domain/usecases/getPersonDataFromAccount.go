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

type getPersonDataFromAccountUsecase struct {
	getAccountFromAuthDataRepo             accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountMemberAccessRepository       accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository
	getPersonDataFromAccountRepository     accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository
	getPersonDataFromAccountAccessIdentity *model.MemberAccessRefOptionsInput
}

func NewGetPersonDataFromAccountUsecase(
	getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountMemberAccessRepository accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
	getPersonDataFromAccountRepository accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
) (accountpresentationusecaseinterfaces.GetPersonDataFromAccountUsecase, error) {
	return &getPersonDataFromAccountUsecase{
		getAccountFromAuthDataRepo,
		getAccountMemberAccessRepository,
		getPersonDataFromAccountRepository,
		&model.MemberAccessRefOptionsInput{
			AccountAccesses: &model.AccountAccessesInput{
				AccountViewProfile: func(b bool) *bool { return &b }(true),
			},
		},
	}, nil
}

func (getPersonDataFromAccountUsecase *getPersonDataFromAccountUsecase) validation(
	input accountpresentationusecasetypes.GetPersonDataFromAccountInput,
) (*accountpresentationusecasetypes.GetPersonDataFromAccountInput, error) {
	if &input.User == nil && input.ViewProfileMode {
		return &accountpresentationusecasetypes.GetPersonDataFromAccountInput{},
			horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/getPersonDataFromAccount",
				nil,
			)
	}

	return &input, nil
}

func (getPersonDataFromAccountUsecase *getPersonDataFromAccountUsecase) Execute(
	input accountpresentationusecasetypes.GetPersonDataFromAccountInput,
) (*model.Person, error) {
	validatedInput, err := getPersonDataFromAccountUsecase.validation(input)
	if err != nil {
		return nil, err
	}

	if validatedInput.ViewProfileMode {
		account, err := getPersonDataFromAccountUsecase.getAccountFromAuthDataRepo.Execute(
			accountdomainrepositorytypes.GetAccountFromAuthDataInput{
				User:    validatedInput.User,
				Context: validatedInput.Context,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/getPersonDataFromAccount",
				err,
			)
		}
		if account == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationTokenNotExist,
				401,
				"/getPersonDataFromAccount",
				nil,
			)
		}
		validatedInput.Account = account

		_, err = getPersonDataFromAccountUsecase.getAccountMemberAccessRepository.Execute(
			accountdomainrepositorytypes.GetAccountMemberAccessInput{
				Account:                validatedInput.Account,
				MemberAccessRefType:    model.MemberAccessRefTypeAccountsBasics,
				MemberAccessRefOptions: *getPersonDataFromAccountUsecase.getPersonDataFromAccountAccessIdentity,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				"/getPersonDataFromAccount",
				err,
			)
		}
	}

	person, err := getPersonDataFromAccountUsecase.getPersonDataFromAccountRepository.Execute(validatedInput.Account)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getPersonDataFromAccount",
			err,
		)
	}
	return person, nil
}
