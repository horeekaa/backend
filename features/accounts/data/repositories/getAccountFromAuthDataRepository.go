package accountdomainrepositories

import (
	firebaseauthcoretypes "github.com/horeekaa/backend/core/authentication/firebase/types"
	authenticationcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/interfaces"
	mongomarshaler "github.com/horeekaa/backend/core/databaseClient/mongodb/modelMarshalers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAccountFromAuthDataRepository struct {
	accountDataSource  databaseaccountdatasourceinterfaces.AccountDataSource
	firebaseDataSource authenticationcoreclientinterfaces.AuthenticationRepo
}

func NewGetAccountFromAuthDataRepository(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	firebaseDataSource authenticationcoreclientinterfaces.AuthenticationRepo,
) (accountdomainrepositoryinterfaces.GetAccountFromAuthData, error) {
	return &getAccountFromAuthDataRepository{
		accountDataSource,
		firebaseDataSource,
	}, nil
}

func (getAccFromAuthDataRepo *getAccountFromAuthDataRepository) Execute(
	input accountdomainrepositorytypes.GetAccountFromAuthDataInput,
) (*model.Account, error) {
	storedAccountID := input.User.FirebaseUser.CustomClaims[firebaseauthcoretypes.FirebaseCustomClaimsAccountIDKey]
	if &storedAccountID != nil {
		storedAccountID = (storedAccountID).(string)
		unmarshaledAccountID, _ := mongomarshaler.UnmarshalObjectID(storedAccountID)

		account, err := getAccFromAuthDataRepo.accountDataSource.GetMongoDataSource().FindByID(
			unmarshaledAccountID,
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/getAccountFromAuthDataRepository",
				err,
			)
		}
		return account, nil
	}

	account, err := getAccFromAuthDataRepo.accountDataSource.GetMongoDataSource().FindOne(
		map[string]interface{}{
			"email": input.User.FirebaseUser.Email,
		},
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAccountFromAuthDataRepository",
			err,
		)
	}

	if account != nil {
		_, err = getAccFromAuthDataRepo.firebaseDataSource.SetRoleInAuthUserData(
			input.Context,
			input.User.FirebaseUser.UID,
			model.AccountTypePerson.String(),
			account.ID.String(),
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/getAccountFromAuthDataRepository",
				err,
			)
		}
	}
	return account, nil
}
