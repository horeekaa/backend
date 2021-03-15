package accountservicescoordinators

import (
	"context"
	"strings"

	servicerepodependencies "github.com/horeekaa/backend/dependencies/services/repos"
	"github.com/horeekaa/backend/model"
	firebaseauthclients "github.com/horeekaa/backend/repositories/authentication/firebase"
	mongomarshaler "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/modelMarshaler"
	authserviceclients "github.com/horeekaa/backend/services/authentication"
	authserviceclientinterfaces "github.com/horeekaa/backend/services/authentication/interfaces"
	authserviceoperations "github.com/horeekaa/backend/services/authentication/operations"
	accountservicecoordinatorinterfaces "github.com/horeekaa/backend/services/coordinators/interfaces/accounts"
	databaseservicerepointerfaces "github.com/horeekaa/backend/services/database/interfaces/repos"
	databaseservicetransactioninterfaces "github.com/horeekaa/backend/services/database/interfaces/transaction"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
	databaseservicetransactions "github.com/horeekaa/backend/services/database/transactions"
)

type manageAccountAuthenticationServiceComponent struct {
	personService                               databaseservicerepointerfaces.PersonService
	accountService                              databaseservicerepointerfaces.AccountService
	authenticationService                       authserviceclientinterfaces.AuthenticationService
	manageAccountAuthenticationUsecaseComponent accountservicecoordinatorinterfaces.ManageAccountAuthenticationUsecaseComponent
}

type manageAccountAuthenticationService struct {
	databaseServiceTransaction databaseservicetransactioninterfaces.DatabaseServiceTransaction
}

func NewManageAccountAuthenticationService(manageAccountAuthenticationUsecaseComponent accountservicecoordinatorinterfaces.ManageAccountAuthenticationUsecaseComponent, context context.Context) (accountservicecoordinatorinterfaces.ManageAccountAuthenticationService, error) {
	personService, _ := servicerepodependencies.InitializePersonService()
	accountService, _ := servicerepodependencies.InitializeAccountService()

	firebaseRepo, _ := firebaseauthclients.NewFirebaseAuthentication(&context)
	authenticationService, _ := authserviceclients.NewAuthenticationService(&firebaseRepo)

	defaultTitle := "ManageAccountAuthentication"

	databaseServiceTransaction, _ := databaseservicetransactions.NewDatabaseServiceTransaction(
		&manageAccountAuthenticationServiceComponent{
			personService:         personService,
			accountService:        accountService,
			authenticationService: authenticationService,
			manageAccountAuthenticationUsecaseComponent: manageAccountAuthenticationUsecaseComponent,
		},
		&defaultTitle,
	)

	return &manageAccountAuthenticationService{
		databaseServiceTransaction: databaseServiceTransaction,
	}, nil
}

func (msgAccAuthServ *manageAccountAuthenticationService) RunTransaction(authHeader string) (*model.Account, error) {
	output, err := msgAccAuthServ.databaseServiceTransaction.RunTransaction(authHeader)
	return output.(*model.Account), err
}

func (msgAccAuthServComponent *manageAccountAuthenticationServiceComponent) PreTransaction(authHeader interface{}) (interface{}, error) {
	return msgAccAuthServComponent.manageAccountAuthenticationUsecaseComponent.Validation(authHeader.(string))
}

func (mgsAccAuthServComponent *manageAccountAuthenticationServiceComponent) TransactionBody(session *databaseserviceoperations.ServiceOptions, authHeader interface{}) (interface{}, error) {
	splitted := strings.Split(authHeader.(string), " ")
	authToken := splitted[len(splitted)-1]

	userChannel, errChannel := mgsAccAuthServComponent.authenticationService.VerifyTokenAndGetUser(authToken)
	select {
	case user := <-userChannel:
		storedAccountID := user.ServiceUser.RepoUser.CustomClaims[authserviceoperations.CustomClaimsAccountIDKey]
		if &storedAccountID != nil {
			storedAccountID = (storedAccountID).(string)
			unmarshaledAccountID, _ := mongomarshaler.UnmarshalObjectID(storedAccountID)

			accountChannel, errChannel := mgsAccAuthServComponent.accountService.FindByID(
				unmarshaledAccountID,
				session,
			)
			account := &model.Account{}
			select {
			case account = <-accountChannel:
				return account, nil
			case err := <-errChannel:
				return nil, err
			}
		}
		accountChannel, errChannel := mgsAccAuthServComponent.accountService.FindOne(
			map[string]interface{}{
				"email": user.ServiceUser.RepoUser.Email,
			},
			session,
		)
		select {
		case account := <-accountChannel:
			mgsAccAuthServComponent.authenticationService.SetRoleInAuthUserData(
				user.ServiceUser.RepoUser.UID,
				model.AccountTypePerson.String(),
				account.ID.String(),
			)
			return account, nil
		case _ = <-errChannel:
			fullNameSplit := strings.Split(user.ServiceUser.RepoUser.DisplayName, " ")
			firstName := fullNameSplit[0]
			lastName := fullNameSplit[len(fullNameSplit)-1]
			if firstName == lastName {
				lastName = ""
			}
			defaultNoOfRecentTransaction := 15

			personChannel, errChannel := mgsAccAuthServComponent.personService.Create(
				&model.CreatePerson{
					FirstName:                   firstName,
					LastName:                    lastName,
					PhoneNumber:                 user.ServiceUser.RepoUser.PhoneNumber,
					Email:                       user.ServiceUser.RepoUser.Email,
					DeviceTokens:                []*string{},
					NoOfRecentTransactionToKeep: &defaultNoOfRecentTransaction,
				},
				session,
			)
			select {
			case person := <-personChannel:
				accountChannel, errChannel := mgsAccAuthServComponent.accountService.Create(
					&model.CreateAccount{
						Type: model.AccountTypePerson,
						Person: &model.ObjectIDOnly{
							ID: &person.ID,
						},
					},
					session,
				)
				select {
				case account := <-accountChannel:
					return account, nil
				case err := <-errChannel:
					return nil, err
				}
			case err := <-errChannel:
				return nil, err
			}
		}
	case err := <-errChannel:
		return nil, err
	}
}