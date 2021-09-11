package moudomainrepositoryutilities

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	moudomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories/utils"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type partyLoader struct {
	accountDataSource      databaseaccountdatasourceinterfaces.AccountDataSource
	personDataSource       databaseaccountdatasourceinterfaces.PersonDataSource
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource
}

func NewPartyLoader(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
) (moudomainrepositoryutilityinterfaces.PartyLoader, error) {
	return &partyLoader{
		accountDataSource,
		personDataSource,
		organizationDataSource,
	}, nil
}

func (partyLoader *partyLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.PartyInput,
	output *model.InternalPartyInput,
) (bool, error) {
	if input != nil {
		currentTime := time.Now()
		output.CreatedAt = &currentTime
		output.UpdatedAt = &currentTime
		organizationLoadedChan := make(chan bool)
		accountLoadedChan := make(chan bool)
		signatureLoadedChan := make(chan bool)
		errChan := make(chan error)
		go func() {
			if input.Organization == nil {
				organizationLoadedChan <- true
				return
			}
			organization, err := partyLoader.organizationDataSource.GetMongoDataSource().FindByID(
				*input.Organization.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}

			jsonTemp, _ := json.Marshal(organization)
			json.Unmarshal(jsonTemp, &output.Organization)
			organizationLoadedChan <- true
		}()

		go func() {
			if input.AccountInCharge == nil {
				accountLoadedChan <- true
				return
			}
			account, err := partyLoader.accountDataSource.GetMongoDataSource().FindByID(
				*input.AccountInCharge.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}

			jsonTemp, _ := json.Marshal(account)
			json.Unmarshal(jsonTemp, &output.AccountInCharge)

			person, err := partyLoader.personDataSource.GetMongoDataSource().FindByID(
				output.AccountInCharge.Person.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}

			jsonTemp, _ = json.Marshal(person)
			json.Unmarshal(jsonTemp, &output.AccountInCharge.Person)
			accountLoadedChan <- true
		}()

		go func() {
			if !funk.GetOrElse(
				funk.Get(input, "Signature.ConfirmSign"),
				false,
			).(bool) || input.AccountInCharge == nil {
				signatureLoadedChan <- true
				return
			}
			output.Signature = &model.InternalSignatureInput{
				ID:        partyLoader.organizationDataSource.GetMongoDataSource().GenerateObjectID(),
				CreatedAt: time.Now(),
			}

			account, err := partyLoader.accountDataSource.GetMongoDataSource().FindByID(
				*input.AccountInCharge.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}

			jsonTemp, _ := json.Marshal(account)
			json.Unmarshal(jsonTemp, &output.Signature.Account)

			person, err := partyLoader.personDataSource.GetMongoDataSource().FindByID(
				output.Signature.Account.Person.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}

			jsonTemp, _ = json.Marshal(person)
			json.Unmarshal(jsonTemp, &output.Signature.Account.Person)
			signatureLoadedChan <- true
		}()

		for i := 0; i < 3; {
			select {
			case err := <-errChan:
				return false, err
			case _ = <-accountLoadedChan:
				i++
			case _ = <-organizationLoadedChan:
				i++
			case _ = <-signatureLoadedChan:
				i++
			}
		}
	}
	return true, nil
}
