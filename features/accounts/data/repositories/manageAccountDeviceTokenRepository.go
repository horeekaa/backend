package accountdomainrepositories

import (
	horeekaafailuretoerror "github.com/horeekaa/backend/core/_errors/usecaseErrors/_failureToError"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type manageAccountDeviceTokenRepository struct {
	accountDataSource                        databaseaccountdatasourceinterfaces.AccountDataSource
	manageAccountDeviceTokenUsecaseComponent accountdomainrepositoryinterfaces.ManageAccountDeviceTokenUsecaseComponent
}

func NewManageAccountDeviceTokenRepository(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
) (accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository, error) {
	return &manageAccountDeviceTokenRepository{
		accountDataSource: accountDataSource,
	}, nil
}

func (mgsAccDevToken *manageAccountDeviceTokenRepository) SetValidation(
	usecaseComponent accountdomainrepositoryinterfaces.ManageAccountDeviceTokenUsecaseComponent,
) (bool, error) {
	mgsAccDevToken.manageAccountDeviceTokenUsecaseComponent = usecaseComponent
	return true, nil
}

func (mgsAccDevToken *manageAccountDeviceTokenRepository) preExecute(input accountdomainrepositorytypes.ManageAccountDeviceTokenInput) (accountdomainrepositorytypes.ManageAccountDeviceTokenInput, error) {
	if mgsAccDevToken.manageAccountDeviceTokenUsecaseComponent == nil {
		return input, nil
	}
	return mgsAccDevToken.manageAccountDeviceTokenUsecaseComponent.Validation(input)
}

func (mgsAccDevToken *manageAccountDeviceTokenRepository) Execute(input accountdomainrepositorytypes.ManageAccountDeviceTokenInput) (*model.Account, error) {
	_, err := mgsAccDevToken.preExecute(input)
	if err != nil {
		return nil, err
	}
	account, err := mgsAccDevToken.accountDataSource.GetMongoDataSource().FindByID(input.Account.ID, &mongodbcoretypes.OperationOptions{})
	if err != nil {
		return nil, horeekaafailuretoerror.ConvertFailure(
			"/manageAccountDeviceTokenRepository",
			err,
		)
	}

	switch input.ManageAccountDeviceTokenAction {
	case accountdomainrepositorytypes.ManageAccountDeviceTokenActionInsert:
		if !funk.Contains(account.DeviceTokens, input.DeviceToken) {
			account.DeviceTokens = append(account.DeviceTokens, &input.DeviceToken)
		}
		break

	case accountdomainrepositorytypes.ManageAccountDeviceTokenActionRemove:
		index := funk.IndexOf(account.DeviceTokens, input.DeviceToken)
		account.DeviceTokens = append(account.DeviceTokens[:index], account.DeviceTokens[index+1:]...)
		break
	}

	updatedAccount, err := mgsAccDevToken.accountDataSource.GetMongoDataSource().Update(
		account.ID,
		&model.UpdateAccount{
			DeviceTokens: account.DeviceTokens,
		},
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaafailuretoerror.ConvertFailure(
			"/manageAccountDeviceTokenRepository",
			err,
		)
	}
	return updatedAccount, nil
}
