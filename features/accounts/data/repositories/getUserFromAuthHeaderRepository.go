package accountdomainrepositories

import (
	"strings"

	authenticationcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/interfaces"
	authenticationcoremodels "github.com/horeekaa/backend/core/authentication/models"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
)

type getUserFromAuthHeaderRepository struct {
	firebaseDataSource authenticationcoreclientinterfaces.AuthenticationRepo
}

func NewGetUserFromAuthHeaderRepository(
	firebaseDataSource authenticationcoreclientinterfaces.AuthenticationRepo,
) (accountdomainrepositoryinterfaces.GetUserFromAuthHeaderRepository, error) {
	return &getUserFromAuthHeaderRepository{
		firebaseDataSource,
	}, nil
}

func (getUsrFromAuthHeaderRepo *getUserFromAuthHeaderRepository) Execute(
	getUserFromAuthHeaderInput accountdomainrepositorytypes.GetUserFromAuthHeaderInput,
) (*authenticationcoremodels.AuthUserWrap, error) {
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
		token.FirebaseToken.UID,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getUserFromAuthHeader",
			err,
		)
	}

	return user, nil
}
