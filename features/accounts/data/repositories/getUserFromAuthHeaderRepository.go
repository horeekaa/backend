package accountdomainrepositories

import (
	"strings"

	"firebase.google.com/go/v4/auth"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	firebaseauthdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/authentication/interfaces"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
)

type getUserFromAuthHeaderRepository struct {
	firebaseDataSource firebaseauthdatasourceinterfaces.FirebaseAuthRepo
}

func NewGetUserFromAuthHeaderRepository(
	firebaseDataSource firebaseauthdatasourceinterfaces.FirebaseAuthRepo,
) (accountdomainrepositoryinterfaces.GetUserFromAuthHeaderRepository, error) {
	return &getUserFromAuthHeaderRepository{
		firebaseDataSource,
	}, nil
}

func (getUsrFromAuthHeaderRepo *getUserFromAuthHeaderRepository) Execute(
	getUserFromAuthHeaderInput accountdomainrepositorytypes.GetUserFromAuthHeaderInput,
) (*auth.UserRecord, error) {
	splitted := strings.Split(getUserFromAuthHeaderInput.AuthHeader, " ")
	authToken := splitted[len(splitted)-1]

	token, err := getUsrFromAuthHeaderRepo.firebaseDataSource.VerifyAndDecodeToken(
		getUserFromAuthHeaderInput.Context,
		authToken,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getUserFromAuthHeader",
			err,
		)
	}
	user, err := getUsrFromAuthHeaderRepo.firebaseDataSource.GetAuthUserDataById(
		getUserFromAuthHeaderInput.Context,
		token.UID,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getUserFromAuthHeader",
			err,
		)
	}

	return user, nil
}
