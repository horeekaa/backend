package accountdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type manageAccountDeviceTokenRepository struct {
	accountDataSource                        databaseaccountdatasourceinterfaces.AccountDataSource
	manageAccountDeviceTokenUsecaseComponent accountdomainrepositoryinterfaces.ManageAccountDeviceTokenUsecaseComponent
	pathIdentity                             string
}

func NewManageAccountDeviceTokenRepository(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
) (accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository, error) {
	return &manageAccountDeviceTokenRepository{
		accountDataSource: accountDataSource,
		pathIdentity:      "ManageAccountDeviceTokenRepository",
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
	validatedInput, err := mgsAccDevToken.preExecute(input)
	if err != nil {
		return nil, err
	}

	switch input.ManageAccountDeviceTokenAction {
	case accountdomainrepositorytypes.ManageAccountDeviceTokenActionInsert:
		if !funk.Contains(validatedInput.Account.DeviceTokens, validatedInput.DeviceToken) {
			validatedInput.Account.DeviceTokens = append(validatedInput.Account.DeviceTokens, validatedInput.DeviceToken)
		}
		break

	case accountdomainrepositorytypes.ManageAccountDeviceTokenActionRemove:
		index := funk.IndexOf(validatedInput.Account.DeviceTokens, validatedInput.DeviceToken)
		validatedInput.Account.DeviceTokens = append(validatedInput.Account.DeviceTokens[:index], validatedInput.Account.DeviceTokens[index+1:]...)
		break
	}

	updatedAccount, err := mgsAccDevToken.accountDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": validatedInput.Account.ID,
		},
		&model.DatabaseUpdateAccount{
			DeviceTokens: validatedInput.Account.DeviceTokens,
		},
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			mgsAccDevToken.pathIdentity,
			err,
		)
	}
	return updatedAccount, nil
}
