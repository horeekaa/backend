package accountdomainrepositories

import (
	"strings"

	"firebase.google.com/go/v4/auth"
	authenticationcoremodels "github.com/horeekaa/backend/core/authentication/models"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type createAccountFromAuthDataTransactionComponent struct {
	personDataSource  databaseaccountdatasourceinterfaces.PersonDataSource
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource
	pathIdentity      string
}

func NewCreateAccountFromAuthDataTransactionComponent(
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
) (accountdomainrepositoryinterfaces.CreateAccountFromAuthDataTransactionComponent, error) {
	return &createAccountFromAuthDataTransactionComponent{
		personDataSource:  personDataSource,
		accountDataSource: accountDataSource,
		pathIdentity:      "CreateAccountFromAuthDataComponent",
	}, nil
}

func (createAccFromAuthDataCom *createAccountFromAuthDataTransactionComponent) PreTransaction(
	crtAccFromAuthDataInput accountdomainrepositorytypes.CreateAccountFromAuthDataInput,
) (accountdomainrepositorytypes.CreateAccountFromAuthDataInput, error) {
	return crtAccFromAuthDataInput, nil
}

func (createAccFromAuthDataCom *createAccountFromAuthDataTransactionComponent) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	createAccFrmAuthDataInput accountdomainrepositorytypes.CreateAccountFromAuthDataInput,
) (*model.Account, error) {
	user := createAccFrmAuthDataInput.Context.Value(authenticationcoremodels.UserContextKey)
	if user == nil {
		return nil, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.AuthenticationTokenFailed,
			createAccFromAuthDataCom.pathIdentity,
			nil,
		)
	}

	fullNameSplit := strings.Split(user.(*auth.UserRecord).DisplayName, " ")
	firstName := fullNameSplit[0]
	lastName := fullNameSplit[len(fullNameSplit)-1]
	if firstName == lastName {
		lastName = ""
	}
	if firstName == "" {
		firstName = strings.Split(user.(*auth.UserRecord).Email, "@")[0]
	}

	person, err := createAccFromAuthDataCom.personDataSource.GetMongoDataSource().Create(
		&model.DatabaseCreatePerson{
			FirstName:   firstName,
			LastName:    lastName,
			PhoneNumber: user.(*auth.UserRecord).PhoneNumber,
		},
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createAccFromAuthDataCom.pathIdentity,
			err,
		)
	}

	account, err := createAccFromAuthDataCom.accountDataSource.GetMongoDataSource().Create(
		&model.DatabaseCreateAccount{
			Email: user.(*auth.UserRecord).Email,
			Type:  model.AccountTypePerson,
			Person: &model.ObjectIDOnly{
				ID: &person.ID,
			},
		},
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createAccFromAuthDataCom.pathIdentity,
			err,
		)
	}
	return account, nil
}
