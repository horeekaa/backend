package accountpresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"

	"github.com/horeekaa/backend/model"
)

type getAccountUsecase struct {
	getAccountRepository accountdomainrepositoryinterfaces.GetAccountRepository
}

func NewGetAccountUsecase(
	getAccountRepository accountdomainrepositoryinterfaces.GetAccountRepository,
) (accountpresentationusecaseinterfaces.GetAccountUsecase, error) {
	return &getAccountUsecase{
		getAccountRepository,
	}, nil
}

func (getAccFromMemberAccessRefUsecase *getAccountUsecase) validation(
	input *model.AccountFilterFields,
) (*model.AccountFilterFields, error) {
	return input, nil
}

func (getAccFromMemberAccessRefUsecase *getAccountUsecase) Execute(
	filterFields *model.AccountFilterFields,
) (*model.Account, error) {
	validatedFilterFields, err := getAccFromMemberAccessRefUsecase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	account, err := getAccFromMemberAccessRefUsecase.getAccountRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAccount",
			err,
		)
	}
	return account, nil
}
