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

type getAccountUsecase struct {
	getAccountFromAuthDataRepository accountdomainrepositoryinterfaces.GetAccountFromAuthData
	getAccountRepository             accountdomainrepositoryinterfaces.GetAccountRepository
	getAccountMemberAccessRepository memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository

	getAccountAccessIdentity *model.MemberAccessRefOptionsInput
	pathIdentity             string
}

func NewGetAccountUsecase(
	getAccountFromAuthDataRepository accountdomainrepositoryinterfaces.GetAccountFromAuthData,
	getAccountRepository accountdomainrepositoryinterfaces.GetAccountRepository,
	getAccountMemberAccessRepository memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
) (accountpresentationusecaseinterfaces.GetAccountUsecase, error) {
	return &getAccountUsecase{
		getAccountFromAuthDataRepository,
		getAccountRepository,
		getAccountMemberAccessRepository,
		&model.MemberAccessRefOptionsInput{
			AccountAccesses: &model.AccountAccessesInput{
				AccountRead: func(b bool) *bool { return &b }(true),
			},
		},
		"GetAccountUsecase",
	}, nil
}

func (getAccUcase *getAccountUsecase) validation(
	input accountpresentationusecasetypes.GetAccountInput,
) (accountpresentationusecasetypes.GetAccountInput, error) {
	return input, nil
}

func (getAccUcase *getAccountUsecase) Execute(
	input accountpresentationusecasetypes.GetAccountInput,
) (*model.Account, error) {
	validatedInput, err := getAccUcase.validation(input)
	if err != nil {
		return nil, err
	}

	if validatedInput.Context != nil {
		account, err := getAccUcase.getAccountFromAuthDataRepository.Execute(
			accountdomainrepositorytypes.GetAccountFromAuthDataInput{
				Context: validatedInput.Context,
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				getAccUcase.pathIdentity,
				err,
			)
		}
		if account == nil {
			return nil, horeekaacoreerror.NewErrorObject(
				horeekaacoreerrorenums.AuthenticationError,
				getAccUcase.pathIdentity,
				nil,
			)
		}

		memberAccessRefTypeOrganizationsBased := model.MemberAccessRefTypeOrganizationsBased
		_, err = getAccUcase.getAccountMemberAccessRepository.Execute(
			memberaccessdomainrepositorytypes.GetAccountMemberAccessInput{
				MemberAccessFilterFields: &model.InternalMemberAccessFilterFields{
					Account:             &model.ObjectIDOnly{ID: &account.ID},
					MemberAccessRefType: &memberAccessRefTypeOrganizationsBased,
					Access:              getAccUcase.getAccountAccessIdentity,
				},
			},
		)
		if err != nil {
			return nil, horeekaacorefailuretoerror.ConvertFailure(
				getAccUcase.pathIdentity,
				err,
			)
		}
	}

	account, err := getAccUcase.getAccountRepository.Execute(
		validatedInput.FilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAccUcase.pathIdentity,
			err,
		)
	}
	return account, nil
}
