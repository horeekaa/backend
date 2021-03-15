package databasereposervices

import (
	horeekaaexceptiontofailure "github.com/horeekaa/backend/_errors/serviceFailures/_exceptionToFailure"
	model "github.com/horeekaa/backend/model"
	databaseinstancereferences "github.com/horeekaa/backend/repositories/databaseClient/instanceReferences/repos"
	databaseservicerepointerfaces "github.com/horeekaa/backend/services/database/interfaces/repos"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type accountService struct {
	accountRepo *databaseinstancereferences.AccountRepo
}

func NewAccountService(accountRepo databaseinstancereferences.AccountRepo) (databaseservicerepointerfaces.AccountService, error) {
	return &accountService{
		&accountRepo,
	}, nil
}

func (accountSvc *accountService) FindByID(ID interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.Account, chan error) {
	accountChn := make(chan *model.Account)
	errorChn := make(chan error)

	go func() {
		account, err := (*accountSvc.accountRepo.Instance).FindByID(ID, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/accountService/FindByID",
				&err,
			)
			return
		}

		accountChn <- account
	}()

	return accountChn, errorChn
}

func (accountSvc *accountService) FindOne(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.Account, chan error) {
	accountChn := make(chan *model.Account)
	errorChn := make(chan error)

	go func() {
		account, err := (*accountSvc.accountRepo.Instance).FindOne(query, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/accountService/FindOne",
				&err,
			)
			return
		}

		accountChn <- account
	}()

	return accountChn, errorChn
}

func (accountSvc *accountService) Find(query map[string]interface{}, serviceOptions *databaseserviceoperations.ServiceOptions) (chan []*model.Account, chan error) {
	accountsChn := make(chan []*model.Account)
	errorChn := make(chan error)

	go func() {
		accounts, err := (*accountSvc.accountRepo.Instance).Find(query, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/accountService/Find",
				&err,
			)
			return
		}

		accountsChn <- accounts
	}()

	return accountsChn, errorChn
}

func (accountSvc *accountService) Create(input *model.CreateAccount, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.Account, chan error) {
	accountChn := make(chan *model.Account)
	errorChn := make(chan error)

	go func() {
		account, err := (*accountSvc.accountRepo.Instance).Create(input, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/accountService/Create",
				&err,
			)
			return
		}

		accountChn <- account
	}()

	return accountChn, errorChn
}

func (accountSvc *accountService) Update(ID interface{}, updateData *model.UpdateAccount, serviceOptions *databaseserviceoperations.ServiceOptions) (chan *model.Account, chan error) {
	accountChn := make(chan *model.Account)
	errorChn := make(chan error)

	go func() {
		account, err := (*accountSvc.accountRepo.Instance).Update(ID, updateData, (*serviceOptions).OperationOptions)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/accountService/Update",
				&err,
			)
			return
		}

		accountChn <- account
	}()

	return accountChn, errorChn
}
