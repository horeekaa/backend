package memberaccessdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type invitationPayloadLoader struct {
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource
	personDataSource  databaseaccountdatasourceinterfaces.PersonDataSource
}

func NewInvitationPayloadLoader(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
) (memberaccessdomainrepositoryutilityinterfaces.InvitationPayloadLoader, error) {
	return &invitationPayloadLoader{
		accountDataSource: accountDataSource,
		personDataSource:  personDataSource,
	}, nil
}

func (invitationPayload *invitationPayloadLoader) Execute(
	notification *model.InternalCreateNotification,
) (bool, error) {
	memberAccess := notification.PayloadOptions.MemberAccessInvitationPayload.MemberAccess

	submittingAccountLoadedChan := make(chan bool)
	invitedAccountLoadedChan := make(chan bool)
	errChan := make(chan error)

	go func() {
		submittingAcc, err := invitationPayload.accountDataSource.GetMongoDataSource().FindByID(
			memberAccess.SubmittingAccount.ID,
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			errChan <- err
		}
		jsonTemp, _ := json.Marshal(submittingAcc)
		json.Unmarshal(jsonTemp, &memberAccess.SubmittingAccount)

		submittingPerson, err := invitationPayload.personDataSource.GetMongoDataSource().FindByID(
			memberAccess.SubmittingAccount.Person.ID,
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			errChan <- err
		}

		jsonTemp, _ = json.Marshal(submittingPerson)
		json.Unmarshal(jsonTemp, &memberAccess.SubmittingAccount.Person)

		submittingAccountLoadedChan <- true
	}()

	go func() {
		invitedAcc, err := invitationPayload.accountDataSource.GetMongoDataSource().FindByID(
			memberAccess.Account.ID,
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			errChan <- err
		}
		jsonTemp, _ := json.Marshal(invitedAcc)
		json.Unmarshal(jsonTemp, &memberAccess.Account)

		invitedPerson, err := invitationPayload.personDataSource.GetMongoDataSource().FindByID(
			memberAccess.Account.Person.ID,
			&mongodbcoretypes.OperationOptions{},
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
