package notificationdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type masterPayloadLoader struct {
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource
	personDataSource  databaseaccountdatasourceinterfaces.PersonDataSource
}

func NewMasterPayloadLoader(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
) (notificationdomainrepositoryutilityinterfaces.MasterPayloadLoader, error) {
	return &masterPayloadLoader{
		accountDataSource,
		personDataSource,
	}, nil
}

func (masterPayload *masterPayloadLoader) TransactionBody(
	operationOptions *mongodbcoretypes.OperationOptions,
	notification *model.DatabaseCreateNotification,
) (bool, error) {
	recipientAccount, err := masterPayload.accountDataSource.GetMongoDataSource().FindByID(
		notification.RecipientAccount.ID,
		operationOptions,
	)
	if err != nil {
		return false, err
	}
	jsonRecipientAccount, _ := json.Marshal(recipientAccount)
	json.Unmarshal(jsonRecipientAccount, &notification.RecipientAccount)

	recipientPerson, err := masterPayload.personDataSource.GetMongoDataSource().FindByID(
		notification.RecipientAccount.Person.ID,
		operationOptions,
	)
	if err != nil {
		return false, err
	}
	jsonRecipientPerson, _ := json.Marshal(recipientPerson)
	json.Unmarshal(jsonRecipientPerson, &notification.RecipientAccount.Person)

	return true, nil
}
