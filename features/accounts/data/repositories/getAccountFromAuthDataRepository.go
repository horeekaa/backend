package accountdomainrepositories

import (
	"firebase.google.com/go/v4/auth"
	firebaseauthcoretypes "github.com/horeekaa/backend/core/authentication/firebase/types"
	authenticationcoremodels "github.com/horeekaa/backend/core/authentication/models"
	mongomarshaler "github.com/horeekaa/backend/core/databaseClient/mongodb/modelMarshalers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	firebaseauthdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/authentication/interfaces"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
)

type getAccountFromAuthDataRepository struct {
	accountDataSource      databaseaccountdatasourceinterfaces.AccountDataSource
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	firebaseDataSource     firebaseauthdatasourceinterfaces.FirebaseAuthRepo
	pathIdentity           string
}

func NewGetAccountFromAuthDataRepository(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	firebaseDataSource firebaseauthdatasourceinterfaces.FirebaseAuthRepo,
) (accountdomainrepositoryinterfaces.GetAccountFromAuthData, error) {
	return &getAccountFromAuthDataRepository{
		accountDataSource,
		memberAccessDataSource,
		firebaseDataSource,
		"GetAccountFromAuthDataRepository",
	}, nil
}

func (getAccFromAuthDataRepo *getAccountFromAuthDataRepository) Execute(
	input accountdomainrepositorytypes.GetAccountFromAuthDataInput,
) (*model.Account, error) {
	user := input.Context.Value(authenticationcoremodels.UserContextKey)
	if user == nil {
		return nil, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.AuthenticationTokenFailed,
			getAccFromAuthDataRepo.pathIdentity,
			nil,
		)
	}

	storedAccountID := user.(*auth.UserRecord).CustomClaims[firebaseauthcoretypes.FirebaseCustomClaimsAccountIDKey]
	if storedAccountID != nil {
		storedAccountID = (storedAccountID).(string)
		unmarshaledAccountID, _ := mongomarshaler.UnmarshalObjectID(storedAccountID)

		account, err := getAccFromAuthDataRepo.accountDataSource.GetMongoDataSource().FindByID(
			unmarshaledAccountID,
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				getAccFromAuthDataRepo.pathIdentity,
				err,
			)
		}
		return account, nil
	}

	account, err := getAccFromAuthDataRepo.accountDataSource.GetMongoDataSource().FindOne(
		map[string]interface{}{
			"email": user.(*auth.UserRecord).Email,
		},
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAccFromAuthDataRepo.pathIdentity,
			err,
		)
	}

	if account != nil {
		_, err = getAccFromAuthDataRepo.firebaseDataSource.SetRoleInAuthUserData(
			input.Context,
			user.(*auth.UserRecord).UID,
			model.AccountTypePerson.String(),
			account.ID.Hex(),
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				getAccFromAuthDataRepo.pathIdentity,
				err,
			)
		}

		_, err = getAccFromAuthDataRepo.memberAccessDataSource.GetMongoDataSource().Update(
			map[string]interface{}{
				"account._id":         account.ID,
				"memberAccessRefType": model.MemberAccessRefTypeOrganizationsBased,
			},
			&model.DatabaseUpdateMemberAccess{
				Status: func(s model.MemberAccessStatus) *model.MemberAccessStatus {
					return &s
				}(model.MemberAccessStatusInactive),
			},
			&mongodbcoretypes.OperationOptions{},
		)
	}
	return account, nil
}
