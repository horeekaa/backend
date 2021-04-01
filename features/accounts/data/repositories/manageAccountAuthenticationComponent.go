package accountdomainrepositories

import (
	"strings"

	firebaseauthcoretypes "github.com/horeekaa/backend/core/authentication/firebase/types"
	mongomarshaler "github.com/horeekaa/backend/core/databaseClient/mongoDB/modelMarshalers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	firebaseauthdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/authentication/interfaces"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type manageAccountAuthenticationTransactionComponent struct {
	personDataSource                            databaseaccountdatasourceinterfaces.PersonDataSource
	accountDataSource                           databaseaccountdatasourceinterfaces.AccountDataSource
	firebaseDataSource                          firebaseauthdatasourceinterfaces.FirebaseAuthRepo
	manageAccountAuthenticationUsecaseComponent accountdomainrepositoryinterfaces.ManageAccountAuthenticationUsecaseComponent
}

func NewManageAccountAuthenticationTransactionComponent(
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	firebaseDataSource firebaseauthdatasourceinterfaces.FirebaseAuthRepo,
) (accountdomainrepositoryinterfaces.ManageAccountAuthenticationTransactionComponent, error) {
	return &manageAccountAuthenticationTransactionComponent{
		personDataSource:   personDataSource,
		accountDataSource:  accountDataSource,
		firebaseDataSource: firebaseDataSource,
	}, nil
}

func (msgAccAuthTrx *manageAccountAuthenticationTransactionComponent) SetValidation(
	usecaseComponent accountdomainrepositoryinterfaces.ManageAccountAuthenticationUsecaseComponent,
) (bool, error) {
	msgAccAuthTrx.manageAccountAuthenticationUsecaseComponent = usecaseComponent
	return true, nil
}

func (msgAccAuthTrx *manageAccountAuthenticationTransactionComponent) PreTransaction(
	manageAccountAuthInput accountdomainrepositorytypes.ManageAccountAuthenticationInput,
) (accountdomainrepositorytypes.ManageAccountAuthenticationInput, error) {
	if msgAccAuthTrx.manageAccountAuthenticationUsecaseComponent == nil {
		return manageAccountAuthInput, nil
	}
	return msgAccAuthTrx.manageAccountAuthenticationUsecaseComponent.Validation(manageAccountAuthInput)
}

func (msgAccAuthTrx *manageAccountAuthenticationTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	manageAccountAuthInput accountdomainrepositorytypes.ManageAccountAuthenticationInput,
) (*model.Account, error) {
	splitted := strings.Split(manageAccountAuthInput.AuthHeader, " ")
	authToken := splitted[len(splitted)-1]

	token, err := msgAccAuthTrx.firebaseDataSource.VerifyAndDecodeToken(
		manageAccountAuthInput.Context,
		authToken,
	)
	if err != nil {
		return nil, err
	}
	user, err := msgAccAuthTrx.firebaseDataSource.GetAuthUserDataById(
		manageAccountAuthInput.Context,
		token.UID,
	)

	storedAccountID := user.CustomClaims[firebaseauthcoretypes.FirebaseCustomClaimsAccountIDKey]
	if &storedAccountID != nil {
		storedAccountID = (storedAccountID).(string)
		unmarshaledAccountID, _ := mongomarshaler.UnmarshalObjectID(storedAccountID)

		account, err := msgAccAuthTrx.accountDataSource.GetMongoDataSource().FindByID(
			unmarshaledAccountID,
			session,
		)
		if err != nil {
			return nil, err
		}
		return account, nil
	}

	account, err := msgAccAuthTrx.accountDataSource.GetMongoDataSource().FindOne(
		map[string]interface{}{
			"email": user.Email,
		},
		session,
	)
	if account == nil {
		fullNameSplit := strings.Split(user.DisplayName, " ")
		firstName := fullNameSplit[0]
		lastName := fullNameSplit[len(fullNameSplit)-1]
		if firstName == lastName {
			lastName = ""
		}
		defaultNoOfRecentTransaction := 15

		person, err := msgAccAuthTrx.personDataSource.GetMongoDataSource().Create(
			&model.CreatePerson{
				FirstName:                   firstName,
				LastName:                    lastName,
				PhoneNumber:                 user.PhoneNumber,
				Email:                       user.Email,
				NoOfRecentTransactionToKeep: &defaultNoOfRecentTransaction,
			},
			session,
		)
		if err != nil {
			return nil, err
		}

		account, err := msgAccAuthTrx.accountDataSource.GetMongoDataSource().Create(
			&model.CreateAccount{
				Type: model.AccountTypePerson,
				Person: &model.ObjectIDOnly{
					ID: &person.ID,
				},
			},
			session,
		)
		if err != nil {
			return nil, err
		}
		return account, nil
	}

	msgAccAuthTrx.firebaseDataSource.SetRoleInAuthUserData(
		manageAccountAuthInput.Context,
		user.UID,
		model.AccountTypePerson.String(),
		account.ID.String(),
	)
	return account, nil
}
