package notificationdomainrepositoryloaderutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryloaderutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils/payloadLoaders"
	"github.com/horeekaa/backend/model"
)

type invitationPayloadLoader struct {
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource
	personDataSource  databaseaccountdatasourceinterfaces.PersonDataSource
}

func NewInvitationPayloadLoader(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
) (notificationdomainrepositoryloaderutilityinterfaces.InvitationPayloadLoader, error) {
	return &invitationPayloadLoader{
		accountDataSource: accountDataSource,
		personDataSource:  personDataSource,
	}, nil
}

func (invitationPyload *invitationPayloadLoader) TransactionBody(
	operationOptions *mongodbcoretypes.OperationOptions,
	notification *model.DatabaseCreateNotification,
) (bool, error) {
	memberAccess := &model.MemberAccessForNotifPayloadInput{}
	switch notification.NotificationCategory {
	case model.NotificationCategoryOrgInvitationRequest:
		memberAccess = notification.PayloadOptions.InvitationRequestPayload.MemberAccess
		break

	case model.NotificationCategoryOrgInvitationAccepted:
		memberAccess = notification.PayloadOptions.InvitationAcceptedPayload.MemberAccess
		break

	default:
		return false, nil
	}

	submittingAccountLoadedChan := make(chan bool)
	invitedAccountLoadedChan := make(chan bool)
	errChan := make(chan error)

	go func() {
		submittingAcc, err := invitationPyload.accountDataSource.GetMongoDataSource().FindByID(
			memberAccess.SubmittingAccount.ID,
			operationOptions,
		)
		if err != nil {
			errChan <- err
		}
		jsonTemp, _ := json.Marshal(submittingAcc)
		json.Unmarshal(jsonTemp, &memberAccess.SubmittingAccount)

		submittingPerson, err := invitationPyload.personDataSource.GetMongoDataSource().FindByID(
			memberAccess.SubmittingAccount.Person.ID,
			operationOptions,
		)
		if err != nil {
			errChan <- err
		}

		jsonTemp, _ = json.Marshal(submittingPerson)
		json.Unmarshal(jsonTemp, &memberAccess.SubmittingAccount.Person)

		submittingAccountLoadedChan <- true
	}()

	go func() {
		invitedAcc, err := invitationPyload.accountDataSource.GetMongoDataSource().FindByID(
			memberAccess.Account.ID,
			operationOptions,
		)
		if err != nil {
			errChan <- err
		}
		jsonTemp, _ := json.Marshal(invitedAcc)
		json.Unmarshal(jsonTemp, &memberAccess.Account)

		invitedPerson, err := invitationPyload.personDataSource.GetMongoDataSource().FindByID(
			memberAccess.Account.Person.ID,
			operationOptions,
		)
		if err != nil {
			errChan <- err
		}

		jsonTemp, _ = json.Marshal(invitedPerson)
		json.Unmarshal(jsonTemp, &memberAccess.Account.Person)

		invitedAccountLoadedChan <- true
	}()

	for i := 0; i < 2; {
		select {
		case err := <-errChan:
			return false, err
		case _ = <-submittingAccountLoadedChan:
			i++
		case _ = <-invitedAccountLoadedChan:
			i++
		}
	}

	return true, nil
}
