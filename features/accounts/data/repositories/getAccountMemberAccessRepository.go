package accountdomainrepositories

import (
	"errors"

	horeekaacorefailure "github.com/horeekaa/backend/core/_errors/serviceFailures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/_errors/serviceFailures/_enums"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/_errors/usecaseErrors/_failureToError"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAccountMemberAccessRepository struct {
	accountDataSource                      databaseaccountdatasourceinterfaces.AccountDataSource
	memberAccessDataSource                 databaseaccountdatasourceinterfaces.MemberAccessDataSource
	getAccountMemberAccessUsecaseComponent accountdomainrepositoryinterfaces.GetAccountMemberAccessUsecaseComponent
}

func NewGetAccountMemberAccessRepository(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	memberAccessDataSource databaseaccountdatasourceinterfaces.MemberAccessDataSource,
) (accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository, error) {
	return &getAccountMemberAccessRepository{
		accountDataSource:      accountDataSource,
		memberAccessDataSource: memberAccessDataSource,
	}, nil
}

func (getAccountMemberAccess *getAccountMemberAccessRepository) SetValidation(
	usecaseComponent accountdomainrepositoryinterfaces.GetAccountMemberAccessUsecaseComponent,
) (bool, error) {
	getAccountMemberAccess.getAccountMemberAccessUsecaseComponent = usecaseComponent
	return true, nil
}

func (getAccountMemberAccess *getAccountMemberAccessRepository) preExecute(
	input accountdomainrepositorytypes.GetAccountMemberAccessInput,
) (accountdomainrepositorytypes.GetAccountMemberAccessInput, error) {
	if &input.Account.ID == nil {
		return accountdomainrepositorytypes.GetAccountMemberAccessInput{}, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.AccountIDNeededToRetrievePersonData,
			"/getAccountMemberAccess",
			errors.New(horeekaacorefailureenums.AccountIDNeededToRetrievePersonData),
		)
	}
	if getAccountMemberAccess.getAccountMemberAccessUsecaseComponent == nil {
		return input, nil
	}
	return getAccountMemberAccess.getAccountMemberAccessUsecaseComponent.Validation(input)
}

func (getAccountMemberAccess *getAccountMemberAccessRepository) Execute(input accountdomainrepositorytypes.GetAccountMemberAccessInput) (*model.MemberAccess, error) {
	preExecuteOutput, err := getAccountMemberAccess.preExecute(input)
	if err != nil {
		return nil, err
	}

	account, err := getAccountMemberAccess.accountDataSource.GetMongoDataSource().FindByID(preExecuteOutput.Account.ID, &mongodbcoretypes.OperationOptions{})
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getPersonDataFromAccount",
			err,
		)
	}

	memberAccess, err := getAccountMemberAccess.memberAccessDataSource.GetMongoDataSource().FindOne(
		map[string]interface{}{
			"account":             &model.Account{ID: account.ID},
			"memberAccessRefType": input.MemberAccessRefType,
			"access":              input.MemberAccessRefOptions,
			"status":              model.MemberAccessStatusActive,
		},
		&mongodbcoretypes.OperationOptions{},
	)
	if memberAccess == nil {
		return nil, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.FeatureNotAccessibleByAccount,
			"/getPersonDataFromAccount",
			errors.New(horeekaacorefailureenums.FeatureNotAccessibleByAccount),
		)
	}
	return memberAccess, nil
}
