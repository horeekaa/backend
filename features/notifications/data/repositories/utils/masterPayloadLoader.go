package notificationdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
	notificationdomainrepositoryloaderutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils/payloadLoaders"
	"github.com/horeekaa/backend/model"
)

type masterPayloadLoader struct {
	accountDataSource       databaseaccountdatasourceinterfaces.AccountDataSource
	personDataSource        databaseaccountdatasourceinterfaces.PersonDataSource
	invitationPayloadLoader notificationdomainrepositoryloaderutilityinterfaces.InvitationPayloadLoader
}

func NewMasterPayloadLoader(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
	invitationPayloadLoader notificationdomainrepositoryloaderutilityinterfaces.InvitationPayloadLoader,
) (notificationdomainrepositoryutilityinterfaces.MasterPayloadLoader, error) {
	return &masterPayloadLoader{
		accountDataSource,
		personDataSource,
		invitationPayloadLoader,
	}, nil
}

func (masterPayload *masterPayloadLoader) TransactionBody(
	operationOptions *mongodbcoretypes.OperationOptions,
	notification *model.DatabaseCreateNotification,
) (bool, error) {
	recipientLoaded := make(chan bool)
	payloadLoaded := make(chan bool)
	errChan := make(chan error)

	go func() {
		recipientAccount, err := masterPayload.accountDataSource.GetMongoDataSource().FindByID(
			notification.RecipientAccount.ID,
			operationOptions,
		)
		if err != nil {
			errChan <- err
			return
		}
		jsonRecipientAccount, _ := json.Marshal(recipientAccount)
		json.Unmarshal(jsonRecipientAccount, &notification.RecipientAccount)

		recipientPerson, err := masterPayload.personDataSource.GetMongoDataSource().FindByID(
			notification.RecipientAccount.Person.ID,
			operationOptions,
		)
		if err != nil {
			errChan <- err
			return
		}
		jsonRecipientPerson, _ := json.Marshal(recipientPerson)
		json.Unmarshal(jsonRecipientPerson, &notification.RecipientAccount.Person)

		recipientLoaded <- true
	}()

	go func() {
		_, err := masterPayload.invitationPayloadLoader.TransactionBody(
			operationOptions,
			notification,
		)
		if err != nil {
			errChan <- err
			return
		}
		payloadLoaded <- true
	}()

	for i := 0; i < 2; {
		select {
		case err := <-errChan:
			return false, err
		case _ = <-recipientLoaded:
			i++
		case _ = <-payloadLoaded:
			i++
		}
	}

	return true, nil
}
